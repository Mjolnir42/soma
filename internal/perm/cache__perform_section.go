/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package perm

import "github.com/mjolnir42/soma/internal/msg"

// These are the per-Section methods used in Cache.Perform

func (c *Cache) performAction(q *msg.Request) {
	switch q.Action {
	case msg.ActionAdd:
		c.performActionAdd(q)
	case msg.ActionRemove:
		c.performActionRemove(q)
	}
}

func (c *Cache) performBucket(q *msg.Request) {
	switch q.Action {
	case msg.ActionCreate:
		c.performBucketCreate(q)
	case msg.ActionDestroy:
		c.performBucketDestroy(q)
	}
}

func (c *Cache) performCategory(q *msg.Request) {
	switch q.Action {
	case msg.ActionRemove:
		c.performCategoryRemove(q)
	}
}

func (c *Cache) performCluster(q *msg.Request) {
	switch q.Action {
	case msg.ActionCreate:
		c.performClusterCreate(q)
	case msg.ActionDestroy:
		c.performClusterDestroy(q)
	}
}

func (c *Cache) performGroup(q *msg.Request) {
	switch q.Action {
	case msg.ActionCreate:
		c.performGroupCreate(q)
	case msg.ActionDestroy:
		c.performGroupDestroy(q)
	}
}

func (c *Cache) performNode(q *msg.Request) {
	switch q.Action {
	case msg.ActionAssign:
		c.performNodeAssign(q)
	case msg.ActionUnassign:
		c.performNodeUnassign(q)
	}
}

func (c *Cache) performPermission(q *msg.Request) {
	switch q.Action {
	case msg.ActionAdd:
		c.performPermissionAdd(q)
	case msg.ActionRemove:
		c.performPermissionRemove(q)
	case msg.ActionMap:
		c.performPermissionMap(q)
	case msg.ActionUnmap:
		c.performPermissionUnmap(q)
	}
}

func (c *Cache) performRepository(q *msg.Request) {
	switch q.Action {
	case msg.ActionCreate:
		c.performRepositoryCreate(q)
	case msg.ActionDestroy:
		c.performRepositoryDestroy(q)
	}
}

func (c *Cache) performRight(q *msg.Request) {
	switch q.Action {
	case msg.ActionGrant:
		c.performRightGrant(q)
	case msg.ActionRevoke:
		c.performRightRevoke(q)
	}
}

func (c *Cache) performSection(q *msg.Request) {
	switch q.Action {
	case msg.ActionAdd:
		c.performSectionAdd(q)
	case msg.ActionRemove:
		c.performSectionRemove(q)
	}
}

func (c *Cache) performTeam(q *msg.Request) {
	switch q.Action {
	case msg.ActionAdd:
		c.performTeamAdd(q)
	case msg.ActionRemove:
		c.performTeamRemove(q)
	case msg.ActionUpdate:
		// XXX TODO
	}
}

func (c *Cache) performUser(q *msg.Request) {
	switch q.Action {
	case msg.ActionAdd:
		c.performUserAdd(q)
	case msg.ActionRemove, msg.ActionPurge:
		c.performUserRemove(q)
	case msg.ActionUpdate:
		// XXX TODO
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
