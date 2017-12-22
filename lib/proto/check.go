/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Check struct {
	CheckID       string `json:"checkId,omitempty"`
	SourceCheckID string `json:"sourceCheckID,omitempty"`
	CheckConfigID string `json:"checkConfigID,omitempty"`
	SourceType    string `json:"sourceType,omitempty"`
	IsInherited   bool   `json:"isInherited,omitempty"`
	InheritedFrom string `json:"inheritedFrom,omitempty"`
	Inheritance   bool   `json:"inheritance,omitempty"`
	ChildrenOnly  bool   `json:"childrenOnly,omitempty"`
	RepositoryID  string `json:"repositoryId,omitempty"`
	BucketID      string `json:"bucketId,omitempty"`
	CapabilityID  string `json:"capabilityId,omitempty"`
}

func (t *Check) DeepCompare(a *Check) bool {
	if t.CheckID != a.CheckID || t.SourceCheckID != a.SourceCheckID ||
		t.CheckConfigID != a.CheckConfigID || t.SourceType != a.SourceType ||
		t.IsInherited != a.IsInherited || t.InheritedFrom != a.InheritedFrom ||
		t.Inheritance != a.Inheritance || t.ChildrenOnly != a.ChildrenOnly ||
		t.RepositoryID != a.RepositoryID || t.BucketID != a.BucketID ||
		t.CapabilityID != a.CapabilityID {
		return false
	}
	return true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
