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

func registerModes(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// modes
			{
				Name:        `mode`,
				Usage:       `SUBCOMMANDS for monitoring system modes`,
				Description: help.Text(`mode::`),
				Subcommands: []cli.Command{
					{
						Name:        `add`,
						Usage:       `Add a new monitoring system mode`,
						Description: help.Text(`mode::add`),
						Action:      runtime(cmdModeAdd),
					},
					{
						Name:        `remove`,
						Usage:       `Remove a monitoring system mode`,
						Description: help.Text(`mode::remove`),
						Action:      runtime(cmdModeRemove),
					},
					{
						Name:        `list`,
						Usage:       `List monitoring system modes`,
						Description: help.Text(`mode::list`),
						Action:      runtime(cmdModeList),
					},
					{
						Name:        `show`,
						Usage:       `Show details about a monitoring mode`,
						Description: help.Text(`mode::show`),
						Action:      runtime(cmdModeShow),
					},
				},
			}, // end modes
		}...,
	)
	return &app
}

// cmdModeAdd function
// soma mode add ${mode}
func cmdModeAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.NewModeRequest()
	req.Mode.Mode = c.Args().First()

	return adm.Perform(`postbody`, `/mode/`, `command`, req, c)
}

// cmdModeRemove function
// soma mode remove ${mode}
func cmdModeRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/mode/%s", url.QueryEscape(c.Args().First()))
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// cmdModeList function
// soma mode list
func cmdModeList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/mode/`, `list`, nil, c)
}

// cmdModeShow function
// soma mode show ${mode}
func cmdModeShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/mode/%s", url.QueryEscape(c.Args().First()))
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
