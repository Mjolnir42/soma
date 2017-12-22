/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Capability struct {
	ID           string                  `json:"ID,omitempty"`
	Name         string                  `json:"name,omitempty"`
	MonitoringID string                  `json:"monitoringID,omitempty"`
	Metric       string                  `json:"metric,omitempty"`
	View         string                  `json:"view,omitempty"`
	Thresholds   uint64                  `json:"thresholds,omitempty"`
	Demux        *[]string               `json:"demux,omitempty"`
	Constraints  *[]CapabilityConstraint `json:"constraints,omitempty"`
	Details      *CapabilityDetails      `json:"details,omitempty"`
}

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

type CapabilityFilter struct {
	MonitoringID   string `json:"monitoringID,omitempty"`
	MonitoringName string `json:"monitoringName,omitempty"`
	Metric         string `json:"metric,omitempty"`
	View           string `json:"view,omitempty"`
}

type CapabilityConstraint struct {
	Type  string
	Name  string
	Value string
}

type CapabilityDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
}

func (c *CapabilityDetails) Clone() *CapabilityDetails {
	clone := &CapabilityDetails{}
	if c.Creation != nil {
		clone.Creation = c.Creation.Clone()
	}
	return clone
}

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
	/*
		if c.Demux != nil {
		demuxloop:
			for _, str := range *c.Demux {
				if c.DeepCompareSlice(str, a.Demux) {
					continue demuxloop
				}
				return false
			}
		}
		if a.Demux != nil {
		revdemuxloop:
			for _, str := range *a.Demux {
				if c.DeepCompareSlice(str, c.Demux) {
					continue revdemuxloop
				}
				return false
			}
		}
		if c.Constraints != nil {
		constraintloop:
			for _, cstr := range *c.Constraints {
				if cstr.DeepCompareSlice(a.Constraints) {
					continue constraintloop
				}
				return false
			}
		}
		if a.Constraints != nil {
		revconstraintloop:
			for _, cstr := range *a.Constraints {
				if cstr.DeepCompareSlice(c.Constraints) {
					continue revconstraintloop
				}
				return false
			}
		}
	*/
	return true
}

func (c *Capability) DeepCompareSlice(s string, a *[]string) bool {
	if a == nil || *a == nil {
		return false
	}

	for _, str := range *a {
		if s == str {
			return true
		}
	}
	return false
}

func (c *CapabilityConstraint) DeepCompare(a *CapabilityConstraint) bool {
	if a == nil {
		return false
	}

	if c.Type != a.Type || c.Name != a.Type || c.Value != a.Value {
		return false
	}
	return true
}

func (c *CapabilityConstraint) DeepCompareSlice(a *[]CapabilityConstraint) bool {
	if a == nil || *a == nil {
		return false
	}

	for _, cstr := range *a {
		if c.DeepCompare(&cstr) {
			return true
		}
	}
	return false
}

func NewCapabilityRequest() Request {
	return Request{
		Flags: &Flags{},
		Capability: &Capability{
			Demux:       &[]string{},
			Constraints: &[]CapabilityConstraint{},
		},
	}
}

func NewCapabilityFilter() Request {
	return Request{
		Filter: &Filter{
			Capability: &CapabilityFilter{},
		},
	}
}

func NewCapabilityResult() Result {
	return Result{
		Errors:       &[]string{},
		Capabilities: &[]Capability{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
