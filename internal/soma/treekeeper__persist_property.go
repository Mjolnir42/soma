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

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/tree"
)

func (tk *TreeKeeper) txProperty(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	switch a.Action {
	case tree.ActionPropertyNew:
		return tk.txPropertyNew(a, stm)
	case tree.ActionPropertyDelete:
		return tk.txPropertyDelete(a, stm)
	case tree.ActionPropertyUpdate: // XXX BUG
		return fmt.Errorf(`TreeKeeper: MISSING TX HANDLER FOR tree.ActionPropertyUpdate`)
	default:
		return fmt.Errorf("Illegal property action: %s", a.Action)
	}
}

//
// PROPERTY NEW
func (tk *TreeKeeper) txPropertyNew(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	if _, err := stm[`PropertyInstanceCreate`].Exec(
		a.Property.InstanceID,
		a.Property.RepositoryID,
		a.Property.SourceInstanceID,
		a.Property.SourceType,
		a.Property.InheritedFrom,
	); err != nil {
		return err
	}

	switch a.Property.Type {
	case msg.PropertyCustom:
		return tk.txPropertyNewCustom(a, stm)
	case msg.PropertySystem:
		return tk.txPropertyNewSystem(a, stm)
	case msg.PropertyService:
		return tk.txPropertyNewService(a, stm)
	case msg.PropertyOncall:
		return tk.txPropertyNewOncall(a, stm)
	}
	return fmt.Errorf(`Impossible property type`)
}

func (tk *TreeKeeper) txPropertyNewCustom(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	var (
		err       error
		statement *sql.Stmt
		id        string
	)
	switch a.Type {
	case msg.EntityRepository:
		statement = stm[`RepositoryPropertyCustomCreate`]
		id = a.Property.Custom.RepositoryID
	case msg.EntityBucket:
		statement = stm[`BucketPropertyCustomCreate`]
		id = a.Bucket.ID
	case msg.EntityGroup:
		statement = stm[`GroupPropertyCustomCreate`]
		id = a.Group.ID
	case msg.EntityCluster:
		statement = stm[`ClusterPropertyCustomCreate`]
		id = a.Cluster.ID
	case msg.EntityNode:
		statement = stm[`NodePropertyCustomCreate`]
		id = a.Node.ID
	}
	_, err = statement.Exec(
		a.Property.InstanceID,
		a.Property.SourceInstanceID,
		id,
		a.Property.View,
		a.Property.Custom.ID,
		a.Property.Inheritance,
		a.Property.ChildrenOnly,
		a.Property.Custom.Value,
	)
	return err
}

func (tk *TreeKeeper) txPropertyNewSystem(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	var (
		err       error
		statement *sql.Stmt
		id        string
	)
	switch a.Type {
	case msg.EntityRepository:
		statement = stm[`RepositoryPropertySystemCreate`]
		id = a.Repository.ID
	case msg.EntityBucket:
		statement = stm[`BucketPropertySystemCreate`]
		id = a.Bucket.ID
	case msg.EntityGroup:
		statement = stm[`GroupPropertySystemCreate`]
		id = a.Group.ID
	case msg.EntityCluster:
		statement = stm[`ClusterPropertySystemCreate`]
		id = a.Cluster.ID
	case msg.EntityNode:
		statement = stm[`NodePropertySystemCreate`]
		id = a.Node.ID
	}
	_, err = statement.Exec(
		a.Property.InstanceID,
		a.Property.SourceInstanceID,
		id,
		a.Property.View,
		a.Property.System.Name,
		a.Property.SourceType,
		a.Property.RepositoryID,
		a.Property.Inheritance,
		a.Property.ChildrenOnly,
		a.Property.System.Value,
		a.Property.IsInherited,
	)
	return err
}

func (tk *TreeKeeper) txPropertyNewService(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	var (
		err       error
		statement *sql.Stmt
		id        string
	)
	switch a.Type {
	case msg.EntityRepository:
		statement = stm[`RepositoryPropertyServiceCreate`]
		id = a.Repository.ID
	case msg.EntityBucket:
		statement = stm[`BucketPropertyServiceCreate`]
		id = a.Bucket.ID
	case msg.EntityGroup:
		statement = stm[`GroupPropertyServiceCreate`]
		id = a.Group.ID
	case msg.EntityCluster:
		statement = stm[`ClusterPropertyServiceCreate`]
		id = a.Cluster.ID
	case msg.EntityNode:
		statement = stm[`NodePropertyServiceCreate`]
		id = a.Node.ID
	}
	_, err = statement.Exec(
		a.Property.InstanceID,
		a.Property.SourceInstanceID,
		id,
		a.Property.View,
		a.Property.Service.Name,
		a.Property.Service.TeamID,
		a.Property.RepositoryID,
		a.Property.Inheritance,
		a.Property.ChildrenOnly,
	)
	return err
}

func (tk *TreeKeeper) txPropertyNewOncall(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	var (
		err       error
		statement *sql.Stmt
		id        string
	)
	switch a.Type {
	case msg.EntityRepository:
		statement = stm[`RepositoryPropertyOncallCreate`]
		id = a.Repository.ID
	case msg.EntityBucket:
		statement = stm[`BucketPropertyOncallCreate`]
		id = a.Bucket.ID
	case msg.EntityGroup:
		statement = stm[`GroupPropertyOncallCreate`]
		id = a.Group.ID
	case msg.EntityCluster:
		statement = stm[`ClusterPropertyOncallCreate`]
		id = a.Cluster.ID
	case msg.EntityNode:
		statement = stm[`NodePropertyOncallCreate`]
		id = a.Node.ID
	}
	_, err = statement.Exec(
		a.Property.InstanceID,
		a.Property.SourceInstanceID,
		id,
		a.Property.View,
		a.Property.Oncall.ID,
		a.Property.RepositoryID,
		a.Property.Inheritance,
		a.Property.ChildrenOnly,
	)
	return err
}

//
// PROPERTY DELETE
func (tk *TreeKeeper) txPropertyDelete(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	if _, err := stm[`PropertyInstanceDelete`].Exec(
		a.Property.InstanceID,
	); err != nil {
		return err
	}

	var statement *sql.Stmt
	switch a.Property.Type {
	case `custom`:
		switch a.Type {
		case msg.EntityRepository:
			statement = stm[`RepositoryPropertyCustomDelete`]
		case msg.EntityBucket:
			statement = stm[`BucketPropertyCustomDelete`]
		case msg.EntityGroup:
			statement = stm[`GroupPropertyCustomDelete`]
		case msg.EntityCluster:
			statement = stm[`ClusterPropertyCustomDelete`]
		case msg.EntityNode:
			statement = stm[`NodePropertyCustomDelete`]
		}
	case `system`:
		switch a.Type {
		case msg.EntityRepository:
			statement = stm[`RepositoryPropertySystemDelete`]
		case msg.EntityBucket:
			statement = stm[`BucketPropertySystemDelete`]
		case msg.EntityGroup:
			statement = stm[`GroupPropertySystemDelete`]
		case msg.EntityCluster:
			statement = stm[`ClusterPropertySystemDelete`]
		case msg.EntityNode:
			statement = stm[`NodePropertySystemDelete`]
		}
	case `service`:
		switch a.Type {
		case msg.EntityRepository:
			statement = stm[`RepositoryPropertyServiceDelete`]
		case msg.EntityBucket:
			statement = stm[`BucketPropertyServiceDelete`]
		case msg.EntityGroup:
			statement = stm[`GroupPropertyServiceDelete`]
		case msg.EntityCluster:
			statement = stm[`ClusterPropertyServiceDelete`]
		case msg.EntityNode:
			statement = stm[`NodePropertyServiceDelete`]
		}
	case `oncall`:
		switch a.Type {
		case msg.EntityRepository:
			statement = stm[`RepositoryPropertyOncallDelete`]
		case msg.EntityBucket:
			statement = stm[`BucketPropertyOncallDelete`]
		case msg.EntityGroup:
			statement = stm[`GroupPropertyOncallDelete`]
		case msg.EntityCluster:
			statement = stm[`ClusterPropertyOncallDelete`]
		case msg.EntityNode:
			statement = stm[`NodePropertyOncallDelete`]
		}
	}
	_, err := statement.Exec(
		a.Property.InstanceID,
	)
	return err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
