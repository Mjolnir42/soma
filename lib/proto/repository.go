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
	ID         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	TeamID     string      `json:"teamId,omitempty"`
	IsDeleted  bool        `json:"isDeleted,omitempty"`
	IsActive   bool        `json:"isActive,omitempty"`
	Members    *[]Bucket   `json:"members,omitempty"`
	Details    *Details    `json:"details,omitempty"`
	Properties *[]Property `json:"properties,omitempty"`
}

// Clone function
func (r *Repository) Clone() Repository {
	clone := Repository{
		ID:        r.ID,
		Name:      r.Name,
		TeamID:    r.TeamID,
		IsDeleted: r.IsDeleted,
		IsActive:  r.IsActive,
	}
	return clone
}

type RepositoryFilter struct {
	ID        string `json:"ID,omitempty"`
	Name      string `json:"name,omitempty"`
	TeamID    string `json:"teamId,omitempty"`
	IsDeleted bool   `json:"isDeleted,omitempty"`
	IsActive  bool   `json:"isActive,omitempty"`
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
