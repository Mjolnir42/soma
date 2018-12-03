/*-
 * Copyright (c) 2015-2018, 1&1 Internet SE
 * Copyright (c) 2015-2018, Jörg Pernfuß <code.jpe@gmail.com>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

type Repository struct {
	ID         string             `json:"id,omitempty"`
	Name       string             `json:"name,omitempty"`
	TeamID     string             `json:"teamId,omitempty"`
	IsDeleted  bool               `json:"isDeleted,omitempty"`
	IsActive   bool               `json:"isActive,omitempty"`
	Members    *[]Bucket          `json:"members,omitempty"`
	Details    *RepositoryDetails `json:"details,omitempty"`
	Properties *[]Property        `json:"properties,omitempty"`
}

// Clone returns a copy of r
func (r *Repository) Clone() Repository {
	clone := Repository{
		ID:        r.ID,
		Name:      r.Name,
		TeamID:    r.TeamID,
		IsDeleted: r.IsDeleted,
		IsActive:  r.IsActive,
	}
	if r.Members != nil {
		b := make([]Bucket, 0)
		for i := range *r.Members {
			b = append(b, (*r.Members)[i].Clone())
		}
		clone.Members = &b
	}
	if r.Details != nil {
		clone.Details = r.Details.Clone()
	}
	if r.Properties != nil && *r.Properties != nil {
		p := make([]Property, 0)
		for i := range *r.Properties {
			p = append(p, (*r.Properties)[i].Clone())
		}
		clone.Properties = &p
	}
	return clone
}

// RepositoryDetails contains metadata about a Repository
type RepositoryDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
}

// Clone returns a copy of r
func (r *RepositoryDetails) Clone() *RepositoryDetails {
	clone := &RepositoryDetails{}
	if r.Creation != nil {
		clone.Creation = r.Creation.Clone()
	}
	return clone
}

type RepositoryFilter struct {
	ID                string `json:"ID,omitempty"`
	Name              string `json:"name,omitempty"`
	TeamID            string `json:"teamId,omitempty"`
	IsDeleted         bool   `json:"isDeleted"`
	IsActive          bool   `json:"isActive"`
	FilterOnIsDeleted bool   `json:"filterOnIsDeleted"`
	FilterOnIsActive  bool   `json:"filterOnIsActive"`
}

func NewRepositoryRequest() Request {
	return Request{
		Flags:      &Flags{},
		Repository: &Repository{},
	}
}

func NewRepositoryFilter() Request {
	return Request{
		Filter: &Filter{
			Repository: &RepositoryFilter{},
		},
	}
}

func NewRepositoryResult() Result {
	return Result{
		Errors:       &[]string{},
		Repositories: &[]Repository{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
