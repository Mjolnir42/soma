/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Group struct {
	ID             string      `json:"id,omitempty"`
	Name           string      `json:"name,omitempty"`
	RepositoryID   string      `json:"repositoryId,omitempty"`
	BucketID       string      `json:"bucketId,omitempty"`
	ObjectState    string      `json:"objectState,omitempty"`
	TeamID         string      `json:"teamId,omitempty"`
	MemberGroups   *[]Group    `json:"memberGroups,omitempty"`
	MemberClusters *[]Cluster  `json:"memberClusters,omitempty"`
	MemberNodes    *[]Node     `json:"memberNodes,omitempty"`
	Details        *Details    `json:"details,omitempty"`
	Properties     *[]Property `json:"properties,omitempty"`
}

func (g *Group) Clone() Group {
	clone := Group{
		ID:           g.ID,
		Name:         g.Name,
		RepositoryID: g.RepositoryID,
		BucketID:     g.BucketID,
		ObjectState:  g.ObjectState,
		TeamID:       g.TeamID,
	}
	if g.Details != nil {
		clone.Details = g.Details.Clone()
	}
	if g.MemberGroups != nil {
		*clone.MemberGroups = make([]Group, len(*g.MemberGroups))
		for i := range *g.MemberGroups {
			(*clone.MemberGroups)[i] = (*g.MemberGroups)[i].Clone()
		}
	}
	if g.MemberClusters != nil {
		*clone.MemberClusters = make([]Cluster, len(*g.MemberClusters))
		for i := range *g.MemberClusters {
			(*clone.MemberClusters)[i] = (*g.MemberClusters)[i].Clone()
		}
	}
	if g.MemberNodes != nil {
		*clone.MemberNodes = make([]Node, len(*g.MemberNodes))
		for i := range *g.MemberNodes {
			(*clone.MemberNodes)[i] = (*g.MemberNodes)[i].Clone()
		}
	}
	if g.Properties != nil {
		*clone.Properties = make([]Property, len(*g.Properties))
		for i := range *g.Properties {
			(*clone.Properties)[i] = (*g.Properties)[i].Clone()
		}
	}
	return clone
}

type GroupFilter struct {
	Name         string `json:"name,omitempty"`
	BucketID     string `json:"bucketId,omitempty"`
	RepositoryID string `json:"repositoryID,omitempty"`
}

//
func (g *Group) DeepCompare(a *Group) bool {
	if a == nil {
		return false
	}
	if g.ID != a.ID || g.Name != a.Name || g.BucketID != a.BucketID ||
		g.ObjectState != a.ObjectState || g.TeamID != a.TeamID {
		return false
	}
	return true
}

func NewGroupRequest() Request {
	return Request{
		Flags: &Flags{},
		Group: &Group{},
	}
}

func NewGroupFilter() Request {
	return Request{
		Filter: &Filter{
			Group: &GroupFilter{},
		},
	}
}

func NewGroupResult() Result {
	return Result{
		Errors: &[]string{},
		Groups: &[]Group{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
