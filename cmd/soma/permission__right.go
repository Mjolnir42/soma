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
				Name:        `right`,
				Usage:       `SUBCOMMANDS for permission grant`,
				Description: help.Text(`right::`),
				Subcommands: []cli.Command{
					{
						Name:         `grant`,
						Usage:        `Grant a permission`,
						Action:       runtime(rightGrant),
						Description:  help.Text(`right::grant`),
						BashComplete: cmpl.TripleToOn,
					},
					{
						Name:         `revoke`,
						Usage:        `Revoke a permission grant`,
						Action:       runtime(rightRevoke),
						Description:  help.Text(`right::revoke`),
						BashComplete: cmpl.TripleFromOn,
					},
					{
						Name:         `list`,
						Usage:        `List all grants of a permission`,
						Action:       runtime(rightList),
						Description:  help.Text(`right::list`),
						BashComplete: cmpl.None,
					},
					/*
						{
							Name:   `show`,
							Usage:  `Show a permission grant for a recipient`,
							Action: runtime(cmdRightShow),
							// BashComplete: cmpl.Triple_For,
							Description: help.Text(`RightsShow`),
						},
					*/
				},
			},
		}...,
	)
	return &app
}

// rightGrant function
// soma right grant $category::$permission
//            to user|admin|team $name
//           [on repository|bucket|monitoring|team $name]
func rightGrant(c *cli.Context) error {
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

	var err error
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
	case msg.CategoryGlobal,
		msg.CategoryIdentity,
		msg.CategoryOperation,
		msg.CategoryPermission,
		msg.CategorySelf,
		msg.CategorySystem:
		fallthrough
	case msg.CategoryGrantGlobal,
		msg.CategoryGrantIdentity,
		msg.CategoryGrantOperation,
		msg.CategoryGrantPermission,
		msg.CategoryGrantSelf:
		if len(opts[`on`]) != 0 {
			return fmt.Errorf("Permissions in category %s are global"+
				" and require no 'on' keyword target.",
				req.Grant.Category)
		}
	case msg.CategoryMonitoring,
		msg.CategoryRepository,
		msg.CategoryTeam:
		fallthrough
	case msg.CategoryGrantMonitoring,
		msg.CategoryGrantRepository,
		msg.CategoryGrantTeam:
		if len(opts[`on`]) != 1 {
			return fmt.Errorf("Permissions in category %s require a"+
				" target, specified via 'on' keyword.",
				req.Grant.Category)
		}
	}

	// lookup permissionID
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
		req.Grant.RecipientType = `admin`
		if req.Grant.RecipientID, err = adm.LookupAdminID(
			opts[`to`][0][1]); err != nil {
			return err
		}
	case `team`:
		req.Grant.RecipientType = `team`
		if err = adm.LookupTeamID(
			opts[`to`][0][1],
			&req.Grant.RecipientID,
		); err != nil {
			return err
		}
	case `tool`:
		return fmt.Errorf(`Tool permissions are not implemented.`)
	}

	if len(opts[`on`]) == 1 {
		switch req.Grant.Category {
		case msg.CategoryRepository,
			msg.CategoryGrantRepository:
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
		case msg.CategoryTeam,
			msg.CategoryGrantTeam:
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
		case msg.CategoryMonitoring,
			msg.CategoryGrantMonitoring:
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
	return adm.Perform(`postbody`, path, `right::grant`, req, c)
}

// rightRevoke function
// soma right revoke $category::$permission
//            from user|admin|team $name
//           [on repository|bucket|monitoring|team $name]
func rightRevoke(c *cli.Context) error {
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

	var err error
	var grantID string
	req := proto.NewGrantFilter()

	permissionSlice := strings.Split(c.Args().First(), `::`)
	if len(permissionSlice) != 2 {
		return fmt.Errorf("Invalid split of permission into %s",
			permissionSlice)
	}
	req.Filter.Grant.Category = permissionSlice[0]

	// validate category
	if err = adm.ValidateCategory(req.Filter.Grant.Category); err != nil {
		return err
	}

	// check optional argument chain
	switch req.Filter.Grant.Category {
	case msg.CategoryGlobal,
		msg.CategoryIdentity,
		msg.CategoryOperation,
		msg.CategoryPermission,
		msg.CategorySelf,
		msg.CategorySystem:
		fallthrough
	case msg.CategoryGrantGlobal,
		msg.CategoryGrantIdentity,
		msg.CategoryGrantOperation,
		msg.CategoryGrantPermission,
		msg.CategoryGrantSelf:
		if len(opts[`on`]) != 0 {
			return fmt.Errorf("Permissions in category %s are global"+
				" and require no 'on' keyword target.",
				req.Filter.Grant.Category)
		}
	case msg.CategoryMonitoring,
		msg.CategoryRepository,
		msg.CategoryTeam:
		fallthrough
	case msg.CategoryGrantMonitoring,
		msg.CategoryGrantRepository,
		msg.CategoryGrantTeam:
		if len(opts[`on`]) != 1 {
			return fmt.Errorf("Permissions in category %s require a"+
				" target, specified via 'on' keyword.",
				req.Filter.Grant.Category)
		}
	default:
		return fmt.Errorf("Unknown category: %s", permissionSlice[0])
	}

	// lookup permissionID
	if err = adm.LookupPermIDRef(
		permissionSlice[1],
		req.Filter.Grant.Category,
		&req.Filter.Grant.PermissionID,
	); err != nil {
		return err
	}

	// lookup recipientID
	req.Filter.Grant.RecipientType = opts[`from`][0][0]
	switch req.Filter.Grant.RecipientType {
	case `user`:
		if req.Filter.Grant.RecipientID, err = adm.LookupUserID(
			opts[`from`][0][1],
		); err != nil {
			return err
		}
	case `admin`:
		if req.Filter.Grant.RecipientID, err = adm.LookupAdminID(
			opts[`from`][0][1],
		); err != nil {
			return err
		}
	case `team`:
		req.Filter.Grant.RecipientType = `team`
		if err = adm.LookupTeamID(
			opts[`from`][0][1],
			&req.Filter.Grant.RecipientID,
		); err != nil {
			return err
		}
	case `tool`:
		return fmt.Errorf(`Tool permissions are not implemented.`)
	}

	// parse optional object spec
	if len(opts[`on`]) == 1 {
		switch req.Filter.Grant.Category {
		case msg.CategoryRepository,
			msg.CategoryGrantRepository:
			switch opts[`on`][0][0] {
			case msg.EntityRepository:
				req.Filter.Grant.ObjectType = msg.EntityRepository
				if req.Filter.Grant.ObjectID, err = adm.LookupRepoID(
					opts[`on`][0][1],
				); err != nil {
					return err
				}
			case msg.EntityBucket:
				req.Filter.Grant.ObjectType = msg.EntityBucket
				if req.Filter.Grant.ObjectID, err = adm.LookupBucketID(
					opts[`on`][0][1],
				); err != nil {
					return err
				}
			default:
				return fmt.Errorf(`Invalid`)
			}
		case msg.CategoryTeam,
			msg.CategoryGrantTeam:
			switch opts[`on`][0][0] {
			case msg.EntityTeam:
				req.Filter.Grant.ObjectType = msg.EntityTeam
				if err = adm.LookupTeamID(
					opts[`on`][0][1],
					&req.Filter.Grant.ObjectID,
				); err != nil {
					return err
				}
			default:
				return fmt.Errorf(`Invalid`)
			}
		case msg.CategoryMonitoring,
			msg.CategoryGrantMonitoring:
			switch opts[`on`][0][0] {
			case msg.EntityMonitoring:
				req.Filter.Grant.ObjectType = msg.EntityMonitoring
				if req.Filter.Grant.ObjectID, err = adm.LookupMonitoringID(
					opts[`on`][0][1],
				); err != nil {
					return err
				}
			default:
				return fmt.Errorf(`Invalid`)
			}
		}
	}

	// lookup grantID
	if err = adm.LookupGrantIDRef(
		req,
		&grantID,
	); err != nil {
		return err
	}

	path := fmt.Sprintf("/category/%s/permission/%s/grant/%s",
		req.Filter.Grant.Category,
		req.Filter.Grant.PermissionID,
		grantID,
	)
	return adm.Perform(`delete`, path, `right::revoke`, nil, c)
}

// rightList function
// soma right list $category::$permission
func rightList(c *cli.Context) error {
	var permissionID string

	permissionSlice := strings.Split(c.Args().First(), `::`)
	if len(permissionSlice) != 2 {
		return fmt.Errorf("Invalid split of permission into %s",
			permissionSlice)
	}

	// validate category
	if err := adm.ValidateCategory(permissionSlice[0]); err != nil {
		return err
	}

	// lookup permissionID
	if err := adm.LookupPermIDRef(
		permissionSlice[1],
		permissionSlice[0],
		&permissionID,
	); err != nil {
		return err
	}

	path := fmt.Sprintf("/category/%s/permission/%s/grant/",
		permissionSlice[0],
		permissionID,
	)
	return adm.Perform(`get`, path, `right::list`, nil, c)
}

func cmdRightShow(c *cli.Context) error {
	return fmt.Errorf(`Not implemented - TODO`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
