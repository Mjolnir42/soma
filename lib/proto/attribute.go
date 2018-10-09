/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// Attribute is a key in a service specification
type Attribute struct {
	Name        string            `json:"name,omitempty"`
	Cardinality string            `json:"cardinality,omitempty"`
	Details     *AttributeDetails `json:"details,omitempty"`
}

// Clone returns a copy of e
func (a *Attribute) Clone() Attribute {
	clone := Attribute{
		Name:        a.Name,
		Cardinality: a.Cardinality,
	}
	if a.Details != nil {
		clone.Details = a.Details.Clone()
	}
	return clone
}

// AttributeDetails contains metadata about an attribute
type AttributeDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
}

// Clone retrurns a copy of e
func (a *AttributeDetails) Clone() *AttributeDetails {
	clone := &AttributeDetails{}
	if a.Creation != nil {
		clone.Creation = a.Creation.Clone()
	}
	return clone
}

// NewAttributeRequest returns a new Request with fields preallocated
// for filling in Attribute data, ensuring no nilptr-deref takes place.
func NewAttributeRequest() Request {
	return Request{
		Flags:     &Flags{},
		Attribute: &Attribute{},
	}
}

// NewAttributeResult returns a new Result with fields preallocated
// for filling in Attribute data, ensuring no nilptr-deref takes place.
func NewAttributeResult() Result {
	return Result{
		Errors:     &[]string{},
		Attributes: &[]Attribute{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
