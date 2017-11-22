/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Details struct {
	CreatedAt    string         `json:"createdAt,omitempty"`
	CreatedBy    string         `json:"createdBy,omitempty"`
	Server       Server         `json:"server,omitempty"`
	CheckConfigs *[]CheckConfig `json:"checkConfigs,omitempty"`
}

func (d *Details) Clone() *Details {
	clone := &Details{
		CreatedAt: d.CreatedAt,
		CreatedBy: d.CreatedBy,
		Server:    d.Server.Clone(),
	}
	// XXX CheckConfigs
	return clone
}

type DetailsCreation struct {
	CreatedAt string `json:"createdAt,omitempty"`
	CreatedBy string `json:"createdBy,omitempty"`
}

func (d *DetailsCreation) Clone() *DetailsCreation {
	return &DetailsCreation{
		CreatedAt: d.CreatedAt,
		CreatedBy: d.CreatedBy,
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
