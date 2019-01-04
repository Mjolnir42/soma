/*-
 * Copyright (c) 2015-2019, Jörg Pernfuß
 * Copyright (c) 2018-2019, 1&1 IONOS SE
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main // import "github.com/mjolnir42/soma/cmd/soma"

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/lib/proto"
)

// bucketCreate function
// soma bucket create ${bucket} in ${repository} environment ${env}
func bucketCreate(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`in`, `environment`}
	mandatoryOptions := []string{`in`, `environment`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	// check the length of the bucket name
	if err := adm.ValidateRuneCountRange(c.Args().First(), 4, 512); err != nil {
		return err
	}

	// fetch list of environments from SOMA to check if a valid
	// environment was requested
	if err := adm.ValidateEnvironment(opts[`environment`][0]); err != nil {
		return err
	}

	var err error
	var repositoryID, repositoryName, teamID string
	if repositoryID, err = adm.LookupRepoID(opts[`in`][0]); err != nil {
		return err
	}
	if err = adm.LookupRepoName(repositoryID, &repositoryName); err != nil {
		return err
	}
	if err = adm.LookupTeamByRepo(repositoryID, &teamID); err != nil {
		return err
	}

	// check if the prefix constraint if fulfilled
	if !strings.HasPrefix(c.Args().First(), repositoryName) {
		return fmt.Errorf("Repository name %s must be a prefix of bucket name %s",
			repositoryName, c.Args().First())
	}

	req := proto.NewBucketRequest()
	req.Bucket.Name = c.Args().First()
	req.Bucket.RepositoryID = repositoryID
	req.Bucket.TeamID = teamID
	req.Bucket.Environment = opts[`environment`][0]

	path := fmt.Sprintf("/repository/%s/bucket/", repositoryID)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

// bucketDestroy function
// soma bucket destroy ${bucket} [in ${repository}]
func bucketDestroy(c *cli.Context) error {
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

	var err error
	var repositoryID, bucketID, repositoryControlID string
	if bucketID, err = adm.LookupBucketID(c.Args().First()); err != nil {
		return err
	}
	if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}

	// optional argument, must be correct if provided
	if _, ok := opts[`in`]; ok {
		if repositoryControlID, err = adm.LookupRepoID(opts[`in`][0]); err != nil {
			return err
		} else if repositoryControlID != repositoryID {
			return fmt.Errorf("bucket %s is not in repository %s", c.Args().First(), opts[`in`][0])
		}
	}

	path := fmt.Sprintf(
		"/repository/%s/bucket/%s",
		url.QueryEscape(repositoryID),
		url.QueryEscape(bucketID),
	)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// bucketList function
// soma bucket list in ${repository}
func bucketList(c *cli.Context) error {
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
	var repositoryID string
	if repositoryID, err = adm.LookupRepoID(opts[`in`][0]); err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/bucket/", repositoryID)
	return adm.Perform(`get`, path, `list`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
