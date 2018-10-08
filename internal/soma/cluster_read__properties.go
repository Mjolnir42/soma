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

// oncallProperties adds the oncall properties of the cluster
func (r *ClusterRead) oncallProperties(cluster *proto.Cluster) error {
	var (
		instanceID, sourceInstanceID string
		view, oncallID, oncallName   string
		rows                         *sql.Rows
		err                          error
	)

	if rows, err = r.stmtPropOncall.Query(
		cluster.ID,
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
		*cluster.Properties = append(*cluster.Properties,
			proto.Property{
				Type:             `oncall`,
				BucketID:         cluster.BucketID,
				InstanceID:       instanceID,
				SourceInstanceID: sourceInstanceID,
				View:             view,
				Oncall: &proto.PropertyOncall{
					ID:   oncallID,
					Name: oncallName,
				},
			},
		)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

// serviceProperties adds the service properties of the cluster
func (r *ClusterRead) serviceProperties(cluster *proto.Cluster) error {
	var (
		instanceID, sourceInstanceID string
		serviceID, view              string
		rows                         *sql.Rows
		err                          error
	)

	if rows, err = r.stmtPropService.Query(
		cluster.ID,
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
		*cluster.Properties = append(*cluster.Properties,
			proto.Property{
				Type:             `service`,
				BucketID:         cluster.BucketID,
				InstanceID:       instanceID,
				SourceInstanceID: sourceInstanceID,
				View:             view,
				Service: &proto.PropertyService{
					ID: serviceID,
				},
			},
		)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

// systemProperties adds the system properties of the cluster
func (r *ClusterRead) systemProperties(cluster *proto.Cluster) error {
	var (
		instanceID, sourceInstanceID, view string
		systemProp, value                  string
		rows                               *sql.Rows
		err                                error
	)

	if rows, err = r.stmtPropSystem.Query(
		cluster.ID,
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
		*cluster.Properties = append(*cluster.Properties,
			proto.Property{
				Type:             `system`,
				BucketID:         cluster.BucketID,
				InstanceID:       instanceID,
				SourceInstanceID: sourceInstanceID,
				View:             view,
				System: &proto.PropertySystem{
					Name:  systemProp,
					Value: value,
				},
			},
		)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

// customProperties adds the custom properties of the cluster
func (r *ClusterRead) customProperties(cluster *proto.Cluster) error {
	var (
		instanceID, sourceInstanceID, view string
		customID, value, customProp        string
		rows                               *sql.Rows
		err                                error
	)

	if rows, err = r.stmtPropCustom.Query(
		cluster.ID,
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
		*cluster.Properties = append(*cluster.Properties,
			proto.Property{
				Type:             `custom`,
				BucketID:         cluster.BucketID,
				InstanceID:       instanceID,
				SourceInstanceID: sourceInstanceID,
				View:             view,
				Custom: &proto.PropertyCustom{
					ID:    customID,
					Name:  customProp,
					Value: value,
				},
			},
		)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
