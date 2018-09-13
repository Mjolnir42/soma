package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerClusters(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// clusters
			{
				Name:  "cluster",
				Usage: "SUBCOMMANDS for clusters",
				Subcommands: []cli.Command{
					{
						Name:         "create",
						Usage:        "Create a new cluster in a bucket",
						Action:       runtime(cmdClusterCreate),
						BashComplete: cmpl.In,
					},
					{
						Name:         "destroy",
						Usage:        "Destroy a cluster in a bucket",
						Action:       runtime(cmdClusterDestroy),
						BashComplete: cmpl.In,
					},
					{
						Name:         "rename",
						Usage:        "Rename a cluster",
						Action:       runtime(cmdClusterRename),
						BashComplete: cmpl.InTo,
					},
					{
						Name:         "list",
						Usage:        "List all clusters in a bucket",
						Action:       runtime(cmdClusterList),
						BashComplete: cmpl.In,
					},
					{
						Name:         "show",
						Usage:        "Show details about a cluster",
						Action:       runtime(cmdClusterShow),
						BashComplete: cmpl.In,
					},
					{
						Name:         "tree",
						Usage:        "Display the cluster as tree",
						Action:       runtime(cmdClusterTree),
						BashComplete: cmpl.In,
					},
					{
						Name:  "member",
						Usage: "SUBCOMMANDS for cluster members",
						Subcommands: []cli.Command{
							{
								Name:         "assign",
								Usage:        "Assign a node to a cluster",
								Action:       runtime(cmdClusterMemberAssign),
								BashComplete: cmpl.InTo,
							},
							{
								Name:         "unassign",
								Usage:        "Unassign a node from a cluster",
								Action:       runtime(cmdClusterMemberUnassign),
								BashComplete: cmpl.InFrom,
							},
							{
								Name:         "list",
								Usage:        "List members of a cluster",
								Action:       runtime(cmdClusterMemberList),
								BashComplete: cmpl.In,
							},
						},
					},
					{
						Name:  "property",
						Usage: "SUBCOMMANDS for properties",
						Subcommands: []cli.Command{
							{
								Name:  "add",
								Usage: "SUBCOMMANDS for property add",
								Subcommands: []cli.Command{
									{
										Name:         "system",
										Usage:        "Add a system property to a cluster",
										Action:       runtime(cmdClusterSystemPropertyAdd),
										BashComplete: cmpl.PropertyAddValue,
									},
									{
										Name:         "service",
										Usage:        "Add a service property to a cluster",
										Action:       runtime(cmdClusterServicePropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
									{
										Name:         `oncall`,
										Usage:        `Add an oncall property to a cluster`,
										Action:       runtime(cmdClusterOncallPropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
									{
										Name:         `custom`,
										Usage:        `Add a custom property to a cluster`,
										Action:       runtime(cmdClusterCustomPropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
								},
							},
							{
								Name:  `delete`,
								Usage: `SUBCOMMANDS for property delete`,
								Subcommands: []cli.Command{
									{
										Name:         `system`,
										Usage:        `Delete a system property from a cluster`,
										Action:       runtime(cmdClusterSystemPropertyDelete),
										BashComplete: cmpl.InFromView,
									},
									{
										Name:         `service`,
										Usage:        `Delete a service property from a cluster`,
										Action:       runtime(cmdClusterServicePropertyDelete),
										BashComplete: cmpl.InFromView,
									},
									{
										Name:         `oncall`,
										Usage:        `Delete an oncall property from a cluster`,
										Action:       runtime(cmdClusterOncallPropertyDelete),
										BashComplete: cmpl.InFromView,
									},
									{
										Name:         `custom`,
										Usage:        `Delete a custom property from a cluster`,
										Action:       runtime(cmdClusterCustomPropertyDelete),
										BashComplete: cmpl.InFromView,
									},
								},
							},
						},
					},
				},
			}, // end clusters
		}...,
	)
	return &app
}

func cmdClusterList(c *cli.Context) error {
	uniqKeys := []string{`in`}
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
		err                    error
		repositoryID, bucketID string
	)
	if bucketID, err = adm.LookupBucketID(opts["in"][0]); err != nil {
		return err
	}
	if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/bucket/%s/cluster/",
		repositoryID, bucketID)
	return adm.Perform(`get`, path, `list`, nil, c)
}

func cmdClusterShow(c *cli.Context) error {
	uniqKeys := []string{`in`}
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

	path := fmt.Sprintf("/repository/%s/bucket/%s/cluster/%s",
		repositoryID, bucketID, clusterID)
	return adm.Perform(`get`, path, `show`, nil, c)
}

func cmdClusterTree(c *cli.Context) error {
	uniqKeys := []string{`in`}
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
	)
	if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
		return err
	}
	if clusterID, err = adm.LookupClusterID(c.Args().First(),
		bucketID); err != nil {
		return err
	}
	if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/bucket/%s/cluster/%s/tree",
		repositoryID, bucketID, clusterID)
	return adm.Perform(`get`, path, `tree`, nil, c)
}

func cmdClusterCreate(c *cli.Context) error {
	uniqKeys := []string{`in`}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(opts,
		[]string{},
		uniqKeys,
		uniqKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	bucketID, err := adm.LookupBucketID(opts[`in`][0])
	if err != nil {
		return err
	}
	repositoryID, err := adm.LookupRepoByBucket(bucketID)
	if err != nil {
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

	path := fmt.Sprintf("/repository/%s/bucket/%s/cluster/", repositoryID, bucketID)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

func cmdClusterDestroy(c *cli.Context) error {
	uniqKeys := []string{`in`}
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
	)
	if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
		return err
	}
	if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}
	if clusterID, err = adm.LookupClusterID(c.Args().First(),
		bucketID); err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/bucket/%s/cluster/%s",
		repositoryID, bucketID, clusterID)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdClusterMemberList(c *cli.Context) error {
	uniqKeys := []string{`in`}
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
	)
	if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
		return err
	}
	if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}
	if clusterID, err = adm.LookupClusterID(c.Args().First(),
		bucketID); err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/bucket/%s/cluster/%s/member/",
		repositoryID, bucketID, clusterID)
	return adm.Perform(`get`, path, `list`, nil, c)
}

func cmdClusterMemberAssign(c *cli.Context) error {
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
		err                                       error
		nodeConfig                                *proto.NodeConfig
		nodeID, repositoryID, bucketID, clusterID string
	)
	if nodeID, err = adm.LookupNodeID(c.Args().First()); err != nil {
		return err
	}
	if nodeConfig, err = adm.LookupNodeConfig(nodeID); err != nil {
		return err
	}
	if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
		return err
	}
	if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}
	if repositoryID != nodeConfig.RepositoryID || bucketID != nodeConfig.BucketID {
		return fmt.Errorf(`Mismatching Repository or Bucket IDs`)
	}
	if clusterID, err = adm.LookupClusterID(opts[`to`][0],
		bucketID); err != nil {
		return err
	}

	req := proto.NewClusterRequest()
	node := proto.Node{
		ID:     nodeID,
		Config: nodeConfig,
	}
	req.Cluster = &proto.Cluster{
		ID:           clusterID,
		RepositoryID: repositoryID,
		BucketID:     bucketID,
		Members:      &[]proto.Node{},
	}
	*req.Cluster.Members = append(*req.Cluster.Members, node)

	path := fmt.Sprintf("/repository/%s/bucket/%s/cluster/%s/member/",
		repositoryID, bucketID, clusterID)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

func cmdClusterMemberUnassign(c *cli.Context) error {
	uniqKeys := []string{`from`, `in`}
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
		err                                       error
		nodeConfig                                *proto.NodeConfig
		nodeID, repositoryID, bucketID, clusterID string
	)
	if nodeID, err = adm.LookupNodeID(c.Args().First()); err != nil {
		return err
	}
	if nodeConfig, err = adm.LookupNodeConfig(nodeID); err != nil {
		return err
	}
	if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
		return err
	}
	if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}
	if repositoryID != nodeConfig.RepositoryID || bucketID != nodeConfig.BucketID {
		return fmt.Errorf(`Mismatching Repository or Bucket IDs`)
	}
	if clusterID, err = adm.LookupClusterID(opts[`from`][0],
		bucketID); err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/bucket/%s/cluster/%s/member/%s/%s",
		repositoryID, bucketID, clusterID, `node`, nodeID)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdClusterSystemPropertyAdd(c *cli.Context) error {
	return cmdClusterPropertyAdd(c, `system`)
}

func cmdClusterServicePropertyAdd(c *cli.Context) error {
	return cmdClusterPropertyAdd(c, `service`)
}

func cmdClusterOncallPropertyAdd(c *cli.Context) error {
	return cmdClusterPropertyAdd(c, `oncall`)
}

func cmdClusterCustomPropertyAdd(c *cli.Context) error {
	return cmdClusterPropertyAdd(c, `custom`)
}

func cmdClusterPropertyAdd(c *cli.Context, pType string) error {
	return cmdPropertyAdd(c, pType, `cluster`)
}

func cmdClusterSystemPropertyDelete(c *cli.Context) error {
	return cmdClusterPropertyDelete(c, `system`)
}

func cmdClusterServicePropertyDelete(c *cli.Context) error {
	return cmdClusterPropertyDelete(c, `service`)
}

func cmdClusterOncallPropertyDelete(c *cli.Context) error {
	return cmdClusterPropertyDelete(c, `oncall`)
}

func cmdClusterCustomPropertyDelete(c *cli.Context) error {
	return cmdClusterPropertyDelete(c, `custom`)
}

func cmdClusterPropertyDelete(c *cli.Context, pType string) error {
	multiple := []string{}
	unique := []string{`from`, `view`, `in`}
	required := []string{`from`, `view`, `in`}
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		multiple,
		unique,
		required,
		c.Args().Tail()); err != nil {
		return err
	}
	var (
		err                                         error
		repositoryID, bucketID, clusterID, sourceID string
	)
	if bucketID, err = adm.LookupBucketID(opts["in"][0]); err != nil {
		return err
	}
	if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}
	if clusterID, err = adm.LookupClusterID(opts[`from`][0],
		bucketID); err != nil {
		return err
	}

	if pType == `system` {
		if err := adm.ValidateSystemProperty(
			c.Args().First()); err != nil {
			return err
		}
	}
	if err := adm.FindClusterPropSrcID(pType, c.Args().First(),
		opts[`view`][0], clusterID, &sourceID); err != nil {
		return err
	}

	req := proto.NewClusterRequest()
	req.Cluster.ID = clusterID
	req.Cluster.BucketID = bucketID

	path := fmt.Sprintf("/repository/%s/bucket/%s/cluster/%s/property/%s/%s",
		repositoryID, bucketID, clusterID, pType, sourceID)
	return adm.Perform(`deletebody`, path, `command`, req, c)
}

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
