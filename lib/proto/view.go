/*-
 * Copyright (c) 2015-2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// View represents a data dimension within a configuration tree
type View struct {
	Name    string       `json:"name,omitempty"`
	Details *ViewDetails `json:"details,omitempty"`
}

// Clone returns a copy of v
func (v *View) Clone() View {
	clone := View{
		Name: v.Name,
	}
	if v.Details != nil {
		clone.Details = v.Details.Clone()
	}
	return clone
}

// ViewDetails contains metadata about a View
type ViewDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
}

// Clone returns a copy of v
func (v *ViewDetails) Clone() *ViewDetails {
	clone := &ViewDetails{}
	if v.Creation != nil {
		clone.Creation = v.Creation.Clone()
	}
	return clone
}

// NewViewRequest returns a new Request with fields preallocated
// for filling in View data, ensuring no nilptr-deref takes place.
func NewViewRequest() Request {
	return Request{
		Flags: &Flags{},
		View:  &View{},
	}
}

// NewViewResult returns a new Result with fields preallocated
// for filling in View data, ensuring no nilptr-deref takes place.
func NewViewResult() Result {
	return Result{
		Errors: &[]string{},
		Views:  &[]View{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
