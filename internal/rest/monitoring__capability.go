/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
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
	"github.com/mjolnir42/soma/lib/proto"
)

// CapabilityList function
func (x *Rest) CapabilityList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionCapability
	request.Action = msg.ActionList

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// CapabilitySearch function
func (x *Rest) CapabilitySearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewCapabilityFilter()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	if cReq.Filter.Capability.MonitoringID == `` {
		dispatchBadRequest(&w,
			fmt.Errorf(`CapabilitySearch request missing MonitoringID`))
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionCapability
	request.Action = msg.ActionSearch

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply

	// XXX BUG filter in SQL statement
	filtered := []proto.Capability{}
	for _, i := range result.Capability {
		if i.MonitoringID == cReq.Filter.Capability.MonitoringID &&
			i.Metric == cReq.Filter.Capability.Metric &&
			i.View == cReq.Filter.Capability.View {
			filtered = append(filtered, i)
		}
	}
	result.Capability = filtered

	// XXX BUG do not return these fields for search
	// cleanup reply, only keep ID and Name
	for i := range result.Capability {
		result.Capability[i].MonitoringID = ``
		result.Capability[i].Metric = ``
		result.Capability[i].View = ``
	}
	x.send(&w, &result)
}

// CapabilityShow function
func (x *Rest) CapabilityShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionCapability
	request.Action = msg.ActionShow
	request.Capability = proto.Capability{
		ID: params.ByName(`capabilityID`),
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// CapabilityDeclare function
func (x *Rest) CapabilityDeclare(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewCapabilityRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionCapability
	request.Action = msg.ActionDeclare
	request.Capability = cReq.Capability.Clone()

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// CapabilityRevoke function
func (x *Rest) CapabilityRevoke(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionCapability
	request.Action = msg.ActionRevoke
	request.Capability = proto.Capability{
		ID: params.ByName(`capabilityID`),
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
