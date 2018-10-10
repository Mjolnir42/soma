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
	"github.com/mjolnir42/soma/internal/help"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerStatus(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// status
			{
				Name:        `status`,
				Usage:       `SUBCOMMANDS for check instance status`,
				Description: help.Text(`status::`),
				Subcommands: []cli.Command{
					{
						Name:        `add`,
						Usage:       `Add a check instance status`,
						Description: help.Text(`status::add`),
						Action:      runtime(cmdStatusAdd),
					},
					{
						Name:        `remove`,
						Usage:       `Remove a check instance status`,
						Description: help.Text(`status::remove`),
						Action:      runtime(cmdStatusRemove),
					},
					{
						Name:        `list`,
						Usage:       `List check instance status`,
						Description: help.Text(`status::list`),
						Action:      runtime(cmdStatusList),
					},
					{
						Name:        `show`,
						Usage:       `Show details about a check instance status`,
						Description: help.Text(`status::show`),
						Action:      runtime(cmdStatusShow),
					},
				},
			}, // end status
		}...,
	)
	return &app
}

// cmdStatusAdd function
// soma status add ${status}
func cmdStatusAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	req := proto.NewStatusRequest()
	req.Status.Name = c.Args().First()

	return adm.Perform(`postbody`, `/status/`, `command`, req, c)
}

// cmdStatusRemove function
// soma status remove ${status}
func cmdStatusRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/status/%s", esc)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// cmdStatusList function
// soma status list
func cmdStatusList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/status/`, `list`, nil, c)
}

// cmdStatusShow function
// soma status show ${status}
func cmdStatusShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/status/%s", esc)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
