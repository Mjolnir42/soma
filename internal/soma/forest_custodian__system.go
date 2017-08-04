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
	"time"

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/proto"
)

// sysProcess is the request dispatcher for privileged requests
func (f *ForestCustodian) sysProcess(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(f.reqLog, q)

	switch q.Action {
	case msg.ActionRepoRebuild:
		f.rebuild(q, &result)
	case msg.ActionRepoRestart:
		f.restart(q, &result)
	case msg.ActionRepoStop:
		f.stop(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// stop shuts down a TreeKeeper
func (f *ForestCustodian) stop(q *msg.Request, mr *msg.Result) {
	var (
		repoName, teamID, keeper string
		err                      error
	)

	// look up name of the repository
	if err = f.stmtRepoName.QueryRow(
		q.System.RepositoryId,
	).Scan(
		&repoName,
		&teamID,
	); err == sql.ErrNoRows {
		mr.NotFound(err)
		return
	} else if err != nil {
		mr.ServerError(err)
		return
	}

	// get the treekeeper for the repository
	keeper = fmt.Sprintf("treekeeper_%s", repoName)
	if !f.soma.handlerMap.Exists(keeper) {
		mr.NotFound(fmt.Errorf("Handler %s not found", keeper))
		return
	}

	// stop the handler
	handler := f.soma.handlerMap.Get(keeper).(*TreeKeeper)
	if !handler.isStopped() {
		close(handler.Stop)
	}
	mr.OK()

	// this was a stop request -> done
	if q.Action == msg.ActionRepoStop {
		return
	}

	// store for later in rebuild/restart handlers
	q.Repository = proto.Repository{
		Id:        q.System.RepositoryId,
		Name:      repoName,
		TeamId:    teamID,
		IsDeleted: false,
		IsActive:  true,
	}

	// stop was called for restart or rebuild, give the handler time
	// to drain its channels
	<-time.After(5 * time.Second)

	// remove handler from lookup table
	f.soma.handlerMap.Del(keeper)

	// fully shut down the handler
	close(handler.Shutdown)
}

// restart launches a fresh TreeKeeper for a repository
func (f *ForestCustodian) restart(q *msg.Request, mr *msg.Result) {
	// stop the running TreeKeeper
	if f.stop(q, mr); !mr.IsOK() {
		return
	}
	mr.Code = 0 // reset result code

	// load the tree again
	if err := f.loadSomaTree(q); err != nil {
		mr.ServerError(err)
		return
	}
	mr.OK()
}

// rebuild recreates checks/instances for a repository
func (f *ForestCustodian) rebuild(q *msg.Request, mr *msg.Result) {
	var err error

	// rebuilds are not allowed in observer mode
	if f.soma.conf.Observer {
		mr.Forbidden(fmt.Errorf(`Attempted rebuild in observer mode`))
		return
	}

	// stop the running TreeKeeper
	if f.stop(q, mr); !mr.IsOK() {
		return
	}
	mr.Code = 0 // reset result code

	// mark all existing check instances as deleted - instances
	// are deleted for both rebuild levels checks and instances
	if _, err = f.stmtRebuildInstance.Exec(
		q.System.RepositoryId,
	); err != nil {
		mr.ServerError(err)
		return
	}

	// only delete checks for rebuild level checks
	if q.System.RebuildLevel == `checks` {
		if _, err = f.stmtRebuildCheck.Exec(
			q.System.RepositoryId,
		); err != nil {
			mr.ServerError(err)
			return
		}
	}

	// load the tree again, with requested rebuild active
	q.Flag = msg.Flags{
		Rebuild:      true,
		RebuildLevel: q.System.RebuildLevel,
	}
	if err = f.loadSomaTree(q); err != nil {
		mr.ServerError(err)
		return
	}

	// rebuild has finished, restart the tree. If the rebuild did not
	// work, this will simply be a broken tree once more
	q.Flag = msg.Flags{}
	if err := f.loadSomaTree(q); err != nil {
		mr.ServerError(err)
		return
	}
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
