/*-
 * Copyright (c) 2015-2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Team struct {
	ID       string       `json:"id,omitempty"`
	Name     string       `json:"name,omitempty"`
	LdapID   string       `json:"ldapId,omitempty"`
	IsSystem bool         `json:"isSystem,omitempty"`
	Details  *TeamDetails `json:"details,omitempty"`
}

func (t *Team) Clone() Team {
	clone := Team{
		ID:       t.ID,
		Name:     t.Name,
		LdapID:   t.LdapID,
		IsSystem: t.IsSystem,
	}
	if t.Details != nil {
		clone.Details = t.Details.Clone()
	}
	return clone
}

type TeamDetails struct {
	Creation       *DetailsCreation `json:"creation,omitempty"`
	DictionaryID   string           `json:"dictionaryID,omitempty"`
	DictionaryName string           `json:"dictionaryName,omitempty"`
}

func (t *TeamDetails) Clone() *TeamDetails {
	clone := &TeamDetails{
		DictionaryID:   t.DictionaryID,
		DictionaryName: t.DictionaryName,
	}
	if t.Creation != nil {
		clone.Creation = t.Creation.Clone()
	}
	return clone
}

type TeamFilter struct {
	Name     string `json:"name,omitempty"`
	LdapID   string `json:"ldapId,omitempty"`
	IsSystem bool   `json:"isSystem,omitempty"`
}

func (t *Team) DeepCompare(a *Team) bool {
	if t.ID != a.ID || t.Name != a.Name || t.LdapID != a.LdapID || t.IsSystem != a.IsSystem {
		return false
	}
	return true
}

func NewTeamRequest() Request {
	return Request{
		Flags: &Flags{},
		Team:  &Team{},
	}
}

func NewTeamFilter() Request {
	return Request{
		Filter: &Filter{
			Team: &TeamFilter{},
		},
	}
}

func NewTeamResult() Result {
	return Result{
		Errors: &[]string{},
		Teams:  &[]Team{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
