/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
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

// WorkflowWrite handles write requests to modify workflows
type WorkflowWrite struct {
	Input                      chan msg.Request
	Shutdown                   chan struct{}
	conn                       *sql.DB
	stmtRetryDeployment        *sql.Stmt
	stmtTriggerAvailableUpdate *sql.Stmt
	stmtSet                    *sql.Stmt
	appLog                     *logrus.Logger
	reqLog                     *logrus.Logger
	errLog                     *logrus.Logger
}

// newWorkflowWrite return a new WorkflowWrite handler with
// input buffer of length
func newWorkflowWrite(length int) (w *WorkflowWrite) {
	w = &WorkflowWrite{}
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (w *WorkflowWrite) register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// run is the event loop for WorkflowWrite
func (w *WorkflowWrite) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.WorkflowRetry:           w.stmtRetryDeployment,
		stmt.WorkflowUpdateAvailable: w.stmtTriggerAvailableUpdate,
		stmt.WorkflowSet:             w.stmtSet,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`workflow_r`, err, stmt.Name(statement))
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
func (w *WorkflowWrite) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(w.reqLog, q)

	switch q.Action {
	case msg.ActionRetry:
		w.retry(q, &result)
	case msg.ActionSet:
		w.set(q, &result)
	default:
		result.UnknownRequest(q)
	}

	q.Reply <- result
}

// retry reschedules a failed deployment task
func (w *WorkflowWrite) retry(q *msg.Request, mr *msg.Result) {
	var (
		err error
		tx  *sql.Tx
		res sql.Result
	)
	txMap := map[string]*sql.Stmt{}
	if tx, err = w.conn.Begin(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for name, statement := range map[string]string{
		`retry`:  stmt.WorkflowRetry,
		`update`: stmt.WorkflowUpdateAvailable,
	} {
		if txMap[name], err = tx.Prepare(statement); err != nil {
			// tx.Rollback() closes open prepared statements
			tx.Rollback()
			mr.ServerError(err, q.Section)
			return
		}
	}

	if res, err = txMap[`retry`].Exec(
		q.Workflow.InstanceId,
	); err != nil {
		tx.Rollback()
		mr.ServerError(err, q.Section)
		return
	}
	if !mr.RowCnt(res.RowsAffected()) {
		tx.Rollback()
		return
	}

	if res, err = txMap[`update`].Exec(
		q.Workflow.InstanceId,
	); err != nil {
		tx.Rollback()
		mr.ServerError(err, q.Section)
		return
	}
	if !mr.RowCnt(res.RowsAffected()) {
		tx.Rollback()
		return
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		mr.ServerError(err, q.Section)
		return
	}
	mr.Workflow = append(mr.Workflow, q.Workflow)
	mr.OK()
}

// set updates the workflow state of a deployment task to a user
// supplied value
func (w *WorkflowWrite) set(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)
	if res, err = w.stmtSet.Exec(
		q.Workflow.InstanceConfigId,
		q.Workflow.Status,
		q.Workflow.NextStatus,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Workflow = append(mr.Workflow, q.Workflow)
		mr.OK()
	}
}

// shutdownNow signals the handler to shut down
func (w *WorkflowWrite) shutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
