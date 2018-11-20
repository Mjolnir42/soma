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
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/handler"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/proto"
)

// RepositoryRead handles read requests for buckets
type RepositoryRead struct {
	Input           chan msg.Request
	Shutdown        chan struct{}
	handlerName     string
	conn            *sql.DB
	stmtList        *sql.Stmt
	stmtPropCustom  *sql.Stmt
	stmtPropOncall  *sql.Stmt
	stmtPropService *sql.Stmt
	stmtPropSystem  *sql.Stmt
	stmtSearch      *sql.Stmt
	stmtShow        *sql.Stmt
	appLog          *logrus.Logger
	reqLog          *logrus.Logger
	errLog          *logrus.Logger
}

// newBucketRead returns a new BucketRead handler with input
// buffer of length
func newRepositoryRead(length int) (string, *RepositoryRead) {
	r := &RepositoryRead{}
	r.handlerName = generateHandlerName() + `_r`
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return r.handlerName, r
}

// Register initializes resources provided by the Soma app
func (r *RepositoryRead) Register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (r *RepositoryRead) RegisterRequests(hmap *handler.Map) {
	for _, action := range []string{
		msg.ActionList,
		msg.ActionSearch,
		msg.ActionShow,
	} {
		hmap.Request(msg.SectionRepositoryConfig, action, r.handlerName)
	}
	for _, action := range []string{
		msg.ActionAudit,
		msg.ActionList,
		msg.ActionSearch,
		msg.ActionShow,
	} {
		hmap.Request(msg.SectionRepository, action, r.handlerName)
	}
}

// Intake exposes the Input channel as part of the handler interface
func (r *RepositoryRead) Intake() chan msg.Request {
	return r.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (r *RepositoryRead) PriorityIntake() chan msg.Request {
	return r.Intake()
}

// Run is the event loop for RepositoryRead
func (r *RepositoryRead) Run() {
	var err error

	for statement, prepStmt := range map[string]**sql.Stmt{
		stmt.AuthorizedRepositoryList:   &r.stmtList,
		stmt.AuthorizedRepositorySearch: &r.stmtSearch,
		stmt.AuthorizedRepositoryShow:   &r.stmtShow,
		stmt.RepoCstProps:               &r.stmtPropCustom,
		stmt.RepoOncProps:               &r.stmtPropOncall,
		stmt.RepoSvcProps:               &r.stmtPropService,
		stmt.RepoSysProps:               &r.stmtPropSystem,
	} {
		if *prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`repository`, err, stmt.Name(statement))
		}
		defer (*prepStmt).Close()
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

// process is the event dispatcher for RepositoryRead
func (r *RepositoryRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	logRequest(r.reqLog, q)

	switch q.Action {
	case msg.ActionList:
		r.list(q, &result)
	case msg.ActionSearch:
		r.search(q, &result)
	case msg.ActionShow:
		r.show(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// list returns all repositories
func (r *RepositoryRead) list(q *msg.Request, mr *msg.Result) {
	var (
		repoID, repoName string
		rows             *sql.Rows
		err              error
	)

	if rows, err = r.stmtList.Query(
		q.Section,
		q.Action,
		q.AuthUser,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&repoID,
			&repoName,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Repository = append(mr.Repository, proto.Repository{
			ID:   repoID,
			Name: repoName,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// show returns the details of a specific repository
func (r *RepositoryRead) show(q *msg.Request, mr *msg.Result) {
	var (
		repoID, repoName, teamID, createdBy string
		isDeleted, isActive                 bool
		createdAt                           time.Time
		err                                 error
	)

	if err = r.stmtShow.QueryRow(
		q.Section,
		q.Action,
		q.AuthUser,
		q.Repository.ID,
		q.Repository.TeamID,
	).Scan(
		&repoID,
		&repoName,
		&isDeleted,
		&isActive,
		&teamID,
		&createdBy,
		&createdAt,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	repo := proto.Repository{
		ID:        repoID,
		Name:      repoName,
		TeamID:    teamID,
		IsDeleted: isDeleted,
		IsActive:  isActive,
		Details: &proto.RepositoryDetails{
			Creation: &proto.DetailsCreation{
				CreatedAt: createdAt.UTC().Format(time.RFC3339),
				CreatedBy: createdBy,
			},
		},
	}

	// add properties
	repo.Properties = &[]proto.Property{}

	if err = r.oncallProperties(&repo); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if err = r.serviceProperties(&repo); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if err = r.systemProperties(&repo); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if err = r.customProperties(&repo); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if len(*repo.Properties) == 0 {
		// trigger ,omitempty in JSON export
		repo.Properties = nil
	}

	mr.Repository = append(mr.Repository, repo)
	mr.OK()
}

// search looks up a repository
func (r *RepositoryRead) search(q *msg.Request, mr *msg.Result) {
	var (
		repoID, repoName, teamID, createdBy            string
		rows                                           *sql.Rows
		err                                            error
		searchRepoID, searchRepoName, searchRepoTeamID sql.NullString
		isDeleted, isActive                            bool
		createdAt                                      time.Time
		searchActive, searchDeleted                    sql.NullBool
	)

	if q.Search.Repository.ID != `` {
		searchRepoID.String = q.Search.Repository.ID
		searchRepoID.Valid = true
	}
	if q.Search.Repository.Name != `` {
		searchRepoName.String = q.Search.Repository.Name
		searchRepoName.Valid = true
	}
	if q.Search.Repository.TeamID != `` {
		searchRepoTeamID.String = q.Search.Repository.TeamID
		searchRepoTeamID.Valid = true
	}

	searchDeleted.Bool = q.Search.Repository.IsDeleted
	searchDeleted.Valid = q.Search.Repository.FilterOnIsDeleted

	searchActive.Bool = q.Search.Repository.IsActive
	searchActive.Valid = q.Search.Repository.FilterOnIsActive

	if rows, err = r.stmtSearch.Query(
		q.Section,
		q.Action,
		q.AuthUser,
		searchRepoID,
		searchRepoName,
		searchRepoTeamID,
		searchDeleted,
		searchActive,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&repoID,
			&repoName,
			&isDeleted,
			&isActive,
			&teamID,
			&createdBy,
			&createdAt,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Repository = append(mr.Repository, proto.Repository{
			ID:        repoID,
			Name:      repoName,
			TeamID:    teamID,
			IsDeleted: isDeleted,
			IsActive:  isActive,
			Details: &proto.RepositoryDetails{
				Creation: &proto.DetailsCreation{
					CreatedAt: createdAt.UTC().Format(time.RFC3339),
					CreatedBy: createdBy,
				},
			},
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// ShutdownNow signals the handler to shut down
func (r *RepositoryRead) ShutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
