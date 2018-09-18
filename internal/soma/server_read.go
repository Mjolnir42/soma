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

// ServerRead handles read requests for server
type ServerRead struct {
	Input       chan msg.Request
	Shutdown    chan struct{}
	handlerName string
	conn        *sql.DB
	stmtList    *sql.Stmt
	stmtShow    *sql.Stmt
	stmtSync    *sql.Stmt
	stmtSearch  *sql.Stmt
	appLog      *logrus.Logger
	reqLog      *logrus.Logger
	errLog      *logrus.Logger
}

// newServerRead return a new ServerRead handler with input buffer of length
func newServerRead(length int) (string, *ServerRead) {
	r := &ServerRead{}
	r.handlerName = generateHandlerName() + `_r`
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return r.handlerName, r
}

// Register initializes resources provided by the Soma app
func (r *ServerRead) Register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (r *ServerRead) RegisterRequests(hmap *handler.Map) {
	for _, action := range []string{
		msg.ActionList,
		msg.ActionShow,
		msg.ActionSync,
		msg.ActionSearch,
	} {
		hmap.Request(msg.SectionServer, action, r.handlerName)
	}
}

// Run is the event loop for ServerRead
func (r *ServerRead) Run() {
	var err error

	for statement, prepStmt := range map[string]**sql.Stmt{
		stmt.ListServers:  &r.stmtList,
		stmt.ShowServers:  &r.stmtShow,
		stmt.SyncServers:  &r.stmtSync,
		stmt.SearchServer: &r.stmtSearch,
	} {
		if *prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`servers`, err, stmt.Name(statement))
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

// Intake exposes the Input channel as part of the handler interface
func (r *ServerRead) Intake() chan msg.Request {
	return r.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (r *ServerRead) PriorityIntake() chan msg.Request {
	return r.Intake()
}

// process is the request dispatcher
func (r *ServerRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	logRequest(r.reqLog, q)

	switch q.Action {
	case msg.ActionList:
		r.list(q, &result)
	case msg.ActionShow:
		r.show(q, &result)
	case msg.ActionSync:
		r.sync(q, &result)
	case msg.ActionSearch:
		r.search(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// list returns all servers
func (r *ServerRead) list(q *msg.Request, mr *msg.Result) {
	var (
		serverID, serverName string
		serverAssetID        int
		rows                 *sql.Rows
		err                  error
	)

	if rows, err = r.stmtList.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&serverID,
			&serverName,
			&serverAssetID,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Server = append(mr.Server, proto.Server{
			ID:      serverID,
			Name:    serverName,
			AssetID: uint64(serverAssetID),
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// show returns the details of a specific server
func (r *ServerRead) show(q *msg.Request, mr *msg.Result) {
	var (
		err                         error
		serverID, serverDc          string
		serverDcLoc, serverName     string
		serverAssetID               int
		serverOnline, serverDeleted bool
	)

	if err = r.stmtShow.QueryRow(
		q.Server.ID,
	).Scan(
		&serverID,
		&serverAssetID,
		&serverDc,
		&serverDcLoc,
		&serverName,
		&serverOnline,
		&serverDeleted,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Server = append(mr.Server, proto.Server{
		ID:         serverID,
		AssetID:    uint64(serverAssetID),
		Datacenter: serverDc,
		Location:   serverDcLoc,
		Name:       serverName,
		IsOnline:   serverOnline,
		IsDeleted:  serverDeleted,
	})
	mr.OK()
}

// sync returns details for all servers suitable for sync processing
func (r *ServerRead) sync(q *msg.Request, mr *msg.Result) {
	var (
		err                         error
		serverID, serverDc          string
		serverDcLoc, serverName     string
		serverAssetID               int
		serverOnline, serverDeleted bool
		rows                        *sql.Rows
	)

	if rows, err = r.stmtSync.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&serverID,
			&serverAssetID,
			&serverDc,
			&serverDcLoc,
			&serverName,
			&serverOnline,
			&serverDeleted,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}

		mr.Server = append(mr.Server, proto.Server{
			ID:         serverID,
			AssetID:    uint64(serverAssetID),
			Datacenter: serverDc,
			Location:   serverDcLoc,
			Name:       serverName,
			IsOnline:   serverOnline,
			IsDeleted:  serverDeleted,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// search looks up a server's ID by name or assetID
func (r *ServerRead) search(q *msg.Request, mr *msg.Result) {
	var (
		err                  error
		serverID, serverName string
		serverAssetID        int
		nullName             sql.NullString
		nullAssetID          sql.NullInt64
	)

	if q.Search.Server.Name != `` {
		nullName.String = q.Search.Server.Name
		nullName.Valid = true
	}

	if q.Search.Server.AssetID != 0 {
		nullAssetID.Int64 = int64(q.Search.Server.AssetID)
		nullAssetID.Valid = true
	}

	if err = r.stmtSearch.QueryRow(
		nullName,
		nullAssetID,
	).Scan(
		&serverID,
		&serverName,
		&serverAssetID,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Server = append(mr.Server, proto.Server{
		ID:      serverID,
		Name:    serverName,
		AssetID: uint64(serverAssetID),
	})
	mr.OK()
}

// ShutdownNow signals the handler to shut down
func (r *ServerRead) ShutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
