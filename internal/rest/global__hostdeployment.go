/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/soma/internal/rest"

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/proto"
)

// HostDeploymentFetch function
func (x *Rest) HostDeploymentFetch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)
	var (
		err     error
		assetID uint64
	)

	request := msg.New(r, params)
	request.Section = msg.SectionHostDeployment
	request.Action = msg.ActionGet

	if err = checkStringIsUUID(params.ByName(`monitoringID`)); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	if assetID, err = strconv.ParseUint(params.ByName(`assetID`), 10, 64); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	request.Monitoring.ID = params.ByName(`monitoringID`)
	request.Node.AssetID = assetID

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// HostDeploymentAssemble function
func (x *Rest) HostDeploymentAssemble(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)
	var (
		err     error
		assetID uint64
	)

	request := msg.New(r, params)
	request.Section = msg.SectionHostDeployment
	request.Action = msg.ActionAssemble

	if err = checkStringIsUUID(params.ByName(`monitoringID`)); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}
	if assetID, err = strconv.ParseUint(params.ByName(`assetID`), 10, 64); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	cReq := proto.NewHostDeploymentRequest()
	if err = decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}
	if cReq.HostDeployment == nil {
		x.replyBadRequest(&w, &request, fmt.Errorf(`HostDeployment section missing`))
		return
	}

	request.Monitoring.ID = params.ByName(`monitoringID`)
	request.Node.AssetID = assetID
	request.DeploymentIDs = cReq.HostDeployment.CurrentCheckInstanceIDList

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
