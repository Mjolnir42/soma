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

func registerValidity(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:        `validity`,
				Usage:       `SUBCOMMANDS for system property validity`,
				Description: help.Text(`validity::`),
				Subcommands: []cli.Command{
					{
						Name:         `add`,
						Usage:        `Add a system property validity definition`,
						Description:  help.Text(`validity::add`),
						Action:       runtime(validityAdd),
						BashComplete: cmpl.ValidityAdd,
					},
					{
						Name:        `remove`,
						Usage:       `Remove a system property validity definition`,
						Description: help.Text(`validity::remove`),
						Action:      runtime(validityRemove),
					},
					{
						Name:        `list`,
						Usage:       `List system property validity definitions`,
						Description: help.Text(`validity::list`),
						Action:      runtime(validityList),
					},
					{
						Name:        `show`,
						Usage:       `Show details about a system property validity`,
						Description: help.Text(`validity::show`),
						Action:      runtime(validityShow),
					},
				},
			},
		}...,
	)
	return &app
}

// validityAdd function
// soma validity add ${property} on ${entity} \
//      direct ${bool} inherited ${bool}
func validityAdd(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`on`, `direct`, `inherited`}
	mandatoryOptions := []string{`on`, `direct`, `inherited`}

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

	req := proto.NewValidityRequest()
	req.Validity.SystemProperty = c.Args().First()
	req.Validity.Entity = opts[`on`][0]
	if err := adm.ValidateBool(opts[`direct`][0],
		&req.Validity.Direct); err != nil {
		return err
	}
	if err := adm.ValidateBool(opts[`inherited`][0],
		&req.Validity.Inherited); err != nil {
		return err
	}

	return adm.Perform(`postbody`, `/validity/`, `command`, req, c)
}

// validityRemove function
// soma validity remove ${property}
func validityRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	path := fmt.Sprintf("/validity/%s",
		url.QueryEscape(c.Args().First()),
	)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// validityList function
// soma validity list
func validityList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/validity/`, `list`, nil, c)
}

// validityShow function
// soma validity show ${property}
func validityShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	path := fmt.Sprintf("/validity/%s",
		url.QueryEscape(c.Args().First()),
	)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
