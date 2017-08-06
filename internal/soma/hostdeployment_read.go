package soma

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/proto"
)

// HostDeploymentRead handles requests for nodes inquiring about
// their local check instance rollouts
type HostDeploymentRead struct {
	Input                   chan msg.Request
	Shutdown                chan struct{}
	conn                    *sql.DB
	stmtInstancesForNode    *sql.Stmt
	stmtLastInstanceVersion *sql.Stmt
	appLog                  *logrus.Logger
	reqLog                  *logrus.Logger
	errLog                  *logrus.Logger
}

// newHostDeploymentRead return a new HostDeploymentRead handler
// with input buffer of length
func newHostDeploymentRead(length int) (r *HostDeploymentRead) {
	r = &HostDeploymentRead{}
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (r *HostDeploymentRead) register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// run is the event loop for HostDeploymentRead
func (self *HostDeploymentRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.DeploymentInstancesForNode:    self.stmtInstancesForNode,
		stmt.DeploymentLastInstanceVersion: self.stmtLastInstanceVersion,
	} {
		if prepStmt, err = self.conn.Prepare(statement); err != nil {
			self.errLog.Fatal(`hostdeployment`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-self.Shutdown:
			break runloop
		case req := <-self.Input:
			go func() {
				self.process(&req)
			}()
		}
	}
}

// process is the request dispatcher
func (self *HostDeploymentRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(self.reqLog, q)

	switch q.Action {
	case msg.ActionGet:
		self.get(q, &result)
	case msg.ActionAssemble:
		self.assemble(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// get returns all local deployments for a node
func (self *HostDeploymentRead) get(q *msg.Request, mr *msg.Result) {
	var (
		checkInstanceID, deploymentDetails, status string
		idList                                     *sql.Rows
		err                                        error
	)

	if idList, err = self.stmtInstancesForNode.Query(
		q.Node.AssetId,
		q.Monitoring.Id,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	defer idList.Close()

	for idList.Next() {
		if err = idList.Scan(
			&checkInstanceID,
		); err != nil {
			mr.ServerError(err, q.Section)
			return
		}

		if err = self.stmtLastInstanceVersion.QueryRow(
			checkInstanceID,
		).Scan(
			&deploymentDetails,
			&status,
		); err != nil {
			mr.ServerError(err, q.Section)
			return
		}

		depl := proto.Deployment{}
		if err = json.Unmarshal(
			[]byte(deploymentDetails),
			&depl,
		); err != nil {
			mr.ServerError(err, q.Section)
			return
		}

		switch status {
		case proto.DeploymentAwaitingRollout,
			proto.DeploymentRolloutInProgress,
			proto.DeploymentActive,
			proto.DeploymentRolloutFailed:
			depl.Task = proto.TaskRollout
		case proto.DeploymentAwaitingDeprovision,
			proto.DeploymentDeprovisionInProgress,
			proto.DeploymentDeprovisioned,
			proto.DeploymentDeprovisionFailed:
			depl.Task = proto.TaskDeprovision
		default:
			depl.Task = proto.TaskPending
		}

		// remove credentials from the hostapi
	skiploop:
		for i, _ := range depl.Service.Attributes {
			if strings.HasPrefix(
				depl.Service.Attributes[i].Name,
				`credential_`,
			) {
				// remove element from slice
				depl.Service.Attributes = append(
					depl.Service.Attributes[:i],
					depl.Service.Attributes[i+1:]...,
				)
				goto skiploop // reset loop
			}
		}

		mr.Deployment = append(mr.Deployment, depl)
	}
	if err = idList.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// assemble calculates update instructions for a node
func (self *HostDeploymentRead) assemble(q *msg.Request, mr *msg.Result) {
	var (
		checkInstanceID, deploymentDetails, status string
		idList                                     *sql.Rows
		err                                        error
	)

	if idList, err = self.stmtInstancesForNode.Query(
		q.Node.AssetId,
		q.Monitoring.Id,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	defer idList.Close()
	idMap := map[string]bool{}

assembleloop:
	for idList.Next() {
		if err = idList.Scan(
			&checkInstanceID,
		); err != nil {
			mr.ServerError(err, q.Section)
			return
		}
		idMap[checkInstanceID] = true

		if err = self.stmtLastInstanceVersion.QueryRow(
			checkInstanceID,
		).Scan(
			&deploymentDetails,
			&status,
		); err != nil {
			mr.ServerError(err, q.Section)
			return
		}

		depl := proto.Deployment{}
		if err = json.Unmarshal(
			[]byte(deploymentDetails),
			&depl,
		); err != nil {
			mr.ServerError(err, q.Section)
			return
		}

		switch status {
		case proto.DeploymentAwaitingRollout,
			proto.DeploymentRolloutInProgress,
			proto.DeploymentActive,
			proto.DeploymentRolloutFailed,
			proto.DeploymentBlocked:
			depl.Task = proto.TaskRollout
		case proto.DeploymentAwaitingDeprovision,
			proto.DeploymentDeprovisionInProgress,
			proto.DeploymentDeprovisionFailed:
			depl.Task = proto.TaskDeprovision
		default:
			// bump this id to the delete list
			delete(idMap, checkInstanceID)
			continue assembleloop
		}

		// remove credentials from the hostapi
	skiploop:
		for i, _ := range depl.Service.Attributes {
			if strings.HasPrefix(
				depl.Service.Attributes[i].Name,
				`credential_`,
			) {
				// remove element from slice
				depl.Service.Attributes = append(
					depl.Service.Attributes[:i],
					depl.Service.Attributes[i+1:]...,
				)
				goto skiploop // reset loop
			}
		}

		mr.Deployment = append(mr.Deployment, depl)
	}
	if err = idList.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	// assemble delete list
	for _, delId := range q.DeploymentIDs {
		if _, ok := idMap[delId]; !ok {
			// submitted ID is not in the list of IDs for which
			// deployments were retrieved
			mr.HostDeployment = append(
				mr.HostDeployment,
				proto.HostDeployment{
					CheckInstanceId: delId,
					DeleteInstance:  true,
				},
			)
		}
	}
	mr.OK()
}

// shutdownNow signals the handler to shut down
func (self *HostDeploymentRead) shutdownNow() {
	close(self.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
