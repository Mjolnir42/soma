/*-
 * Copyright (c) 2015-2018, Jörg Pernfuß
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

// DatacenterWrite handles write requests for datacenters
type DatacenterWrite struct {
	Input       chan msg.Request
	Shutdown    chan struct{}
	handlerName string
	conn        *sql.DB
	stmtAdd     *sql.Stmt
	stmtRemove  *sql.Stmt
	stmtRename  *sql.Stmt
	appLog      *logrus.Logger
	reqLog      *logrus.Logger
	errLog      *logrus.Logger
}

// newDatacenterWrite return a new DatacenterWrite handler with input
// buffer of length
func newDatacenterWrite(length int) (string, *DatacenterWrite) {
	w := &DatacenterWrite{}
	w.handlerName = generateHandlerName() + `_w`
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	return w.handlerName, w
}

// Register initializes resources provided by the Soma app
func (w *DatacenterWrite) Register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (w *DatacenterWrite) RegisterRequests(hmap *handler.Map) {
	for _, action := range []string{
		msg.ActionAdd,
		msg.ActionRemove,
		msg.ActionRename,
	} {
		hmap.Request(msg.SectionDatacenter, action, w.handlerName)
	}
}

// Intake exposes the Input channel as part of the handler interface
func (w *DatacenterWrite) Intake() chan msg.Request {
	return w.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (w *DatacenterWrite) PriorityIntake() chan msg.Request {
	return w.Intake()
}

// Run is the event loop for DatacenterWrite
func (w *DatacenterWrite) Run() {
	var err error

	for statement, prepStmt := range map[string]**sql.Stmt{
		stmt.DatacenterAdd:    &w.stmtAdd,
		stmt.DatacenterDel:    &w.stmtRemove,
		stmt.DatacenterRename: &w.stmtRename,
	} {
		if *prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`datacenter`, err, stmt.Name(statement))
		}
		defer (*prepStmt).Close()
	}

	for {
		select {
		case <-w.Shutdown:
			break
		case req := <-w.Input:
			w.process(&req)
		}
	}
}

// process is the request dispatcher
func (w *DatacenterWrite) process(q *msg.Request) {
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

// add inserts a new datacenter
func (w *DatacenterWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
	)

	if res, err = w.stmtAdd.Exec(
		q.Datacenter.LoCode,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Datacenter = append(mr.Datacenter, q.Datacenter)
	}
}

// remove deletes a datacenter
func (w *DatacenterWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
	)

	if res, err = w.stmtRemove.Exec(
		q.Datacenter.LoCode,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Datacenter = append(mr.Datacenter, q.Datacenter)
	}
}

// rename changes a datacenter's locode
func (w *DatacenterWrite) rename(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
	)

	if res, err = w.stmtRename.Exec(
		q.Update.Datacenter.LoCode,
		q.Datacenter.LoCode,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Datacenter = append(mr.Datacenter, q.Datacenter)
	}
}

// ShutdownNow signals the handler to shut down
func (w *DatacenterWrite) ShutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
