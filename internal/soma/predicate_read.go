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

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/handler"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/proto"
)

// PredicateRead handles read requests for predicates
type PredicateRead struct {
	Input       chan msg.Request
	Shutdown    chan struct{}
	handlerName string
	conn        *sql.DB
	stmtList    *sql.Stmt
	stmtShow    *sql.Stmt
	appLog      *logrus.Logger
	reqLog      *logrus.Logger
	errLog      *logrus.Logger
}

// newPredicateRead return a new PredicateRead handler with input buffer of length
func newPredicateRead(length int) (string, *PredicateRead) {
	r := &PredicateRead{}
	r.handlerName = generateHandlerName() + `_r`
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return r.handlerName, r
}

// Register initializes resources provided by the Soma app
func (r *PredicateRead) Register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (r *PredicateRead) RegisterRequests(hmap *handler.Map) {
	for _, action := range []string{
		msg.ActionList,
		msg.ActionShow,
	} {
		hmap.Request(msg.SectionPredicate, action, r.handlerName)
	}
}

// Intake exposes the Input channel as part of the handler interface
func (r *PredicateRead) Intake() chan msg.Request {
	return r.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (r *PredicateRead) PriorityIntake() chan msg.Request {
	return r.Intake()
}

// Run is the event loop for PredicateRead
func (r *PredicateRead) Run() {
	var err error

	for statement, prepStmt := range map[string]**sql.Stmt{
		stmt.PredicateList: &r.stmtList,
		stmt.PredicateShow: &r.stmtShow,
	} {
		if *prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`predicate`, err, stmt.Name(statement))
		}
		defer (*prepStmt).Close()
	}

runloop:
	for {
		select {
		case <-r.Shutdown:
			break runloop
		case req := <-r.Input:
			go func() {
				r.process(&req)
			}()
		}
	}
}

// process is the request dispatcher
func (r *PredicateRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	logRequest(r.reqLog, q)

	switch q.Action {
	case msg.ActionList:
		r.list(q, &result)
	case msg.ActionShow:
		r.show(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// list returns all predicates
func (r *PredicateRead) list(q *msg.Request, mr *msg.Result) {
	var (
		predicate string
		rows      *sql.Rows
		err       error
	)

	if rows, err = r.stmtList.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(&predicate); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Predicate = append(mr.Predicate, proto.Predicate{
			Symbol: predicate,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// show returns the details of a specific predicate
func (r *PredicateRead) show(q *msg.Request, mr *msg.Result) {
	var (
		predicate string
		err       error
	)

	if err = r.stmtShow.QueryRow(
		q.Predicate.Symbol,
	).Scan(
		&predicate,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Predicate = append(mr.Predicate, proto.Predicate{
		Symbol: predicate,
	})
	mr.OK()
}

// ShutdownNow signals the handler to shut down
func (r *PredicateRead) ShutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
