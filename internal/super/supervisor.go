/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"database/sql"
	"encoding/hex"

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/config"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/perm"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/auth"
)

var (
	// only one supervisor instance will be set up by New
	initialized = false
	singleton   *Supervisor
)

// Supervisor handles all AAA requests
type Supervisor struct {
	Input                             chan msg.Request
	Update                            chan msg.Request
	Shutdown                          chan struct{}
	conn                              *sql.DB
	seed                              []byte
	key                               []byte
	readonly                          bool
	tokenExpiry                       uint64
	kexExpiry                         uint64
	credExpiry                        uint64
	activation                        string
	rootDisabled                      bool
	rootRestricted                    bool
	kex                               *svKexMap
	tokens                            *svTokenMap
	credentials                       *svCredMap
	permCache                         *perm.Cache
	stmtTokenSelect                   *sql.Stmt
	stmtFindUserID                    *sql.Stmt
	stmtCheckUserActive               *sql.Stmt
	stmtCategoryList                  *sql.Stmt
	stmtCategoryShow                  *sql.Stmt
	stmtSectionList                   *sql.Stmt
	stmtSectionShow                   *sql.Stmt
	stmtSectionSearch                 *sql.Stmt
	stmtSectionAdd                    *sql.Stmt
	stmtActionList                    *sql.Stmt
	stmtActionShow                    *sql.Stmt
	stmtActionSearch                  *sql.Stmt
	stmtActionAdd                     *sql.Stmt
	stmtRevokeAuthorizationGlobal     *sql.Stmt
	stmtRevokeAuthorizationRepository *sql.Stmt
	stmtRevokeAuthorizationTeam       *sql.Stmt
	stmtRevokeAuthorizationMonitoring *sql.Stmt
	stmtGrantAuthorizationGlobal      *sql.Stmt
	stmtGrantAuthorizationRepository  *sql.Stmt
	stmtGrantAuthorizationTeam        *sql.Stmt
	stmtGrantAuthorizationMonitoring  *sql.Stmt
	stmtSearchAuthorizationGlobal     *sql.Stmt
	stmtSearchAuthorizationRepository *sql.Stmt
	stmtSearchAuthorizationTeam       *sql.Stmt
	stmtSearchAuthorizationMonitoring *sql.Stmt
	stmtPermissionList                *sql.Stmt
	stmtPermissionSearch              *sql.Stmt
	stmtPermissionMapEntry            *sql.Stmt
	stmtPermissionUnmapEntry          *sql.Stmt
	appLog                            *logrus.Logger
	reqLog                            *logrus.Logger
	errLog                            *logrus.Logger
	conf                              *config.Config
}

// New returns a new supervisor if none have been initialized yet,
// or the already initialized supervisor if it has.
// It will panic if config has broken cryptographic seeds
func New(c *config.Config) *Supervisor {
	var err error
	if initialized {
		return singleton
	}

	s := &Supervisor{}
	s.conf = c
	s.Input = make(chan msg.Request, s.conf.QueueLen)
	s.Update = make(chan msg.Request, s.conf.QueueLen)
	s.Shutdown = make(chan struct{})
	s.readonly = s.conf.ReadOnly
	if s.seed, err = hex.DecodeString(
		s.conf.Auth.TokenSeed,
	); err != nil {
		panic(err)
	}
	if len(s.seed) == 0 {
		panic(`token.seed has length 0`)
	}
	if s.key, err = hex.DecodeString(
		s.conf.Auth.TokenKey,
	); err != nil {
		panic(err)
	}
	if len(s.key) == 0 {
		panic(`token.key has length 0`)
	}
	s.tokenExpiry = s.conf.Auth.TokenExpirySeconds
	s.kexExpiry = s.conf.Auth.KexExpirySeconds
	s.credExpiry = s.conf.Auth.CredentialExpiryDays
	s.activation = s.conf.Auth.Activation

	// set package variable config for functions
	cfg = c
	// set singleton supervisor instance
	singleton = s
	initialized = true
	return s
}

// Register initializes resources provided by the Soma app
func (s *Supervisor) Register(c *sql.DB, l ...*logrus.Logger) {
	s.conn = c
	s.appLog = l[0]
	s.reqLog = l[1]
	s.errLog = l[2]
}

// Run is the event loop for Supervisor
func (s *Supervisor) Run() {
	var err error

	// set library options
	auth.TokenExpirySeconds = s.tokenExpiry
	auth.KexExpirySeconds = s.kexExpiry

	// initialize maps
	s.tokens = s.newTokenMap()
	s.credentials = s.newCredentialMap()
	s.kex = s.newKexMap()

	// start permission cache
	s.permCache = perm.New()

	// load from database
	s.startupLoad()

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.SelectToken:                   s.stmtTokenSelect,
		stmt.FindUserID:                    s.stmtFindUserID,
		stmt.CategoryList:                  s.stmtCategoryList,
		stmt.CategoryShow:                  s.stmtCategoryShow,
		stmt.PermissionList:                s.stmtPermissionList,
		stmt.PermissionSearchByName:        s.stmtPermissionSearch,
		stmt.SectionList:                   s.stmtSectionList,
		stmt.SectionShow:                   s.stmtSectionShow,
		stmt.SectionSearch:                 s.stmtSectionSearch,
		stmt.ActionList:                    s.stmtActionList,
		stmt.ActionShow:                    s.stmtActionShow,
		stmt.ActionSearch:                  s.stmtActionSearch,
		stmt.SearchGlobalAuthorization:     s.stmtSearchAuthorizationGlobal,
		stmt.SearchRepositoryAuthorization: s.stmtSearchAuthorizationRepository,
		stmt.SearchTeamAuthorization:       s.stmtSearchAuthorizationTeam,
		stmt.SearchMonitoringAuthorization: s.stmtSearchAuthorizationMonitoring,
	} {
		if prepStmt, err = s.conn.Prepare(statement); err != nil {
			s.errLog.Fatal(`supervisor`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

	if !s.readonly {
		for statement, prepStmt := range map[string]*sql.Stmt{
			stmt.CheckUserActive:               s.stmtCheckUserActive,
			stmt.SectionAdd:                    s.stmtSectionAdd,
			stmt.ActionAdd:                     s.stmtActionAdd,
			stmt.RevokeGlobalAuthorization:     s.stmtRevokeAuthorizationGlobal,
			stmt.RevokeRepositoryAuthorization: s.stmtRevokeAuthorizationRepository,
			stmt.RevokeTeamAuthorization:       s.stmtRevokeAuthorizationTeam,
			stmt.RevokeMonitoringAuthorization: s.stmtRevokeAuthorizationMonitoring,
			stmt.GrantGlobalAuthorization:      s.stmtGrantAuthorizationGlobal,
			stmt.GrantRepositoryAuthorization:  s.stmtGrantAuthorizationRepository,
			stmt.GrantTeamAuthorization:        s.stmtGrantAuthorizationTeam,
			stmt.GrantMonitoringAuthorization:  s.stmtGrantAuthorizationMonitoring,
			stmt.PermissionMapEntry:            s.stmtPermissionMapEntry,
			stmt.PermissionUnmapEntry:          s.stmtPermissionUnmapEntry,
		} {
			if prepStmt, err = s.conn.Prepare(statement); err != nil {
				s.errLog.Fatal(`supervisor`, err, stmt.Name(statement))
			}
			defer prepStmt.Close()
		}
	}

runloop:
	for {
		// handle cache updates before handling user requests
		select {
		case req := <-s.Update:
			s.process(&req)
			continue runloop
		default:
			// this empty default case makes this select non-blocking
		}

		// handle whatever request comes in
		select {
		case <-s.Shutdown:
			break runloop
		case req := <-s.Update:
			s.process(&req)
		case req := <-s.Input:
			s.process(&req)
		}
	}
}

func (s *Supervisor) process(q *msg.Request) {
	switch q.Section {
	case msg.SectionSupervisor:
		switch q.Action {
		case msg.ActionKex:
			go func() { s.kexInit(q) }()
		case msg.ActionBootstrap:
			s.bootstrapRoot(q)
		case msg.ActionAuthenticate:
			go func() { s.validateBasicAuth(q) }()
		case msg.ActionToken:
			go func() { s.issueToken(q) }()
		case msg.ActionActivate:
			go func() { s.activateUser(q) }()
		case msg.ActionPassword:
			go func() { s.userPassword(q) }()
		case msg.ActionAuthorize:
			go func() { s.authorize(q) }()
		case msg.ActionCacheUpdate:
			s.cache(q)
		}
	case msg.SectionCategory:
		s.category(q)
	case msg.SectionPermission:
		s.permission(q)
	case msg.SectionRight:
		s.right(q)
	case msg.SectionSection:
		s.section(q)
	case msg.SectionAction:
		s.action(q)
	}
}

// ShutdownNow signals the handler to stop
func (s *Supervisor) ShutdownNow() {
	close(s.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
