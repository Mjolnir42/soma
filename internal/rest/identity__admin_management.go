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
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/proto"
)

// AdminMgmtAdd function
func (x *Rest) AdminMgmtAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionAdminMgmt
	request.Action = msg.ActionAdd

	cReq := proto.NewAdminRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}
	if strings.Contains(cReq.Admin.UserName, `:`) {
		x.replyBadRequest(&w, &request, fmt.Errorf(
			`Invalid username containing : character`))
		return
	}
	request.Admin.UserName = cReq.Admin.UserName

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// AdminMgmtRemove function
func (x *Rest) AdminMgmtRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionAdminMgmt
	request.Action = msg.ActionRemove

	request.Admin.ID = params.ByName(`adminID`)
	request.Admin.UserID = params.ByName(`userID`)

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	switch request.Action {
	case msg.ActionRemove:
		// expire admin's credentials
		credInval := msg.New(r, params)
		credInval.Section = msg.SectionSupervisor
		credInval.Action = msg.ActionPassword
		credInval.Super = &msg.Supervisor{
			Task:        msg.TaskRevoke,
			RevokeForID: params.ByName(`adminID`),
		}
		x.handlerMap.MustLookup(&request).Intake() <- credInval
		<-credInval.Reply

		// expire admin's tokens
		tokenInval := msg.New(r, params)
		tokenInval.Section = msg.SectionSupervisor
		tokenInval.Action = msg.ActionToken
		tokenInval.Super = &msg.Supervisor{
			Task:        msg.TaskInvalidateAccount,
			RevokeForID: params.ByName(`adminID`),
		}
		x.handlerMap.MustLookup(&request).Intake() <- tokenInval
		<-tokenInval.Reply
	default:
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

func (x *Rest) AdminMgmtShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionAdminMgmt
	request.Action = msg.ActionShow

	request.Admin.UserID = params.ByName(`userID`)

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
