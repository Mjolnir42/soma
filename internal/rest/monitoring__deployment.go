/*-
 * Copyright (c) 2016-2018, Jörg Pernfuß
 * Copyright (c) 2018, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/soma/internal/rest"

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/soma/internal/msg"
)

// DeploymentShow function
func (x *Rest) DeploymentShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionDeployment
	request.Action = msg.ActionShow

	if err := checkStringIsUUID(params.ByName(`deploymentID`)); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}
	request.Deployment.ID = params.ByName(`deploymentID`)

	// BUG	if !x.isAuthorized(&request) {
	// BUG		x.replyForbidden(&w, &request)
	// BUG		return
	// BUG	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// DeploymentUpdate function
func (x *Rest) DeploymentUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionDeployment
	switch params.ByName(`action`) {
	case msg.ActionSuccess:
		request.Action = msg.ActionSuccess
	case msg.ActionFailed:
		request.Action = msg.ActionFailed
	default:
		x.replyBadRequest(&w, &request, fmt.Errorf("Unknown action: %s", params.ByName(`action`)))
		return
	}

	if err := checkStringIsUUID(params.ByName(`deploymentID`)); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}
	request.Deployment.ID = params.ByName(`deploymentID`)

	if params.ByName(`monitoringID`) != `` {
		if err := checkStringIsUUID(params.ByName(`monitoringID`)); err != nil {
			x.replyBadRequest(&w, &request, err)
			return
		}
		request.Monitoring.ID = params.ByName(`monitoringID`)
	}

	// BUG	if !x.isAuthorized(&request) {
	// BUG		x.replyForbidden(&w, &request)
	// BUG		return
	// BUG	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// DeploymentList function
func (x *Rest) DeploymentList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionDeployment
	request.Action = msg.ActionList

	if err := checkStringIsUUID(params.ByName(`monitoringID`)); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}
	request.Monitoring.ID = params.ByName(`monitoringID`)

	// BUG	if !x.isAuthorized(&request) {
	// BUG		x.replyForbidden(&w, &request)
	// BUG		return
	// BUG	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// DeploymentPending function
func (x *Rest) DeploymentPending(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionDeployment
	request.Action = msg.ActionPending

	if err := checkStringIsUUID(params.ByName(`monitoringID`)); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}
	request.Monitoring.ID = params.ByName(`monitoringID`)

	// BUG	if !x.isAuthorized(&request) {
	// BUG		x.replyForbidden(&w, &request)
	// BUG		return
	// BUG	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// DeploymentFilter function
func (x *Rest) DeploymentFilter(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionDeployment
	request.Action = msg.ActionFilter

	x.replyNotImplemented(&w, &request, nil)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
