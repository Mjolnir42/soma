package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/internal/help"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerChecks(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  "checks",
				Usage: "SUBCOMMANDS for check configurations",
				Subcommands: []cli.Command{
					{
						Name:         "create",
						Usage:        "Create a new check configuration",
						Description:  help.Text(`ChecksCreate`),
						Action:       runtime(cmdCheckAdd),
						BashComplete: cmpl.CheckAdd,
					},
					{
						Name:         `delete`,
						Usage:        "Delete a check configuration",
						Description:  help.Text(`ChecksDelete`),
						Action:       runtime(cmdCheckDelete),
						BashComplete: cmpl.In,
					},
					{
						Name:         "list",
						Usage:        "List check configurations",
						Description:  help.Text(`ChecksList`),
						Action:       runtime(cmdCheckList),
						BashComplete: cmpl.In,
					},
					{
						Name:         "show",
						Usage:        "Show details about a check configuration",
						Description:  help.Text(`ChecksShow`),
						Action:       runtime(cmdCheckShow),
						BashComplete: cmpl.In,
					},
				},
			},
		}...,
	)
	return &app
}

func cmdCheckAdd(c *cli.Context) error {
	var (
		err    error
		teamID string
	)
	opts := make(map[string][]string)
	constraints := []proto.CheckConfigConstraint{}
	thresholds := []proto.CheckConfigThreshold{}

	if err = adm.ParseVariadicCheckArguments(
		opts,
		constraints,
		thresholds,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	req := proto.NewCheckConfigRequest()
	if err = adm.ValidateLBoundUint64(opts[`interval`][0],
		&req.CheckConfig.Interval, 1); err != nil {
		return err
	}
	if err = adm.ValidateRuneCount(c.Args().First(), 256); err != nil {
		return err
	}
	if req.CheckConfig.CapabilityID, err = adm.LookupCapabilityID(
		opts[`with`][0]); err != nil {
		return err
	}
	req.CheckConfig.ObjectType = opts[`on/type`][0]
	req.CheckConfig.Name = c.Args().First()
	req.CheckConfig.BucketID, err = adm.LookupBucketID(opts[`in`][0])
	if err != nil {
		return err
	}
	if req.CheckConfig.RepositoryID, err = adm.LookupRepoByBucket(
		req.CheckConfig.BucketID); err != nil {
		return err
	}
	if req.CheckConfig.ObjectID, err = adm.LookupCheckObjectID(
		opts[`on/type`][0],
		opts[`on/object`][0],
		req.CheckConfig.BucketID,
	); err != nil {
		return err
	}

	// clear bucketid if check is on a repository
	if req.CheckConfig.ObjectType == `repository` {
		req.CheckConfig.BucketID = ``
	}

	// optional argument: inheritance
	if iv, ok := opts[`inheritance`]; ok {
		if err = adm.ValidateBool(iv[0],
			&req.CheckConfig.Inheritance); err != nil {
			return err
		}
	} else {
		// inheritance defaults to true
		req.CheckConfig.Inheritance = true
	}

	// optional argument: childrenonly
	if co, ok := opts[`childrenonly`]; ok {
		if err = adm.ValidateBool(co[0],
			&req.CheckConfig.ChildrenOnly); err != nil {
			return err
		}
	} else {
		// childrenonly defaults to false
		req.CheckConfig.ChildrenOnly = false
	}

	// optional argument: extern
	if ex, ok := opts[`extern`]; ok {
		if err = adm.ValidateRuneCount(ex[0], 64); err != nil {
			return err
		}
		req.CheckConfig.ExternalID = ex[0]
	}

	if teamID, err = adm.LookupTeamByRepo(
		req.CheckConfig.RepositoryID); err != nil {
		return err
	}

	if req.CheckConfig.Thresholds, err = adm.ValidateThresholds(
		thresholds,
	); err != nil {
		return err
	}

	if req.CheckConfig.Constraints, err = adm.ValidateCheckConstraints(
		req.CheckConfig.RepositoryID,
		teamID,
		constraints,
	); err != nil {
		return err
	}

	path := fmt.Sprintf("/checks/%s/", req.CheckConfig.RepositoryID)
	return adm.Perform(`postbody`, path, `command`, req, nil)
}

func cmdCheckDelete(c *cli.Context) error {
	unique := []string{`in`}
	required := []string{`in`}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		unique,
		required,
		c.Args().Tail(),
	); err != nil {
		return err
	}
	var (
		err                       error
		bucketID, repoID, checkID string
	)
	if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
		return err
	}
	if repoID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}
	if checkID, _, err = adm.LookupCheckConfigID(c.Args().First(),
		repoID, ``); err != nil {
		return err
	}

	path := fmt.Sprintf("/checks/%s/%s", repoID, checkID)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdCheckList(c *cli.Context) error {
	unique := []string{`in`}
	required := []string{`in`}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		unique,
		required,
		adm.AllArguments(c),
	); err != nil {
		return err
	}
	var (
		err              error
		bucketID, repoID string
	)
	bucketID, err = adm.LookupBucketID(opts[`in`][0])
	if err != nil {
		return err
	}
	if repoID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}

	path := fmt.Sprintf("/checks/%s/", repoID)
	return adm.Perform(`get`, path, `list`, nil, c)
}

func cmdCheckShow(c *cli.Context) error {
	unique := []string{`in`}
	required := []string{`in`}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		unique,
		required,
		c.Args().Tail()); err != nil {
		return err
	}
	var (
		err                       error
		bucketID, repoID, checkID string
	)
	if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
		return err
	}
	if repoID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}
	if checkID, _, err = adm.LookupCheckConfigID(c.Args().First(),
		repoID, ``); err != nil {
		return err
	}

	path := fmt.Sprintf("/checks/%s/%s", repoID, checkID)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix