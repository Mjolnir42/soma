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

// LevelRead handles read requests for alert levels
type LevelRead struct {
	Input       chan msg.Request
	Shutdown    chan struct{}
	handlerName string
	conn        *sql.DB
	stmtList    *sql.Stmt
	stmtSearch  *sql.Stmt
	stmtShow    *sql.Stmt
	appLog      *logrus.Logger
	reqLog      *logrus.Logger
	errLog      *logrus.Logger
}

// newLevelRead return a new LevelRead handler with input buffer of
// length
func newLevelRead(length int) (string, *LevelRead) {
	r := &LevelRead{}
	r.handlerName = generateHandlerName() + `_r`
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return r.handlerName, r
}

// Register initializes resources provided by the Soma app
func (r *LevelRead) Register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (r *LevelRead) RegisterRequests(hmap *handler.Map) {
	for _, action := range []string{
		msg.ActionList,
		msg.ActionSearch,
		msg.ActionShow,
	} {
		hmap.Request(msg.SectionLevel, action, r.handlerName)
	}
}

// Intake exposes the Input channel as part of the handler interface
func (r *LevelRead) Intake() chan msg.Request {
	return r.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (r *LevelRead) PriorityIntake() chan msg.Request {
	return r.Intake()
}

// Run is the event loop for LevelRead
func (r *LevelRead) Run() {
	var err error

	for statement, prepStmt := range map[string]**sql.Stmt{
		stmt.LevelList:   &r.stmtList,
		stmt.LevelSearch: &r.stmtSearch,
		stmt.LevelShow:   &r.stmtShow,
	} {
		if *prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`level`, err, stmt.Name(statement))
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
func (r *LevelRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	logRequest(r.reqLog, q)

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

// list returns all alert levels
func (r *LevelRead) list(q *msg.Request, mr *msg.Result) {
	var (
		level, short string
		rows         *sql.Rows
		err          error
	)

	if rows, err = r.stmtList.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(&level, &short); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Level = append(mr.Level, proto.Level{
			Name:      level,
			ShortName: short,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.OK()
}

// show returns the details of a specific alert levels
func (r *LevelRead) show(q *msg.Request, mr *msg.Result) {
	var (
		level, short string
		numeric      uint16
		err          error
	)

	if err = r.stmtShow.QueryRow(
		q.Level.Name,
	).Scan(
		&level,
		&short,
		&numeric,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Level = append(mr.Level, proto.Level{
		Name:      level,
		ShortName: short,
		Numeric:   numeric,
	})
	mr.OK()
}

// search returns a specific alert level
func (r *LevelRead) search(q *msg.Request, mr *msg.Result) {
	var (
		searchName, searchShort sql.NullString
		level, short            string
		rows                    *sql.Rows
		err                     error
	)

	if q.Search.Level.Name != `` {
		searchName.String = q.Search.Level.Name
		searchName.Valid = true
	}

	if q.Search.Level.ShortName != `` {
		searchShort.String = q.Search.Level.ShortName
		searchShort.Valid = true
	}

	if rows, err = r.stmtSearch.Query(
		searchName,
		searchShort,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(&level, &short); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Level = append(mr.Level, proto.Level{
			Name:      level,
			ShortName: short,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.OK()
}

// ShutdownNow signals the handler to shut down
func (r *LevelRead) ShutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
