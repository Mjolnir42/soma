/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"database/sql"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/handler"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/proto"
)

// MetricWrite handles write requests for metrics
type MetricWrite struct {
	Input              chan msg.Request
	Shutdown           chan struct{}
	handlerName        string
	conn               *sql.DB
	stmtAdd            *sql.Stmt
	stmtPkgAdd         *sql.Stmt
	stmtPkgRemove      *sql.Stmt
	stmtRemove         *sql.Stmt
	stmtVerifyProvider *sql.Stmt
	stmtVerifyUnit     *sql.Stmt
	appLog             *logrus.Logger
	reqLog             *logrus.Logger
	errLog             *logrus.Logger
}

// newMetricWrite return a new MetricWrite handler with input buffer of
// length
func newMetricWrite(length int) (string, *MetricWrite) {
	w := &MetricWrite{}
	w.handlerName = generateHandlerName() + `_w`
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	return w.handlerName, w
}

// Register initializes resources provided by the Soma app
func (w *MetricWrite) Register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (w *MetricWrite) RegisterRequests(hmap *handler.Map) {
	for _, action := range []string{
		msg.ActionAdd,
		msg.ActionRemove,
	} {
		hmap.Request(msg.SectionMetric, action, w.handlerName)
	}
}

// Intake exposes the Input channel as part of the handler interface
func (w *MetricWrite) Intake() chan msg.Request {
	return w.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (w *MetricWrite) PriorityIntake() chan msg.Request {
	return w.Intake()
}

// Run is the event loop for MetricWrite
func (w *MetricWrite) Run() {
	var err error

	for statement, prepStmt := range map[string]**sql.Stmt{
		stmt.MetricAdd:      &w.stmtAdd,
		stmt.MetricDel:      &w.stmtRemove,
		stmt.MetricPkgAdd:   &w.stmtPkgAdd,
		stmt.MetricPkgDel:   &w.stmtPkgRemove,
		stmt.ProviderVerify: &w.stmtVerifyProvider,
		stmt.UnitVerify:     &w.stmtVerifyUnit,
	} {
		if *prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`metric`, err, stmt.Name(statement))
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
func (w *MetricWrite) process(q *msg.Request) {
	result := msg.FromRequest(q)
	logRequest(w.reqLog, q)

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

// add inserts a new metric
func (w *MetricWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		res      sql.Result
		err      error
		tx       *sql.Tx
		pkg      proto.MetricPackage
		inputVal string
	)

	// test the referenced unit exists
	if err = w.stmtVerifyUnit.QueryRow(
		q.Metric.Unit,
	).Scan(
		&inputVal,
	); err == sql.ErrNoRows {
		mr.BadRequest(
			fmt.Errorf("Unit %s is not registered",
				q.Metric.Unit),
			q.Section,
		)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	// test the referenced providers exist
	if q.Metric.Packages != nil && *q.Metric.Packages != nil {
		for _, pkg = range *q.Metric.Packages {
			if w.stmtVerifyProvider.QueryRow(
				pkg.Provider,
			).Scan(
				&inputVal,
			); err == sql.ErrNoRows {
				mr.BadRequest(
					fmt.Errorf(
						"Provider %s is not registered",
						pkg.Provider),
					q.Section,
				)
				return
			} else if err != nil {
				mr.ServerError(err, q.Section)
				return
			}
		}
	}

	// start transaction
	if tx, err = w.conn.Begin(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	defer tx.Rollback()

	// insert metric
	if res, err = tx.Stmt(w.stmtAdd).Exec(
		q.Metric.Path,
		q.Metric.Unit,
		q.Metric.Description,
	); err != nil {
		tx.Rollback()
		mr.ServerError(err, q.Section)
		return
	}

	// get row count while still within the transaction
	if !mr.RowCnt(res.RowsAffected()) {
		tx.Rollback()
		return
	}

	// insert all provider package information
	if q.Metric.Packages != nil && *q.Metric.Packages != nil {
		for _, pkg = range *q.Metric.Packages {
			if res, err = tx.Stmt(w.stmtPkgAdd).Exec(
				q.Metric.Path,
				pkg.Provider,
				pkg.Name,
			); err != nil {
				tx.Rollback()
				mr.ServerError(err, q.Section)
				return
			}
			if !mr.RowCnt(res.RowsAffected()) {
				tx.Rollback()
				return
			}
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Metric = append(mr.Metric, q.Metric)
	mr.OK()
}

// remove deletes a metric
func (w *MetricWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
		tx  *sql.Tx
	)

	// start transaction
	if tx, err = w.conn.Begin(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	// delete provider package information for this metric
	if res, err = tx.Stmt(w.stmtPkgRemove).Exec(
		q.Metric.Path,
	); err != nil {
		tx.Rollback()
		mr.ServerError(err, q.Section)
		return
	}

	// delete metric
	if res, err = tx.Stmt(w.stmtRemove).Exec(
		q.Metric.Path,
	); err != nil {
		tx.Rollback()
		mr.ServerError(err, q.Section)
		return
	}

	// get row count while still within the transaction
	if !mr.RowCnt(res.RowsAffected()) {
		tx.Rollback()
		return
	}

	if err = tx.Commit(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Metric = append(mr.Metric, q.Metric)
	mr.OK()
}

// ShutdownNow signals the handler to shut down
func (w *MetricWrite) ShutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
