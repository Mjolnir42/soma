/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/proto"
)

// DeploymentWrite handles requests for generated deployment
// details
type DeploymentWrite struct {
	Input                    chan msg.Request
	Shutdown                 chan struct{}
	conn                     *sql.DB
	stmtGet                  *sql.Stmt
	stmtSetStatusUpdate      *sql.Stmt
	stmtGetStatus            *sql.Stmt
	stmtActivate             *sql.Stmt
	stmtList                 *sql.Stmt
	stmtAll                  *sql.Stmt
	stmtClearFlag            *sql.Stmt
	stmtDeprovision          *sql.Stmt
	stmtDeprovisionForUpdate *sql.Stmt
	appLog                   *logrus.Logger
	reqLog                   *logrus.Logger
	errLog                   *logrus.Logger
}

// newDeploymentWrite return a new DeploymentWrite handler with
// input buffer of length
func newDeploymentWrite(length int) (w *DeploymentWrite) {
	w = &DeploymentWrite{}
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (w *DeploymentWrite) register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// run is the event loop for DeploymentWrite
func (self *DeploymentWrite) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.DeploymentGet:              self.stmtGet,
		stmt.DeploymentUpdate:           self.stmtSetStatusUpdate,
		stmt.DeploymentStatus:           self.stmtGetStatus,
		stmt.DeploymentActivate:         self.stmtActivate,
		stmt.DeploymentList:             self.stmtList,
		stmt.DeploymentListAll:          self.stmtAll,
		stmt.DeploymentClearFlag:        self.stmtClearFlag,
		stmt.DeploymentDeprovision:      self.stmtDeprovision,
		stmt.DeploymentDeprovisionStyle: self.stmtDeprovisionForUpdate,
	} {
		if prepStmt, err = self.conn.Prepare(statement); err != nil {
			self.errLog.Fatal(`deployment`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-self.Shutdown:
			break runloop
		case req := <-self.Input:
			self.process(&req)
		}
	}
}

// process is the request dispatcher
func (self *DeploymentWrite) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(self.reqLog, q)

	switch q.Action {
	case msg.ActionGet:
		self.get(q, &result)
	case msg.ActionSuccess:
		self.success(q, &result)
	case msg.ActionFailed:
		self.failed(q, &result)
	case msg.ActionList:
		self.listPending(q, &result)
	case msg.ActionAll:
		self.listAll(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// get retrieves a single deployment, adds the correct current task to
// the stored deployment and advances the deployment workflow as required
func (self *DeploymentWrite) get(q *msg.Request, mr *msg.Result) {
	var (
		instanceConfigID, status, nextStatus                      string
		newCurrentStatus, details, newNextStatus, deprovisionTask string
		statusUpdateRequired, hasUpdate                           bool
		err                                                       error
		res                                                       sql.Result
	)

	if err = self.stmtGet.QueryRow(
		q.Deployment.ID,
	).Scan(
		&instanceConfigID,
		&status,
		&nextStatus,
		&details,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	depl := proto.Deployment{}
	if err = json.Unmarshal([]byte(details), &depl); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	// returns true if there is a updated version blocked, ie.
	// after this deprovisioning a new version will be rolled out
	// -- statement always returns true or false, never null
	if err = self.stmtDeprovisionForUpdate.QueryRow(
		q.Deployment.ID,
	).Scan(
		&hasUpdate,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	switch hasUpdate {
	case false:
		deprovisionTask = proto.TaskDelete
	default:
		deprovisionTask = proto.TaskDeprovision
	}

	switch status {
	case proto.DeploymentAwaitingRollout:
		newCurrentStatus = proto.DeploymentRolloutInProgress
		newNextStatus = proto.DeploymentActive
		depl.Task = proto.TaskRollout
		statusUpdateRequired = true
	case proto.DeploymentRolloutInProgress:
		depl.Task = proto.TaskRollout
		statusUpdateRequired = false
	case proto.DeploymentActive:
		depl.Task = proto.TaskRollout
		statusUpdateRequired = false
	case proto.DeploymentRolloutFailed:
		newCurrentStatus = proto.DeploymentRolloutInProgress
		newNextStatus = proto.DeploymentActive
		depl.Task = proto.TaskRollout
		statusUpdateRequired = true
	case proto.DeploymentAwaitingDeprovision:
		newCurrentStatus = proto.DeploymentDeprovisionInProgress
		newNextStatus = proto.DeploymentDeprovisioned
		depl.Task = deprovisionTask
		statusUpdateRequired = true
	case proto.DeploymentDeprovisionInProgress:
		depl.Task = deprovisionTask
		statusUpdateRequired = false
	case proto.DeploymentDeprovisionFailed:
		newCurrentStatus = proto.DeploymentDeprovisionInProgress
		newNextStatus = proto.DeploymentDeprovisioned
		depl.Task = deprovisionTask
		statusUpdateRequired = true
	default:
		// the SQL query filters for the above statuses, a different
		// status should never appear
		mr.ServerError(
			fmt.Errorf(
				"Impossible deployment state %s encountered",
				status,
			),
			q.Section,
		)
		return
	}

	if statusUpdateRequired {
		if res, err = self.stmtSetStatusUpdate.Exec(
			newCurrentStatus,
			newNextStatus,
			instanceConfigID,
		); err != nil {
			mr.ServerError(err, q.Section)
			return
		}
	}

	if mr.RowCnt(res.RowsAffected()) {
		mr.Deployment = append(mr.Deployment, depl)
	}
}

// success marks a rollout as successfully completed
func (self *DeploymentWrite) success(q *msg.Request, mr *msg.Result) {
	var (
		instanceConfigID, status, next, task string
		err                                  error
		res                                  sql.Result
	)

	if err = self.stmtGetStatus.QueryRow(
		q.Deployment.ID,
	).Scan(
		&instanceConfigID,
		&status,
		&next,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	switch status {
	case proto.DeploymentRolloutInProgress:
		if res, err = self.stmtActivate.Exec(
			next,
			proto.DeploymentNone,
			time.Now().UTC(),
			instanceConfigID,
		); err != nil {
			mr.ServerError(err, q.Section)
			return
		}

		task = proto.TaskRollout
	case proto.DeploymentDeprovisionInProgress:
		if task, err = self.deprovisionForUpdate(
			q,
		); err != nil {
			mr.ServerError(err, q.Section)
			return
		}

		if res, err = self.stmtDeprovision.Exec(
			next,
			proto.DeploymentNone,
			time.Now().UTC(),
			instanceConfigID,
		); err != nil {
			mr.ServerError(err, q.Section)
			return
		}
	default:
		mr.ServerError(
			fmt.Errorf(
				"Illegal current state (%s)"+
					" for requested state update.",
				status,
			),
			q.Section,
		)
		return
	}

	if mr.RowCnt(res.RowsAffected()) {
		mr.Deployment = append(mr.Deployment, proto.Deployment{
			ID:   q.Deployment.ID,
			Task: task,
		})
	}
}

// failed marks a rollout as failed
func (self *DeploymentWrite) failed(q *msg.Request, mr *msg.Result) {
	var (
		instanceConfigID, status, next, task string
		err                                  error
		res                                  sql.Result
	)

	if err = self.stmtGetStatus.QueryRow(
		q.Deployment.ID,
	).Scan(
		&instanceConfigID,
		&status,
		&next,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	switch status {
	case proto.DeploymentRolloutInProgress:
		if res, err = self.stmtSetStatusUpdate.Exec(
			proto.DeploymentRolloutFailed,
			proto.DeploymentNone,
			instanceConfigID,
		); err != nil {
			mr.ServerError(err, q.Section)
			return
		}
		task = proto.TaskRollout
	case proto.DeploymentDeprovisionInProgress:
		if task, err = self.deprovisionForUpdate(
			q,
		); err != nil {
			mr.ServerError(err, q.Section)
			return
		}

		if res, err = self.stmtSetStatusUpdate.Exec(
			proto.DeploymentDeprovisionFailed,
			proto.DeploymentNone,
			instanceConfigID,
		); err != nil {
			mr.ServerError(err, q.Section)
			return
		}
	default:
		mr.ServerError(
			fmt.Errorf(
				"Illegal current state (%s)"+
					" for requested state update.",
				status,
			),
			q.Section,
		)
		return
	}

	if mr.RowCnt(res.RowsAffected()) {
		mr.Deployment = append(mr.Deployment, proto.Deployment{
			ID:   q.Deployment.ID,
			Task: task,
		})
	}
}

// listPending returns all deployment IDs for a monitoring system that have
// a pending update that has not yet been fetched
func (self *DeploymentWrite) listPending(q *msg.Request, mr *msg.Result) {
	var (
		instanceID string
		err        error
		list       *sql.Rows
	)

	if list, err = self.stmtList.Query(
		q.Monitoring.Id,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	defer list.Close()

	for list.Next() {
		if err = list.Scan(
			&instanceID,
		); err != nil {
			mr.ServerError(err, q.Section)
			return
		}

		mr.Deployment = append(mr.Deployment, proto.Deployment{
			ID: instanceID,
		})

		// XXX BUG requires a manual transaction
		self.stmtClearFlag.Exec(instanceID)
	}
	if err = list.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// listAll returns all deployment IDs for a monitoring system
func (self *DeploymentWrite) listAll(q *msg.Request, mr *msg.Result) {
	var (
		instanceID string
		err        error
		all        *sql.Rows
	)

	if all, err = self.stmtAll.Query(
		q.Deployment.ID,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	defer all.Close()

	for all.Next() {
		if err = all.Scan(
			&instanceID,
		); err != nil {
			mr.ServerError(err, q.Section)
			return
		}

		mr.Deployment = append(mr.Deployment, proto.Deployment{
			ID: instanceID,
		})

		// XXX BUG requires a manual transaction
		self.stmtClearFlag.Exec(instanceID)
	}
	if err = all.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// shutdownNow signals the handler to shut down
func (self *DeploymentWrite) shutdownNow() {
	close(self.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
