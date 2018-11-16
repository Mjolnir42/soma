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
	uuid "github.com/satori/go.uuid"
)

// UserWrite handles write requests for views
type UserWrite struct {
	Input       chan msg.Request
	Shutdown    chan struct{}
	handlerName string
	conn        *sql.DB
	stmtAdd     *sql.Stmt
	stmtRemove  *sql.Stmt
	stmtPurge   *sql.Stmt
	stmtUpdate  *sql.Stmt
	appLog      *logrus.Logger
	reqLog      *logrus.Logger
	errLog      *logrus.Logger
	soma        *Soma
}

// newUserWrite return a new UserWrite handler with input buffer of
// length
func newUserWrite(length int, s *Soma) (string, *UserWrite) {
	w := &UserWrite{}
	w.handlerName = generateHandlerName() + `_w`
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	w.soma = s
	return w.handlerName, w
}

// Register initializes resources provided by the Soma app
func (w *UserWrite) Register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (w *UserWrite) RegisterRequests(hmap *handler.Map) {
	for _, action := range []string{
		msg.ActionAdd,
		msg.ActionRemove,
		msg.ActionPurge,
		msg.ActionUpdate,
	} {
		hmap.Request(msg.SectionUserMgmt, action, w.handlerName)
	}
}

// Intake exposes the Input channel as part of the handler interface
func (w *UserWrite) Intake() chan msg.Request {
	return w.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (w *UserWrite) PriorityIntake() chan msg.Request {
	return w.Intake()
}

// Run is the event loop for UserWrite
func (w *UserWrite) Run() {
	var err error

	for statement, prepStmt := range map[string]**sql.Stmt{
		stmt.UserAdd:    &w.stmtAdd,
		stmt.UserPurge:  &w.stmtPurge,
		stmt.UserRemove: &w.stmtRemove,
		stmt.UserUpdate: &w.stmtUpdate,
	} {
		if *prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`user`, err, stmt.Name(statement))
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
func (w *UserWrite) process(q *msg.Request) {
	result := msg.FromRequest(q)
	logRequest(w.reqLog, q)

	switch q.Action {
	case msg.ActionAdd:
		w.add(q, &result)
	case msg.ActionUpdate:
		w.update(q, &result)
	case msg.ActionRemove:
		w.remove(q, &result)
	case msg.ActionPurge:
		w.purge(q, &result)
	default:
		result.UnknownRequest(q)
	}

	if result.IsOK() {
		// supervisor must be notified of user change
		go func() {
			super := w.soma.getSupervisor()
			super.Update <- msg.CacheUpdateFromRequest(q)
		}()
	}
	q.Reply <- result
}

// add inserts a new user
func (w *UserWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	q.User.ID = uuid.Must(uuid.NewV4()).String()
	if res, err = w.stmtAdd.Exec(
		q.User.ID,
		q.User.UserName,
		q.User.FirstName,
		q.User.LastName,
		q.User.EmployeeNumber,
		q.User.MailAddress,
		false,
		q.User.IsSystem,
		false,
		q.User.TeamID,
		q.AuthUser,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.User = append(mr.User, q.User)
	}
}

// update refreshes a user's information
func (w *UserWrite) update(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtUpdate.Exec(
		q.Update.User.UserName,
		q.Update.User.FirstName,
		q.Update.User.LastName,
		q.Update.User.EmployeeNumber,
		q.Update.User.MailAddress,
		q.Update.User.IsDeleted,
		q.Update.User.TeamID,
		q.User.ID,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.User = append(mr.User, q.User)
	}
}

// remove marks a user as deleted
func (w *UserWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtRemove.Exec(
		q.User.ID,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.User = append(mr.User, q.User)
	}
}

// purge deletes users marked as deleted from the database
func (w *UserWrite) purge(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtPurge.Exec(
		q.User.ID,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.User = append(mr.User, q.User)
	}
}

// ShutdownNow signals the handler to shut down
func (w *UserWrite) ShutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
