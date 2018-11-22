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

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/proto"
)

// WorkflowSummary returns information about the current workflow
// status distribution
func (x *Rest) WorkflowSummary(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionWorkflow
	request.Action = msg.ActionSummary

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// WorkflowList function
func (x *Rest) WorkflowList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionWorkflow
	request.Action = msg.ActionList

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// WorkflowSearch function
func (x *Rest) WorkflowSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewWorkflowFilter()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}
	if cReq.Filter.Workflow.Status == `` {
		dispatchBadRequest(&w, fmt.Errorf(
			`No workflow status specified`))
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionWorkflow
	request.Action = msg.ActionSearch
	request.Workflow = proto.Workflow{
		Status: cReq.Filter.Workflow.Status,
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// WorkflowRetry function
func (x *Rest) WorkflowRetry(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewWorkflowRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}
	if cReq.Workflow.InstanceID == `` {
		dispatchBadRequest(&w, fmt.Errorf(
			`No instanceID for retry specified`))
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionWorkflow
	request.Action = msg.ActionRetry
	request.Workflow = proto.Workflow{
		InstanceID: cReq.Workflow.InstanceID,
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// WorkflowSet function
func (x *Rest) WorkflowSet(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewWorkflowRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	// set requests must be fully specified
	switch {
	case cReq.Workflow.Status == ``:
		fallthrough
	case cReq.Workflow.NextStatus == ``:
		fallthrough
	case params.ByName(`instanceconfigID`) == ``:
		dispatchBadRequest(&w, fmt.Errorf(
			`Incomplete status information specified`))
		return
	default:
	}

	// It's dangerous out there, take this -f
	if !cReq.Flags.Forced {
		dispatchBadRequest(&w, fmt.Errorf(
			`WorkflowSet request declined, force required.`))
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionWorkflow
	request.Action = msg.ActionSet
	request.Workflow = proto.Workflow{
		InstanceConfigID: params.ByName(`instanceconfigID`),
		Status:           cReq.Workflow.Status,
		NextStatus:       cReq.Workflow.NextStatus,
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
