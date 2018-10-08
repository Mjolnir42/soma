/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/tree"
	"github.com/mjolnir42/soma/lib/proto"
	"github.com/satori/go.uuid"
)

func (tk *TreeKeeper) addCheck(config *proto.CheckConfig) error {
	var err error
	var chk *tree.Check
	if chk, err = tk.convertCheck(config); err == nil {
		tk.tree.Find(tree.FindRequest{
			ElementType: config.ObjectType,
			ElementID:   config.ObjectID,
		}, true).SetCheck(*chk)
		return nil
	}
	return err
}

func (tk *TreeKeeper) rmCheck(config *proto.CheckConfig) error {
	var err error
	var chk *tree.Check
	if chk, err = tk.convertCheckForDelete(config); err == nil {
		tk.tree.Find(tree.FindRequest{
			ElementType: config.ObjectType,
			ElementID:   config.ObjectID,
		}, true).DeleteCheck(*chk)
		return nil
	}
	return err
}

func (tk *TreeKeeper) convertCheck(conf *proto.CheckConfig) (*tree.Check, error) {
	treechk := &tree.Check{
		ID:            uuid.Nil,
		SourceID:      uuid.Nil,
		InheritedFrom: uuid.Nil,
		Inheritance:   conf.Inheritance,
		ChildrenOnly:  conf.ChildrenOnly,
		Interval:      conf.Interval,
	}
	treechk.CapabilityID, _ = uuid.FromString(conf.CapabilityID)
	treechk.ConfigID, _ = uuid.FromString(conf.ID)
	if err := tk.stmtGetView.QueryRow(conf.CapabilityID).Scan(&treechk.View); err != nil {
		return &tree.Check{}, err
	}

	treechk.Thresholds = make([]tree.CheckThreshold, len(conf.Thresholds))
	for i, thr := range conf.Thresholds {
		nthr := tree.CheckThreshold{
			Predicate: thr.Predicate.Symbol,
			Level:     uint8(thr.Level.Numeric),
			Value:     thr.Value,
		}
		treechk.Thresholds[i] = nthr
	}

	treechk.Constraints = make([]tree.CheckConstraint, len(conf.Constraints))
	for i, constr := range conf.Constraints {
		ncon := tree.CheckConstraint{
			Type: constr.ConstraintType,
		}
		switch constr.ConstraintType {
		case msg.ConstraintNative:
			ncon.Key = constr.Native.Name
			ncon.Value = constr.Native.Value
		case msg.ConstraintOncall:
			ncon.Key = `OncallId`
			ncon.Value = constr.Oncall.ID
		case msg.ConstraintCustom:
			ncon.Key = constr.Custom.ID
			ncon.Value = constr.Custom.Value
		case msg.ConstraintSystem:
			ncon.Key = constr.System.Name
			ncon.Value = constr.System.Value
		case msg.ConstraintService:
			ncon.Key = `id`
			ncon.Value = constr.Service.ID
		case msg.ConstraintAttribute:
			ncon.Key = constr.Attribute.Name
			ncon.Value = constr.Attribute.Value
		}
		treechk.Constraints[i] = ncon
	}
	return treechk, nil
}

func (tk *TreeKeeper) convertCheckForDelete(conf *proto.CheckConfig) (*tree.Check, error) {
	var err error
	treechk := &tree.Check{
		ID:            uuid.Nil,
		InheritedFrom: uuid.Nil,
	}
	if treechk.SourceID, err = uuid.FromString(conf.ExternalID); err != nil {
		return nil, err
	}
	if treechk.ConfigID, err = uuid.FromString(conf.ID); err != nil {
		return nil, err
	}
	return treechk, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
