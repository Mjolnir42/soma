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

func (tk *TreeKeeper) addProperty(q *msg.Request) {
	prop, id := tk.convProperty(`add`, q)
	tk.tree.Find(tree.FindRequest{
		ElementType: q.TargetEntity,
		ElementID:   id,
	}, true).(tree.Propertier).SetProperty(prop)
}

func (tk *TreeKeeper) rmProperty(q *msg.Request) {
	prop, id := tk.convProperty(`rm`, q)
	tk.tree.Find(tree.FindRequest{
		ElementType: q.TargetEntity,
		ElementID:   id,
	}, true).(tree.Propertier).DeleteProperty(prop)
}

func (tk *TreeKeeper) updateProperty(q *msg.Request) {
	prop, id := tk.convProperty(`update`, q)
	tk.tree.Find(tree.FindRequest{
		ElementType: q.TargetEntity,
		ElementID:   id,
	}, true).(tree.Propertier).UpdateProperty(prop)
}

func (tk *TreeKeeper) convProperty(task string, q *msg.Request) (
	tree.Property, string) {

	var prop tree.Property
	var id string

	switch q.Section {
	case msg.SectionNodeConfig:
		id = q.Node.ID
		prop = tk.pTT(task, (*q.Node.Properties)[0])
	case msg.SectionCluster:
		id = q.Cluster.ID
		prop = tk.pTT(task, (*q.Cluster.Properties)[0])
	case msg.SectionGroup:
		id = q.Group.ID
		prop = tk.pTT(task, (*q.Group.Properties)[0])
	case msg.SectionBucket:
		id = q.Bucket.ID
		prop = tk.pTT(task, (*q.Bucket.Properties)[0])
	case msg.SectionRepositoryConfig:
		id = q.Repository.ID
		prop = tk.pTT(task, (*q.Repository.Properties)[0])
	}
	return prop, id
}

func (tk *TreeKeeper) pTT(task string, pp proto.Property) tree.Property {
	switch pp.Type {
	case msg.PropertyCustom:
		customID, _ := uuid.FromString(pp.Custom.ID)
		switch task {
		case `add`:
			return &tree.PropertyCustom{
				ID:           uuid.Must(uuid.NewV4()),
				CustomID:     customID,
				Inheritance:  pp.Inheritance,
				ChildrenOnly: pp.ChildrenOnly,
				View:         pp.View,
				Key:          pp.Custom.Name,
				Value:        pp.Custom.Value,
			}
		case `rm`:
			srcUUID, _ := uuid.FromString(pp.SourceInstanceID)
			return &tree.PropertyCustom{
				SourceID: srcUUID,
				CustomID: customID,
				View:     pp.View,
				Key:      pp.Custom.Name,
				Value:    pp.Custom.Value,
			}
		case `update`:
			srcUUID, _ := uuid.FromString(pp.SourceInstanceID)
			return &tree.PropertyCustom{
				ID:           srcUUID,
				SourceID:     srcUUID,
				CustomID:     customID,
				Inheritance:  pp.Inheritance,
				ChildrenOnly: pp.ChildrenOnly,
				View:         pp.View,
				Key:          pp.Custom.Name,
				Value:        pp.Custom.Value,
			}
		}
	case msg.PropertyOncall:
		oncallID, _ := uuid.FromString(pp.Oncall.ID)
		switch task {
		case `add`:
			return &tree.PropertyOncall{
				ID:           uuid.Must(uuid.NewV4()),
				OncallID:     oncallID,
				Inheritance:  pp.Inheritance,
				ChildrenOnly: pp.ChildrenOnly,
				View:         pp.View,
				Name:         pp.Oncall.Name,
				Number:       pp.Oncall.Number,
			}
		case `rm`:
			srcUUID, _ := uuid.FromString(pp.SourceInstanceID)
			return &tree.PropertyOncall{
				SourceID: srcUUID,
				OncallID: oncallID,
				View:     pp.View,
				Name:     pp.Oncall.Name,
				Number:   pp.Oncall.Number,
			}
		case `update`:
			srcUUID, _ := uuid.FromString(pp.SourceInstanceID)
			return &tree.PropertyOncall{
				ID:           srcUUID,
				SourceID:     srcUUID,
				OncallID:     oncallID,
				Inheritance:  pp.Inheritance,
				ChildrenOnly: pp.ChildrenOnly,
				View:         pp.View,
				Name:         pp.Oncall.Name,
				Number:       pp.Oncall.Number,
			}
		}
	case msg.PropertyService:
		switch task {
		case `add`:
			return &tree.PropertyService{
				ID:           uuid.Must(uuid.NewV4()),
				Inheritance:  pp.Inheritance,
				ChildrenOnly: pp.ChildrenOnly,
				View:         pp.View,
				Service:      pp.Service.Name,
				Attributes:   pp.Service.Attributes,
			}
		case `rm`:
			srcUUID, _ := uuid.FromString(pp.SourceInstanceID)
			return &tree.PropertyService{
				SourceID: srcUUID,
				View:     pp.View,
				Service:  pp.Service.Name,
			}
		case `update`:
			srcUUID, _ := uuid.FromString(pp.SourceInstanceID)
			return &tree.PropertyService{
				ID:           srcUUID,
				SourceID:     srcUUID,
				Inheritance:  pp.Inheritance,
				ChildrenOnly: pp.ChildrenOnly,
				View:         pp.View,
				Service:      pp.Service.Name,
				Attributes:   pp.Service.Attributes,
			}
		}
	case msg.PropertySystem:
		switch task {
		case `add`:
			return &tree.PropertySystem{
				ID:           uuid.Must(uuid.NewV4()),
				Inheritance:  pp.Inheritance,
				ChildrenOnly: pp.ChildrenOnly,
				View:         pp.View,
				Key:          pp.System.Name,
				Value:        pp.System.Value,
			}
		case `rm`:
			srcUUID, _ := uuid.FromString(pp.SourceInstanceID)
			return &tree.PropertySystem{
				SourceID: srcUUID,
				View:     pp.View,
				Key:      pp.System.Name,
				Value:    pp.System.Value,
			}
		case `update`:
			srcUUID, _ := uuid.FromString(pp.SourceInstanceID)
			return &tree.PropertySystem{
				ID:           srcUUID,
				SourceID:     srcUUID,
				Inheritance:  pp.Inheritance,
				ChildrenOnly: pp.ChildrenOnly,
				View:         pp.View,
				Key:          pp.System.Name,
				Value:        pp.System.Value,
			}
		}
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
