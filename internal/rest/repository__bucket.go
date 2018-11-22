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
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/proto"
)

// BucketList function
func (x *Rest) BucketList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionBucket
	request.Action = msg.ActionList

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// BucketSearch function
func (x *Rest) BucketSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionBucket
	request.Action = msg.ActionList

	cReq := proto.NewBucketFilter()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	if cReq.Filter.Bucket.Name == `` && cReq.Filter.Bucket.ID == `` {
		x.replyBadRequest(&w, &request,
			fmt.Errorf(`BucketSearch request without condition`))
		return
	}
	request.Search.Bucket.ID = cReq.Filter.Bucket.ID
	request.Search.Bucket.Name = cReq.Filter.Bucket.Name

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply

	// XXX BUG filter in SQL statement
	filtered := []proto.Bucket{}
	for _, i := range result.Bucket {
		if (i.Name == cReq.Filter.Bucket.Name) || (i.ID == cReq.Filter.Bucket.ID) {
			filtered = append(filtered, i)
		}
	}
	result.Bucket = filtered
	x.send(&w, &result)
}

// BucketShow function
func (x *Rest) BucketShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionBucket
	request.Action = msg.ActionShow
	request.Bucket = proto.Bucket{
		ID: params.ByName(`bucket`),
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// BucketTree function
func (x *Rest) BucketTree(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionBucket
	request.Action = msg.ActionTree
	request.Tree = proto.Tree{
		ID:   params.ByName(`bucketID`),
		Type: msg.EntityBucket,
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// BucketCreate function
func (x *Rest) BucketCreate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionBucket
	request.Action = msg.ActionCreate

	cReq := proto.NewBucketRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	if cReq.Bucket.Name == `` || cReq.Bucket.Environment == `` ||
		cReq.Bucket.TeamID == `` || cReq.Bucket.RepositoryID == `` {
		x.replyBadRequest(&w, &request,
			fmt.Errorf(`Incomplete Bucket.Create request`))
		return
	}

	nameLen := utf8.RuneCountInString(cReq.Bucket.Name)
	if nameLen < 4 || nameLen > 512 {
		x.replyBadRequest(&w, &request,
			fmt.Errorf(`Illegal bucket name length (4 < x <= 512)`))
		return
	}
	request.Bucket = cReq.Bucket.Clone()

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// BucketDestroy function
func (x *Rest) BucketDestroy(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	// TODO
}

// BucketRename function
func (x *Rest) BucketRename(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionBucket
	request.Action = msg.ActionRename

	cReq := proto.NewBucketRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Update.Bucket.Name = cReq.Bucket.Name

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// BucketMemberList function
func (x *Rest) BucketMemberList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	// TODO
}

// BucketMemberAssign function
func (x *Rest) BucketMemberAssign(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	// TODO
}

// BucketMemberUnassign function
func (x *Rest) BucketMemberUnassign(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	// TODO
}

// BucketPropertyCreate function
func (x *Rest) BucketPropertyCreate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionBucket
	request.Action = msg.ActionPropertyCreate

	cReq := proto.NewBucketRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	switch {
	case params.ByName(`bucketID`) != cReq.Bucket.ID:
		x.replyBadRequest(&w, &request,
			fmt.Errorf("Mismatched bucket ids: %s, %s",
				params.ByName(`bucket`),
				cReq.Bucket.ID))
		return
	case len(*cReq.Bucket.Properties) != 1:
		x.replyBadRequest(&w, &request,
			fmt.Errorf("Expected property count 1, actual count: %d",
				len(*cReq.Bucket.Properties)))
		return
	case (*cReq.Bucket.Properties)[0].Type == `service` && (*cReq.Bucket.Properties)[0].Service.Name == ``:
		x.replyBadRequest(&w, &request,
			fmt.Errorf(`Empty service name is invalid`))
		return
	}
	request.TargetEntity = msg.EntityBucket
	request.Bucket = cReq.Bucket.Clone()
	request.Property.Type = params.ByName(`propertyType`)

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// BucketPropertyDestroy function
func (x *Rest) BucketPropertyDestroy(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionBucket
	request.Action = msg.ActionPropertyDestroy
	request.TargetEntity = msg.EntityBucket
	request.Property.Type = params.ByName(`propertyType`)
	request.Bucket = proto.Bucket{
		ID: params.ByName(`bucket`),
		Properties: &[]proto.Property{
			proto.Property{
				Type:             params.ByName(`propertyType`),
				BucketID:         params.ByName(`bucketID`),
				SourceInstanceID: params.ByName(`sourceID`),
			},
		},
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// BucketPropertyUpdate function
func (x *Rest) BucketPropertyUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	// XXX BUG TODO
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
