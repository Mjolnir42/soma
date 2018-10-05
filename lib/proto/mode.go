/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// Mode represents a monitoring system mode
type Mode struct {
	Mode    string       `json:"mode,omitempty"`
	Details *ModeDetails `json:"details,omitempty"`
}

// Clone returns a copy of m
func (m *Mode) Clone() Mode {
	clone := Mode{
		Mode: m.Mode,
	}
	if m.Details != nil {
		clone.Details = m.Details.Clone()
	}
	return clone
}

// ModeDetails contains metadata about a Mode
type ModeDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
}

// Clone returns a copy of m
func (m *ModeDetails) Clone() *ModeDetails {
	clone := &ModeDetails{}
	if m.Creation != nil {
		clone.Creation = m.Creation.Clone()
	}
	return clone
}

// NewModeRequest returns a new Request with fields preallocated
// for filling in Mode data, ensuring no nilptr-deref takes place.
func NewModeRequest() Request {
	return Request{
		Flags: &Flags{},
		Mode:  &Mode{},
	}
}

// NewModeResult returns a new Result with fields preallocated
// for filling in Mode data, ensuring no nilptr-deref takes place.
func NewModeResult() Result {
	return Result{
		Errors: &[]string{},
		Modes:  &[]Mode{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
