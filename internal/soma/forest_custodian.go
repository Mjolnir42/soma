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
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/internal/tree"
	uuid "github.com/satori/go.uuid"
)

// ForestCustodian handles requests for new repositories by creating
// additional TreeKeeper instances
type ForestCustodian struct {
	Input               chan msg.Request
	System              chan msg.Request
	Shutdown            chan struct{}
	conn                *sql.DB
	stmtAdd             *sql.Stmt
	stmtLoad            *sql.Stmt
	stmtRepoName        *sql.Stmt
	stmtRebuildCheck    *sql.Stmt
	stmtRebuildInstance *sql.Stmt
	appLog              *logrus.Logger
	reqLog              *logrus.Logger
	errLog              *logrus.Logger
	soma                *Soma
}

// newForestCustodian returns a new ForestCustodian handler
// with input buffer of length
func newForestCustodian(length int, s *Soma) (f *ForestCustodian) {
	f = &ForestCustodian{}
	f.Input = make(chan msg.Request, length)
	f.System = make(chan msg.Request, length)
	f.Shutdown = make(chan struct{})
	f.soma = s
	return
}

// register initializes resources provided by the Soma app
func (f *ForestCustodian) register(c *sql.DB, l ...*logrus.Logger) {
	f.conn = c
	f.appLog = l[0]
	f.reqLog = l[1]
	f.errLog = l[2]
}

// run is the event loop for ForestCustodian
func (f *ForestCustodian) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ForestAddRepository:          f.stmtAdd,
		stmt.ForestLoadRepository:         f.stmtLoad,
		stmt.ForestRepoNameById:           f.stmtRepoName,
		stmt.ForestRebuildDeleteChecks:    f.stmtRebuildCheck,
		stmt.ForestRebuildDeleteInstances: f.stmtRebuildInstance,
	} {
		if prepStmt, err = f.conn.Prepare(statement); err != nil {
			f.errLog.Fatal(`forestcustodian`, err,
				stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

	f.initialLoad()

	if f.soma.conf.Observer {
		f.appLog.Println(`ForestCustodian entered observer mode`)
		for {
			select {
			case <-f.Shutdown:
				goto exit
			case req := <-f.System:
				f.sysProcess(&req)
			}
		}
	}

runloop:
	for {
		select {
		case <-f.Shutdown:
			break runloop
		case req := <-f.Input:
			f.process(&req)
		case req := <-f.System:
			f.sysProcess(&req)
		}
	}
exit:
}

// process is the request dispatcher
func (f *ForestCustodian) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(f.reqLog, q)

	switch q.Action {
	case msg.ActionCreate:
		f.create(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// create spawns a new repository tree
func (f *ForestCustodian) create(q *msg.Request, mr *msg.Result) {
	var (
		res        sql.Result
		err        error
		sTree      *tree.Tree
		actionChan chan *tree.Action
		errChan    chan *tree.Error
	)
	actionChan = make(chan *tree.Action, 1024000)
	errChan = make(chan *tree.Error, 1024000)

	if q.Repository.TeamId == `` {
		mr.BadRequest(
			fmt.Errorf("Team has not been set prior to spawning TreeKeeper for repo: %s", q.Repository.Name),
			q.Section,
		)
		return
	}

	q.Repository.Id = uuid.NewV4().String()

	sTree = tree.New(tree.TreeSpec{
		Id:     uuid.NewV4().String(),
		Name:   fmt.Sprintf("root_%s", q.Repository.Name),
		Action: actionChan,
		Log:    f.appLog,
	})
	sTree.SetError(errChan)

	tree.NewRepository(tree.RepositorySpec{
		Id:      q.Repository.Id,
		Name:    q.Repository.Name,
		Team:    q.Repository.TeamId,
		Deleted: false,
		Active:  q.Repository.IsActive,
	}).Attach(tree.AttachRequest{
		Root:       sTree,
		ParentType: "root",
		ParentId:   sTree.GetID(),
	})

	// there should not be anything on the error channel
	// during tree creation
	for i := len(errChan); i > 0; i-- {
		e := <-errChan
		mr.ServerError(e.Error(), q.Section)
		return
	}

	for i := len(actionChan); i > 0; i-- {
		action := <-actionChan
		switch action.Action {
		case msg.ActionCreate:
			if action.Type == `fault` {
				continue
			}
			if action.Type == msg.SectionRepository {
				if res, err = f.stmtAdd.Exec(
					action.Repository.Id,
					action.Repository.Name,
					action.Repository.IsActive,
					false,
					action.Repository.TeamId,
					q.AuthUser,
				); err != nil {
					mr.ServerError(err, q.Section)
					return
				}
				if !mr.RowCnt(res.RowsAffected()) {
					return
				}
			}
		case `attached`:
			// ignored
		default:
			mr.NotImplemented(
				fmt.Errorf("Unknown requested action: %s",
					action.Action,
				))
			return
		}
	}

	// start the handler routine
	if err = f.spawnTreeKeeper(q, sTree, errChan,
		actionChan, q.Repository.TeamId); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	mr.Repository = append(mr.Repository, q.Repository)
	mr.OK()
}

// shutdownNow signals the handler to shut down
func (f *ForestCustodian) shutdownNow() {
	close(f.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
