/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
)

func registerInstanceMgmt(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  `instance-mgmt`,
				Usage: `SUBCOMMANDS for check instance management`,
				Subcommands: []cli.Command{
					{
						Name:   `list`,
						Usage:  `List all check instances`,
						Action: runtime(cmdInstanceMgmtAll),
					},
					{
						Name:   `show`,
						Usage:  `Show details about a check instance`,
						Action: runtime(cmdInstanceMgmtShow),
					},
					{
						Name:   `versions`,
						Usage:  `Show version history of a check instance`,
						Action: runtime(cmdInstanceMgmtVersion),
					},
				},
			},
		}...,
	)
	return &app
}

func cmdInstanceMgmtAll(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/instance/`, `list`, nil, c)
}

func cmdInstanceMgmtShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	if !adm.IsUUID(c.Args().First()) {
		return fmt.Errorf("Argument is not a UUID: %s",
			c.Args().First())
	}

	path := fmt.Sprintf("/instance/%s", c.Args().First())
	return adm.Perform(`get`, path, `show`, nil, c)
}

func cmdInstanceMgmtVersion(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	if !adm.IsUUID(c.Args().First()) {
		return fmt.Errorf("Argument is not a UUID: %s",
			c.Args().First())
	}

	path := fmt.Sprintf("/instance/%s/versions", c.Args().First())
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
