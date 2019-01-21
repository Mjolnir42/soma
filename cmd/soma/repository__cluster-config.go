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
						Name:         `dumptree`,
						Usage:        `Display the cluster as tree`,
						Description:  help.Text(`cluster-config::tree`),
						Action:       runtime(clusterConfigTree),
						BashComplete: cmpl.In,
					},
					{
						Name:         `destroy`,
						Usage:        `Destroy a cluster`,
						Description:  help.Text(`cluster-config::destroy`),
						Action:       runtime(clusterConfigDestroy),
						BashComplete: cmpl.In,
					},
					{
						Name:        `member`,
						Usage:       `SUBCOMMANDS for cluster membership management`,
						Description: help.Text(`cluster-config::`),
						Subcommands: []cli.Command{
							{
								Name:         `assign`,
								Usage:        `Assign a node to a cluster`,
								Description:  help.Text(`cluster-config::member-assign`),
								Action:       runtime(clusterConfigMemberAssign),
								BashComplete: cmpl.InTo,
							},
							{
								Name:         `unassign`,
								Usage:        `Unassign a node from a cluster`,
								Description:  help.Text(`cluster-config::member-unassign`),
								Action:       runtime(clusterConfigMemberUnassign),
								BashComplete: cmpl.InFrom,
							},
							{
								Name:         `list`,
								Usage:        `List member nodes of a cluster`,
								Description:  help.Text(`cluster-config::member-list`),
								Action:       runtime(clusterConfigMemberList),
								BashComplete: cmpl.DirectInOf,
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

// clusterConfigDestroy function
// soma cluster destroy ${cluster} in ${bucket}
func clusterConfigDestroy(c *cli.Context) error {
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
	if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
		return err
	}
	if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}
	if clusterID, err = adm.LookupClusterID(c.Args().First(), bucketID); err != nil {
		return err
	}

	path := fmt.Sprintf(
		"/repository/%s/bucket/%s/cluster/%s",
		url.QueryEscape(repositoryID),
		url.QueryEscape(bucketID),
		url.QueryEscape(clusterID),
	)
	return adm.Perform(`delete`, path, `cluster-config::destroy`, nil, c)
}

// clusterConfigTree function
// soma cluster dumptree ${cluster} in ${bucket}
func clusterConfigTree(c *cli.Context) error {
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

	path := fmt.Sprintf(
		"/repository/%s/bucket/%s/cluster/%s/tree",
		url.QueryEscape(repositoryID),
		url.QueryEscape(bucketID),
		url.QueryEscape(clusterID),
	)
	return adm.Perform(`get`, path, `cluster-config::tree`, nil, c)
}

// clusterConfigMemberAssign function
// soma cluster member assign ${node} to ${cluster} [in ${bucket}]
func clusterConfigMemberAssign(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`to`, `in`}
	mandatoryOptions := []string{`to`}

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
	var nodeID, clusterID, bucketID, repositoryID, controlID string
	nodeConfig := &proto.NodeConfig{}

	if nodeID, err = adm.LookupNodeID(c.Args().First()); err != nil {
		return err
	}
	if nodeConfig, err = adm.LookupNodeConfig(nodeID); err != nil {
		// unassigned node can not be assigned to a cluster, this is
		// currently intended
		return err
	}
	bucketID = nodeConfig.BucketID
	repositoryID = nodeConfig.RepositoryID

	// optional argument must be correct if provided
	if _, ok := opts[`in`]; ok {
		if controlID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
			return err
		} else if controlID != bucketID {
			return fmt.Errorf(
				"Invalid bucket %s(%s), node %s is assigned to %s",
				opts[`in`][0],
				controlID,
				c.Args().First(),
				bucketID,
			)
		}
	}

	if clusterID, err = adm.LookupClusterID(opts[`to`][0], bucketID); err != nil {
		return err
	}
	req := proto.NewClusterRequest()
	req.Cluster.ID = clusterID
	req.Cluster.RepositoryID = repositoryID
	req.Cluster.BucketID = bucketID
	req.Cluster.Members = &[]proto.Node{proto.Node{
		ID:     nodeID,
		Config: nodeConfig,
	}}

	path := fmt.Sprintf("/repository/%s/bucket/%s/cluster/%s/member/",
		url.QueryEscape(repositoryID),
		url.QueryEscape(bucketID),
		url.QueryEscape(clusterID),
	)
	return adm.Perform(`postbody`, path, `cluster-config::member-assign`, req, c)
}

// clusterConfigMemberUnassign function
// soma cluster member unassign ${node} from ${cluster} [in ${bucket}]
func clusterConfigMemberUnassign(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`from`, `in`}
	mandatoryOptions := []string{`from`}

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
	var nodeID, clusterID, bucketID, repositoryID, controlID string
	nodeConfig := &proto.NodeConfig{}

	if nodeID, err = adm.LookupNodeID(c.Args().First()); err != nil {
		return err
	}
	if nodeConfig, err = adm.LookupNodeConfig(nodeID); err != nil {
		// unassigned node can not be assigned to a cluster
		return err
	}
	bucketID = nodeConfig.BucketID
	repositoryID = nodeConfig.RepositoryID

	// optional argument must be correct if provided
	if _, ok := opts[`in`]; ok {
		if controlID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
			return err
		} else if controlID != bucketID {
			return fmt.Errorf(
				"Invalid bucket %s(%s), node %s is assigned to %s",
				opts[`in`][0],
				controlID,
				c.Args().First(),
				bucketID,
			)
		}
	}

	if clusterID, err = adm.LookupClusterID(opts[`from`][0], bucketID); err != nil {
		return err
	}
	req := proto.NewClusterRequest()
	req.Cluster.ID = clusterID
	req.Cluster.RepositoryID = repositoryID
	req.Cluster.BucketID = bucketID
	req.Cluster.Members = &[]proto.Node{proto.Node{
		ID:     nodeID,
		Config: nodeConfig,
	}}

	path := fmt.Sprintf("/repository/%s/bucket/%s/cluster/%s/member/%s/%s",
		url.QueryEscape(repositoryID),
		url.QueryEscape(bucketID),
		url.QueryEscape(clusterID),
		url.QueryEscape(proto.EntityNode),
		url.QueryEscape(nodeID),
	)
	return adm.Perform(`delete`, path, `cluster-config::member-unassign`, nil, c)
}

// clusterConfigMemberList function
// soma cluster member list of ${cluster} in ${bucket}
func clusterConfigMemberList(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`of`, `in`}
	mandatoryOptions := []string{`of`, `in`}

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
	var clusterID, bucketID, repositoryID string
	if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
		return err
	}
	if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}
	if clusterID, err = adm.LookupClusterID(opts[`of`][0], bucketID); err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/bucket/%s/cluster/%s/member/",
		url.QueryEscape(repositoryID),
		url.QueryEscape(bucketID),
		url.QueryEscape(clusterID),
	)
	return adm.Perform(`get`, path, `cluster-config::member-list`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
