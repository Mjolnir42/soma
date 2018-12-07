/*-
 * Copyright (c) 2015-2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

// Bucket type
type Bucket struct {
	ID             string      `json:"ID,omitempty"`
	Name           string      `json:"name,omitempty"`
	RepositoryID   string      `json:"repositoryID,omitempty"`
	TeamID         string      `json:"teamID,omitempty"`
	Environment    string      `json:"environment,omitempty"`
	IsDeleted      bool        `json:"isDeleted,omitempty"`
	IsFrozen       bool        `json:"isFrozen,omitempty"`
	MemberGroups   *[]Group    `json:"memberGroups,omitempty"`
	MemberClusters *[]Cluster  `json:"memberClusters,omitempty"`
	MemberNodes    *[]Node     `json:"memberNodes,omitempty"`
	Details        *Details    `json:"details,omitempty"`
	Properties     *[]Property `json:"properties,omitempty"`
}

// Clone returns a copy of b without Details
func (b *Bucket) Clone() Bucket {
	clone := Bucket{
		ID:           b.ID,
		Name:         b.Name,
		RepositoryID: b.RepositoryID,
		TeamID:       b.TeamID,
		Environment:  b.Environment,
		IsDeleted:    b.IsDeleted,
		IsFrozen:     b.IsFrozen,
	}
	if b.MemberGroups != nil && *b.MemberGroups != nil {
		g := make([]Group, 0)
		for i := range *b.MemberGroups {
			g = append(g, (*b.MemberGroups)[i].Clone())
		}
		clone.MemberGroups = &g
	}
	if b.MemberClusters != nil && *b.MemberClusters != nil {
		c := make([]Cluster, 0)
		for i := range *b.MemberClusters {
			c = append(c, (*b.MemberClusters)[i].Clone())
		}
		clone.MemberClusters = &c
	}
	if b.MemberNodes != nil && *b.MemberNodes != nil {
		n := make([]Node, 0)
		for i := range *b.MemberNodes {
			n = append(n, (*b.MemberNodes)[i].Clone())
		}
		clone.MemberNodes = &n
	}
	if b.Properties != nil && *b.Properties != nil {
		p := make([]Property, 0)
		for i := range *b.Properties {
			p = append(p, (*b.Properties)[i].Clone())
		}
		clone.Properties = &p
	}
	return clone
}

// BucketFilter type
type BucketFilter struct {
	Name         string `json:"name,omitempty"`
	ID           string `json:"ID,omitempty"`
	RepositoryID string `json:"repositoryID,omitempty"`
	IsDeleted    bool   `json:"isDeleted,omitempty"`
	IsFrozen     bool   `json:"isFrozen,omitempty"`
}

// NewBucketRequest function
func NewBucketRequest() Request {
	return Request{
		Flags:  &Flags{},
		Bucket: &Bucket{},
	}
}

// NewBucketFilter function
func NewBucketFilter() Request {
	return Request{
		Filter: &Filter{
			Bucket: &BucketFilter{},
		},
	}
}

// NewBucketResult function
func NewBucketResult() Result {
	return Result{
		Errors:  &[]string{},
		Buckets: &[]Bucket{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
