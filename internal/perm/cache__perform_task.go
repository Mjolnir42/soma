/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package perm

import "github.com/mjolnir42/soma/internal/msg"

// These are the methods used when an action can have variations.
// Cache locking is performed by the action methods, tasks do not
// lock the cache!

// performActionRemoveTask implements performActionRemove in a
// reusable way, ie. without locking
func (c *Cache) performActionRemoveTask(sectionID, actionID string) {
	// unmap the action from all permissions
	permIDs := c.pmap.getActionPermissionID(
		sectionID,
		actionID,
	)
	for i := range permIDs {
		c.pmap.unmapAction(
			sectionID,
			actionID,
			permIDs[i],
		)
	}
	// remove the action
	c.action.rmActionByID(
		sectionID,
		actionID,
	)
}

// performCategoryRemoveTask implements performCategoryRemove
// without locking
func (c *Cache) performCategoryRemoveTask(category string) {
	// retrieve permissions in this category
	permIDs := c.pmap.getCategoryPermissionID(category)

	// remove all permissions in this category
	for _, perm := range permIDs {
		c.performPermissionRemoveTask(perm.ID)
	}

	// retrieve all sections and actions in this category
	sectionIDs := c.section.getCategory(category)
	for _, section := range sectionIDs {
		c.performSectionRemoveTask(section.ID)
	}
	// no category to remove as categories are tracked implicitly
}

// performPermissionMapAction maps action to a permission
func (c *Cache) performPermissionMapAction(q *msg.Request) {
	for _, a := range *q.Permission.Actions {
		c.pmap.mapAction(
			a.SectionID,
			a.ID,
			q.Permission.ID,
		)
	}
}

// performPermissionMapSection maps section to a permission
func (c *Cache) performPermissionMapSection(q *msg.Request) {
	for _, s := range *q.Permission.Sections {
		c.pmap.mapSection(
			s.ID,
			q.Permission.ID,
		)
	}
}

// performPermissionRemoveTask implements performPermissionRemove
// without locking
func (c *Cache) performPermissionRemoveTask(permID string) {
	// retrieve category for this permission
	category := c.pmap.getCategory(permID)

	// revoke all permission grants
	switch category {
	case `omnipotence`, `system`, `global`, `permission`, `operations`:
		grantIDs := c.grantGlobal.getPermissionGrantID(permID)
		for _, grantID := range grantIDs {
			c.grantGlobal.revoke(grantID)
		}
	case `repository`:
		grantIDs := c.grantRepository.getPermissionGrantID(permID)
		for _, grantID := range grantIDs {
			c.grantRepository.revoke(grantID)
		}
	case `team`:
		grantIDs := c.grantTeam.getPermissionGrantID(permID)
		for _, grantID := range grantIDs {
			c.grantTeam.revoke(grantID)
		}
	case `monitoring`:
		grantIDs := c.grantMonitoring.getPermissionGrantID(permID)
		for _, grantID := range grantIDs {
			c.grantMonitoring.revoke(grantID)
		}
	}

	// remove permission incl. all mappings
	c.pmap.removePermission(permID)
}

// performPermissionUnmapAction unmaps an action from a permission
func (c *Cache) performPermissionUnmapAction(q *msg.Request) {
	for _, a := range *q.Permission.Actions {
		c.pmap.unmapAction(
			a.SectionID,
			a.ID,
			q.Permission.ID,
		)
	}
}

// performPermissionUnmapSection unmaps a section from a permission
func (c *Cache) performPermissionUnmapSection(q *msg.Request) {
	for _, s := range *q.Permission.Sections {
		c.pmap.unmapSection(
			s.ID,
			q.Permission.ID,
		)
	}
}

// performRightGrantScopeMonitoring grants a monitoring-scoped
// permission
func (c *Cache) performRightGrantScopeMonitoring(q *msg.Request) {
	switch q.Grant.ObjectType {
	case `monitoring`:
		c.grantMonitoring.grant(
			q.Grant.RecipientType,
			q.Grant.RecipientID,
			q.Grant.Category,
			q.Grant.ObjectID,
			q.Grant.PermissionID,
			q.Grant.ID,
		)
	}
}

// performRightGrantScopeRepository grants a repo-scoped permission
func (c *Cache) performRightGrantScopeRepository(q *msg.Request) {
	switch q.Grant.ObjectType {
	case `repository`, `bucket`:
		c.grantRepository.grant(
			q.Grant.RecipientType,
			q.Grant.RecipientID,
			q.Grant.Category,
			q.Grant.ObjectID,
			q.Grant.PermissionID,
			q.Grant.ID,
		)
	}
}

// performRightGrantScopeTeam grants a team-scoped permission
func (c *Cache) performRightGrantScopeTeam(q *msg.Request) {
	switch q.Grant.ObjectType {
	case `team`:
		c.grantTeam.grant(
			q.Grant.RecipientType,
			q.Grant.RecipientID,
			q.Grant.Category,
			q.Grant.ObjectID,
			q.Grant.PermissionID,
			q.Grant.ID,
		)
	}
}

// performRightGrantUnscoped grants a unscoped permission
func (c *Cache) performRightGrantUnscoped(q *msg.Request) {
	c.grantGlobal.grant(
		q.Grant.RecipientType,
		q.Grant.RecipientID,
		q.Grant.Category,
		q.Grant.PermissionID,
		q.Grant.ID,
	)
}

// performRightRevokeScopeMonitoring revokes a monitoring-scoped grant
func (c *Cache) performRightRevokeScopeMonitoring(q *msg.Request) {
	c.grantMonitoring.revoke(q.Grant.ID)
}

// performRightRevokeScopeRepository revokes a repo-scoped grant
func (c *Cache) performRightRevokeScopeRepository(q *msg.Request) {
	c.grantRepository.revoke(q.Grant.ID)
}

// performRightRevokeScopeTeam revokes a team-scoped grant
func (c *Cache) performRightRevokeScopeTeam(q *msg.Request) {
	c.grantTeam.revoke(q.Grant.ID)
}

// performRightRevokeUnscoped revokes a global grant
func (c *Cache) performRightRevokeUnscoped(q *msg.Request) {
	c.grantGlobal.revoke(q.Grant.ID)
}

// performSectionRemoveTask implements performSectionRemove without
// locking
func (c *Cache) performSectionRemoveTask(sectionID string) {
	// delete all member actions
	actionIDs := c.action.getActionsBySectionID(sectionID)
	for _, actionID := range actionIDs {
		c.performActionRemoveTask(sectionID, actionID)
	}
	// delete section from action lookup
	c.action.rmSectionByID(sectionID)

	// unmap the section from all permissions
	permIDs := c.pmap.getSectionPermissionID(
		sectionID,
	)
	for i := range permIDs {
		c.pmap.unmapSection(
			sectionID,
			permIDs[i],
		)
	}
	// remove the section
	c.section.rmByID(sectionID)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
