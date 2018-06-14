/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import (
	"sync"

	"github.com/satori/go.uuid"
)

// Implementation of the `Propertier` interface

//
// Propertier:> Add Property

func (teg *Group) SetProperty(p Property) {
	// if deleteOK is true, then prop is the property that can be
	// deleted
	if dupe, deleteOK, prop := teg.checkDuplicate(p); dupe && !deleteOK {
		teg.Fault.Error <- &Error{Action: `duplicate_set_property`}
		return
	} else if dupe && deleteOK {
		srcUUID, _ := uuid.FromString(prop.GetSourceInstance())
		switch prop.GetType() {
		case `custom`:
			cstUUID, _ := uuid.FromString(prop.GetKey())
			teg.deletePropertyInherited(&PropertyCustom{
				SourceID:  srcUUID,
				View:      prop.GetView(),
				Inherited: true,
				CustomID:  cstUUID,
				Key:       prop.(*PropertyCustom).GetKeyField(),
				Value:     prop.(*PropertyCustom).GetValueField(),
			})
		case `service`:
			// GetValue for serviceproperty returns the uuid to never
			// match, we do not set it
			teg.deletePropertyInherited(&PropertyService{
				SourceID:  srcUUID,
				View:      prop.GetView(),
				Inherited: true,
				Service:   prop.GetKey(),
			})
		case `system`:
			teg.deletePropertyInherited(&PropertySystem{
				SourceID:  srcUUID,
				View:      prop.GetView(),
				Inherited: true,
				Key:       prop.GetKey(),
				Value:     prop.GetValue(),
			})
		case `oncall`:
			oncUUID, _ := uuid.FromString(prop.GetKey())
			teg.deletePropertyInherited(&PropertyOncall{
				SourceID:  srcUUID,
				View:      prop.GetView(),
				Inherited: true,
				OncallID:  oncUUID,
				Name:      prop.(*PropertyOncall).GetName(),
				Number:    prop.(*PropertyOncall).GetNumber(),
			})
		}
	}
	p.SetID(p.GetInstanceID(teg.Type, teg.ID, teg.log))
	if p.Equal(uuid.Nil) {
		p.SetID(uuid.Must(uuid.NewV4()))
	}
	// this property is the source instance
	p.SetInheritedFrom(teg.ID)
	p.SetInherited(false)
	p.SetSourceType(teg.Type)
	if i, e := uuid.FromString(p.GetID()); e == nil {
		p.SetSourceID(i)
	}
	// send a scrubbed copy down
	f := p.Clone()
	f.SetInherited(true)
	f.SetID(uuid.UUID{})
	if f.hasInheritance() {
		teg.setPropertyOnChildren(f)
	}
	// scrub instance startup information prior to storing
	p.clearInstances()
	teg.addProperty(p)
	teg.actionPropertyNew(p.MakeAction())
}

func (teg *Group) setPropertyInherited(p Property) {
	f := p.Clone()
	f.SetID(f.GetInstanceID(teg.Type, teg.ID, teg.log))
	if f.Equal(uuid.Nil) {
		f.SetID(uuid.Must(uuid.NewV4()))
	}
	f.clearInstances()

	if !f.GetIsInherited() {
		teg.Fault.Error <- &Error{
			Action: `group.setPropertyInherited on inherited=false`}
		return
	}
	if dupe, deleteOK, _ := teg.checkDuplicate(p); dupe && deleteOK {
		// we received an inherited SetProperty from above us in the
		// tree for a property that is duplicate, but we are not the
		// source of the duplicate -> corrupt tree
		teg.Fault.Error <- &Error{
			Action: `group.setPropertyInherited corruption detected`}
		return
	} else if dupe && !deleteOK {
		// we received an inherited SetProperty from above us in the
		// tree for a property that is duplicate; we have a locally
		// set property -> stop inheritance, no error
		return
	}
	teg.addProperty(f)
	p.SetID(uuid.UUID{})
	teg.setPropertyOnChildren(p)
	teg.actionPropertyNew(f.MakeAction())
}

func (teg *Group) setPropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	for child := range teg.Children {
		wg.Add(1)
		go func(stp Property, c string) {
			defer wg.Done()
			teg.Children[c].setPropertyInherited(stp)
		}(p, child)
	}
	wg.Wait()
}

func (teg *Group) addProperty(p Property) {
	teg.hasUpdate = true
	switch p.GetType() {
	case `custom`:
		teg.PropertyCustom[p.GetID()] = p
	case `system`:
		teg.PropertySystem[p.GetID()] = p
	case `service`:
		teg.PropertyService[p.GetID()] = p
	case `oncall`:
		teg.PropertyOncall[p.GetID()] = p
	default:
		teg.hasUpdate = false
		teg.Fault.Error <- &Error{Action: `group.addProperty unknown type`}
	}
}

//
// Propertier:> Update Property

func (teg *Group) UpdateProperty(p Property) {
	if !teg.verifySourceInstance(
		p.GetSourceInstance(),
		p.GetType(),
	) {
		teg.Fault.Error <- &Error{Action: `update_property_on_non_source`}
		return
	}

	// keep a copy for ourselves, no shared pointers
	p.SetInheritedFrom(teg.ID)
	p.SetSourceType(teg.Type)
	p.SetInherited(true)
	f := p.Clone()
	f.SetInherited(false)
	if teg.switchProperty(f) {
		teg.updatePropertyOnChildren(p)
	}
}

func (teg *Group) updatePropertyInherited(p Property) {
	// keep a copy for ourselves, no shared pointers
	f := p.Clone()
	if !f.GetIsInherited() {
		teg.Fault.Error <- &Error{
			Action: `group.updatePropertyInherited on inherited=false`}
		return
	}
	if teg.switchProperty(f) {
		teg.updatePropertyOnChildren(p)
	}
}

func (teg *Group) updatePropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	for child := range teg.Children {
		wg.Add(1)
		go func(stp Property, c string) {
			defer wg.Done()
			teg.Children[c].updatePropertyInherited(stp)
		}(p, child)
	}
	wg.Wait()
}

func (teg *Group) switchProperty(p Property) bool {
	uid := teg.findIDForSource(
		p.GetSourceInstance(),
		p.GetType(),
	)
	if uid == `` {
		// we do not have the property for which we received an update
		if dupe, deleteOK, _ := teg.checkDuplicate(p); dupe && !deleteOK {
			// the update is duplicate to an property for which we
			// have the source instance, ie we just received an update
			// for which we have an overwrite. Ignore it and do not
			// inherit it further down
			return false
		}
		teg.Fault.Error <- &Error{
			Action: `group.switchProperty property not found`}
		return false
	}
	updID, _ := uuid.FromString(uid)
	p.SetID(updID)
	curr := teg.getCurrentProperty(p)
	if curr == nil {
		return false
	}
	teg.addProperty(p)
	teg.actionPropertyUpdate(p.MakeAction())

	if !p.hasInheritance() && curr.hasInheritance() {
		// replacing inheritance with !inheritance:
		// call deletePropertyOnChildren(curr) to clean up
		srcUUID, _ := uuid.FromString(curr.GetSourceInstance())
		switch curr.GetType() {
		case `custom`:
			cstUUID, _ := uuid.FromString(curr.GetKey())
			teg.deletePropertyOnChildren(&PropertyCustom{
				SourceID:    srcUUID,
				View:        curr.GetView(),
				Inherited:   true,
				CustomID:    cstUUID,
				Key:         curr.(*PropertyCustom).GetKeyField(),
				Value:       curr.(*PropertyCustom).GetValueField(),
				Inheritance: true,
			})
		case `service`:
			// GetValue for serviceproperty returns the uuid to never
			// match, we do not set it
			teg.deletePropertyOnChildren(&PropertyService{
				SourceID:    srcUUID,
				View:        curr.GetView(),
				Inherited:   true,
				Service:     curr.GetKey(),
				Inheritance: true,
			})
		case `system`:
			teg.deletePropertyOnChildren(&PropertySystem{
				SourceID:    srcUUID,
				View:        curr.GetView(),
				Inherited:   true,
				Key:         curr.GetKey(),
				Value:       curr.GetValue(),
				Inheritance: true,
			})
		case `oncall`:
			oncUUID, _ := uuid.FromString(curr.GetKey())
			teg.deletePropertyOnChildren(&PropertyOncall{
				SourceID:    srcUUID,
				View:        curr.GetView(),
				Inherited:   true,
				OncallID:    oncUUID,
				Name:        curr.(*PropertyOncall).GetName(),
				Number:      curr.(*PropertyOncall).GetNumber(),
				Inheritance: true,
			})
		}
	}
	if p.hasInheritance() && !curr.hasInheritance() {
		// replacing !inheritance with inheritance:
		// call setPropertyonChildren(p) to propagate
		t := p.Clone()
		t.SetInherited(true)
		teg.setPropertyOnChildren(t)
	}
	return p.hasInheritance() && curr.hasInheritance()
}

func (teg *Group) getCurrentProperty(p Property) Property {
	switch p.GetType() {
	case `custom`:
		return teg.PropertyCustom[p.GetID()].Clone()
	case `system`:
		return teg.PropertySystem[p.GetID()].Clone()
	case `service`:
		return teg.PropertyService[p.GetID()].Clone()
	case `oncall`:
		return teg.PropertyOncall[p.GetID()].Clone()
	}
	teg.Fault.Error <- &Error{
		Action: `group.getCurrentProperty unknown type`}
	return nil
}

//
// Propertier:> Delete Property

func (teg *Group) DeleteProperty(p Property) {
	if !teg.verifySourceInstance(
		p.GetSourceInstance(),
		p.GetType(),
	) {
		teg.Fault.Error <- &Error{Action: `group.DeleteProperty on !source`}
		return
	}

	var flow Property
	resync := false
	delID := teg.findIDForSource(
		p.GetSourceInstance(),
		p.GetType(),
	)
	if delID != `` {
		// this is a delete for a locally set property. It might be a
		// delete for an overwrite property, in which case we need to
		// ask the parent to sync it to us again.
		// If it was an overwrite, the parent should have a property
		// we would consider a dupe if it were to be passed down to
		// us.
		// If p is considered a dupe, then flow is set to the prop we
		// need to inherit.
		var delProp Property
		switch p.GetType() {
		case `custom`:
			delProp = teg.PropertyCustom[delID]
		case `system`:
			delProp = teg.PropertySystem[delID]
		case `service`:
			delProp = teg.PropertyService[delID]
		case `oncall`:
			delProp = teg.PropertyOncall[delID]
		}
		resync, _, flow = teg.Parent.(Propertier).checkDuplicate(
			delProp,
		)
	}

	p.SetInherited(false)
	if teg.rmProperty(p) {
		p.SetInherited(true)
		teg.deletePropertyOnChildren(p)
	}

	// now that the property is deleted from us and our children,
	// request resync if required
	if resync {
		teg.Parent.(Propertier).resyncProperty(
			flow.GetSourceInstance(),
			p.GetType(),
			teg.ID.String(),
		)
	}
}

func (teg *Group) deletePropertyInherited(p Property) {
	if teg.rmProperty(p) {
		teg.deletePropertyOnChildren(p)
	}
}

func (teg *Group) deletePropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	for child := range teg.Children {
		wg.Add(1)
		go func(stp Property, c string) {
			defer wg.Done()
			teg.Children[c].deletePropertyInherited(stp)
		}(p, child)
	}
	wg.Wait()
}

func (teg *Group) deletePropertyAllInherited() {
	for _, p := range teg.PropertyCustom {
		if !p.GetIsInherited() {
			continue
		}
		teg.deletePropertyInherited(p.Clone())
	}
	for _, p := range teg.PropertySystem {
		if !p.GetIsInherited() {
			continue
		}
		teg.deletePropertyInherited(p.Clone())
	}
	for _, p := range teg.PropertyService {
		if !p.GetIsInherited() {
			continue
		}
		teg.deletePropertyInherited(p.Clone())
	}
	for _, p := range teg.PropertyOncall {
		if !p.GetIsInherited() {
			continue
		}
		teg.deletePropertyInherited(p.Clone())
	}
}

func (teg *Group) deletePropertyAllLocal() {
	for _, p := range teg.PropertyCustom {
		if p.GetIsInherited() {
			continue
		}
		teg.DeleteProperty(p.Clone())
	}
	for _, p := range teg.PropertySystem {
		if p.GetIsInherited() {
			continue
		}
		teg.DeleteProperty(p.Clone())
	}
	for _, p := range teg.PropertyService {
		if p.GetIsInherited() {
			continue
		}
		teg.DeleteProperty(p.Clone())
	}
	for _, p := range teg.PropertyOncall {
		if p.GetIsInherited() {
			continue
		}
		teg.DeleteProperty(p.Clone())
	}
}

func (teg *Group) rmProperty(p Property) bool {
	delID := teg.findIDForSource(
		p.GetSourceInstance(),
		p.GetType(),
	)
	if delID == `` {
		// we do not have the property for which we received a delete
		if dupe, deleteOK, _ := teg.checkDuplicate(p); dupe && !deleteOK {
			// the delete is duplicate to a property for which we
			// have the source instance, ie we just received a delete
			// for which we have an overwrite. Ignore it and do not
			// inherit it further down
			return false
		}

		teg.Fault.Error <- &Error{
			Action: `group.rmProperty property not found`}
		return false
	}

	hasInheritance := false
	teg.hasUpdate = true
	switch p.GetType() {
	case `custom`:
		teg.actionPropertyDelete(
			teg.PropertyCustom[delID].MakeAction(),
		)
		hasInheritance = teg.PropertyCustom[delID].hasInheritance()
		delete(teg.PropertyCustom, delID)
	case `service`:
		teg.actionPropertyDelete(
			teg.PropertyService[delID].MakeAction(),
		)
		hasInheritance = teg.PropertyService[delID].hasInheritance()
		delete(teg.PropertyService, delID)
	case `system`:
		teg.actionPropertyDelete(
			teg.PropertySystem[delID].MakeAction(),
		)
		hasInheritance = teg.PropertySystem[delID].hasInheritance()
		delete(teg.PropertySystem, delID)
	case `oncall`:
		teg.actionPropertyDelete(
			teg.PropertyOncall[delID].MakeAction(),
		)
		hasInheritance = teg.PropertyOncall[delID].hasInheritance()
		delete(teg.PropertyOncall, delID)
	default:
		teg.hasUpdate = false
		teg.Fault.Error <- &Error{Action: `group.rmProperty unknown type`}
		return false
	}
	return hasInheritance
}

//
// Propertier:> Utility

//
func (teg *Group) verifySourceInstance(id, prop string) bool {
	switch prop {
	case `custom`:
		if _, ok := teg.PropertyCustom[id]; !ok {
			goto bailout
		}
		return teg.PropertyCustom[id].GetSourceInstance() == id
	case `service`:
		if _, ok := teg.PropertyService[id]; !ok {
			goto bailout
		}
		return teg.PropertyService[id].GetSourceInstance() == id
	case `system`:
		if _, ok := teg.PropertySystem[id]; !ok {
			goto bailout
		}
		return teg.PropertySystem[id].GetSourceInstance() == id
	case `oncall`:
		if _, ok := teg.PropertyOncall[id]; !ok {
			goto bailout
		}
		return teg.PropertyOncall[id].GetSourceInstance() == id
	}

bailout:
	teg.Fault.Error <- &Error{
		Action: `group.verifySourceInstance not found`}
	return false
}

//
func (teg *Group) findIDForSource(source, prop string) string {
	switch prop {
	case `custom`:
		for id := range teg.PropertyCustom {
			if teg.PropertyCustom[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	case `system`:
		for id := range teg.PropertySystem {
			if teg.PropertySystem[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	case `service`:
		for id := range teg.PropertyService {
			if teg.PropertyService[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	case `oncall`:
		for id := range teg.PropertyOncall {
			if teg.PropertyOncall[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	}
	return ``
}

//
func (teg *Group) resyncProperty(srcID, pType, childID string) {
	pID := teg.findIDForSource(srcID, pType)
	if pID == `` {
		return
	}

	var f Property
	switch pType {
	case `custom`:
		f = teg.PropertyCustom[pID].(*PropertyCustom).Clone()
	case `oncall`:
		f = teg.PropertyOncall[pID].(*PropertyOncall).Clone()
	case `service`:
		f = teg.PropertyService[pID].(*PropertyService).Clone()
	case `system`:
		f = teg.PropertySystem[pID].(*PropertySystem).Clone()
	}
	if !f.hasInheritance() {
		return
	}
	f.SetInherited(true)
	f.SetID(uuid.UUID{})
	f.clearInstances()
	teg.Children[childID].setPropertyInherited(f)
}

// when a child attaches, it calls self.Parent.syncProperty(self.ID)
// to get get all properties of that part of the tree
func (teg *Group) syncProperty(childID string) {
customloop:
	for prop := range teg.PropertyCustom {
		if !teg.PropertyCustom[prop].hasInheritance() {
			continue customloop
		}
		f := teg.PropertyCustom[prop].(*PropertyCustom).Clone()
		f.SetInherited(true)
		f.SetID(uuid.UUID{})
		f.clearInstances()
		teg.Children[childID].setPropertyInherited(f)
	}
oncallloop:
	for prop := range teg.PropertyOncall {
		if !teg.PropertyOncall[prop].hasInheritance() {
			continue oncallloop
		}
		f := teg.PropertyOncall[prop].(*PropertyOncall).Clone()
		f.SetInherited(true)
		f.SetID(uuid.UUID{})
		f.clearInstances()
		teg.Children[childID].setPropertyInherited(f)
	}
serviceloop:
	for prop := range teg.PropertyService {
		if !teg.PropertyService[prop].hasInheritance() {
			continue serviceloop
		}
		f := teg.PropertyService[prop].(*PropertyService).Clone()
		f.SetInherited(true)
		f.SetID(uuid.UUID{})
		f.clearInstances()
		teg.Children[childID].setPropertyInherited(f)
	}
systemloop:
	for prop := range teg.PropertySystem {
		if !teg.PropertySystem[prop].hasInheritance() {
			continue systemloop
		}
		f := teg.PropertySystem[prop].(*PropertySystem).Clone()
		f.SetInherited(true)
		f.SetID(uuid.UUID{})
		f.clearInstances()
		teg.Children[childID].setPropertyInherited(f)
	}
}

// function to be used by a child to check if the parent has a
// specific Property
func (teg *Group) checkProperty(propType string, propID string) bool {
	switch propType {
	case "custom":
		if _, ok := teg.PropertyCustom[propID]; ok {
			return true
		}
	case "service":
		if _, ok := teg.PropertyService[propID]; ok {
			return true
		}
	case "system":
		if _, ok := teg.PropertySystem[propID]; ok {
			return true
		}
	case "oncall":
		if _, ok := teg.PropertyOncall[propID]; ok {
			return true
		}
	}
	return false
}

// Checks if this property is already defined on this node, and
// whether it was inherited, ie. can be deleted so it can be
// overwritten
func (teg *Group) checkDuplicate(p Property) (bool, bool, Property) {
	var dupe, deleteOK bool
	var prop Property

propswitch:
	switch p.GetType() {
	case "custom":
		for _, pVal := range teg.PropertyCustom {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "service":
		for _, pVal := range teg.PropertyService {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "oncall":
		for _, pVal := range teg.PropertyOncall {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "system":
		for _, pVal := range teg.PropertySystem {
			// tags are only dupes if the value is the same as well
			if p.GetKey() != `tag` {
				dupe, deleteOK, prop = isDupe(pVal, p)
				if dupe {
					break propswitch
				}
			} else if p.GetValue() == pVal.GetValue() {
				// tag and same value, can be a dupe
				dupe, deleteOK, prop = isDupe(pVal, p)
				if dupe {
					break propswitch
				}
			}
			// tag + different value => pass
		}
	default:
		// trigger error path
		teg.Fault.Error <- &Error{Action: `group.checkDuplicate unknown type`}
		dupe = true
		deleteOK = false
	}
	return dupe, deleteOK, prop
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
