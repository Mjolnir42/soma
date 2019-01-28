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

// NodeWrite handles write requests for nodes
type NodeWrite struct {
	Input       chan msg.Request
	Shutdown    chan struct{}
	handlerName string
	conn        *sql.DB
	stmtAdd     *sql.Stmt
	stmtPurge   *sql.Stmt
	stmtRemove  *sql.Stmt
	stmtUpdate  *sql.Stmt
	appLog      *logrus.Logger
	reqLog      *logrus.Logger
	errLog      *logrus.Logger
}

// newNodeWrite return a new NodeWrite handler with input buffer of
// length
func newNodeWrite(length int) (string, *NodeWrite) {
	w := &NodeWrite{}
	w.handlerName = generateHandlerName() + `_w`
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	return w.handlerName, w
}

// Register initializes resources provided by the Soma app
func (w *NodeWrite) Register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (w *NodeWrite) RegisterRequests(hmap *handler.Map) {
	for _, action := range []string{
		msg.ActionAdd,
		msg.ActionRemove,
		msg.ActionPurge,
		msg.ActionUpdate,
	} {
		hmap.Request(msg.SectionNodeMgmt, action, w.handlerName)
	}
}

// Intake exposes the Input channel as part of the handler interface
func (w *NodeWrite) Intake() chan msg.Request {
	return w.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (w *NodeWrite) PriorityIntake() chan msg.Request {
	return w.Intake()
}

// Run is the event loop for NodeWrite
func (w *NodeWrite) Run() {
	var err error

	for statement, prepStmt := range map[string]**sql.Stmt{
		stmt.NodeAdd:    &w.stmtAdd,
		stmt.NodeUpdate: &w.stmtUpdate,
		stmt.NodeRemove: &w.stmtRemove,
		stmt.NodePurge:  &w.stmtPurge,
	} {
		if *prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`node`, err, stmt.Name(statement))
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
func (w *NodeWrite) process(q *msg.Request) {
	result := msg.FromRequest(q)
	logRequest(w.reqLog, q)

	switch q.Action {
	case msg.ActionAdd:
		w.add(q, &result)
	case msg.ActionRemove:
		w.remove(q, &result)
	case msg.ActionUpdate:
		w.update(q, &result)
	case msg.ActionPurge:
		w.purge(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// add inserts a new node
func (w *NodeWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	q.Node.ID = uuid.Must(uuid.NewV4()).String()
	if q.Node.ServerID == `` {
		q.Node.ServerID = `00000000-0000-0000-0000-000000000000`
	}
	if res, err = w.stmtAdd.Exec(
		q.Node.ID,
		q.Node.AssetID,
		q.Node.Name,
		q.Node.TeamID,
		q.Node.ServerID,
		q.Node.State,
		q.Node.IsOnline,
		false,
		q.AuthUser,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Node = append(mr.Node, q.Node)
	}
}

// remove delete a node
func (w *NodeWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtRemove.Exec(
		q.Node.ID,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Node = append(mr.Node, q.Node)
	}
}

// update refreshes a node
func (w *NodeWrite) update(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtUpdate.Exec(
		q.Update.Node.AssetID,
		q.Update.Node.Name,
		q.Update.Node.TeamID,
		q.Update.Node.ServerID,
		q.Update.Node.IsOnline,
		q.Update.Node.IsDeleted,
		q.Node.ID,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Node = append(mr.Node, q.Node)
	}
}

// purge removes a node flagged as deleted
func (w *NodeWrite) purge(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtPurge.Exec(
		q.Node.ID,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Node = append(mr.Node, q.Node)
	}
}

// ShutdownNow signals the handler to shut down
func (w *NodeWrite) ShutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
