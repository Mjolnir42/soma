/*-
 * Copyright (c) 2015-2016, 1&1 Internet SE
 * Copyright (c) 2018, 1&1 IONOS SE
 * Copyright (c) 2015-2018, Jörg Pernfuß <code.jpe@gmail.com>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// Oncall defines an oncall duty team
type Oncall struct {
	ID      string          `json:"id,omitempty"`
	Name    string          `json:"name,omitempty"`
	Number  string          `json:"number,omitempty"`
	Members *[]OncallMember `json:"members,omitempty"`
	Details *OncallDetails  `json:"details,omitempty"`
}

// Clone returns a copy of o
func (o *Oncall) Clone() Oncall {
	clone := Oncall{
		ID:      o.ID,
		Name:    o.Name,
		Number:  o.Number,
		Members: &[]OncallMember{},
	}
	if o.Details != nil {
		clone.Details = o.Details.Clone()
	}
	if o.Members != nil {
		for i := range *o.Members {
			*clone.Members = append(*clone.Members, (*o.Members)[i].Clone())
		}
	}
	if len(*clone.Members) == 0 {
		clone.Members = nil
	}
	return clone
}

// Sanitize resets some fields in o
func (o *Oncall) Sanitize() {
	o.ID = ``
	o.Members = nil
	o.Details = nil
}

// DeepCompare returns true if o and a are equal, excluding details
func (o *Oncall) DeepCompare(a *Oncall) bool {
	if o.ID != a.ID || o.Name != a.Name || o.Number != a.Number {
		return false
	}

memberloop:
	for _, member := range *o.Members {
		if member.DeepCompareSlice(a.Members) {
			continue memberloop
		}
		return false
	}

revmemberloop:
	for _, member := range *a.Members {
		if member.DeepCompareSlice(o.Members) {
			continue revmemberloop
		}
		return false
	}

	return true
}

// OncallDetails contains metadata about an oncall duty team
type OncallDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
}

// Clone returns a copy of o
func (o *OncallDetails) Clone() *OncallDetails {
	clone := &OncallDetails{}
	if o.Creation != nil {
		clone.Creation = o.Creation.Clone()
	}
	return clone
}

// OncallMember describes a member of an oncall duty team
type OncallMember struct {
	UserID   string `json:"userID,omitempty"`
	UserName string `json:"userName,omitempty"`
}

// Clone returns a copy of o
func (o *OncallMember) Clone() OncallMember {
	return OncallMember{
		UserName: o.UserName,
		UserID:   o.UserID,
	}
}

// DeepCompare returns true if o and a are equal
func (o *OncallMember) DeepCompare(a *OncallMember) bool {
	if o.UserName != a.UserName || o.UserID != a.UserID {
		return false
	}
	return true
}

// DeepCompareSlice returns true if o is equal to an oncall member
// contained in slice a
func (o *OncallMember) DeepCompareSlice(a *[]OncallMember) bool {
	if a == nil || *a == nil {
		return false
	}
	for _, member := range *a {
		if o.DeepCompare(&member) {
			return true
		}
	}
	return false
}

// OncallFilter defines by which attributes an oncall duty team can be
// searched for
type OncallFilter struct {
	Name   string `json:"name,omitempty"`
	Number string `json:"number,omitempty"`
}

// NewOncallRequest returns a new Request with fields preallocated
// for filling in Oncall data, ensuring no nilptr-deref takes place.
func NewOncallRequest() Request {
	return Request{
		Flags:  &Flags{},
		Oncall: &Oncall{},
	}
}

// NewOncallFilter returns a new Request with fields preallocated
// for filling in an Oncall filter, ensuring no nilptr-deref takes place.
func NewOncallFilter() Request {
	return Request{
		Filter: &Filter{
			Oncall: &OncallFilter{},
		},
	}
}

// NewOncallResult returns a new Result with fields preallocated
// for filling in Oncall data, ensuring no nilptr-deref takes place.
func NewOncallResult() Result {
	return Result{
		Errors:  &[]string{},
		Oncalls: &[]Oncall{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
