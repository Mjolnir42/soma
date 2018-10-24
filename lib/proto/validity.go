/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// Validity defines in which ways it is valid for an entity to
// receive a specific system property
type Validity struct {
	SystemProperty string           `json:"systemProperty,omitempty"`
	Entity         string           `json:"entity,omitempty"`
	Direct         bool             `json:"direct,string,omitempty"`
	Inherited      bool             `json:"inherited,string,omitempty"`
	Details        *ValidityDetails `json:"details,omitempty"`
}

// Clone returns a copy of v
func (v *Validity) Clone() Validity {
	clone := Validity{
		SystemProperty: v.SystemProperty,
		Entity:         v.Entity,
		Direct:         v.Direct,
		Inherited:      v.Inherited,
	}
	if v.Details != nil {
		clone.Details = v.Details.Clone()
	}
	return clone
}

// ValidityDetails contains metadata about a Validity definition
type ValidityDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
	DetailsCreation
}

// Clone returns a copy of v
func (v *ValidityDetails) Clone() *ValidityDetails {
	clone := &ValidityDetails{}
	if v.Creation != nil {
		clone.Creation = v.Creation.Clone()
	}
	return clone
}

// NewValidityRequest returns a new Request with fields preallocated
// for filling in Validity data, ensuring no nilptr-deref takes place.
func NewValidityRequest() Request {
	return Request{
		Flags:    &Flags{},
		Validity: &Validity{},
	}
}

// NewValidityResult returns a new Result with fields preallocated
// for filling in Validity data, ensuring no nilptr-deref takes place.
func NewValidityResult() Result {
	return Result{
		Errors:     &[]string{},
		Validities: &[]Validity{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
