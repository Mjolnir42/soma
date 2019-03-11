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

// bashCompNode feeds cmpl with information retrieved via the
// full execution runtime, allowing completion of node names
func bashCompNode(c *cli.Context) {
	cmpl.GenericDataOnly(c, filterNodeName(nodeFetch()))
}

// bashCompNodeAssign calls the completion for node::assign commands
// with keywords and node name data
func bashCompNodeAssign(c *cli.Context) {
	cmpl.Augmented(c, `to`, filterNodeName(nodeFetch()))
}

// bashCompNodeUnassign calls the completion for node::unassign commands
// with keywords and node name data
func bashCompNodeUnassign(c *cli.Context) {
	cmpl.Augmented(c, `from`, filterNodeName(nodeFetch()))
}

// bashCompNodeRepossess calls the completion for node-mgmt::repossess
// commands with keywords and name data
func bashCompNodeRepossess(c *cli.Context) {
	cmpl.Augmented(c, `to`, filterNodeName(nodeFetch()))
}

// bashCompNodeRename calls the completion for node-mgmt::rename
// commands with keywords and name data
func bashCompNodeRename(c *cli.Context) {
	cmpl.Augmented(c, `to`, filterNodeName(nodeFetch()))
}

// bashCompNodeRelocate calls the completion for node-mgmt::relocate
// commands with keywords and name data
func bashCompNodeRelocate(c *cli.Context) {
	cmpl.Augmented(c, `to`, filterNodeName(nodeFetch()))
}

// bashCompNodeConfigTree calls the completion for node-config::tree
// commands with keywords and node name data
func bashCompNodeConfigTree(c *cli.Context) {
	cmpl.Augmented(c, `in`, filterNodeName(nodeFetch()))
}

// nodeFetch returns a slice of possible nodes, queried from the SOMA
// server. On errors an empty slice is returned and the error is pushed
// onto the error stack.
func nodeFetch() []proto.Node {
	var err error
	var resp *resty.Response
	res := &proto.Result{}

	// fetch nodes list
	if resp, err = adm.GetReq(`/node/`); err != nil {
		pushError(err)
		return []proto.Node{}
	}

	// extract result
	if err = adm.DecodedResponse(resp, res); err != nil {
		pushError(err)
		return []proto.Node{}
	}

	// check result
	if res == nil {
		pushError(fmt.Errorf(`adm.DecodedResponse returned nil object`))
		return []proto.Node{}
	}

	if res.Nodes == nil {
		pushError(fmt.Errorf(`SOMA server returned no nodes for parsing`))
		return []proto.Node{}
	}

	if len(*res.Nodes) == 0 {
		return []proto.Node{}
	}

	return *res.Nodes
}

// filterNodeName returns the node names contained in data
func filterNodeName(data []proto.Node) (res []string) {
	res = make([]string, len(data))
	for i := range data {
		res[i] = data[i].Name
	}
	return
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
