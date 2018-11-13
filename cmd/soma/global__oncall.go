/*-
 * Copyright (c) 2015-2018, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2018, 1&1 IONOS SE
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

func registerOncall(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:        `oncall`,
				Usage:       `SUBCOMMANDS for oncall duty team management`,
				Description: help.Text(`oncall::`),
				Subcommands: []cli.Command{
					{
						Name:         `add`,
						Usage:        `Register a new oncall duty team`,
						Description:  help.Text(`oncall::add`),
						Action:       runtime(oncallAdd),
						BashComplete: cmpl.OncallAdd,
					},
					{
						Name:        `remove`,
						Usage:       `Remove an existing oncall duty team`,
						Description: help.Text(`oncall::remove`),
						Action:      runtime(oncallRemove),
					},
					{
						Name:         `update`,
						Usage:        `Update phone number or name of an existing oncall duty team`,
						Description:  help.Text(`oncall::update`),
						Action:       runtime(oncallUpdate),
						BashComplete: cmpl.OncallUpdate,
					},
					{
						Name:        `list`,
						Usage:       `List all registered oncall duty teams`,
						Description: help.Text(`oncall::list`),
						Action:      runtime(oncallList),
					},
					{
						Name:        `show`,
						Usage:       `Show information about a specific oncall duty team`,
						Description: help.Text(`oncall::show`),
						Action:      runtime(oncallShow),
					},
					{
						Name:        `member`,
						Usage:       `SUBCOMMANDS to manipulate oncall duty membership`,
						Description: help.Text(`oncall::`),
						Subcommands: []cli.Command{
							{
								Name:         `assign`,
								Usage:        `Assign a user to an oncall duty team`,
								Description:  help.Text(`oncall::member-assign`),
								Action:       runtime(oncallMemberAssign),
								BashComplete: cmpl.To,
							},
							{
								Name:         `unassign`,
								Usage:        `Unassign a user from an oncall duty team`,
								Description:  help.Text(`oncall::member-unassign`),
								Action:       runtime(oncallMemberUnassign),
								BashComplete: cmpl.From,
							},
							{
								Name:        `list`,
								Usage:       `List the users assigned to an oncall duty team`,
								Description: help.Text(`oncall::member-list`),
								Action:      runtime(oncallMemberList),
							},
						},
					},
				},
			},
		}...,
	)
	return &app
}

// oncallAdd function
// soma oncall add ${name} phone ${extension}
func oncallAdd(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`phone`}
	mandatoryOptions := []string{`phone`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	if err := adm.ValidateOncallNumber(opts[`phone`][0]); err != nil {
		return err
	}
	if err := adm.ValidateNotUUID(c.Args().First()); err != nil {
		return err
	}

	req := proto.NewOncallRequest()
	req.Oncall.Name = c.Args().First()
	req.Oncall.Number = opts[`phone`][0]

	return adm.Perform(`postbody`, `/oncall/`, `command`, req, c)
}

// oncallRemove function
// soma oncall remove ${name}
func oncallRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	oncallID, err := adm.LookupOncallID(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/oncall/%s", url.QueryEscape(oncallID))
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// oncallUpdate function
// soma oncall update ${name} [phone ${extension}] [name ${new-name}]
func oncallUpdate(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`phone`, `name`}
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

	oncallID, err := adm.LookupOncallID(c.Args().First())
	if err != nil {
		return err
	}

	req := proto.NewOncallRequest()
	validUpdate := false
	// both arguments are optional, but one of them must be given
	if _, ok := opts[`phone`]; ok {
		if err := adm.ValidateOncallNumber(opts[`phone`][0]); err != nil {
			return err
		}
		req.Oncall.Number = opts[`phone`][0]
		validUpdate = true
	}
	if _, ok := opts[`name`]; ok {
		if err := adm.ValidateNotUUID(opts[`name`][0]); err != nil {
			return err
		}
		req.Oncall.Name = opts[`name`][0]
		validUpdate = true
	}

	if !validUpdate {
		return fmt.Errorf("Syntax error: specify either name or phone" +
			" extension for update")
	}

	path := fmt.Sprintf("/oncall/%s", url.QueryEscape(oncallID))
	return adm.Perform(`patchbody`, path, `command`, req, c)
}

// oncallList function
// soma oncall list
func oncallList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/oncall/`, `list`, nil, c)
}

// oncallShow function
// soma oncall show ${name}
func oncallShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	oncallID, err := adm.LookupOncallID(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/oncall/%s", url.QueryEscape(oncallID))
	return adm.Perform(`get`, path, `show`, nil, c)
}

// oncallMemberAssign function
// soma oncall member assign ${user} to ${oncall}
func oncallMemberAssign(c *cli.Context) error {
	var err error
	var userID, oncallID string

	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`to`}
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

	if userID, err = adm.LookupUserID(c.Args().First()); err != nil {
		return err
	}
	if oncallID, err = adm.LookupOncallID(opts[`to`][0]); err != nil {
		return err
	}

	req := proto.NewOncallRequest()
	req.Oncall.ID = oncallID
	req.Oncall.Members = &[]proto.OncallMember{
		proto.OncallMember{UserID: userID},
	}

	path := fmt.Sprintf("/oncall/%s/member/", url.QueryEscape(oncallID))
	return adm.Perform(`patchbody`, path, `command`, req, c)
}

// oncallMemberUnassign function
// soma oncall member unassign ${user} from ${oncall}
func oncallMemberUnassign(c *cli.Context) error {
	var err error
	var userID, oncallID string

	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`from`}
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

	if userID, err = adm.LookupUserID(c.Args().First()); err != nil {
		return err
	}
	if oncallID, err = adm.LookupOncallID(opts[`from`][0]); err != nil {
		return err
	}

	path := fmt.Sprintf("/oncall/%s/member/%s",
		url.QueryEscape(oncallID),
		url.QueryEscape(userID),
	)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// oncallMemberList function
// soma oncall member list ${oncall}
func oncallMemberList(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	oncallID, err := adm.LookupOncallID(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/oncall/%s/member/", url.QueryEscape(oncallID))
	return adm.Perform(`get`, path, `list`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
