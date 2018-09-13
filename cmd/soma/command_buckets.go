package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/internal/help"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerBuckets(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// buckets
			{
				Name:  "bucket",
				Usage: "SUBCOMMANDS for buckets",
				Subcommands: []cli.Command{
					{
						Name:         "create",
						Usage:        "Create a new bucket inside a repository",
						Action:       runtime(cmdBucketCreate),
						BashComplete: cmpl.BucketCreate,
					},
					{
						Name:         "destroy",
						Usage:        "Mark an existing bucket as deleted",
						Action:       runtime(cmdBucketDestroy),
						BashComplete: cmpl.In,
					},
					{
						Name:         "rename",
						Usage:        "Rename an existing bucket",
						Action:       runtime(cmdBucketRename),
						BashComplete: cmpl.BucketRename,
					},
					{
						Name:         "list",
						Usage:        "List existing buckets",
						Action:       runtime(cmdBucketList),
						BashComplete: cmpl.In,
					},
					{
						Name:   "show",
						Usage:  "Show information about a specific bucket",
						Action: runtime(cmdBucketShow),
					},
					{
						Name:   `tree`,
						Usage:  `Display the bucket as tree`,
						Action: runtime(cmdBucketTree),
					},
					{
						Name:   `instances`,
						Usage:  `List check instances for a bucket`,
						Action: runtime(cmdBucketInstance),
					},
					{
						Name:  "property",
						Usage: "SUBCOMMANDS for properties",
						Subcommands: []cli.Command{
							{
								Name:        `create`,
								Usage:       `SUBCOMMANDS for property create`,
								Description: help.Text(`BucketsPropertyCreate`),
								Subcommands: []cli.Command{
									{
										Name:         "system",
										Usage:        "Add a system property to a bucket",
										Action:       runtime(cmdBucketSystemPropertyAdd),
										BashComplete: cmpl.PropertyAddValue,
									},
									{
										Name:         "service",
										Usage:        "Add a service property to a bucket",
										Action:       runtime(cmdBucketServicePropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
									{
										Name:         "oncall",
										Usage:        "Add an oncall property to a bucket",
										Action:       runtime(cmdBucketOncallPropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
									{
										Name:         "custom",
										Usage:        "Add a custom property to a bucket",
										Action:       runtime(cmdBucketCustomPropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
								},
							},
							{
								Name:  `destroy`,
								Usage: `SUBCOMMANDS for property destroy`,
								Subcommands: []cli.Command{
									{
										Name:         `system`,
										Usage:        `Delete a system property from a bucket`,
										Action:       runtime(cmdBucketSystemPropertyDelete),
										BashComplete: cmpl.FromView,
									},
									{
										Name:         `service`,
										Usage:        `Delete a service property from a bucket`,
										Action:       runtime(cmdBucketServicePropertyDelete),
										BashComplete: cmpl.FromView,
									},
									{
										Name:         `oncall`,
										Usage:        `Delete an oncall property from a bucket`,
										Action:       runtime(cmdBucketOncallPropertyDelete),
										BashComplete: cmpl.FromView,
									},
									{
										Name:         `custom`,
										Usage:        `Delete a custom property from a bucket`,
										Action:       runtime(cmdBucketCustomPropertyDelete),
										BashComplete: cmpl.FromView,
									},
								},
							},
						},
					},
				},
			}, // end buckets
		}...,
	)
	return &app
}

func cmdBucketList(c *cli.Context) error {
	var err error
	var repositoryID, path string

	if err = adm.VerifyNoArgument(c); err != nil {
		path = `/bucket/`
	} else {
		uniqKeys := []string{`in`}
		opts := map[string][]string{}

		if err = adm.ParseVariadicArguments(
			opts,
			[]string{},
			uniqKeys,
			uniqKeys,
			c.Args().Tail()); err != nil {
			return err
		}

		if repositoryID, err = adm.LookupRepoByBucket(opts[`in`][0]); err != nil {
			return err
		}

		path = fmt.Sprintf("/repository/%s/bucket/", repositoryID)
	}

	return adm.Perform(`get`, path, `list`, nil, c)
}

func cmdBucketShow(c *cli.Context) error {
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

	path := fmt.Sprintf("/repository/%s/bucket/%s",
		repositoryID, bucketID)
	return adm.Perform(`get`, path, `show`, nil, c)
}

func cmdBucketTree(c *cli.Context) error {
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

	path := fmt.Sprintf("/repository/%s/bucket/%s/tree",
		repositoryID, bucketID)
	return adm.Perform(`get`, path, `tree`, nil, c)
}

func cmdBucketCreate(c *cli.Context) error {
	var err error
	var repositoryID string

	if err = adm.ValidateRuneCountRange(c.Args().First(), 4, 512); err != nil {
		return err
	}

	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`in`, `environment`}
	mandatoryOptions := []string{`in`, `environment`}

	if err = adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail()); err != nil {
		return err
	}

	if repositoryID, err = adm.LookupRepoID(opts[`in`][0]); err != nil {
		return err
	}

	// fetch list of environments from SOMA to check if a valid
	// environment was requested
	if err = adm.ValidateEnvironment(opts[`environment`][0]); err != nil {
		return err
	}

	req := proto.NewBucketRequest()
	req.Bucket.Name = c.Args().First()
	req.Bucket.RepositoryID = repositoryID
	req.Bucket.Environment = opts[`environment`][0]

	path := fmt.Sprintf("/repository/%s/bucket/", repositoryID)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

func cmdBucketDestroy(c *cli.Context) error {
	var err error
	var repositoryID, controlID, bucketID string

	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`in`}
	mandatoryOptions := []string{}

	if err = adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail()); err != nil {
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

	path := fmt.Sprintf("/repository/%s/bucket/%s", repositoryID, bucketID)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

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

///

func cmdBucketSystemPropertyAdd(c *cli.Context) error {
	return cmdBucketPropertyAdd(c, `system`)
}

func cmdBucketServicePropertyAdd(c *cli.Context) error {
	return cmdBucketPropertyAdd(c, `service`)
}

func cmdBucketOncallPropertyAdd(c *cli.Context) error {
	return cmdBucketPropertyAdd(c, `oncall`)
}

func cmdBucketCustomPropertyAdd(c *cli.Context) error {
	return cmdBucketPropertyAdd(c, `custom`)
}

func cmdBucketPropertyAdd(c *cli.Context, pType string) error {
	return cmdPropertyAdd(c, pType, `bucket`)
}

func cmdBucketSystemPropertyDelete(c *cli.Context) error {
	return cmdBucketPropertyDelete(c, `system`)
}

func cmdBucketServicePropertyDelete(c *cli.Context) error {
	return cmdBucketPropertyDelete(c, `service`)
}

func cmdBucketOncallPropertyDelete(c *cli.Context) error {
	return cmdBucketPropertyDelete(c, `oncall`)
}

func cmdBucketCustomPropertyDelete(c *cli.Context) error {
	return cmdBucketPropertyDelete(c, `custom`)
}

func cmdBucketPropertyDelete(c *cli.Context, pType string) error {
	unique := []string{`from`, `view`}
	required := []string{`from`, `view`}
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
	bucketID, err := adm.LookupBucketID(opts[`from`][0])
	if err != nil {
		return err
	}

	if pType == `system` {
		if err := adm.ValidateSystemProperty(
			c.Args().First()); err != nil {
			return err
		}
	}
	var sourceID string
	if err := adm.FindBucketPropSrcID(pType, c.Args().First(),
		opts[`view`][0], bucketID, &sourceID); err != nil {
		return err
	}

	path := fmt.Sprintf("/bucket/%s/property/%s/%s",
		bucketID, pType, sourceID)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
