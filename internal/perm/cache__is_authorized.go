/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package perm // import "github.com/mjolnir42/soma/internal/perm"

import (
	"strings"

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/proto"
)

// isAuthorized implements Cache.IsAuthorized and checks if the
// request is authorized
func (c *Cache) isAuthorized(q *msg.Request) msg.Result {
	result := msg.FromRequest(q)
	// default action is to deny
	result.Super.Verdict = 403
	result.Super.Audit = result.Super.Audit.
		WithField(`permCache::status`, `received`).
		WithField(`permCache::result`, `forbidden`).
		WithField(`permCache::error`, `none`)

	var user *proto.User
	var subjType, category, actionID, sectionID string
	var sectionPermIDs, actionPermIDs, mergedPermIDs []string
	var any bool

	// determine type of the request subject
	switch {
	case strings.HasPrefix(q.Super.Authorize.AuthUser, `admin_`):
		subjType = `admin`
	case strings.HasPrefix(q.Super.Authorize.AuthUser, `tool_`):
		subjType = `tool`
	default:
		subjType = `user`
	}
	result.Super.Audit = result.Super.Audit.WithField(`permCache::subjectType`, subjType)

	// set readlock on the cache
	c.lock.RLock()
	defer c.lock.RUnlock()

	// look up the user, also handles admin and tool accounts
	result.Super.Audit = result.Super.Audit.WithField(`permCache::subject`, q.Super.Authorize.AuthUser)
	if user = c.user.getByName(q.Super.Authorize.AuthUser); user == nil {
		result.Super.Audit = result.Super.Audit.WithField(`permCache::status`, `processing`).
			WithField(`permCache::result`, `InternalError`).
			WithField(`permCache::error`, `UserNotFound`)
		goto dispatch
	}

	// check if the subject has omnipotence
	if c.checkOmnipotence(subjType, user.ID, &result) {
		result.Super.Audit = result.Super.Audit.
			WithField(`permCache::status`, `evaluated`).
			WithField(`permCache::result`, `omnipotent`)
		result.Super.Verdict = 200
		goto dispatch
	}
	result.Super.Audit = result.Super.Audit.WithField(`permCache::omnipotence`, `false`)

	// root can not receive additional permissions
	switch q.Super.Authorize.AuthUser {
	case `root`:
		result.Super.Audit = result.Super.Audit.
			WithField(`permCache::status`, `exhausted`)
		goto dispatch
	}

	// extract category
	// XXX BUG nilptr crash if section is not found
	category = c.section.getByName(q.Super.Authorize.Section).Category

	// lookup sectionID and actionID of the Request, abort for
	// unknown actions
	if action := c.action.getByName(
		q.Super.Authorize.Section,
		q.Super.Authorize.Action,
	); action == nil {
		result.Super.Audit = result.Super.Audit.
			WithField(`permCache::status`, `processing`).
			WithField(`permCache::result`, `InternalError`).
			WithField(`permCache::error`, `ActionNotFound`)
		goto dispatch
	} else {
		sectionID = action.SectionID
		actionID = action.ID
	}

	// check if the user has the correct system permission
	if ok, invalid := c.checkSystem(category, subjType,
		user.ID, &result); invalid {
		result.Super.Audit = result.Super.Audit.
			WithField(`permCache::status`, `processing`).
			WithField(`permCache::result`, `InternalError`).
			WithField(`permCache::error`, `InvalidMissingSystemPermission`)
		goto dispatch
	} else if ok {
		result.Super.Audit = result.Super.Audit.
			WithField(`permCache::status`, `evaluated`).
			WithField(`permCache::result`, `systempermission`)
		result.Super.Verdict = 200
		goto dispatch
	}
	result.Super.Audit = result.Super.Audit.WithField(`permCache::system`, `false`)

	// admin accounts can not receive regular permissions
	switch subjType {
	case `admin`:
		result.Super.Audit = result.Super.Audit.
			WithField(`permCache::status`, `exhausted`)
		goto dispatch
	}

	// TODO if (       q.Super.Authorize.Section == msg.SectionRight )
	// TODO    and (   q.Super.Authorize.Action  == msg.ActionGrant
	// TODO         or q.Super.Authorize.Action  == msg.ActionRevoke )
	// TODO then the user could also be authorized by a permission from
	// TODO the :grant category

	// lookup all permissionIDs that map either section or action
	sectionPermIDs = c.pmap.getSectionPermissionID(sectionID)
	actionPermIDs = c.pmap.getActionPermissionID(sectionID, actionID)
	mergedPermIDs = append(sectionPermIDs, actionPermIDs...)

	// check if we care about the specific object
	switch q.Super.Authorize.Action {
	case `list`, `search`:
		any = true
	}

	// check if the user has one the permissions that map the
	// requested action
	if c.checkPermission(mergedPermIDs, any, q.Super.Authorize, subjType, user.ID,
		category, &result) {
		result.Super.Audit = result.Super.Audit.
			WithField(`permCache::status`, `evaluated`).
			WithField(`permCache::result`, `userpermission`)
		result.Super.Verdict = 200
		goto dispatch
	}
	result.Super.Audit = result.Super.Audit.WithField(`permCache::direct-user`, `false`)

	// admin and tool accounts do not inherit team rights,
	// authorization check ends here
	switch subjType {
	case `admin`, `tool`:
		result.Super.Audit = result.Super.Audit.
			WithField(`permCache::status`, `exhausted`)
		goto dispatch
	}

	// check if the user's team has a specific grant for the action
	if c.checkPermission(mergedPermIDs, any, q.Super.Authorize, `team`, user.TeamID,
		category, &result) {
		result.Super.Audit = result.Super.Audit.
			WithField(`permCache::status`, `evaluated`).
			WithField(`permCache::result`, `teampermission`)
		result.Super.Verdict = 200
		goto dispatch
	}
	result.Super.Audit = result.Super.Audit.
		WithField(`permCache::inherited-team`, `false`).
		WithField(`permCache::status`, `exhausted`)

dispatch:
	return result
}

// checkOmnipotence returns true if the subject is omnipotent
func (c *Cache) checkOmnipotence(subjectType, subjectID string, result *msg.Result) bool {
	return c.grantGlobal.assess(
		subjectType,
		subjectID,
		`omnipotence`,
		`00000000-0000-0000-0000-000000000000`,
		result,
	)
}

// checkSystem returns true,false if the subject has the system
// permission for the category. If no system permission exists it
// returns false,true
func (c *Cache) checkSystem(category, subjectType,
	subjectID string, result *msg.Result) (bool, bool) {
	permID := c.pmap.getIDByName(`system`, category)
	if permID == `` {
		// there must be a system permission for every category,
		// refuse authorization since the permission cache is broken
		return false, true
	}
	return c.grantGlobal.assess(
		subjectType,
		subjectID,
		`system`,
		permID,
		result,
	), false
}

// checkPermission returns true if the subject has a grant for the
// requested action
func (c *Cache) checkPermission(permIDs []string, any bool,
	q *msg.Request, subjectType, subjectID, category string, result *msg.Result) bool {
	var objID string

permloop:
	for _, permID := range permIDs {
		// determine objID
		if any {
			// invalid uuid
			objID = msg.InvalidObjectID
		} else {
			switch q.Section {
			// per-monitoring scope
			case msg.SectionMonitoring, msg.SectionCapability, msg.SectionDeployment:
				objID = q.Monitoring.ID
			// per-team scope
			case msg.SectionPropertyService:
				objID = q.Property.Service.TeamID
			case msg.SectionNode:
				objID = q.Node.TeamID
			case msg.SectionRepository:
				objID = q.Repository.TeamID
			// per-repository scope
			case msg.SectionInstance, msg.SectionNodeConfig, msg.SectionPropertyCustom,
				msg.SectionRepositoryConfig:
				objID = q.Repository.ID
			case msg.SectionBucket, msg.SectionCluster, msg.SectionCheckConfig,
				msg.SectionGroup:
				objID = q.Bucket.ID
			// global scope
			default:
				// invalid uuid
				objID = msg.InvalidObjectID
			}
		}

		// check authorization
		switch q.Section {
		case msg.SectionMonitoring, msg.SectionCapability, msg.SectionDeployment:
			// per-monitoring sections
			if c.grantMonitoring.assess(subjectType, subjectID,
				category, objID, permID, any, result) {
				return true
			}
		case msg.SectionBucket, msg.SectionCheckConfig, msg.SectionCluster,
			msg.SectionGroup, msg.SectionInstance, msg.SectionNodeConfig,
			msg.SectionPropertyCustom, msg.SectionRepositoryConfig:
			// per-repository sections
			if c.grantRepository.assess(subjectType, subjectID,
				category, objID, permID, any, result) {
				return true
			}
			switch q.Section {
			case msg.SectionBucket, msg.SectionCluster, msg.SectionCheckConfig,
				msg.SectionGroup:
				// permission could be on the repository
				objID = c.object.repoForBucket(q.Bucket.ID)
				if objID == `` {
					continue permloop
				}
				if c.grantRepository.assess(subjectType, subjectID,
					category, objID, permID, any, result) {
					return true
				}
			}
		case msg.SectionNode, msg.SectionPropertyService, msg.SectionRepository:
			// per-team sections
			if c.grantTeam.assess(subjectType, subjectID,
				category, objID, permID, any, result) {
				return true
			}
		default:
			// global sections
			if c.grantGlobal.assess(subjectType, subjectID, category,
				permID, result) {
				return true
			}
		}
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
