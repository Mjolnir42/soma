/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// Category represents a permission scope
type Category struct {
	Name    string           `json:"name,omitempty"`
	Details *CategoryDetails `json:"details,omitempty"`
}

// Clone returns a copy of c
func (c *Category) Clone() Category {
	clone := Category{
		Name: c.Name,
	}
	if c.Details != nil {
		clone.Details = c.Details.Clone()
	}
	return clone
}

// CategoryDetails contains metadata about a Category
type CategoryDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
}

// Clone returns a copy of c
func (c *CategoryDetails) Clone() *CategoryDetails {
	clone := &CategoryDetails{}
	if c.Creation != nil {
		clone.Creation = c.Creation.Clone()
	}
	return clone
}

// NewCategoryRequest returns a new Request with fields preallocated
// for filling in Category data, ensuring no nilptr-deref takes place.
func NewCategoryRequest() Request {
	return Request{
		Flags:    &Flags{},
		Category: &Category{},
	}
}

// NewCategoryResult returns a new Result with fields preallocated
// for filling in Category data, ensuring no nilptr-deref takes place.
func NewCategoryResult() Result {
	return Result{
		Errors:     &[]string{},
		Categories: &[]Category{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
