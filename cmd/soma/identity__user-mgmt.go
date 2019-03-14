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

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/internal/help"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerUserMgmt(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:        `user-mgmt`,
				Usage:       `SUBCOMMANDS for account management`,
				Description: help.Text(`user-mgmt::`),
				Subcommands: []cli.Command{
					{
						Name:         `add`,
						Usage:        `Add a new user account`,
						Description:  help.Text(`user-mgmt::add`),
						Action:       runtime(userMgmtAdd),
						BashComplete: cmpl.UserMgmtAdd,
					},
					{
						Name:         `update`,
						Usage:        `Set/change user information`,
						Description:  help.Text(`user-mgmt::update`),
						Action:       runtime(userMgmtUpdate),
						BashComplete: cmpl.UserMgmtUpdate,
					},
					{
						Name:        `remove`,
						Usage:       `Flag a user account as deleted`,
						Description: help.Text(`user-mgmt::remove`),
						Action:      runtime(userMgmtRemove),
					},
					{
						Name:        `purge`,
						Usage:       `Purge a removed user from the system`,
						Description: help.Text(`user-mgmt::purge`),
						Action:      runtime(userMgmtPurge),
					},
					{
						Name:        `list`,
						Usage:       `List all registered users`,
						Description: help.Text(`user-mgmt::list`),
						Action:      runtime(userMgmtList),
					},
					{
						Name:        `show`,
						Usage:       `Show information about a specific user`,
						Description: help.Text(`user-mgmt::show`),
						Action:      runtime(userMgmtShow),
					},
					{
						Name:        `sync`,
						Usage:       `List all registered users suitable for sync`,
						Description: help.Text(`user-mgmt::sync`),
						Action:      runtime(userMgmtSync),
					},
					{
						Name:        `activate`,
						Usage:       `Activate an inactive user account`,
						Description: help.Text(`supervisor::activate`),
						Action:      supervisorActivate,
					},
					{
						Name:        `admin`,
						Usage:       `SUBCOMMANDS for admin account management`,
						Description: help.Text(`admin-mgmt::`),
						Subcommands: []cli.Command{
							{
								Name:        `grant`,
								Usage:       `Grant admin account`,
								Description: help.Text(`admin-mgmt::add`),
								Action:      runtime(adminMgmtAdd),
							},
							{
								Name:        `revoke`,
								Usage:       `Revoke admin account`,
								Description: help.Text(`admin-mgmt::revoke`),
								Action:      runtime(adminMgmtRemove),
							},
						},
					},
					{
						Name:        `password`,
						Usage:       `SUBCOMMANDS for password management`,
						Description: help.Text(`supervisor::password`),
						Subcommands: []cli.Command{
							{
								Name:        `update`,
								Usage:       `Update the password of one's own account`,
								Description: help.Text(`supervisor::password-update`),
								Action:      boottime(supervisorPasswordUpdate),
							},
							{
								Name:        `reset`,
								Usage:       `Reset the password of one's own account via activation credentials`,
								Description: help.Text(`supervisor::password-reset`),
								Action:      boottime(supervisorPasswordReset),
							},
						},
					},
				},
			},
		}...,
	)
	return &app
}

// userMgmtAdd function
// soma user-mgmt add ${username} \
//      firstname ${fname} \
//      lastname ${lname} \
//      employeenr ${num} \
//      mailaddr ${addr} \
//      team ${team} \
//      [system ${bool}]
func userMgmtAdd(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{
		`firstname`,
		`lastname`,
		`employeenr`,
		`mailaddr`,
		`team`,
		`system`,
	}
	mandatoryOptions := []string{
		`firstname`,
		`lastname`,
		`employeenr`,
		`mailaddr`,
		`team`,
	}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	if err := adm.ValidateEmployeeNumber(opts[`employeenr`][0]); err != nil {
		return err
	}
	if err := adm.ValidateMailAddress(opts[`mailaddr`][0]); err != nil {
		return err
	}
	if err := adm.ValidateNotUUID(c.Args().First()); err != nil {
		return err
	}

	req := proto.NewUserRequest()
	req.User.UserName = c.Args().First()
	req.User.FirstName = opts[`firstname`][0]
	req.User.LastName = opts[`lastname`][0]
	req.User.MailAddress = opts[`mailaddr`][0]
	req.User.EmployeeNumber = opts[`employeenr`][0]
	if err := adm.LookupTeamID(
		opts[`team`][0],
		&req.User.TeamID,
	); err != nil {
		return err
	}
	req.User.IsDeleted = false
	req.User.IsActive = false
	req.User.IsSystem = false

	// optional argument
	if _, ok := opts[`system`]; ok {
		if err := adm.ValidateBool(
			opts[`system`][0],
			&req.User.IsSystem,
		); err != nil {
			return fmt.Errorf(
				"Syntax error, system argument not boolean: %s, %s",
				opts[`system`][0],
				err.Error(),
			)
		}
	}

	return adm.Perform(`postbody`, `/user/`, `command`, req, c)
}

// userMgmtUpdate function
// soma user-mgmt update ${userID} \
//      username ${name} \
//      firstname ${fname} \
//      lastname ${lname} \
//      employeenr ${num} \
//      mailaddr ${addr} \
//      team ${team} \
//      [deleted ${bool}]
func userMgmtUpdate(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{
		`username`,
		`firstname`,
		`lastname`,
		`employeenr`,
		`mailaddr`,
		`team`,
		`deleted`,
	}
	mandatoryOptions := []string{
		`username`,
		`firstname`,
		`lastname`,
		`employeenr`,
		`mailaddr`,
		`team`,
	}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	if err := adm.ValidateEmployeeNumber(opts[`employeenr`][0]); err != nil {
		return err
	}
	if err := adm.ValidateMailAddress(opts[`mailaddr`][0]); err != nil {
		return err
	}
	if !adm.IsUUID(c.Args().First()) {
		return fmt.Errorf(
			`user-mgmt update requiress UUID as first argument`,
		)
	}

	req := proto.NewUserRequest()
	req.User.ID = c.Args().First()
	req.User.UserName = opts[`username`][0]
	req.User.FirstName = opts[`firstname`][0]
	req.User.LastName = opts[`lastname`][0]
	req.User.MailAddress = opts[`mailaddr`][0]
	req.User.EmployeeNumber = opts[`employeenr`][0]
	if err := adm.LookupTeamID(
		opts[`team`][0],
		&req.User.TeamID,
	); err != nil {
		return err
	}

	// optional argument
	if _, ok := opts[`deleted`]; ok {
		if err := adm.ValidateBool(
			opts[`deleted`][0],
			&req.User.IsDeleted,
		); err != nil {
			return fmt.Errorf(
				"Syntax error, deleted argument not boolean: %s, %s",
				opts[`deleted`][0],
				err.Error(),
			)
		}
	}

	path := fmt.Sprintf("/user/%s", url.QueryEscape(req.User.ID))
	return adm.Perform(`putbody`, path, `command`, req, c)
}

// userMgmtRemove function
// soma user-mgmt remove ${username}
func userMgmtRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	userID, err := adm.LookupUserID(c.Args().First())
	if err != nil {
		return err
	}
	req := proto.Request{
		Flags: &proto.Flags{
			Purge: false,
		},
	}

	path := fmt.Sprintf("/user/%s", url.QueryEscape(userID))
	return adm.Perform(`deletebody`, path, `command`, req, c)
}

// userMgmtPurge function
// soma user-mgmt purge ${username}
func userMgmtPurge(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	userID, err := adm.LookupUserID(c.Args().First())
	if err != nil {
		return err
	}

	req := proto.Request{
		Flags: &proto.Flags{
			Purge: true,
		},
	}

	path := fmt.Sprintf("/user/%s", userID)
	return adm.Perform(`deletebody`, path, `command`, req, c)
}

// userMgmtList function
// soma user-mgmt list
func userMgmtList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/user/`, `list`, nil, c)
}

// userMgmtShow function
// soma user-mgmt show ${username}
func userMgmtShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	var (
		err error
		id  string
	)
	if id, err = adm.LookupUserID(c.Args().First()); err != nil {
		return err
	}

	path := fmt.Sprintf("/user/%s", id)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// userMgmtSync function
// soma user-mgmt sync
func userMgmtSync(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/sync/user/`, `list`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
