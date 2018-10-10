/*-
 * Copyright (c) 2016-2018, Jörg Pernfuß
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

func registerUnits(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:        `unit`,
				Usage:       `SUBCOMMANDS for metric units`,
				Description: help.Text(`unit::`),
				Subcommands: []cli.Command{
					{
						Name:         `add`,
						Usage:        `Add a new metric unit`,
						Description:  help.Text(`unit::add`),
						Action:       runtime(cmdUnitAdd),
						BashComplete: cmpl.Name,
					},
					{
						Name:        `remove`,
						Usage:       `Remove a metric unit`,
						Description: help.Text(`unit::remove`),
						Action:      runtime(cmdUnitRemove),
					},
					{
						Name:        `list`,
						Usage:       `List metric units`,
						Description: help.Text(`unit::list`),
						Action:      runtime(cmdUnitList),
					},
					{
						Name:        `show`,
						Usage:       `Show details about a metric unit`,
						Description: help.Text(`unit::show`),
						Action:      runtime(cmdUnitShow),
					},
				},
			},
		}...,
	)
	return &app
}

// cmdUnitAdd function
// soma unit add ${unit} name ${name}
func cmdUnitAdd(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`name`}
	mandatoryOptions := []string{`name`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	req := proto.NewUnitRequest()
	req.Unit.Unit = c.Args().First()
	req.Unit.Name = opts[`name`][0]

	return adm.Perform(`postbody`, `/unit/`, `command`, req, c)
}

// cmdUnitRemove function
// soma unit remove ${unit}
func cmdUnitRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/unit/%s", esc)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// cmdUnitList function
// soma unit list
func cmdUnitList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/unit/`, `list`, nil, c)
}

// cmdUnitShow function
// soma unit show ${unit}
func cmdUnitShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/unit/%s", esc)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
