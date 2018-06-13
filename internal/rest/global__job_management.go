/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/soma/internal/rest"

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/soma/internal/msg"
)

// JobMgmtWait function
func (x *Rest) JobMgmtWait(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionJobMgmt
	request.Action = msg.ActionWait
	request.Job.ID = params.ByName(`jobID`)
	request.Flag.Unscoped = true

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`job_block`)
	handler.Intake() <- request
	<-request.Reply
	dispatchNoContent(&w)
}

// JobMgmtList function
func (x *Rest) JobMgmtList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionJobMgmt
	request.Action = msg.ActionList
	request.Flag.Unscoped = true

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`job_r`)
	handler.Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)

}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
