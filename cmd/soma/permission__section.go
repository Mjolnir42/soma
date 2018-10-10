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
				Name:        `section`,
				Usage:       `SUBCOMMANDS for permission sections`,
				Description: help.Text(`section::`),
				Subcommands: []cli.Command{
					{
						Name:         `add`,
						Usage:        `Add a permission section`,
						Description:  help.Text(`section::add`),
						Action:       runtime(cmdSectionAdd),
						BashComplete: cmpl.To,
					},
					{
						Name:         `remove`,
						Usage:        `Remove a permission section`,
						Description:  help.Text(`section::remove`),
						Action:       runtime(cmdSectionRemove),
						BashComplete: cmpl.From,
					},
					{
						Name:         `list`,
						Usage:        `List permission sections`,
						Description:  help.Text(`section::list`),
						Action:       runtime(cmdSectionList),
						BashComplete: cmpl.DirectIn,
					},
					{
						Name:         `show`,
						Usage:        `Show details about permission section`,
						Description:  help.Text(`section::show`),
						Action:       runtime(cmdSectionShow),
						BashComplete: cmpl.In,
					},
				},
			},
		}...,
	)
	return &app
}

// cmdSectionAdd function -- somaadm section add ${name} to ${category}
func cmdSectionAdd(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`to`}
	mandatoryOptions := []string{`to`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail()); err != nil {
		return err
	}

	if err := adm.ValidateNoColon(c.Args().First()); err != nil {
		return err
	}

	if err := adm.ValidateCategory(opts[`to`][0]); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(opts[`to`][0]); err != nil {
		return err
	}

	req := proto.NewSectionRequest()
	req.Section.Name = c.Args().First()
	req.Section.Category = opts[`to`][0]
	path := fmt.Sprintf("/category/%s/section/", req.Section.Category)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

// cmdSectionRemove -- somaadm section remove ${name} [from ${category}]
func cmdSectionRemove(c *cli.Context) error {
	var (
		err                 error
		category, sectionID string
	)

	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`from`}
	mandatoryOptions := []string{}
	if err = adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail()); err != nil {
		return err
	}

	// lookup ${sectionID} for ${name}
	if sectionID, err = adm.LookupSectionID(
		c.Args().First()); err != nil {
		return err
	}

	// lookup ${category} for ${sectionID}
	if category, err = adm.LookupCategoryBySection(sectionID); err != nil {
		return err
	}

	// if [from  ${category}] was given on the command line, the given
	// category must be valid and must be the correct one for ${name}
	if _, ok := opts[`from`]; ok {
		if err := adm.ValidateCategory(opts[`from`][0]); err != nil {
			return err
		}
		if opts[`from`][0] != category {
			return fmt.Errorf("Category mismatch: %s vs %s",
				opts[`from`][0],
				category,
			)
		}
	}

	if err := adm.ValidateNoSlash(category); err != nil {
		return err
	}

	path := fmt.Sprintf("/category/%s/section/%s", category, sectionID)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// cmdSectionList function -- somaadm section list in ${category}
func cmdSectionList(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`in`}
	mandatoryOptions := []string{`in`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		adm.AllArguments(c),
	); err != nil {
		return err
	}

	if err := adm.ValidateCategory(opts[`in`][0]); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(opts[`in`][0]); err != nil {
		return err
	}

	path := fmt.Sprintf("/category/%s/section/", opts[`in`][0])
	return adm.Perform(`get`, path, `list`, nil, c)
}

// cmdSectionShow function -- somaadm section show ${name} [in ${category}]
func cmdSectionShow(c *cli.Context) error {
	var (
		err                 error
		sectionID, category string
	)
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`in`}
	mandatoryOptions := []string{}

	if err = adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail()); err != nil {
		return err
	}

	// lookup ${sectionID} for ${name}
	if sectionID, err = adm.LookupSectionID(
		c.Args().First()); err != nil {
		return err
	}

	// lookup ${category} for ${sectionID}
	if category, err = adm.LookupCategoryBySection(sectionID); err != nil {
		return err
	}

	// if [in  ${category}] was given on the command line, the given
	// category must be valid and must be the correct one for ${name}
	if _, ok := opts[`in`]; ok {
		if err := adm.ValidateCategory(opts[`in`][0]); err != nil {
			return err
		}
		if opts[`in`][0] != category {
			return fmt.Errorf("Category mismatch: %s vs %s",
				opts[`in`][0],
				category,
			)
		}
	}

	if err := adm.ValidateNoSlash(category); err != nil {
		return err
	}

	path := fmt.Sprintf("/category/%s/section/%s", category, sectionID)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
