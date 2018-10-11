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

func registerProviders(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// providers
			{
				Name:        `provider`,
				Usage:       `SUBCOMMANDS for metric providers`,
				Description: help.Text(`provider::`),
				Subcommands: []cli.Command{
					{
						Name:        `add`,
						Usage:       `Add a new metric provider`,
						Description: help.Text(`provider::add`),
						Action:      runtime(cmdProviderAdd),
					},
					{
						Name:        `remove`,
						Usage:       `Remove a metric provider`,
						Description: help.Text(`provider::remove`),
						Action:      runtime(cmdProviderRemove),
					},
					{
						Name:        `list`,
						Usage:       `List metric providers`,
						Description: help.Text(`provider::list`),
						Action:      runtime(cmdProviderList),
					},
					{
						Name:        `show`,
						Usage:       `Show details about a metric provider`,
						Description: help.Text(`provider::show`),
						Action:      runtime(cmdProviderShow),
					},
				},
			}, // end providers
		}...,
	)
	return &app
}

// cmdProviderAdd function
// soma provider add ${provider}
func cmdProviderAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.NewProviderRequest()
	req.Provider.Name = c.Args().First()

	if err := adm.ValidateNoSlash(req.Provider.Name); err != nil {
		return err
	}

	return adm.Perform(`postbody`, `/provider/`, `command`, req, c)
}

// cmdProviderRemove function
// soma provider remove ${provider}
func cmdProviderRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/provider/%s", esc)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// cmdProviderList function
// soma provider list
func cmdProviderList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/provider/`, `list`, nil, c)
}

// cmdProviderShow function
// soma provider show ${provider}
func cmdProviderShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/provider/%s", esc)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
