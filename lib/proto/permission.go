/*-
 * Copyright (c) 2015-2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// Permission represents a set of sections and actions that can be
// granted at once
type Permission struct {
	ID       string             `json:"id,omitempty"`
	Name     string             `json:"name,omitempty"`
	Category string             `json:"category,omitempty"`
	Actions  *[]Action          `json:"actions,omitempty"`
	Sections *[]Section         `json:"sections,omitempty"`
	Details  *PermissionDetails `json:"details,omitempty"`
}

// Clone returns a copy of p
func (p *Permission) Clone() Permission {
	clone := Permission{
		ID:       p.ID,
		Name:     p.Name,
		Category: p.Category,
	}
	if p.Actions != nil && len(*p.Actions) != 0 {
		action := make([]Action, 0, len(*p.Actions))
		for i := range *p.Actions {
			action[i] = (*p.Actions)[i].Clone()
		}
		clone.Actions = &action
	}
	if p.Sections != nil && len(*p.Sections) != 0 {
		section := make([]Section, 0, len(*p.Sections))
		for i := range *p.Sections {
			section[i] = (*p.Sections)[i].Clone()
		}
		clone.Sections = &section
	}
	if p.Details != nil {
		clone.Details = p.Details.Clone()
	}
	return clone
}

// PermissionDetails contains metadata about a Permission
type PermissionDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
}

// Clone retrurns a copy of p
func (p *PermissionDetails) Clone() *PermissionDetails {
	clone := &PermissionDetails{}
	if p.Creation != nil {
		clone.Creation = p.Creation.Clone()
	}
	return clone
}

// PermissionFilter represents parts of a permission that a permission
// can be searched by
type PermissionFilter struct {
	Name     string `json:"name,omitempty"`
	Category string `json:"category,omitempty"`
}

// NewPermissionRequest returns a new Request with fields preallocated
// for filling in Permission data, ensuring no nilptr-deref takes place.
func NewPermissionRequest() Request {
	return Request{
		Flags:      &Flags{},
		Permission: &Permission{},
	}
}

// NewPermissionFilter returns a new Request with fields preallocated
// for filling in a Permission filter, ensuring no nilptr-deref takes place.
func NewPermissionFilter() Request {
	return Request{
		Filter: &Filter{
			Permission: &PermissionFilter{},
		},
	}
}

// NewPermissionResult returns a new Result with fields preallocated
// for filling in Permission data, ensuring no nilptr-deref takes place.
func NewPermissionResult() Result {
	return Result{
		Errors:      &[]string{},
		Permissions: &[]Permission{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
