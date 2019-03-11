/*-
 * Copyright (c) 2016-2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma // import "github.com/mjolnir42/soma/internal/soma"

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/handler"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/proto"
)

// DeploymentWrite handles requests for generated deployment
// details
type DeploymentWrite struct {
	Input                    chan msg.Request
	Shutdown                 chan struct{}
	handlerName              string
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
func newDeploymentWrite(length int) (string, *DeploymentWrite) {
	w := &DeploymentWrite{}
	w.handlerName = generateHandlerName() + `_w`
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	return w.handlerName, w
}

// Register initializes resources provided by the Soma app
func (w *DeploymentWrite) Register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (w *DeploymentWrite) RegisterRequests(hmap *handler.Map) {
	for _, action := range []string{
		msg.ActionShow,
		msg.ActionSuccess,
		msg.ActionFailed,
		msg.ActionList,
		msg.ActionPending,
	} {
		hmap.Request(msg.SectionDeployment, action, w.handlerName)
	}
}

// Intake exposes the Input channel as part of the handler interface
func (w *DeploymentWrite) Intake() chan msg.Request {
	return w.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (w *DeploymentWrite) PriorityIntake() chan msg.Request {
	return w.Intake()
}

// Run is the event loop for DeploymentWrite
func (w *DeploymentWrite) Run() {
	var err error

	for statement, prepStmt := range map[string]**sql.Stmt{
		stmt.DeploymentGet:              &w.stmtGet,
		stmt.DeploymentUpdate:           &w.stmtSetStatusUpdate,
		stmt.DeploymentStatus:           &w.stmtGetStatus,
		stmt.DeploymentActivate:         &w.stmtActivate,
		stmt.DeploymentList:             &w.stmtList,
		stmt.DeploymentListAll:          &w.stmtAll,
		stmt.DeploymentClearFlag:        &w.stmtClearFlag,
		stmt.DeploymentDeprovision:      &w.stmtDeprovision,
		stmt.DeploymentDeprovisionStyle: &w.stmtDeprovisionForUpdate,
	} {
		if *prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`deployment`, err, stmt.Name(statement))
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
func (w *DeploymentWrite) process(q *msg.Request) {
	result := msg.FromRequest(q)
	logRequest(w.reqLog, q)

	switch q.Action {
	case msg.ActionShow:
		w.show(q, &result)
	case msg.ActionSuccess:
		w.success(q, &result)
	case msg.ActionFailed:
		w.failed(q, &result)
	case msg.ActionPending:
		w.pending(q, &result)
	case msg.ActionList:
		w.list(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// show retrieves a single deployment, adds the correct current task to
// the stored deployment and advances the deployment workflow as required
func (w *DeploymentWrite) show(q *msg.Request, mr *msg.Result) {
	var (
		instanceConfigID, status, nextStatus                      string
		newCurrentStatus, details, newNextStatus, deprovisionTask string
		statusUpdateRequired, hasUpdate                           bool
		err                                                       error
		res                                                       sql.Result
	)

	if err = w.stmtGet.QueryRow(
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
	if err = w.stmtDeprovisionForUpdate.QueryRow(
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
		if res, err = w.stmtSetStatusUpdate.Exec(
			newCurrentStatus,
			newNextStatus,
			instanceConfigID,
		); err != nil {
			mr.ServerError(err, q.Section)
			return
		}
		if mr.RowCnt(res.RowsAffected()) {
			mr.Deployment = append(mr.Deployment, depl)
		}
	} else {
		mr.Deployment = append(mr.Deployment, depl)
		mr.OK()
	}
}

// success marks a rollout as successfully completed
func (w *DeploymentWrite) success(q *msg.Request, mr *msg.Result) {
	var (
		instanceConfigID, status, next, task string
		err                                  error
		res                                  sql.Result
	)

	if err = w.stmtGetStatus.QueryRow(
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
		if res, err = w.stmtActivate.Exec(
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
		if task, err = w.deprovisionForUpdate(
			q,
		); err != nil {
			mr.ServerError(err, q.Section)
			return
		}

		if res, err = w.stmtDeprovision.Exec(
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
func (w *DeploymentWrite) failed(q *msg.Request, mr *msg.Result) {
	var (
		instanceConfigID, status, next, task string
		err                                  error
		res                                  sql.Result
	)

	if err = w.stmtGetStatus.QueryRow(
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
		if res, err = w.stmtSetStatusUpdate.Exec(
			proto.DeploymentRolloutFailed,
			proto.DeploymentNone,
			instanceConfigID,
		); err != nil {
			mr.ServerError(err, q.Section)
			return
		}
		task = proto.TaskRollout
	case proto.DeploymentDeprovisionInProgress:
		if task, err = w.deprovisionForUpdate(
			q,
		); err != nil {
			mr.ServerError(err, q.Section)
			return
		}

		if res, err = w.stmtSetStatusUpdate.Exec(
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

// pending returns all deployment IDs for a monitoring system that have
// a pending update that has not yet been fetched
func (w *DeploymentWrite) pending(q *msg.Request, mr *msg.Result) {
	var (
		instanceID string
		err        error
		list       *sql.Rows
	)

	if list, err = w.stmtList.Query(
		q.Monitoring.ID,
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
		w.stmtClearFlag.Exec(instanceID)
	}
	if err = list.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// list returns all deployment IDs for a monitoring system
func (w *DeploymentWrite) list(q *msg.Request, mr *msg.Result) {
	var (
		instanceID string
		err        error
		all        *sql.Rows
	)

	if all, err = w.stmtAll.Query(
		q.Monitoring.ID,
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
		w.stmtClearFlag.Exec(instanceID)
	}
	if err = all.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// ShutdownNow signals the handler to shut down
func (w *DeploymentWrite) ShutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
