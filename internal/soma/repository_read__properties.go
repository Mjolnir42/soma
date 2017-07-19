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

// oncallProperties adds the oncall properties of the repository
func (r *RepositoryRead) oncallProperties(repo *proto.Repository) error {
	var (
		instanceID, sourceInstanceID string
		view, oncallID, oncallName   string
		rows                         *sql.Rows
		err                          error
	)

	if rows, err = r.stmtPropOncall.Query(
		repo.Id,
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
		*repo.Properties = append(*repo.Properties,
			proto.Property{
				Type:             `oncall`,
				RepositoryId:     repo.Id,
				InstanceId:       instanceID,
				SourceInstanceId: sourceInstanceID,
				View:             view,
				Oncall: &proto.PropertyOncall{
					Id:   oncallID,
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

// serviceProperties adds the service properties of the repository
func (r *RepositoryRead) serviceProperties(repo *proto.Repository) error {
	var (
		instanceID, sourceInstanceID string
		serviceName, view            string
		rows                         *sql.Rows
		err                          error
	)

	if rows, err = r.stmtPropService.Query(
		repo.Id,
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
		*repo.Properties = append(*repo.Properties,
			proto.Property{
				Type:             `service`,
				RepositoryId:     repo.Id,
				InstanceId:       instanceID,
				SourceInstanceId: sourceInstanceID,
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

// systemProperties adds the system properties of the repository
func (r *RepositoryRead) systemProperties(repo *proto.Repository) error {
	var (
		instanceID, sourceInstanceID, view string
		systemProp, value                  string
		rows                               *sql.Rows
		err                                error
	)

	if rows, err = r.stmtPropSystem.Query(
		repo.Id,
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
		*repo.Properties = append(*repo.Properties,
			proto.Property{
				Type:             `system`,
				RepositoryId:     repo.Id,
				InstanceId:       instanceID,
				SourceInstanceId: sourceInstanceID,
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

// customProperties adds the custom properties of the repository
func (r *RepositoryRead) customProperties(repo *proto.Repository) error {
	var (
		instanceID, sourceInstanceID, view string
		customID, value, customProp        string
		rows                               *sql.Rows
		err                                error
	)

	if rows, err = r.stmtPropCustom.Query(
		repo.Id,
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
		*repo.Properties = append(*repo.Properties,
			proto.Property{
				Type:             `custom`,
				RepositoryId:     repo.Id,
				InstanceId:       instanceID,
				SourceInstanceId: sourceInstanceID,
				View:             view,
				Custom: &proto.PropertyCustom{
					Id:    customID,
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
