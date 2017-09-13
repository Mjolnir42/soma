/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
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
	uuid "github.com/satori/go.uuid"
)

func (s *Supervisor) rightWrite(q *msg.Request, mr *msg.Result) {
	// admin accounts can only receive system permissions
	if q.Grant.RecipientType == msg.SubjectAdmin && q.Grant.Category != msg.CategorySystem {
		mr.BadRequest(fmt.Errorf(
			"Admin accounts can not receive grants"+
				" in category %s", q.Grant.Category))
		return
	}

	switch q.Action {
	case msg.ActionGrant:
		switch q.Grant.Category {
		case msg.CategorySystem,
			msg.CategoryGlobal, msg.CategoryGrantGlobal,
			msg.CategoryPermission, msg.CategoryGrantPermission,
			msg.CategoryOperation, msg.CategoryGrantOperation:
			s.rightGrantGlobal(q, mr)
		case msg.CategoryRepository, msg.CategoryGrantRepository:
			s.rightGrantRepository(q, mr)
		case msg.CategoryTeam, msg.CategoryGrantTeam:
			s.rightGrantTeam(q, mr)
		case msg.CategoryMonitoring, msg.CategoryGrantMonitoring:
			s.rightGrantMonitoring(q, mr)
		}
	case msg.ActionRevoke:
		switch q.Grant.Category {
		case msg.CategorySystem,
			msg.CategoryGlobal, msg.CategoryGrantGlobal,
			msg.CategoryPermission, msg.CategoryGrantPermission,
			msg.CategoryOperation, msg.CategoryGrantOperation:
			s.rightRevokeGlobal(q, mr)
		case msg.CategoryRepository, msg.CategoryGrantRepository:
			s.rightRevokeRepository(q, mr)
		case msg.CategoryTeam, msg.CategoryGrantTeam:
			s.rightRevokeTeam(q, mr)
		case msg.CategoryMonitoring, msg.CategoryGrantMonitoring:
			s.rightRevokeMonitoring(q, mr)
		}
	}
}

func (s *Supervisor) rightGrantGlobal(q *msg.Request, mr *msg.Result) {
	var (
		err                             error
		res                             sql.Result
		adminID, userID, toolID, teamID sql.NullString
	)

	if q.Grant.ObjectType != `` || q.Grant.ObjectId != `` {
		mr.BadRequest(fmt.Errorf(
			`Invalid global grant specification`))
		return
	}

	switch q.Grant.RecipientType {
	case msg.SubjectAdmin:
		adminID.String = q.Grant.RecipientId
		adminID.Valid = true
	case msg.SubjectUser:
		userID.String = q.Grant.RecipientId
		userID.Valid = true
	case msg.SubjectTool:
		toolID.String = q.Grant.RecipientId
		toolID.Valid = true
	case msg.SubjectTeam:
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
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Grant = append(mr.Grant, q.Grant)
	}
}

func (s *Supervisor) rightGrantRepository(q *msg.Request, mr *msg.Result) {
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
		); err == sql.ErrNoRows {
			mr.NotFound(err, q.Section)
			return
		} else if err != nil {
			mr.ServerError(err, q.Section)
			return
		}

		bucketID.String = q.Grant.ObjectId
		bucketID.Valid = true
	case msg.EntityGroup, msg.EntityCluster, msg.EntityNode:
		mr.NotImplemented(fmt.Errorf(
			`Deep repository grants currently not implemented.`))
		return
	default:
		mr.BadRequest(fmt.Errorf(
			`Invalid repository grant specification`))
		return
	}

	switch q.Grant.RecipientType {
	case msg.SubjectUser:
		userID.String = q.Grant.RecipientId
		userID.Valid = true
	case msg.SubjectTool:
		toolID.String = q.Grant.RecipientId
		toolID.Valid = true
	case msg.SubjectTeam:
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
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Grant = append(mr.Grant, q.Grant)
	}
}

func (s *Supervisor) rightGrantTeam(q *msg.Request, mr *msg.Result) {
	var (
		err                    error
		res                    sql.Result
		userID, toolID, teamID sql.NullString
	)

	switch q.Grant.RecipientType {
	case msg.SubjectUser:
		userID.String = q.Grant.RecipientId
		userID.Valid = true
	case msg.SubjectTool:
		toolID.String = q.Grant.RecipientId
		toolID.Valid = true
	case msg.SubjectTeam:
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
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Grant = append(mr.Grant, q.Grant)
	}
}

func (s *Supervisor) rightGrantMonitoring(q *msg.Request, mr *msg.Result) {
	var (
		err                    error
		res                    sql.Result
		userID, toolID, teamID sql.NullString
	)

	switch q.Grant.RecipientType {
	case msg.SubjectUser:
		userID.String = q.Grant.RecipientId
		userID.Valid = true
	case msg.SubjectTool:
		toolID.String = q.Grant.RecipientId
		toolID.Valid = true
	case msg.SubjectTeam:
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
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Grant = append(mr.Grant, q.Grant)
	}
}

func (s *Supervisor) rightRevokeGlobal(q *msg.Request, mr *msg.Result) {
	var err error
	var res sql.Result

	if res, err = s.stmtRevokeAuthorizationGlobal.Exec(
		q.Grant.Id,
		q.Grant.PermissionId,
		q.Grant.Category,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Grant = append(mr.Grant, q.Grant)
	}
}

func (s *Supervisor) rightRevokeRepository(q *msg.Request, mr *msg.Result) {
	var err error
	var res sql.Result

	if res, err = s.stmtRevokeAuthorizationRepository.Exec(
		q.Grant.Id,
		q.Grant.PermissionId,
		q.Grant.Category,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Grant = append(mr.Grant, q.Grant)
	}
}

func (s *Supervisor) rightRevokeTeam(q *msg.Request, mr *msg.Result) {
	var err error
	var res sql.Result

	if res, err = s.stmtRevokeAuthorizationTeam.Exec(
		q.Grant.Id,
		q.Grant.PermissionId,
		q.Grant.Category,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Grant = append(mr.Grant, q.Grant)
	}
}

func (s *Supervisor) rightRevokeMonitoring(q *msg.Request, mr *msg.Result) {
	var err error
	var res sql.Result

	if res, err = s.stmtRevokeAuthorizationMonitoring.Exec(
		q.Grant.Id,
		q.Grant.PermissionId,
		q.Grant.Category,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Grant = append(mr.Grant, q.Grant)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
