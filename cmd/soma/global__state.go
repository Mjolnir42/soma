/*-
 * Copyright (c) 2015-2018, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main // import "github.com/mjolnir42/soma/cmd/somaadm"

import (
	"fmt"
	"net/url"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/internal/help"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerStates(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:        `state`,
				Usage:       `SUBCOMMANDS for states`,
				Description: help.Text(`state`),
				Subcommands: []cli.Command{
					{
						Name:        `add`,
						Usage:       `Add a new object state`,
						Description: help.Text(`state::add`),
						Action:      runtime(cmdStateAdd),
					},
					{
						Name:        `remove`,
						Usage:       `Remove an existing object state`,
						Description: help.Text(`state::remove`),
						Action:      runtime(cmdStateRemove),
					},
					{
						Name:         `rename`,
						Usage:        `Rename an existing object state`,
						Description:  help.Text(`state::rename`),
						Action:       runtime(cmdStateRename),
						BashComplete: cmpl.To,
					},
					{
						Name:        `list`,
						Usage:       `List all object states`,
						Description: help.Text(`state::list`),
						Action:      runtime(cmdStateList),
					},
					{
						Name:        `show`,
						Usage:       `Show information about an object states`,
						Description: help.Text(`state::show`),
						Action:      runtime(cmdStateShow),
					},
				},
			}, // end states
		}...,
	)
	return &app
}

// cmdStateAdd function
// somaadm state add ${state}
func cmdStateAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.NewStateRequest()
	req.State.Name = c.Args().First()

	return adm.Perform(`postbody`, `/state/`, `command`, req, c)
}

// cmdStateRemove function
// somaadm state remove ${state}
func cmdStateRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/state/%s", esc)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// cmdStateRename function
// somaadm state rename ${state} to ${new-state}
func cmdStateRename(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`to`}
	mandatoryOptions := []string{`to`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	req := proto.NewStateRequest()
	req.State.Name = opts[`to`][0]

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/state/%s", esc)
	return adm.Perform(`putbody`, path, `command`, req, c)
}

// cmdStateList function
// somaadm state list
func cmdStateList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/state/`, `list`, nil, c)
}

// cmdStateShow function
// somaadm state show ${state}
func cmdStateShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/state/%s", esc)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
