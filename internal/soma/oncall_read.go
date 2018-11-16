/*-
 * Copyright (c) 2016-2018, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2018, 1&1 IONOS SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma // import "github.com/mjolnir42/soma/internal/soma"

import (
	"database/sql"
	"strconv"
	"time"

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
	stmtMembers *sql.Stmt
	stmtSearch  *sql.Stmt
	stmtShow    *sql.Stmt
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
		msg.ActionMemberList,
		msg.ActionSearch,
		msg.ActionShow,
	} {
		hmap.Request(msg.SectionOncall, action, r.handlerName)
	}
}

// Intake exposes the Input channel as part of the handler interface
func (r *OncallRead) Intake() chan msg.Request {
	return r.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (r *OncallRead) PriorityIntake() chan msg.Request {
	return r.Intake()
}

// Run is the event loop for OncallRead
func (r *OncallRead) Run() {
	var err error

	for statement, prepStmt := range map[string]**sql.Stmt{
		stmt.OncallList:       &r.stmtList,
		stmt.OncallMemberList: &r.stmtMembers,
		stmt.OncallSearch:     &r.stmtSearch,
		stmt.OncallShow:       &r.stmtShow,
	} {
		if *prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`oncall`, err, stmt.Name(statement))
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
func (r *OncallRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	logRequest(r.reqLog, q)

	switch q.Action {
	case msg.ActionList:
		r.list(q, &result)
	case msg.ActionSearch:
		r.search(q, &result)
	case msg.ActionShow:
		r.show(q, &result)
	case msg.ActionMemberList:
		r.members(q, &result)
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
		oncallID, oncallName        string
		dictID, dictName, createdBy string
		createdAt                   time.Time
		oncallNumber                int
		err                         error
	)

	if err = r.stmtShow.QueryRow(
		q.Oncall.ID,
	).Scan(
		&oncallID,
		&oncallName,
		&oncallNumber,
		&dictID,
		&dictName,
		&createdBy,
		&createdAt,
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
		Details: &proto.OncallDetails{
			Creation: &proto.DetailsCreation{
				CreatedAt: createdAt.Format(msg.RFC3339Milli),
				CreatedBy: createdBy,
			},
			DictionaryID:   dictID,
			DictionaryName: dictName,
		},
	})
	mr.OK()
}

// members returns the members of an oncall duty
func (r *OncallRead) members(q *msg.Request, mr *msg.Result) {
	var (
		userID, userUID string
		rows            *sql.Rows
		err             error
	)

	if rows, err = r.stmtMembers.Query(
		q.Oncall.ID,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	oncall := proto.Oncall{}
	oncall.ID = q.Oncall.ID
	oncall.Members = &[]proto.OncallMember{}

	for rows.Next() {
		if err = rows.Scan(
			&userID,
			&userUID,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		*oncall.Members = append(*oncall.Members, proto.OncallMember{
			UserID:   userID,
			UserName: userUID,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Oncall = append(mr.Oncall, oncall)
	mr.OK()
}

// ShutdownNow signals the handler to shut down
func (r *OncallRead) ShutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
