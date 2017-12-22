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

func (tk *TreeKeeper) txProperty(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	switch a.Action {
	case `property_new`:
		return tk.txPropertyNew(a, stm)
	case `property_delete`:
		return tk.txPropertyDelete(a, stm)
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
	case `custom`:
		return tk.txPropertyNewCustom(a, stm)
	case `system`:
		return tk.txPropertyNewSystem(a, stm)
	case `service`:
		return tk.txPropertyNewService(a, stm)
	case `oncall`:
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
	case `repository`:
		statement = stm[`RepositoryPropertyCustomCreate`]
		id = a.Property.Custom.RepositoryID
	case `bucket`:
		statement = stm[`BucketPropertyCustomCreate`]
		id = a.Bucket.ID
	case `group`:
		statement = stm[`GroupPropertyCustomCreate`]
		id = a.Group.ID
	case `cluster`:
		statement = stm[`ClusterPropertyCustomCreate`]
		id = a.Cluster.ID
	case `node`:
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
	case `repository`:
		statement = stm[`RepositoryPropertySystemCreate`]
		id = a.Repository.ID
	case `bucket`:
		statement = stm[`BucketPropertySystemCreate`]
		id = a.Bucket.ID
	case `group`:
		statement = stm[`GroupPropertySystemCreate`]
		id = a.Group.ID
	case `cluster`:
		statement = stm[`ClusterPropertySystemCreate`]
		id = a.Cluster.ID
	case `node`:
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
	case `repository`:
		statement = stm[`RepositoryPropertyServiceCreate`]
		id = a.Repository.ID
	case `bucket`:
		statement = stm[`BucketPropertyServiceCreate`]
		id = a.Bucket.ID
	case `group`:
		statement = stm[`GroupPropertyServiceCreate`]
		id = a.Group.ID
	case `cluster`:
		statement = stm[`ClusterPropertyServiceCreate`]
		id = a.Cluster.ID
	case `node`:
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
	case `repository`:
		statement = stm[`RepositoryPropertyOncallCreate`]
		id = a.Repository.ID
	case `bucket`:
		statement = stm[`BucketPropertyOncallCreate`]
		id = a.Bucket.ID
	case `group`:
		statement = stm[`GroupPropertyOncallCreate`]
		id = a.Group.ID
	case `cluster`:
		statement = stm[`ClusterPropertyOncallCreate`]
		id = a.Cluster.ID
	case `node`:
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
		case `repository`:
			statement = stm[`RepositoryPropertyCustomDelete`]
		case `bucket`:
			statement = stm[`BucketPropertyCustomDelete`]
		case `group`:
			statement = stm[`GroupPropertyCustomDelete`]
		case `cluster`:
			statement = stm[`ClusterPropertyCustomDelete`]
		case `node`:
			statement = stm[`NodePropertyCustomDelete`]
		}
	case `system`:
		switch a.Type {
		case `repository`:
			statement = stm[`RepositoryPropertySystemDelete`]
		case `bucket`:
			statement = stm[`BucketPropertySystemDelete`]
		case `group`:
			statement = stm[`GroupPropertySystemDelete`]
		case `cluster`:
			statement = stm[`ClusterPropertySystemDelete`]
		case `node`:
			statement = stm[`NodePropertySystemDelete`]
		}
	case `service`:
		switch a.Type {
		case `repository`:
			statement = stm[`RepositoryPropertyServiceDelete`]
		case `bucket`:
			statement = stm[`BucketPropertyServiceDelete`]
		case `group`:
			statement = stm[`GroupPropertyServiceDelete`]
		case `cluster`:
			statement = stm[`ClusterPropertyServiceDelete`]
		case `node`:
			statement = stm[`NodePropertyServiceDelete`]
		}
	case `oncall`:
		switch a.Type {
		case `repository`:
			statement = stm[`RepositoryPropertyOncallDelete`]
		case `bucket`:
			statement = stm[`BucketPropertyOncallDelete`]
		case `group`:
			statement = stm[`GroupPropertyOncallDelete`]
		case `cluster`:
			statement = stm[`ClusterPropertyOncallDelete`]
		case `node`:
			statement = stm[`NodePropertyOncallDelete`]
		}
	}
	_, err := statement.Exec(
		a.Property.InstanceID,
	)
	return err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
