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
	"net/url"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/help"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerCategories(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:        `category`,
				Usage:       `SUBCOMMANDS for permission scope categories`,
				Description: help.Text(`category`),
				Subcommands: []cli.Command{
					{
						Name:        `add`,
						Usage:       `Register a new permission scope category`,
						Description: help.Text(`category::add`),
						Action:      runtime(cmdPermissionCategoryAdd),
					},
					{
						Name:        `remove`,
						Usage:       `Remove an existing permission scope category`,
						Description: help.Text(`category::remove`),
						Action:      runtime(cmdPermissionCategoryRemove),
					},
					{
						Name:        `list`,
						Usage:       `List all permission scope categories`,
						Description: help.Text(`category::list`),
						Action:      runtime(cmdPermissionCategoryList),
					},
					{
						Name:        `show`,
						Usage:       `Show details for a permission scope category`,
						Description: help.Text(`category::show`),
						Action:      runtime(cmdPermissionCategoryShow),
					},
				},
			},
		}...,
	)
	return &app
}

// cmdPermissionCategoryAdd function -- somaadm category add ${name}
func cmdPermissionCategoryAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	if err := adm.ValidateNoColon(c.Args().First()); err != nil {
		return err
	}
	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	req := proto.NewCategoryRequest()
	req.Category.Name = c.Args().First()

	return adm.Perform(`postbody`, `/category/`, `command`, req, c)
}

// cmdPermissionCategoryRemove function -- somaadm category remove ${name}
func cmdPermissionCategoryRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/category/%s", esc)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// cmdPermissionCategoryList function -- somaadm category list
func cmdPermissionCategoryList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/category/`, `list`, nil, c)
}

// cmdPermissionCategoryShow function -- somaadm category show ${name}
func cmdPermissionCategoryShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/category/%s", esc)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
