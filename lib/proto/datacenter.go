/*-
 * Copyright (c) 2015-2016, 1&1 Internet SE
 * Copyright (c) 2015-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// Datacenter is the definition of a datacenter
type Datacenter struct {
	LoCode  string             `json:"loCode,omitempty"`
	Details *DatacenterDetails `json:"details,omitempty"`
}

// Clone returns a copy of d
func (d *Datacenter) Clone() Datacenter {
	clone := Datacenter{
		LoCode: d.LoCode,
	}
	if d.Details != nil {
		clone.Details = d.Details.Clone()
	}
	return clone
}

// DatacenterDetails contains metadata about a datacenter
type DatacenterDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
}

// Clone returns a copy of d
func (d *DatacenterDetails) Clone() *DatacenterDetails {
	clone := &DatacenterDetails{}
	if d.Creation != nil {
		clone.Creation = d.Creation.Clone()
	}
	return clone
}

// NewDatacenterRequest returns a new Request with fields preallocated
// for filling in Datacenter data, ensuring no nilptr-deref takes place.
func NewDatacenterRequest() Request {
	return Request{
		Flags:      &Flags{},
		Datacenter: &Datacenter{},
	}
}

// NewDatacenterResult returns a new Result with fields preallocated
// for filling in Datacenter data, ensuring no nilptr-deref takes place.
func NewDatacenterResult() Result {
	return Result{
		Errors:      &[]string{},
		Datacenters: &[]Datacenter{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
