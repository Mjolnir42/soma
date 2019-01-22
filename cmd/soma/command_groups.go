package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/lib/proto"
)

func cmdGroupMemberAddCluster(c *cli.Context) error {
	uniqKeys := []string{"to", "in"}
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
		err                           error
		bucketID, mClusterID, groupID string
		req                           proto.Request
		cluster                       proto.Cluster
	)
	if bucketID, err = adm.LookupBucketID(
		opts["in"][0]); err != nil {
		return err
	}
	if mClusterID, err = adm.LookupGroupID(c.Args().First(),
		bucketID); err != nil {
		return err
	}
	if groupID, err = adm.LookupClusterID(opts["to"][0],
		bucketID); err != nil {
		return err
	}

	cluster.ID = mClusterID
	req.Group = &proto.Group{}
	req.Group.ID = groupID
	req.Group.BucketID = bucketID
	*req.Group.MemberClusters = append(
		*req.Group.MemberClusters, cluster)

	path := fmt.Sprintf("/groups/%s/members/", groupID)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

func cmdGroupMemberAddNode(c *cli.Context) error {
	uniqKeys := []string{"to", "in"}
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
		err                        error
		bucketID, groupID, mNodeID string
		req                        proto.Request
		node                       proto.Node
	)
	if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
		return err
	}
	if mNodeID, err = adm.LookupNodeID(c.Args().First()); err != nil {
		return err
	}
	if groupID, err = adm.LookupGroupID(opts[`to`][0],
		bucketID); err != nil {
		return err
	}

	node.ID = mNodeID
	req.Group = &proto.Group{}
	req.Group.ID = groupID
	req.Group.BucketID = bucketID
	*req.Group.MemberNodes = append(*req.Group.MemberNodes, node)

	path := fmt.Sprintf("/groups/%s/members/", groupID)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

func cmdGroupMemberDeleteGroup(c *cli.Context) error {
	uniqKeys := []string{"from", "in"}
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
		err                         error
		bucketID, mGroupID, groupID string
	)
	if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
		return err
	}
	if mGroupID, err = adm.LookupGroupID(c.Args().First(),
		bucketID); err != nil {
		return err
	}
	if groupID, err = adm.LookupGroupID(opts[`from`][0],
		bucketID); err != nil {
		return err
	}

	path := fmt.Sprintf("/groups/%s/members/%s", groupID,
		mGroupID)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdGroupMemberDeleteCluster(c *cli.Context) error {
	uniqKeys := []string{"from", "in"}
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
		err                           error
		bucketID, mClusterID, groupID string
	)
	if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
		return err
	}
	if mClusterID, err = adm.LookupClusterID(c.Args().First(),
		bucketID); err != nil {
		return err
	}
	if groupID, err = adm.LookupGroupID(opts[`from`][0],
		bucketID); err != nil {
		return err
	}

	path := fmt.Sprintf("/groups/%s/members/%s", groupID,
		mClusterID)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdGroupMemberDeleteNode(c *cli.Context) error {
	uniqKeys := []string{"from", "in"}
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
		err                        error
		bucketID, groupID, mNodeID string
	)
	if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
		return err
	}
	if mNodeID, err = adm.LookupNodeID(c.Args().First()); err != nil {
		return err
	}
	if groupID, err = adm.LookupGroupID(opts[`from`][0],
		bucketID); err != nil {
		return err
	}

	path := fmt.Sprintf("/groups/%s/members/%s", groupID,
		mNodeID)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdGroupMemberList(c *cli.Context) error {
	uniqKeys := []string{"in"}
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
		err               error
		bucketID, groupID string
	)
	if bucketID, err = adm.LookupBucketID(opts["in"][0]); err != nil {
		return err
	}
	if groupID, err = adm.LookupGroupID(c.Args().First(),
		bucketID); err != nil {
		return err
	}

	path := fmt.Sprintf("/groups/%s/members/", groupID)
	return adm.Perform(`get`, path, `list`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
