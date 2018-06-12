/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
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
						Name:   `all`,
						Usage:  `List all check instances`,
						Action: runtime(cmdInstanceMgmtAll),
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
