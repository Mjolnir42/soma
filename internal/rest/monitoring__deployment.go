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

	if err := checkStringIsUUID(params.ByName(`deploymentID`)); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionDeployment
	request.Action = msg.ActionShow
	request.Deployment.ID = params.ByName(`deploymentID`)

	// BUG	if !x.isAuthorized(&request) {
	// BUG		x.replyForbidden(&w, &request, nil)
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

	if err := checkStringIsUUID(params.ByName(`deploymentID`)); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionDeployment
	request.Deployment.ID = params.ByName(`deploymentID`)

	switch params.ByName(`action`) {
	case msg.ActionSuccess:
		request.Action = msg.ActionSuccess
	case msg.ActionFailed:
		request.Action = msg.ActionFailed
	default:
		dispatchBadRequest(&w, fmt.Errorf("Unknown action: %s", params.ByName(`action`)))
		return
	}

	if params.ByName(`monitoringID`) != `` {
		if err := checkStringIsUUID(params.ByName(`monitoringID`)); err != nil {
			dispatchBadRequest(&w, err)
			return
		}
		request.Monitoring.ID = params.ByName(`monitoringID`)
	}

	// BUG	if !x.isAuthorized(&request) {
	// BUG		x.replyForbidden(&w, &request, nil)
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

	if err := checkStringIsUUID(params.ByName(`monitoringID`)); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionDeployment
	request.Action = msg.ActionList
	request.Monitoring.ID = params.ByName(`monitoringID`)

	// BUG	if !x.isAuthorized(&request) {
	// BUG		x.replyForbidden(&w, &request, nil)
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

	if err := checkStringIsUUID(params.ByName(`monitoringID`)); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionDeployment
	request.Action = msg.ActionPending
	request.Monitoring.ID = params.ByName(`monitoringID`)

	// BUG	if !x.isAuthorized(&request) {
	// BUG		x.replyForbidden(&w, &request, nil)
	// BUG		return
	// BUG	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// DeploymentFilter function
func (x *Rest) DeploymentFilter(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	dispatchNotImplemented(&w, nil)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
