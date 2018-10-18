/*-
 * Copyright (c) 2015-2018, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main // import "github.com/mjolnir42/soma/cmd/soma"

import (
	"fmt"

	resty "gopkg.in/resty.v0"

	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/lib/proto"
)

// attributeFetch returns a slice of possible service attributes queried
// from the SOMA server. On errors an empty slice is returned and the
// error is pushed onto the error stack.
func attributeFetch() []proto.Attribute {
	var (
		err  error
		resp *resty.Response
	)
	res := &proto.Result{}

	// fetch list of possible service attributes
	if resp, err = adm.GetReq(`/attribute/`); err != nil {
		pushError(err)
		return []proto.Attribute{}
	}

	// extract result
	if err = adm.DecodedResponse(resp, res); err != nil {
		pushError(err)
		return []proto.Attribute{}
	}

	// check result
	if res == nil {
		pushError(fmt.Errorf(`adm.DecodedResponse returned nil object`))
		return []proto.Attribute{}
	}

	if res.Attributes == nil || len(*res.Attributes) == 0 {
		pushError(fmt.Errorf(`server returned no attributes for parsing`))
		return []proto.Attribute{}
	}

	return *res.Attributes
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
