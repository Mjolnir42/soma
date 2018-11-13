/*-
 * Copyright (c) 2015-2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// Job is the specification of an asynchronously executed SOMA job
type Job struct {
	ID           string      `json:"id,omitempty"`
	Status       string      `json:"status,omitempty"`
	Result       string      `json:"result,omitempty"`
	Type         string      `json:"type,omitempty"`
	Serial       int         `json:"serial,omitempty"`
	RepositoryID string      `json:"repositoryId,omitempty"`
	UserID       string      `json:"userId,omitempty"`
	TeamID       string      `json:"teamId,omitempty"`
	TsQueued     string      `json:"queued,omitempty"`
	TsStarted    string      `json:"started,omitempty"`
	TsFinished   string      `json:"finished,omitempty"`
	Error        string      `json:"error,omitempty"`
	Details      *JobDetails `json:"details,omitempty"`
}

// Clone returns a cope of j
func (j *Job) Clone() Job {
	clone := Job{
		ID:           j.ID,
		Status:       j.Status,
		Result:       j.Result,
		Type:         j.Type,
		Serial:       j.Serial,
		RepositoryID: j.RepositoryID,
		UserID:       j.UserID,
		TeamID:       j.TeamID,
		TsQueued:     j.TsQueued,
		TsStarted:    j.TsStarted,
		TsFinished:   j.TsFinished,
	}
	if j.Details != nil {
		clone.Details = j.Details.Clone()
	}
	return clone
}

// JobDetails contains metadata about a JobType
type JobDetails struct {
	CreatedAt     string `json:"createdAt,omitempty"`
	CreatedBy     string `json:"createdBy,omitempty"`
	Specification string `json:"specification,omitempty"`
}

// Clone returns a cope of j
func (j *JobDetails) Clone() *JobDetails {
	return &JobDetails{
		CreatedAt:     j.CreatedAt,
		CreatedBy:     j.CreatedBy,
		Specification: j.Specification,
	}
}

// JobFilter represents parts of a Job that it can be searched by
type JobFilter struct {
	User   string   `json:"user,omitempty"`
	Team   string   `json:"team,omitempty"`
	Status string   `json:"status,omitempty"`
	Result string   `json:"result,omitempty"`
	Since  string   `json:"since,omitempty"`
	IDList []string `json:"idlist,omitempty"`
}

// NewJobFilter returns a new Request with fields preallocated
// for filling in a Job filter, ensuring no nilptr-deref takes place.
func NewJobFilter() Request {
	return Request{
		Flags: &Flags{},
		Filter: &Filter{
			Job: &JobFilter{},
		},
	}
}

// NewJobResult returns a new Result with fields preallocated
// for filling in Job data, ensuring no nilptr-deref takes place.
func NewJobResult() Result {
	return Result{
		Errors: &[]string{},
		Jobs:   &[]Job{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
