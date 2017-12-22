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
	"strconv"
	"time"

	"github.com/mjolnir42/soma/internal/tree"
	"github.com/mjolnir42/soma/lib/proto"
)

func (tk *TreeKeeper) txCheckConfig(conf proto.CheckConfig,
	stm map[string]*sql.Stmt) error {
	var (
		nullBucket sql.NullString
		err        error
	)
	if conf.BucketID != "" {
		nullBucket = sql.NullString{
			String: conf.BucketID,
			Valid:  true,
		}
	} else {
		nullBucket = sql.NullString{String: "", Valid: false}
	}
	if _, err = stm[`CreateCheckConfigurationBase`].Exec(
		conf.ID,
		conf.Name,
		int64(conf.Interval),
		conf.RepositoryID,
		nullBucket,
		conf.CapabilityID,
		conf.ObjectID,
		conf.ObjectType,
		conf.IsActive,
		conf.IsEnabled,
		conf.Inheritance,
		conf.ChildrenOnly,
		conf.ExternalID,
	); err != nil {
		return err
	}

threshloop:
	for _, thr := range conf.Thresholds {
		if _, err = stm[`CreateCheckConfigurationThreshold`].Exec(
			conf.ID,
			thr.Predicate.Symbol,
			strconv.FormatInt(thr.Value, 10),
			thr.Level.Name,
		); err != nil {
			break threshloop
		}
	}
	if err != nil {
		return err
	}

constrloop:
	for _, constr := range conf.Constraints {
		switch constr.ConstraintType {
		case "native":
			if _, err = stm[`CreateCheckConfigurationConstraintNative`].Exec(
				conf.ID,
				constr.Native.Name,
				constr.Native.Value,
			); err != nil {
				break constrloop
			}
		case "oncall":
			if _, err = stm[`CreateCheckConfigurationConstraintOncall`].Exec(
				conf.ID,
				constr.Oncall.ID,
			); err != nil {
				break constrloop
			}
		case "custom":
			if _, err = stm[`CreateCheckConfigurationConstraintCustom`].Exec(
				conf.ID,
				constr.Custom.ID,
				constr.Custom.RepositoryID,
				constr.Custom.Value,
			); err != nil {
				break constrloop
			}
		case "system":
			if _, err = stm[`CreateCheckConfigurationConstraintSystem`].Exec(
				conf.ID,
				constr.System.Name,
				constr.System.Value,
			); err != nil {
				break constrloop
			}
		case "service":
			if constr.Service.TeamID != tk.meta.teamID {
				err = fmt.Errorf(
					"Service constraint has mismatched TeamID values: %s/%s",
					tk.meta.teamID, constr.Service.TeamID)
				break constrloop
			}
			if _, err = stm[`CreateCheckConfigurationConstraintService`].Exec(
				conf.ID,
				tk.meta.teamID,
				constr.Service.Name,
			); err != nil {
				break constrloop
			}
		case "attribute":
			if _, err = stm[`CreateCheckConfigurationConstraintAttribute`].Exec(
				conf.ID,
				constr.Attribute.Name,
				constr.Attribute.Value,
			); err != nil {
				break constrloop
			}
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (tk *TreeKeeper) txCheck(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	switch a.Action {
	case `check_new`:
		return tk.txCheckNew(a, stm)
	case `check_removed`:
		return tk.txCheckRemoved(a, stm)
	default:
		return fmt.Errorf("Illegal check action: %s", a.Action)
	}
}

func (tk *TreeKeeper) txCheckNew(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	var id string
	bucket := sql.NullString{String: a.Bucket.ID, Valid: true}
	switch a.Type {
	case `repository`:
		id = a.Repository.ID
		bucket = sql.NullString{String: "", Valid: false}
	case `bucket`:
		id = a.Bucket.ID
	case `group`:
		id = a.Group.ID
	case `cluster`:
		id = a.Cluster.ID
	case `node`:
		id = a.Node.ID
	}
	statement := stm[`CreateCheck`]
	_, err := statement.Exec(
		a.Check.CheckID,
		a.Check.RepositoryID,
		bucket,
		a.Check.SourceCheckID,
		a.Check.SourceType,
		a.Check.InheritedFrom,
		a.Check.CheckConfigID,
		a.Check.CapabilityID,
		id,
		a.Type,
	)
	return err
}

func (tk *TreeKeeper) txCheckRemoved(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	statement := stm[`DeleteCheck`]
	_, err := statement.Exec(a.Check.CheckID)
	return err
}

func (tk *TreeKeeper) txCheckInstance(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	switch a.Type {
	case `repository`, `bucket`:
		return fmt.Errorf("Illegal check instance on %s", a.Type)
	}

	switch a.Action {
	case `check_instance_create`:
		if err := tk.txCheckInstanceCreate(a, stm); err != nil {
			return err
		}
		// for a new check instance, the first instance
		// configuration must be created alongside it
		fallthrough
	case `check_instance_update`:
		return tk.txCheckInstanceConfigCreate(a, stm)
	case `check_instance_delete`:
		return tk.txCheckInstanceDelete(a, stm)
	default:
		return fmt.Errorf("Illegal check instance action: %s", a.Action)
	}
}

func (tk *TreeKeeper) txCheckInstanceCreate(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	statement := stm[`CreateCheckInstance`]
	_, err := statement.Exec(
		a.CheckInstance.InstanceID,
		a.CheckInstance.CheckID,
		a.CheckInstance.ConfigID,
		`00000000-0000-0000-0000-000000000000`,
		time.Now().UTC(),
	)
	return err
}

func (tk *TreeKeeper) txCheckInstanceConfigCreate(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	statement := stm[`CreateCheckInstanceConfiguration`]
	_, err := statement.Exec(
		a.CheckInstance.InstanceConfigID,
		a.CheckInstance.Version,
		a.CheckInstance.InstanceID,
		a.CheckInstance.ConstraintHash,
		a.CheckInstance.ConstraintValHash,
		a.CheckInstance.InstanceService,
		a.CheckInstance.InstanceSvcCfgHash,
		a.CheckInstance.InstanceServiceConfig,
		time.Now().UTC(),
		`awaiting_computation`,
		`none`,
		false,
		"{}",
	)
	return err
}

func (tk *TreeKeeper) txCheckInstanceDelete(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	statement := stm[`DeleteCheckInstance`]
	_, err := statement.Exec(
		a.CheckInstance.InstanceID,
	)
	return err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
