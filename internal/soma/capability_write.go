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
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/handler"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	uuid "github.com/satori/go.uuid"
)

// CapabilityWrite handles write requests for capabilities
type CapabilityWrite struct {
	Input                chan msg.Request
	Shutdown             chan struct{}
	handlerName          string
	conn                 *sql.DB
	stmtAdd              *sql.Stmt
	stmtRemove           *sql.Stmt
	stmtVerifyMetric     *sql.Stmt
	stmtVerifyMonitoring *sql.Stmt
	stmtVerifyView       *sql.Stmt
	appLog               *logrus.Logger
	reqLog               *logrus.Logger
	errLog               *logrus.Logger
}

// newCapabilityWrite return a new CapabilityWrite handler with
// input buffer of length
func newCapabilityWrite(length int) (string, *CapabilityWrite) {
	w := &CapabilityWrite{}
	w.handlerName = generateHandlerName() + `_w`
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	return w.handlerName, w
}

// Register initializes resources provided by the Soma app
func (w *CapabilityWrite) Register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (w *CapabilityWrite) RegisterRequests(hmap *handler.Map) {
	for _, action := range []string{
		msg.ActionAdd,
		msg.ActionRemove,
	} {
		hmap.Request(msg.SectionCapability, action, w.handlerName)
	}
}

// Intake exposes the Input channel as part of the handler interface
func (w *CapabilityWrite) Intake() chan msg.Request {
	return w.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (w *CapabilityWrite) PriorityIntake() chan msg.Request {
	return w.Intake()
}

// Run is the event loop for CapabilityWrite
func (w *CapabilityWrite) Run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.AddCapability:          w.stmtAdd,
		stmt.DelCapability:          w.stmtRemove,
		stmt.MetricVerify:           w.stmtVerifyMetric,
		stmt.VerifyMonitoringSystem: w.stmtVerifyMonitoring,
		stmt.ViewVerify:             w.stmtVerifyView,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`capability`, err, stmt.Name(statement))
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
func (w *CapabilityWrite) process(q *msg.Request) {
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

// add inserts a new capability
func (w *CapabilityWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		inputVal string
		res      sql.Result
		err      error
	)

	// input validation: MonitoringID
	if w.stmtVerifyMonitoring.QueryRow(
		q.Capability.MonitoringID,
	).Scan(
		&inputVal,
	); err == sql.ErrNoRows {
		mr.NotFound(fmt.Errorf(
			"Monitoring system with ID %s is not registered",
			q.Capability.MonitoringID),
			q.Section,
		)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	// input validation: metric
	if w.stmtVerifyMetric.QueryRow(
		q.Capability.Metric,
	).Scan(
		&inputVal,
	); err == sql.ErrNoRows {
		mr.NotFound(fmt.Errorf("Metric %s is not registered",
			q.Capability.Metric),
			q.Section,
		)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	// input validation: view
	if w.stmtVerifyView.QueryRow(
		q.Capability.View,
	).Scan(
		&inputVal,
	); err == sql.ErrNoRows {
		mr.NotFound(fmt.Errorf("View %s is not registered",
			q.Capability.View),
			q.Section,
		)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	q.Capability.ID = uuid.Must(uuid.NewV4()).String()
	if res, err = w.stmtAdd.Exec(
		q.Capability.ID,
		q.Capability.MonitoringID,
		q.Capability.Metric,
		q.Capability.View,
		q.Capability.Thresholds,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Capability = append(mr.Capability, q.Capability)
	}
}

// remove deletes a capability
func (w *CapabilityWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
	)

	if res, err = w.stmtRemove.Exec(
		q.Capability.ID,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Capability = append(mr.Capability, q.Capability)
	}
}

// ShutdownNow signals the handler to shut down
func (w *CapabilityWrite) ShutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
