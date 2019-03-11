/*-
 * Copyright (c) 2016-2018, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2018, 1&1 IONOS SE
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/soma/internal/rest"

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/auth"
	"github.com/mjolnir42/soma/lib/proto"
)

// send is the output function for all requests that did not
// fail input validation and got processes by the application.
func (x *Rest) send(w *http.ResponseWriter, r *msg.Result) {
	var (
		bjson  []byte
		err    error
		k      auth.Kex
		result proto.Result
	)

	// build RequestLog entry
	logEntry := x.reqLog.WithField(`RequestID`, r.ID.String()).
		WithField(`Section`, r.Section).
		WithField(`Action`, r.Action).
		WithField(`Request`, fmt.Sprintf("%s::%s", r.Section, r.Action)).
		WithField(`RequestURI`, r.RequestURI).
		WithField(`Phase`, `result`)

	// this is central error command, proceeding to ErrorLog while
	// updating the RequestLog metadata
	if r.Error != nil {
		x.errLog.WithField(`RequestID`, r.ID.String()).
			WithField(`Section`, r.Section).
			WithField(`Action`, r.Action).
			WithField(`Phase`, `result`).
			WithField(`Code`, r.Code).
			Errorln(r.Error.Error())
		logEntry = logEntry.WithField(`HasError`, `true`)
	} else {
		logEntry = logEntry.WithField(`HasError`, `false`)
	}

	// copy result data into the output object
	switch r.Section {
	// simple result data with straight forward copying

	case msg.SectionAction:
		result = proto.NewActionResult()
		*result.Actions = append(*result.Actions, r.ActionObj...)
	case msg.SectionAttribute:
		result = proto.NewAttributeResult()
		*result.Attributes = append(*result.Attributes, r.Attribute...)
	case msg.SectionCapability:
		result = proto.NewCapabilityResult()
		*result.Capabilities = append(*result.Capabilities, r.Capability...)
	case msg.SectionCategory:
		result = proto.NewCategoryResult()
		*result.Categories = append(*result.Categories, r.Category...)
	case msg.SectionCheckConfig:
		result = proto.NewCheckConfigResult()
		*result.CheckConfigs = append(*result.CheckConfigs, r.CheckConfig...)
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
	case msg.SectionHostDeployment:
		result = proto.NewHostDeploymentResult()
		*result.Deployments = append(*result.Deployments, r.Deployment...)
		*result.HostDeployments = append(*result.HostDeployments, r.HostDeployment...)
	case msg.SectionJobResultMgmt:
		result = proto.NewJobResultResult()
		*result.JobResults = append(*result.JobResults, r.JobResult...)
	case msg.SectionJobStatusMgmt:
		result = proto.NewJobStatusResult()
		*result.JobStatus = append(*result.JobStatus, r.JobStatus...)
	case msg.SectionJobTypeMgmt:
		result = proto.NewJobTypeResult()
		*result.JobTypes = append(*result.JobTypes, r.JobType...)
	case msg.SectionLevel:
		result = proto.NewLevelResult()
		*result.Levels = append(*result.Levels, r.Level...)
	case msg.SectionMetric:
		result = proto.NewMetricResult()
		*result.Metrics = append(*result.Metrics, r.Metric...)
	case msg.SectionMode:
		result = proto.NewModeResult()
		*result.Modes = append(*result.Modes, r.Mode...)
	case msg.SectionOncall:
		result = proto.NewOncallResult()
		*result.Oncalls = append(*result.Oncalls, r.Oncall...)
	case msg.SectionPermission:
		result = proto.NewPermissionResult()
		*result.Permissions = append(*result.Permissions, r.Permission...)
	case msg.SectionPredicate:
		result = proto.NewPredicateResult()
		*result.Predicates = append(*result.Predicates, r.Predicate...)
	case msg.SectionProvider:
		result = proto.NewProviderResult()
		*result.Providers = append(*result.Providers, r.Provider...)
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
	case msg.SectionUnit:
		result = proto.NewUnitResult()
		*result.Units = append(*result.Units, r.Unit...)
	case msg.SectionValidity:
		result = proto.NewValidityResult()
		*result.Validities = append(*result.Validities, r.Validity...)
	case msg.SectionView:
		result = proto.NewViewResult()
		*result.Views = append(*result.Views, r.View...)
	case msg.SectionWorkflow:
		result = proto.NewWorkflowResult()
		*result.Workflows = append(*result.Workflows, r.Workflow...)

	// result data with multiple permission scopes combines different
	// sections

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
	case msg.SectionMonitoring:
		fallthrough
	case msg.SectionMonitoringMgmt:
		result = proto.NewMonitoringResult()
		*result.Monitorings = append(*result.Monitorings, r.Monitoring...)
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
	case msg.SectionTeam:
		fallthrough
	case msg.SectionTeamMgmt:
		switch r.Action {
		case msg.ActionMemberList:
			result = proto.NewUserResult()
			*result.Users = append(*result.Users, r.User...)
		default:
			result = proto.NewTeamResult()
			*result.Teams = append(*result.Teams, r.Team...)
		}
	case msg.SectionUser:
		fallthrough
	case msg.SectionUserMgmt:
		result = proto.NewUserResult()
		*result.Users = append(*result.Users, r.User...)
	case msg.SectionAdminMgmt:
		result = proto.NewAdminResult()
		*result.Admins = append(*result.Admins, r.Admin...)

	// tree configuration results have different result data based on
	// the action and may have multiple scopes

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
	case msg.SectionBucket:
		switch r.Action {
		case msg.ActionTree:
			result = proto.NewTreeResult()
			*result.Tree = r.Tree
		default:
			result = proto.NewBucketResult()
			*result.Buckets = append(*result.Buckets, r.Bucket...)
		}
	case msg.SectionGroup:
		switch r.Action {
		case msg.ActionTree:
			result = proto.NewTreeResult()
			*result.Tree = r.Tree
		default:
			result = proto.NewGroupResult()
			*result.Groups = append(*result.Groups, r.Group...)
		}
	case msg.SectionCluster:
		switch r.Action {
		case msg.ActionTree:
			result = proto.NewTreeResult()
			*result.Tree = r.Tree
		default:
			result = proto.NewClusterResult()
			*result.Clusters = append(*result.Clusters, r.Cluster...)
		}
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

	// system results have generally no output, or a supervisor actions
	// that require special handling

	case msg.SectionSystem:
		result = proto.NewResult()
		result.RequestID = r.ID.String()
		switch r.Action {
		case msg.ActionRepoRebuild:
		case msg.ActionRepoRestart:
		case msg.ActionRepoStop:
		case msg.ActionShutdown:
		case msg.ActionToken:
			// system::token is a supervisor action

			// check supervisor task
			switch r.Super.Task {
			case msg.TaskInvalidateAccount:
			case msg.TaskInvalidateGlobal:
			default:
				logEntry.WithField(`Code`, r.Code).
					WithField(`Masked`, 403).
					WithField(`Task`, r.Super.Task).
					Warnf(`Unhandled supervisor task`)
				result.Forbidden(nil)
				goto buildJSON
			}

			// check supervisor verdict
			if r.Code == 200 && r.Super.Verdict == 200 {
				result.OK()
				logEntry.WithField(`Code`, r.Code).Info(`OK`)
				goto buildJSON
			}

			// mask as 403/Forbidden
			logEntry.WithField(`Code`, r.Code).
				WithField(`Masked`, 403).
				Warnf(`Forbidden`)
			result.Forbidden(nil)
			goto buildJSON

		default:
			logEntry.WithField(`Code`, r.Code).
				WithField(`Masked`, 500).
				Warnf(`Unhandled system action`)
			result.Forbidden(nil)
			goto buildJSON
		}

	// supervisor results handle AAA data and mask internal error codes
	// to avoid information leaks

	case msg.SectionSupervisor:
		switch r.Action {
		case msg.ActionToken:
			result = proto.NewResult()
			result.RequestID = r.ID.String()

			switch r.Super.Task {
			// invalidate Token requests
			case msg.TaskInvalidate:
				fallthrough
			case msg.TaskInvalidateAccount:
				// check supervisor verdict
				if r.Code == 200 && r.Super.Verdict == 200 {
					result.OK()
					logEntry.WithField(`Code`, r.Code).Info(`OK`)
					goto buildJSON
				}

				// mask as 403/Forbidden
				logEntry.WithField(`Code`, r.Code).
					WithField(`Masked`, 403).
					WithField(`Task`, r.Super.Task).
					Warnf(`Forbidden`)
				result.Forbidden(nil)
				goto buildJSON

			// token generation request - encrypted payload
			case msg.TaskRequest:
				// check supervisor verdict
				if r.Code == 200 && r.Super.Verdict == 200 {
					logEntry.WithField(`Code`, r.Code).Info(`OK`)
					goto dispatchOCTET
				}

				// mask as 403/Forbidden
				logEntry.WithField(`Code`, r.Code).
					WithField(`Masked`, 403).
					WithField(`Task`, r.Super.Task).
					Warnf(`Forbidden`)
				result.Forbidden(nil)
				goto buildJSON

			default: // switch r.Super.Task
				logEntry.WithField(`Code`, r.Code).
					WithField(`Masked`, 403).
					WithField(`Task`, r.Super.Task).
					Warnf(`Unhandled supervisor task`)
				result.Forbidden(nil)
				goto buildJSON

			}
			// unreachable

		case msg.ActionKex:
			// Key Exchange request -- encrypted payload
			k = r.Super.Kex
			// special KEX payload is generated here
			if bjson, err = json.Marshal(&k); err != nil {
				logEntry.WithField(`Code`, 500).
					WithField(`HasError`, `true`).
					Warn(`ServerError`)
				x.errLog.WithField(`RequestID`, r.ID.String()).
					WithField(`Phase`, `json`).
					Error(err)
				// KEX has no regular application payload
				x.hardServerError(w)
				return
			}
			logEntry.WithField(`Code`, r.Code).Info(`OK`)
			goto dispatchJSON

		case msg.ActionPassword:
			// Password manipulation request -- encrypted payload

			// check upervisor task
			switch r.Super.Task {
			case msg.TaskChange:
			case msg.TaskReset:
			default:
				logEntry.WithField(`Code`, r.Code).
					WithField(`Masked`, 403).
					WithField(`Task`, r.Super.Task).
					Warnf(`Unhandled supervisor task`)
				result.Forbidden(nil)
				goto buildJSON
			}

			// check supervisor verdict
			if r.Code == 200 && r.Super.Verdict == 200 {
				result.OK()
				logEntry.WithField(`Code`, r.Code).Info(`OK`)
				goto dispatchOCTET
			}

			// mask as 403/Forbidden
			logEntry.WithField(`Code`, r.Code).
				WithField(`Masked`, 403).
				Warnf(`Forbidden`)
			result.Forbidden(nil)
			goto buildJSON

		case msg.ActionActivate:
			// Account activation request -- encrypted payload

			// check upervisor task
			switch r.Super.Task {
			case msg.SubjectRoot:
			case msg.SubjectUser:
			default:
				logEntry.WithField(`Code`, r.Code).
					WithField(`Masked`, 403).
					WithField(`Task`, r.Super.Task).
					Warnf(`Unhandled supervisor task`)
				result.Forbidden(nil)
				goto buildJSON
			}

			// check supervisor verdict
			if r.Code == 200 && r.Super.Verdict == 200 {
				result.OK()
				logEntry.WithField(`Code`, r.Code).Info(`OK`)
				goto dispatchOCTET
			}

			// check policy violation
			if r.Code == 406 {
				// request failed due to a policy constraint, do not
				// mask the error and return the full detail error message
				logEntry.WithField(`Code`, r.Code).Warn(r.Error)
				x.hardConflict(w, r.Error)
				return
			}

			// mask as 403/Forbidden
			logEntry.WithField(`Code`, r.Code).
				WithField(`Masked`, 403).
				Warnf(`Forbidden`)
			result.Forbidden(nil)
			goto buildJSON

		default:
			logEntry.WithField(`Code`, r.Code).
				WithField(`Masked`, 403).
				Warnf(`Unhandled supervisor action`)
			result.Forbidden(nil)
			goto buildJSON
		}

	default:
		logEntry.WithField(`Code`, r.Code).
			WithField(`Masked`, 500).
			Warnf(`Result for unhandled section`)
		result.Error(nil)
		// prevent any data leak
		result.DataClean()
		goto buildJSON
	}
	result.RequestID = r.ID.String()

	logEntry = logEntry.WithField(`Code`, r.Code)

	switch r.Code {
	case 200:
		result.OK()
		logEntry.WithField(`Code`, r.Code).Info(`OK`)
		if r.Error != nil {
			result.Errors = &[]string{r.Error.Error()}
		}
	case 202:
		result.JobID = r.JobID
		result.Accepted()
		logEntry.WithField(`Code`, r.Code).
			WithField(`JobID`, r.JobID).Info(`Accepted`)
	case 400:
		result.BadRequest(r.Error)
		logEntry.WithField(`Code`, r.Code).Warn(`BadRequest`)
	case 403:
		result.Forbidden(r.Error)
		logEntry.WithField(`Code`, r.Code).Warn(`Forbidden`)
	case 404:
		result.NotFoundErr(r.Error)
		logEntry.WithField(`Code`, r.Code).Warn(`NotFound`)
	case 500:
		result.Error(r.Error)
		logEntry.WithField(`Code`, r.Code).Warn(`ServerError`)
	case 501:
		result.NotImplemented()
		logEntry.WithField(`Code`, r.Code).Warn(`NotImplemented`)
	case 503:
		result.Unavailable()
		logEntry.WithField(`Code`, r.Code).Warn(`ServiceUnavailable`)
	default:
		logEntry.WithField(`Code`, r.Code).
			WithField(`Masked`, 500).
			Warn(`Unhandled internal result code`)
		result.Error(nil)
		// prevent any data leak
		result.DataClean()
	}
	goto buildJSON

dispatchOCTET:
	x.writeReplyOctetStream(w, &r.Super.Encrypted.Data)
	return

buildJSON:
	if bjson, err = json.Marshal(&result); err != nil {
		x.errLog.WithField(`RequestID`, r.ID.String()).
			WithField(`Phase`, `json`).
			Error(err)
		x.hardServerError(w)
		return
	}

dispatchJSON:
	x.writeReplyJSON(w, &bjson)
	return
}

// writeReplyOctetStream writes out b as the reply with content-type set
// to application/octet-stream
func (x *Rest) writeReplyOctetStream(w *http.ResponseWriter, b *[]byte) {
	(*w).Header().Set(`Content-Type`, `application/octet-stream`)
	(*w).WriteHeader(http.StatusOK)
	(*w).Write(*b)
}

// writeReplyJSON writes out b as the reply with content-type
// set to application/json
func (x *Rest) writeReplyJSON(w *http.ResponseWriter, b *[]byte) {
	(*w).Header().Set(`Content-Type`, `application/json`)
	(*w).WriteHeader(http.StatusOK)
	(*w).Write(*b)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
