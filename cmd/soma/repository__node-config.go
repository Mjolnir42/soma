/*-
 * Copyright (c) 2019, Jörg Pernfuß
 * Copyright (c) 2019, 1&1 IONOS SE
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

// nodeConfigTree function
// soma node dumptree ${node} [in ${bucket}]
func nodeConfigTree(c *cli.Context) error {
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
	var nodeID string
	config := &proto.NodeConfig{}

	if nodeID, err = adm.LookupNodeID(c.Args().First()); err != nil {
		return err
	}
	if config, err = adm.LookupNodeConfig(nodeID); err != nil {
		return err
	}

	// optional argument, must be correct if given
	if _, ok := opts[`in`]; ok {
		if bucketID, err := adm.LookupBucketID(opts[`in`][0]); err != nil {
			return err
		} else if bucketID != config.BucketID {
			return fmt.Errorf("Invalid request: node %s is in bucket %s, not %s",
				c.Args().First(),
				config.BucketID,
				bucketID,
			)
		}
	}

	path := fmt.Sprintf("/repository/%s/bucket/%s/node/%s/tree",
		url.QueryEscape(config.RepositoryID),
		url.QueryEscape(config.BucketID),
		url.QueryEscape(nodeID),
	)
	return adm.Perform(`get`, path, `node-config::tree`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
