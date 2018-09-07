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
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/auth"
	"github.com/mjolnir42/soma/lib/proto"
)

// send is the output function for all requests that did not
// fail input validation and got processes by the application.
func send(w *http.ResponseWriter, r *msg.Result) {
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

	// copy result data into the output object
	switch r.Section {
	case msg.SectionAction:
		result = proto.NewActionResult()
		*result.Actions = append(*result.Actions, r.ActionObj...)

	case msg.SectionAttribute:
		result = proto.NewAttributeResult()
		*result.Attributes = append(*result.Attributes, r.Attribute...)

	case msg.SectionBucket:
		switch r.Action {
		case msg.ActionTree:
			result = proto.NewTreeResult()
			*result.Tree = r.Tree
		default:
			result = proto.NewBucketResult()
			*result.Buckets = append(*result.Buckets, r.Bucket...)
		}

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
		switch r.Action {
		case msg.ActionTree:
			result = proto.NewTreeResult()
			*result.Tree = r.Tree
		default:
			result = proto.NewClusterResult()
			*result.Clusters = append(*result.Clusters, r.Cluster...)
		}

	case msg.SectionDatacenter:
		result = proto.NewDatacenterResult()
		*result.Datacenters = append(*result.Datacenters, r.Datacenter...)

	case msg.SectionDeployment:
		result = proto.NewDeploymentResult()
		*result.Deployments = append(*result.Deployments, r.Deployment...)

	case msg.SectionEntity:
		result = proto.NewEntityResult()
		*result.Entities = append(*result.Entities, r.Entity...)

	case msg.SectionEnvironment:
		result = proto.NewEnvironmentResult()
		*result.Environments = append(*result.Environments, r.Environment...)

	case msg.SectionGroup:
		switch r.Action {
		case msg.ActionTree:
			result = proto.NewTreeResult()
			*result.Tree = r.Tree
		default:
			result = proto.NewGroupResult()
			*result.Groups = append(*result.Groups, r.Group...)
		}

	case msg.SectionHostDeployment:
		result = proto.NewHostDeploymentResult()
		*result.Deployments = append(*result.Deployments, r.Deployment...)
		*result.HostDeployments = append(*result.HostDeployments, r.HostDeployment...)

	case msg.SectionInstance:
		fallthrough
	case msg.SectionInstanceMgmt:
		result = proto.NewInstanceResult()
		*result.Instances = append(*result.Instances, r.Instance...)

	case msg.SectionJob:
		fallthrough
	case msg.SectionJobMgmt:
		result = proto.NewJobResult()
		*result.Jobs = append(*result.Jobs, r.Job...)

	case msg.SectionLevel:
		result = proto.NewLevelResult()
		*result.Levels = append(*result.Levels, r.Level...)

	case msg.SectionMetric:
		result = proto.NewMetricResult()
		*result.Metrics = append(*result.Metrics, r.Metric...)

	case msg.SectionMode:
		result = proto.NewModeResult()
		*result.Modes = append(*result.Modes, r.Mode...)

	case msg.SectionMonitoring:
		fallthrough
	case msg.SectionMonitoringMgmt:
		result = proto.NewMonitoringResult()
		*result.Monitorings = append(*result.Monitorings, r.Monitoring...)

	case msg.SectionNode:
		fallthrough
	case msg.SectionNodeConfig:
		fallthrough
	case msg.SectionNodeMgmt:
		switch r.Action {
		case msg.ActionTree:
			result = proto.NewTreeResult()
			*result.Tree = r.Tree
		default:
			result = proto.NewNodeResult()
			*result.Nodes = append(*result.Nodes, r.Node...)
		}

	case msg.SectionOncall:
		result = proto.NewOncallResult()
		*result.Oncalls = append(*result.Oncalls, r.Oncall...)

	case msg.SectionPermission:
		result = proto.NewPermissionResult()
		*result.Permissions = append(*result.Permissions, r.Permission...)

	case msg.SectionPredicate:
		result = proto.NewPredicateResult()
		*result.Predicates = append(*result.Predicates, r.Predicate...)

	case msg.SectionPropertyMgmt:
		fallthrough
	case msg.SectionPropertyCustom:
		fallthrough
	case msg.SectionPropertyNative:
		fallthrough
	case msg.SectionPropertyService:
		fallthrough
	case msg.SectionPropertySystem:
		fallthrough
	case msg.SectionPropertyTemplate:
		result = proto.NewPropertyResult()
		*result.Properties = append(*result.Properties, r.Property...)

	case msg.SectionProvider:
		result := proto.NewProviderResult()
		*result.Providers = append(*result.Providers, r.Provider...)

	case msg.SectionRepository:
		fallthrough
	case msg.SectionRepositoryConfig:
		fallthrough
	case msg.SectionRepositoryMgmt:
		switch r.Action {
		case msg.ActionTree:
			result = proto.NewTreeResult()
			*result.Tree = r.Tree
		default:
			result = proto.NewRepositoryResult()
			*result.Repositories = append(*result.Repositories, r.Repository...)
		}

	case msg.SectionRight:
		result = proto.NewGrantResult()
		*result.Grants = append(*result.Grants, r.Grant...)

	case msg.SectionSection:
		result = proto.NewSectionResult()
		*result.Sections = append(*result.Sections, r.SectionObj...)

	case msg.SectionServer:
		result = proto.NewServerResult()
		*result.Servers = append(*result.Servers, r.Server...)

	case msg.SectionState:
		result = proto.NewStateResult()
		*result.States = append(*result.States, r.State...)

	case msg.SectionStatus:
		result = proto.NewStatusResult()
		*result.Status = append(*result.Status, r.Status...)

	case msg.SectionSupervisor:
		switch r.Action {
		case msg.ActionToken: // switch r.Action
			switch r.Super.Task {
			case msg.TaskInvalidate: // switch r.Super.Task
				fallthrough
			case msg.TaskInvalidateAccount: // switch r.Super.Task
				result = proto.NewResult()
				if r.Code == 200 && r.Super.Verdict == 200 {
					log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 200)
					result.OK()
					goto buildJSON
				}
				dispatchForbidden(w, nil)
			case msg.TaskRequest: // switch r.Super.Task
				// token requests are encrypted
				result = proto.NewResult()
				if r.Code == 200 && r.Super.Verdict == 200 {
					log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 200)
					goto dispatchOCTET
				}
			default: // switch r.Super.Task
				log.Printf(msg.LogStrErr, r.Section, fmt.Sprintf("%s/%s", r.Action, r.Super.Task), 0,
					`Result for unhandled supervisor task`)
				dispatchForbidden(w, nil)
			}
			return
		case msg.ActionKex: // switch r.Action
			k = r.Super.Kex
			if bjson, err = json.Marshal(&k); err != nil {
				log.Printf(msg.LogStrErr, r.Section, r.Action, r.Code, err.Error())
				dispatchInternalError(w, nil)
				return
			}
			goto dispatchJSON
		case msg.ActionPassword: // switch r.Action
			switch r.Super.Task {
			case msg.TaskChange: // switch r.Super.Task
			case msg.TaskReset: // switch r.Super.Task
			default: // switch r.Super.Task
				log.Printf(msg.LogStrErr, r.Section, fmt.Sprintf("%s/%s", r.Action, r.Super.Task), 0,
					`Result for unhandled supervisor task`)
				dispatchForbidden(w, nil)
				return
			}
			switch r.Code {
			case 200: // switch r.Code
				if r.Super.Verdict == 200 {
					log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 200)
					goto dispatchOCTET
				}
				log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 403)
				dispatchForbidden(w, nil)
			default: // switch r.Code
				log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 403)
				dispatchForbidden(w, nil)
			}
			return
		case msg.ActionActivate: // switch r.Action
			switch r.Super.Task {
			case msg.SubjectRoot: // switch r.Super.Task
			case msg.SubjectUser: // switch r.Super.Task
			default: // switch r.Super.Task
				log.Printf(msg.LogStrErr, r.Section, fmt.Sprintf("%s/%s", r.Action, r.Super.Task), 0,
					`Result for unhandled supervisor task subject`)
				dispatchForbidden(w, nil)
				return
			}
			switch r.Code {
			case 200: // switch r.Code
				if r.Super.Verdict == 200 {
					log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 200)
					goto dispatchOCTET
				}
				log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 403)
				dispatchForbidden(w, nil)
			case 406: // switch r.Code
				log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 406)
				// request failed due to a policy constraint, do not
				// mask the error and return the full detail error message
				dispatchConflict(w, r.Error)
			default: // switch r.Code
				log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 403)
				dispatchForbidden(w, nil)
			}
			return
		default: // switch r.Action
			log.Printf(msg.LogStrErr, r.Section, r.Action, 0, `Result for unhandled supervisor action`)
			dispatchForbidden(w, nil)
			return
		}

	case msg.SectionSystem:
		result = proto.NewResult()
		switch r.Action {
		case msg.ActionRepoRebuild:
		case msg.ActionRepoRestart:
		case msg.ActionRepoStop:
		case msg.ActionShutdown:
		case msg.ActionToken:
			// Supervisor interactions are masked
			switch r.Super.Task {
			case msg.TaskInvalidateAccount:
			case msg.TaskInvalidateGlobal:
			default:
				log.Printf(msg.LogStrErr, r.Section, fmt.Sprintf("%s/%s", r.Action, r.Super.Task), 0,
					`Result for unhandled supervisor task`)
				dispatchForbidden(w, nil)
				return
			}
			if r.Code == 200 && r.Super.Verdict == 200 {
				log.Printf(msg.LogStrOK, r.Section, r.Action, r.Code, 200)
				result.OK()
				goto buildJSON
			}
			dispatchForbidden(w, nil)
			return
		default:
			log.Printf(msg.LogStrErr, r.Section, r.Action, 0, `Result for unhandled system action`)
			dispatchInternalError(w, nil)
			return
		}

	case msg.SectionTeam:
		fallthrough
	case msg.SectionTeamMgmt:
		result = proto.NewTeamResult()
		*result.Teams = append(*result.Teams, r.Team...)

	case msg.SectionUnit:
		result = proto.NewUnitResult()
		*result.Units = append(*result.Units, r.Unit...)

	case msg.SectionUser:
		fallthrough
	case msg.SectionUserMgmt:
		result = proto.NewUserResult()
		*result.Users = append(*result.Users, r.User...)

	case msg.SectionValidity:
		result = proto.NewValidityResult()
		*result.Validities = append(*result.Validities, r.Validity...)

	case msg.SectionView:
		result = proto.NewViewResult()
		*result.Views = append(*result.Views, r.View...)

	case msg.SectionWorkflow:
		result = proto.NewWorkflowResult()
		*result.Workflows = append(*result.Workflows, r.Workflow...)

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
