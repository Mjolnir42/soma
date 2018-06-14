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

func (s *Supervisor) sectionWrite(q *msg.Request, mr *msg.Result) {
	if s.readonly {
		mr.ReadOnly()
		mr.Super.Audit.
			WithField(`Code`, mr.Code).
			Warningln(mr.Error)
		return
	}

	switch q.Action {
	case msg.ActionAdd:
		s.sectionAdd(q, mr)
	case msg.ActionRemove:
		s.sectionRemove(q, mr)
	}

	if mr.IsOK() {
		go func() {
			s.Update <- msg.CacheUpdateFromRequest(q)
		}()
	}
}

func (s *Supervisor) sectionAdd(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	q.SectionObj.ID = uuid.Must(uuid.NewV4()).String()
	if res, err = s.stmtSectionAdd.Exec(
		q.SectionObj.ID,
		q.SectionObj.Name,
		q.SectionObj.Category,
		q.AuthUser,
	); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.SectionObj = append(mr.SectionObj, q.SectionObj)
		mr.Super.Audit.WithField(`Code`, mr.Code).Infoln(fmt.Sprintf(
			"Successfully added section %s", q.SectionObj.Name))
		return
	}
	mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
}

func (s *Supervisor) sectionRemove(q *msg.Request, mr *msg.Result) {
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
		`action_tx_remove`:     stmt.ActionRemove,
		`action_tx_removeMap`:  stmt.ActionRemoveFromMap,
		`section_tx_remove`:    stmt.SectionRemove,
		`section_tx_removeMap`: stmt.SectionRemoveFromMap,
		`section_tx_actlist`:   stmt.SectionListActions,
	} {
		if txMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("s.SectionTx.Prepare(%s) error: %s",
				name, err.Error())
			mr.ServerError(err, q.Section)
			mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
			tx.Rollback()
			return
		}
	}

	if res, err = s.sectionRemoveTx(q.SectionObj.ID,
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
			"Successfully removed section %s", q.SectionObj.ID))
	return
}

func (s *Supervisor) sectionRemoveTx(id string,
	txMap map[string]*sql.Stmt) (sql.Result, error) {
	var (
		err      error
		res      sql.Result
		rows     *sql.Rows
		actionID string
		affected int64
	)

	// remove all actions in this section
	if rows, err = txMap[`section_tx_actlist`].Query(
		id); err != nil {
		return res, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&actionID,
		); err != nil {
			rows.Close()
			return res, err
		}
		if res, err = s.actionRemoveTx(actionID, txMap); err != nil {
			rows.Close()
			return res, err
		}
		if affected, err = res.RowsAffected(); err != nil {
			rows.Close()
			return res, err
		} else if affected != 1 {
			rows.Close()
			return res, fmt.Errorf("Delete statement caught %d rows"+
				" of actions instead of 1 (actionID=%s)", affected,
				actionID)
		}
	}
	if err = rows.Err(); err != nil {
		return res, err
	}

	// remove section from all permissions
	if res, err = txMap[`section_tx_removeMap`].Exec(id); err != nil {
		return res, err
	}

	// remove section
	return txMap[`section_tx_remove`].Exec(id)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
