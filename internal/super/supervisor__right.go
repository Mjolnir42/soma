/*-
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"database/sql"
	"fmt"

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

func (s *Supervisor) right(q *msg.Request) {
	result := msg.FromRequest(q)

	s.requestLog(q)

	if q.Grant.RecipientType != `user` {
		result.NotImplemented(fmt.Errorf("Rights for recipient type"+
			" %s are currently not implemented",
			q.Grant.RecipientType))
		goto abort
	}

	switch q.Action {
	case msg.ActionGrant, msg.ActionRevoke:
		if s.readonly {
			result.Conflict(fmt.Errorf(`Readonly instance`))
			goto abort
		}
		s.rightWrite(q)
	case msg.ActionSearch:
		go func() { s.rightRead(q) }()
	default:
		result.UnknownRequest(q)
		goto abort
	}
	return

abort:
	q.Reply <- result
}

func (s *Supervisor) rightWrite(q *msg.Request) {
	switch q.Action {
	case msg.ActionGrant:
		switch q.Grant.Category {
		case msg.CategorySystem,
			msg.CategoryGlobal,
			msg.CategoryGrantGlobal,
			msg.CategoryPermission,
			msg.CategoryGrantPermission,
			msg.CategoryOperation,
			msg.CategoryGrantOperation:
			s.rightGrantGlobal(q)
		case msg.CategoryRepository,
			msg.CategoryGrantRepository:
			s.rightGrantRepository(q)
		case msg.CategoryTeam,
			msg.CategoryGrantTeam:
			s.rightGrantTeam(q)
		case msg.CategoryMonitoring,
			msg.CategoryGrantMonitoring:
			s.rightGrantMonitoring(q)
		}
	case msg.ActionRevoke:
		switch q.Grant.Category {
		case msg.CategorySystem,
			msg.CategoryGlobal,
			msg.CategoryGrantGlobal,
			msg.CategoryPermission,
			msg.CategoryGrantPermission,
			msg.CategoryOperation,
			msg.CategoryGrantOperation:
			s.rightRevokeGlobal(q)
		case msg.CategoryRepository,
			msg.CategoryGrantRepository:
			s.rightRevokeRepository(q)
		case msg.CategoryTeam,
			msg.CategoryGrantTeam:
			s.rightRevokeTeam(q)
		case msg.CategoryMonitoring,
			msg.CategoryGrantMonitoring:
			s.rightRevokeMonitoring(q)
		}
	}
}

func (s *Supervisor) rightGrantGlobal(q *msg.Request) {
	result := msg.FromRequest(q)
	var (
		err                             error
		res                             sql.Result
		adminID, userID, toolID, teamID sql.NullString
	)

	if q.Grant.ObjectType != `` || q.Grant.ObjectId != `` {
		result.BadRequest(fmt.Errorf(
			`Invalid global grant specification`))
		goto dispatch
	}

	switch q.Grant.RecipientType {
	case `admin`:
		adminID.String = q.Grant.RecipientId
		adminID.Valid = true
	case `user`:
		userID.String = q.Grant.RecipientId
		userID.Valid = true
	case `tool`:
		toolID.String = q.Grant.RecipientId
		toolID.Valid = true
	case `team`:
		teamID.String = q.Grant.RecipientId
		teamID.Valid = true
	}

	q.Grant.Id = uuid.NewV4().String()
	if res, err = s.stmtGrantAuthorizationGlobal.Exec(
		q.Grant.Id,
		adminID,
		userID,
		toolID,
		teamID,
		q.Grant.PermissionId,
		q.Grant.Category,
		q.AuthUser,
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if result.RowCnt(res.RowsAffected()) {
		result.Grant = []proto.Grant{q.Grant}
	}

dispatch:
	q.Reply <- result
}

func (s *Supervisor) rightGrantRepository(q *msg.Request) {
	result := msg.FromRequest(q)
	var (
		err                       error
		res                       sql.Result
		userID, toolID, teamID    sql.NullString
		repoID, bucketID, groupID sql.NullString
		clusterID, nodeID         sql.NullString
		repoName                  string
	)

	switch q.Grant.ObjectType {
	case msg.EntityRepository:
		repoID.String = q.Grant.ObjectId
		repoID.Valid = true
	case msg.EntityBucket:
		if err = s.conn.QueryRow(
			stmt.RepoByBucketId,
			q.Grant.ObjectId,
		).Scan(
			repoID,
			repoName,
		); err != nil {
			result.ServerError(err)
			goto dispatch
		}

		bucketID.String = q.Grant.ObjectId
		bucketID.Valid = true
	case msg.EntityGroup, msg.EntityCluster, msg.EntityNode:
		result.NotImplemented(fmt.Errorf(
			`Deep repository grants currently not implemented.`))
		goto dispatch
	default:
		result.BadRequest(fmt.Errorf(
			`Invalid repository grant specification`))
		goto dispatch
	}

	switch q.Grant.RecipientType {
	case `user`:
		userID.String = q.Grant.RecipientId
		userID.Valid = true
	case `tool`:
		toolID.String = q.Grant.RecipientId
		toolID.Valid = true
	case `team`:
		teamID.String = q.Grant.RecipientId
		teamID.Valid = true
	}

	q.Grant.Id = uuid.NewV4().String()
	if res, err = s.stmtGrantAuthorizationRepository.Exec(
		q.Grant.Id,
		userID,
		toolID,
		teamID,
		q.Grant.Category,
		q.Grant.PermissionId,
		q.Grant.ObjectType,
		repoID,
		bucketID,
		groupID,
		clusterID,
		nodeID,
		q.AuthUser,
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if result.RowCnt(res.RowsAffected()) {
		result.Grant = []proto.Grant{q.Grant}
	}

dispatch:
	q.Reply <- result
}

func (s *Supervisor) rightGrantTeam(q *msg.Request) {
	result := msg.FromRequest(q)
	var (
		err                    error
		res                    sql.Result
		userID, toolID, teamID sql.NullString
	)

	switch q.Grant.RecipientType {
	case `user`:
		userID.String = q.Grant.RecipientId
		userID.Valid = true
	case `tool`:
		toolID.String = q.Grant.RecipientId
		toolID.Valid = true
	case `team`:
		teamID.String = q.Grant.RecipientId
		teamID.Valid = true
	}

	q.Grant.Id = uuid.NewV4().String()
	if res, err = s.stmtGrantAuthorizationTeam.Exec(
		q.Grant.Id,
		userID,
		toolID,
		teamID,
		q.Grant.Category,
		q.Grant.PermissionId,
		q.Grant.ObjectId,
		q.AuthUser,
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if result.RowCnt(res.RowsAffected()) {
		result.Grant = []proto.Grant{q.Grant}
	}

dispatch:
	q.Reply <- result
}

func (s *Supervisor) rightGrantMonitoring(q *msg.Request) {
	result := msg.FromRequest(q)
	var (
		err                    error
		res                    sql.Result
		userID, toolID, teamID sql.NullString
	)

	switch q.Grant.RecipientType {
	case `user`:
		userID.String = q.Grant.RecipientId
		userID.Valid = true
	case `tool`:
		toolID.String = q.Grant.RecipientId
		toolID.Valid = true
	case `team`:
		teamID.String = q.Grant.RecipientId
		teamID.Valid = true
	}

	q.Grant.Id = uuid.NewV4().String()
	if res, err = s.stmtGrantAuthorizationMonitoring.Exec(
		q.Grant.Id,
		userID,
		toolID,
		teamID,
		q.Grant.Category,
		q.Grant.PermissionId,
		q.Grant.ObjectId,
		q.AuthUser,
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if result.RowCnt(res.RowsAffected()) {
		result.Grant = []proto.Grant{q.Grant}
	}

dispatch:
	q.Reply <- result
}

func (s *Supervisor) rightRevokeGlobal(q *msg.Request) {
	result := msg.FromRequest(q)
	var err error
	var res sql.Result

	if res, err = s.stmtRevokeAuthorizationGlobal.Exec(
		q.Grant.Id,
		q.Grant.PermissionId,
		q.Grant.Category,
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if result.RowCnt(res.RowsAffected()) {
		result.Grant = []proto.Grant{q.Grant}
	}

dispatch:
	q.Reply <- result
}

func (s *Supervisor) rightRevokeRepository(q *msg.Request) {
	result := msg.FromRequest(q)
	var err error
	var res sql.Result

	if res, err = s.stmtRevokeAuthorizationRepository.Exec(
		q.Grant.Id,
		q.Grant.PermissionId,
		q.Grant.Category,
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if result.RowCnt(res.RowsAffected()) {
		result.Grant = []proto.Grant{q.Grant}
	}

dispatch:
	q.Reply <- result
}

func (s *Supervisor) rightRevokeTeam(q *msg.Request) {
	result := msg.FromRequest(q)
	var err error
	var res sql.Result

	if res, err = s.stmtRevokeAuthorizationTeam.Exec(
		q.Grant.Id,
		q.Grant.PermissionId,
		q.Grant.Category,
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if result.RowCnt(res.RowsAffected()) {
		result.Grant = []proto.Grant{q.Grant}
	}

dispatch:
	q.Reply <- result
}

func (s *Supervisor) rightRevokeMonitoring(q *msg.Request) {
	result := msg.FromRequest(q)
	var err error
	var res sql.Result

	if res, err = s.stmtRevokeAuthorizationMonitoring.Exec(
		q.Grant.Id,
		q.Grant.PermissionId,
		q.Grant.Category,
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if result.RowCnt(res.RowsAffected()) {
		result.Grant = []proto.Grant{q.Grant}
	}

dispatch:
	q.Reply <- result
}

func (s *Supervisor) rightRead(q *msg.Request) {
	switch q.Action {
	case msg.ActionSearch:
		switch q.Grant.Category {
		case msg.CategorySystem,
			msg.CategoryGlobal,
			msg.CategoryGrantGlobal,
			msg.CategoryPermission,
			msg.CategoryGrantPermission,
			msg.CategoryOperation,
			msg.CategoryGrantOperation:
			s.rightSearchGlobal(q)
		case msg.CategoryRepository,
			msg.CategoryGrantRepository,
			msg.CategoryTeam,
			msg.CategoryGrantTeam,
			msg.CategoryMonitoring,
			msg.CategoryGrantMonitoring:
			s.rightSearchScoped(q)
		}
	}
}

func (s *Supervisor) rightSearchGlobal(q *msg.Request) {
	result := msg.FromRequest(q)
	var (
		err     error
		grantID string
	)
	if err = s.stmtSearchAuthorizationGlobal.QueryRow(
		q.Grant.PermissionId,
		q.Grant.Category,
		q.Grant.RecipientId,
		q.Grant.RecipientType,
	).Scan(&grantID); err == sql.ErrNoRows {
		result.NotFound(err)
		goto dispatch
	} else if err != nil {
		result.ServerError(err)
		goto dispatch
	}
	result.Grant = []proto.Grant{proto.Grant{
		Id:            grantID,
		PermissionId:  q.Grant.PermissionId,
		Category:      q.Grant.Category,
		RecipientId:   q.Grant.RecipientId,
		RecipientType: q.Grant.RecipientType,
	}}

dispatch:
	q.Reply <- result
}

func (s *Supervisor) rightSearchScoped(q *msg.Request) {
	result := msg.FromRequest(q)
	var (
		err     error
		grantID string
		scope   *sql.Stmt
	)
	switch q.Grant.Category {
	case msg.CategoryRepository,
		msg.CategoryGrantRepository:
		scope = s.stmtSearchAuthorizationRepository
	case msg.CategoryTeam,
		msg.CategoryGrantTeam:
		scope = s.stmtSearchAuthorizationTeam
	case msg.CategoryMonitoring,
		msg.CategoryGrantMonitoring:
		scope = s.stmtSearchAuthorizationMonitoring
	}
	if err = scope.QueryRow(
		q.Grant.PermissionId,
		q.Grant.Category,
		q.Grant.RecipientId,
		q.Grant.RecipientType,
		q.Grant.ObjectType,
		q.Grant.ObjectId,
	).Scan(&grantID); err == sql.ErrNoRows {
		result.NotFound(err)
		goto dispatch
	} else if err != nil {
		result.ServerError(err)
		goto dispatch
	}
	result.Grant = []proto.Grant{proto.Grant{
		Id:            grantID,
		PermissionId:  q.Grant.PermissionId,
		Category:      q.Grant.Category,
		RecipientId:   q.Grant.RecipientId,
		RecipientType: q.Grant.RecipientType,
		ObjectType:    q.Grant.ObjectType,
		ObjectId:      q.Grant.ObjectId,
	}}

dispatch:
	q.Reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
