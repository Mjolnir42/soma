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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/handler"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/internal/tree"
	metrics "github.com/rcrowley/go-metrics"
	uuid "github.com/satori/go.uuid"
)

// Metrics is the map of runtime metric registries
var Metrics = make(map[string]metrics.Registry)

// TreeKeeper handles the repository tree structure
type TreeKeeper struct {
	Input               chan msg.Request
	Shutdown            chan struct{}
	Stop                chan struct{}
	errors              chan *tree.Error
	actions             chan *tree.Action
	conn                *sql.DB
	tree                *tree.Tree
	stmtGetView         *sql.Stmt
	stmtStartJob        *sql.Stmt
	stmtCapMonMetric    *sql.Stmt
	stmtCheck           *sql.Stmt
	stmtCheckConfig     *sql.Stmt
	stmtCheckInstance   *sql.Stmt
	stmtCluster         *sql.Stmt
	stmtClusterCustProp *sql.Stmt
	stmtClusterOncall   *sql.Stmt
	stmtClusterService  *sql.Stmt
	stmtClusterSysProp  *sql.Stmt
	stmtDefaultDC       *sql.Stmt
	stmtDelDuplicate    *sql.Stmt
	stmtGetComputed     *sql.Stmt
	stmtGetPrevious     *sql.Stmt
	stmtGroup           *sql.Stmt
	stmtGroupCustProp   *sql.Stmt
	stmtGroupOncall     *sql.Stmt
	stmtGroupService    *sql.Stmt
	stmtGroupSysProp    *sql.Stmt
	stmtList            *sql.Stmt
	stmtNode            *sql.Stmt
	stmtNodeCustProp    *sql.Stmt
	stmtNodeOncall      *sql.Stmt
	stmtNodeService     *sql.Stmt
	stmtNodeSysProp     *sql.Stmt
	stmtPkgs            *sql.Stmt
	stmtTeam            *sql.Stmt
	stmtThreshold       *sql.Stmt
	stmtUpdate          *sql.Stmt
	appLog              *logrus.Logger
	treeLog             *logrus.Logger
	startLog            *logrus.Logger
	meta                struct {
		repoID   string
		repoName string
		teamID   string
	}
	status struct {
		isBroken        bool
		isReady         bool
		isStopped       bool
		isFrozen        bool
		requiresRebuild bool
		rebuildLevel    string
	}
	soma *Soma
}

// newTreeKeeper returns a new TreeKeeper handler with input buffer
// of length
func newTreeKeeper(length int) (tk *TreeKeeper) {
	tk = &TreeKeeper{}
	tk.Input = make(chan msg.Request, length)
	tk.Shutdown = make(chan struct{})
	tk.Stop = make(chan struct{})
	return
}

// Register is only implemented by TreeKeeper to fulfill the Handler
// interface. It is expected to by started by ForestCustodian which
// will fully initialize it.
func (tk *TreeKeeper) Register(c *sql.DB, l ...*logrus.Logger) {
	// TreeKeeper receives its own db connection
	tk.appLog = l[0]
	// TreeKeeper does not use the global request log
	// TreeKeeper does not use the global error log
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes. It is implemented by TreeKeeper to fulfill the Handler
// interface
func (tk *TreeKeeper) RegisterRequests(hmap *handler.Map) {
}

// Intake exposes the Input channel as part of the handler interface
func (tk *TreeKeeper) Intake() chan msg.Request {
	return tk.Input
}

// Run is the method a treeKeeper executes in its background
// go-routine. It checks and handles the input channels and reacts
// appropriately.
func (tk *TreeKeeper) Run() {
	tk.appLog.Printf(
		"Starting TreeKeeper for Repo %s (%s)",
		tk.meta.repoName,
		tk.meta.repoID,
	)

	// adjust the number of running treekeeper instances
	c := metrics.GetOrRegisterCounter(
		`.treekeeper.count`, Metrics[`soma`])
	c.Inc(1)
	defer c.Dec(1)

	// set the tree to the startup logger and load
	tk.tree.SwitchLogger(tk.startLog)
	tk.startupLoad()

	// reset the tree to the regular logger
	tk.tree.SwitchLogger(tk.treeLog)
	// render the startup logger inert without risking
	// a nilptr dereference later
	tk.startLog = logrus.New()
	tk.startLog.Out = ioutil.Discard

	// close the startup logfile
	func() {
		fh := tk.soma.logMap.Get(
			fmt.Sprintf("startup_repository_%s", tk.meta.repoName),
		)
		if fh == nil {
			return
		}
		tk.soma.logMap.Del(fmt.Sprintf(
			"startup_repository_%s", tk.meta.repoName,
		))
		fh.Close()
	}()

	// deferred close the regular logfile
	defer func() {
		fh := tk.soma.logMap.Get(fmt.Sprintf("repository_%s", tk.meta.repoName))
		if fh == nil {
			return
		}
		tk.soma.logMap.Del(fmt.Sprintf("repository_%s", tk.meta.repoName))
		fh.Close()
	}()

	var err error

	// treekeepers have a dedicated connection pool
	defer tk.conn.Close()

	// if this was a rebuild, simply return if it failed
	if tk.status.isBroken && tk.status.requiresRebuild {
		return
	}

	// rebuild was successful, process events from initial loading
	// then exit. We issue a fake job for this.
	if tk.status.requiresRebuild {
		req := msg.Request{
			Section: `rebuild`,
			Action:  `rebuild`,
			JobID:   uuid.NewV4(),
		}
		tk.process(&req)
		tk.buildDeploymentDetails()
		tk.orderDeploymentDetails()
		tk.conn.Close()
		return
	}

	// there was an error during startupLoad(), the repository is
	// considered broken.
broken:
	if tk.status.isBroken {
		b := metrics.GetOrRegisterCounter(
			`.treekeeper.broken.count`, Metrics[`soma`])
		b.Inc(1)
		defer b.Dec(1)

		tickTack := time.NewTicker(time.Second * 10).C
	hoverloop:
		for {
			select {
			case <-tickTack:
				tk.appLog.Printf(
					"TK[%s]: BROKEN REPOSITORY %s flying"+
						" holding patterns!\n",
					tk.meta.repoName, tk.meta.repoID)

			case <-tk.Shutdown:
				break hoverloop
			case <-tk.Stop:
				tk.stop()
				goto stopsign
			}
		}
		return
	}

	// prepare statements
	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.TreekeeperDeleteDuplicateDetails:          tk.stmtDelDuplicate,
		stmt.TxDeployDetailClusterCustProp:             tk.stmtClusterCustProp,
		stmt.TxDeployDetailClusterSysProp:              tk.stmtClusterSysProp,
		stmt.TxDeployDetailDefaultDatacenter:           tk.stmtDefaultDC,
		stmt.TxDeployDetailNodeCustProp:                tk.stmtNodeCustProp,
		stmt.TxDeployDetailNodeSysProp:                 tk.stmtNodeSysProp,
		stmt.TxDeployDetailsCapabilityMonitoringMetric: tk.stmtCapMonMetric,
		stmt.TxDeployDetailsCheck:                      tk.stmtCheck,
		stmt.TxDeployDetailsCheckConfig:                tk.stmtCheckConfig,
		stmt.TxDeployDetailsCheckConfigThreshold:       tk.stmtThreshold,
		stmt.TxDeployDetailsCheckInstance:              tk.stmtCheckInstance,
		stmt.TxDeployDetailsCluster:                    tk.stmtCluster,
		stmt.TxDeployDetailsClusterOncall:              tk.stmtClusterOncall,
		stmt.TxDeployDetailsClusterService:             tk.stmtClusterService,
		stmt.TxDeployDetailsComputeList:                tk.stmtList,
		stmt.TxDeployDetailsGroup:                      tk.stmtGroup,
		stmt.TxDeployDetailsGroupCustProp:              tk.stmtGroupCustProp,
		stmt.TxDeployDetailsGroupOncall:                tk.stmtGroupOncall,
		stmt.TxDeployDetailsGroupService:               tk.stmtGroupService,
		stmt.TxDeployDetailsGroupSysProp:               tk.stmtGroupSysProp,
		stmt.TxDeployDetailsNode:                       tk.stmtNode,
		stmt.TxDeployDetailsNodeOncall:                 tk.stmtNodeOncall,
		stmt.TxDeployDetailsNodeService:                tk.stmtNodeService,
		stmt.TxDeployDetailsProviders:                  tk.stmtPkgs,
		stmt.TxDeployDetailsTeam:                       tk.stmtTeam,
		stmt.TxDeployDetailsUpdate:                     tk.stmtUpdate,
		stmt.TreekeeperGetComputedDeployments:          tk.stmtGetComputed,
		stmt.TreekeeperGetPreviousDeployment:           tk.stmtGetPrevious,
		stmt.TreekeeperGetViewFromCapability:           tk.stmtGetView,
		stmt.TreekeeperStartJob:                        tk.stmtStartJob,
	} {
		if prepStmt, err = tk.conn.Prepare(statement); err != nil {
			tk.treeLog.Println("Error preparing SQL statement: ", err)
			tk.treeLog.Println("Failed statement: ", statement)
			tk.status.isBroken = true
			goto broken
		}
		defer prepStmt.Close()
	}

	tk.appLog.Printf("TK[%s]: ready for service!\n", tk.meta.repoName)
	tk.status.isReady = true

	// in observer mode, the TreeKeeper does nothing after loading
	// the tree
	if tk.soma.conf.Observer {
		tk.appLog.Printf(
			"TreeKeeper [%s] entered observer mode\n", tk.meta.repoName)

		select {
		case <-tk.Stop:
			tk.stop()
			goto stopsign
		case <-tk.Shutdown:
			goto exit
		}
	}

stopsign:
	if tk.status.isStopped {
		// drain the input channel, it could be currently full and
		// writers blocked on it. Future writers will check
		// isStopped() before writing (and/or remove this tree from
		// the handlerMap)
	drain:
		for i := len(tk.Input); i > 0; i-- {
			<-tk.Input
		}
		if len(tk.Input) > 0 {
			// there were blocked writers on a full buffered channel
			goto drain
		}

		tk.appLog.Printf("TreeKeeper [%s] has stopped", tk.meta.repoName)
		for {
			select {
			case <-tk.Shutdown:
				goto exit
			case <-tk.Stop:
			}
		}
	}
runloop:
	for {
		select {
		case <-tk.Shutdown:
			break runloop
		case <-tk.Stop:
			tk.stop()
			goto stopsign
		case req := <-tk.Input:
			tk.process(&req)
			tk.soma.handlerMap.Get(`job_block`).(*JobBlock).Notify <- req.JobID.String()
			if !tk.status.isFrozen {
				// buildDeploymentDetails and orderDeploymentDetails can
				// both mark the tree as broken if there was an error
				// preparing required SQL statements
				tk.buildDeploymentDetails()
				if tk.status.isBroken {
					goto broken
				}
				tk.orderDeploymentDetails()
				if tk.status.isBroken {
					goto broken
				}
			}
		}
	}
exit:
}

func (tk *TreeKeeper) isReady() bool {
	return tk.status.isReady
}

func (tk *TreeKeeper) isBroken() bool {
	return tk.status.isBroken
}

func (tk *TreeKeeper) stop() {
	tk.status.isStopped = true
	tk.status.isReady = false
	tk.status.isBroken = false
}

func (tk *TreeKeeper) isStopped() bool {
	return tk.status.isStopped
}

func (tk *TreeKeeper) process(q *msg.Request) {
	var (
		err                                   error
		hasErrors, hasJobLog, jobNeverStarted bool
		tx                                    *sql.Tx
		stm                                   map[string]*sql.Stmt
		jobLog                                *logrus.Logger
		lfh                                   *os.File
	)

	if !tk.status.requiresRebuild {
		_, err = tk.stmtStartJob.Exec(q.JobID.String(), time.Now().UTC())
		if err != nil {
			tk.treeLog.Printf("Failed starting job %s: %s\n",
				q.JobID.String(),
				err)
			jobNeverStarted = true
			goto bailout
		}
		tk.appLog.Printf("Processing job: %s\n", q.JobID.String())
	} else {
		tk.appLog.Printf("Processing rebuild job: %s\n", q.JobID.String())
	}
	if lfh, err = os.Create(filepath.Join(
		tk.soma.conf.LogPath,
		`job`,
		fmt.Sprintf("%s_%s_%s.log",
			time.Now().UTC().Format(msg.RFC3339Milli),
			tk.meta.repoName,
			q.JobID.String(),
		),
	)); err != nil {
		tk.treeLog.Printf("Failed opening joblog %s: %s\n",
			q.JobID.String(),
			err)
	}
	defer lfh.Close()
	defer lfh.Sync()
	jobLog = logrus.New()
	jobLog.Out = lfh
	hasJobLog = true

	tk.tree.Begin()

	// q.Action == `rebuild` will fall through switch
	switch q.Action {
	// XXX CONVERT to msg.Request.Section / msg.Request.Action

	//
	// TREE MANIPULATION REQUESTS
	case
		`create_bucket`:
		tk.treeBucket(q)

	case
		`create_group`,
		`delete_group`,
		`reset_group_to_bucket`,
		`add_group_to_group`:
		tk.treeGroup(q)

	case
		`create_cluster`,
		`delete_cluster`,
		`reset_cluster_to_bucket`,
		`add_cluster_to_group`:
		tk.treeCluster(q)

	case
		"assign_node",
		"delete_node",
		"reset_node_to_bucket",
		"add_node_to_group",
		"add_node_to_cluster":
		tk.treeNode(q)

	//
	// PROPERTY MANIPULATION REQUESTS
	case
		`add_system_property_to_repository`,
		`add_system_property_to_bucket`,
		`add_system_property_to_group`,
		`add_system_property_to_cluster`,
		`add_system_property_to_node`,
		`add_service_property_to_repository`,
		`add_service_property_to_bucket`,
		`add_service_property_to_group`,
		`add_service_property_to_cluster`,
		`add_service_property_to_node`,
		`add_oncall_property_to_repository`,
		`add_oncall_property_to_bucket`,
		`add_oncall_property_to_group`,
		`add_oncall_property_to_cluster`,
		`add_oncall_property_to_node`,
		`add_custom_property_to_repository`,
		`add_custom_property_to_bucket`,
		`add_custom_property_to_group`,
		`add_custom_property_to_cluster`,
		`add_custom_property_to_node`:
		tk.addProperty(q)

	case
		`delete_system_property_from_repository`,
		`delete_system_property_from_bucket`,
		`delete_system_property_from_group`,
		`delete_system_property_from_cluster`,
		`delete_system_property_from_node`,
		`delete_service_property_from_repository`,
		`delete_service_property_from_bucket`,
		`delete_service_property_from_group`,
		`delete_service_property_from_cluster`,
		`delete_service_property_from_node`,
		`delete_oncall_property_from_repository`,
		`delete_oncall_property_from_bucket`,
		`delete_oncall_property_from_group`,
		`delete_oncall_property_from_cluster`,
		`delete_oncall_property_from_node`,
		`delete_custom_property_from_repository`,
		`delete_custom_property_from_bucket`,
		`delete_custom_property_from_group`,
		`delete_custom_property_from_cluster`,
		`delete_custom_property_from_node`:
		tk.rmProperty(q)

	//
	// CHECK MANIPULATION REQUESTS
	case
		`add_check_to_repository`,
		`add_check_to_bucket`,
		`add_check_to_group`,
		`add_check_to_cluster`,
		`add_check_to_node`:
		err = tk.addCheck(&q.CheckConfig)

	case
		`remove_check_from_repository`,
		`remove_check_from_bucket`,
		`remove_check_from_group`,
		`remove_check_from_cluster`,
		`remove_check_from_node`:
		err = tk.rmCheck(&q.CheckConfig)
	}

	// check if we accumulated an error in one of the switch cases
	if err != nil {
		goto bailout
	}

	// recalculate check instances
	tk.tree.ComputeCheckInstances()

	// open multi-statement transaction
	if tx, stm, err = tk.startTx(); err != nil {
		goto bailout
	}
	defer tx.Rollback()

	// defer constraint checks
	if _, err = tx.Exec(stmt.TxDeferAllConstraints); err != nil {
		tk.treeLog.Println("Failed to exec: tkStmtDeferAllConstraints")
		goto bailout
	}

	// save the check configuration as part of the transaction before
	// processing the action channel
	if strings.Contains(q.Action, "add_check_to_") {
		if err = tk.txCheckConfig(q.CheckConfig,
			stm); err != nil {
			goto bailout
		}
	}

	// mark the check configuration as deleted
	if strings.HasPrefix(q.Action, `remove_check_from_`) {
		if _, err = tx.Exec(
			stmt.TxMarkCheckConfigDeleted,
			q.CheckConfig.ID,
		); err != nil {
			goto bailout
		}
	}

	// if the error channel has entries, we can fully ignore the
	// action channel
	for i := len(tk.errors); i > 0; i-- {
		e := <-tk.errors
		if hasJobLog {
			b, _ := json.Marshal(e)
			jobLog.Println(string(b))
		}
		hasErrors = true
		if err == nil {
			err = fmt.Errorf(e.Action)
		}
	}
	if hasErrors {
		goto bailout
	}

actionloop:
	for i := len(tk.actions); i > 0; i-- {
		a := <-tk.actions

		// log all actions for the job
		if hasJobLog {
			b, _ := json.Marshal(a)
			jobLog.Println(string(b))
		}

		// only check and check_instance actions are relevant during
		// a rebuild, everything else is ignored. Even some deletes are
		// valid, for example when a property overwrites inheritance of
		// another property, the first will generate deletes.
		// Other deletes should not occur, like node/delete, but will be
		// sorted later. TODO
		if tk.status.requiresRebuild {
			if tk.status.rebuildLevel == `instances` {
				switch a.Action {
				case `check_new`, `check_removed`:
					// ignore only in instance-rebuild mode
					continue actionloop
				}
			}
			switch a.Action {
			case `property_new`, `property_delete`,
				`create`, `update`, `delete`,
				`node_assignment`,
				`member_new`, `member_removed`:
				// ignore in all rebuild modes
				continue actionloop
			}
		}

		switch a.Action {
		case `property_new`, `property_delete`:
			if err = tk.txProperty(a, stm); err != nil {
				break actionloop
			}
		case `check_new`, `check_removed`:
			if err = tk.txCheck(a, stm); err != nil {
				break actionloop
			}
		case `check_instance_create`,
			`check_instance_update`,
			`check_instance_delete`:
			if err = tk.txCheckInstance(a, stm); err != nil {
				break actionloop
			}
		case `create`, `update`, `delete`, `node_assignment`,
			`member_new`, `member_removed`:
			if err = tk.txTree(a, stm, q.AuthUser); err != nil {
				break actionloop
			}
		default:
			err = fmt.Errorf(
				"Unhandled message in action stream: %s/%s",
				a.Type,
				a.Action,
			)
			break actionloop
		}

		switch a.Type {
		case "errorchannel":
			continue actionloop
		}
	}
	if err != nil {
		goto bailout
	}

	if !tk.status.requiresRebuild {
		// mark job as finished
		if _, err = tx.Exec(
			stmt.TxFinishJob,
			q.JobID.String(),
			time.Now().UTC(),
			"success",
			``, // empty error field
		); err != nil {
			goto bailout
		}
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		goto bailout
	}
	tk.appLog.Printf("SUCCESS - Finished job: %s\n", q.JobID.String())

	// accept tree changes
	tk.tree.Commit()

	// update permission cache
	switch q.Section {
	case msg.SectionRepository, msg.SectionBucket, msg.SectionGroup, msg.SectionCluster:
		switch q.Action {
		case msg.ActionCreate, msg.ActionDestroy:
			go func() {
				super := tk.soma.getSupervisor()
				super.Update <- msg.CacheUpdateFromRequest(q)
			}()
		}
	case msg.SectionNodeConfig:
		switch q.Action {
		case msg.ActionAssign, msg.ActionUnassign:
			go func() {
				super := tk.soma.getSupervisor()
				super.Update <- msg.CacheUpdateFromRequest(q)
			}()
		}
	}
	return

bailout:
	tk.appLog.Printf("FAILED - Finished job: %s\n", q.JobID.String())
	tk.treeLog.Printf("Job-Error(%s): %s\n", q.JobID.String(), err)
	if hasJobLog {
		jobLog.Printf("Aborting error: %s\n", err)
	}

	// if this was a rebuild, the tree will not persist and the
	// job is faked. Also if the job never actually started, then it
	// should never be rolled back nor attempted to mark failed.
	if tk.status.requiresRebuild || jobNeverStarted {
		return
	}

	tk.tree.Rollback()
	tx.Rollback()
	tk.conn.Exec(
		stmt.TxFinishJob,
		q.JobID.String(),
		time.Now().UTC(),
		"failed",
		err.Error(),
	)
	for i := len(tk.actions); i > 0; i-- {
		a := <-tk.actions
		jB, _ := json.Marshal(a)
		if hasJobLog {
			jobLog.Printf("Cleaned message: %s\n", string(jB))
		}
	}
	return
}

// ShutdownNow signals the handler to shut down
func (tk *TreeKeeper) ShutdownNow() {
	if !tk.isStopped() {
		tk.stopNow()
	}
	close(tk.Shutdown)
}

func (tk *TreeKeeper) stopNow() {
	close(tk.Stop)
}

func (tk *TreeKeeper) drain(s string) (j int) {
	switch s {
	case "action":
		j = len(tk.actions)
		for i := j; i > 0; i-- {
			<-tk.actions
		}
	case "error":
		j = len(tk.errors)
		for i := j; i > 0; i-- {
			<-tk.errors
		}
	default:
		panic("Requested drain for unknown")
	}
	return j
}

func (tk *TreeKeeper) startTx() (
	*sql.Tx, map[string]*sql.Stmt, error) {

	var err error
	var tx *sql.Tx
	open := false
	stMap := map[string]*sql.Stmt{}

	if tx, err = tk.conn.Begin(); err != nil {
		goto bailout
	}
	open = true

	//
	// PROPERTY STATEMENTS
	for name, statement := range map[string]string{
		`PropertyInstanceCreate`:          stmt.TxPropertyInstanceCreate,
		`PropertyInstanceDelete`:          stmt.TxPropertyInstanceDelete,
		`RepositoryPropertyOncallCreate`:  stmt.TxRepositoryPropertyOncallCreate,
		`RepositoryPropertyOncallDelete`:  stmt.TxRepositoryPropertyOncallDelete,
		`RepositoryPropertyServiceCreate`: stmt.TxRepositoryPropertyServiceCreate,
		`RepositoryPropertyServiceDelete`: stmt.TxRepositoryPropertyServiceDelete,
		`RepositoryPropertySystemCreate`:  stmt.TxRepositoryPropertySystemCreate,
		`RepositoryPropertySystemDelete`:  stmt.TxRepositoryPropertySystemDelete,
		`RepositoryPropertyCustomCreate`:  stmt.TxRepositoryPropertyCustomCreate,
		`RepositoryPropertyCustomDelete`:  stmt.TxRepositoryPropertyCustomDelete,
		`BucketPropertyOncallCreate`:      stmt.TxBucketPropertyOncallCreate,
		`BucketPropertyOncallDelete`:      stmt.TxBucketPropertyOncallDelete,
		`BucketPropertyServiceCreate`:     stmt.TxBucketPropertyServiceCreate,
		`BucketPropertyServiceDelete`:     stmt.TxBucketPropertyServiceDelete,
		`BucketPropertySystemCreate`:      stmt.TxBucketPropertySystemCreate,
		`BucketPropertySystemDelete`:      stmt.TxBucketPropertySystemDelete,
		`BucketPropertyCustomCreate`:      stmt.TxBucketPropertyCustomCreate,
		`BucketPropertyCustomDelete`:      stmt.TxBucketPropertyCustomDelete,
		`GroupPropertyOncallCreate`:       stmt.TxGroupPropertyOncallCreate,
		`GroupPropertyOncallDelete`:       stmt.TxGroupPropertyOncallDelete,
		`GroupPropertyServiceCreate`:      stmt.TxGroupPropertyServiceCreate,
		`GroupPropertyServiceDelete`:      stmt.TxGroupPropertyServiceDelete,
		`GroupPropertySystemCreate`:       stmt.TxGroupPropertySystemCreate,
		`GroupPropertySystemDelete`:       stmt.TxGroupPropertySystemDelete,
		`GroupPropertyCustomCreate`:       stmt.TxGroupPropertyCustomCreate,
		`GroupPropertyCustomDelete`:       stmt.TxGroupPropertyCustomDelete,
		`ClusterPropertyOncallCreate`:     stmt.TxClusterPropertyOncallCreate,
		`ClusterPropertyOncallDelete`:     stmt.TxClusterPropertyOncallDelete,
		`ClusterPropertyServiceCreate`:    stmt.TxClusterPropertyServiceCreate,
		`ClusterPropertyServiceDelete`:    stmt.TxClusterPropertyServiceDelete,
		`ClusterPropertySystemCreate`:     stmt.TxClusterPropertySystemCreate,
		`ClusterPropertySystemDelete`:     stmt.TxClusterPropertySystemDelete,
		`ClusterPropertyCustomCreate`:     stmt.TxClusterPropertyCustomCreate,
		`ClusterPropertyCustomDelete`:     stmt.TxClusterPropertyCustomDelete,
		`NodePropertyOncallCreate`:        stmt.TxNodePropertyOncallCreate,
		`NodePropertyOncallDelete`:        stmt.TxNodePropertyOncallDelete,
		`NodePropertyServiceCreate`:       stmt.TxNodePropertyServiceCreate,
		`NodePropertyServiceDelete`:       stmt.TxNodePropertyServiceDelete,
		`NodePropertySystemCreate`:        stmt.TxNodePropertySystemCreate,
		`NodePropertySystemDelete`:        stmt.TxNodePropertySystemDelete,
		`NodePropertyCustomCreate`:        stmt.TxNodePropertyCustomCreate,
		`NodePropertyCustomDelete`:        stmt.TxNodePropertyCustomDelete,
	} {
		if stMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("tk.Prepare(%s) error: %s",
				name, err.Error())
			delete(stMap, name)
			goto bailout
		}
	}

	//
	// CHECK STATEMENTS
	for name, statement := range map[string]string{
		`CreateCheck`: stmt.TxCreateCheck,
		`DeleteCheck`: stmt.TxMarkCheckDeleted,
	} {
		if stMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("tk.Prepare(%s) error: %s",
				name, err.Error())
			delete(stMap, name)
			goto bailout
		}
	}

	//
	// CHECK INSTANCE STATEMENTS
	for name, statement := range map[string]string{
		`CreateCheckInstance`:              stmt.TxCreateCheckInstance,
		`CreateCheckInstanceConfiguration`: stmt.TxCreateCheckInstanceConfiguration,
		`DeleteCheckInstance`:              stmt.TxMarkCheckInstanceDeleted,
	} {
		if stMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("tk.Prepare(%s) error: %s",
				name, err.Error())
			delete(stMap, name)
			goto bailout
		}
	}

	//
	// CHECK CONFIGURATION STATEMENTS
	for name, statement := range map[string]string{
		`CreateCheckConfigurationBase`:                stmt.TxCreateCheckConfigurationBase,
		`CreateCheckConfigurationThreshold`:           stmt.TxCreateCheckConfigurationThreshold,
		`CreateCheckConfigurationConstraintSystem`:    stmt.TxCreateCheckConfigurationConstraintSystem,
		`CreateCheckConfigurationConstraintNative`:    stmt.TxCreateCheckConfigurationConstraintNative,
		`CreateCheckConfigurationConstraintOncall`:    stmt.TxCreateCheckConfigurationConstraintOncall,
		`CreateCheckConfigurationConstraintCustom`:    stmt.TxCreateCheckConfigurationConstraintCustom,
		`CreateCheckConfigurationConstraintService`:   stmt.TxCreateCheckConfigurationConstraintService,
		`CreateCheckConfigurationConstraintAttribute`: stmt.TxCreateCheckConfigurationConstraintAttribute,
	} {
		if stMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("tk.Prepare(%s) error: %s",
				name, err.Error())
			delete(stMap, name)
			goto bailout
		}
	}

	//
	// TREE MANIPULATION STATEMENTS
	for name, statement := range map[string]string{
		`BucketAssignNode`:         stmt.TxBucketAssignNode,
		`ClusterCreate`:            stmt.TxClusterCreate,
		`ClusterDelete`:            stmt.TxClusterDelete,
		`ClusterMemberNew`:         stmt.TxClusterMemberNew,
		`ClusterMemberRemove`:      stmt.TxClusterMemberRemove,
		`ClusterUpdate`:            stmt.TxClusterUpdate,
		`CreateBucket`:             stmt.TxCreateBucket,
		`GroupCreate`:              stmt.TxGroupCreate,
		`GroupDelete`:              stmt.TxGroupDelete,
		`GroupMemberNewCluster`:    stmt.TxGroupMemberNewCluster,
		`GroupMemberNewGroup`:      stmt.TxGroupMemberNewGroup,
		`GroupMemberNewNode`:       stmt.TxGroupMemberNewNode,
		`GroupMemberRemoveCluster`: stmt.TxGroupMemberRemoveCluster,
		`GroupMemberRemoveGroup`:   stmt.TxGroupMemberRemoveGroup,
		`GroupMemberRemoveNode`:    stmt.TxGroupMemberRemoveNode,
		`GroupUpdate`:              stmt.TxGroupUpdate,
		`NodeUnassignFromBucket`:   stmt.TxNodeUnassignFromBucket,
		`UpdateNodeState`:          stmt.TxUpdateNodeState,
	} {
		if stMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("tk.Prepare(%s) error: %s",
				name, err.Error())
			delete(stMap, name)
			goto bailout
		}
	}

	return tx, stMap, nil

bailout:
	if open {
		// if the transaction was opened, then tx.Rollback() will close all
		// prepared statements. If the transaction was not opened yet, then
		// no statements have been prepared inside it - there is nothing to
		// close
		defer tx.Rollback()
	}
	return nil, nil, err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
