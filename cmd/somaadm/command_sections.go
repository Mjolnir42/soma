/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/internal/help"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerSection(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  `sections`,
				Usage: `SUBCOMMANDS for permission sections`,
				Subcommands: []cli.Command{
					{
						Name:         `add`,
						Usage:        `Add a permission section`,
						Description:  help.Text(`SectionsAdd`),
						Action:       runtime(cmdSectionAdd),
						BashComplete: cmpl.To,
					},
					{
						Name:        `remove`,
						Usage:       `Remove a permission section`,
						Description: help.Text(`SectionsRemove`),
						Action:      runtime(cmdSectionRemove),
					},
					{
						Name:        `list`,
						Usage:       `List permission sections`,
						Description: help.Text(`SectionsList`),
						Action:      runtime(cmdSectionList),
					},
					{
						Name:        `show`,
						Usage:       `Show details about permission section`,
						Description: help.Text(`SectionsShow`),
						Action:      runtime(cmdSectionShow),
					},
				},
			},
		}...,
	)
	return &app
}

func cmdSectionAdd(c *cli.Context) error {
	unique := []string{`to`}
	required := []string{`to`}
	opts := make(map[string][]string)
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		unique,
		required,
		c.Args().Tail(),
	); err != nil {
		return err
	}
	if err := adm.ValidateNoColon(c.Args().First()); err != nil {
		return err
	}

	if err := adm.ValidateCategory(opts[`to`][0]); err != nil {
		return err
	}

	req := proto.NewSectionRequest()
	req.Section.Name = c.Args().First()
	req.Section.Category = opts[`to`][0]
	path := fmt.Sprintf("/category/%s/section/", req.Section.Category)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

func cmdSectionRemove(c *cli.Context) error {
	var (
		err                 error
		category, sectionID string
	)
	if err = adm.VerifySingleArgument(c); err != nil {
		return err
	}
	if sectionID, err = adm.LookupSectionID(
		c.Args().First()); err != nil {
		return err
	}
	if category, err = adm.LookupCategoryBySection(
		sectionID); err != nil {
		return err
	}

	path := fmt.Sprintf("/category/%s/section/%s", category, sectionID)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdSectionList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/section/`, `list`, nil, c)
}

func cmdSectionShow(c *cli.Context) error {
	var (
		err       error
		sectionID string
	)
	if err = adm.VerifySingleArgument(c); err != nil {
		return err
	}
	if sectionID, err = adm.LookupSectionID(
		c.Args().First()); err != nil {
		return err
	}

	path := fmt.Sprintf("/section/%s", sectionID)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
