/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2018, Jörg Pernfuß
 * Copyright (c) 2018, 1&1 IONOS SE
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// Capability defines a metric that can be monitored by a specific
// monitoring system
type Capability struct {
	ID           string                   `json:"ID,omitempty"`
	Name         string                   `json:"name,omitempty"`
	MonitoringID string                   `json:"monitoringID,omitempty"`
	Metric       string                   `json:"metric,omitempty"`
	View         string                   `json:"view,omitempty"`
	Thresholds   uint64                   `json:"thresholds,omitempty"`
	Demux        *[]string               `json:"demux,omitempty"`
	Constraints  *[]CapabilityConstraint `json:"constraints,omitempty"`
	Details      *CapabilityDetails      `json:"details,omitempty"`
}

// Clone returns a copy of c
func (c *Capability) Clone() Capability {
	clone := Capability{
		ID:           c.ID,
		Name:         c.Name,
		MonitoringID: c.MonitoringID,
		Metric:       c.Metric,
		View:         c.View,
		Thresholds:   c.Thresholds,
	}
	if c.Details != nil {
		clone.Details = c.Details.Clone()
	}
	// XXX Demux
	// XXX Constraints
	return clone
}

type CapabilityConstraint struct {
	Type  string
	Name  string
	Value string
}

// DeepCompare returns true if c and a are equal, excluding details
func (c *Capability) DeepCompare(a *Capability) bool {
	if c.ID != a.ID {
		return false
	}
	if c.Name != a.Name {
		return false
	}
	if c.MonitoringID != a.MonitoringID {
		return false
	}
	if c.Metric != a.Metric {
		return false
	}
	if c.View != a.View {
		return false
	}
	if c.Thresholds != a.Thresholds {
		return false
	}
	if c.Demux != nil {
	demuxloop:
		for i := range *c.Demux {
			if (*c.Demux)[i].DeepCompareSlice(a.Demux) {
				continue demuxloop
			}
			return false
		}
	}
	if a.Demux != nil {
	revdemuxloop:
		for i := range *a.Demux {
			if (*a.Demux)[i].DeepCompareSlice(c.Demux) {
				continue revdemuxloop
			}
			return false
		}
	}
	if c.Constraints != nil {
	constraintloop:
		for i := range *c.Constraints {
			if (*c.Constraints)[i].DeepCompareSlice(*a.Constraints) {
				continue constraintloop
			}
			return false
		}
	}
	if a.Constraints != nil {
	revconstraintloop:
		for i := range *a.Constraints {
			if (*a.Constraints)[i].DeepCompareSlice(*c.Constraints) {
				continue revconstraintloop
			}
			return false
		}
	}
	return true
}

// CapabilityDetails contains metadata about a capability
type CapabilityDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
}

// Clone returns a copy of c
func (c *CapabilityDetails) Clone() *CapabilityDetails {
	clone := &CapabilityDetails{}
	if c.Creation != nil {
		clone.Creation = c.Creation.Clone()
	}
	return clone
}

// CapabilityFilter defines by what a capability can be searched by
type CapabilityFilter struct {
	MonitoringID   string `json:"monitoringID,omitempty"`
	MonitoringName string `json:"monitoringName,omitempty"`
	Metric         string `json:"metric,omitempty"`
	View           string `json:"view,omitempty"`
}

// NewCapabilityRequest returns a new Request with fields preallocated
// for filling in Capability data, ensuring no nilptr-deref takes place.
func NewCapabilityRequest() Request {
	return Request{
		Flags: &Flags{},
		Capability: &Capability{
			Demux:       &[]Attribute{},
			Constraints: &[]CheckConfigConstraint{},
		},
	}
}

// NewCapabilityFilter returns a new Request with fields preallocated
// for filling in an Capability filter, ensuring no nilptr-deref takes place.
func NewCapabilityFilter() Request {
	return Request{
		Filter: &Filter{
			Capability: &CapabilityFilter{},
		},
	}
}

// NewCapabilityResult returns a new Result with fields preallocated
// for filling in Capability data, ensuring no nilptr-deref takes place.
func NewCapabilityResult() Result {
	return Result{
		Errors:       &[]string{},
		Capabilities: &[]Capability{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
