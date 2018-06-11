/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2017, Jörg Pernfuß
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

// AttributeWrite handles write requests for alert levels
type AttributeWrite struct {
	Input      chan msg.Request
	Shutdown   chan struct{}
	conn       *sql.DB
	stmtAdd    *sql.Stmt
	stmtRemove *sql.Stmt
	appLog     *logrus.Logger
	reqLog     *logrus.Logger
	errLog     *logrus.Logger
}

// newAttributeWrite return a new AttributeWrite handler with input buffer of
// length
func newAttributeWrite(length int) (w *AttributeWrite) {
	w = &AttributeWrite{}
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	return
}

// Register initializes resources provided by the Soma app
func (w *AttributeWrite) Register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// Intake exposes the Input channel as part of the handler interface
func (w *AttributeWrite) Intake() chan msg.Request {
	return w.Input
}

// Run is the event loop for AttributeWrite
func (w *AttributeWrite) Run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.AttributeAdd:    w.stmtAdd,
		stmt.AttributeRemove: w.stmtRemove,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`attribute`, err, stmt.Name(statement))
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
func (w *AttributeWrite) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(w.reqLog, q)

	switch q.Action {
	case msg.ActionAdd:
		w.add(q, &result)
	case msg.ActionRemove:
		w.remove(q, &result)
	default:
		result.UnknownRequest(q)
	}

	q.Reply <- result
}

// add inserts a new attribute
func (w *AttributeWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtAdd.Exec(
		q.Attribute.Name,
		q.Attribute.Cardinality,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Attribute = append(mr.Attribute, q.Attribute)
	}
}

// remove deletes an attribute
func (w *AttributeWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtRemove.Exec(
		q.Attribute.Name,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Attribute = append(mr.Attribute, q.Attribute)
	}
}

// ShutdownNow signals the handler to shut down
func (w *AttributeWrite) ShutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
