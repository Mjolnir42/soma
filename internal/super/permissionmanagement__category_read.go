/*-
 * Copyright (c) 2016, 1&1 Internet SE
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

func (s *Supervisor) categoryRead(q *msg.Request, mr *msg.Result) {
	switch q.Action {
	case msg.ActionList:
		s.categoryList(q, mr)
	case msg.ActionShow:
		s.categoryShow(q, mr)
	}
}

func (s *Supervisor) categoryList(q *msg.Request, mr *msg.Result) {
	var (
		err      error
		rows     *sql.Rows
		category string
	)

	if rows, err = s.stmtCategoryList.Query(); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&category,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
			return
		}
		mr.Category = append(mr.Category, proto.Category{
			Name: category,
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

func (s *Supervisor) categoryShow(q *msg.Request, mr *msg.Result) {
	var (
		err            error
		category, user string
		ts             time.Time
	)

	if err = s.stmtCategoryShow.QueryRow(
		q.Category.Name,
	).Scan(
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
	mr.Category = append(mr.Category, proto.Category{
		Name: category,
		Details: &proto.CategoryDetails{
			Creation: &proto.DetailsCreation{
				CreatedAt: ts.Format(msg.RFC3339Milli),
				CreatedBy: user,
			},
		},
	})
	mr.OK()
	mr.Super.Audit.WithField(`Code`, mr.Code).Infoln(`OK`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
