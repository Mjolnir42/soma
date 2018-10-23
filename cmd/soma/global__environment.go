/*-
 * Copyright (c) 2015-2018, Jörg Pernfuß
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

func registerEnvironments(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// environments
			{
				Name:        `environment`,
				Usage:       `SUBCOMMANDS for environment definitions`,
				Description: help.Text(`environment::`),
				Subcommands: []cli.Command{
					{
						Name:        `add`,
						Usage:       `Define a new environment`,
						Description: help.Text(`environment::add`),
						Action:      runtime(environmentAdd),
					},
					{
						Name:        `remove`,
						Usage:       `Remove an existing unused environment`,
						Description: help.Text(`environment::remove`),
						Action:      runtime(environmentRemove),
					},
					{
						Name:         `rename`,
						Usage:        `Rename an existing environment`,
						Description:  help.Text(`environment::rename`),
						Action:       runtime(environmentRename),
						BashComplete: cmpl.To,
					},
					{
						Name:        `list`,
						Usage:       `List all available environments`,
						Description: help.Text(`environment::list`),
						Action:      runtime(environmentList),
					},
					{
						Name:        `show`,
						Usage:       `Show information about a specific environment`,
						Description: help.Text(`environment::show`),
						Action:      runtime(environmentShow),
					},
				},
			},
		}...,
	)
	return &app
}

// environmentAdd function
// soma environment add ${environment}
func environmentAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	req := proto.NewEnvironmentRequest()
	req.Environment.Name = c.Args().First()

	return adm.Perform(`postbody`, `/environment/`, `command`, req, c)
}

// environmentRemove function
// soma environment remove ${environment}
func environmentRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	path := fmt.Sprintf(
		"/environment/%s",
		url.QueryEscape(c.Args().First()),
	)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// environmentRename function
// soma environment rename ${old} to ${new}
func environmentRename(c *cli.Context) error {
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

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(opts[`to`][0]); err != nil {
		return err
	}

	req := proto.NewEnvironmentRequest()
	req.Environment.Name = opts[`to`][0]

	path := fmt.Sprintf(
		"/environment/%s",
		url.QueryEscape(c.Args().First()),
	)
	return adm.Perform(`putbody`, path, `command`, req, c)
}

// environmentList function
// soma environment list
func environmentList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/environment/`, `list`, nil, c)
}

// environmentShow function
// soma environment show ${environment}
func environmentShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf(
		"/environment/%s",
		url.QueryEscape(c.Args().First()),
	)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
