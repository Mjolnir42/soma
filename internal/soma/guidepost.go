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

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

// GuidePost handles the request routing to the correct
// TreeKeeper instance that will process the request
type GuidePost struct {
	Input                     chan msg.Request
	System                    chan msg.Request
	Shutdown                  chan struct{}
	conn                      *sql.DB
	stmtJobSave               *sql.Stmt
	stmtRepoForBucketID       *sql.Stmt
	stmtRepoNameByID          *sql.Stmt
	stmtNodeDetails           *sql.Stmt
	stmtServiceLookup         *sql.Stmt
	stmtServiceAttributes     *sql.Stmt
	stmtCapabilityThresholds  *sql.Stmt
	stmtCheckDetailsForDelete *sql.Stmt
	stmtBucketForNodeID       *sql.Stmt
	stmtBucketForClusterID    *sql.Stmt
	stmtBucketForGroupID      *sql.Stmt
	appLog                    *logrus.Logger
	reqLog                    *logrus.Logger
	errLog                    *logrus.Logger
	soma                      *Soma
}

// newGuidePost returns a new GuidePost handler
// with input buffer of length
func newGuidePost(length int, s *Soma) (g *GuidePost) {
	g = &GuidePost{}
	g.Input = make(chan msg.Request, length)
	g.System = make(chan msg.Request, length)
	g.Shutdown = make(chan struct{})
	g.soma = s
	return
}

// Register initializes resources provided by the Soma app
func (g *GuidePost) Register(c *sql.DB, l ...*logrus.Logger) {
	g.conn = c
	g.appLog = l[0]
	g.reqLog = l[1]
	g.errLog = l[2]
}

// Intake exposes the Input channel as part of the handler interface
func (g *GuidePost) Intake() chan msg.Request {
	return g.Input
}

// Run is the event loop for GuidePost
func (g *GuidePost) Run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.JobSave:               g.stmtJobSave,
		stmt.RepoByBucketId:        g.stmtRepoForBucketID,
		stmt.NodeDetails:           g.stmtNodeDetails,
		stmt.RepoNameById:          g.stmtRepoNameByID,
		stmt.ServiceLookup:         g.stmtServiceLookup,
		stmt.ServiceAttributes:     g.stmtServiceAttributes,
		stmt.CapabilityThresholds:  g.stmtCapabilityThresholds,
		stmt.CheckDetailsForDelete: g.stmtCheckDetailsForDelete,
		stmt.NodeBucketId:          g.stmtBucketForNodeID,
		stmt.ClusterBucketId:       g.stmtBucketForClusterID,
		stmt.GroupBucketId:         g.stmtBucketForGroupID,
	} {
		if prepStmt, err = g.conn.Prepare(statement); err != nil {
			g.errLog.Fatal(`guidepost`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

	if g.soma.conf.Observer {
		g.appLog.Println(`GuidePost entered observer mode`)
	observerloop:
		for {
			select {
			case <-g.Shutdown:
				goto exit
			case req := <-g.System:
				g.sysprocess(&req)
				continue observerloop
			}
		}
	}

runloop:
	for {
		select {
		case <-g.Shutdown:
			break runloop
		case req := <-g.Input:
			g.process(&req)
		case req := <-g.System:
			g.sysprocess(&req)
		}
	}
exit:
}

// process saves and forwards the Request
func (g *GuidePost) process(q *msg.Request) {
	var (
		res                      sql.Result
		err                      error
		j                        []byte
		repoID, repoName, keeper string
		nf                       bool
		handler                  *TreeKeeper
		rowCnt                   int64
	)
	result := msg.FromRequest(q)

	// to which tree this request must be forwarded
	if repoID, repoName, nf, err = g.extractRouting(q); err != nil {
		goto bailout
	}

	// verify we can process the request
	if nf, err = g.validateRequest(q); err != nil {
		goto bailout
	}

	// fill in required data for the request
	if nf, err = g.fillReqData(q); err != nil {
		goto bailout
	}

	// check we have a treekeeper for that repository
	if nf, err = g.validateKeeper(repoName); err != nil {
		goto bailout
	}
	keeper = fmt.Sprintf("repository_%s", repoName)
	handler = g.soma.handlerMap.Get(keeper).(*TreeKeeper)

	// store job in database
	q.JobID = uuid.NewV4()
	g.appLog.Infof("Saving job %s (%s/%s) for %s",
		q.JobID.String(), q.Section, q.Action,
		q.AuthUser)
	j, _ = json.Marshal(q)
	if res, err = g.stmtJobSave.Exec(
		q.JobID.String(),
		`queued`,
		`pending`,
		q.Action,
		repoID,
		q.AuthUser,
		string(j),
	); err != nil {
		goto bailout
	}
	// insert can have 0 rows affected if the where clause could
	// not find the user
	rowCnt, _ = res.RowsAffected()
	if rowCnt == 0 {
		err = fmt.Errorf("No rows affected while saving job for user %s",
			q.AuthUser)
		nf = false
		goto bailout
	}

	handler.Input <- *q
	result.JobId = q.JobID.String()

	switch q.Section {
	case msg.SectionRepository:
		result.Repository = append(result.Repository,
			q.Repository)
	case msg.SectionBucket:
		result.Bucket = append(result.Bucket,
			q.Bucket)
	case msg.SectionGroup:
		result.Group = append(result.Group,
			q.Group)
	case msg.SectionCluster:
		result.Cluster = append(result.Cluster,
			q.Cluster)
	case msg.SectionNodeConfig:
		result.Node = append(result.Node,
			q.Node)
	case msg.SectionCheckConfig:
		result.CheckConfig = append(result.CheckConfig,
			q.CheckConfig)
	}

bailout:
	if err != nil {
		if nf {
			result.NotFound(err, q.Section)
		} else {
			result.ServerError(err, q.Section)
		}
	}
	q.Reply <- result
}

// sysprocess handles admin actions
func (g *GuidePost) sysprocess(q *msg.Request) {
	var (
		repoName, repoID, keeper string
		err                      error
		handler                  *TreeKeeper
	)
	result := msg.FromRequest(q)
	result.System = []proto.SystemOperation{q.System}

	switch q.System.Request {
	case msg.ActionRepoStop:
		repoID = q.System.RepositoryId
	default:
		result.UnknownRequest(q)
		goto exit
	}

	if err = g.stmtRepoNameByID.QueryRow(
		repoID,
	).Scan(
		&repoName,
	); err == sql.ErrNoRows {
		result.NotFound(fmt.Errorf(`No such repository`))
		goto exit
	} else if err != nil {
		result.ServerError(err)
		goto exit
	}

	// check we have a treekeeper for that repository
	keeper = fmt.Sprintf("repository_%s", repoName)
	if _, ok := g.soma.handlerMap.Get(keeper).(*TreeKeeper); !ok {
		// no handler running, nothing to stop
		result.OK()
		goto exit
	}

	// might already be stopped
	handler = g.soma.handlerMap.Get(keeper).(*TreeKeeper)
	if handler.isStopped() {
		result.OK()
		goto exit
	}

	// check the treekeeper is ready for system requests
	if !(handler.isReady() || handler.isBroken()) {
		result.Unavailable(
			fmt.Errorf("Repository %s not fully loaded yet.",
				repoName),
		)
		goto exit
	}

	switch q.System.Request {
	case msg.ActionRepoStop:
		if !handler.isStopped() {
			close(handler.Stop)
		}
		result.OK()
	}

exit:
	q.Reply <- result
}

// ShutdownNow signals the handler to shut down
func (g *GuidePost) ShutdownNow() {
	close(g.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
