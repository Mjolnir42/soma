/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2017, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/auth"
	"github.com/mjolnir42/soma/lib/proto"
)

// sendMsgResult is the output function for all requests that did not
// fail input validation and got processes by the application.
func sendMsgResult(w *http.ResponseWriter, r *msg.Result) {
	var (
		bjson  []byte
		err    error
		k      auth.Kex
		result proto.Result
	)

	// this is central error command, proceeding to log
	if r.Error != nil {
		log.Printf(msg.LogStrErr, r.Section, r.Action, r.Code, r.Error.Error())
	}

	switch r.Section {
	case msg.SectionSupervisor:
		switch r.Action {
		case msg.ActionToken:
			switch r.Super.Task {
			case msg.TaskInvalidate:
				result = proto.NewResult()
			}
		}

	case msg.SectionSystemOperation:
		switch r.Action {
		case msg.ActionRevokeTokens:
			switch r.Super.Task {
			case msg.TaskInvalidateAccount:
				result = proto.NewResult()
			}
		}

	case msg.SectionAction:
		result = proto.NewActionResult()
		*result.Actions = append(*result.Actions, r.ActionObj...)

	case msg.SectionAttribute:
		result = proto.NewAttributeResult()
		*result.Attributes = append(*result.Attributes, r.Attribute...)

	case msg.SectionBucket:
		result = proto.NewBucketResult()
		*result.Buckets = append(*result.Buckets, r.Bucket...)

	case msg.SectionCapability:
		result = proto.NewCapabilityResult()
		*result.Capabilities = append(*result.Capabilities, r.Capability...)

	case msg.SectionCategory:
		result = proto.NewCategoryResult()
		*result.Categories = append(*result.Categories, r.Category...)

	case msg.SectionCheckConfig:
		result = proto.NewCheckConfigResult()
		*result.CheckConfigs = append(*result.CheckConfigs, r.CheckConfig...)

	case msg.SectionCluster:
		result = proto.NewClusterResult()
		*result.Clusters = append(*result.Clusters, r.Cluster...)

	case msg.SectionDatacenter:
		result = proto.NewDatacenterResult()
		*result.Datacenters = append(*result.Datacenters, r.Datacenter...)

	// XXX below are legacy definitions
	case `kex`:
		k = r.Super.Kex
		if bjson, err = json.Marshal(&k); err != nil {
			log.Printf(msg.LogStrErr, r.Section, r.Action, r.Code, err.Error())
			dispatchInternalError(w, nil)
			return
		}
		goto dispatchJSON
	case `bootstrap`, `activate`, `token`, `password`:
		// for these request types, response codes are masked. they are also
		// not behind basic auth
		switch r.Code {
		case 200:
			if r.Super.Verdict == 200 {
				log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 200)
				goto dispatchOCTET
			}
			log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 403)
			dispatchForbidden(w, nil)
		case 406:
			log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 406)
			dispatchConflict(w, r.Error)
		default:
			log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 403)
			dispatchForbidden(w, nil)
		}
		return
	case `permission`:
		result = proto.NewPermissionResult()
		*result.Permissions = append(*result.Permissions, r.Permission...)
	case `right`:
		result = proto.NewGrantResult()
		*result.Grants = append(*result.Grants, r.Grant...)
	case `section`:
		result = proto.NewSectionResult()
		*result.Sections = append(*result.Sections, r.SectionObj...)
	case `environment`:
		result = proto.NewEnvironmentResult()
		*result.Environments = append(*result.Environments, r.Environment...)
	case `job`:
		result = proto.NewJobResult()
		*result.Jobs = append(*result.Jobs, r.Job...)
	case `tree`:
		result = proto.NewTreeResult()
		*result.Tree = r.Tree
	case `runtime`:
		switch r.Action {
		case `instance_list_all`:
			result = proto.NewInstanceResult()
			*result.Instances = append(*result.Instances, r.Instance...)
		case `job_list_all`:
			result = proto.NewJobResult()
			*result.Jobs = append(*result.Jobs, r.Job...)
		default:
			result = proto.NewSystemOperationResult()
			*result.SystemOperations = append(*result.SystemOperations, r.System...)
		}
	case `instance`:
		result = proto.NewInstanceResult()
		*result.Instances = append(*result.Instances, r.Instance...)
	case `workflow`:
		result = proto.NewWorkflowResult()
		*result.Workflows = append(*result.Workflows, r.Workflow...)
	case `state`:
		result = proto.NewStateResult()
		*result.States = append(*result.States, r.State...)
	case `entity`:
		result = proto.NewEntityResult()
		*result.Entities = append(*result.Entities, r.Entity...)
	case `monitoringsystem`:
		result = proto.NewMonitoringResult()
		*result.Monitorings = append(*result.Monitorings, r.Monitoring...)
	case `node-mgmt`:
		result = proto.NewNodeResult()
		*result.Nodes = append(*result.Nodes, r.Node...)
	default:
		log.Printf(msg.LogStrErr, r.Section, r.Action, 0, `Result from unhandled subsystem`)
		dispatchInternalError(w, nil)
		return
	}

	switch r.Code {
	case 200:
		log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 200)
		if r.Error != nil {
			result.Error(r.Error)
		}
		result.OK()
	case 202:
		log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 202)
		result.JobID = r.JobID
		result.Accepted()
	case 400:
		log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 200)
		result.BadRequest(r.Error)
	case 403:
		log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 200)
		result.Forbidden(r.Error)
	case 404:
		log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 200)
		result.NotFoundErr(r.Error)
	case 500:
		log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 200)
		result.Error(r.Error)
	case 501:
		log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 200)
		result.NotImplemented()
	case 503:
		log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 200)
		result.Unavailable()
	default:
		log.Printf(msg.LogStrErr, r.Section, r.Action, r.Code, `Unhandled internal result code`)
		dispatchInternalError(w, nil)
		return
	}
	goto buildJSON

dispatchOCTET:
	dispatchOctetReply(w, &r.Super.Encrypted.Data)
	return

buildJSON:
	if bjson, err = json.Marshal(&result); err != nil {
		log.Printf(msg.LogStrErr, r.Section, r.Action, r.Code, err)
		dispatchInternalError(w, nil)
		return
	}

dispatchJSON:
	dispatchJSONReply(w, &bjson)
	return
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
