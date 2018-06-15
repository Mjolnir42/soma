/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma // import "github.com/mjolnir42/soma/internal/soma"

import (
	"database/sql"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/lib/pq"
	"github.com/mjolnir42/soma/internal/handler"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/proto"
)

// WorkflowRead handles read request for workflow progress
// information
type WorkflowRead struct {
	Input       chan msg.Request
	Shutdown    chan struct{}
	handlerName string
	conn        *sql.DB
	stmtSummary *sql.Stmt
	stmtList    *sql.Stmt
	appLog      *logrus.Logger
	reqLog      *logrus.Logger
	errLog      *logrus.Logger
}

// newWorkflowRead return a new WorkflowRead handler with input buffer
// of length
func newWorkflowRead(length int) (string, *WorkflowRead) {
	r := &WorkflowRead{}
	r.handlerName = generateHandlerName() + `_r`
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return r.handlerName, r
}

// Register initializes resources provided by the Soma app
func (r *WorkflowRead) Register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (r *WorkflowRead) RegisterRequests(hmap *handler.Map) {
	for _, action := range []string{
		msg.ActionSummary,
		msg.ActionList,
	} {
		hmap.Request(msg.SectionWorkflow, action, r.handlerName)
	}
}

// Intake exposes the Input channel as part of the handler interface
func (r *WorkflowRead) Intake() chan msg.Request {
	return r.Input
}

// Run is the event loop for WorkflowRead
func (r *WorkflowRead) Run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.WorkflowSummary: r.stmtSummary,
		stmt.WorkflowList:    r.stmtList,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`workflow_r`, err, stmt.Name(statement))
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
func (r *WorkflowRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(r.reqLog, q)

	switch q.Action {
	case msg.ActionSummary:
		r.summary(q, &result)
	case msg.ActionList:
		r.list(q, &result)
	default:
		result.UnknownRequest(q)
		return
	}

	q.Reply <- result
}

// list returns information on all deployments in a specific
// workflow status
func (r *WorkflowRead) list(q *msg.Request, mr *msg.Result) {
	var (
		err                                           error
		status, instanceID, checkID, repoID, configID string
		instanceConfigID                              string
		version                                       int64
		rows                                          *sql.Rows
		activatedNull, deprovisionedNull              pq.NullTime
		updatedNull, notifiedNull                     pq.NullTime
		created                                       time.Time
		isInherited                                   bool
	)

	workflow := proto.Workflow{
		Instances: &[]proto.Instance{},
	}

	if rows, err = r.stmtList.Query(
		q.Workflow.Status,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	for rows.Next() {
		if err = rows.Scan(
			&instanceID,
			&checkID,
			&repoID,
			&configID,
			&instanceConfigID,
			&version,
			&status,
			&created,
			&activatedNull,
			&deprovisionedNull,
			&updatedNull,
			&notifiedNull,
			&isInherited,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		instance := proto.Instance{
			ID:               instanceID,
			CheckID:          checkID,
			RepositoryID:     repoID,
			ConfigID:         configID,
			InstanceConfigID: instanceConfigID,
			Version:          uint64(version),
			CurrentStatus:    status,
			IsInherited:      isInherited,
			Info: &proto.InstanceVersionInfo{
				CreatedAt: created.UTC().Format(msg.RFC3339Milli),
			},
		}
		if activatedNull.Valid {
			instance.Info.ActivatedAt = activatedNull.
				Time.UTC().Format(msg.RFC3339Milli)
		}
		if deprovisionedNull.Valid {
			instance.Info.DeprovisionedAt = deprovisionedNull.
				Time.UTC().Format(msg.RFC3339Milli)
		}
		if updatedNull.Valid {
			instance.Info.StatusLastUpdatedAt = updatedNull.
				Time.UTC().Format(msg.RFC3339Milli)
		}
		if notifiedNull.Valid {
			instance.Info.NotifiedAt = notifiedNull.
				Time.UTC().Format(msg.RFC3339Milli)
		}
		*workflow.Instances = append(*workflow.Instances,
			instance)
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Workflow = append(mr.Workflow, workflow)
	mr.OK()
}

// summary returns counts for the workflow status distribution
func (r *WorkflowRead) summary(q *msg.Request, mr *msg.Result) {
	var (
		err    error
		status string
		count  int64
		rows   *sql.Rows
	)
	summary := proto.WorkflowSummary{}

	if rows, err = r.stmtSummary.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	for rows.Next() {
		if err = rows.Scan(
			&status,
			&count,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}

		switch status {
		case proto.DeploymentAwaitingComputation:
			summary.AwaitingComputation = uint64(count)
		case proto.DeploymentComputed:
			summary.Computed = uint64(count)
		case proto.DeploymentAwaitingRollout:
			summary.AwaitingRollout = uint64(count)
		case proto.DeploymentRolloutInProgress:
			summary.RolloutInProgress = uint64(count)
		case proto.DeploymentRolloutFailed:
			summary.RolloutFailed = uint64(count)
		case proto.DeploymentActive:
			summary.Active = uint64(count)
		case proto.DeploymentAwaitingDeprovision:
			summary.AwaitingDeprovision = uint64(count)
		case proto.DeploymentDeprovisionInProgress:
			summary.DeprovisionInProgress = uint64(count)
		case proto.DeploymentDeprovisionFailed:
			summary.DeprovisionFailed = uint64(count)
		case proto.DeploymentDeprovisioned:
			summary.Deprovisioned = uint64(count)
		case proto.DeploymentAwaitingDeletion:
			summary.AwaitingDeletion = uint64(count)
		case proto.DeploymentBlocked:
			summary.Blocked = uint64(count)
		}
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Workflow = append(mr.Workflow, proto.Workflow{
		Summary: &summary,
	})
	mr.OK()
}

// ShutdownNow signals the handler to shut down
func (r *WorkflowRead) ShutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
