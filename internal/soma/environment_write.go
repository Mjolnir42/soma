/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2015-2017, Jörg Pernfuß
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

// EnvironmentWrite handles write requests for environments
type EnvironmentWrite struct {
	Input      chan msg.Request
	Shutdown   chan struct{}
	conn       *sql.DB
	stmtCreate *sql.Stmt
	stmtDelete *sql.Stmt
	stmtRename *sql.Stmt
	appLog     *logrus.Logger
	reqLog     *logrus.Logger
	errLog     *logrus.Logger
}

// newEnvironmentWrite return a new EnvironmentWrite handler with
// input buffer of length
func newEnvironmentWrite(length int) (w *EnvironmentWrite) {
	w = &EnvironmentWrite{}
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	return
}

// Register initializes resources provided by the Soma app
func (w *EnvironmentWrite) Register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// Intake exposes the Input channel as part of the handler interface
func (w *EnvironmentWrite) Intake() chan msg.Request {
	return w.Input
}

// Run is the event loop for EnvironmentWrite
func (w *EnvironmentWrite) Run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.EnvironmentAdd:    w.stmtCreate,
		stmt.EnvironmentRemove: w.stmtDelete,
		stmt.EnvironmentRename: w.stmtRename,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`environment`, err, stmt.Name(statement))
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
func (w *EnvironmentWrite) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(w.reqLog, q)

	switch q.Action {
	case msg.ActionAdd:
		w.add(q, &result)
	case msg.ActionRemove:
		w.remove(q, &result)
	case msg.ActionRename:
		w.rename(q, &result)
	default:
		result.UnknownRequest(q)
	}

	q.Reply <- result
}

// add inserts a new environment
func (w *EnvironmentWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtCreate.Exec(
		q.Environment.Name,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Environment = append(mr.Environment, q.Environment)
	}
}

// remove deletes an environment
func (w *EnvironmentWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtDelete.Exec(
		q.Environment.Name,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Environment = append(mr.Environment, q.Environment)
	}
}

// rename changes the name of an environment
func (w *EnvironmentWrite) rename(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtRename.Exec(
		q.Update.Environment.Name,
		q.Environment.Name,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Environment = append(mr.Environment, q.Update.Environment)
	}
}

// ShutdownNow signals the handler to shut down
func (w *EnvironmentWrite) ShutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
