/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// Section represents an object handler inside a category scope
type Section struct {
	ID       string          `json:"id,omitempty"`
	Name     string          `json:"name,omitempty"`
	Category string          `json:"category,omitempty"`
	Details  *SectionDetails `json:"details,omitempty"`
}

// Clone returns a copy of s
func (s *Section) Clone() Section {
	clone := Section{
		ID:       s.ID,
		Name:     s.Name,
		Category: s.Category,
	}
	if s.Details != nil {
		clone.Details = s.Details.Clone()
	}
	return clone
}

// SectionDetails contains metadata about a Section
type SectionDetails struct {
	Creation *DetailsCreation `json:"details,omitempty"`
}

// Clone returns a copy of s
func (s *SectionDetails) Clone() *SectionDetails {
	clone := &SectionDetails{}
	if s.Creation != nil {
		clone.Creation = s.Creation.Clone()
	}
	return clone
}

// SectionFilter represents parts of a section that a section
// can be searched by
type SectionFilter struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Category string `json:"category,omitempty"`
}

// NewSectionRequest returns a new Request with fields preallocated
// for filling in Section data, ensuring no nilptr-deref takes place.
func NewSectionRequest() Request {
	return Request{
		Flags:   &Flags{},
		Section: &Section{},
	}
}

// NewSectionFilter returns a new Request with fields preallocated
// for filling in a Section filter, ensuring no nilptr-deref takes place.
func NewSectionFilter() Request {
	return Request{
		Filter: &Filter{
			Section: &SectionFilter{},
		},
	}
}

// NewSectionResult returns a new Result with fields preallocated
// for filling in Section data, ensuring no nilptr-deref takes place.
func NewSectionResult() Result {
	return Result{
		Errors:   &[]string{},
		Sections: &[]Section{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
