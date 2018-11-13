/*-
 * Copyright (c) 2018, Jörg Pernfuß
 * Copyright (c) 2018, 1&1 IONOS SE
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// JobType defines the type of a job
type JobType struct {
	ID      string          `json:"id,omitempty"`
	Name    string          `json:"name,omitempty"`
	Details *JobTypeDetails `json:"details,omitempty"`
}

// Clone returns a copy of j
func (j *JobType) Clone() JobType {
	clone := JobType{
		ID:   j.ID,
		Name: j.Name,
	}
	if j.Details != nil {
		clone.Details = j.Details.Clone()
	}
	return clone
}

// JobTypeDetails contains metadata about a JobType
type JobTypeDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
}

// Clone returns a copy of j
func (j *JobTypeDetails) Clone() *JobTypeDetails {
	clone := &JobTypeDetails{}
	if j.Creation != nil {
		clone.Creation = j.Creation.Clone()
	}
	return clone
}

// JobTypeFilter represents parts of a JobType that it can be
// searched by
type JobTypeFilter struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// NewJobTypeRequest returns a new Request with fields preallocated
// for filling in JobType data, ensuring no nilptr-deref takes place.
func NewJobTypeRequest() Request {
	return Request{
		Flags:   &Flags{},
		JobType: &JobType{},
	}
}

// NewJobTypeFilter returns a new Request with fields preallocated
// for filling in a JobType filter, ensuring no nilptr-deref takes place.
func NewJobTypeFilter() Request {
	return Request{
		Filter: &Filter{
			JobType: &JobTypeFilter{},
		},
	}
}

// NewJobTypeResult returns a new Result with fields preallocated
// for filling in JobType data, ensuring no nilptr-deref takes place.
func NewJobTypeResult() Result {
	return Result{
		Errors:   &[]string{},
		JobTypes: &[]JobType{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
