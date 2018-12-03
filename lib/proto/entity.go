/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2017, Jörg Pernfuß
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// Entity is a Type without the golang keyword problem
type Entity struct {
	Name    string         `json:"name,omitempty"`
	Details *EntityDetails `json:"details,omitempty"`
}

// Clone retrurns a copy of e
func (e *Entity) Clone() Entity {
	clone := Entity{
		Name: e.Name,
	}
	if e.Details != nil {
		clone.Details = e.Details.Clone()
	}
	return clone
}

// EntityDetails contains metadata about an Entity
type EntityDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
}

// Clone retrurns a copy of e
func (e *EntityDetails) Clone() *EntityDetails {
	clone := &EntityDetails{}
	if e.Creation != nil {
		clone.Creation = e.Creation.Clone()
	}
	return clone
}

// NewEntityRequest returns a new Request with fields preallocated
// for filling in Entity data, ensuring no nilptr-deref takes place.
func NewEntityRequest() Request {
	return Request{
		Flags:  &Flags{},
		Entity: &Entity{},
	}
}

// NewEntityResult returns a new Result with fields preallocated
// for filling in Entity data, ensuring no nilptr-deref takes place.
func NewEntityResult() Result {
	return Result{
		Errors:   &[]string{},
		Entities: &[]Entity{},
	}
}

// Entity types
const (
	EntityRepository = `repository`
	EntityBucket     = `bucket`
	EntityGroup      = `group`
	EntityCluster    = `cluster`
	EntityNode       = `node`
)

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
