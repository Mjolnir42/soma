/*-
 * Copyright (c) 2015-2019, Jörg Pernfuß
 * Copyright (c) 2019, 1&1 IONOS SE
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main // import "github.com/mjolnir42/soma/cmd/soma"

import (
	"fmt"
	"net/url"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/internal/help"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerNodes(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// nodes
			{
				Name:        `node`,
				Usage:       `SUBCOMMANDS for node management`,
				Description: help.Text(`node::`),
				Subcommands: []cli.Command{
					{
						Name:         `add`,
						Usage:        `Register a new node with SOMA`,
						Description:  help.Text(`node-mgmt::add`),
						Action:       runtime(nodeMgmtAdd),
						BashComplete: cmpl.NodeAdd,
					},
					{
						Name:         `remove`,
						Usage:        `Mark a node as deleted`,
						Description:  help.Text(`node-mgmt::remove`),
						Action:       runtime(nodeMgmtRemove),
						BashComplete: comptime(bashCompNode),
					},
					{
						Name:         `purge`,
						Usage:        `Purge a node marked as deleted`,
						Description:  help.Text(`node-mgmt::purge`),
						Action:       runtime(nodeMgmtPurge),
						BashComplete: comptime(bashCompNode),
						Flags: []cli.Flag{
							cli.BoolFlag{
								Name:  `all, a`,
								Usage: `Purge all deleted nodes`,
							},
						},
					},
					{
						Name:         `update`,
						Usage:        `Update a node's information`,
						Description:  help.Text(`node-mgmt::update`),
						Action:       runtime(nodeMgmtUpdate),
						BashComplete: cmpl.NodeUpdate,
					},
					{
						Name:         `repossess`,
						Usage:        `Repossess a node to a different team`,
						Description:  help.Text(`node-mgmt::repossess`),
						Action:       runtime(nodeMgmtRepossess),
						BashComplete: comptime(bashCompNodeRepossess),
					},
					{
						Name:         `rename`,
						Usage:        `Rename a node`,
						Description:  help.Text(`node-mgmt::rename`),
						Action:       runtime(nodeMgmtRename),
						BashComplete: comptime(bashCompNodeRename),
					},
					{
						Name:         `relocate`,
						Usage:        `Relocate a node to a different server`,
						Description:  help.Text(`node-mgmt::relocate`),
						Action:       runtime(nodeMgmtRelocate),
						BashComplete: comptime(bashCompNodeRelocate),
					},
					{
						Name:         `list`,
						Usage:        `List all nodes`,
						Description:  help.Text(`node::list`),
						Action:       runtime(nodeList),
						BashComplete: cmpl.None,
					},
					{
						Name:         `show`,
						Usage:        `Show the details for a specifc node`,
						Description:  help.Text(`node::show`),
						Action:       runtime(nodeShow),
						BashComplete: comptime(bashCompNode),
					},
					{
						Name:         `sync`,
						Usage:        `List all nodes with data suitable for sync operations`,
						Description:  help.Text(`node-mgmt::sync`),
						Action:       runtime(nodeMgmtSync),
						BashComplete: cmpl.None,
					},
					{
						Name:         `config`,
						Usage:        `Show the repository/bucket assignment of a specific node`,
						Description:  help.Text(`node::show-config`),
						Action:       runtime(nodeShowConfig),
						BashComplete: comptime(bashCompNode),
					},
					{
						Name:         `assign`,
						Usage:        `Assign a node to configuration bucket`,
						Description:  help.Text(`node::assign`),
						Action:       runtime(nodeAssign),
						BashComplete: comptime(bashCompNodeAssign),
					},
					{
						Name:         `unassign`,
						Usage:        `Unassign a node from its configuration bucket`,
						Description:  help.Text(`node::unassign`),
						Action:       runtime(nodeUnassign),
						BashComplete: comptime(bashCompNodeUnassign),
					},
					{
						Name:         `dumptree`,
						Usage:        `List the node as a tree`,
						Description:  help.Text(`node-config::tree`),
						Action:       runtime(nodeConfigTree),
						BashComplete: comptime(bashCompNodeConfigTree),
					},
					{
						Name:        `property`,
						Usage:       `SUBCOMMANDS for properties on nodes`,
						Description: help.Text(`node-config::`),
						Subcommands: []cli.Command{
							{
								Name:        `create`,
								Usage:       `SUBCOMMANDS to create properties`,
								Description: help.Text(`node-config::property-create`),
								Subcommands: []cli.Command{
									{
										Name:         `system`,
										Usage:        `Add a system property to a node`,
										Description:  help.Text(`node-config::property-create`),
										Action:       runtime(nodeConfigPropertyCreateSystem),
										BashComplete: cmpl.PropertyCreateInValue,
									},
									{
										Name:         `service`,
										Usage:        `Add a service property to a node`,
										Description:  help.Text(`node-config::property-create`),
										Action:       runtime(nodeConfigPropertyCreateService),
										BashComplete: cmpl.PropertyCreateInValue,
									},
									{
										Name:         `oncall`,
										Usage:        `Add an oncall property to a node`,
										Description:  help.Text(`node-config::property-create`),
										Action:       runtime(nodeConfigPropertyCreateOncall),
										BashComplete: cmpl.PropertyCreateIn,
									},
									{
										Name:         `custom`,
										Usage:        `Add a custom property to a node`,
										Description:  help.Text(`node-config::property-create`),
										Action:       runtime(nodeConfigPropertyCreateCustom),
										BashComplete: cmpl.PropertyCreateIn,
									},
								},
							},
							{
								Name:        `destroy`,
								Usage:       `SUBCOMMANDS to destroy properties`,
								Description: help.Text(`node-config::property-destroy`),
								Subcommands: []cli.Command{
									{
										Name:         `system`,
										Usage:        `Delete a system property from a node`,
										Description:  help.Text(`node-config::property-destroy`),
										Action:       runtime(nodeConfigPropertyDestroySystem),
										BashComplete: cmpl.PropertyOnInView,
									},
									{
										Name:         `service`,
										Usage:        `Delete a service property from a node`,
										Description:  help.Text(`node-config::property-destroy`),
										Action:       runtime(nodeConfigPropertyDestroyService),
										BashComplete: cmpl.PropertyOnInView,
									},
									{
										Name:         `oncall`,
										Usage:        `Delete an oncall property from a node`,
										Description:  help.Text(`node-config::property-destroy`),
										Action:       runtime(nodeConfigPropertyDestroyOncall),
										BashComplete: cmpl.PropertyOnInView,
									},
									{
										Name:         `custom`,
										Usage:        `Delete a custom property from a node`,
										Description:  help.Text(`node-config::property-destroy`),
										Action:       runtime(nodeConfigPropertyDestroyCustom),
										BashComplete: cmpl.PropertyOnInView,
									},
								},
							},
						},
					},
				},
			},
		}...,
	)
	return &app
}

// nodeMgmtAdd function
// soma node add ${node}   \
//      assetid ${id}      \
//      team ${team}       \
//      [server ${server}] \
//      [online ${isOnline}]
func nodeMgmtAdd(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`assetid`, `team`, `server`, `online`}
	mandatoryOptions := []string{`assetid`, `team`}

	var err error
	if err = adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	req := proto.NewNodeRequest()

	// optional arguments, defaults to true
	if _, ok := opts[`online`]; ok {
		if err := adm.ValidateBool(opts[`online`][0],
			&req.Node.IsOnline); err != nil {
			return err
		}
	} else {
		req.Node.IsOnline = true
	}

	// optional argument, default assignment is handled on the server
	if _, ok := opts[`server`]; ok {
		if req.Node.ServerID, err = adm.LookupServerID(
			opts[`server`][0]); err != nil {
			return err
		}
	}

	req.Node.Name = c.Args().First()
	if err := adm.ValidateNotUUID(req.Node.Name); err != nil {
		return err
	}
	if err = adm.LookupTeamID(
		opts[`team`][0],
		&req.Node.TeamID,
	); err != nil {
		return nil
	}

	if err = adm.ValidateLBoundUint64(opts[`assetid`][0],
		&req.Node.AssetID, 1); err != nil {
		return err
	}

	return adm.Perform(`postbody`, `/node/`, `node-mgmt::add`, req, c)
}

// nodeMgmtRemove function
// soma node remove ${node}
func nodeMgmtRemove(c *cli.Context) (err error) {
	// check deferred errors
	if err := popError(); err != nil {
		return err
	}
	req := proto.Request{Flags: &proto.Flags{
		Purge: false,
	}}

	if err = adm.VerifySingleArgument(c); err != nil {
		return err
	}
	var id, path string
	if id, err = adm.LookupNodeID(c.Args().First()); err != nil {
		return err
	}
	path = fmt.Sprintf("/node/%s", id)

	return adm.Perform(`deletebody`, path, `node-mgmt::remove`, req, c)
}

// nodeMgmtPurge function
// soma node purge [-a|--all] [${node}]
func nodeMgmtPurge(c *cli.Context) (err error) {
	// check deferred errors
	if err := popError(); err != nil {
		return err
	}

	var path string
	req := proto.Request{Flags: &proto.Flags{
		Purge: true,
	}}

	switch c.Bool(`all`) {
	case true:
		if err := adm.VerifyNoArgument(c); err != nil {
			return err
		}
		path = `node`
	default:
		if err := adm.VerifySingleArgument(c); err != nil {
			return err
		}
		nodeID, err := adm.LookupNodeID(c.Args().First())
		if err != nil {
			return err
		}
		path = fmt.Sprintf("/node/%s",
			url.QueryEscape(nodeID),
		)
	}
	return adm.Perform(`deletebody`, path, `node-mgmt::purge`, req, c)
}

// nodeMgmtUpdate function
// soma node update ${nodeUUID} \
//      name ${name}            \
//      assetid ${id}           \
//      team ${team}            \
//      server ${server}        \
//      online ${isOnline}      \
//      deleted ${isDeleted}
func nodeMgmtUpdate(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`name`, `assetid`, `team`, `server`, `online`, `deleted`}
	mandatoryOptions := []string{`name`, `assetid`, `team`, `server`, `online`, `deleted`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	if err := adm.ValidateUUID(c.Args().First()); err != nil {
		return err
	}
	opts[`nodeID`] = []string{c.Args().First()}

	return nodeMgmtVariadicUpdate(c, opts)
}

// nodeMgmtRepossess function
// soma node repossess ${node} to ${team}
func nodeMgmtRepossess(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`to`}
	mandatoryOptions := []string{`to`}

	// check deferred errors
	if err := popError(); err != nil {
		return err
	}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	opts[`nodeName`] = []string{c.Args().First()}
	opts[`team`] = opts[`to`]
	opts[`to`] = []string{}

	return nodeMgmtVariadicUpdate(c, opts)
}

// nodeMgmtRename function
// soma node rename ${node} to ${name}
func nodeMgmtRename(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`to`}
	mandatoryOptions := []string{`to`}

	// check deferred errors
	if err := popError(); err != nil {
		return err
	}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	opts[`nodeName`] = []string{c.Args().First()}
	opts[`name`] = opts[`to`]
	opts[`to`] = []string{}

	return nodeMgmtVariadicUpdate(c, opts)
}

// nodeMgmtRelocate function
// soma node relocate ${node} to ${server}
func nodeMgmtRelocate(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`to`}
	mandatoryOptions := []string{`to`}

	// check deferred errors
	if err := popError(); err != nil {
		return err
	}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	opts[`nodeName`] = []string{c.Args().First()}
	opts[`server`] = opts[`to`]
	opts[`to`] = []string{}

	return nodeMgmtVariadicUpdate(c, opts)
}

// nodeMgmtVariadicUpdate function
// this function implements the various commands that update
// node data
func nodeMgmtVariadicUpdate(c *cli.Context, opts map[string][]string) error {
	req := proto.NewNodeRequest()
	// request is a full update
	if _, ok := opts[`nodeID`]; ok {
		req.Node.ID = opts[`nodeID`][0]
	}
	// request is a partial update
	if _, ok := opts[`nodeName`]; ok {
		nodeID, err := adm.LookupNodeID(opts[`nodeName`][0])
		if err != nil {
			return err
		}
		req.Node.ID = nodeID
	}
	node, err := adm.LookupNode(req.Node.ID)
	if err != nil {
		return err
	}

	if _, ok := opts[`assetid`]; ok {
		if err = adm.ValidateLBoundUint64(opts[`assetid`][0],
			&req.Node.AssetID, 1); err != nil {
			return err
		}
	} else {
		req.Node.AssetID = node.AssetID
	}

	if _, ok := opts[`name`]; ok {
		if err = adm.ValidateNotUUID(opts[`name`][0]); err != nil {
			return err
		}
		req.Node.Name = opts[`name`][0]
	} else {
		req.Node.Name = node.Name
	}

	if _, ok := opts[`team`]; ok {
		if err = adm.LookupTeamID(opts[`team`][0],
			&req.Node.TeamID,
		); err != nil {
			return err
		}
	} else {
		req.Node.TeamID = node.TeamID
	}

	if _, ok := opts[`server`]; ok {
		if req.Node.ServerID, err = adm.LookupServerID(
			opts[`server`][0]); err != nil {
			return err
		}
	} else {
		req.Node.ServerID = node.ServerID
	}

	if _, ok := opts[`online`]; ok {
		if err = adm.ValidateBool(opts[`online`][0],
			&req.Node.IsOnline); err != nil {
			return err
		}
	} else {
		req.Node.IsOnline = node.IsOnline
	}

	if _, ok := opts[`deleted`]; ok {
		if err = adm.ValidateBool(opts[`deleted`][0],
			&req.Node.IsDeleted); err != nil {
			return err
		}
	} else {
		req.Node.IsDeleted = node.IsDeleted
	}

	path := fmt.Sprintf("/node/%s",
		url.QueryEscape(req.Node.ID),
	)
	return adm.Perform(`putbody`, path, `node-mgmt::update`, req, c)
}

// nodeMgmtSync function
// soma node sync
func nodeMgmtSync(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/sync/node/`, `node-mgmt::sync`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
