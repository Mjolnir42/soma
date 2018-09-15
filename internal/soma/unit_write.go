/*-
 * Copyright (c) 2016-2018, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma // import "github.com/mjolnir42/soma/internal/soma"

import (
	"database/sql"

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/handler"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
)

// UnitWrite handles write requests for units
type UnitWrite struct {
	Input       chan msg.Request
	Shutdown    chan struct{}
	handlerName string
	conn        *sql.DB
	stmtAdd     *sql.Stmt
	stmtRemove  *sql.Stmt
	appLog      *logrus.Logger
	reqLog      *logrus.Logger
	errLog      *logrus.Logger
}

// newUnitWrite return a new UnitWrite handler with input buffer of
// length
func newUnitWrite(length int) (string, *UnitWrite) {
	w := &UnitWrite{}
	w.handlerName = generateHandlerName() + `_w`
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	return w.handlerName, w
}

// Register initializes resources provided by the Soma app
func (w *UnitWrite) Register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (w *UnitWrite) RegisterRequests(hmap *handler.Map) {
	for _, action := range []string{
		msg.ActionAdd,
		msg.ActionRemove,
	} {
		hmap.Request(msg.SectionUnit, action, w.handlerName)
	}
}

// Intake exposes the Input channel as part of the handler interface
func (w *UnitWrite) Intake() chan msg.Request {
	return w.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (w *UnitWrite) PriorityIntake() chan msg.Request {
	return w.Intake()
}

// Run is the event loop for UnitWrite
func (w *UnitWrite) Run() {
	var err error

	for statement, prepStmt := range map[string]**sql.Stmt{
		stmt.UnitAdd: &w.stmtAdd,
		stmt.UnitDel: &w.stmtRemove,
	} {
		if *prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`unit`, err, stmt.Name(statement))
		}
		defer (*prepStmt).Close()()
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
func (w *UnitWrite) process(q *msg.Request) {
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

// add inserts a new unit
func (w *UnitWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtAdd.Exec(
		q.Unit.Unit,
		q.Unit.Name,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Unit = append(mr.Unit, q.Unit)
	}
}

// remove deletes a unit
func (w *UnitWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtRemove.Exec(
		q.Unit.Unit,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Unit = append(mr.Unit, q.Unit)
	}
}

// ShutdownNow signals the handler to shut down
func (w *UnitWrite) ShutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
