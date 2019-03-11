/*-
 * Copyright (c) 2019, Jörg Pernfuß
 * Copyright (c) 2019, 1&1 IONOS SE
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
	"github.com/mjolnir42/soma/lib/proto"
)

// AdminRead handles write requests for views
type AdminRead struct {
	Input       chan msg.Request
	Shutdown    chan struct{}
	handlerName string
	conn        *sql.DB
	stmtShow    *sql.Stmt
	appLog      *logrus.Logger
	reqLog      *logrus.Logger
	errLog      *logrus.Logger
}

// newAdminRead return a new AdminRead handler with input buffer of
// length
func newAdminRead(length int) (string, *AdminRead) {
	r := &AdminRead{}
	r.handlerName = generateHandlerName() + `_r`
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return r.handlerName, r
}

// Register initializes resources provided by the Soma app
func (r *AdminRead) Register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (r *AdminRead) RegisterRequests(hmap *handler.Map) {
	for _, action := range []string{
		msg.ActionShow,
	} {
		hmap.Request(msg.SectionAdminMgmt, action, r.handlerName)
	}
}

// Intake exposes the Input channel as part of the handler interface
func (r *AdminRead) Intake() chan msg.Request {
	return r.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (r *AdminRead) PriorityIntake() chan msg.Request {
	return r.Intake()
}

// Run is the event loop for UserWrite
func (r *AdminRead) Run() {
	var err error

	for statement, prepStmt := range map[string]**sql.Stmt{
		stmt.AdminShow: &r.stmtShow,
	} {
		if *prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`admin`, err, stmt.Name(statement))
		}
		defer (*prepStmt).Close()
	}

runloop:
	for {
		select {
		case <-r.Shutdown:
			break runloop
		case req := <-r.Input:
			r.process(&req)
		}
	}
}

// process is the request dispatcher
func (r *AdminRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	logRequest(r.reqLog, q)

	switch q.Action {
	case msg.ActionShow:
		r.show(q, &result)
	default:
		result.UnknownRequest(q)
	}

	q.Reply <- result
}

// show retrieves details about an admin
func (r *AdminRead) show(q *msg.Request, mr *msg.Result) {
	var (
		err                                  error
		adminID, adminName, userID, userName string
	)

	if err = r.stmtShow.QueryRow(
		q.Admin.UserID,
	).Scan(
		&adminID,
		&adminName,
		&userID,
		&userName,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Admin = append(mr.Admin, proto.Admin{
		ID:       adminID,
		Name:     adminName,
		UserID:   userID,
		UserName: userName,
	})
	mr.OK()
}

// ShutdownNow signals the handler to shut down
func (r *AdminRead) ShutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
