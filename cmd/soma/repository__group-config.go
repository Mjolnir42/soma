/*-
 * Copyright (c) 2016-2019, Jörg Pernfuß
 * Copyright (c) 2019, 1&1 IONOS SE
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

func registerGroups(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:        `group`,
				Usage:       `SUBCOMMANDS for group management`,
				Description: help.Text(`group-config::`),
				Subcommands: []cli.Command{
					{
						Name:         `create`,
						Usage:        `Create a new group`,
						Description:  help.Text(`group-config::create`),
						Action:       runtime(groupConfigCreate),
						BashComplete: cmpl.In,
					},
					{
						Name:         `destroy`,
						Usage:        `Destroy an existing group inside a tree bucket`,
						Description:  help.Text(`group-config::destroy`),
						Action:       runtime(groupConfigDestroy),
						BashComplete: cmpl.In,
					},
					{
						Name:         `list`,
						Usage:        `List all groups in a bucket`,
						Description:  help.Text(`group-config::create`),
						Action:       runtime(groupConfigList),
						BashComplete: cmpl.DirectIn,
					},
					{
						Name:         `show`,
						Usage:        `Show full details about a specific group`,
						Description:  help.Text(`group-config::show`),
						Action:       runtime(groupConfigShow),
						BashComplete: cmpl.In,
					},
					{
						Name:         `dumptree`,
						Usage:        `Display the group as tree`,
						Description:  help.Text(`group-config::tree`),
						Action:       runtime(groupConfigTree),
						BashComplete: cmpl.In,
					},
					{
						Name:        `property`,
						Usage:       `SUBCOMMANDS for properties on groups`,
						Description: help.Text(`group-config::`),
						Subcommands: []cli.Command{
							{
								Name:        `create`,
								Usage:       `SUBCOMMANDS to create properties`,
								Description: help.Text(`group-config::property-create`),
								Subcommands: []cli.Command{
									{
										Name:         `system`,
										Usage:        `Add a system property to a group`,
										Description:  help.Text(`group-config::property-create`),
										Action:       runtime(groupConfigPropertyCreateSystem),
										BashComplete: cmpl.PropertyCreateInValue,
									},
									{
										Name:         `service`,
										Usage:        `Add a service property to a group`,
										Description:  help.Text(`group-config::property-create`),
										Action:       runtime(groupConfigPropertyCreateService),
										BashComplete: cmpl.PropertyCreateInValue,
									},
									{
										Name:         `oncall`,
										Usage:        `Add an oncall property to a group`,
										Description:  help.Text(`group-config::property-create`),
										Action:       runtime(groupConfigPropertyCreateOncall),
										BashComplete: cmpl.PropertyCreateIn,
									},
									{
										Name:         `custom`,
										Usage:        `Add a custom property to a group`,
										Description:  help.Text(`group-config::property-create`),
										Action:       runtime(groupConfigPropertyCreateCustom),
										BashComplete: cmpl.PropertyCreateIn,
									},
								},
							},
							{
								Name:        `destroy`,
								Usage:       `SUBCOMMANDS to destroy properties`,
								Description: help.Text(`group-config::property-destroy`),
								Subcommands: []cli.Command{
									{
										Name:         `system`,
										Usage:        `Delete a system property from a group`,
										Description:  help.Text(`group-config::property-destroy`),
										Action:       runtime(groupConfigPropertyDestroySystem),
										BashComplete: cmpl.PropertyOnInView,
									},
									{
										Name:         `service`,
										Usage:        `Delete a service property from a group`,
										Description:  help.Text(`group-config::property-destroy`),
										Action:       runtime(groupConfigPropertyDestroyService),
										BashComplete: cmpl.PropertyOnInView,
									},
									{
										Name:         `oncall`,
										Usage:        `Delete an oncall property from a group`,
										Description:  help.Text(`group-config::property-destroy`),
										Action:       runtime(groupConfigPropertyDestroyOncall),
										BashComplete: cmpl.PropertyOnInView,
									},
									{
										Name:         `custom`,
										Usage:        `Delete a custom property from a group`,
										Description:  help.Text(`group-config::property-destroy`),
										Action:       runtime(groupConfigPropertyDestroyCustom),
										BashComplete: cmpl.PropertyOnInView,
									},
								},
							},
						},
					},
					{
						Name:        `member`,
						Usage:       `SUBCOMMANDS for group membership management`,
						Description: help.Text(`group-config::`),
						Subcommands: []cli.Command{
							{
								Name:         `list`,
								Usage:        `List all members of a group`,
								Description:  help.Text(`group-config::member-list`),
								Action:       runtime(groupConfigMemberList),
								BashComplete: cmpl.DirectInOf,
							},
							{
								Name:        `assign`,
								Usage:       `SUBCOMMANDS for assigning objects to a group`,
								Description: help.Text(`group-config::member-assign`),
								Subcommands: []cli.Command{
									{
										Name:         `group`,
										Usage:        `Assign a group to another group`,
										Description:  help.Text(`group-config::member-assign`),
										Action:       runtime(groupConfigMemberAssignGroup),
										BashComplete: cmpl.InTo,
									},
									{
										Name:         `cluster`,
										Usage:        `Assign a cluster to a group`,
										Description:  help.Text(`group-config::member-assign`),
										Action:       runtime(groupConfigMemberAssignCluster),
										BashComplete: cmpl.InTo,
									},
									{
										Name:         `node`,
										Usage:        `Assign a node to a group`,
										Description:  help.Text(`group-config::member-assign`),
										Action:       runtime(groupConfigMemberAssignNode),
										BashComplete: cmpl.InTo,
									},
								},
							},
							{
								Name:        `unassign`,
								Usage:       `SUBCOMMANDS to unassign objects from a group`,
								Description: help.Text(`group-config::member-unassign`),
								Subcommands: []cli.Command{
									{
										Name:         `group`,
										Usage:        `Unassign a child group from its parent group`,
										Description:  help.Text(`group-config::member-unassign`),
										Action:       runtime(groupConfigMemberUnassignGroup),
										BashComplete: cmpl.InFrom,
									},
									{
										Name:         `cluster`,
										Usage:        `Unassign a child cluster from its parent group`,
										Description:  help.Text(`group-config::member-unassign`),
										Action:       runtime(groupConfigMemberUnassignCluster),
										BashComplete: cmpl.InFrom,
									},
									{
										Name:         `node`,
										Usage:        `Unassign a node from its parent group`,
										Description:  help.Text(`group-config::member-unassign`),
										Action:       runtime(groupConfigMemberUnassignNode),
										BashComplete: cmpl.InFrom,
									},
								},
							},
						},
					},
				},
			},
		}...,
	)
	return &app
}

// groupConfigCreate function
// soma group create ${group} in ${bucket}
func groupConfigCreate(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`in`}
	mandatoryOptions := []string{`in`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	var err error
	var repositoryID, bucketID string
	if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
		return err
	}

	if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}

	req := proto.NewGroupRequest()
	req.Group.Name = c.Args().First()
	req.Group.RepositoryID = repositoryID
	req.Group.BucketID = bucketID

	if err := adm.ValidateRuneCountRange(
		req.Group.Name, 2, 256); err != nil {
		return err
	}
	if err := adm.ValidateNotUUID(req.Group.Name); err != nil {
		return err
	}

	path := fmt.Sprintf(
		"/repository/%s/bucket/%s/group/",
		url.QueryEscape(repositoryID),
		url.QueryEscape(bucketID),
	)
	return adm.Perform(`postbody`, path, `group-config::create`, req, c)
}

// groupConfigDestroy function
// soma group destroy ${group} in ${bucket}
func groupConfigDestroy(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`in`}
	mandatoryOptions := []string{`in`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	var err error
	var repositoryID, bucketID, groupID string
	if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
		return err
	}
	if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}
	if groupID, err = adm.LookupGroupID(c.Args().First(), bucketID); err != nil {
		return err
	}

	path := fmt.Sprintf(
		"/repository/%s/bucket/%s/group/%s",
		url.QueryEscape(repositoryID),
		url.QueryEscape(bucketID),
		url.QueryEscape(groupID),
	)
	return adm.Perform(`delete`, path, `group-config::destroy`, nil, c)
}

// groupConfigList function
// soma group list in ${bucket}
func groupConfigList(c *cli.Context) error {
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

	var err error
	var repositoryID, bucketID string
	if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
		return err
	}

	if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}

	path := fmt.Sprintf(
		"/repository/%s/bucket/%s/group/",
		url.QueryEscape(repositoryID),
		url.QueryEscape(bucketID),
	)
	return adm.Perform(`get`, path, `group-config::list`, nil, c)
}

// groupConfigShow function
// soma group show ${group} in ${bucket}
func groupConfigShow(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`in`}
	mandatoryOptions := []string{`in`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	var err error
	var repositoryID, bucketID, groupID string
	if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
		return err
	}
	if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}
	if groupID, err = adm.LookupGroupID(c.Args().First(), bucketID); err != nil {
		return err
	}

	path := fmt.Sprintf(
		"/repository/%s/bucket/%s/group/%s",
		url.QueryEscape(repositoryID),
		url.QueryEscape(bucketID),
		url.QueryEscape(groupID),
	)
	return adm.Perform(`get`, path, `group-config::show`, nil, c)
}

// groupConfigTree function
// soma group dumptree ${group} in ${bucket}
func groupConfigTree(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`in`}
	mandatoryOptions := []string{`in`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	var err error
	var repositoryID, bucketID, groupID string
	if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
		return err
	}
	if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}
	if groupID, err = adm.LookupGroupID(c.Args().First(), bucketID); err != nil {
		return err
	}

	path := fmt.Sprintf(
		"/repository/%s/bucket/%s/group/%s/tree",
		url.QueryEscape(repositoryID),
		url.QueryEscape(bucketID),
		url.QueryEscape(groupID),
	)
	return adm.Perform(`get`, path, `group-config::tree`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
