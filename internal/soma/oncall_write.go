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

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/handler"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

// OncallWrite handles write requests for oncall
type OncallWrite struct {
	Input               chan msg.Request
	Shutdown            chan struct{}
	handlerName         string
	conn                *sql.DB
	stmtAdd             *sql.Stmt
	stmtAssign          *sql.Stmt
	stmtRemove          *sql.Stmt
	stmtUnassign        *sql.Stmt
	stmtUpdate          *sql.Stmt
	stmtGetAllInstances *sql.Stmt
	appLog              *logrus.Logger
	reqLog              *logrus.Logger
	errLog              *logrus.Logger
	handlerMap          *handler.Map
}

// newOncallWrite return a new OncallWrite handler with input buffer of
// length
func newOncallWrite(length int, appHandlerMap *handler.Map) (string, *OncallWrite) {
	w := &OncallWrite{}
	w.handlerName = generateHandlerName() + `_w`
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	w.handlerMap = appHandlerMap
	return w.handlerName, w
}

// Register initializes resources provided by the Soma app
func (w *OncallWrite) Register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (w *OncallWrite) RegisterRequests(hmap *handler.Map) {
	for _, action := range []string{
		msg.ActionAdd,
		msg.ActionMemberAssign,
		msg.ActionMemberUnassign,
		msg.ActionRemove,
		msg.ActionUpdate,
	} {
		hmap.Request(msg.SectionOncall, action, w.handlerName)
	}
}

// Intake exposes the Input channel as part of the handler interface
func (w *OncallWrite) Intake() chan msg.Request {
	return w.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (w *OncallWrite) PriorityIntake() chan msg.Request {
	return w.Intake()
}

// Run is the event loop for OncallWrite
func (w *OncallWrite) Run() {
	var err error

	for statement, prepStmt := range map[string]**sql.Stmt{
		stmt.OncallAdd:             &w.stmtAdd,
		stmt.OncallMemberAssign:    &w.stmtAssign,
		stmt.OncallMemberUnassign:  &w.stmtUnassign,
		stmt.OncallRemove:          &w.stmtRemove,
		stmt.OncallUpdate:          &w.stmtUpdate,
		stmt.OncallGetAllInstances: &w.stmtGetAllInstances,
	} {
		if *prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`oncall`, err, stmt.Name(statement))
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
func (w *OncallWrite) process(q *msg.Request) {
	result := msg.FromRequest(q)
	logRequest(w.reqLog, q)

	switch q.Action {
	case msg.ActionAdd:
		w.add(q, &result)
	case msg.ActionMemberAssign:
		w.assign(q, &result)
	case msg.ActionMemberUnassign:
		w.unassign(q, &result)
	case msg.ActionRemove:
		w.remove(q, &result)
	case msg.ActionUpdate:
		w.update(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// add inserts a new oncall
func (w *OncallWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
	)

	q.Oncall.ID = uuid.Must(uuid.NewV4()).String()
	if res, err = w.stmtAdd.Exec(
		q.Oncall.ID,
		q.Oncall.Name,
		q.Oncall.Number,
		q.AuthUser,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Oncall = append(mr.Oncall, q.Oncall)
	}
}

// assign adds a user to an oncall duty team
func (w *OncallWrite) assign(q *msg.Request, mr *msg.Result) {
	var res sql.Result
	var err error

	if res, err = w.stmtAssign.Exec(
		q.Oncall.ID,
		(*q.Oncall.Members)[0].UserID,
		q.AuthUser,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Oncall = append(mr.Oncall, q.Oncall)
	}
}

// unassign removes a user from an oncall duty team
func (w *OncallWrite) unassign(q *msg.Request, mr *msg.Result) {
	var res sql.Result
	var err error

	if res, err = w.stmtUnassign.Exec(
		q.Oncall.ID,
		(*q.Oncall.Members)[0].UserID,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Oncall = append(mr.Oncall, q.Oncall)
	}
}

// remove removes an oncall entry
func (w *OncallWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
	)

	if res, err = w.stmtRemove.Exec(
		q.Oncall.ID,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Oncall = append(mr.Oncall, q.Oncall)
	}
}

// update refreshes an oncall entry
func (w *OncallWrite) update(q *msg.Request, mr *msg.Result) {
	var (
		name                                                                   sql.NullString
		number                                                                 sql.NullInt64
		res                                                                    sql.Result
		rows                                                                   *sql.Rows
		n                                                                      int // ensure err not redeclared in if block
		err                                                                    error
		repositoryID, bucketID, entityID, entityType, propertyInstanceID, view string
		inheritance                                                            bool
	)

	// our update statement uses NULL to check which of the values
	// should be updated - can be both
	if q.Update.Oncall.Name != `` {
		name = sql.NullString{String: q.Update.Oncall.Name, Valid: true}
	} else {
		name = sql.NullString{String: ``, Valid: false}
	}

	if q.Update.Oncall.Number != `` {
		if n, err = strconv.Atoi(q.Update.Oncall.Number); err != nil {
			mr.ServerError(err, q.Section)
			return
		}
		number = sql.NullInt64{Int64: int64(n), Valid: true}
	} else {
		number = sql.NullInt64{Int64: 0, Valid: false}
	}
	if res, err = w.stmtUpdate.Exec(
		name,
		number,
		q.Oncall.ID,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Oncall = append(mr.Oncall, q.Oncall)
	}

	if rows, err = w.stmtGetAllInstances.Query(
		q.Oncall.ID,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(&bucketID, &entityID, &entityType, &repositoryID, &propertyInstanceID, &view, &inheritance); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		//send a request to the treekeeper to update the property within the tree
		w.updateTreekeeper(q.AuthUser, q.Oncall.ID, q.Oncall.Name, q.Oncall.Number, view, repositoryID, bucketID, entityID, entityType, propertyInstanceID, inheritance)
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

}

func (w *OncallWrite) updateTreekeeper(authUser, oncallID, oncallName, oncallNumber, view, repositoryID, bucketID, entityID, entityType, propertyInstanceID string, inheritance bool) error {
	returnChannel := make(chan msg.Result, 1)
	request := msg.Request{}
	request.Reply = returnChannel
	request.AuthUser = authUser
	request.Action = msg.ActionPropertyUpdate
	request.Property.Type = proto.PropertyTypeOncall

	//Create Service object
	prop := proto.Property{
		Type:        proto.PropertyTypeOncall,
		View:        view,
		Inheritance: inheritance,
	}
	prop.Oncall = &proto.PropertyOncall{
		ID:     oncallID,
		Name:   oncallName,
		Number: oncallNumber,
	}
	prop.SourceInstanceID = propertyInstanceID
	//Forge the request
	switch entityType {
	case msg.EntityRepository:
		request.Section = msg.SectionRepositoryConfig
		request.TargetEntity = msg.EntityRepository
		req := proto.NewRepositoryRequest()
		req.Repository.ID = entityID
		req.Repository.Properties = &[]proto.Property{prop}
		request.Repository = req.Repository.Clone()
	case msg.EntityBucket:
		request.Section = msg.SectionBucket
		request.TargetEntity = msg.EntityBucket
		req := proto.NewBucketRequest()
		req.Bucket.ID = entityID
		req.Bucket.RepositoryID = repositoryID
		req.Bucket.Properties = &[]proto.Property{prop}
		request.Bucket = req.Bucket.Clone()
	case msg.EntityGroup:
		request.Section = msg.SectionGroup
		request.TargetEntity = msg.EntityGroup
		req := proto.NewGroupRequest()
		req.Group.ID = entityID
		req.Group.RepositoryID = repositoryID
		req.Group.BucketID = bucketID
		req.Group.Properties = &[]proto.Property{prop}
		request.Group = req.Group.Clone()
	case msg.EntityCluster:
		request.Section = msg.SectionCluster
		request.TargetEntity = msg.EntityCluster
		req := proto.NewClusterRequest()
		req.Cluster.ID = entityID
		req.Cluster.RepositoryID = repositoryID
		req.Cluster.BucketID = bucketID
		req.Cluster.Properties = &[]proto.Property{prop}
		request.Cluster = req.Cluster.Clone()
	case msg.EntityNode:
		request.Section = msg.SectionNode
		request.TargetEntity = msg.EntityNode
		req := proto.NewNodeRequest()
		req.Node.ID = entityID
		req.Node.Properties = &[]proto.Property{prop}
		request.Node = req.Node.Clone()
	}
	w.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	if !result.IsOK() {
		return result.Error
	}
	return nil
}

// ShutdownNow signals the handler to shut down
func (w *OncallWrite) ShutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
