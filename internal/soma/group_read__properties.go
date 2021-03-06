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

// oncallProperties adds the oncall properties of the group
func (r *GroupRead) oncallProperties(group *proto.Group) error {
	var (
		instanceID, sourceInstanceID string
		view, oncallID, oncallName   string
		rows                         *sql.Rows
		err                          error
	)

	if rows, err = r.stmtPropOncall.Query(
		group.ID,
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
		*group.Properties = append(*group.Properties,
			proto.Property{
				Type:             `oncall`,
				BucketID:         group.BucketID,
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

// serviceProperties adds the service properties of the group
func (r *GroupRead) serviceProperties(group *proto.Group) error {
	var (
		instanceID, sourceInstanceID string
		serviceID, view              string
		rows                         *sql.Rows
		err                          error
	)

	if rows, err = r.stmtPropService.Query(
		group.ID,
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
		*group.Properties = append(*group.Properties,
			proto.Property{
				Type:             `service`,
				BucketID:         group.BucketID,
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

// systemProperties adds the system properties of the group
func (r *GroupRead) systemProperties(group *proto.Group) error {
	var (
		instanceID, sourceInstanceID, view string
		systemProp, value                  string
		rows                               *sql.Rows
		err                                error
	)

	if rows, err = r.stmtPropSystem.Query(
		group.ID,
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
		*group.Properties = append(*group.Properties,
			proto.Property{
				Type:             `system`,
				BucketID:         group.BucketID,
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

// customProperties adds the custom properties of the group
func (r *GroupRead) customProperties(group *proto.Group) error {
	var (
		instanceID, sourceInstanceID, view string
		customID, value, customProp        string
		rows                               *sql.Rows
		err                                error
	)

	if rows, err = r.stmtPropCustom.Query(
		group.ID,
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
		*group.Properties = append(*group.Properties,
			proto.Property{
				Type:             `custom`,
				BucketID:         group.BucketID,
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
