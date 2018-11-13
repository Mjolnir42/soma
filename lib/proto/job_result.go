/*-
 * Copyright (c) 2018, Jörg Pernfuß
 * Copyright (c) 2018, 1&1 IONOS SE
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// JobResult defines a possible processing result of a job
type JobResult struct {
	ID      string            `json:"id,omitempty"`
	Name    string            `json:"name,omitempty"`
	Details *JobResultDetails `json:"details,omitempty"`
}

// Clone returns a copy of j
func (j *JobResult) Clone() JobResult {
	clone := JobResult{
		ID:   j.ID,
		Name: j.Name,
	}
	if j.Details != nil {
		clone.Details = j.Details.Clone()
	}
	return clone
}

// JobResultDetails contains metadata about a JobResult
type JobResultDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
}

// Clone returns a copy of j
func (j *JobResultDetails) Clone() *JobResultDetails {
	clone := &JobResultDetails{}
	if j.Creation != nil {
		clone.Creation = j.Creation.Clone()
	}
	return clone
}

// JobResultFilter represents parts of a JobResult that it can be
// searched by
type JobResultFilter struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// NewJobResultRequest returns a new Request with fields preallocated
// for filling in JobResult data, ensuring no nilptr-deref takes place.
func NewJobResultRequest() Request {
	return Request{
		Flags:     &Flags{},
		JobResult: &JobResult{},
	}
}

// NewJobResultFilter returns a new Request with fields preallocated
// for filling in a JobResult filter, ensuring no nilptr-deref takes place.
func NewJobResultFilter() Request {
	return Request{
		Filter: &Filter{
			JobResult: &JobResultFilter{},
		},
	}
}

// NewJobResultResult returns a new Result with fields preallocated
// for filling in JobResult data, ensuring no nilptr-deref takes place.
func NewJobResultResult() Result {
	return Result{
		Errors:     &[]string{},
		JobResults: &[]JobResult{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
