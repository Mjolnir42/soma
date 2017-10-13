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

func (s *Supervisor) sectionRead(q *msg.Request, mr *msg.Result) {
	switch q.Action {
	case msg.ActionList:
		s.sectionList(q, mr)
	case msg.ActionShow:
		s.sectionShow(q, mr)
	case msg.ActionSearch:
		s.sectionSearch(q, mr)
	}
}

func (s *Supervisor) sectionList(q *msg.Request, mr *msg.Result) {
	var (
		err                    error
		rows                   *sql.Rows
		sectionID, sectionName string
	)

	if rows, err = s.stmtSectionList.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&sectionID,
			&sectionName,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.SectionObj = append(mr.SectionObj, proto.Section{
			Id:   sectionID,
			Name: sectionName,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

func (s *Supervisor) sectionShow(q *msg.Request, mr *msg.Result) {
	var (
		err                                    error
		sectionID, sectionName, category, user string
		ts                                     time.Time
	)

	if err = s.stmtSectionShow.QueryRow(
		q.SectionObj.Id,
	).Scan(
		&sectionID,
		&sectionName,
		&category,
		&user,
		&ts,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.SectionObj = append(mr.SectionObj, proto.Section{
		Id:       sectionID,
		Name:     sectionName,
		Category: category,
		Details: &proto.DetailsCreation{
			CreatedAt: ts.Format(msg.RFC3339Milli),
			CreatedBy: user,
		},
	})
	mr.OK()
}

func (s *Supervisor) sectionSearch(q *msg.Request, mr *msg.Result) {
	var (
		err                    error
		rows                   *sql.Rows
		sectionID, sectionName string
	)

	if rows, err = s.stmtSectionSearch.Query(
		q.SectionObj.Name,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&sectionID,
			&sectionName,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.SectionObj = append(mr.SectionObj, proto.Section{
			Id:   sectionID,
			Name: sectionName,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix