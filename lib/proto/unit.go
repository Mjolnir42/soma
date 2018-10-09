/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// Unit is a measurement unit that metrics can be collected in
type Unit struct {
	Unit    string       `json:"unit,omitempty"`
	Name    string       `json:"name,omitempty"`
	Details *UnitDetails `json:"details,omitempty"`
}

// Clone returns a copy of u
func (u *Unit) Clone() Unit {
	clone := Unit{
		Unit: u.Unit,
		Name: u.Name,
	}
	if u.Details != nil {
		clone.Details = u.Details.Clone()
	}
	return clone
}

// UnitFilter represents parts of a unit that a unit can be searched by
type UnitFilter struct {
	Unit string `json:"unit,omitempty"`
	Name string `json:"name,omitempty"`
}

// UnitDetails contains metadata about an attribute
type UnitDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
}

// Clone retrurns a copy of e
func (u *UnitDetails) Clone() *UnitDetails {
	clone := &UnitDetails{}
	if u.Creation != nil {
		clone.Creation = u.Creation.Clone()
	}
	return clone
}

// DeepCompare returns true if u and a are equal, excluding details
func (u *Unit) DeepCompare(a *Unit) bool {
	if u.Unit != a.Unit || u.Name != a.Name {
		return false
	}
	return true
}

// NewUnitRequest returns a new Request with fields preallocated
// for filling in Unit data, ensuring no nilptr-deref takes place.
func NewUnitRequest() Request {
	return Request{
		Flags: &Flags{},
		Unit:  &Unit{},
	}
}

// NewUnitFilter returns a new Request with fields preallocated
// for filling in a Unit filter, ensuring no nilptr-deref takes place.
func NewUnitFilter() Request {
	return Request{
		Filter: &Filter{
			Unit: &UnitFilter{},
		},
	}
}

// NewUnitResult returns a new Result with fields preallocated
// for filling in Unit data, ensuring no nilptr-deref takes place.
func NewUnitResult() Result {
	return Result{
		Errors: &[]string{},
		Units:  &[]Unit{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
