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
	"os"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/internal/help"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerProperty(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:        `property-mgmt`,
				Usage:       `SUBCOMMANDS for property management`,
				Description: help.Text(`property-mgmt::`),
				Subcommands: []cli.Command{
					{
						Name:        `system`,
						Usage:       `SUBCOMMANDS for global system property management`,
						Description: help.Text(`property-system::`),
						Subcommands: []cli.Command{
							{
								Name:        `add`,
								Usage:       `Add a new global system property`,
								Description: help.Text(`property-system::add`),
								Action:      runtime(propertyMgmtSystemAdd),
							},
							{
								Name:        `remove`,
								Usage:       `Remove a global system property`,
								Description: help.Text(`property-system::remove`),
								Action:      runtime(propertyMgmtSystemRemove),
							},
							{
								Name:        `show`,
								Usage:       `Show details about a global system property`,
								Description: help.Text(`property-system::show`),
								Action:      runtime(propertyMgmtSystemShow),
							},
							{
								Name:        `list`,
								Usage:       `List all global system properties`,
								Description: help.Text(`property-system::list`),
								Action:      runtime(propertyMgmtSystemList),
							},
						},
					},
				},
			},
		}...,
	)
	return &app
}

// SYSTEM PROPERTIES

// propertyMgmtSystemAdd function
// soma property-mgmt system add ${property}
func propertyMgmtSystemAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	req := proto.NewSystemPropertyRequest()
	req.Property.System.Name = c.Args().First()

	path := fmt.Sprintf("/property-mgmt/%s/", proto.PropertyTypeSystem)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

// propertyMgmtSystemRemove function
// soma property-mgmt system remove ${property}
func propertyMgmtSystemRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/property-mgmt/%s/%s",
		proto.PropertyTypeSystem,
		esc,
	)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// propertyMgmtSystemShow function
// soma property-mgmt system show ${property}
func propertyMgmtSystemShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/property-mgmt/%s/%s",
		proto.PropertyTypeSystem,
		esc,
	)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// propertyMgmtSystemList function
// soma property-mgmt system list
func propertyMgmtSystemList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/property-mgmt/%s/", proto.PropertyTypeSystem)
	return adm.Perform(`get`, path, `list`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
