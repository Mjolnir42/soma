/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/soma/internal/rest"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/mjolnir42/soma/lib/auth"
	"github.com/mjolnir42/soma/lib/proto"
)

// peekJSONBody unmarshals a copy of the JSON request body from r into
// s, leaving r.Body ready for another reader
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

// decodeJSONBody unmarshals the JSON request body from r into s
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
