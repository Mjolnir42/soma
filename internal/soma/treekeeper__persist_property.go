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
	case tree.ActionPropertyUpdate:
		return tk.txPropertyUpdate(a, stm)
	default:
		return fmt.Errorf("Illegal property action: %s", a.Action)
	}
}

//
// PROPERTY NEW
func (tk *TreeKeeper) txPropertyNew(a *tree.Action,
	stm map[string]*sql.Stmt) error {

	var (
		err              error
		statement        *sql.Stmt
		entity, entityID string
	)

	if _, err = stm[`PropertyInstanceCreate`].Exec(
		a.Property.InstanceID,
		a.Property.RepositoryID,
		a.Property.SourceInstanceID,
		a.Property.SourceType,
		a.Property.InheritedFrom,
	); err != nil {
		return err
	}

	entity, entityID = getEntityData(a)

	switch a.Property.Type {
	case msg.PropertyCustom:
		statement = stm[entity+"PropertyCustomCreate"]
		_, err = statement.Exec(
			a.Property.InstanceID,
			a.Property.SourceInstanceID,
			entityID,
			a.Property.View,
			a.Property.Custom.ID,
			a.Property.RepositoryID,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
			a.Property.Custom.Value,
		)
		return err
	case msg.PropertySystem:
		statement = stm[entity+"PropertySystemCreate"]
		_, err = statement.Exec(
			a.Property.InstanceID,
			a.Property.SourceInstanceID,
			entityID,
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
	case msg.PropertyService:
		statement = stm[entity+"PropertyServiceCreate"]
		_, err = statement.Exec(
			a.Property.InstanceID,
			a.Property.SourceInstanceID,
			entityID,
			a.Property.View,
			a.Property.Service.ID,
			a.Property.Service.TeamID,
			a.Property.RepositoryID,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
		)
		return err
	case msg.PropertyOncall:
		statement = stm[entity+"PropertyOncallCreate"]
		_, err = statement.Exec(
			a.Property.InstanceID,
			a.Property.SourceInstanceID,
			entityID,
			a.Property.View,
			a.Property.Oncall.ID,
			a.Property.RepositoryID,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
		)
		return err
	}
	return fmt.Errorf(`Impossible property type`)
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

//
// PROPERTY UPDATE
func (tk *TreeKeeper) txPropertyUpdate(a *tree.Action,
	stm map[string]*sql.Stmt) error {

	var (
		err              error
		statement        *sql.Stmt
		entity, entityID string
	)
	entity, entityID = getEntityData(a)
	switch a.Property.Type {
	case msg.PropertyCustom:
		statement = stm[entity+"PropertyCustomUpdate"]
		_, err = statement.Exec(
			a.Property.InstanceID,
			a.Property.SourceInstanceID,
			entityID,
			a.Property.View,
			a.Property.Custom.ID,
			a.Property.Inheritance,
			a.Property.Custom.Value,
		)
	case msg.PropertySystem:
		statement = stm[entity+"PropertySystemUpdate"]
		_, err = statement.Exec(
			a.Property.InstanceID,
			a.Property.SourceInstanceID,
			entityID,
			a.Property.View,
			a.Property.System.Name,
			a.Property.Inheritance,
			a.Property.System.Value,
		)
		return err
	case msg.PropertyService, msg.PropertyOncall:
		return err
	}
	return fmt.Errorf(`Impossible property type`)
}

func getEntityData(a *tree.Action) (object, id string) {
	switch a.Type {
	case msg.EntityRepository:
		return "Repository", a.Repository.ID
	case msg.EntityBucket:
		return "Bucket", a.Bucket.ID
	case msg.EntityGroup:
		return "Group", a.Group.ID
	case msg.EntityCluster:
		return "Cluster", a.Cluster.ID
	case msg.EntityNode:
		return "Node", a.Node.ID
	}
	return "", ""
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
