/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// Metric is a definition of a measurement metric
type Metric struct {
	Path        string           `json:"path,omitempty"`
	Unit        string           `json:"unit,omitempty"`
	Description string           `json:"description,omitempty"`
	Packages    *[]MetricPackage `json:"packages,omitempty"`
	Details     *MetricDetails   `json:"details,omitempty"`
}

// Clone returns a copy of m
func (m *Metric) Clone() Metric {
	clone := Metric{
		Path:        m.Path,
		Unit:        m.Unit,
		Description: m.Description,
		Packages:    &[]MetricPackage{},
	}
	if m.Details != nil {
		clone.Details = m.Details.Clone()
	}
	if m.Packages != nil {
		for i := range *m.Packages {
			*clone.Packages = append(*clone.Packages, (*m.Packages)[i].Clone())
		}
	}
	if len(*clone.Packages) == 0 {
		clone.Packages = nil
	}
	return clone
}

// DeepCompare returns true if m and a are equal, excluding details
func (m *Metric) DeepCompare(a *Metric) bool {
	if m.Path != a.Path || m.Unit != a.Unit || m.Description != a.Description {
		return false
	}

packageloop:
	for _, pkg := range *m.Packages {
		if pkg.DeepCompareSlice(a.Packages) {
			continue packageloop
		}
		return false
	}

revpackageloop:
	for _, pkg := range *a.Packages {
		if pkg.DeepCompareSlice(m.Packages) {
			continue revpackageloop
		}
		return false
	}

	return true
}

// MetricDetails contains metadata about a metric
type MetricDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
}

// Clone returns a copy of m
func (m *MetricDetails) Clone() *MetricDetails {
	clone := &MetricDetails{}
	if m.Creation != nil {
		clone.Creation = m.Creation.Clone()
	}
	return clone
}

// MetricPackage specifies a package that can provide a metric
type MetricPackage struct {
	Provider string `json:"provider,omitempty"`
	Name     string `json:"name,omitempty"`
}

// Clone returns a copy of m
func (m *MetricPackage) Clone() MetricPackage {
	return MetricPackage{
		Provider: m.Provider,
		Name:     m.Name,
	}
}

// DeepCompare returns true if m and a are equal, excluding details
func (m *MetricPackage) DeepCompare(a *MetricPackage) bool {
	if m.Provider != a.Provider || m.Name != a.Name {
		return false
	}
	return true
}

// DeepCompareSlice returns true if m is equal to a metric package contained
// in slice a
func (m *MetricPackage) DeepCompareSlice(a *[]MetricPackage) bool {
	if a == nil || *a == nil {
		return false
	}
	for _, pkg := range *a {
		if m.DeepCompare(&pkg) {
			return true
		}
	}
	return false
}

// MetricFilter represents parts of a metric that a metric can be
// searched by
type MetricFilter struct {
	Unit     string `json:"unit,omitempty"`
	Provider string `json:"provider,omitempty"`
	Package  string `json:"package,omitempty"`
}

// NewMetricRequest returns a new Request with fields preallocated
// for filling in Metric data, ensuring no nilptr-deref takes place.
func NewMetricRequest() Request {
	return Request{
		Flags:  &Flags{},
		Metric: &Metric{},
	}
}

// NewMetricFilter returns a new Request with fields preallocated
// for filling in a Metric filter, ensuring no nilptr-deref takes place.
func NewMetricFilter() Request {
	return Request{
		Filter: &Filter{
			Metric: &MetricFilter{},
		},
	}
}

// NewMetricResult returns a new Result with fields preallocated
// for filling in Metric data, ensuring no nilptr-deref takes place.
func NewMetricResult() Result {
	return Result{
		Errors:  &[]string{},
		Metrics: &[]Metric{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
