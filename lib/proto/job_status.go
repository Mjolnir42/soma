/*-
 * Copyright (c) 2018, Jörg Pernfuß
 * Copyright (c) 2018, 1&1 IONOS SE
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// JobStatus defines the processing state of a job
type JobStatus struct {
	ID      string            `json:"id,omitempty"`
	Name    string            `json:"name,omitempty"`
	Details *JobStatusDetails `json:"details,omitempty"`
}

// Clone returns a copy of j
func (j *JobStatus) Clone() JobStatus {
	clone := JobStatus{
		ID:   j.ID,
		Name: j.Name,
	}
	if j.Details != nil {
		clone.Details = j.Details.Clone()
	}
	return clone
}

// JobStatusDetails contains metadata about a JobStatus
type JobStatusDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
}

// Clone returns a copy of j
func (j *JobStatusDetails) Clone() *JobStatusDetails {
	clone := &JobStatusDetails{}
	if j.Creation != nil {
		clone.Creation = j.Creation.Clone()
	}
	return clone
}

// JobStatusFilter represents parts of a JobStatus that it can be
// searched by
type JobStatusFilter struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// NewJobStatusRequest returns a new Request with fields preallocated
// for filling in JobStatus data, ensuring no nilptr-deref takes place.
func NewJobStatusRequest() Request {
	return Request{
		Flags:     &Flags{},
		JobStatus: &JobStatus{},
	}
}

// NewJobStatusFilter returns a new Request with fields preallocated
// for filling in a JobStatus filter, ensuring no nilptr-deref takes place.
func NewJobStatusFilter() Request {
	return Request{
		Filter: &Filter{
			JobStatus: &JobStatusFilter{},
		},
	}
}

// NewJobStatusResult returns a new Result with fields preallocated
// for filling in JobStatus data, ensuring no nilptr-deref takes place.
func NewJobStatusResult() Result {
	return Result{
		Errors:    &[]string{},
		JobStatus: &[]JobStatus{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
