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

func (teb *Bucket) SetProperty(p Property) {
	// if deleteOK is true, then prop is the property that can be
	// deleted
	if dupe, deleteOK, prop := teb.checkDuplicate(p); dupe && !deleteOK {
		teb.Fault.Error <- &Error{Action: `duplicate_set_property`}
		return
	} else if dupe && deleteOK {
		srcUUID, _ := uuid.FromString(prop.GetSourceInstance())
		switch prop.GetType() {
		case `custom`:
			cstUUID, _ := uuid.FromString(prop.GetKey())
			teb.deletePropertyInherited(&PropertyCustom{
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
			teb.deletePropertyInherited(&PropertyService{
				SourceID:  srcUUID,
				View:      prop.GetView(),
				Inherited: true,
				Service:   prop.GetKey(),
			})
		case `system`:
			teb.deletePropertyInherited(&PropertySystem{
				SourceID:  srcUUID,
				View:      prop.GetView(),
				Inherited: true,
				Key:       prop.GetKey(),
				Value:     prop.GetValue(),
			})
		case `oncall`:
			oncUUID, _ := uuid.FromString(prop.GetKey())
			teb.deletePropertyInherited(&PropertyOncall{
				SourceID:  srcUUID,
				View:      prop.GetView(),
				Inherited: true,
				OncallID:  oncUUID,
				Name:      prop.(*PropertyOncall).GetName(),
				Number:    prop.(*PropertyOncall).GetNumber(),
			})
		}
	}
	p.SetID(p.GetInstanceID(teb.Type, teb.ID, teb.log))
	if p.Equal(uuid.Nil) {
		p.SetID(uuid.NewV4())
	}
	// this property is the source instance
	p.SetInheritedFrom(teb.ID)
	p.SetInherited(false)
	p.SetSourceType(teb.Type)
	if i, e := uuid.FromString(p.GetID()); e == nil {
		p.SetSourceID(i)
	}
	// send a scrubbed copy down
	f := p.Clone()
	f.SetInherited(true)
	f.SetID(uuid.UUID{})
	if f.hasInheritance() {
		teb.setPropertyOnChildren(f)
	}
	// scrub instance startup information prior to storing
	p.clearInstances()
	teb.addProperty(p)
	teb.actionPropertyNew(p.MakeAction())
}

func (teb *Bucket) setPropertyInherited(p Property) {
	f := p.Clone()
	f.SetID(f.GetInstanceID(teb.Type, teb.ID, teb.log))
	if f.Equal(uuid.Nil) {
		f.SetID(uuid.NewV4())
	}
	f.clearInstances()

	if !f.GetIsInherited() {
		teb.Fault.Error <- &Error{
			Action: `bucket.setPropertyInherited on inherited=false`}
		return
	}
	if dupe, deleteOK, _ := teb.checkDuplicate(p); dupe && deleteOK {
		// we received an inherited SetProperty from above us in the
		// tree for a property that is duplicate, but we are not the
		// source of the duplicate -> corrupt tree
		teb.Fault.Error <- &Error{
			Action: `bucket.setPropertyInherited corruption detected`}
		return
	} else if dupe && !deleteOK {
		// we received an inherited SetProperty from above us in the
		// tree for a property that is duplicate; we have a locally
		// set property -> stop inheritance, no error
		return
	}
	teb.addProperty(f)
	p.SetID(uuid.UUID{})
	teb.setPropertyOnChildren(p)
	teb.actionPropertyNew(f.MakeAction())
}

func (teb *Bucket) setPropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	for child := range teb.Children {
		wg.Add(1)
		go func(stp Property, c string) {
			defer wg.Done()
			teb.Children[c].setPropertyInherited(stp)
		}(p, child)
	}
	wg.Wait()
}

func (teb *Bucket) addProperty(p Property) {
	switch p.GetType() {
	case `custom`:
		teb.PropertyCustom[p.GetID()] = p
	case `system`:
		teb.PropertySystem[p.GetID()] = p
	case `service`:
		teb.PropertyService[p.GetID()] = p
	case `oncall`:
		teb.PropertyOncall[p.GetID()] = p
	default:
		teb.Fault.Error <- &Error{Action: `bucket.addProperty unknown type`}
	}
}

//
// Propertier:> Update Property

func (teb *Bucket) UpdateProperty(p Property) {
	if !teb.verifySourceInstance(
		p.GetSourceInstance(),
		p.GetType(),
	) {
		teb.Fault.Error <- &Error{Action: `update_property_on_non_source`}
		return
	}

	// keep a copy for ourselves, no shared pointers
	p.SetInheritedFrom(teb.ID)
	p.SetSourceType(teb.Type)
	p.SetInherited(true)
	f := p.Clone()
	f.SetInherited(false)
	if teb.switchProperty(f) {
		teb.updatePropertyOnChildren(p)
	}
}

func (teb *Bucket) updatePropertyInherited(p Property) {
	// keep a copy for ourselves, no shared pointers
	f := p.Clone()
	if !f.GetIsInherited() {
		teb.Fault.Error <- &Error{
			Action: `bucket.updatePropertyInherited on inherited=false`}
		return
	}
	if teb.switchProperty(f) {
		teb.updatePropertyOnChildren(p)
	}
}

func (teb *Bucket) updatePropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	for child := range teb.Children {
		wg.Add(1)
		go func(stp Property, c string) {
			defer wg.Done()
			teb.Children[c].updatePropertyInherited(stp)
		}(p, child)
	}
	wg.Wait()
}

func (teb *Bucket) switchProperty(p Property) bool {
	uid := teb.findIDForSource(
		p.GetSourceInstance(),
		p.GetType(),
	)
	if uid == `` {
		// we do not have the property for which we received an update
		if dupe, deleteOK, _ := teb.checkDuplicate(p); dupe && !deleteOK {
			// the update is duplicate to an property for which we
			// have the source instance, ie we just received an update
			// for which we have an overwrite. Ignore it and do not
			// inherit it further down
			return false
		}
		teb.Fault.Error <- &Error{
			Action: `bucket.switchProperty property not found`}
		return false
	}
	updID, _ := uuid.FromString(uid)
	p.SetID(updID)
	curr := teb.getCurrentProperty(p)
	if curr == nil {
		return false
	}
	teb.addProperty(p)
	teb.actionPropertyUpdate(p.MakeAction())

	if !p.hasInheritance() && curr.hasInheritance() {
		// replacing inheritance with !inheritance:
		// call deletePropertyOnChildren(curr) to clean up
		srcUUID, _ := uuid.FromString(curr.GetSourceInstance())
		switch curr.GetType() {
		case `custom`:
			cstUUID, _ := uuid.FromString(curr.GetKey())
			teb.deletePropertyOnChildren(&PropertyCustom{
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
			teb.deletePropertyOnChildren(&PropertyService{
				SourceID:    srcUUID,
				View:        curr.GetView(),
				Inherited:   true,
				Service:     curr.GetKey(),
				Inheritance: true,
			})
		case `system`:
			teb.deletePropertyOnChildren(&PropertySystem{
				SourceID:    srcUUID,
				View:        curr.GetView(),
				Inherited:   true,
				Key:         curr.GetKey(),
				Value:       curr.GetValue(),
				Inheritance: true,
			})
		case `oncall`:
			oncUUID, _ := uuid.FromString(curr.GetKey())
			teb.deletePropertyOnChildren(&PropertyOncall{
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
		teb.setPropertyOnChildren(t)
	}
	return p.hasInheritance() && curr.hasInheritance()
}

func (teb *Bucket) getCurrentProperty(p Property) Property {
	switch p.GetType() {
	case `custom`:
		return teb.PropertyCustom[p.GetID()].Clone()
	case `system`:
		return teb.PropertySystem[p.GetID()].Clone()
	case `service`:
		return teb.PropertyService[p.GetID()].Clone()
	case `oncall`:
		return teb.PropertyOncall[p.GetID()].Clone()
	}
	teb.Fault.Error <- &Error{
		Action: `bucket.getCurrentProperty unknown type`}
	return nil
}

//
// Propertier:> Delete Property

func (teb *Bucket) DeleteProperty(p Property) {
	if !teb.verifySourceInstance(
		p.GetSourceInstance(),
		p.GetType(),
	) {
		teb.Fault.Error <- &Error{Action: `bucket.DeleteProperty on !source`}
		return
	}

	var flow Property
	resync := false
	delID := teb.findIDForSource(
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
			delProp = teb.PropertyCustom[delID]
		case `system`:
			delProp = teb.PropertySystem[delID]
		case `service`:
			delProp = teb.PropertyService[delID]
		case `oncall`:
			delProp = teb.PropertyOncall[delID]
		}
		resync, _, flow = teb.Parent.(Propertier).checkDuplicate(
			delProp,
		)
	}

	p.SetInherited(false)
	if teb.rmProperty(p) {
		p.SetInherited(true)
		teb.deletePropertyOnChildren(p)
	}

	// now that the property is deleted from us and our children,
	// request resync if required
	if resync {
		teb.Parent.(Propertier).resyncProperty(
			flow.GetSourceInstance(),
			p.GetType(),
			teb.ID.String(),
		)
	}
}

func (teb *Bucket) deletePropertyInherited(p Property) {
	if teb.rmProperty(p) {
		teb.deletePropertyOnChildren(p)
	}
}

func (teb *Bucket) deletePropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	for child := range teb.Children {
		wg.Add(1)
		go func(stp Property, c string) {
			defer wg.Done()
			teb.Children[c].deletePropertyInherited(stp)
		}(p, child)
	}
	wg.Wait()
}

func (teb *Bucket) deletePropertyAllInherited() {
	for _, p := range teb.PropertyCustom {
		if !p.GetIsInherited() {
			continue
		}
		teb.deletePropertyInherited(p.Clone())
	}
	for _, p := range teb.PropertySystem {
		if !p.GetIsInherited() {
			continue
		}
		teb.deletePropertyInherited(p.Clone())
	}
	for _, p := range teb.PropertyService {
		if !p.GetIsInherited() {
			continue
		}
		teb.deletePropertyInherited(p.Clone())
	}
	for _, p := range teb.PropertyOncall {
		if !p.GetIsInherited() {
			continue
		}
		teb.deletePropertyInherited(p.Clone())
	}
}

func (teb *Bucket) deletePropertyAllLocal() {
	for _, p := range teb.PropertyCustom {
		if p.GetIsInherited() {
			continue
		}
		teb.DeleteProperty(p.Clone())
	}
	for _, p := range teb.PropertySystem {
		if p.GetIsInherited() {
			continue
		}
		teb.DeleteProperty(p.Clone())
	}
	for _, p := range teb.PropertyService {
		if p.GetIsInherited() {
			continue
		}
		teb.DeleteProperty(p.Clone())
	}
	for _, p := range teb.PropertyOncall {
		if p.GetIsInherited() {
			continue
		}
		teb.DeleteProperty(p.Clone())
	}
}

func (teb *Bucket) rmProperty(p Property) bool {
	delID := teb.findIDForSource(
		p.GetSourceInstance(),
		p.GetType(),
	)
	if delID == `` {
		// we do not have the property for which we received a delete
		if dupe, deleteOK, _ := teb.checkDuplicate(p); dupe && !deleteOK {
			// the delete is duplicate to a property for which we
			// have the source instance, ie we just received a delete
			// for which we have an overwrite. Ignore it and do not
			// inherit it further down
			return false
		}

		teb.Fault.Error <- &Error{
			Action: `bucket.rmProperty property not found`}
		return false
	}

	hasInheritance := false
	switch p.GetType() {
	case `custom`:
		teb.actionPropertyDelete(
			teb.PropertyCustom[delID].MakeAction(),
		)
		hasInheritance = teb.PropertyCustom[delID].hasInheritance()
		delete(teb.PropertyCustom, delID)
	case `service`:
		teb.actionPropertyDelete(
			teb.PropertyService[delID].MakeAction(),
		)
		hasInheritance = teb.PropertyService[delID].hasInheritance()
		delete(teb.PropertyService, delID)
	case `system`:
		teb.actionPropertyDelete(
			teb.PropertySystem[delID].MakeAction(),
		)
		hasInheritance = teb.PropertySystem[delID].hasInheritance()
		delete(teb.PropertySystem, delID)
	case `oncall`:
		teb.actionPropertyDelete(
			teb.PropertyOncall[delID].MakeAction(),
		)
		hasInheritance = teb.PropertyOncall[delID].hasInheritance()
		delete(teb.PropertyOncall, delID)
	default:
		teb.Fault.Error <- &Error{Action: `bucket.rmProperty unknown type`}
		return false
	}
	return hasInheritance
}

//
// Propertier:> Utility

// used to verify this is a source instance
func (teb *Bucket) verifySourceInstance(id, prop string) bool {
	switch prop {
	case `custom`:
		if _, ok := teb.PropertyCustom[id]; !ok {
			goto bailout
		}
		return teb.PropertyCustom[id].GetSourceInstance() == id
	case `service`:
		if _, ok := teb.PropertyService[id]; !ok {
			goto bailout
		}
		return teb.PropertyService[id].GetSourceInstance() == id
	case `system`:
		if _, ok := teb.PropertySystem[id]; !ok {
			goto bailout
		}
		return teb.PropertySystem[id].GetSourceInstance() == id
	case `oncall`:
		if _, ok := teb.PropertyOncall[id]; !ok {
			goto bailout
		}
		return teb.PropertyOncall[id].GetSourceInstance() == id
	}

bailout:
	teb.Fault.Error <- &Error{
		Action: `bucket.verifySourceInstance not found`}
	return false
}

//
func (teb *Bucket) findIDForSource(source, prop string) string {
	switch prop {
	case `custom`:
		for id := range teb.PropertyCustom {
			if teb.PropertyCustom[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	case `system`:
		for id := range teb.PropertySystem {
			if teb.PropertySystem[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	case `service`:
		for id := range teb.PropertyService {
			if teb.PropertyService[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	case `oncall`:
		for id := range teb.PropertyOncall {
			if teb.PropertyOncall[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	}
	return ``
}

//
func (teb *Bucket) resyncProperty(srcID, pType, childID string) {
	pID := teb.findIDForSource(srcID, pType)
	if pID == `` {
		return
	}

	var f Property
	switch pType {
	case `custom`:
		f = teb.PropertyCustom[pID].(*PropertyCustom).Clone()
	case `oncall`:
		f = teb.PropertyOncall[pID].(*PropertyOncall).Clone()
	case `service`:
		f = teb.PropertyService[pID].(*PropertyService).Clone()
	case `system`:
		f = teb.PropertySystem[pID].(*PropertySystem).Clone()
	}
	if !f.hasInheritance() {
		return
	}
	f.SetInherited(true)
	f.SetID(uuid.UUID{})
	f.clearInstances()
	teb.Children[childID].setPropertyInherited(f)
}

// when a child attaches, it calls self.Parent.syncProperty(self.ID)
// to get get all properties of that part of the tree
func (teb *Bucket) syncProperty(childID string) {
customloop:
	for prop := range teb.PropertyCustom {
		if !teb.PropertyCustom[prop].hasInheritance() {
			continue customloop
		}
		f := teb.PropertyCustom[prop].(*PropertyCustom).Clone()
		f.SetInherited(true)
		f.SetID(uuid.UUID{})
		f.clearInstances()
		teb.Children[childID].setPropertyInherited(f)
	}
oncallloop:
	for prop := range teb.PropertyOncall {
		if !teb.PropertyOncall[prop].hasInheritance() {
			continue oncallloop
		}
		f := teb.PropertyOncall[prop].(*PropertyOncall).Clone()
		f.SetInherited(true)
		f.SetID(uuid.UUID{})
		f.clearInstances()
		teb.Children[childID].setPropertyInherited(f)
	}
serviceloop:
	for prop := range teb.PropertyService {
		if !teb.PropertyService[prop].hasInheritance() {
			continue serviceloop
		}
		f := teb.PropertyService[prop].(*PropertyService).Clone()
		f.SetInherited(true)
		f.SetID(uuid.UUID{})
		f.clearInstances()
		teb.Children[childID].setPropertyInherited(f)
	}
systemloop:
	for prop := range teb.PropertySystem {
		if !teb.PropertySystem[prop].hasInheritance() {
			continue systemloop
		}
		f := teb.PropertySystem[prop].(*PropertySystem).Clone()
		f.SetInherited(true)
		f.SetID(uuid.UUID{})
		f.clearInstances()
		teb.Children[childID].setPropertyInherited(f)
	}
}

// function to be used by a child to check if the parent has a
// specific Property
func (teb *Bucket) checkProperty(propType string, propID string) bool {
	switch propType {
	case "custom":
		if _, ok := teb.PropertyCustom[propID]; ok {
			return true
		}
	case "service":
		if _, ok := teb.PropertyService[propID]; ok {
			return true
		}
	case "system":
		if _, ok := teb.PropertySystem[propID]; ok {
			return true
		}
	case "oncall":
		if _, ok := teb.PropertyOncall[propID]; ok {
			return true
		}
	}
	return false
}

// Checks if this property is already defined on this node, and
// whether it was inherited, ie. can be deleted so it can be
// overwritten
func (teb *Bucket) checkDuplicate(p Property) (bool, bool, Property) {
	var dupe, deleteOK bool
	var prop Property

propswitch:
	switch p.GetType() {
	case "custom":
		for _, pVal := range teb.PropertyCustom {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "service":
		for _, pVal := range teb.PropertyService {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "oncall":
		for _, pVal := range teb.PropertyOncall {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "system":
		for _, pVal := range teb.PropertySystem {
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
		teb.Fault.Error <- &Error{Action: `bucket.checkDuplicate unknown type`}
		dupe = true
		deleteOK = false
	}
	return dupe, deleteOK, prop
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
