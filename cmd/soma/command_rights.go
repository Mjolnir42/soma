package main

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/internal/help"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerRights(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  "rights",
				Usage: "SUBCOMMANDS for rights",
				Subcommands: []cli.Command{
					{
						Name:         "grant",
						Usage:        "Grant a permission",
						Action:       runtime(cmdRightGrant),
						Description:  help.Text(`RightsGrant`),
						BashComplete: cmpl.TripleToOn,
					},
					{
						Name:         "revoke",
						Usage:        "Revoke a permission",
						Action:       runtime(cmdRightRevoke),
						Description:  help.Text(`RightsRevoke`),
						BashComplete: cmpl.TripleFromOn,
					},
					{
						Name:        `list`,
						Usage:       `List all grants of a permission`,
						Action:      runtime(cmdRightList),
						Description: help.Text(`RightsList`),
					},
					{
						Name:   `show`,
						Usage:  `Show a permission grant for a recipient`,
						Action: runtime(cmdRightShow),
						// BashComplete: cmpl.Triple_For,
						Description: help.Text(`RightsShow`),
					},
				},
			},
		}...,
	)
	return &app
}

func cmdRightGrant(c *cli.Context) error {
	opts := map[string][][2]string{}
	if err := adm.ParseVariadicTriples(
		opts,
		[]string{},
		[]string{`to`, `on`},
		[]string{`to`},
		c.Args().Tail(),
	); err != nil {
		return err
	}
	var (
		err error
	)
	req := proto.NewGrantRequest()

	permissionSlice := strings.Split(c.Args().First(), `::`)
	if len(permissionSlice) != 2 {
		return fmt.Errorf("Invalid split of permission into %s",
			permissionSlice)
	}
	// validate category
	req.Grant.Category = permissionSlice[0]
	if err = adm.ValidateCategory(req.Grant.Category); err != nil {
		return err
	}

	// check optional argument chain
	switch req.Grant.Category {
	case msg.CategoryGlobal, msg.CategoryIdentity, msg.CategoryOperation,
		msg.CategoryPermission, msg.CategorySelf, msg.CategorySystem:
		fallthrough
	case msg.CategoryGrantGlobal, msg.CategoryGrantIdentity, msg.CategoryGrantOperation,
		msg.CategoryGrantPermission, msg.CategoryGrantSelf:
		if len(opts[`on`]) != 0 {
			return fmt.Errorf("Permissions in category %s are global"+
				" and require no 'on' keyword target.",
				req.Grant.Category)
		}
	case msg.CategoryMonitoring, msg.CategoryRepository, msg.CategoryTeam:
		fallthrough
	case msg.CategoryGrantMonitoring, msg.CategoryGrantRepository, msg.CategoryGrantTeam:
		if len(opts[`on`]) != 1 {
			return fmt.Errorf("Permissions in category %s require a"+
				" target, specified via 'on' keyword.",
				req.Grant.Category)
		}
	}

	// lookup permissionid
	if err = adm.LookupPermIDRef(
		permissionSlice[1],
		req.Grant.Category,
		&req.Grant.PermissionID,
	); err != nil {
		return err
	}

	// check that the permission is granted to a valid entity
	if err = adm.VerifyPermissionTarget(opts[`to`][0][0]); err != nil {
		return err
	}
	switch opts[`to`][0][0] {
	case `user`:
		req.Grant.RecipientType = `user`
		if req.Grant.RecipientID, err = adm.LookupUserID(
			opts[`to`][0][1]); err != nil {
			return err
		}
	case `admin`:
		return fmt.Errorf(`Admin permissions are not implemented.`)
	case `tool`:
		return fmt.Errorf(`Tool permissions are not implemented.`)
	case `team`:
		return fmt.Errorf(`Team permissions are not implemented.`)
	}

	if len(opts[`on`]) == 1 {
		switch req.Grant.Category {
		case msg.CategoryRepository, msg.CategoryGrantRepository:
			switch opts[`on`][0][0] {
			case msg.EntityRepository:
				req.Grant.ObjectType = msg.EntityRepository
				if req.Grant.ObjectID, err = adm.LookupRepoID(
					opts[`on`][0][1],
				); err != nil {
					return err
				}
			case msg.EntityBucket:
				req.Grant.ObjectType = msg.EntityBucket
				if req.Grant.ObjectID, err = adm.LookupBucketID(
					opts[`on`][0][1],
				); err != nil {
					return err
				}
			default:
				return fmt.Errorf(`Invalid`)
			}
		case msg.CategoryTeam, msg.CategoryGrantTeam:
			switch opts[`on`][0][0] {
			case msg.EntityTeam:
				req.Grant.ObjectType = msg.EntityTeam
				if err = adm.LookupTeamID(
					opts[`on`][0][1],
					&req.Grant.ObjectID,
				); err != nil {
					return err
				}
			default:
				return fmt.Errorf(`Invalid`)
			}
		case msg.CategoryMonitoring, msg.CategoryGrantMonitoring:
			switch opts[`on`][0][0] {
			case msg.EntityMonitoring:
				req.Grant.ObjectType = msg.EntityMonitoring
				if req.Grant.ObjectID, err = adm.LookupMonitoringID(
					opts[`on`][0][1],
				); err != nil {
					return err
				}
			default:
				return fmt.Errorf(`Invalid`)
			}
		}
	}

	path := fmt.Sprintf("/category/%s/permission/%s/grant/",
		req.Grant.Category, req.Grant.PermissionID)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

func cmdRightRevoke(c *cli.Context) error {
	// XXX TODO
	return fmt.Errorf(`Not implemented - TODO`)
	/*
		opts := map[string][][2]string{}
		if err := adm.ParseVariadicTriples(
			opts,
			[]string{},
			[]string{`from`, `on`},
			[]string{`from`},
			c.Args().Tail(),
		); err != nil {
			return err
		}

		var (
			err                     error
			userId, permId, grantId string
		)
		if err = adm.LookupPermIDRef(c.Args().First(),
			`foobar`, // dummy value for new structs
			&permId); err != nil {
			return err
		}
		if userId, err = adm.LookupUserID(opts[`from`][0][1]); err != nil {
			return err
		}
		if err = adm.LookupGrantIdRef(`user`, userId, permId, `category`,
			&grantId); err != nil {
			return err
		}

		path := fmt.Sprintf( "/category/%s/permission/%s/grant/%s"
			`category`, permId, grantId)
		return adm.Perform(`delete`, path, `command`, nil, c)
	*/
}

func cmdRightList(c *cli.Context) error {
	return fmt.Errorf(`Not implemented - TODO`)
}

func cmdRightShow(c *cli.Context) error {
	return fmt.Errorf(`Not implemented - TODO`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
