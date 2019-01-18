package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/lib/proto"
)

func cmdBucketRename(c *cli.Context) error {
	var err error
	var repositoryID, controlID, bucketID string

	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`in`, `to`}
	mandatoryOptions := []string{`to`}

	if err = adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail()); err != nil {
		return err
	}

	if err = adm.ValidateRuneCountRange(opts[`to`][0], 4, 512); err != nil {
		return err
	}
	if bucketID, err = adm.LookupBucketID(c.Args().First()); err != nil {
		return err
	}
	if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}
	if _, ok := opts[`in`]; ok {
		if controlID, err = adm.LookupRepoID(opts[`in`][0]); err != nil {
			return err
		} else if controlID != repositoryID {
			return fmt.Errorf("bucket %s is not in repository %s", c.Args().First(), opts[`in`][0])
		}
	}

	req := proto.NewBucketRequest()
	req.Bucket.ID = bucketID
	req.Bucket.RepositoryID = repositoryID
	req.Bucket.Name = opts[`to`][0]

	path := fmt.Sprintf("/repository/%s/bucket/%s", repositoryID, bucketID)
	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdBucketInstance(c *cli.Context) error {
	var err error
	var repositoryID, bucketID string

	if err = adm.VerifySingleArgument(c); err != nil {
		return err
	}
	if bucketID, err = adm.LookupBucketID(c.Args().First()); err != nil {
		return err
	}
	if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/bucket/%s/instance/",
		repositoryID, bucketID)
	return adm.Perform(`get`, path, `list`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
