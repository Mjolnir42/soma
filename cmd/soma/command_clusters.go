package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/internal/help"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerClusters(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:        `cluster`,
				Usage:       `SUBCOMMANDS for bucket cluster management`,
				Description: help.Text(`cluster-config::`),
				Subcommands: []cli.Command{
					{
						Name:         `list`,
						Usage:        `List all clusters in a bucket`,
						Description:  help.Text(`cluster-config::list`),
						Action:       runtime(clusterConfigList),
						BashComplete: cmpl.DirectIn,
					},
					{
						Name:         `show`,
						Usage:        `Show details about a cluster`,
						Description:  help.Text(`cluster-config::show`),
						Action:       runtime(clusterConfigShow),
						BashComplete: cmpl.In,
					},
					{
						Name:         `create`,
						Usage:        `Create a new cluster in a bucket`,
						Description:  help.Text(`cluster-config::create`),
						Action:       runtime(clusterConfigCreate),
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
						Name:        `property`,
						Usage:       `SUBCOMMANDS for properties on clusters`,
						Description: help.Text(`cluster-config::`),
						Subcommands: []cli.Command{
							{
								Name:        `create`,
								Usage:       `SUBCOMMANDS to create properties`,
								Description: help.Text(`cluster-config::property-create`),
								Subcommands: []cli.Command{
									{
										Name:         `system`,
										Usage:        `Add a system property to a cluster`,
										Description:  help.Text(`cluster-config::property-create`),
										Action:       runtime(clusterConfigPropertyCreateSystem),
										BashComplete: cmpl.PropertyCreateInValue,
									},
									{
										Name:         `service`,
										Usage:        `Add a service property to a cluster`,
										Description:  help.Text(`cluster-config::property-create`),
										Action:       runtime(clusterConfigPropertyCreateService),
										BashComplete: cmpl.PropertyCreateInValue,
									},
									{
										Name:         `oncall`,
										Usage:        `Add an oncall property to a cluster`,
										Description:  help.Text(`cluster-config::property-create`),
										Action:       runtime(clusterConfigPropertyCreateOncall),
										BashComplete: cmpl.PropertyCreateIn,
									},
									{
										Name:         `custom`,
										Usage:        `Add a custom property to a cluster`,
										Description:  help.Text(`cluster-config::property-create`),
										Action:       runtime(clusterConfigPropertyCreateCustom),
										BashComplete: cmpl.PropertyCreateIn,
									},
								},
							},
							{
								Name:        `destroy`,
								Usage:       `SUBCOMMANDS to destroy properties`,
								Description: help.Text(`cluster-config::property-destroy`),
								Subcommands: []cli.Command{
									{
										Name:         `system`,
										Usage:        `Delete a system property from a cluster`,
										Description:  help.Text(`cluster-config::property-destroy`),
										Action:       runtime(clusterConfigPropertyDestroySystem),
										BashComplete: cmpl.PropertyOnInView,
									},
									{
										Name:         `service`,
										Usage:        `Delete a service property from a cluster`,
										Description:  help.Text(`cluster-config::property-destroy`),
										Action:       runtime(clusterConfigPropertyDestroyService),
										BashComplete: cmpl.PropertyOnInView,
									},
									{
										Name:         `oncall`,
										Usage:        `Delete an oncall property from a cluster`,
										Description:  help.Text(`cluster-config::property-destroy`),
										Action:       runtime(clusterConfigPropertyDestroyOncall),
										BashComplete: cmpl.PropertyOnInView,
									},
									{
										Name:         `custom`,
										Usage:        `Delete a custom property from a cluster`,
										Description:  help.Text(`cluster-config::property-destroy`),
										Action:       runtime(clusterConfigPropertyDestroyCustom),
										BashComplete: cmpl.PropertyOnInView,
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
