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

func (s *Supervisor) actionWrite(q *msg.Request, mr *msg.Result) {
	if s.readonly {
		mr.ReadOnly()
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return
	}

	switch q.Action {
	case msg.ActionAdd:
		s.actionAdd(q, mr)
	case msg.ActionRemove:
		s.actionRemove(q, mr)
	}

	if mr.IsOK() {
		go func() {
			s.Update <- msg.CacheUpdateFromRequest(q)
		}()
	}
}

func (s *Supervisor) actionAdd(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	q.ActionObj.ID = uuid.Must(uuid.NewV4()).String()
	if res, err = s.stmtActionAdd.Exec(
		q.ActionObj.ID,
		q.ActionObj.Name,
		q.ActionObj.SectionID,
		q.AuthUser,
	); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.ActionObj = append(mr.ActionObj, q.ActionObj)
		mr.Super.Audit.WithField(`Code`, mr.Code).Infoln(fmt.Sprintf(
			"Successfully added action %s", q.ActionObj.Name))
		return
	}
	mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
}

func (s *Supervisor) actionRemove(q *msg.Request, mr *msg.Result) {
	var (
		err error
		tx  *sql.Tx
		res sql.Result
	)
	txMap := map[string]*sql.Stmt{}

	// open multi-statement transaction
	if tx, err = s.conn.Begin(); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
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
			mr.ServerError(err, q.Section)
			mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
			tx.Rollback()
			return
		}
	}

	if res, err = s.actionRemoveTx(q.ActionObj.ID,
		txMap); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		tx.Rollback()
		return
	}

	// sets r.OK()
	if !mr.RowCnt(res.RowsAffected()) {
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		tx.Rollback()
		return
	}

	// close transaction
	if err = tx.Commit(); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	mr.ActionObj = append(mr.ActionObj, q.ActionObj)
	mr.Super.Audit.WithField(`Code`, mr.Code).
		Infoln(fmt.Sprintf(
			"Successfully removed action %s", q.ActionObj.ID))
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
