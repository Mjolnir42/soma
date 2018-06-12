/*-
 * Copyright (c) 2015-2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Oncall struct {
	ID      string          `json:"id,omitempty"`
	Name    string          `json:"name,omitempty"`
	Number  string          `json:"number,omitempty"`
	Members *[]OncallMember `json:"members,omitempty"`
	Details *OncallDetails  `json:"details,omitempty"`
}

func (o *Oncall) Clone() Oncall {
	clone := Oncall{
		ID:      o.ID,
		Name:    o.Name,
		Number:  o.Number,
		Members: &[]OncallMember{},
		Details: o.Details.Clone(),
	}
	if o.Members != nil {
		for i := range *o.Members {
			*clone.Members = append(*clone.Members, (*o.Members)[i].Clone())
		}
	}
	if len(*clone.Members) == 0 {
		clone.Members = nil
	}
	return clone
}

func (o *Oncall) Sanitize() {
	o.ID = ``
	o.Members = nil
	o.Details = nil
}

type OncallDetails struct {
	DetailsCreation
}

func (d *OncallDetails) Clone() *OncallDetails {
	return &OncallDetails{
		DetailsCreation: *d.DetailsCreation.Clone(),
	}
}

type OncallMember struct {
	UserName string `json:"userName,omitempty"`
	UserID   string `json:"userId,omitempty"`
}

func (m *OncallMember) Clone() OncallMember {
	return OncallMember{
		UserName: m.UserName,
		UserID:   m.UserID,
	}
}

type OncallFilter struct {
	Name   string `json:"name,omitempty"`
	Number string `json:"number,omitempty"`
}

func (p *Oncall) DeepCompare(a *Oncall) bool {
	if p.ID != a.ID || p.Name != a.Name || p.Number != a.Number {
		return false
	}
	return true
}

func NewOncallRequest() Request {
	return Request{
		Flags:  &Flags{},
		Oncall: &Oncall{},
	}
}

func NewOncallFilter() Request {
	return Request{
		Filter: &Filter{
			Oncall: &OncallFilter{},
		},
	}
}

func NewOncallResult() Result {
	return Result{
		Errors:  &[]string{},
		Oncalls: &[]Oncall{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
