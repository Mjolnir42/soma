/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"database/sql"
	"fmt"

	"github.com/mjolnir42/soma/internal/tree"
)

func (tk *TreeKeeper) txTree(a *tree.Action,
	stm map[string]*sql.Stmt, user string) error {
	switch a.Action {
	case `create`:
		return tk.txTreeCreate(a, stm, user)
	case `update`:
		return tk.txTreeUpdate(a, stm)
	case `delete`:
		return tk.txTreeDelete(a, stm)
	case `member_new`, `node_assignment`:
		return tk.txTreeMemberNew(a, stm)
	case `member_removed`:
		return tk.txTreeMemberRemoved(a, stm)
	default:
		return fmt.Errorf("Illegal tree action: %s", a.Action)
	}
}

func (tk *TreeKeeper) txTreeCreate(a *tree.Action,
	stm map[string]*sql.Stmt, user string) error {
	var err error
	switch a.Type {
	case `bucket`:
		_, err = stm[`CreateBucket`].Exec(
			a.Bucket.ID,
			a.Bucket.Name,
			a.Bucket.IsFrozen,
			a.Bucket.IsDeleted,
			a.Bucket.RepositoryID,
			a.Bucket.Environment,
			a.Bucket.TeamID,
			user,
		)
	case `group`:
		_, err = stm[`GroupCreate`].Exec(
			a.Group.Id,
			a.Group.BucketId,
			a.Group.Name,
			a.Group.ObjectState,
			a.Group.TeamId,
			user,
		)
	case `cluster`:
		_, err = stm[`ClusterCreate`].Exec(
			a.Cluster.ID,
			a.Cluster.Name,
			a.Cluster.BucketID,
			a.Cluster.ObjectState,
			a.Cluster.TeamID,
			user,
		)
	}
	return err
}

func (tk *TreeKeeper) txTreeUpdate(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	var (
		err          error
		statement    *sql.Stmt
		id, newState string
	)
	switch a.Type {
	case `group`:
		statement = stm[`GroupUpdate`]
		id = a.Group.Id
		newState = a.Group.ObjectState
	case `cluster`:
		statement = stm[`ClusterUpdate`]
		id = a.Cluster.ID
		newState = a.Cluster.ObjectState
	case `node`:
		statement = stm[`UpdateNodeState`]
		id = a.Node.ID
		newState = a.Node.State
	}
	_, err = statement.Exec(
		id,
		newState,
	)
	return err
}

func (tk *TreeKeeper) txTreeDelete(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	var err error
	switch a.Type {
	case `group`:
		_, err = stm[`GroupDelete`].Exec(
			a.Group.Id,
		)
	case `cluster`:
		_, err = stm[`ClusterDelete`].Exec(
			a.Cluster.ID,
		)
	case `node`:
		if _, err = stm[`NodeUnassignFromBucket`].Exec(
			a.Node.ID,
			a.Node.Config.BucketID,
			a.Node.TeamID,
		); err != nil {
			return err
		}
		// node unassign requires state update
		err = tk.txTreeUpdate(a, stm)
	}
	return err
}

func (tk *TreeKeeper) txTreeMemberNew(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	var (
		err               error
		id, child, bucket string
		statement         *sql.Stmt
	)
	switch a.Type {
	case `bucket`:
		_, err = stm[`BucketAssignNode`].Exec(
			a.ChildNode.ID,
			a.Bucket.ID,
			a.Bucket.TeamID,
		)
		goto exit
	case `group`:
		id = a.Group.Id
		bucket = a.Group.BucketId
		switch a.ChildType {
		case `group`:
			statement = stm[`GroupMemberNewGroup`]
			child = a.ChildGroup.Id
		case `cluster`:
			statement = stm[`GroupMemberNewCluster`]
			child = a.ChildCluster.ID
		case `node`:
			statement = stm[`GroupMemberNewNode`]
			child = a.ChildNode.ID
		}
	case `cluster`:
		id = a.Cluster.ID
		bucket = a.Cluster.BucketID
		child = a.ChildNode.ID
		statement = stm[`ClusterMemberNew`]
	}
	_, err = statement.Exec(
		id,
		child,
		bucket,
	)
exit:
	return err
}

func (tk *TreeKeeper) txTreeMemberRemoved(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	var (
		err       error
		id, child string
		statement *sql.Stmt
	)
	switch a.Type {
	case `group`:
		id = a.Group.Id
		switch a.ChildType {
		case `group`:
			statement = stm[`GroupMemberRemoveGroup`]
			child = a.ChildGroup.Id
		case `cluster`:
			statement = stm[`GroupMemberRemoveCluster`]
			child = a.ChildCluster.ID
		case `node`:
			statement = stm[`GroupMemberRemoveNode`]
			child = a.ChildNode.ID
		}
	case `cluster`:
		id = a.Cluster.ID
		child = a.ChildNode.ID
		statement = stm[`ClusterMemberRemove`]
	}
	_, err = statement.Exec(
		id,
		child,
	)
	return err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
