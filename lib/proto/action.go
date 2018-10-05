/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// Action represents an action that a section can perform
type Action struct {
	ID          string         `json:"ID,omitempty"`
	Name        string         `json:"name,omitempty"`
	SectionID   string         `json:"sectionID,omitempty"`
	SectionName string         `json:"sectionName,omitempty"`
	Category    string         `json:"category,omitempty"`
	Details     *ActionDetails `json:"details,omitempty"`
}

// Clone returns a copy of a
func (a *Action) Clone() Action {
	clone := Action{
		ID:          a.ID,
		Name:        a.Name,
		SectionID:   a.SectionID,
		SectionName: a.SectionName,
		Category:    a.Category,
	}
	if a.Details != nil {
		clone.Details = a.Details.Clone()
	}
	return clone
}

// ActionDetails contains metadata about an Action
type ActionDetails struct {
	Creation *DetailsCreation `json:"details,omitempty"`
}

// Clone returns a copy of a
func (a *ActionDetails) Clone() *ActionDetails {
	clone := &ActionDetails{}
	if a.Creation != nil {
		clone.Creation = a.Creation.Clone()
	}
	return clone
}

// ActionFilter represents parts of a permission that a permission
// can be searched by
type ActionFilter struct {
	Name      string `json:"name,omitempty"`
	SectionID string `json:"sectionID,omitempty"`
}

// NewActionRequest returns a new Request with fields preallocated
// for filling in Action data, ensuring no nilptr-deref takes place.
func NewActionRequest() Request {
	return Request{
		Flags:  &Flags{},
		Action: &Action{},
	}
}

// NewActionFilter returns a new Request with fields preallocated
// for filling in a Action filter, ensuring no nilptr-deref takes place.
func NewActionFilter() Request {
	return Request{
		Filter: &Filter{
			Action: &ActionFilter{},
		},
	}
}

// NewActionResult returns a new Result with fields preallocated
// for filling in Action data, ensuring no nilptr-deref takes place.
func NewActionResult() Result {
	return Result{
		Errors:  &[]string{},
		Actions: &[]Action{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
