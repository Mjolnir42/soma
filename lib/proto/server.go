/*-
 * Copyright (c) 2015-2018, Jörg Pernfuß
 * Copyright (c) 2015-2016, 1&1 Internet SE
 * Copyright (c) 2018, 1&1 IONOS SE
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// Server defines a physical server
type Server struct {
	ID         string         `json:"ID,omitempty"`
	AssetID    uint64         `json:"assetID,omitempty"`
	Datacenter string         `json:"datacenter,omitempty"`
	Location   string         `json:"location,omitempty"`
	Name       string         `json:"name,omitempty"`
	IsOnline   bool           `json:"isOnline,omitempty"`
	IsDeleted  bool           `json:"isDeleted,omitempty"`
	Details    *ServerDetails `json:"details,omitempty"`
}

// Clone returns a copy of s
func (s *Server) Clone() Server {
	clone := Server{
		ID:         s.ID,
		AssetID:    s.AssetID,
		Datacenter: s.Datacenter,
		Location:   s.Location,
		Name:       s.Name,
		IsOnline:   s.IsOnline,
		IsDeleted:  s.IsDeleted,
	}
	if s.Details != nil {
		clone.Details = s.Details.Clone()
	}
	return clone
}

// DeepCompare returns true if s and a are equal, excluding details
func (s *Server) DeepCompare(a *Server) bool {
	if s.ID != a.ID || s.AssetID != a.AssetID || s.Datacenter != a.Datacenter ||
		s.Location != a.Location || s.Name != a.Name || s.IsOnline != a.IsOnline ||
		s.IsDeleted != a.IsDeleted {
		return false
	}
	return true
}

// ServerDetails contains metadata about a server
type ServerDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
}

// Clone returns a copy of s
func (s *ServerDetails) Clone() *ServerDetails {
	return &ServerDetails{
		Creation: d.Creation.Clone(),
	}
}

// ServerFilter defines by which attributes a server ca be searched
type ServerFilter struct {
	IsOnline   bool   `json:"isOnline,omitempty"`
	NotOnline  bool   `json:"notOnline,omitempty"`
	Deleted    bool   `json:"Deleted,omitempty"`
	NotDeleted bool   `json:"notDeleted,omitempty"`
	Datacenter string `json:"datacenter,omitempty"`
	Name       string `json:"name,omitempty"`
	AssetID    uint64 `json:"assetID,omitempty"`
}

// NewServerRequest returns a new Request with fields preallocated
// for filling in Server data, ensuring no nilptr-deref takes place.
func NewServerRequest() Request {
	return Request{
		Flags:  &Flags{},
		Server: &Server{},
	}
}

// NewServerFilter returns a new Request with fields preallocated
// for filling in an Server filter, ensuring no nilptr-deref takes place.
func NewServerFilter() Request {
	return Request{
		Filter: &Filter{
			Server: &ServerFilter{},
		},
	}
}

// NewServerResult returns a new Result with fields preallocated
// for filling in Server data, ensuring no nilptr-deref takes place.
func NewServerResult() Result {
	return Result{
		Errors:  &[]string{},
		Servers: &[]Server{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
