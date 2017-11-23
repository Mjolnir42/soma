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

	request := newRequest(r, params)
	request.Section = msg.SectionBucket
	request.Action = msg.ActionList

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`bucket_r`)
	handler.Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// BucketSearch function
func (x *Rest) BucketSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewBucketRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	if cReq.Filter.Bucket.Name == `` && cReq.Filter.Bucket.ID == `` {
		dispatchBadRequest(&w,
			fmt.Errorf(`BucketSearch request without condition`))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionBucket
	request.Action = msg.ActionList

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`bucket_r`)
	handler.Intake() <- request
	result := <-request.Reply

	// XXX BUG filter in SQL statement
	filtered := []proto.Bucket{}
	for _, i := range result.Bucket {
		if (i.Name == cReq.Filter.Bucket.Name) || (i.ID == cReq.Filter.Bucket.ID) {
			filtered = append(filtered, i)
		}
	}
	result.Bucket = filtered
	sendMsgResult(&w, &result)
}

// BucketShow function
func (x *Rest) BucketShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionBucket
	request.Action = msg.ActionShow
	request.Bucket = proto.Bucket{
		ID: params.ByName(`bucket`),
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`bucket_r`)
	handler.Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// BucketCreate function
func (x *Rest) BucketCreate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewBucketRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	if cReq.Bucket.Name == `` || cReq.Bucket.Environment == `` ||
		cReq.Bucket.TeamID == `` || cReq.Bucket.RepositoryID == `` {
		dispatchBadRequest(&w,
			fmt.Errorf(`Incomplete Bucket.Create request`))
		return
	}

	nameLen := utf8.RuneCountInString(cReq.Bucket.Name)
	if nameLen < 4 || nameLen > 512 {
		dispatchBadRequest(&w,
			fmt.Errorf(`Illegal bucket name length (4 < x <= 512)`))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionBucket
	request.Action = msg.ActionCreate
	request.Bucket = cReq.Bucket.Clone()

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`guidepost`)
	handler.Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// BucketPropertyCreate function
func (x *Rest) BucketPropertyCreate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewBucketRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	switch {
	case params.ByName(`bucket`) != cReq.Bucket.ID:
		dispatchBadRequest(&w,
			fmt.Errorf("Mismatched bucket ids: %s, %s",
				params.ByName(`bucket`),
				cReq.Bucket.ID))
		return
	case len(*cReq.Bucket.Properties) != 1:
		dispatchBadRequest(&w,
			fmt.Errorf("Expected property count 1, actual count: %d",
				len(*cReq.Bucket.Properties)))
		return
	case params.ByName(`type`) != (*cReq.Bucket.Properties)[0].Type:
		dispatchBadRequest(&w,
			fmt.Errorf("Mismatched property types: %s, %s",
				params.ByName(`type`),
				(*cReq.Bucket.Properties)[0].Type))
		return
	case (params.ByName(`type`) == `service`) && (*cReq.Bucket.Properties)[0].Service.Name == ``:
		dispatchBadRequest(&w,
			fmt.Errorf(`Empty service name is invalid`))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionBucket
	request.Action = msg.ActionPropertyCreate
	request.Bucket = cReq.Bucket.Clone()

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`guidepost`)
	handler.Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// BucketPropertyDestroy function
func (x *Rest) BucketPropertyDestroy(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionBucket
	request.Action = msg.ActionPropertyDestroy
	request.Bucket = proto.Bucket{
		ID: params.ByName(`bucket`),
		Properties: &[]proto.Property{
			proto.Property{
				Type:             params.ByName(`type`),
				BucketID:         params.ByName(`bucket`),
				SourceInstanceID: params.ByName(`source`),
			},
		},
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`guidepost`)
	handler.Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
