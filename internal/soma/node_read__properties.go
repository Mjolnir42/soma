/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"database/sql"

	"github.com/mjolnir42/soma/lib/proto"
)

// oncallProperties adds the oncall properties of the node
func (r *NodeRead) oncallProperties(node *proto.Node) error {
	var (
		rows                               *sql.Rows
		err                                error
		instanceID, sourceInstanceID, view string
		oncallID, oncallName               string
	)

	if rows, err = r.stmtPropOncall.Query(
		node.ID,
	); err != nil {
		return err
	}

	for rows.Next() {
		if err = rows.Scan(
			&instanceID,
			&sourceInstanceID,
			&view,
			&oncallID,
			&oncallName,
		); err != nil {
			rows.Close()
			return err
		}
		*node.Properties = append(*node.Properties, proto.Property{
			Type:             `oncall`,
			RepositoryID:     node.Config.RepositoryID,
			BucketID:         node.Config.BucketID,
			InstanceID:       instanceID,
			SourceInstanceID: sourceInstanceID,
			View:             view,
			Oncall: &proto.PropertyOncall{
				ID:   oncallID,
				Name: oncallName,
			},
		})
	}
	err = rows.Err()
	return err
}

// serviceProperties adds the service properties of the node
func (r *NodeRead) serviceProperties(node *proto.Node) error {
	var (
		rows                         *sql.Rows
		err                          error
		instanceID, sourceInstanceID string
		serviceID, view              string
	)

	if rows, err = r.stmtPropService.Query(
		node.ID,
	); err != nil {
		return err
	}

	for rows.Next() {
		if err = rows.Scan(
			&instanceID,
			&sourceInstanceID,
			&view,
			&serviceID,
		); err != nil {
			rows.Close()
			return err
		}
		*node.Properties = append(*node.Properties, proto.Property{
			Type:             `service`,
			RepositoryID:     node.Config.RepositoryID,
			BucketID:         node.Config.BucketID,
			InstanceID:       instanceID,
			SourceInstanceID: sourceInstanceID,
			View:             view,
			Service: &proto.PropertyService{
				ID: serviceID,
			},
		})
	}
	err = rows.Err()
	return err
}

// systemProperties adds the system properties of the node
func (r *NodeRead) systemProperties(node *proto.Node) error {
	var (
		rows                               *sql.Rows
		err                                error
		instanceID, sourceInstanceID, view string
		value, systemProp                  string
	)

	if rows, err = r.stmtPropSystem.Query(
		node.ID,
	); err != nil {
		return err
	}

	for rows.Next() {
		if err = rows.Scan(
			&instanceID,
			&sourceInstanceID,
			&view,
			&systemProp,
			&value,
		); err != nil {
			rows.Close()
			return err
		}
		*node.Properties = append(*node.Properties, proto.Property{
			Type:             `system`,
			RepositoryID:     node.Config.RepositoryID,
			BucketID:         node.Config.BucketID,
			InstanceID:       instanceID,
			SourceInstanceID: sourceInstanceID,
			View:             view,
			System: &proto.PropertySystem{
				Name:  systemProp,
				Value: value,
			},
		})
	}
	err = rows.Err()
	return err
}

// customProperties adds the custom properties of the node
func (r *NodeRead) customProperties(node *proto.Node) error {
	var (
		rows                               *sql.Rows
		err                                error
		instanceID, sourceInstanceID, view string
		value, customProp, customID        string
	)

	if rows, err = r.stmtPropCustom.Query(
		node.ID,
	); err != nil {
		return err
	}

	for rows.Next() {
		if err = rows.Scan(
			&instanceID,
			&sourceInstanceID,
			&view,
			&customID,
			&value,
			&customProp,
		); err != nil {
			rows.Close()
			return err
		}
		*node.Properties = append(*node.Properties, proto.Property{
			Type:             `custom`,
			RepositoryID:     node.Config.RepositoryID,
			BucketID:         node.Config.BucketID,
			InstanceID:       instanceID,
			SourceInstanceID: sourceInstanceID,
			View:             view,
			Custom: &proto.PropertyCustom{
				ID:    customID,
				Name:  customProp,
				Value: value,
			},
		})
	}
	err = rows.Err()
	return err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
