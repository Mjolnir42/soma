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
	"github.com/mjolnir42/soma/lib/proto"
)

func (s *Supervisor) rightRead(q *msg.Request, mr *msg.Result) {
	switch q.Action {
	case msg.ActionList:
		switch q.Grant.RecipientType {
		case msg.SubjectUser:
		case msg.SubjectAdmin:
		case msg.SubjectTeam:
		default:
			mr.NotImplemented(
				fmt.Errorf("Rights for recipient type"+
					" %s are currently not implemented",
					q.Grant.RecipientType))
			mr.Super.Audit.
				WithField(`Code`, mr.Code).
				Warningln(mr.Error)
			return
		}

		switch q.Grant.Category {
		case msg.CategorySystem,
			msg.CategoryGlobal,
			msg.CategoryGrantGlobal,
			msg.CategoryIdentity,
			msg.CategoryGrantIdentity,
			msg.CategorySelf,
			msg.CategoryGrantSelf,
			msg.CategoryPermission,
			msg.CategoryGrantPermission,
			msg.CategoryOperation,
			msg.CategoryGrantOperation:
			s.rightListGlobal(q, mr)
		case msg.CategoryRepository,
			msg.CategoryGrantRepository,
			msg.CategoryTeam,
			msg.CategoryGrantTeam,
			msg.CategoryMonitoring,
			msg.CategoryGrantMonitoring:
			s.rightListScoped(q, mr)
		}
	case msg.ActionShow:
		switch q.Grant.RecipientType {
		case msg.SubjectUser:
		case msg.SubjectAdmin:
		case msg.SubjectTeam:
		default:
			mr.NotImplemented(
				fmt.Errorf("Rights for recipient type"+
					" %s are currently not implemented",
					q.Grant.RecipientType))
			mr.Super.Audit.
				WithField(`Code`, mr.Code).
				Warningln(mr.Error)
			return
		}

		switch q.Grant.Category {
		case msg.CategorySystem,
			msg.CategoryGlobal,
			msg.CategoryGrantGlobal,
			msg.CategoryIdentity,
			msg.CategoryGrantIdentity,
			msg.CategorySelf,
			msg.CategoryGrantSelf,
			msg.CategoryPermission,
			msg.CategoryGrantPermission,
			msg.CategoryOperation,
			msg.CategoryGrantOperation:
			s.rightShowGlobal(q, mr)
		case msg.CategoryRepository,
			msg.CategoryGrantRepository,
			msg.CategoryTeam,
			msg.CategoryGrantTeam,
			msg.CategoryMonitoring,
			msg.CategoryGrantMonitoring:
			s.rightShowScoped(q, mr)
		}
	case msg.ActionSearch:
		switch q.Search.Grant.RecipientType {
		case msg.SubjectUser:
		case msg.SubjectAdmin:
		case msg.SubjectTeam:
		default:
			mr.NotImplemented(
				fmt.Errorf("Rights for recipient type"+
					" %s are currently not implemented",
					q.Search.Grant.RecipientType))
			mr.Super.Audit.
				WithField(`Code`, mr.Code).
				Warningln(mr.Error)
			return
		}

		switch q.Search.Grant.Category {
		case msg.CategorySystem,
			msg.CategoryGlobal,
			msg.CategoryGrantGlobal,
			msg.CategoryIdentity,
			msg.CategoryGrantIdentity,
			msg.CategorySelf,
			msg.CategoryGrantSelf,
			msg.CategoryPermission,
			msg.CategoryGrantPermission,
			msg.CategoryOperation,
			msg.CategoryGrantOperation:
			s.rightSearchGlobal(q, mr)
		case msg.CategoryRepository,
			msg.CategoryGrantRepository,
			msg.CategoryTeam,
			msg.CategoryGrantTeam,
			msg.CategoryMonitoring,
			msg.CategoryGrantMonitoring:
			s.rightSearchScoped(q, mr)
		}
	}
}

func (s *Supervisor) rightListGlobal(q *msg.Request, mr *msg.Result) {
	var (
		err                             error
		rows                            *sql.Rows
		grantID                         string
		adminID, userID, toolID, teamID sql.NullString
	)

	if rows, err = s.stmtListAuthorizationGlobal.Query(
		q.Search.Grant.PermissionID,
		q.Search.Grant.Category,
	); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&grantID,
			&adminID,
			&userID,
			&toolID,
			&teamID,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
			return
		}
		grant := proto.Grant{
			ID:           grantID,
			PermissionID: q.Search.Grant.PermissionID,
			Category:     q.Search.Grant.Category,
		}
		switch {
		case adminID.Valid:
			grant.RecipientID = adminID.String
			grant.RecipientType = msg.SubjectAdmin
		case userID.Valid:
			grant.RecipientID = userID.String
			grant.RecipientType = msg.SubjectUser
		case toolID.Valid:
			grant.RecipientID = toolID.String
			grant.RecipientType = msg.SubjectTool
		case teamID.Valid:
			grant.RecipientID = teamID.String
			grant.RecipientType = msg.SubjectTeam
		}
		mr.Grant = append(mr.Grant, grant)
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}
	mr.OK()
	mr.Super.Audit.WithField(`Code`, mr.Code).Infoln(`OK`)
}

func (s *Supervisor) rightListScoped(q *msg.Request, mr *msg.Result) {
	// XXX BUG TODO
}

func (s *Supervisor) rightShowGlobal(q *msg.Request, mr *msg.Result) {
	// XXX BUG TODO
}

func (s *Supervisor) rightShowScoped(q *msg.Request, mr *msg.Result) {
	// XXX BUG TODO
}

func (s *Supervisor) rightSearchGlobal(q *msg.Request, mr *msg.Result) {
	var (
		err     error
		grantID string
	)

	if err = s.stmtSearchAuthorizationGlobal.QueryRow(
		q.Search.Grant.PermissionID,
		q.Search.Grant.Category,
		q.Search.Grant.RecipientID,
		q.Search.Grant.RecipientType,
	).Scan(
		&grantID,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}
	mr.Grant = append(mr.Grant, proto.Grant{
		ID:            grantID,
		PermissionID:  q.Search.Grant.PermissionID,
		Category:      q.Search.Grant.Category,
		RecipientID:   q.Search.Grant.RecipientID,
		RecipientType: q.Search.Grant.RecipientType,
	})
	mr.OK()
	mr.Super.Audit.WithField(`Code`, mr.Code).Infoln(`OK`)
}

func (s *Supervisor) rightSearchScoped(q *msg.Request, mr *msg.Result) {
	var (
		err     error
		grantID string
		scope   *sql.Stmt
	)

	switch q.Search.Grant.Category {
	case msg.CategoryRepository,
		msg.CategoryGrantRepository:
		scope = s.stmtSearchAuthorizationRepository
	case msg.CategoryTeam,
		msg.CategoryGrantTeam:
		scope = s.stmtSearchAuthorizationTeam
	case msg.CategoryMonitoring,
		msg.CategoryGrantMonitoring:
		scope = s.stmtSearchAuthorizationMonitoring
	default:
		err = fmt.Errorf("Unhandled search category: %s", q.Search.Grant.Category)
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}
	if err = scope.QueryRow(
		q.Search.Grant.PermissionID,
		q.Search.Grant.Category,
		q.Search.Grant.RecipientID,
		q.Search.Grant.RecipientType,
		q.Search.Grant.ObjectType,
		q.Search.Grant.ObjectID,
	).Scan(
		&grantID,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}
	mr.Grant = append(mr.Grant, proto.Grant{
		ID:            grantID,
		PermissionID:  q.Search.Grant.PermissionID,
		Category:      q.Search.Grant.Category,
		RecipientID:   q.Search.Grant.RecipientID,
		RecipientType: q.Search.Grant.RecipientType,
		ObjectType:    q.Search.Grant.ObjectType,
		ObjectID:      q.Search.Grant.ObjectID,
	})
	mr.OK()
	mr.Super.Audit.WithField(`Code`, mr.Code).Infoln(`OK`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
