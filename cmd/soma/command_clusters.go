package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/lib/proto"
)

func cmdClusterRename(c *cli.Context) error {
	uniqKeys := []string{`to`, `in`}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		uniqKeys,
		uniqKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	var (
		err                               error
		repositoryID, bucketID, clusterID string
		req                               proto.Request
	)
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

	req.Cluster = &proto.Cluster{}
	req.Cluster.Name = opts[`to`][0]

	path := fmt.Sprintf("/repository/%s/bucket/%s/cluster/%s",
		repositoryID, bucketID, clusterID)
	return adm.Perform(`patchbody`, path, `command`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
