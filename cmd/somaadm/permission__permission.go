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
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/internal/help"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerPermissions(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  `permission`,
				Usage: `SUBCOMMANDS for permissions`,
				Subcommands: []cli.Command{
					{
						Name:         `add`,
						Usage:        `Register a new permission`,
						Description:  help.Text(`permission::add`),
						Action:       runtime(cmdPermissionAdd),
						BashComplete: cmpl.To,
					},
					{
						Name:         `remove`,
						Usage:        `Remove a permission from a category`,
						Description:  help.Text(`permission::remove`),
						Action:       runtime(cmdPermissionRemove),
						BashComplete: cmpl.From,
					},
					{
						Name:         `list`,
						Usage:        `List all permissions in a category`,
						Description:  help.Text(`permission::list`),
						Action:       runtime(cmdPermissionList),
						BashComplete: cmpl.DirectIn,
					},
					{
						Name:         `show`,
						Usage:        `Show details for a permission`,
						Description:  help.Text(`permission::show`),
						Action:       runtime(cmdPermissionShow),
						BashComplete: cmpl.In,
					},
					{
						Name:         `map`,
						Usage:        `Map an action to a permission`,
						Description:  help.Text(`permission::map`),
						Action:       runtime(cmdPermissionMap),
						BashComplete: cmpl.To,
					},
					{
						Name:         `unmap`,
						Usage:        `Unmap an action from a permission`,
						Description:  help.Text(`permission::unmap`),
						Action:       runtime(cmdPermissionUnmap),
						BashComplete: cmpl.From,
					},
				},
			}, // end permissions
		}...,
	)
	return &app
}

// cmdPermissionAdd function
// somaadm permission add ${permission} to ${category}
func cmdPermissionAdd(c *cli.Context) error {
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

	if err := adm.ValidateNoColon(c.Args().First()); err != nil {
		return err
	}
	if err := adm.ValidateCategory(opts[`to`][0]); err != nil {
		return err
	}

	esc := url.QueryEscape(opts[`to`][0])
	req := proto.NewPermissionRequest()
	req.Permission.Name = c.Args().First()
	req.Permission.Category = opts[`to`][0]
	path := fmt.Sprintf("/category/%s/permission/", esc)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

// cmdPermissionRemove function
// somaadm permission remove ${permission} from ${category}
// somaadm permission remove ${category}::${permission}
func cmdPermissionRemove(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`from`}
	mandatoryOptions := []string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	var permission, category, permissionID string
	permissionSlice := strings.Split(c.Args().First(), `::`)
	switch len(permissionSlice) {
	case 1:
		permission = permissionSlice[0]
	case 2:
		permission = permissionSlice[0]
		category = permissionSlice[1]
		if err := adm.ValidateCategory(category); err != nil {
			return err
		}
	}
	if _, ok := opts[`from`]; !ok {
		// if the optional argument was not provided, then category must
		// have been set via splitting permissionSlice
		if category == `` {
			return fmt.Errorf(`Missing category information`)
		}
	}
	if category == `` {
		category = opts[`from`][0]
		if err := adm.ValidateCategory(category); err != nil {
			return err
		}
	} else {
		if category != opts[`from`][0] {
			// example: somaadm permission remove self::information from global
			return fmt.Errorf("Mismatching category information: %s vs %s",
				category,
				opts[`from`][0],
			)
		}
	}

	if err := adm.LookupPermIDRef(permission, category,
		&permissionID,
	); err != nil {
		return err
	}

	path := fmt.Sprintf("/category/%s/permission/%s",
		url.QueryEscape(category),
		url.QueryEscape(permissionID))
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// cmdPermissionList function
// somaadm permission list in ${category}
func cmdPermissionList(c *cli.Context) error {
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

	esc := url.QueryEscape(opts[`in`][0])
	path := fmt.Sprintf("/category/%s/permission/", esc)
	return adm.Perform(`get`, path, `list`, nil, c)
}

// cmdPermissionShow function
// somaadm permission show ${permission} in ${category}
// somaadm permission show ${category}::${permission}
func cmdPermissionShow(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`in`}
	mandatoryOptions := []string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	var permission, category, permissionID string
	permissionSlice := strings.Split(c.Args().First(), `::`)
	switch len(permissionSlice) {
	case 1:
		permission = permissionSlice[0]
	case 2:
		permission = permissionSlice[0]
		category = permissionSlice[1]
		if err := adm.ValidateCategory(category); err != nil {
			return err
		}
	}
	if _, ok := opts[`in`]; !ok {
		// if the optional argument was not provided, then category must
		// have been set via splitting permissionSlice
		if category == `` {
			return fmt.Errorf(`Missing category information`)
		}
	}
	if category == `` {
		category = opts[`in`][0]
		if err := adm.ValidateCategory(category); err != nil {
			return err
		}
	} else {
		if category != opts[`in`][0] {
			return fmt.Errorf("Mismatching category information: %s vs %s",
				category,
				opts[`in`][0],
			)
		}
	}

	if err := adm.LookupPermIDRef(permission, category,
		&permissionID,
	); err != nil {
		return err
	}

	path := fmt.Sprintf("/category/%s/permission/%s",
		url.QueryEscape(category),
		url.QueryEscape(permissionID))
	return adm.Perform(`get`, path, `show`, nil, c)
}

// cmdPermissionMap function
// somaadm permission map ${section}::${action} to ${category}::${permission}
// somaadm permission map ${section} to ${category}::${permission}
func cmdPermissionMap(c *cli.Context) error {
	return cmdPermissionEdit(c, `map`)
}

// cmdPermissionUnmap function
// somaadm permission unmap ${section}::${action} from ${category}::${permission}
// somaadm permission unmap ${section} from ${category}::${permission}
func cmdPermissionUnmap(c *cli.Context) error {
	return cmdPermissionEdit(c, `unmap`)
}

// cmdPermissionEdit function
func cmdPermissionEdit(c *cli.Context, cmd string) error {
	var (
		err                                              error
		section, action, category, permission, sCategory string
		sectionID, actionID, permissionID, syn           string
		sectionMapping                                   bool
	)
	switch cmd {
	case `map`:
		syn = `to`
	case `unmap`:
		syn = `from`
	}
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{syn}
	mandatoryOptions := []string{syn}

	if err = adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	actionSlice := strings.Split(c.Args().First(), `::`)
	switch len(actionSlice) {
	case 1:
		section = actionSlice[0]
		sectionMapping = true
	case 2:
		section = actionSlice[0]
		action = actionSlice[1]
	default:
		return fmt.Errorf("Not a valid {section}::{action}"+
			" specifier: %s", c.Args().First())
	}

	permissionSlice := strings.Split(opts[syn][0], `::`)
	switch len(permissionSlice) {
	case 2:
		category = permissionSlice[0]
		permission = permissionSlice[1]
	default:
		return fmt.Errorf("Not a valid {category}::{permission}"+
			" specifier: %s", opts[syn][0])
	}
	// validate category
	if err = adm.ValidateCategory(category); err != nil {
		return err
	}
	// lookup permissionID
	if err = adm.LookupPermIDRef(
		permission,
		category,
		&permissionID,
	); err != nil {
		return err
	}
	// lookup sectionID
	if sectionID, err = adm.LookupSectionID(
		section,
	); err != nil {
		return err
	}
	// lookup ${category} for ${sectionID}
	if sCategory, err = adm.LookupCategoryBySection(
		sectionID,
	); err != nil {
		return err
	}
	// mapped section's category must match permission's category
	if sCategory != category {
		return fmt.Errorf("Category mismatch. Section %s from category %s can not be mapped to permission %s in category %s",
			section, sCategory, permission, category)
	}

	// lookup actionID
	if !sectionMapping {
		if actionID, err = adm.LookupActionID(
			action,
			sectionID,
		); err != nil {
			return err
		}
	}

	req := proto.NewPermissionRequest()
	switch cmd {
	case `map`:
		req.Flags.Add = true
	case `unmap`:
		req.Flags.Remove = true
	}
	req.Permission.ID = permissionID
	req.Permission.Name = permission
	req.Permission.Category = category
	if !sectionMapping {
		req.Permission.Actions = &[]proto.Action{
			proto.Action{
				ID:        actionID,
				Name:      action,
				SectionID: sectionID,
				Category:  category,
			},
		}
	} else {
		req.Permission.Sections = &[]proto.Section{
			proto.Section{
				ID:       sectionID,
				Name:     section,
				Category: category,
			},
		}
	}

	esc := url.QueryEscape(category)
	path := fmt.Sprintf("/category/%s/permission/%s",
		esc, permissionID)
	return adm.Perform(`patchbody`, path, `command`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
