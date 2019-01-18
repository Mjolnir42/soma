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
	"github.com/mjolnir42/soma/lib/proto"
)

// clusterConfigList function
// soma cluster list in ${bucket}
func clusterConfigList(c *cli.Context) error {
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
		"/repository/%s/bucket/%s/cluster/",
		url.QueryEscape(repositoryID),
		url.QueryEscape(bucketID),
	)
	return adm.Perform(`get`, path, `cluster-config::list`, nil, c)
}

// clusterConfigShow function
// soma cluster show ${cluster} in ${bucket}
func clusterConfigShow(c *cli.Context) error {
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
	var repositoryID, bucketID, clusterID string
	if bucketID, err = adm.LookupBucketID(opts["in"][0]); err != nil {
		return err
	}
	if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}
	if clusterID, err = adm.LookupClusterID(c.Args().First(),
		bucketID); err != nil {
		return err
	}

	path := fmt.Sprintf(
		"/repository/%s/bucket/%s/cluster/%s",
		url.QueryEscape(repositoryID),
		url.QueryEscape(bucketID),
		url.QueryEscape(clusterID),
	)
	return adm.Perform(`get`, path, `cluster-config::show`, nil, c)
}

// clusterConfigCreate function
// soma cluster create ${cluster} in ${bucket}
func clusterConfigCreate(c *cli.Context) error {
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

	req := proto.NewClusterRequest()
	req.Cluster.Name = c.Args().First()
	req.Cluster.RepositoryID = repositoryID
	req.Cluster.BucketID = bucketID

	if err := adm.ValidateRuneCountRange(
		req.Cluster.Name, 4, 256); err != nil {
		return err
	}

	path := fmt.Sprintf(
		"/repository/%s/bucket/%s/cluster/",
		url.QueryEscape(repositoryID),
		url.QueryEscape(bucketID),
	)
	return adm.Perform(`postbody`, path, `cluster-config::create`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
