/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type CheckInstance struct {
	InstanceID            string `json:"instanceId,omitempty"`
	CheckID               string `json:"checkId,omitempty"`
	ConfigID              string `json:"configId,omitempty"`
	InstanceConfigID      string `json:"instanceConfigID,omitempty"`
	Version               uint64 `json:"version"`
	ConstraintHash        string `json:"constraintHash,omitempty"`
	ConstraintValHash     string `json:"constraintValHash,omitempty"`
	InstanceSvcCfgHash    string `json:"instanceSvcCfghash,omitempty"`
	InstanceService       string `json:"instanceService,omitempty"`
	InstanceServiceConfig string `json:"instanceServiceCfg,omitempty"`
}

func (t *CheckInstance) DeepCompare(a *CheckInstance) bool {
	if t.InstanceID != a.InstanceID || t.CheckID != a.CheckID || t.ConfigID != a.ConfigID ||
		t.ConstraintHash != a.ConstraintHash || t.ConstraintValHash != a.ConstraintValHash ||
		t.InstanceSvcCfgHash != a.InstanceSvcCfgHash || t.InstanceService != a.InstanceService {
		// - InstanceConfigID is a randomly generated uuid on every instance calculation
		// - Version is incremented on every instance calculation
		// - InstanceServiceConfig is compared as deploymentdetails.Service
		return false
	}
	return true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
