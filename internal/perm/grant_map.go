/*-
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package perm

import "fmt"

// unscopedGrantMap is the cache data structure for global permission
// grants. It covers the categories 'omnipotence', 'system', 'global',
// 'permission' and 'operations'.
type unscopedGrantMap struct {
	// subjectId -> category -> permissionID -> grantID
	// The subjectID is built as follows:
	//	user:${userUUID}
	//	admin:${adminUUID}
	//	tool:${toolUUID}
	//	team:${teamUUID}
	grants map[string]map[string]map[string]string
}

// newUnscopedGrantMap returns an initialized unscopedGrantMap
func newUnscopedGrantMap() *unscopedGrantMap {
	u := unscopedGrantMap{}
	u.grants = map[string]map[string]map[string]string{}
	return &u
}

// grant records a grant of a permission to a subject into the cache
func (m *unscopedGrantMap) grant(subjType, subjID, category,
	permissionID, grantID string) {
	// only accept these four types
	switch subjType {
	case `user`, `admin`, `tool`, `team`:
	default:
		return
	}
	subject := fmt.Sprintf("%s:%s", subjType, subjID)

	//
	if _, ok := m.grants[subject]; !ok {
		m.grants[subject] = map[string]map[string]string{}
	}
	//
	if _, ok := m.grants[subject][category]; !ok {
		m.grants[subject][category] = map[string]string{}
	}
	m.grants[subject][category][permissionID] = grantID
}

// scopedGrantMap is the cache data structure for permission grants
// on an object.
type scopedGrantMap struct {
	// subjectID -> category -> permissionID -> objectID -> grantID
	// The subjectID is built as follows:
	//	user:${userUUID}
	//	admin:${adminUUID}
	//	tool:${toolUUID}
	//	team:${teamUUID}
	grants map[string]map[string]map[string]map[string]string
}

// newScopedGrantMap return ans initialized scopedGrantMap
func newScopedGrantMap() *scopedGrantMap {
	s := scopedGrantMap{}
	s.grants = map[string]map[string]map[string]map[string]string{}
	return &s
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix