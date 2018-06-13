/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2017, Jörg Pernfuß
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
	"github.com/mjolnir42/soma/lib/proto"
)

// MonitoringRead handles read requests for monitoring systems
type MonitoringRead struct {
	Input            chan msg.Request
	Shutdown         chan struct{}
	conn             *sql.DB
	stmtListAll      *sql.Stmt
	stmtListScoped   *sql.Stmt
	stmtShow         *sql.Stmt
	stmtSearchAll    *sql.Stmt
	stmtSearchScoped *sql.Stmt
	appLog           *logrus.Logger
	reqLog           *logrus.Logger
	errLog           *logrus.Logger
}

// newMonitoringRead return a new MonitoringRead handler with
// input buffer of length
func newMonitoringRead(length int) (r *MonitoringRead) {
	r = &MonitoringRead{}
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return
}

// Register initializes resources provided by the Soma app
func (r *MonitoringRead) Register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// Run is the event loop for MonitoringRead
func (r *MonitoringRead) Run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ListAllMonitoringSystems:      r.stmtListAll,
		stmt.ListScopedMonitoringSystems:   r.stmtListScoped,
		stmt.ShowMonitoringSystem:          r.stmtShow,
		stmt.SearchAllMonitoringSystems:    r.stmtSearchAll,
		stmt.SearchScopedMonitoringSystems: r.stmtSearchScoped,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`monitoring`, err, stmt.Name(statement))
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

// Intake exposes the Input channel as part of the handler interface
func (r *MonitoringRead) Intake() chan msg.Request {
	return r.Input
}

// process is the request dispatcher
func (r *MonitoringRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(r.reqLog, q)

	switch q.Action {
	case msg.ActionList:
		r.listScoped(q, &result)
	case msg.ActionAll:
		r.listAll(q, &result)
	case msg.ActionSearch:
		r.searchScoped(q, &result)
	case msg.ActionSearchAll:
		r.searchAll(q, &result)
	case msg.ActionShow:
		r.show(q, &result)
	default:
		result.UnknownRequest(q)
	}

	q.Reply <- result
}

// listAll returns all monitoring systems
func (r *MonitoringRead) listAll(q *msg.Request, mr *msg.Result) {
	var (
		err            error
		monitoringID   string
		monitoringName string
		rows           *sql.Rows
	)
	if rows, err = r.stmtListAll.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&monitoringID,
			&monitoringName,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Monitoring = append(mr.Monitoring, proto.Monitoring{
			ID:   monitoringID,
			Name: monitoringName,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// listScoped returns all monitoring systems the user has access to
func (r *MonitoringRead) listScoped(q *msg.Request, mr *msg.Result) {
	var (
		err            error
		monitoringID   string
		monitoringName string
		rows           *sql.Rows
	)
	if rows, err = r.stmtListScoped.Query(
		q.AuthUser,
	); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&monitoringID,
			&monitoringName,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Monitoring = append(mr.Monitoring, proto.Monitoring{
			ID:   monitoringID,
			Name: monitoringName,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// show returns details about a specific monitoring system
func (r *MonitoringRead) show(q *msg.Request, mr *msg.Result) {
	var (
		err                      error
		monitoringID, name, mode string
		contact, teamID          string
		callbackNull             sql.NullString
		callback                 string
	)
	if err = r.stmtShow.QueryRow(
		q.Monitoring.ID,
	).Scan(
		&monitoringID,
		&name,
		&mode,
		&contact,
		&teamID,
		&callbackNull,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	if callbackNull.Valid {
		callback = callbackNull.String
	}
	mr.Monitoring = append(mr.Monitoring, proto.Monitoring{
		ID:       monitoringID,
		Name:     name,
		Mode:     mode,
		Contact:  contact,
		TeamID:   teamID,
		Callback: callback,
	})
	mr.OK()
}

// searchAll looks up a monitoring systems ID by name
func (r *MonitoringRead) searchAll(q *msg.Request, mr *msg.Result) {
	var (
		err            error
		monitoringID   string
		monitoringName string
	)
	// search condition has unique constraint
	if err = r.stmtSearchAll.QueryRow(
		q.Search.Monitoring.Name,
	).Scan(
		&monitoringID,
		&monitoringName,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Monitoring = append(mr.Monitoring, proto.Monitoring{
		ID:   monitoringID,
		Name: monitoringName,
	})
	mr.OK()
}

// searchScoped looks up a monitoring systems ID by name, restricted
// to the user's access
func (r *MonitoringRead) searchScoped(q *msg.Request, mr *msg.Result) {
	var (
		err            error
		monitoringID   string
		monitoringName string
	)
	// search condition has unique constraint
	if err = r.stmtSearchScoped.QueryRow(
		q.AuthUser,
		q.Search.Monitoring.Name,
	).Scan(
		&monitoringID,
		&monitoringName,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Monitoring = append(mr.Monitoring, proto.Monitoring{
		ID:   monitoringID,
		Name: monitoringName,
	})
	mr.OK()
}

// ShutdownNow signals the handler to shut down
func (r *MonitoringRead) ShutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
