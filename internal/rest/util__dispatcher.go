/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"runtime/debug"

	"github.com/mjolnir42/soma/lib/auth"
	"github.com/mjolnir42/soma/lib/proto"
)

func panicCatcher(w http.ResponseWriter) {
	if r := recover(); r != nil {
		log.Printf("%s\n", debug.Stack())
		msg := fmt.Sprintf("PANIC! %s", r)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}

func peekJSONBody(r *http.Request, s interface{}) error {
	var err error
	body, _ := ioutil.ReadAll(r.Body)

	decoder := json.NewDecoder(
		ioutil.NopCloser(bytes.NewReader(body)),
	)
	r.Body = ioutil.NopCloser(bytes.NewReader(body))

	switch s.(type) {
	case *proto.Request:
		c := s.(*proto.Request)
		err = decoder.Decode(c)
	case *auth.Kex:
		c := s.(*auth.Kex)
		err = decoder.Decode(c)
	default:
		rt := reflect.TypeOf(s)
		err = fmt.Errorf("DecodeJsonBody: Unhandled request type: %s", rt)
	}
	return err
}

func decodeJSONBody(r *http.Request, s interface{}) error {
	decoder := json.NewDecoder(r.Body)
	var err error

	switch s.(type) {
	case *proto.Request:
		c := s.(*proto.Request)
		err = decoder.Decode(c)
	case *auth.Kex:
		c := s.(*auth.Kex)
		err = decoder.Decode(c)
	default:
		rt := reflect.TypeOf(s)
		err = fmt.Errorf("DecodeJsonBody: Unhandled request type: %s", rt)
	}
	return err
}

func dispatchBadRequest(w *http.ResponseWriter, err error) {
	if err != nil {
		http.Error(*w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Error(*w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func dispatchJSONReply(w *http.ResponseWriter, b *[]byte) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusOK)
	(*w).Write(*b)
}

func dispatchNoContent(w *http.ResponseWriter) {
	(*w).WriteHeader(http.StatusNoContent)
	(*w).Write(nil)
}

func dispatchInternalError(w *http.ResponseWriter, err error) {
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Error(*w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func dispatchNotImplemented(w *http.ResponseWriter, err error) {
	if err != nil {
		http.Error(*w, err.Error(), http.StatusNotImplemented)
		return
	}
	http.Error(*w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func dispatchUnauthorized(w *http.ResponseWriter, err error) {
	if err != nil {
		http.Error(*w, err.Error(), http.StatusUnauthorized)
		return
	}
	http.Error(*w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

func dispatchNotFound(w *http.ResponseWriter, err error) {
	if err != nil {
		http.Error(*w, err.Error(), http.StatusNotFound)
		return
	}
	http.Error(*w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func dispatchConflict(w *http.ResponseWriter, err error) {
	if err != nil {
		http.Error(*w, err.Error(), http.StatusConflict)
		return
	}
	http.Error(*w, http.StatusText(http.StatusConflict), http.StatusConflict)
}

func dispatchOctetReply(w *http.ResponseWriter, b *[]byte) {
	(*w).Header().Set("Content-Type", `application/octet-stream`)
	(*w).WriteHeader(http.StatusOK)
	(*w).Write(*b)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
