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
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/handler"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/proto"
)

// OncallRead handles read requests for oncall
type OncallRead struct {
	Input       chan msg.Request
	Shutdown    chan struct{}
	handlerName string
	conn        *sql.DB
	stmtList    *sql.Stmt
	stmtShow    *sql.Stmt
	stmtSearch  *sql.Stmt
	appLog      *logrus.Logger
	reqLog      *logrus.Logger
	errLog      *logrus.Logger
}

// newOncallRead return a new OncallRead handler with input buffer of length
func newOncallRead(length int) (string, *OncallRead) {
	r := &OncallRead{}
	r.handlerName = generateHandlerName() + `_r`
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return r.handlerName, r
}

// Register initializes resources provided by the Soma app
func (r *OncallRead) Register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (r *OncallRead) RegisterRequests(hmap *handler.Map) {
	for _, action := range []string{
		msg.ActionList,
		msg.ActionShow,
		msg.ActionSearch,
	} {
		hmap.Request(msg.SectionOncall, action, r.handlerName)
	}
}

// Intake exposes the Input channel as part of the handler interface
func (r *OncallRead) Intake() chan msg.Request {
	return r.Input
}

// Run is the event loop for OncallRead
func (r *OncallRead) Run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.OncallList:   r.stmtList,
		stmt.OncallSearch: r.stmtSearch,
		stmt.OncallShow:   r.stmtShow,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`oncall`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
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
func (r *OncallRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(r.reqLog, q)

	switch q.Action {
	case msg.ActionList:
		r.list(q, &result)
	case msg.ActionSearch:
		r.search(q, &result)
	case msg.ActionShow:
		r.show(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// list returns all oncall duties
func (r *OncallRead) list(q *msg.Request, mr *msg.Result) {
	var (
		oncallID, oncallName string
		rows                 *sql.Rows
		err                  error
	)

	if rows, err = r.stmtList.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(&oncallID, &oncallName); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Oncall = append(mr.Oncall, proto.Oncall{
			ID:   oncallID,
			Name: oncallName,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// search returns all oncall duties
func (r *OncallRead) search(q *msg.Request, mr *msg.Result) {
	var (
		oncallID, oncallName string
		rows                 *sql.Rows
		err                  error
	)

	if rows, err = r.stmtSearch.Query(
		q.Search.Oncall.Name,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(&oncallID, &oncallName); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Oncall = append(mr.Oncall, proto.Oncall{
			ID:   oncallID,
			Name: oncallName,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// show returns the details of a specific oncall duty
func (r *OncallRead) show(q *msg.Request, mr *msg.Result) {
	var (
		oncallID, oncallName string
		oncallNumber         int
		err                  error
	)

	if err = r.stmtShow.QueryRow(
		q.Oncall.ID,
	).Scan(
		&oncallID,
		&oncallName,
		&oncallNumber,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Oncall = append(mr.Oncall, proto.Oncall{
		ID:     oncallID,
		Name:   oncallName,
		Number: strconv.Itoa(oncallNumber),
	})
	mr.OK()
}

// ShutdownNow signals the handler to shut down
func (r *OncallRead) ShutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
