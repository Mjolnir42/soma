/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

// Package perm implements the permission cache module for the
// SOMA supervisor. It tracks which actions are mapped to permissions
// and which permissions have been granted.
//
// It can be queried whether a given user is authorized to perform
// an action.
package perm // import "github.com/mjolnir42/soma/internal/perm"

import (
	"fmt"
	"sync"

	"github.com/mjolnir42/soma/internal/msg"
)

// Cache is a permission cache for the SOMA supervisor
type Cache struct {
	// the entire cache has one global mutex, since many actions
	// requires updates to multiple data structures. Locking each
	// of them individually could therefor lead to deadlocks.
	// A global lock is more robust than a lock order scheme, which
	// could still be adopted later as a performance improvement.
	lock sync.RWMutex

	// general ID<>name lookup maps
	section *sectionLookup
	action  *actionLookup
	user    *userLookup
	team    *teamLookup

	// semi-flat repository object lookup map
	object *objectLookup

	// keeps track which actions are mapped to which permissions
	pmap *permissionMapping

	// keeps track of permission grants
	grantGlobal     *unscopedGrantMap
	grantRepository *scopedGrantMap
	grantTeam       *scopedGrantMap
	grantMonitoring *scopedGrantMap
}

// New returns a new permission cache
func New() *Cache {
	c := Cache{}
	c.lock = sync.RWMutex{}
	c.section = newSectionLookup()
	c.action = newActionLookup()
	c.user = newUserLookup()
	c.team = newTeamLookup()
	c.object = newObjectLookup()
	c.pmap = newPermissionMapping()
	c.grantGlobal = newUnscopedGrantMap()
	c.grantRepository = newScopedGrantMap(`repository`)
	c.grantTeam = newScopedGrantMap(`team`)
	c.grantMonitoring = newScopedGrantMap(`monitoring`)
	return &c
}

// Perform executes the request on the cache
func (c *Cache) Perform(q *msg.Request) {
	// delegate the request to per-section methods
	switch q.Cache.Section {
	case msg.SectionAction:
		c.performAction(q.Cache)
	case msg.SectionAdminMgmt:
		c.performAdmin(q.Cache)
	case msg.SectionBucket:
		c.performBucket(q.Cache)
	case msg.SectionCategory:
		c.performCategory(q.Cache)
	case msg.SectionCluster:
		c.performCluster(q.Cache)
	case msg.SectionGroup:
		c.performGroup(q.Cache)
	case msg.SectionNodeConfig:
		c.performNode(q.Cache)
	case msg.SectionPermission:
		c.performPermission(q.Cache)
	case msg.SectionRepository, msg.SectionRepositoryMgmt:
		c.performRepository(q.Cache)
	case msg.SectionRight:
		c.performRight(q.Cache)
	case msg.SectionSection:
		c.performSection(q.Cache)
	case msg.SectionTeam, msg.SectionTeamMgmt:
		c.performTeam(q.Cache)
	case msg.SectionUser, msg.SectionUserMgmt:
		c.performUser(q.Cache)
	default:
		panic(fmt.Sprintf(
			"Unhandled permission cache update in section: %s",
			q.Cache.Section,
		))
	}
}

// Compact frees up memory used by arrays that is no longer
// reference by the slice built on top of them
func (c *Cache) Compact() {
	c.lock.Lock()
	c.pmap.compact()
	c.section.compact()
	c.team.compact()
	c.lock.Unlock()
}

// IsAuthorized checks if q describes an authorized request
func (c *Cache) IsAuthorized(q *msg.Request) msg.Result {
	return c.isAuthorized(q)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
