/*-
 * Copyright (c) 2018, Jörg Pernfuß
 * Copyright (c) 2018, 1&1 IONOS SE
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

// JobResultWrite handles write requests for object entities
type JobResultWrite struct {
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

// newJobResultWrite return a new JobResultWrite handler with
// input buffer of length
func newJobResultWrite(length int) (string, *JobResultWrite) {
	w := &JobResultWrite{}
	w.handlerName = generateHandlerName() + `_w`
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	return w.handlerName, w
}

// Register initializes resources provided by the Soma app
func (w *JobResultWrite) Register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (w *JobResultWrite) RegisterRequests(hmap *handler.Map) {
	for _, action := range []string{
		msg.ActionAdd,
		msg.ActionRemove,
	} {
		hmap.Request(msg.SectionJobResultMgmt, action, w.handlerName)
	}
}

// Intake exposes the Input channel as part of the handler interface
func (w *JobResultWrite) Intake() chan msg.Request {
	return w.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (w *JobResultWrite) PriorityIntake() chan msg.Request {
	return w.Intake()
}

// Run is the event loop for JobResultWrite
func (w *JobResultWrite) Run() {
	var err error

	for statement, prepStmt := range map[string]**sql.Stmt{
		stmt.JobResultMgmtAdd:    &w.stmtAdd,
		stmt.JobResultMgmtRemove: &w.stmtRemove,
	} {
		if *prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`jobType`, err, stmt.Name(statement))
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
func (w *JobResultWrite) process(q *msg.Request) {
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

// add inserts a new entity
func (w *JobResultWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	q.JobResult.ID = uuid.Must(uuid.NewV4()).String()
	if res, err = w.stmtAdd.Exec(
		q.JobResult.ID,
		q.JobResult.Name,
		q.AuthUser,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.JobResult = append(mr.JobResult, q.JobResult)
	}
}

// remove deletes an entity
func (w *JobResultWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtRemove.Exec(
		q.JobResult.ID,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.JobResult = append(mr.JobResult, q.JobResult)
	}
}

// ShutdownNow signals the handler to shut down
func (w *JobResultWrite) ShutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
