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

func registerEntities(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:        `entity`,
				Usage:       `SUBCOMMANDS for entities`,
				Description: help.Text(`entity`),
				Subcommands: []cli.Command{
					{
						Name:        `add`,
						Usage:       `Add a new entity`,
						Description: help.Text(`entity::add`),
						Action:      runtime(cmdEntityAdd),
					},
					{
						Name:        `remove`,
						Usage:       `Remove an existing entity`,
						Description: help.Text(`entity::remove`),
						Action:      runtime(cmdEntityRemove),
					},
					{
						Name:         `rename`,
						Usage:        `Rename an existing entity`,
						Description:  help.Text(`entity::rename`),
						Action:       runtime(cmdEntityRename),
						BashComplete: cmpl.To,
					},
					{
						Name:        `list`,
						Usage:       `List all entities`,
						Description: help.Text(`entity::list`),
						Action:      runtime(cmdEntityList),
					},
					{
						Name:        `show`,
						Usage:       `Show information about a specific entity`,
						Description: help.Text(`entity::show`),
						Action:      runtime(cmdEntityShow),
					},
				},
			},
		}...,
	)
	return &app
}

// cmdEntityAdd  function
// somaadm entity add ${entity}
func cmdEntityAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.NewEntityRequest()
	req.Entity.Name = c.Args().First()

	return adm.Perform(`postbody`, `/entity/`, `command`, req, c)
}

// cmdEntityRemove function
// somaadm entity remove ${entity}
func cmdEntityRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/entity/%s", url.QueryEscape(c.Args().First()))
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// cmdEntityRename function
// somaadm entity rename ${entity} to ${new-entity}
func cmdEntityRename(c *cli.Context) error {
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

	req := proto.NewEntityRequest()
	req.Entity.Name = opts[`to`][0]

	path := fmt.Sprintf("/entity/%s", url.QueryEscape(c.Args().First()))
	return adm.Perform(`putbody`, path, `command`, req, c)
}

// cmdEntityList function
// somaadm entity list
func cmdEntityList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/entity/`, `list`, nil, c)
}

// cmdEntityShow function
// somaadm entity show ${entity}
func cmdEntityShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/entity/%s", url.QueryEscape(c.Args().First()))
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
