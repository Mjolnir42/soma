/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"database/sql"

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
)

// EntityWrite handles write requests for object entities
type EntityWrite struct {
	Input      chan msg.Request
	Shutdown   chan struct{}
	conn       *sql.DB
	stmtAdd    *sql.Stmt
	stmtRemove *sql.Stmt
	stmtRename *sql.Stmt
	appLog     *logrus.Logger
	reqLog     *logrus.Logger
	errLog     *logrus.Logger
}

// newEntityWrite return a new EntityWrite handler with
// input buffer of length
func newEntityWrite(length int) (w *EntityWrite) {
	w = &EntityWrite{}
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (w *EntityWrite) register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// run is the event loop for EntityWrite
func (w *EntityWrite) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.EntityAdd:    w.stmtAdd,
		stmt.EntityDel:    w.stmtRemove,
		stmt.EntityRename: w.stmtRename,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`entity`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-w.Shutdown:
			break runloop
		case req := <-w.Input:
			w.process(&req)
		}
	}
}

// process is the request dispatcher
func (w *EntityWrite) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(w.reqLog, q)

	switch q.Action {
	case msg.ActionCreate:
		w.create(q, &result)
	case msg.ActionDelete:
		w.delete(q, &result)
	case msg.ActionRename:
		w.rename(q, &result)
	default:
		result.UnknownRequest(q)
	}

	q.Reply <- result
}

// create adds a new entity
func (w *EntityWrite) create(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtAdd.Exec(
		q.Entity.Name,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Entity = append(mr.Entity, q.Entity)
	}
}

// delete removes an entity
func (w *EntityWrite) delete(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtRemove.Exec(
		q.Entity.Name,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Entity = append(mr.Entity, q.Entity)
	}
}

// rename changes a entity's name
func (w *EntityWrite) rename(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtRemove.Exec(
		q.Update.Entity.Name,
		q.Entity.Name,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Entity = append(mr.Entity, q.Update.Entity)
	}
}

// shutdownNow signals the handler to shut down
func (w *EntityWrite) shutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
