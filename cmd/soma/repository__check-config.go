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

func registerChecks(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  `check-config`,
				Usage: `SUBCOMMANDS for check configuration management`,
				Subcommands: []cli.Command{
					{
						Name:         `create`,
						Usage:        `Create a new check configuration`,
						Description:  help.Text(`check-config::create`),
						Action:       runtime(checkConfigCreate),
						BashComplete: cmpl.CheckConfigCreate,
					},
					{
						Name:         `destroy`,
						Usage:        "Destroy a check configuration",
						Description:  help.Text(`check-config::destroy`),
						Action:       runtime(checkConfigDestroy),
						BashComplete: cmpl.CheckConfigDestroy,
					},
					{
						Name:         `list`,
						Usage:        `List check configurations in a repository`,
						Description:  help.Text(`check-config::list`),
						Action:       runtime(checkConfigList),
						BashComplete: cmpl.DirectIn,
					},
					{
						Name:         `show`,
						Usage:        `Show details about a check configuration`,
						Description:  help.Text(`check-config::show`),
						Action:       runtime(checkConfigShow),
						BashComplete: cmpl.In,
					},
				},
			},
		}...,
	)
	return &app
}

// checkConfigCreate function
// soma check-config create ...
func checkConfigCreate(c *cli.Context) error {
	var err error
	var teamID string
	opts := map[string][]string{}
	constraints := []proto.CheckConfigConstraint{}
	thresholds := []proto.CheckConfigThreshold{}
	req := proto.NewCheckConfigRequest()

	if err := adm.ParseVariadicCheckArguments(
		&constraints,
		&thresholds,
		opts,
		c.Args().Tail(),
	); err != nil {
		return err
	}
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
	if err = adm.ValidateNotUUID(req.CheckConfig.Name); err != nil {
		return err
	}
	fmt.Println("Pre Switch statement")
	switch req.CheckConfig.ObjectType {
	case `repository`:
		if req.CheckConfig.RepositoryID, err = adm.LookupRepoID(opts[`on/object`][0]); err != nil {
			return err
		}
		req.CheckConfig.ObjectID = req.CheckConfig.RepositoryID
	case `bucket`:
		if req.CheckConfig.BucketID, err = adm.LookupBucketID(opts[`on/object`][0]); err != nil {
			return err
		}
		if req.CheckConfig.RepositoryID, err = adm.LookupRepoByBucket(req.CheckConfig.BucketID); err != nil {
			return err
		}
		req.CheckConfig.ObjectID = req.CheckConfig.BucketID
	case `node`:
		if req.CheckConfig.ObjectID, err = adm.LookupNodeID(opts[`on/object`][0]); err != nil {
			return err
		}
		config := &proto.NodeConfig{}
		if config, err = adm.LookupNodeConfig(req.CheckConfig.ObjectID); err != nil {
			return err
		}
		req.CheckConfig.BucketID = config.BucketID
		req.CheckConfig.RepositoryID = config.RepositoryID
	case `group`, `cluster`:
		if req.CheckConfig.BucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
			return err
		}
		if req.CheckConfig.RepositoryID, err = adm.LookupRepoByBucket(req.CheckConfig.BucketID); err != nil {
			return err
		}
		if req.CheckConfig.ObjectID, err = adm.LookupCheckObjectID(
			req.CheckConfig.ObjectType, opts[`on/object`][0],
			req.CheckConfig.BucketID,
		); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown object entity: %s", req.CheckConfig.ObjectType)
	}
	fmt.Println("Post Switch statement")

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
	fmt.Println("Pre LookupTeamByRepo")
	if err = adm.LookupTeamByRepo(
		req.CheckConfig.RepositoryID, &teamID); err != nil {
		return err
	}
	fmt.Println("Pre ValidateThresholds")
	if req.CheckConfig.Thresholds, err = adm.ValidateThresholds(
		thresholds,
	); err != nil {
		return err
	}
	fmt.Println("Pre ValidateCheckConstraints")
	if req.CheckConfig.Constraints, err = adm.ValidateCheckConstraints(
		req.CheckConfig.RepositoryID,
		teamID,
		constraints,
	); err != nil {
		return err
	}
	path := fmt.Sprintf("/checkconfig/%s/",
		url.QueryEscape(req.CheckConfig.RepositoryID),
	)
	fmt.Println("Pre Perform")
	return adm.Perform(`postbody`, path, `check-config::create`, req, nil)
}

// checkConfigDestroy function
// soma check-config destroy ${name} in repository|bucket ${repo|bucket}
func checkConfigDestroy(c *cli.Context) error {
	opts := map[string][][2]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`in`}
	mandatoryOptions := []string{`in`}

	var err error
	if err = adm.ParseVariadicTriples(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	var bucketID, repoID, checkID string

	switch opts[`in`][0][0] {
	case proto.EntityRepository:
		if repoID, err = adm.LookupRepoID(opts[`in`][0][1]); err != nil {
			return err
		}
	case proto.EntityBucket:
		if bucketID, err = adm.LookupBucketID(opts[`in`][0][1]); err != nil {
			return err
		}
		if repoID, err = adm.LookupRepoByBucket(bucketID); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Invalid entity: %s", opts[`in`][0][0])
	}

	if checkID, _, err = adm.LookupCheckConfigID(c.Args().First(), repoID, ``); err != nil {
		return err
	}

	path := fmt.Sprintf("/checkconfig/%s/%s",
		url.QueryEscape(repoID),
		url.QueryEscape(checkID),
	)
	return adm.Perform(`delete`, path, `check-config::destroy`, nil, c)
}

// checkConfigList function
// soma check-config list in ${repository}
func checkConfigList(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`in`}
	mandatoryOptions := []string{`in`}

	var err error
	if err = adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		adm.AllArguments(c),
	); err != nil {
		return err
	}

	var repoID string
	if repoID, err = adm.LookupRepoID(opts[`in`][0]); err != nil {
		return err
	}
	path := fmt.Sprintf("/checkconfig/%s/",
		url.QueryEscape(repoID),
	)
	return adm.Perform(`get`, path, `check-config::list`, nil, c)
}

// checkConfigShow function
// soma check-config show ${name} in ${repository}
func checkConfigShow(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`in`}
	mandatoryOptions := []string{`in`}

	var err error
	if err = adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	var repoID, checkID string
	if repoID, err = adm.LookupRepoID(opts[`in`][0]); err != nil {
		return err
	}
	if checkID, _, err = adm.LookupCheckConfigID(c.Args().First(),
		repoID, ``); err != nil {
		return err
	}

	path := fmt.Sprintf("/checkconfig/%s/%s",
		url.QueryEscape(repoID),
		url.QueryEscape(checkID),
	)
	return adm.Perform(`get`, path, `check-config::list`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
