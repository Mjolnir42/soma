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
	return adm.Perform(`postbody`, path, `bucket::create`, req, c)
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
	return adm.Perform(`delete`, path, `bucket::destroy`, nil, c)
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
	return adm.Perform(`get`, path, `bucket::list`, nil, c)
}

// bucketShow function
// soma bucket show ${bucket} [in ${repository}]
func bucketShow(c *cli.Context) error {
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
	var repositoryID, repositoryControlID, bucketID string

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
	return adm.Perform(`get`, path, `bucket::show`, nil, c)
}

// bucketSearch function
// soma bucket search [id ${uuid}] [name ${bucket}] [repository ${repository}] [environment ${environment}] [deleted ${isDeleted}]
func bucketSearch(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`id`, `name`, `repository`, `environment`, `deleted`}
	mandatoryOptions := []string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		adm.AllArguments(c),
	); err != nil {
		return err
	}

	validCondition := false
	req := proto.NewBucketFilter()

	if _, ok := opts[`id`]; ok {
		req.Filter.Bucket.ID = opts[`id`][0]
		if err := adm.ValidateUUID(req.Filter.Bucket.ID); err != nil {
			return err
		}
		validCondition = true
	}

	if _, ok := opts[`name`]; ok {
		req.Filter.Bucket.Name = opts[`name`][0]
		if err := adm.ValidateNotUUID(req.Filter.Bucket.Name); err != nil {
			return err
		}
		validCondition = true
	}

	if _, ok := opts[`repository`]; ok {
		repositoryID, err := adm.LookupRepoID(opts[`repository`][0])
		if err != nil {
			return err
		}
		req.Filter.Bucket.RepositoryID = repositoryID
		if err := adm.ValidateUUID(req.Filter.Bucket.RepositoryID); err != nil {
			return err
		}
		validCondition = true
	}

	if _, ok := opts[`environment`]; ok {
		req.Filter.Bucket.Environment = opts[`environment`][0]
		if err := adm.ValidateNotUUID(req.Filter.Bucket.Environment); err != nil {
			return err
		}
		if err := adm.ValidateNoSlash(req.Filter.Bucket.Environment); err != nil {
			return err
		}
		validCondition = true
	}

	if _, ok := opts[`deleted`]; ok {
		if err := adm.ValidateBool(opts[`deleted`][0],
			&req.Filter.Bucket.IsDeleted,
		); err != nil {
			return err
		}
		req.Filter.Bucket.FilterOnIsDeleted = true
		validCondition = true
	}

	if !validCondition {
		return fmt.Errorf(`Syntax error: at least one search condition must be specified`)
	}

	return adm.Perform(`postbody`, `/search/bucket/`, `list`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
