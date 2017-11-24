/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Cluster struct {
	ID          string      `json:"ID,omitempty"`
	Name        string      `json:"name,omitempty"`
	BucketID    string      `json:"bucketID,omitempty"`
	ObjectState string      `json:"objectState,omitempty"`
	TeamID      string      `json:"teamID,omitempty"`
	Members     *[]Node     `json:"members,omitempty"`
	Details     *Details    `json:"details,omitempty"`
	Properties  *[]Property `json:"properties,omitempty"`
}

func (c *Cluster) Clone() Cluster {
	clone := Cluster{
		ID:          c.ID,
		Name:        c.Name,
		BucketID:    c.BucketID,
		ObjectState: c.ObjectState,
		TeamID:      c.TeamID,
	}
	if c.Details != nil {
		clone.Details = c.Details.Clone()
	}
	if c.Members != nil {
		*clone.Members = make([]Node, len(*c.Members))
		for i := range *c.Members {
			(*clone.Members)[i] = (*c.Members)[i].Clone()
		}
	}
	if c.Properties != nil {
		*clone.Properties = make([]Property, len(*c.Properties))
		for i := range *c.Properties {
			(*clone.Properties)[i] = (*c.Properties)[i].Clone()
		}
	}
	return clone
}

type ClusterFilter struct {
	Name     string `json:"name,omitempty"`
	BucketID string `json:"bucketID,omitempty"`
	TeamID   string `json:"teamID,omitempty"`
}

func (c *Cluster) DeepCompare(a *Cluster) bool {
	if a == nil {
		return false
	}

	if c.ID != a.ID ||
		c.Name != a.Name ||
		c.BucketID != a.BucketID ||
		c.ObjectState != a.ObjectState ||
		c.TeamID != a.TeamID {
		return false
	}

	if c.Members != nil && a.Members != nil {
	member:
		for i, _ := range *c.Members {
			for j, _ := range *a.Members {
				if (*c.Members)[i].ID == (*a.Members)[j].ID {
					continue member
				}
			}
			return false
		}
		return true
	} else if c.Members != nil && a.Members == nil {
		return false
	} else if c.Members == nil && a.Members != nil {
		return false
	}
	return true
}

func NewClusterRequest() Request {
	return Request{
		Flags:   &Flags{},
		Cluster: &Cluster{},
	}
}

func NewClusterFilter() Request {
	return Request{
		Filter: &Filter{
			Cluster: &ClusterFilter{},
		},
	}
}

func NewClusterResult() Result {
	return Result{
		Errors:   &[]string{},
		Clusters: &[]Cluster{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
