/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"database/sql"
	"fmt"

	"github.com/mjolnir42/soma/internal/msg"
)

// Extract the request routing information
func (g *GuidePost) extractRouting(q *msg.Request) (string, string, bool, error) {
	var repoID, repoName, bucketID string
	var err error

	repoID, bucketID = g.extractID(q)

	// lookup repository by bucket
	if bucketID != `` {
		if err = g.stmtRepoForBucketID.QueryRow(
			bucketID,
		).Scan(
			&repoID,
			&repoName,
		); err != nil {
			if err == sql.ErrNoRows {
				return ``, ``, true, fmt.Errorf(
					"No repository found for bucketID %s",
					bucketID,
				)
			}
			return ``, ``, false, err
		}
	}

	// lookup repository name
	if repoName == `` && repoID != `` {
		if err = g.stmtRepoNameByID.QueryRow(
			repoID,
		).Scan(
			&repoName,
		); err != nil {
			if err == sql.ErrNoRows {
				return ``, ``, true, fmt.Errorf(
					"No repository found with id %s",
					repoID,
				)
			}
			return ``, ``, false, err
		}
	}

	if repoName == `` {
		return ``, ``, true, fmt.Errorf(
			`GuidePost: unable find repository for request`,
		)
	}
	return repoID, repoName, false, nil
}

// Extract embedded IDs that can be used for routing
func (g *GuidePost) extractID(q *msg.Request) (string, string) {
	switch q.Section {
	case msg.SectionNodeConfig:
		switch q.Action {
		case msg.ActionAssign:
		case msg.ActionUnassign:
			return q.Repository.ID, q.Bucket.ID
		case msg.ActionPropertyCreate:
		case msg.ActionPropertyUpdate:
		case msg.ActionPropertyDestroy:
		default:
			return ``, ``
		}
		return q.Node.Config.RepositoryID, q.Node.Config.BucketID
	case msg.SectionCheckConfig:
		switch q.Action {
		case msg.ActionCreate:
		case msg.ActionDestroy:
		default:
			return ``, ``
		}
		return q.CheckConfig.RepositoryID, ``
	case msg.SectionCluster:
		switch q.Action {
		case msg.ActionCreate:
		case msg.ActionDestroy:
			return q.Repository.ID, q.Bucket.ID
		case msg.ActionMemberAssign:
		case msg.ActionMemberUnassign:
		case msg.ActionPropertyCreate:
		case msg.ActionPropertyUpdate:
		case msg.ActionPropertyDestroy:
		default:
			return ``, ``
		}
		return ``, q.Cluster.BucketID
	case msg.SectionGroup:
		switch q.Action {
		case msg.ActionCreate:
		case msg.ActionDestroy:
			return q.Repository.ID, q.Bucket.ID
		case msg.ActionMemberAssign:
		case msg.ActionMemberUnassign:
		case msg.ActionPropertyCreate:
		case msg.ActionPropertyUpdate:
		case msg.ActionPropertyDestroy:
		default:
			return ``, ``
		}
		return ``, q.Group.BucketID
	case msg.SectionBucket:
		switch q.Action {
		case msg.ActionCreate, msg.ActionRename:
			return q.Bucket.RepositoryID, ``
		case msg.ActionDestroy:
			return q.Bucket.RepositoryID, q.Bucket.ID
		}
		switch q.Action {
		case msg.ActionPropertyCreate:
		case msg.ActionPropertyUpdate:
		case msg.ActionPropertyDestroy:
		default:
			return ``, ``
		}
		return ``, q.Bucket.ID
	case msg.SectionRepositoryConfig:
		switch q.Action {
		case msg.ActionPropertyCreate:
		case msg.ActionPropertyUpdate:
		case msg.ActionPropertyDestroy:
		default:
			return ``, ``
		}
		return q.Repository.ID, ``
	case msg.SectionRepository:
		switch q.Action {
		case msg.ActionDestroy:
		case msg.ActionRename:
		case msg.ActionRepossess:
		default:
			return ``, ``
		}
		return q.Repository.ID, ``
	}
	return ``, ``
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
