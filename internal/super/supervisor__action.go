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
	"time"

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

func (s *Supervisor) action(q *msg.Request) {
	result := msg.FromRequest(q)

	s.requestLog(q)

	switch q.Action {
	case `list`, `show`, `search`:
		go func() { s.actionRead(q) }()
	case `add`, `remove`:
		if s.readonly {
			result.Conflict(fmt.Errorf(`Readonly instance`))
			goto abort
		}
		s.actionWrite(q)
	default:
		result.UnknownRequest(q)
		goto abort
	}
	return

abort:
	q.Reply <- result
}

func (s *Supervisor) actionRead(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case `list`:
		s.actionList(q, &result)
	case `show`:
		s.actionShow(q, &result)
	case `search`:
		s.actionSearch(q, &result)
	}

	q.Reply <- result
}

func (s *Supervisor) actionList(q *msg.Request, r *msg.Result) {
	r.ActionObj = []proto.Action{}
	var (
		err                             error
		rows                            *sql.Rows
		actionID, actionName, sectionID string
	)
	if rows, err = s.stmtActionList.Query(); err != nil {
		r.ServerError(err, q.Section)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&actionID,
			&actionName,
			&sectionID,
		); err != nil {
			r.ServerError(err, q.Section)
			return
		}
		r.ActionObj = append(r.ActionObj,
			proto.Action{
				Id:        actionID,
				Name:      actionName,
				SectionId: sectionID,
			})
	}
	if err = rows.Err(); err != nil {
		r.ServerError(err, q.Section)
		return
	}
	r.OK()
}

func (s *Supervisor) actionShow(q *msg.Request, r *msg.Result) {
	var (
		err                             error
		ts                              time.Time
		actionID, actionName, sectionID string
		category, user, sectionName     string
	)
	if err = s.stmtActionShow.QueryRow(q.ActionObj.Id).Scan(
		&actionID,
		&actionName,
		&sectionID,
		&sectionName,
		&category,
		&user,
		&ts,
	); err == sql.ErrNoRows {
		r.NotFound(err, q.Section)
		return
	} else if err != nil {
		r.ServerError(err, q.Section)
		return
	}
	r.ActionObj = []proto.Action{proto.Action{
		Id:          actionID,
		Name:        actionName,
		SectionId:   sectionID,
		SectionName: sectionName,
		Category:    category,
		Details: &proto.DetailsCreation{
			CreatedBy: user,
			CreatedAt: ts.Format(msg.RFC3339Milli),
		},
	}}
	r.OK()
}

func (s *Supervisor) actionSearch(q *msg.Request, r *msg.Result) {
	r.ActionObj = []proto.Action{}
	var (
		err                             error
		rows                            *sql.Rows
		actionID, actionName, sectionID string
	)
	if rows, err = s.stmtActionSearch.Query(
		q.ActionObj.Name,
		q.ActionObj.SectionId,
	); err != nil {
		r.ServerError(err, q.Section)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&actionID,
			&actionName,
			&sectionID,
		); err != nil {
			r.ServerError(err, q.Section)
			return
		}
		r.ActionObj = append(r.ActionObj,
			proto.Action{
				Id:        actionID,
				Name:      actionName,
				SectionId: sectionID,
			})
	}
	if err = rows.Err(); err != nil {
		r.ServerError(err, q.Section)
		return
	}
	r.OK()
}

func (s *Supervisor) actionWrite(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case `add`:
		s.actionAdd(q, &result)
	case `remove`:
		s.actionRemove(q, &result)
	}

	if result.IsOK() {
		s.Update <- msg.CacheUpdateFromRequest(q)
	}

	q.Reply <- result
}

func (s *Supervisor) actionAdd(q *msg.Request, r *msg.Result) {
	var (
		err error
		res sql.Result
	)
	q.ActionObj.Id = uuid.NewV4().String()
	if res, err = s.stmtActionAdd.Exec(
		q.ActionObj.Id,
		q.ActionObj.Name,
		q.ActionObj.SectionId,
		q.AuthUser,
	); err != nil {
		r.ServerError(err, q.Section)
		return
	}
	if r.RowCnt(res.RowsAffected()) {
		r.ActionObj = []proto.Action{q.ActionObj}
	}
}

func (s *Supervisor) actionRemove(q *msg.Request, r *msg.Result) {
	var (
		err error
		tx  *sql.Tx
		res sql.Result
	)
	txMap := map[string]*sql.Stmt{}

	// open multi-statement transaction
	if tx, err = s.conn.Begin(); err != nil {
		r.ServerError(err, q.Section)
		return
	}

	// prepare statements for this transaction
	for name, statement := range map[string]string{
		`action_tx_remove`:    stmt.ActionRemove,
		`action_tx_removeMap`: stmt.ActionRemoveFromMap,
	} {
		if txMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("s.ActionTx.Prepare(%s) error: %s",
				name, err.Error())
			r.ServerError(err, q.Section)
			tx.Rollback()
			return
		}
	}

	if res, err = s.actionRemoveTx(q.ActionObj.Id,
		txMap); err != nil {
		r.ServerError(err, q.Section)
		tx.Rollback()
		return
	}
	// sets r.OK()
	if !r.RowCnt(res.RowsAffected()) {
		tx.Rollback()
		return
	}

	// close transaction
	if err = tx.Commit(); err != nil {
		r.ServerError(err, q.Section)
		return
	}

	r.ActionObj = []proto.Action{q.ActionObj}
}

func (s *Supervisor) actionRemoveTx(id string,
	txMap map[string]*sql.Stmt) (sql.Result, error) {
	var (
		err error
		res sql.Result
	)

	// remove action from all permissions
	if res, err = txMap[`action_tx_removeMap`].Exec(
		id); err != nil {
		return res, err
	}

	// remove action
	return txMap[`action_tx_remove`].Exec(id)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
