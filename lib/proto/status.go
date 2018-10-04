/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

// Status describes the possible status descriptions a check instance
// can be in
type Status struct {
	Name    string         `json:"name,omitempty"`
	Details *StatusDetails `json:"details,omitempty"`
}

// Clone returns a copy of s
func (s *Status) Clone() Status {
	clone := Status{
		Name: s.Name,
	}
	if s.Details != nil {
		clone.Details = s.Details.Clone()
	}
	return clone
}

// StatusDetails contains metadata about a Status
type StatusDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
}

// Clone returns a copy of s
func (s *StatusDetails) Clone() *StatusDetails {
	clone := &StatusDetails{}
	if s.Creation != nil {
		clone.Creation = s.Creation.Clone()
	}
	return clone
}

func NewStatusRequest() Request {
	return Request{
		Flags:  &Flags{},
		Status: &Status{},
	}
}

func NewStatusResult() Result {
	return Result{
		Errors: &[]string{},
		Status: &[]Status{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
