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

// nodeList function
// soma node list
func nodeList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/node/`, `node::list`, nil, c)
}

// nodeShow function
// soma node show ${node}
func nodeShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	// check deferred errors
	if err := popError(); err != nil {
		return err
	}

	nodeID, err := adm.LookupNodeID(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/node/%s", url.QueryEscape(nodeID))
	return adm.Perform(`get`, path, `node::show`, nil, c)
}

// nodeShowConfig function
// soma node config ${node}
func nodeShowConfig(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	// check deferred errors
	if err := popError(); err != nil {
		return err
	}

	nodeID, err := adm.LookupNodeID(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/node/%s/config", url.QueryEscape(nodeID))
	return adm.Perform(`get`, path, `node::show-config`, nil, c)
}

// nodeAssign function
// soma node assign ${node} to ${bucket}
func nodeAssign(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.VariadicArguments(`node::assign`, c, &opts); err != nil {
		return err
	}

	// check deferred errors
	if err := popError(); err != nil {
		return err
	}

	var (
		err                              error
		bucketID, repoID, nodeID         string
		teamIDFromBucket, teamIDFromNode string
	)
	if bucketID, err = adm.LookupBucketID(opts[`to`][0]); err != nil {
		return err
	}
	if repoID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}
	if nodeID, err = adm.LookupNodeID(c.Args().First()); err != nil {
		return err
	}
	if teamIDFromBucket, err = adm.LookupTeamByBucket(bucketID); err != nil {
		return err
	}
	if teamIDFromNode, err = adm.LookupTeamByNode(nodeID); err != nil {
		return err
	}
	if teamIDFromBucket != teamIDFromNode {
		return fmt.Errorf(
			`Cannot assign node since node and bucket belong to` +
				` different teams.`)
	}

	req := proto.NewNodeRequest()
	req.Node.ID = nodeID
	req.Node.Config = &proto.NodeConfig{}
	req.Node.Config.RepositoryID = repoID
	req.Node.Config.BucketID = bucketID
	req.Node.TeamID = teamIDFromNode

	path := fmt.Sprintf("/node/%s/config",
		url.QueryEscape(nodeID),
	)
	return adm.Perform(`putbody`, path, `node::assign`, req, c)
}

// nodeUnassign function
// soma node unassign ${node} [from ${bucket}]
func nodeUnassign(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.VariadicArguments(`node::unassign`, c, &opts); err != nil {
		return err
	}

	// check deferred errors
	if err := popError(); err != nil {
		return err
	}

	var (
		err    error
		nodeID string
	)
	config := &proto.NodeConfig{}
	if nodeID, err = adm.LookupNodeID(c.Args().First()); err != nil {
		return err
	}
	if config, err = adm.LookupNodeConfig(nodeID); err != nil {
		return err
	}
	// optional argument must be correct if given
	if _, ok := opts[`from`]; ok {
		if bucketID, err := adm.LookupBucketID(opts[`from`][0]); err != nil {
			return err
		} else if bucketID != config.BucketID {
			return fmt.Errorf("Invalid request: node %s is in bucket %s, not %s",
				c.Args().First(),
				config.BucketID,
				bucketID,
			)
		}
	}

	path := fmt.Sprintf("/repository/%s/bucket/%s/node/%s/config",
		url.QueryEscape(config.RepositoryID),
		url.QueryEscape(config.BucketID),
		url.QueryEscape(nodeID),
	)
	return adm.Perform(`delete`, path, `node::unassign`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
