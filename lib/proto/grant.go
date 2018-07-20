/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Grant struct {
	ID            string           `json:"id"`
	RecipientType string           `json:"recipientType"`
	RecipientID   string           `json:"recipientId"`
	PermissionID  string           `json:"permissionId"`
	Category      string           `json:"category"`
	ObjectType    string           `json:"objectType"`
	ObjectID      string           `json:"objectId"`
	Details       *DetailsCreation `json:"details,omitempty"`
}

func (g *Grant) Clone() Grant {
	clone := Grant{
		ID:            g.ID,
		RecipientType: g.RecipientType,
		RecipientID:   g.RecipientID,
		PermissionID:  g.PermissionID,
		Category:      g.Category,
		ObjectType:    g.ObjectType,
		ObjectID:      g.ObjectID,
	}
	if g.Details != nil {
		clone.Details = g.Details.Clone()
	}
	return clone
}

type GrantFilter struct {
	RecipientType string `json:"recipientType"`
	RecipientID   string `json:"recipientId"`
	PermissionID  string `json:"permissionId"`
	Category      string `json:"category"`
	ObjectType    string `json:"objectType"`
	ObjectID      string `json:"objectId"`
}

func NewGrantRequest() Request {
	return Request{
		Flags: &Flags{},
		Grant: &Grant{},
	}
}

func NewGrantFilter() Request {
	return Request{
		Filter: &Filter{
			Grant: &GrantFilter{},
		},
	}
}

func NewGrantResult() Result {
	return Result{
		Errors: &[]string{},
		Grants: &[]Grant{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
