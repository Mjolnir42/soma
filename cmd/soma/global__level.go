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

func registerLevels(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:        `level`,
				Usage:       `SUBCOMMANDS for notification level definitions`,
				Description: help.Text(`level::`),
				Subcommands: []cli.Command{
					{
						Name:         `add`,
						Usage:        `Add a new notification level`,
						Description:  help.Text(`level::add`),
						Action:       runtime(levelAdd),
						BashComplete: cmpl.LevelAdd,
					},
					{
						Name:        `remove`,
						Usage:       `Remove a notification level`,
						Description: help.Text(`level::remove`),
						Action:      runtime(levelRemove),
					},
					{
						Name:        `list`,
						Usage:       `List all notification levels`,
						Description: help.Text(`level::list`),
						Action:      runtime(levelList),
					},
					{
						Name:        `show`,
						Usage:       `Show details about a notification level`,
						Description: help.Text(`level::show`),
						Action:      runtime(levelShow),
					},
					{
						Name:        `search`,
						Usage:       `Lookup a notification level by name or shortname`,
						Description: help.Text(`level::search`),
						Action:      runtime(levelSearch),
					},
				},
			},
		}...,
	)
	return &app
}

// levelAdd function
// soma level add ${lvl} shortname ${abbrev} numeric ${num}
func levelAdd(c *cli.Context) error {
	var lvl uint64
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`shortname`, `numeric`}
	mandatoryOptions := []string{`shortname`, `numeric`}

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

	if err := adm.ValidateLBoundUint64(
		opts[`numeric`][0], &lvl, 0,
	); err != nil {
		return err
	}

	req := proto.NewLevelRequest()
	req.Level.Name = c.Args().First()
	req.Level.ShortName = opts[`shortname`][0]
	req.Level.Numeric = uint16(lvl)

	return adm.Perform(`postbody`, `/level/`, `command`, req, c)
}

// levelRemove function
// soma level remove ${lvl}|${abbrev}
func levelRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	name := ``
	if err := adm.LookupLevelName(c.Args().First(), &name); err != nil {
		return err
	}

	path := fmt.Sprintf("/level/%s", url.QueryEscape(name))
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// levelList function
// soma level list
func levelList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/level/`, `list`, nil, c)
}

// levelShow function
// soma level show ${lvl}|${abbrev}
func levelShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	name := ``
	if err := adm.LookupLevelName(c.Args().First(), &name); err != nil {
		return err
	}

	path := fmt.Sprintf("/level/%s", url.QueryEscape(name))
	return adm.Perform(`get`, path, `show`, nil, c)
}

// levelSearch function
// soma level search ${lvl}|${abbrev}
func levelSearch(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	req := proto.NewLevelFilter()
	req.Filter.Level.Name = c.Args().First()
	req.Filter.Level.ShortName = c.Args().First()

	return adm.Perform(`postbody`, `/search/level/`, `list`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
