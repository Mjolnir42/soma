/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"database/sql"
	"time"

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/proto"
)

func (s *Supervisor) actionRead(q *msg.Request, mr *msg.Result) {
	switch q.Action {
	case msg.ActionList:
		s.actionList(q, mr)
	case msg.ActionShow:
		s.actionShow(q, mr)
	case msg.ActionSearch:
		s.actionSearch(q, mr)
	}
}

func (s *Supervisor) actionList(q *msg.Request, mr *msg.Result) {
	var (
		err                             error
		rows                            *sql.Rows
		actionID, actionName, sectionID string
	)

	if rows, err = s.stmtActionList.Query(
		q.ActionObj.SectionID,
	); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&actionID,
			&actionName,
			&sectionID,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
			return
		}
		mr.ActionObj = append(mr.ActionObj, proto.Action{
			ID:        actionID,
			Name:      actionName,
			SectionID: sectionID,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}
	mr.OK()
	mr.Super.Audit.WithField(`Code`, mr.Code).Infoln(`OK`)
}

func (s *Supervisor) actionShow(q *msg.Request, mr *msg.Result) {
	var (
		err                             error
		ts                              time.Time
		actionID, actionName, sectionID string
		category, user, sectionName     string
	)

	if err = s.stmtActionShow.QueryRow(
		q.ActionObj.ID,
	).Scan(
		&actionID,
		&actionName,
		&sectionID,
		&sectionName,
		&category,
		&user,
		&ts,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}
	mr.ActionObj = append(mr.ActionObj, proto.Action{
		ID:          actionID,
		Name:        actionName,
		SectionID:   sectionID,
		SectionName: sectionName,
		Category:    category,
		Details: &proto.DetailsCreation{
			CreatedBy: user,
			CreatedAt: ts.Format(msg.RFC3339Milli),
		},
	})
	mr.OK()
	mr.Super.Audit.WithField(`Code`, mr.Code).Infoln(`OK`)
}

func (s *Supervisor) actionSearch(q *msg.Request, mr *msg.Result) {
	var (
		err                             error
		rows                            *sql.Rows
		actionID, actionName, sectionID string
	)

	if rows, err = s.stmtActionSearch.Query(
		q.Search.ActionObj.Name,
		q.Search.ActionObj.SectionID,
	); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&actionID,
			&actionName,
			&sectionID,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
			return
		}
		mr.ActionObj = append(mr.ActionObj, proto.Action{
			ID:        actionID,
			Name:      actionName,
			SectionID: sectionID,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}
	mr.OK()
	mr.Super.Audit.WithField(`Code`, mr.Code).Infoln(`OK`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
