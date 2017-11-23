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

// oncallProperties adds the oncall properties of the bucket
func (r *BucketRead) oncallProperties(bucket *proto.Bucket) error {
	var (
		instanceID, sourceInstanceID string
		view, oncallID, oncallName   string
		rows                         *sql.Rows
		err                          error
	)

	if rows, err = r.stmtPropOncall.Query(
		bucket.ID,
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
		*bucket.Properties = append(*bucket.Properties,
			proto.Property{
				Type:             `oncall`,
				RepositoryID:     bucket.RepositoryID,
				BucketID:         bucket.ID,
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

// serviceProperties adds the service properties of the bucket
func (r *BucketRead) serviceProperties(bucket *proto.Bucket) error {
	var (
		instanceID, sourceInstanceID string
		serviceName, view            string
		rows                         *sql.Rows
		err                          error
	)

	if rows, err = r.stmtPropService.Query(
		bucket.ID,
	); err != nil {
		return err
	}

	for rows.Next() {
		if err = rows.Scan(
			&instanceID,
			&sourceInstanceID,
			&view,
			&serviceName,
		); err != nil {
			rows.Close()
			return err
		}
		*bucket.Properties = append(*bucket.Properties,
			proto.Property{
				Type:             `service`,
				RepositoryID:     bucket.RepositoryID,
				BucketID:         bucket.ID,
				InstanceID:       instanceID,
				SourceInstanceID: sourceInstanceID,
				View:             view,
				Service: &proto.PropertyService{
					Name: serviceName,
				},
			},
		)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

// systemProperties adds the system properties of the bucket
func (r *BucketRead) systemProperties(bucket *proto.Bucket) error {
	var (
		instanceID, sourceInstanceID, view string
		systemProp, value                  string
		rows                               *sql.Rows
		err                                error
	)

	if rows, err = r.stmtPropSystem.Query(
		bucket.ID,
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
		*bucket.Properties = append(*bucket.Properties,
			proto.Property{
				Type:             `system`,
				RepositoryID:     bucket.RepositoryID,
				BucketID:         bucket.ID,
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

// customProperties adds the custom properties of the bucket
func (r *BucketRead) customProperties(bucket *proto.Bucket) error {
	var (
		instanceID, sourceInstanceID, view string
		customID, value, customProp        string
		rows                               *sql.Rows
		err                                error
	)

	if rows, err = r.stmtPropCustom.Query(
		bucket.ID,
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
		*bucket.Properties = append(*bucket.Properties,
			proto.Property{
				Type:             `custom`,
				RepositoryID:     bucket.RepositoryID,
				BucketID:         bucket.ID,
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
