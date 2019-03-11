/*-
 * Copyright (c) 2019, Jörg Pernfuß
 * Copyright (c) 2019, 1&1 IONOS SE
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// Admin describes an admin account
type Admin struct {
	ID          string            `json:"id,omitempty"`
	Name        string            `json:"name,omitempty"`
	UserID      string            `json:"userID,omitempty"`
	UserName    string            `json:"userName,omitempty"`
	Details     *AdminDetails     `json:"details,omitempty"`
	Credentials *AdminCredentials `json:"credentials,omitempty"`
}

type AdminCredentials struct {
	Reset          bool   `json:"reset,omitempty"`
	ForcedPassword string `json:"forcedPassword,omitempty"`
}

type AdminDetails struct {
	Creation       *DetailsCreation `json:"creation,omitempty"`
	DictionaryID   string           `json:"dictionaryID,omitempty"`
	DictionaryName string           `json:"dictionaryName,omitempty"`
}

func NewAdminRequest() Request {
	return Request{
		Flags: &Flags{},
		Admin: &Admin{},
	}
}

func NewAdminResult() Result {
	return Result{
		Errors: &[]string{},
		Admins: &[]Admin{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
