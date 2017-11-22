/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Action struct {
	ID          string           `json:"ID,omitempty"`
	Name        string           `json:"name,omitempty"`
	SectionID   string           `json:"sectionID,omitempty"`
	SectionName string           `json:"sectionName,omitempty"`
	Category    string           `json:"category,omitempty"`
	Details     *DetailsCreation `json:"details,omitempty"`
}

func (a *Action) Clone() Action {
	return Action{
		ID:          a.ID,
		Name:        a.Name,
		SectionID:   a.SectionID,
		SectionName: a.SectionName,
		Category:    a.Category,
		Details:     a.Details.Clone(),
	}
}

type ActionFilter struct {
	Name      string `json:"name,omitempty"`
	SectionID string `json:"sectionID,omitempty"`
}

func NewActionRequest() Request {
	return Request{
		Flags:  &Flags{},
		Action: &Action{},
	}
}

func NewActionFilter() Request {
	return Request{
		Filter: &Filter{
			Action: &ActionFilter{},
		},
	}
}

func NewActionResult() Result {
	return Result{
		Errors:  &[]string{},
		Actions: &[]Action{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
