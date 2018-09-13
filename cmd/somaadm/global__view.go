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
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/internal/help"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerViews(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  `view`,
				Usage: `SUBCOMMANDS for views`,
				Subcommands: []cli.Command{
					{
						Name:        `add`,
						Usage:       `Register a new view`,
						Description: help.Text(`view::add`),
						Action:      runtime(cmdViewAdd),
					},
					{
						Name:        `remove`,
						Usage:       `Remove an existing view`,
						Description: help.Text(`view::remove`),
						Action:      runtime(cmdViewRemove),
					},
					{
						Name:         `rename`,
						Usage:        `Rename an existing view`,
						Description:  help.Text(`view::rename`),
						Action:       runtime(cmdViewRename),
						BashComplete: cmpl.To,
					},
					{
						Name:        `list`,
						Usage:       `List all registered views`,
						Description: help.Text(`view::list`),
						Action:      runtime(cmdViewList),
					},
					{
						Name:        `show`,
						Usage:       `Show information about a specific view`,
						Description: help.Text(`view::show`),
						Action:      runtime(cmdViewShow),
					},
				},
			}, // end views
		}...,
	)
	return &app
}

// cmdViewAdd function
// somaadm view add ${view}
func cmdViewAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.Request{}
	req.View = &proto.View{}
	req.View.Name = c.Args().First()
	if strings.Contains(req.View.Name, `.`) {
		return fmt.Errorf(`Views must not contain the character '.'`)
	}

	return adm.Perform(`postbody`, `/view/`, `command`, req, c)
}

// cmdViewRemove function
// somaadm view remove ${view}
func cmdViewRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/view/%s", url.QueryEscape(c.Args().First()))
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// cmdViewRename function
// somaadm view rename ${view} to ${new-view}
func cmdViewRename(c *cli.Context) error {
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

	req := proto.NewViewRequest()
	req.View.Name = opts[`to`][0]

	path := fmt.Sprintf("/view/%s", url.QueryEscape(c.Args().First()))
	return adm.Perform(`putbody`, path, `command`, nil, c)
}

// cmdViewList function
// somaadm view list
func cmdViewList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/view/`, `list`, nil, c)
}

// cmdViewShow function
// somaadm view show ${view}
func cmdViewShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/view/%s", url.QueryEscape(c.Args().First()))
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
