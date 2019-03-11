/*-
 * Copyright (c) 2019, Jörg Pernfuß
 * Copyright (c) 2019, 1&1 IONOS SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main // import "github.com/mjolnir42/soma/cmd/soma"

import (
	"fmt"

	resty "gopkg.in/resty.v0"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/lib/proto"
)

// bashCompServer feeds cmpl with information retrieved via the
// full execution runtime.
func bashCompServer(c *cli.Context) {
	cmpl.GenericDataOnly(c, filterServerName(serverFetch()))
}

// serverFetch returns a slice of possible servers, queried from the
// SOMA server. On errors an empty slice is returned and the error is
// pushed onto the error stack.
func serverFetch() []proto.Server {
	var err error
	var resp *resty.Response
	res := &proto.Result{}

	// fetch server list
	if resp, err = adm.GetReq(`/server/`); err != nil {
		pushError(err)
		return []proto.Server{}
	}

	// extract result
	if err = adm.DecodedResponse(resp, res); err != nil {
		pushError(err)
		return []proto.Server{}
	}

	// check result
	if res == nil {
		pushError(fmt.Errorf(`adm.DecodedResponse returned nil object`))
		return []proto.Server{}
	}

	if res.Servers == nil {
		pushError(fmt.Errorf(`SOMA server returned no servers for parsing`))
		return []proto.Server{}
	}

	if len(*res.Servers) == 0 {
		return []proto.Server{}
	}

	return *res.Servers
}

// filterServerName returns the server names contained in data
func filterServerName(data []proto.Server) (res []string) {
	res = make([]string, len(data))
	for i := range data {
		res[i] = data[i].Name
	}
	return
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
