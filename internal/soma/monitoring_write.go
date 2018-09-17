/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2018, Jörg Pernfuß
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
	uuid "github.com/satori/go.uuid"
)

// MonitoringWrite handles write requests for monitoring systems
type MonitoringWrite struct {
	Input       chan msg.Request
	Shutdown    chan struct{}
	handlerName string
	conn        *sql.DB
	stmtCreate  *sql.Stmt
	stmtDelete  *sql.Stmt
	appLog      *logrus.Logger
	reqLog      *logrus.Logger
	errLog      *logrus.Logger
}

// newMonitoringWrite return a new MonitoringWrite handler with
// input buffer of length
func newMonitoringWrite(length int) (string, *MonitoringWrite) {
	w := &MonitoringWrite{}
	w.handlerName = generateHandlerName() + `_w`
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	return w.handlerName, w
}

// Register initializes resources provided by the Soma app
func (w *MonitoringWrite) Register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (w *MonitoringWrite) RegisterRequests(hmap *handler.Map) {
	for _, action := range []string{
		msg.ActionAdd,
		msg.ActionRemove,
	} {
		hmap.Request(msg.SectionMonitoringMgmt, action, w.handlerName)
	}
}

// Intake exposes the Input channel as part of the handler interface
func (w *MonitoringWrite) Intake() chan msg.Request {
	return w.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (w *MonitoringWrite) PriorityIntake() chan msg.Request {
	return w.Intake()
}

// Run is the event loop for MonitoringWrite
func (w *MonitoringWrite) Run() {
	var err error

	for statement, prepStmt := range map[string]**sql.Stmt{
		stmt.MonitoringSystemAdd:    &w.stmtCreate,
		stmt.MonitoringSystemRemove: &w.stmtDelete,
	} {
		if *prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`monitoring`, err, stmt.Name(statement))
		}
		defer (*prepStmt).Close()
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
func (w *MonitoringWrite) process(q *msg.Request) {
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

// add inserts a new monitoring system
func (w *MonitoringWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		err      error
		res      sql.Result
		callback sql.NullString
	)

	q.Monitoring.ID = uuid.Must(uuid.NewV4()).String()
	if q.Monitoring.Callback != `` {
		callback = sql.NullString{
			String: q.Monitoring.Callback,
			Valid:  true,
		}
	}
	if res, err = w.stmtCreate.Exec(
		q.Monitoring.ID,
		q.Monitoring.Name,
		q.Monitoring.Mode,
		q.Monitoring.Contact,
		q.Monitoring.TeamID,
		callback,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Monitoring = append(mr.Monitoring, q.Monitoring)
	}
}

// remove deletes a monitoring system
func (w *MonitoringWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtDelete.Exec(
		q.Monitoring.ID,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Monitoring = append(mr.Monitoring, q.Monitoring)
	}
}

// ShutdownNow signals the handler to shut down
func (w *MonitoringWrite) ShutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
