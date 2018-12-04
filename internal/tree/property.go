/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

type Property interface {
	GetID() string
	GetInstanceID(objType string, objID uuid.UUID, l *log.Logger) uuid.UUID
	GetIsInherited() bool
	GetKey() string
	GetSource() string
	GetSourceInstance() string
	GetSourceType() string
	GetType() string
	GetValue() string
	GetView() string

	SetID(id uuid.UUID)
	SetInherited(inherited bool)
	SetInheritedFrom(id uuid.UUID)
	SetSourceID(id uuid.UUID)
	SetSourceType(s string)

	Clone() Property
	Equal(id uuid.UUID) bool
	MakeAction() Action

	hasInheritance() bool
	isChildrenOnly() bool
	clearInstances()
}

type PropertyInstance struct {
	ObjectID   uuid.UUID
	ObjectType string
	InstanceID uuid.UUID
}

//
// Custom
type PropertyCustom struct {
	// ID of the custom property
	ID uuid.UUID
	// ID of the source custom property this was inherited from
	SourceID uuid.UUID
	// ObjectType the source property was attached to
	SourceType string
	// ID of the custom property type
	CustomID uuid.UUID
	// Indicator if this was inherited
	Inherited bool
	// ID of the object the SourceID property is on
	InheritedFrom uuid.UUID
	// Inheritance is enabled/disabled
	Inheritance bool
	// ChildrenOnly is enabled/disabled
	ChildrenOnly bool
	// View this property is attached in
	View string
	// Property Key
	Key string
	// Property Value
	Value string
	// Filled with IDs during from-DB load to restore with same IDs
	Instances []PropertyInstance
}

func (p *PropertyCustom) GetType() string {
	return "custom"
}

func (p *PropertyCustom) GetID() string {
	return p.ID.String()
}

func (p *PropertyCustom) GetSource() string {
	return p.InheritedFrom.String()
}

func (p *PropertyCustom) hasInheritance() bool {
	return p.Inheritance
}

func (p *PropertyCustom) isChildrenOnly() bool {
	return p.ChildrenOnly
}

func (p *PropertyCustom) GetSourceInstance() string {
	return p.SourceID.String()
}

func (p *PropertyCustom) GetSourceType() string {
	return p.SourceType
}

func (p *PropertyCustom) GetIsInherited() bool {
	return p.Inherited
}

func (p *PropertyCustom) GetView() string {
	return p.View
}

func (p *PropertyCustom) GetKey() string {
	return p.CustomID.String()
}

func (p *PropertyCustom) GetValue() string {
	return p.Key + p.Value
}

func (p *PropertyCustom) GetKeyField() string {
	return p.Key
}

func (p *PropertyCustom) GetValueField() string {
	return p.Value
}

func (p *PropertyCustom) GetInstanceID(objType string, objID uuid.UUID, l *log.Logger) uuid.UUID {
	if !uuid.Equal(p.ID, uuid.Nil) {
		return p.ID
	}
	for _, instance := range p.Instances {
		if objType == instance.ObjectType && uuid.Equal(instance.ObjectID, objID) {
			l.Printf("tree.Property.GetInstanceID() found existing instance: %s\n", instance.InstanceID)
			return instance.InstanceID
		}
	}
	return uuid.Nil
}

func (p *PropertyCustom) SetID(id uuid.UUID) {
	p.ID, _ = uuid.FromString(id.String())
}

func (p *PropertyCustom) Equal(id uuid.UUID) bool {
	return uuid.Equal(p.ID, id)
}

func (p *PropertyCustom) clearInstances() {
	p.Instances = nil
}

func (p *PropertyCustom) SetInheritedFrom(id uuid.UUID) {
	p.InheritedFrom, _ = uuid.FromString(id.String())
}

func (p *PropertyCustom) SetInherited(inherited bool) {
	p.Inherited = inherited
}

func (p *PropertyCustom) SetSourceID(id uuid.UUID) {
	p.SourceID, _ = uuid.FromString(id.String())
}

func (p *PropertyCustom) SetSourceType(s string) {
	p.SourceType = s
}

func (p PropertyCustom) Clone() Property {
	cl := PropertyCustom{
		SourceType:   p.SourceType,
		Inherited:    p.Inherited,
		Inheritance:  p.Inheritance,
		ChildrenOnly: p.ChildrenOnly,
		View:         p.View,
		Key:          p.Key,
		Value:        p.Value,
	}
	cl.ID, _ = uuid.FromString(p.ID.String())
	cl.InheritedFrom, _ = uuid.FromString(p.InheritedFrom.String())
	cl.SourceID, _ = uuid.FromString(p.SourceID.String())
	cl.CustomID, _ = uuid.FromString(p.CustomID.String())
	cl.Instances = make([]PropertyInstance, len(p.Instances))
	copy(cl.Instances, p.Instances)

	return &cl
}

func (p *PropertyCustom) MakeAction() Action {
	return Action{
		Property: proto.Property{
			InstanceID:       p.GetID(),
			SourceInstanceID: p.GetSourceInstance(),
			SourceType:       p.GetSourceType(),
			IsInherited:      p.GetIsInherited(),
			InheritedFrom:    p.GetSource(),
			Type:             p.GetType(),
			Inheritance:      p.hasInheritance(),
			ChildrenOnly:     p.isChildrenOnly(),
			View:             p.GetView(),
			Custom: &proto.PropertyCustom{
				ID:    p.CustomID.String(),
				Name:  p.Key,
				Value: p.Value,
			},
		},
	}
}

//
// Service
type PropertyService struct {
	ID            uuid.UUID
	SourceID      uuid.UUID
	SourceType    string
	Inherited     bool
	InheritedFrom uuid.UUID
	Inheritance   bool
	ChildrenOnly  bool
	View          string
	ServiceID     uuid.UUID
	ServiceName   string
	Attributes    []proto.ServiceAttribute
	Instances     []PropertyInstance
}

func (p *PropertyService) GetType() string {
	return "service"
}

func (p *PropertyService) GetID() string {
	return p.ID.String()
}

func (p *PropertyService) GetSource() string {
	return p.InheritedFrom.String()
}

func (p *PropertyService) hasInheritance() bool {
	return p.Inheritance
}

func (p *PropertyService) isChildrenOnly() bool {
	return p.ChildrenOnly
}

func (p *PropertyService) GetSourceInstance() string {
	return p.SourceID.String()
}

func (p *PropertyService) GetSourceType() string {
	return p.SourceType
}

func (p *PropertyService) GetIsInherited() bool {
	return p.Inherited
}

func (p *PropertyService) GetView() string {
	return p.View
}

func (p *PropertyService) GetKey() string {
	return p.ServiceID.String()
}

func (p *PropertyService) GetValue() string {
	return p.ServiceName
}

func (p *PropertyService) GetInstanceID(objType string, objID uuid.UUID, l *log.Logger) uuid.UUID {
	if !uuid.Equal(p.ID, uuid.Nil) {
		return p.ID
	}
	for _, instance := range p.Instances {
		if objType == instance.ObjectType && uuid.Equal(instance.ObjectID, objID) {
			l.Printf("tree.Property.GetInstanceID() found existing instance: %s\n", instance.InstanceID)
			return instance.InstanceID
		}
	}
	return uuid.Nil
}

func (p *PropertyService) SetID(id uuid.UUID) {
	p.ID, _ = uuid.FromString(id.String())
}

func (p *PropertyService) Equal(id uuid.UUID) bool {
	return uuid.Equal(p.ID, id)
}

func (p *PropertyService) clearInstances() {
	p.Instances = nil
}

func (p *PropertyService) SetInheritedFrom(id uuid.UUID) {
	p.InheritedFrom, _ = uuid.FromString(id.String())
}

func (p *PropertyService) SetInherited(inherited bool) {
	p.Inherited = inherited
}

func (p *PropertyService) SetSourceID(id uuid.UUID) {
	p.SourceID, _ = uuid.FromString(id.String())
}

func (p *PropertyService) SetSourceType(s string) {
	p.SourceType = s
}

func (p PropertyService) Clone() Property {
	cl := PropertyService{
		SourceType:   p.SourceType,
		Inherited:    p.Inherited,
		Inheritance:  p.Inheritance,
		ChildrenOnly: p.ChildrenOnly,
		View:         p.View,
		ServiceName:  p.ServiceName,
	}
	cl.ID = uuid.Must(uuid.FromString(p.ID.String()))
	cl.SourceID = uuid.Must(uuid.FromString(p.SourceID.String()))
	cl.InheritedFrom = uuid.Must(uuid.FromString(p.InheritedFrom.String()))
	cl.ServiceID = uuid.Must(uuid.FromString(p.ServiceID.String()))
	cl.Attributes = make([]proto.ServiceAttribute, 0)
	for _, attr := range p.Attributes {
		a := proto.ServiceAttribute{
			Name:  attr.Name,
			Value: attr.Value,
		}
		cl.Attributes = append(cl.Attributes, a)
	}
	cl.Instances = make([]PropertyInstance, len(p.Instances))
	copy(cl.Instances, p.Instances)

	return &cl
}

func (p *PropertyService) MakeAction() Action {
	a := Action{
		Property: proto.Property{
			InstanceID:       p.GetID(),
			SourceInstanceID: p.GetSourceInstance(),
			SourceType:       p.GetSourceType(),
			IsInherited:      p.GetIsInherited(),
			InheritedFrom:    p.GetSource(),
			Type:             p.GetType(),
			Inheritance:      p.hasInheritance(),
			ChildrenOnly:     p.isChildrenOnly(),
			View:             p.GetView(),
			Service: &proto.PropertyService{
				ID:   p.ServiceID.String(),
				Name: p.ServiceName,
			},
		},
	}
	a.Property.Service.Attributes = make([]proto.ServiceAttribute, len(p.Attributes))
	for i, attr := range p.Attributes {
		t := proto.ServiceAttribute{
			Name:  attr.Name,
			Value: attr.Value,
		}
		a.Property.Service.Attributes[i] = t
	}
	return a
}

//
// System
type PropertySystem struct {
	ID            uuid.UUID
	SourceID      uuid.UUID
	SourceType    string
	Inherited     bool
	InheritedFrom uuid.UUID
	Inheritance   bool
	ChildrenOnly  bool
	View          string
	Key           string
	Value         string
	Instances     []PropertyInstance
}

func (p *PropertySystem) GetType() string {
	return "system"
}

func (p *PropertySystem) GetID() string {
	return p.ID.String()
}

func (p *PropertySystem) GetSource() string {
	return p.InheritedFrom.String()
}

func (p *PropertySystem) hasInheritance() bool {
	return p.Inheritance
}

func (p *PropertySystem) isChildrenOnly() bool {
	return p.ChildrenOnly
}

func (p *PropertySystem) GetSourceInstance() string {
	return p.SourceID.String()
}

func (p *PropertySystem) GetSourceType() string {
	return p.SourceType
}

func (p *PropertySystem) GetIsInherited() bool {
	return p.Inherited
}

func (p *PropertySystem) GetView() string {
	return p.View
}

func (p *PropertySystem) GetKey() string {
	return p.Key
}

func (p *PropertySystem) GetValue() string {
	return p.Value
}

func (p *PropertySystem) GetInstanceID(objType string, objID uuid.UUID, l *log.Logger) uuid.UUID {
	if !uuid.Equal(p.ID, uuid.Nil) {
		return p.ID
	}
	for _, instance := range p.Instances {
		if objType == instance.ObjectType && uuid.Equal(instance.ObjectID, objID) {
			l.Printf("tree.Property.GetInstanceID() found existing instance: %s\n", instance.InstanceID)
			return instance.InstanceID
		}
	}
	return uuid.Nil
}

func (p *PropertySystem) SetID(id uuid.UUID) {
	p.ID, _ = uuid.FromString(id.String())
}

func (p *PropertySystem) Equal(id uuid.UUID) bool {
	return uuid.Equal(p.ID, id)
}

func (p *PropertySystem) clearInstances() {
	p.Instances = nil
}

func (p *PropertySystem) SetInheritedFrom(id uuid.UUID) {
	p.InheritedFrom, _ = uuid.FromString(id.String())
}

func (p *PropertySystem) SetInherited(inherited bool) {
	p.Inherited = inherited
}

func (p *PropertySystem) SetSourceID(id uuid.UUID) {
	p.SourceID, _ = uuid.FromString(id.String())
}

func (p *PropertySystem) SetSourceType(s string) {
	p.SourceType = s
}

func (p PropertySystem) Clone() Property {
	cl := PropertySystem{
		SourceType:   p.SourceType,
		Inherited:    p.Inherited,
		Inheritance:  p.Inheritance,
		ChildrenOnly: p.ChildrenOnly,
		View:         p.View,
		Key:          p.Key,
		Value:        p.Value,
	}
	cl.ID, _ = uuid.FromString(p.ID.String())
	cl.SourceID, _ = uuid.FromString(p.SourceID.String())
	cl.InheritedFrom, _ = uuid.FromString(p.InheritedFrom.String())
	cl.Instances = make([]PropertyInstance, len(p.Instances))
	copy(cl.Instances, p.Instances)

	return &cl
}

func (p *PropertySystem) MakeAction() Action {
	return Action{
		Property: proto.Property{
			InstanceID:       p.GetID(),
			SourceInstanceID: p.GetSourceInstance(),
			SourceType:       p.GetSourceType(),
			IsInherited:      p.GetIsInherited(),
			InheritedFrom:    p.GetSource(),
			Type:             p.GetType(),
			Inheritance:      p.hasInheritance(),
			ChildrenOnly:     p.isChildrenOnly(),
			View:             p.GetView(),
			System: &proto.PropertySystem{
				Name:  p.Key,
				Value: p.Value,
			},
		},
	}
}

//
// Oncall
type PropertyOncall struct {
	ID            uuid.UUID
	SourceID      uuid.UUID
	SourceType    string
	OncallID      uuid.UUID
	Inherited     bool
	InheritedFrom uuid.UUID
	Inheritance   bool
	ChildrenOnly  bool
	View          string
	Name          string
	Number        string
	Instances     []PropertyInstance
}

func (p *PropertyOncall) GetType() string {
	return "oncall"
}

func (p *PropertyOncall) GetID() string {
	return p.ID.String()
}

func (p *PropertyOncall) GetSource() string {
	return p.InheritedFrom.String()
}

func (p *PropertyOncall) hasInheritance() bool {
	return p.Inheritance
}

func (p *PropertyOncall) isChildrenOnly() bool {
	return p.ChildrenOnly
}

func (p *PropertyOncall) GetSourceInstance() string {
	return p.SourceID.String()
}

func (p *PropertyOncall) GetSourceType() string {
	return p.SourceType
}

func (p *PropertyOncall) GetIsInherited() bool {
	return p.Inherited
}

func (p *PropertyOncall) GetView() string {
	return p.View
}

func (p *PropertyOncall) GetKey() string {
	return p.OncallID.String()
}

func (p *PropertyOncall) GetValue() string {
	return p.Name + p.Number
}

func (p *PropertyOncall) GetName() string {
	return p.Name
}

func (p *PropertyOncall) GetNumber() string {
	return p.Number
}

func (p *PropertyOncall) GetInstanceID(objType string, objID uuid.UUID, l *log.Logger) uuid.UUID {
	if !uuid.Equal(p.ID, uuid.Nil) {
		return p.ID
	}
	for _, instance := range p.Instances {
		if objType == instance.ObjectType && uuid.Equal(instance.ObjectID, objID) {
			l.Printf("tree.Property.GetInstanceID() found existing instance: %s\n", instance.InstanceID)
			return instance.InstanceID
		}
	}
	return uuid.Nil
}

func (p *PropertyOncall) SetID(id uuid.UUID) {
	p.ID, _ = uuid.FromString(id.String())
}

func (p *PropertyOncall) Equal(id uuid.UUID) bool {
	return uuid.Equal(p.ID, id)
}

func (p *PropertyOncall) clearInstances() {
	p.Instances = nil
}

func (p *PropertyOncall) SetInheritedFrom(id uuid.UUID) {
	p.InheritedFrom, _ = uuid.FromString(id.String())
}

func (p *PropertyOncall) SetInherited(inherited bool) {
	p.Inherited = inherited
}

func (p *PropertyOncall) SetSourceID(id uuid.UUID) {
	p.SourceID, _ = uuid.FromString(id.String())
}

func (p *PropertyOncall) SetSourceType(s string) {
	p.SourceType = s
}

func (p PropertyOncall) Clone() Property {
	cl := PropertyOncall{
		SourceType:   p.SourceType,
		Inherited:    p.Inherited,
		Inheritance:  p.Inheritance,
		ChildrenOnly: p.ChildrenOnly,
		View:         p.View,
		Name:         p.Name,
		Number:       p.Number,
	}
	cl.ID, _ = uuid.FromString(p.ID.String())
	cl.SourceID, _ = uuid.FromString(p.SourceID.String())
	cl.OncallID, _ = uuid.FromString(p.OncallID.String())
	cl.InheritedFrom, _ = uuid.FromString(p.InheritedFrom.String())
	cl.Instances = make([]PropertyInstance, len(p.Instances))
	copy(cl.Instances, p.Instances)

	return &cl
}

func (p *PropertyOncall) MakeAction() Action {
	return Action{
		Property: proto.Property{
			InstanceID:       p.GetID(),
			SourceInstanceID: p.GetSourceInstance(),
			SourceType:       p.GetSourceType(),
			IsInherited:      p.GetIsInherited(),
			InheritedFrom:    p.GetSource(),
			Type:             p.GetType(),
			Inheritance:      p.hasInheritance(),
			ChildrenOnly:     p.isChildrenOnly(),
			View:             p.GetView(),
			Oncall: &proto.PropertyOncall{
				ID:     p.OncallID.String(),
				Name:   p.Name,
				Number: p.Number,
			},
		},
	}
}

func isDupe(o, n Property) (bool, bool, Property) {
	var dupe, deleteOK bool
	var prop Property

	if o.GetKey() == n.GetKey() {
		// not allowed to replace view any with a more
		// specific view or vice versa. Replacing any with any
		// is fine
		if (o.GetView() == `any` && n.GetView() != `any`) ||
			(o.GetView() != `any` && n.GetView() == `any`) {
			// not actually a dupe, but trigger error path
			dupe = true
			deleteOK = false
		}
		// same view means we have a duplicate
		if o.GetView() == n.GetView() {
			dupe = true
			prop = o.Clone()
			// inherited properties can be deleted and replaced
			if o.GetIsInherited() {
				deleteOK = true
			}
		}
	}
	return dupe, deleteOK, prop
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
