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

func registerAction(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  `action`,
				Usage: `SUBCOMMANDS for permission actions`,
				Subcommands: []cli.Command{
					{
						Name:         `add`,
						Usage:        `Add a permission action to a section`,
						Description:  help.Text(`action::add`),
						Action:       runtime(cmdActionAdd),
						BashComplete: cmpl.InTo,
					},
					{
						Name:         `remove`,
						Usage:        `Remove a permission action from a section`,
						Description:  help.Text(`action::remove`),
						Action:       runtime(cmdActionRemove),
						BashComplete: cmpl.InFrom,
					},
					{
						Name:         `list`,
						Usage:        `List permission actions in a section`,
						Description:  help.Text(`action::list`),
						Action:       runtime(cmdActionList),
						BashComplete: cmpl.DirectInOf,
					},
					{
						Name:         `show`,
						Usage:        `Show details about a permission action`,
						Description:  help.Text(`action::show`),
						Action:       runtime(cmdActionShow),
						BashComplete: cmpl.InFrom,
					},
				},
			},
		}...,
	)
	return &app
}

// cmdActionAdd function
// somaadm action add ${action} to ${section} [in ${category}]
func cmdActionAdd(c *cli.Context) error {
	var (
		err                 error
		sectionID, category string
	)
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`to`, `in`}
	mandatoryOptions := []string{`to`}

	if err = adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	if err = adm.ValidateNoColon(c.Args().First()); err != nil {
		return err
	}

	if sectionID, err = adm.LookupSectionID(
		opts[`to`][0],
	); err != nil {
		return err
	}

	// lookup ${category} for ${sectionID}
	if category, err = adm.LookupCategoryBySection(
		sectionID,
	); err != nil {
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

	req := proto.NewActionRequest()
	req.Action.Name = c.Args().First()
	req.Action.SectionID = sectionID
	path := fmt.Sprintf("/category/%s/section/%s/action/",
		category, sectionID)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

// cmdActionRemove function
// somaadm action remove ${action} from ${section} [in ${category}]
func cmdActionRemove(c *cli.Context) error {
	var (
		err                           error
		category, sectionID, actionID string
	)
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`from`, `in`}
	mandatoryOptions := []string{`from`}

	if err = adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	if sectionID, err = adm.LookupSectionID(
		opts[`from`][0],
	); err != nil {
		return err
	}
	if actionID, err = adm.LookupActionID(
		c.Args().First(),
		sectionID,
	); err != nil {
		return err
	}

	// lookup ${category} for ${sectionID}
	if category, err = adm.LookupCategoryBySection(
		sectionID,
	); err != nil {
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

	path := fmt.Sprintf("/category/%s/section/%s/action/%s",
		category, sectionID, actionID)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// cmdActionList function
// somaadm action list in ${section} [of  ${category}]
func cmdActionList(c *cli.Context) error {
	var (
		err                 error
		category, sectionID string
	)
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`in`, `of`}
	mandatoryOptions := []string{`in`}

	if err = adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	if sectionID, err = adm.LookupSectionID(
		opts[`in`][0],
	); err != nil {
		return err
	}

	// lookup ${category} for ${sectionID}
	if category, err = adm.LookupCategoryBySection(
		sectionID,
	); err != nil {
		return err
	}

	// if [of  ${category}] was given on the command line, the given
	// category must be valid and must be the correct one for ${name}
	if _, ok := opts[`of`]; ok {
		if err := adm.ValidateCategory(opts[`of`][0]); err != nil {
			return err
		}
		if opts[`of`][0] != category {
			return fmt.Errorf("Category mismatch: %s vs %s",
				opts[`of`][0],
				category,
			)
		}
	}

	path := fmt.Sprintf("/category/%s/section/%s/action/",
		category, sectionID)
	return adm.Perform(`get`, path, `list`, nil, c)
}

// cmdActionShow function
// somaadm action show ${action} from ${section} [in ${category}]
func cmdActionShow(c *cli.Context) error {
	var (
		err                           error
		category, sectionID, actionID string
	)
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`from`, `in`}
	mandatoryOptions := []string{`from`}

	if err = adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	if sectionID, err = adm.LookupSectionID(
		opts[`from`][0],
	); err != nil {
		return err
	}
	if actionID, err = adm.LookupActionID(
		c.Args().First(),
		sectionID,
	); err != nil {
		return err
	}
	// lookup ${category} for ${sectionID}
	if category, err = adm.LookupCategoryBySection(
		sectionID,
	); err != nil {
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

	path := fmt.Sprintf("/category/%s/section/%s/action/%s",
		category, sectionID, actionID)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
