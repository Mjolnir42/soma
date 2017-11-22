/*-
 * Copyright (c) 2015-2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

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

func (s *Server) Clone() Server {
	return Server{
		ID:         s.ID,
		AssetID:    s.AssetID,
		Datacenter: s.Datacenter,
		Location:   s.Location,
		Name:       s.Name,
		IsOnline:   s.IsOnline,
		IsDeleted:  s.IsDeleted,
		Details:    s.Details.Clone(),
	}
}

type ServerDetails struct {
	Creation *DetailsCreation
	/*
		Nodes     []string `json:"nodes,omitempty"`
	*/
}

func (d *ServerDetails) Clone() *ServerDetails {
	return &ServerDetails{
		Creation: d.Creation.Clone(),
	}
}

type ServerFilter struct {
	IsOnline   bool   `json:"isOnline,omitempty"`
	NotOnline  bool   `json:"notOnline,omitempty"`
	Deleted    bool   `json:"Deleted,omitempty"`
	NotDeleted bool   `json:"notDeleted,omitempty"`
	Datacenter string `json:"datacenter,omitempty"`
	Name       string `json:"name,omitempty"`
	AssetID    uint64 `json:"assetID,omitempty"`
}

func (s *Server) DeepCompare(a *Server) bool {
	if s.ID != a.ID || s.AssetID != a.AssetID || s.Datacenter != a.Datacenter ||
		s.Location != a.Location || s.Name != a.Name || s.IsOnline != a.IsOnline ||
		s.IsDeleted != a.IsDeleted {
		return false
	}
	return true
}

func NewServerRequest() Request {
	return Request{
		Flags:  &Flags{},
		Server: &Server{},
	}
}

func NewServerFilter() Request {
	return Request{
		Filter: &Filter{
			Server: &ServerFilter{},
		},
	}
}

func NewServerResult() Result {
	return Result{
		Errors:  &[]string{},
		Servers: &[]Server{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
