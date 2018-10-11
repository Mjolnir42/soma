/*-
 * Copyright (c) 2016-2018, Jörg Pernfuß
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

func registerDatacenters(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// datacenters
			{
				Name:        `datacenter`,
				Usage:       `SUBCOMMANDS for datacenters`,
				Description: help.Text(`datacenter::`),
				Subcommands: []cli.Command{
					{
						Name:        `add`,
						Usage:       `Register a new datacenter`,
						Description: help.Text(`datacenter::add`),
						Action:      runtime(cmdDatacenterAdd),
					},
					{
						Name:        `remove`,
						Usage:       `Remove an existing datacenter`,
						Description: help.Text(`datacenter::remove`),
						Action:      runtime(cmdDatacenterRemove),
					},
					{
						Name:         `rename`,
						Usage:        `Rename an existing datacenter`,
						Description:  help.Text(`datacenter::rename`),
						Action:       runtime(cmdDatacenterRename),
						BashComplete: cmpl.To,
					},
					{
						Name:        `list`,
						Usage:       `List all datacenters`,
						Description: help.Text(`datacenter::list`),
						Action:      runtime(cmdDatacenterList),
					},
					{
						Name:        `show`,
						Usage:       `Show information about a specific datacenter`,
						Description: help.Text(`datacenter::show`),
						Action:      runtime(cmdDatacenterShow),
					},
					{
						Name:        `sync`,
						Usage:       `List all datacenters in a format suitable for sync`,
						Description: help.Text(`datacenter::sync`),
						Action:      runtime(cmdDatacenterSync),
					},
				},
			},
		}...,
	)
	return &app
}

// cmdDatacenterAdd function
// soma datacenter add ${datacenter}
func cmdDatacenterAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.NewDatacenterRequest()
	req.Datacenter.LoCode = c.Args().First()

	if err := adm.ValidateNoSlash(req.Datacenter.LoCode); err != nil {
		return err
	}

	return adm.Perform(`postbody`, `/datacenter/`, `command`, req, c)
}

// cmdDatacenterRemove function
// soma datacenter remove ${datacenter}
func cmdDatacenterRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/datacenter/%s", esc)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// cmdDatacenterRename function
// soma datacenter rename ${old} to ${new}
func cmdDatacenterRename(c *cli.Context) error {
	key := []string{`to`}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(opts, key, key, key,
		c.Args().Tail()); err != nil {
		return err
	}

	req := proto.NewDatacenterRequest()
	req.Datacenter.LoCode = opts[`to`][0]

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}
	if err := adm.ValidateNoSlash(opts[`to`][0]); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/datacenter/%s", esc)
	return adm.Perform(`put`, path, `command`, nil, c)
}

// cmdDatacenterList function
// soma datacenter list
func cmdDatacenterList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/datacenter/`, `list`, nil, c)
}

// cmdDatacenterSync function
// soma datacenter sync
func cmdDatacenterSync(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/sync/datacenter/`, `list`, nil, c)
}

// cmdDatacenterShow function
// soma datacenter show ${datacenter}
func cmdDatacenterShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/datacenter/%s", esc)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
