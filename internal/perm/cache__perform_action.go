/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package perm

import (
	"github.com/mjolnir42/soma/internal/msg"
)

// These are the per-Action methods used in Cache.Perform

// performActionAdd registers an action
func (c *Cache) performActionAdd(q *msg.Request) {
	c.lock.Lock()
	c.action.add(
		q.ActionObj.SectionID,
		q.ActionObj.SectionName,
		q.ActionObj.ID,
		q.ActionObj.Name,
		q.ActionObj.Category,
	)
	c.lock.Unlock()
}

// performActionRemove removes an action from the cache after
// removing it from all permission maps
func (c *Cache) performActionRemove(q *msg.Request) {
	c.lock.Lock()
	c.performActionRemoveTask(
		q.ActionObj.SectionID,
		q.ActionObj.ID,
	)
	c.lock.Unlock()
}

// performBucketCreate adds a bucket to the object cache
func (c *Cache) performBucketCreate(q *msg.Request) {
	c.lock.Lock()
	c.object.addBucket(
		q.Bucket.RepositoryID,
		q.Bucket.ID,
	)
	c.lock.Unlock()
}

// performBucketDestroy removes a bucket from the object cache
func (c *Cache) performBucketDestroy(q *msg.Request) {
	c.lock.Lock()
	// revoke all grants on the object to be deleted
	grantIDs := c.grantRepository.getObjectGrantID(q.Bucket.ID)
	for _, grantID := range grantIDs {
		c.grantRepository.revoke(grantID)
	}
	// remove object
	c.object.rmBucket(q.Bucket.ID)
	c.lock.Unlock()
}

// performCategoryRemove removes an entire category from the
// cache
func (c *Cache) performCategoryRemove(q *msg.Request) {
	c.lock.Lock()
	c.performCategoryRemoveTask(
		q.Category.Name,
	)
	c.lock.Unlock()
}

// performClusterCreate adds a cluster to the object cache
func (c *Cache) performClusterCreate(q *msg.Request) {
	c.lock.Lock()
	c.object.addCluster(
		q.Cluster.BucketID,
		q.Cluster.ID,
	)
	c.lock.Unlock()
}

// performClusterDestroy removes a cluster from the object cache
func (c *Cache) performClusterDestroy(q *msg.Request) {
	c.lock.Lock()
	// revoke all grants on the object to be deleted
	grantIDs := c.grantRepository.getObjectGrantID(q.Cluster.ID)
	for _, grantID := range grantIDs {
		c.grantRepository.revoke(grantID)
	}
	// remove object
	c.object.rmCluster(q.Cluster.ID)
	c.lock.Unlock()
}

// performGroupCreate adds a group to the object cache
func (c *Cache) performGroupCreate(q *msg.Request) {
	c.lock.Lock()
	c.object.addGroup(
		q.Group.BucketID,
		q.Group.ID,
	)
	c.lock.Unlock()
}

// performGroupDestroy removes a group from the object cache
func (c *Cache) performGroupDestroy(q *msg.Request) {
	c.lock.Lock()
	// revoke all grants on the object to be deleted
	grantIDs := c.grantRepository.getObjectGrantID(q.Group.ID)
	for _, grantID := range grantIDs {
		c.grantRepository.revoke(grantID)
	}
	// remove object
	c.object.rmGroup(q.Group.ID)
	c.lock.Unlock()
}

// performNodeAssign adds a node to the object cache
func (c *Cache) performNodeAssign(q *msg.Request) {
	c.lock.Lock()
	c.object.addNode(
		q.Node.Config.BucketID,
		q.Node.ID,
	)
	c.lock.Unlock()
}

// performNodeUnassign removes a node from the object cache
func (c *Cache) performNodeUnassign(q *msg.Request) {
	c.lock.Lock()
	// revoke all grants on the object to be deleted
	grantIDs := c.grantRepository.getObjectGrantID(q.Node.ID)
	for _, grantID := range grantIDs {
		c.grantRepository.revoke(grantID)
	}
	// remove object
	c.object.rmNode(q.Node.ID)
	c.lock.Unlock()
}

// performPermissionAdd registers a permission
func (c *Cache) performPermissionAdd(q *msg.Request) {
	c.lock.Lock()
	c.pmap.addPermission(
		q.Permission.ID,
		q.Permission.Name,
		q.Permission.Category,
	)
	c.lock.Unlock()
}

// performPermissionMap maps a section or action to a permission
func (c *Cache) performPermissionMap(q *msg.Request) {
	c.lock.Lock()
	// map request can contain either actions or sections, not a mix
	if q.Permission.Actions != nil {
		c.performPermissionMapAction(q)
	}
	if q.Permission.Sections != nil {
		c.performPermissionMapSection(q)
	}
	c.lock.Unlock()
}

// performPermissionRemove removes a permission from the cache
func (c *Cache) performPermissionRemove(q *msg.Request) {
	c.lock.Lock()
	c.performPermissionRemoveTask(q.Permission.ID)
	c.lock.Unlock()
}

// performPermissionUnmap unmaps a section or action from the
// permission
func (c *Cache) performPermissionUnmap(q *msg.Request) {
	c.lock.Lock()
	// unmap request can contain either actions or sections, not a mix
	if q.Permission.Actions != nil {
		c.performPermissionUnmapAction(q)
	}
	if q.Permission.Sections != nil {
		c.performPermissionUnmapSection(q)
	}
	c.lock.Unlock()
}

// performRepositoryCreate adds a new repository to the object cache
func (c *Cache) performRepositoryCreate(q *msg.Request) {
	c.lock.Lock()
	c.object.addRepository(q.Repository.ID)
	c.lock.Unlock()
}

// performRepositoryDestroy removes a repository from the object cache
func (c *Cache) performRepositoryDestroy(q *msg.Request) {
	c.lock.Lock()
	// revoke all grants on the object to be deleted
	grantIDs := c.grantRepository.getObjectGrantID(q.Repository.ID)
	for _, grantID := range grantIDs {
		c.grantRepository.revoke(grantID)
	}
	// remove object
	c.object.rmRepository(q.Repository.ID)
	c.lock.Unlock()
}

// performRightGrant grants a permission
func (c *Cache) performRightGrant(q *msg.Request) {
	c.lock.Lock()
	switch q.Grant.Category {
	case `omnipotence`, `system`, `global`, `permission`, `operations`:
		c.performRightGrantUnscoped(q)
	case `repository`:
		c.performRightGrantScopeRepository(q)
	case `team`:
		c.performRightGrantScopeTeam(q)
	case `monitoring`:
		c.performRightGrantScopeMonitoring(q)
	}
	c.lock.Unlock()
}

// performRightRevoke revokes a permission grant
func (c *Cache) performRightRevoke(q *msg.Request) {
	c.lock.Lock()
	switch q.Grant.Category {
	case `omnipotence`, `system`, `global`, `permission`, `operations`:
		c.performRightRevokeUnscoped(q)
	case `repository`:
		c.performRightRevokeScopeRepository(q)
	case `team`:
		c.performRightRevokeScopeTeam(q)
	case `monitoring`:
		c.performRightRevokeScopeMonitoring(q)
	}
	c.lock.Unlock()
}

// performSectionAdd registers a section
func (c *Cache) performSectionAdd(q *msg.Request) {
	c.lock.Lock()
	c.section.add(
		q.SectionObj.ID,
		q.SectionObj.Name,
		q.SectionObj.Category,
	)
	c.lock.Unlock()
}

// performSectionRemove removes a section after removing all its
// actions and permission mappings
func (c *Cache) performSectionRemove(q *msg.Request) {
	c.lock.Lock()
	c.performSectionRemoveTask(q.SectionObj.ID)
	c.lock.Unlock()
}

// performTeamAdd registers a team
func (c *Cache) performTeamAdd(q *msg.Request) {
	c.lock.Lock()
	c.team.add(
		q.Team.ID,
		q.Team.Name,
	)
	c.lock.Unlock()
}

// performTeamRemove removes a team
func (c *Cache) performTeamRemove(q *msg.Request) {
	c.lock.Lock()
	// revoke all global grants for the team
	grantIDs := c.grantGlobal.getSubjectGrantID(`team`, q.Team.ID)
	for _, grantID := range grantIDs {
		c.grantGlobal.revoke(grantID)
	}
	// revoke all monitoring grants for the team
	grantIDs = c.grantMonitoring.getSubjectGrantID(`team`, q.Team.ID)
	for _, grantID := range grantIDs {
		c.grantMonitoring.revoke(grantID)
	}
	// revoke all repository grants for the team
	grantIDs = c.grantRepository.getSubjectGrantID(`team`, q.Team.ID)
	for _, grantID := range grantIDs {
		c.grantRepository.revoke(grantID)
	}
	// revoke all team grants for the team
	grantIDs = c.grantTeam.getSubjectGrantID(`team`, q.Team.ID)
	for _, grantID := range grantIDs {
		c.grantTeam.revoke(grantID)
	}
	// remove team
	c.team.rmByID(q.Team.ID)
	c.lock.Unlock()
}

// performUserAdd registers a user
func (c *Cache) performUserAdd(q *msg.Request) {
	c.lock.Lock()
	c.user.add(
		q.User.ID,
		q.User.UserName,
		q.User.TeamID,
	)
	c.team.addMember(
		q.User.TeamID,
		q.User.ID,
	)
	c.lock.Unlock()
}

// performUserRemove removes a user
func (c *Cache) performUserRemove(q *msg.Request) {
	c.lock.Lock()
	u := c.user.getByID(q.User.ID)
	if u == nil {
		c.lock.Unlock()
		return
	}
	// revoke all global grants for the user
	grantIDs := c.grantGlobal.getSubjectGrantID(`user`, u.ID)
	for _, grantID := range grantIDs {
		c.grantGlobal.revoke(grantID)
	}
	// revoke all monitoring grants for the user
	grantIDs = c.grantMonitoring.getSubjectGrantID(`user`, u.ID)
	for _, grantID := range grantIDs {
		c.grantMonitoring.revoke(grantID)
	}
	// revoke all repository grants for the user
	grantIDs = c.grantRepository.getSubjectGrantID(`user`, u.ID)
	for _, grantID := range grantIDs {
		c.grantRepository.revoke(grantID)
	}
	// revoke all team grants for the user
	grantIDs = c.grantTeam.getSubjectGrantID(`user`, u.ID)
	for _, grantID := range grantIDs {
		c.grantTeam.revoke(grantID)
	}
	// remove user from team
	c.team.rmMember(u.TeamID, u.ID)
	// remove user
	c.user.rmByID(u.ID)
	c.lock.Unlock()
}

// performUserUpdate updates a user's information
func (c *Cache) performUserUpdate(q *msg.Request) {
	c.lock.Lock()
	old := c.user.getByID(q.User.ID)
	if old == nil {
		panic(`Update on non-existing user -- supervisor corruption`)
	}
	// user switched teams
	if old.TeamID != q.Update.User.TeamID {
		c.team.rmMember(
			old.TeamID,
			q.User.ID,
		)
		c.team.addMember(
			q.Update.User.TeamID,
			q.User.ID,
		)
	}
	// add methods overwrite and can thus be used for updates
	c.user.add(
		q.User.ID,
		q.Update.User.UserName,
		q.Update.User.TeamID,
	)
	c.lock.Unlock()
	// this update may delete the user
	if q.Update.User.IsDeleted {
		c.performUserRemove(q)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
