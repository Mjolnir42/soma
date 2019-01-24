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
	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/internal/help"
)

func registerNodes(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// nodes
			{
				Name:  "nodes",
				Usage: "SUBCOMMANDS for nodes",
				Subcommands: []cli.Command{
					{
						Name:        `list`,
						Usage:       `List all nodes`,
						Description: help.Text(`node::list`),
						Action:      runtime(nodeList),
					},
					{
						Name:        `show`,
						Usage:       `Show the details for a specifc node`,
						Description: help.Text(`node::show`),
						Action:      runtime(nodeShow),
					},
					{
						Name:        `config`,
						Usage:       `Show the repository/bucket assignment of a specific node`,
						Description: help.Text(`node::show-config`),
						Action:      runtime(nodeShowConfig),
					},
					{
						Name:         `assign`,
						Usage:        `Assign a node to configuration bucket`,
						Description:  help.Text(`node::assign`),
						Action:       runtime(nodeAssign),
						BashComplete: cmpl.To,
					},
					{
						Name:         `unassign`,
						Usage:        `Unassign a node from its configuration bucket`,
						Description:  help.Text(`node::unassign`),
						Action:       runtime(nodeUnassign),
						BashComplete: cmpl.From,
					},
					{
						Name:         `dumptree`,
						Usage:        `List the node as a tree`,
						Description:  help.Text(`node-config::tree`),
						Action:       runtime(nodeConfigTree),
						BashComplete: cmpl.In,
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
