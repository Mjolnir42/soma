package rest // import "github.com/mjolnir42/soma/internal/rest"

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/proto"
)

// SystemOperation function
func (x *Rest) SystemOperation(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewSystemRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	switch cReq.System.Request {
	case msg.ActionRepoRebuild:
	case msg.ActionRepoRestart:
	case msg.ActionRepoStop:
	case msg.ActionShutdown:
	default:
		dispatchBadRequest(&w, fmt.Errorf(
			"Mismatching requests: %s vs %s",
			cReq.System.Request,
			msg.ActionRepoStop,
		))
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionSystem
	request.Action = cReq.System.Request
	request.System = proto.System{
		Request:      cReq.System.Request,
		RepositoryID: cReq.System.RepositoryID,
		RebuildLevel: cReq.System.RebuildLevel,
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).PriorityIntake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// SupervisorTokenInvalidateGlobal is the rest endpoint for admins
// to invalidate all current access tokens
func (x *Rest) SupervisorTokenInvalidateGlobal(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionSystem
	request.Action = msg.ActionToken
	request.Super = &msg.Supervisor{
		Task: msg.TaskInvalidateGlobal,
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// SupervisorTokenInvalidateAccount is the rest endpoint for admins
// to invalidate all current access tokens for another user
func (x *Rest) SupervisorTokenInvalidateAccount(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionSystem
	request.Action = msg.ActionToken
	request.Super = &msg.Supervisor{
		Task:            msg.TaskInvalidateAccount,
		RevokeForName: params.ByName(`account`),
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
