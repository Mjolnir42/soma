/*-
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/mjolnir42/scrypth64"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

func (s *Supervisor) startupLoad() {

	s.startupRoot()

	if !s.readonly {
		s.startupCredentials()
	}

	s.startupTokens()

	s.startupTeam()

	s.startupUser()

	s.startupAdmin()

	s.startupCategory()

	s.startupSection()

	s.startupAction()

	s.startupPermission()

	s.startupPermissionMap()

	s.startupGrantGlobalAuthorization()
	s.startupGrantRepositoryAuthorization()
	s.startupGrantMonitoringAuthorization()
	s.startupGrantTeamAuthorization()
}

func (s *Supervisor) startupRoot() {
	var (
		err                  error
		flag, crypt          string
		mcf                  scrypth64.Mcf
		validFrom, expiresAt time.Time
		state                bool
		rows                 *sql.Rows
	)

	rows, err = s.conn.Query(stmt.LoadRootFlags)
	if err != nil {
		s.errLog.Fatal(`supervisor/load-root-flags,query: `, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&flag,
			&state,
		); err != nil {
			s.errLog.Fatal(`supervisor/load-root-flags,scan: `, err)
		}
		switch flag {
		case `disabled`:
			s.rootDisabled = state
		case `restricted`:
			s.rootRestricted = state
		}
	}
	if err = rows.Err(); err != nil {
		s.errLog.Fatal(`supervisor/load-root-flags,next: `, err)
	}

	// only load root credentials on master instance
	if !s.readonly {
		if err = s.conn.QueryRow(stmt.LoadRootPassword).Scan(
			&crypt,
			&validFrom,
			&expiresAt,
		); err == sql.ErrNoRows {
			// root bootstrap outstanding
			s.appLog.Infoln(`supervisor/startup: root account not bootstrapped yet`)
			return
		} else if err != nil {
			s.errLog.Fatal(`supervisor/load-root-password: `, err)
		}
		if mcf, err = scrypth64.FromString(crypt); err != nil {
			s.errLog.Fatal(`supervisor/string-to-mcf: `, err)
		}
		s.credentials.insert(`root`, uuid.Nil, validFrom.UTC(),
			msg.PosTimeInf.UTC(), mcf)
		s.appLog.Infoln(`supervisor/startup: loaded root account configuration`)
	}
}

func (s *Supervisor) startupCredentials() {
	var (
		err                  error
		rows                 *sql.Rows
		userID, user, crypt  string
		reset                bool
		validFrom, expiresAt time.Time
		id                   uuid.UUID
		mcf                  scrypth64.Mcf
	)

	rows, err = s.conn.Query(stmt.LoadAllUserCredentials)
	if err != nil {
		s.errLog.Fatal(`supervisor/load-credentials,query: `, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&userID,
			&crypt,
			&reset,
			&validFrom,
			&expiresAt,
			&user,
		); err != nil {
			s.errLog.Fatal(`supervisor/load-credentials,scan: `, err)
		}

		if id, err = uuid.FromString(userID); err != nil {
			s.errLog.Fatal(`supervisor/string-to-uuid: `, err)
		}
		if mcf, err = scrypth64.FromString(crypt); err != nil {
			s.errLog.Fatal(`supervisor/string-to-mcf: `, err)
		}

		s.credentials.restore(user, id, validFrom, expiresAt, mcf, reset, true)
	}
	if err = rows.Err(); err != nil {
		s.errLog.Fatal(`supervisor/load-credentials,next: `, err)
	}
}

func (s *Supervisor) startupTokens() {
	var (
		err                         error
		token, salt, valid, expires string
		validFrom, expiresAt        time.Time
		rows                        *sql.Rows
	)

	rows, err = s.conn.Query(stmt.LoadAllTokens)
	if err != nil {
		s.errLog.Fatal(`supervisor/load-tokens,query: `, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&token,
			&salt,
			&validFrom,
			&expiresAt,
		); err != nil {
			s.errLog.Fatal(`supervisor/load-tokens,scan: `, err)
		}
		valid = validFrom.Format(msg.RFC3339Milli)
		expires = expiresAt.Format(msg.RFC3339Milli)

		if err = s.tokens.insert(token, valid, expires, salt); err != nil {
			s.errLog.Fatal(`supervisor/load-tokens,insert: `, err)
		}
	}
	if err = rows.Err(); err != nil {
		s.errLog.Fatal(`supervisor/load-tokens,next: `, err)
	}
}

func (s *Supervisor) startupTeam() {
	var (
		err              error
		teamID, teamName string
		isSystem         bool
		ldapID           int
		rows             *sql.Rows
	)

	rows, err = s.conn.Query(stmt.TeamLoad)
	if err != nil {
		s.errLog.Fatal(`supervisor/load-team,query: `, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&teamID,
			&teamName,
			&ldapID,
			&isSystem,
		); err != nil {
			s.errLog.Fatal(`supervisor/load-team,scan: `, err)
		}
		go func(tID, tName string, lID int, isSys bool) {
			s.Update <- msg.CacheUpdateFromRequest(&msg.Request{
				Section: msg.SectionTeam,
				Action:  msg.ActionAdd,
				Team: proto.Team{
					ID:       tID,
					Name:     tName,
					LdapID:   strconv.Itoa(lID),
					IsSystem: isSys,
				},
			})
		}(teamID, teamName, ldapID, isSystem)
		s.appLog.Infof("supervisor/startup: permCache update - loaded team: %s", teamName)
	}
	if err = rows.Err(); err != nil {
		s.errLog.Fatal(`supervisor/load-team,next: `, err)
	}
}

func (s *Supervisor) startupUser() {
	var (
		err                                  error
		userID, userUID, firstName, lastName string
		mailAddr, teamID                     string
		isActive, isSystem, isDeleted        bool
		employeeNum                          int
		rows                                 *sql.Rows
	)

	rows, err = s.conn.Query(stmt.UserLoad)
	if err != nil {
		s.errLog.Fatal(`supervisor/load-user,query: `, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&userID,
			&userUID,
			&firstName,
			&lastName,
			&employeeNum,
			&mailAddr,
			&isActive,
			&isSystem,
			&isDeleted,
			&teamID,
		); err != nil {
			s.errLog.Fatal(`supervisor/load-user,scan: `, err)
		}
		go func(uID, uUID, fName, lName, mAddr, tID string, eNum int, act, sys, del bool) {
			s.Update <- msg.CacheUpdateFromRequest(&msg.Request{
				Section: msg.SectionUser,
				Action:  msg.ActionAdd,
				User: proto.User{
					ID:             uID,
					UserName:       uUID,
					FirstName:      fName,
					LastName:       lName,
					EmployeeNumber: strconv.Itoa(eNum),
					MailAddress:    mAddr,
					IsActive:       act,
					IsSystem:       sys,
					IsDeleted:      del,
					TeamID:         tID,
				},
			})
		}(userID, userUID, firstName, lastName, mailAddr, teamID, employeeNum, isActive, isSystem, isDeleted)
		s.appLog.Infof("supervisor/startup: permCache update - loaded user: %s", userUID)
	}
	if err = rows.Err(); err != nil {
		s.errLog.Fatal(`supervisor/load-user,next: `, err)
	}
}

func (s *Supervisor) startupAdmin() {
	var (
		err                                  error
		adminID, adminName, userID, userName string
		rows                                 *sql.Rows
	)

	rows, err = s.conn.Query(stmt.AdminLoad)
	if err != nil {
		s.errLog.Fatal(`supervisor/load-admin,query: `, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&adminID,
			&adminName,
			&userID,
			&userName,
		); err != nil {
			s.errLog.Fatal(`supervisor/load-admin,scan: `, err)
		}
		go func(aID, aName, uID, uName string) {
			s.Update <- msg.CacheUpdateFromRequest(&msg.Request{
				Section: msg.SectionAdminMgmt,
				Action:  msg.ActionAdd,
				Admin: proto.Admin{
					ID:       aID,
					Name:     aName,
					UserID:   uID,
					UserName: uName,
				},
			})
		}(adminID, adminName, userID, userName)
		s.appLog.Infof("supervisor/startup: permCache update - loaded admin: %s", adminName)
	}
	if err = rows.Err(); err != nil {
		s.errLog.Fatal(`supervisor/load-admin,next: `, err)
	}
}

func (s *Supervisor) startupCategory() {
	var (
		err      error
		category string
		rows     *sql.Rows
	)

	rows, err = s.conn.Query(stmt.CategoryList)
	if err != nil {
		s.errLog.Fatal(`supervisor/load-category,query: `, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&category,
		); err != nil {
			s.errLog.Fatal(`supervisor/load-category,scan: `, err)
		}
		go func(cat string) {
			s.Update <- msg.CacheUpdateFromRequest(&msg.Request{
				Section: msg.SectionCategory,
				Action:  msg.ActionAdd,
				Category: proto.Category{
					Name: cat,
				},
			})
		}(category)
		s.appLog.Infof("supervisor/startup: permCache update - loaded category: %s", category)
	}
	if err = rows.Err(); err != nil {
		s.errLog.Fatal(`supervisor/load-category,next: `, err)
	}
}

func (s *Supervisor) startupSection() {
	var (
		err                              error
		category, sectionID, sectionName string
		rows                             *sql.Rows
	)

	rows, err = s.conn.Query(stmt.SectionLoad)
	if err != nil {
		s.errLog.Fatal(`supervisor/load-section,query: `, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&sectionID,
			&sectionName,
			&category,
		); err != nil {
			s.errLog.Fatal(`supervisor/load-section,scan: `, err)
		}
		go func(sID, sNam, cat string) {
			s.Update <- msg.CacheUpdateFromRequest(&msg.Request{
				Section: msg.SectionSection,
				Action:  msg.ActionAdd,
				SectionObj: proto.Section{
					ID:       sID,
					Name:     sNam,
					Category: cat,
				},
			})
		}(sectionID, sectionName, category)
		s.appLog.Infof("supervisor/startup: permCache update - loaded section: %s|%s|%s", sectionID, sectionName, category)
	}
	if err = rows.Err(); err != nil {
		s.errLog.Fatal(`supervisor/load-section,next: `, err)
	}
}

func (s *Supervisor) startupAction() {
	var (
		err                    error
		actionID, actionName   string
		sectionID, sectionName string
		category               string
		rows                   *sql.Rows
	)

	rows, err = s.conn.Query(stmt.ActionLoad)
	if err != nil {
		s.errLog.Fatal(`supervisor/load-action,query: `, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&actionID,
			&actionName,
			&sectionID,
			&sectionName,
			&category,
		); err != nil {
			s.errLog.Fatal(`supervisor/load-action,scan: `, err)
		}
		go func(aID, aNam, sID, sNam, cat string) {
			s.Update <- msg.CacheUpdateFromRequest(&msg.Request{
				Section: msg.SectionAction,
				Action:  msg.ActionAdd,
				ActionObj: proto.Action{
					ID:          aID,
					Name:        aNam,
					SectionID:   sID,
					SectionName: sNam,
					Category:    cat,
				},
			})
		}(actionID, actionName, sectionID, sectionName, category)
		s.appLog.Infof("supervisor/startup: permCache update - loaded action: %s|%s|%s|%s|%s", actionID, actionName, sectionID, sectionName, category)
	}
	if err = rows.Err(); err != nil {
		s.errLog.Fatal(`supervisor/load-action,next: `, err)
	}
}

func (s *Supervisor) startupPermission() {
	var (
		err                          error
		permissionID, permissionName string
		category                     string
		rows                         *sql.Rows
	)

	rows, err = s.conn.Query(stmt.PermissionLoad)
	if err != nil {
		s.errLog.Fatal(`supervisor/load-permission,query: `, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&permissionID,
			&permissionName,
			&category,
		); err != nil {
			s.errLog.Fatal(`supervisor/load-permission,scan: `, err)
		}
		go func(pID, pNam, cat string) {
			s.Update <- msg.CacheUpdateFromRequest(&msg.Request{
				Section: msg.SectionPermission,
				Action:  msg.ActionAdd,
				Permission: proto.Permission{
					ID:       pID,
					Name:     pNam,
					Category: cat,
				},
			})
		}(permissionID, permissionName, category)
		s.appLog.Infof("supervisor/startup: permCache update - loaded permission: %s|%s|%s",
			permissionID,
			permissionName,
			category)
	}
	if err = rows.Err(); err != nil {
		s.errLog.Fatal(`supervisor/load-permission,next: `, err)
	}
}

func (s *Supervisor) startupPermissionMap() {
	var (
		err                          error
		permissionID, permissionName string
		mappingID, category          string
		sectionID, sectionName       string
		actionID, actionName         string
		nullActionID, nullActionName sql.NullString
		rows                         *sql.Rows
		sectionMapping               bool
	)

	rows, err = s.conn.Query(stmt.PermissionMapLoad)
	if err != nil {
		s.errLog.Fatal(`supervisor/load-permission-map,query: `, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&mappingID,
			&category,
			&permissionID,
			&permissionName,
			&sectionID,
			&sectionName,
			&nullActionID,
			&nullActionName,
		); err != nil {
			s.errLog.Fatal(`supervisor/load-permission-map,scan: `, err)
		}
		// ID and name must be NULL or NOT NULL at the same time
		if nullActionID.Valid != nullActionName.Valid {
			s.errLog.Fatalf("supervisor/load-permission-map,partial null action %s|%s", nullActionID.String, nullActionName.String)
		}
		switch nullActionID.Valid {
		case true:
			actionID = nullActionID.String
			actionName = nullActionName.String
			sectionMapping = false
		default:
			actionID = ``
			actionName = ``
			sectionMapping = true
		}
		go func(pID, pNam, mID, cat, sID, sNam, aID, aNam string, sm bool) {
			switch sm {
			case true:
				s.Update <- msg.CacheUpdateFromRequest(&msg.Request{
					Section: msg.SectionPermission,
					Action:  msg.ActionMap,
					Permission: proto.Permission{
						ID:       pID,
						Name:     pNam,
						Category: cat,
						Sections: &[]proto.Section{proto.Section{
							ID:       sID,
							Name:     sNam,
							Category: cat,
						}},
					},
				})
			default:
				s.Update <- msg.CacheUpdateFromRequest(&msg.Request{
					Section: msg.SectionPermission,
					Action:  msg.ActionMap,
					Permission: proto.Permission{
						ID:       pID,
						Name:     pNam,
						Category: cat,
						Actions: &[]proto.Action{proto.Action{
							ID:        aID,
							Name:      aNam,
							SectionID: sID,
							Category:  cat,
						}},
					},
				})
			}
		}(permissionID, permissionName, mappingID, category, sectionID, sectionName, actionID, actionName, sectionMapping)

		s.appLog.Infof("supervisor/startup: permCache update - loaded permission map: %s|%s|%s|%s|%s|%s|%s|%s",
			mappingID,
			category,
			permissionID,
			permissionName,
			sectionID,
			sectionName,
			actionID,
			actionName,
		)
	}
	if err = rows.Err(); err != nil {
		s.errLog.Fatal(`supervisor/load-permission-map,next: `, err)
	}
}

func (s *Supervisor) startupGrantGlobalAuthorization() {
	var (
		err                                 error
		grantID, permissionID, category     string
		recipientType, recipientID          string
		nAdminID, nUserID, nToolID, nTeamID sql.NullString
		rows                                *sql.Rows
	)

	rows, err = s.conn.Query(stmt.LoadGlobalAuthorization)
	if err != nil {
		s.errLog.Fatal(`supervisor/load-grant,query: `, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&grantID,
			&nAdminID,
			&nUserID,
			&nToolID,
			&nTeamID,
			&permissionID,
			&category,
		); err != nil {
			s.errLog.Fatal(`supervisor/load-grant,scan: `, err)
		}
		// only one of the can be !NULL, enforced by database check
		// constraint
		switch {
		case nAdminID.Valid:
			recipientType = msg.SubjectAdmin
			recipientID = nAdminID.String
		case nUserID.Valid:
			recipientType = msg.SubjectUser
			recipientID = nUserID.String
		case nToolID.Valid:
			recipientType = msg.SubjectTool
			recipientID = nToolID.String
		case nTeamID.Valid:
			recipientType = msg.SubjectTeam
			recipientID = nTeamID.String
		}
		go func(gID, cat, pID, rTyp, rID string) {
			s.Update <- msg.CacheUpdateFromRequest(&msg.Request{
				Section: msg.SectionRight,
				Action:  msg.ActionGrant,
				Grant: proto.Grant{
					ID:            gID,
					Category:      category,
					PermissionID:  permissionID,
					RecipientType: rTyp,
					RecipientID:   rID,
				},
			})
		}(grantID, category, permissionID, recipientType, recipientID)

		s.appLog.Infof("supervisor/startup: permCache update - loaded right grant: %s|%s|%s|%s|%s",
			grantID,
			category,
			permissionID,
			recipientType,
			recipientID,
		)
	}
	if err = rows.Err(); err != nil {
		s.errLog.Fatal(`supervisor/load-grant,next: `, err)
	}
}

func (s *Supervisor) startupGrantRepositoryAuthorization() {
	var (
		err                                               error
		grantID, permissionID, category                   string
		recipientType, recipientID                        string
		entityType, entityID                              string
		nUserID, nToolID, nTeamID                         sql.NullString
		nRepoID, nBucketID, nGroupID, nClusterID, nNodeID sql.NullString
		rows                                              *sql.Rows
	)

	rows, err = s.conn.Query(stmt.LoadRepositoryAuthorization)
	if err != nil {
		s.errLog.Fatal(`supervisor/load-grant-repository,query: `, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&grantID,
			&nUserID,
			&nToolID,
			&nTeamID,
			&category,
			&permissionID,
			&entityType,
			&nRepoID,
			&nBucketID,
			&nGroupID,
			&nClusterID,
			&nNodeID,
		); err != nil {
			s.errLog.Fatal(`supervisor/load-grant-repository,scan: `, err)
		}
		// only one of the can be !NULL, enforced by database check
		// constraint
		switch {
		case nUserID.Valid:
			recipientType = msg.SubjectUser
			recipientID = nUserID.String
		case nToolID.Valid:
			recipientType = msg.SubjectTool
			recipientID = nToolID.String
		case nTeamID.Valid:
			recipientType = msg.SubjectTeam
			recipientID = nTeamID.String
		}
		switch {
		case nRepoID.Valid:
			if entityType != msg.EntityRepository {
				s.errLog.Fatal(`supervisor/load-grant-repository,validate: `,
					`illegal entity mismatch`)
			}
			entityID = nRepoID.String
		case nBucketID.Valid:
			if entityType != msg.EntityBucket {
				s.errLog.Fatal(`supervisor/load-grant-repository,validate: `,
					`illegal entity mismatch`)
			}
			entityID = nBucketID.String
		case nGroupID.Valid:
			if entityType != msg.EntityGroup {
				s.errLog.Fatal(`supervisor/load-grant-repository,validate: `,
					`illegal entity mismatch`)
			}
			entityID = nGroupID.String
		case nClusterID.Valid:
			if entityType != msg.EntityCluster {
				s.errLog.Fatal(`supervisor/load-grant-repository,validate: `,
					`illegal entity mismatch`)
			}
			entityID = nClusterID.String
		case nNodeID.Valid:
			if entityType != msg.EntityNode {
				s.errLog.Fatal(`supervisor/load-grant-repository,validate: `,
					`illegal entity mismatch`)
			}
			entityID = nNodeID.String
		}
		go func(gID, cat, pID, rTyp, rID, oTyp, oID string) {
			s.Update <- msg.CacheUpdateFromRequest(&msg.Request{
				Section: msg.SectionRight,
				Action:  msg.ActionGrant,
				Grant: proto.Grant{
					ID:            gID,
					Category:      category,
					PermissionID:  permissionID,
					RecipientType: rTyp,
					RecipientID:   rID,
					ObjectType:    oTyp,
					ObjectID:      oID,
				},
			})
		}(grantID, category, permissionID, recipientType, recipientID, entityType, entityID)

		s.appLog.Infof("supervisor/startup: permCache update - loaded repository right grant: %s|%s|%s|%s|%s|%s|%s",
			grantID,
			category,
			permissionID,
			recipientType,
			recipientID,
			entityType,
			entityID,
		)
	}
	if err = rows.Err(); err != nil {
		s.errLog.Fatal(`supervisor/load-grant-repository,next: `, err)
	}
}

func (s *Supervisor) startupGrantMonitoringAuthorization() {
	var (
		err                                           error
		grantID, permissionID, monitoringID, category string
		recipientType, recipientID                    string
		nUserID, nToolID, nTeamID                     sql.NullString
		rows                                          *sql.Rows
	)

	rows, err = s.conn.Query(stmt.LoadMonitoringAuthorization)
	if err != nil {
		s.errLog.Fatal(`supervisor/load-grant-monitoring,query: `, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&grantID,
			&nUserID,
			&nToolID,
			&nTeamID,
			&monitoringID,
			&permissionID,
			&category,
		); err != nil {
			s.errLog.Fatal(`supervisor/load-grant-monitoring,scan: `, err)
		}
		// only one of the can be !NULL, enforced by database check
		// constraint
		switch {
		case nUserID.Valid:
			recipientType = msg.SubjectUser
			recipientID = nUserID.String
		case nToolID.Valid:
			recipientType = msg.SubjectTool
			recipientID = nToolID.String
		case nTeamID.Valid:
			recipientType = msg.SubjectTeam
			recipientID = nTeamID.String
		}
		go func(gID, cat, pID, rTyp, rID, oID string) {
			s.Update <- msg.CacheUpdateFromRequest(&msg.Request{
				Section: msg.SectionRight,
				Action:  msg.ActionGrant,
				Grant: proto.Grant{
					ID:            gID,
					Category:      category,
					PermissionID:  permissionID,
					RecipientType: rTyp,
					RecipientID:   rID,
					ObjectType:    msg.EntityMonitoring,
					ObjectID:      oID,
				},
			})
		}(grantID, category, permissionID, recipientType, recipientID, monitoringID)

		s.appLog.Infof("supervisor/startup: permCache update - loaded monitoring right grant: %s|%s|%s|%s|%s|%s",
			grantID,
			category,
			permissionID,
			recipientType,
			recipientID,
			monitoringID,
		)
	}
	if err = rows.Err(); err != nil {
		s.errLog.Fatal(`supervisor/load-grant-monitoring,next: `, err)
	}
}

func (s *Supervisor) startupGrantTeamAuthorization() {
	var (
		err                                           error
		grantID, permissionID, targetTeamID, category string
		recipientType, recipientID                    string
		nUserID, nToolID, nTeamID                     sql.NullString
		rows                                          *sql.Rows
	)

	rows, err = s.conn.Query(stmt.LoadTeamAuthorization)
	if err != nil {
		s.errLog.Fatal(`supervisor/load-grant-team,query: `, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&grantID,
			&nUserID,
			&nToolID,
			&nTeamID,
			&targetTeamID,
			&permissionID,
			&category,
		); err != nil {
			s.errLog.Fatal(`supervisor/load-grant-team,scan: `, err)
		}
		// only one of the can be !NULL, enforced by database check
		// constraint
		switch {
		case nUserID.Valid:
			recipientType = msg.SubjectUser
			recipientID = nUserID.String
		case nToolID.Valid:
			recipientType = msg.SubjectTool
			recipientID = nToolID.String
		case nTeamID.Valid:
			recipientType = msg.SubjectTeam
			recipientID = nTeamID.String
		}
		go func(gID, cat, pID, rTyp, rID, oID string) {
			s.Update <- msg.CacheUpdateFromRequest(&msg.Request{
				Section: msg.SectionRight,
				Action:  msg.ActionGrant,
				Grant: proto.Grant{
					ID:            gID,
					Category:      category,
					PermissionID:  permissionID,
					RecipientType: rTyp,
					RecipientID:   rID,
					ObjectType:    msg.EntityTeam,
					ObjectID:      oID,
				},
			})
		}(grantID, category, permissionID, recipientType, recipientID, targetTeamID)

		s.appLog.Infof("supervisor/startup: permCache update - loaded team right grant: %s|%s|%s|%s|%s|%s",
			grantID,
			category,
			permissionID,
			recipientType,
			recipientID,
			targetTeamID,
		)
	}
	if err = rows.Err(); err != nil {
		s.errLog.Fatal(`supervisor/load-grant-team,next: `, err)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
