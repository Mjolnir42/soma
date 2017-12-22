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
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/client9/reopen"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/tree"
	"github.com/mjolnir42/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

// initialLoad loads all repositories on startup
func (f *ForestCustodian) initialLoad() {
	var (
		rows                     *sql.Rows
		err                      error
		repoID, repoName, teamID string
		isActive, isDeleted      bool
	)
	f.appLog.Println(`ForestCustodian loading existing repositories`)

	if rows, err = f.stmtLoad.Query(); err != nil {
		f.errLog.Fatal("fc.initialLoad(), loading error: ", err)
	}
	defer rows.Close()

treeloop:
	for rows.Next() {
		if err = rows.Scan(
			&repoID,
			&repoName,
			&isDeleted,
			&isActive,
			&teamID,
		); err != nil {
			f.errLog.Printf("fc.initialLoad, scan error: %s",
				err.Error())
			err = nil
			continue treeloop
		}

		f.appLog.Printf("ForestCustodian loading treekeeper: %s",
			repoName)
		if err = f.loadSomaTree(&msg.Request{
			Repository: proto.Repository{
				ID:        repoID,
				Name:      repoName,
				TeamID:    teamID,
				IsDeleted: isDeleted,
				IsActive:  isActive,
			},
		}); err != nil {
			f.errLog.Printf("fc.loadSomaTree(), error: %s",
				err.Error(),
			)
			err = nil
		}
	}
	if err = rows.Err(); err != nil {
		f.errLog.Printf("fc.initialLoad(), error: %s", err.Error())
	}
}

// loadSomaTree loads an existing repository
func (f *ForestCustodian) loadSomaTree(q *msg.Request) error {
	var err error

	actionChan := make(chan *tree.Action, 1024000)
	errChan := make(chan *tree.Error, 1024000)

	sTree := tree.New(tree.Spec{
		ID:     uuid.NewV4().String(),
		Name:   fmt.Sprintf("root_%s", q.Repository.Name),
		Action: actionChan,
		Log:    f.appLog,
	})
	sTree.SetError(errChan)
	tree.NewRepository(tree.RepositorySpec{
		ID:      q.Repository.ID,
		Name:    q.Repository.Name,
		Team:    q.Repository.TeamID,
		Deleted: q.Repository.IsDeleted,
		Active:  q.Repository.IsActive,
	}).Attach(tree.AttachRequest{
		Root:       sTree,
		ParentType: "root",
		ParentID:   sTree.GetID(),
	})
	// errors during tree loading are not a good sign
	for i := len(errChan); i > 0; i-- {
		e := <-errChan
		// log all errors
		f.errLog.Errorf(
			"ForestCustodian/%s: %s",
			q.Repository.Name,
			e.String(),
		)
		err = fmt.Errorf(e.String())
	}
	for i := len(actionChan); i > 0; i-- {
		// discard actions on initial load
		<-actionChan
	}

	// return last error from error channel
	if err != nil {
		return err
	}

	return f.spawnTreeKeeper(
		q,
		sTree,
		errChan,
		actionChan,
		q.Repository.TeamID,
	)
}

func (f *ForestCustodian) spawnTreeKeeper(q *msg.Request, s *tree.Tree,
	ec chan *tree.Error, ac chan *tree.Action, team string) error {

	// only start the single requested repo
	if f.soma.conf.ObserverRepo != `` && q.Repository.Name != f.soma.conf.ObserverRepo {
		return nil
	}
	var (
		err      error
		db       *sql.DB
		lfh, sfh *reopen.FileWriter
	)

	if db, err = f.soma.newDatabaseConn(); err != nil {
		return err
	}

	keeperName := fmt.Sprintf("repository_%s", q.Repository.Name)
	if lfh, err = reopen.NewFileWriter(filepath.Join(
		f.soma.conf.LogPath,
		`repository`,
		fmt.Sprintf("%s.log", keeperName),
	)); err != nil {
		return err
	}
	if sfh, err = reopen.NewFileWriter(filepath.Join(
		f.soma.conf.LogPath,
		`repository`,
		fmt.Sprintf("startup_%s.log", keeperName),
	)); err != nil {
		return err
	}
	tK := new(TreeKeeper)
	tK.Input = make(chan msg.Request, 1024)
	tK.Shutdown = make(chan struct{})
	tK.Stop = make(chan struct{})
	tK.conn = db
	tK.tree = s
	tK.errors = ec
	tK.actions = ac
	tK.status.isBroken = false
	tK.status.isReady = false
	tK.status.isFrozen = false
	tK.status.isStopped = false
	tK.status.requiresRebuild = q.Flag.Rebuild
	tK.status.rebuildLevel = q.Flag.RebuildLevel
	tK.meta.repoID = q.Repository.ID
	tK.meta.repoName = q.Repository.Name
	tK.meta.teamID = team
	tK.appLog = f.appLog
	tK.treeLog = logrus.New()
	tK.treeLog.Out = lfh
	tK.startLog = logrus.New()
	tK.startLog.Out = sfh
	// startup logs are not rotated, the logrotate map is
	// just used to keep acccess to the filehandle
	f.soma.logMap.Add(keeperName, lfh)
	f.soma.logMap.Add(
		fmt.Sprintf("startup_%s", keeperName),
		sfh,
	)

	// during rebuild the treekeeper will not run in background
	if tK.status.requiresRebuild {
		tK.Run()
	} else {
		// non-rebuild, register TK and detach
		f.soma.handlerMap.Add(keeperName, tK)
		go tK.Run()
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
