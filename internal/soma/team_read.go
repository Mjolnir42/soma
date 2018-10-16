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

// TeamRead handles read requests for teams
type TeamRead struct {
	Input       chan msg.Request
	Shutdown    chan struct{}
	handlerName string
	conn        *sql.DB
	stmtList    *sql.Stmt
	stmtShow    *sql.Stmt
	stmtSync    *sql.Stmt
	appLog      *logrus.Logger
	reqLog      *logrus.Logger
	errLog      *logrus.Logger
}

// newTeamRead return a new TeamRead handler with input buffer of length
func newTeamRead(length int) (string, *TeamRead) {
	r := &TeamRead{}
	r.handlerName = generateHandlerName() + `_r`
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return r.handlerName, r
}

// Register initializes resources provided by the Soma app
func (r *TeamRead) Register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (r *TeamRead) RegisterRequests(hmap *handler.Map) {
	// Category: identity
	for _, action := range []string{
		msg.ActionList,
		msg.ActionSearch,
		msg.ActionShow,
		msg.ActionSync,
	} {
		hmap.Request(msg.SectionTeamMgmt, action, r.handlerName)
	}
	// Category: self
	for _, action := range []string{
		msg.ActionShow,
		msg.ActionSearch,
	} {
		hmap.Request(msg.SectionTeam, action, r.handlerName)
	}
}

// Intake exposes the Input channel as part of the handler interface
func (r *TeamRead) Intake() chan msg.Request {
	return r.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (r *TeamRead) PriorityIntake() chan msg.Request {
	return r.Intake()
}

// Run is the event loop for TeamRead
func (r *TeamRead) Run() {
	var err error

	for statement, prepStmt := range map[string]**sql.Stmt{
		stmt.ListTeams: &r.stmtList,
		stmt.ShowTeams: &r.stmtShow,
		stmt.SyncTeams: &r.stmtSync,
	} {
		if *prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`team`, err, stmt.Name(statement))
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
func (r *TeamRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	logRequest(r.reqLog, q)

	switch q.Action {
	case msg.ActionList, msg.ActionSearch:
		r.list(q, &result)
	case msg.ActionShow:
		r.show(q, &result)
	case msg.ActionSync:
		r.sync(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// list returns all teams
func (r *TeamRead) list(q *msg.Request, mr *msg.Result) {
	var (
		teamID, teamName string
		rows             *sql.Rows
		err              error
	)

	if rows, err = r.stmtList.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&teamID,
			&teamName,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Team = append(mr.Team, proto.Team{
			ID:   teamID,
			Name: teamName,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// show returns the details for a specific team
func (r *TeamRead) show(q *msg.Request, mr *msg.Result) {
	var (
		teamID, teamName string
		ldapID           int
		systemFlag       bool
		err              error
	)

	if err = r.stmtShow.QueryRow(
		q.Team.ID,
	).Scan(
		&teamID,
		&teamName,
		&ldapID,
		&systemFlag,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Team = append(mr.Team, proto.Team{
		ID:       teamID,
		Name:     teamName,
		LdapID:   strconv.Itoa(ldapID),
		IsSystem: systemFlag,
	})
	mr.OK()
}

// sync returns all teams in a format suitable for sync processing
func (r *TeamRead) sync(q *msg.Request, mr *msg.Result) {
	var (
		teamID, teamName string
		ldapID           int
		systemFlag       bool
		rows             *sql.Rows
		err              error
	)

	if rows, err = r.stmtSync.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&teamID,
			&teamName,
			&ldapID,
			&systemFlag,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Team = append(mr.Team, proto.Team{
			ID:       teamID,
			Name:     teamName,
			LdapID:   strconv.Itoa(ldapID),
			IsSystem: systemFlag,
		})
	}
	mr.OK()
}

// ShutdownNow signals the handler to shut down
func (r *TeamRead) ShutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
