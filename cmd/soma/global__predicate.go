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
	"github.com/mjolnir42/soma/internal/help"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerPredicates(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:        `predicate`,
				Usage:       `SUBCOMMANDS for threshold predicates`,
				Description: help.Text(`predicate::`),
				Subcommands: []cli.Command{
					{
						Name:        `add`,
						Usage:       `Add a threshold predicate`,
						Description: help.Text(`predicate::add`),
						Action:      runtime(predicateAdd),
					},
					{
						Name:        `remove`,
						Usage:       `Remove a threshold predicate`,
						Description: help.Text(`predicate::remove`),
						Action:      runtime(predicateRemove),
					},
					{
						Name:        `list`,
						Usage:       `List all threshold predicates`,
						Description: help.Text(`predicate::list`),
						Action:      runtime(predicateList),
					},
					{
						Name:        `show`,
						Usage:       `Show details about a threshold predicate`,
						Description: help.Text(`predicate::show`),
						Action:      runtime(predicateShow),
					},
				},
			},
		}...,
	)
	return &app
}

// predicateAdd function
// soma predicate add ${pred}
func predicateAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	req := proto.NewPredicateRequest()
	req.Predicate.Symbol = c.Args().First()

	return adm.Perform(`postbody`, `/predicate/`, `command`, req, c)
}

// predicateRemove function
// soma predicate remove ${pred}
func predicateRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	path := fmt.Sprintf(
		"/predicate/%s",
		url.QueryEscape(c.Args().First()),
	)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// predicateList function
// soma predicate list
func predicateList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/predicate/`, `list`, nil, c)
}

// predicateShow function
// soma predicate show ${pred}
func predicateShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf(
		"/predicate/%s",
		url.QueryEscape(c.Args().First()),
	)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
