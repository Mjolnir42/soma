/*-
 * Copyright (c) 2015-2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// State represents the states an object inside a configuration tree can
// be in
type State struct {
	Name    string        `json:"name,omitempty"`
	Details *StateDetails `json:"details,omitempty"`
}

// Clone returns a copy of s
func (s *State) Clone() State {
	clone := State{
		Name: s.Name,
	}
	if s.Details != nil {
		clone.Details = s.Details.Clone()
	}
	return clone
}

// StateDetails contains metadata about a State
type StateDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
}

// Clone returns a copy of s
func (s *StateDetails) Clone() *StateDetails {
	clone := &StateDetails{}
	if s.Creation != nil {
		clone.Creation = s.Creation.Clone()
	}
	return clone
}

// NewStatRequest returns a new Request with fields preallocated
// for filling in State data, ensuring no nilptr-deref takes place.
func NewStateRequest() Request {
	return Request{
		Flags: &Flags{},
		State: &State{},
	}
}

// NewStateResult returns a new Result with fields preallocated
// for filling in State data, ensuring no nilptr-deref takes place.
func NewStateResult() Result {
	return Result{
		Errors: &[]string{},
		States: &[]State{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
