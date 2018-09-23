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

func (ter *Repository) SetProperty(p Property) {
	// if deleteOK is true, then prop is the property that can be
	// deleted
	if dupe, deleteOK, prop := ter.checkDuplicate(p); dupe && !deleteOK {
		ter.Fault.Error <- &Error{Action: `duplicate_set_property`}
		return
	} else if dupe && deleteOK {
		srcUUID, _ := uuid.FromString(prop.GetSourceInstance())
		switch prop.GetType() {
		case `custom`:
			cstUUID, _ := uuid.FromString(prop.GetKey())
			ter.deletePropertyInherited(&PropertyCustom{
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
			ter.deletePropertyInherited(&PropertyService{
				SourceID:  srcUUID,
				View:      prop.GetView(),
				Inherited: true,
				Service:   prop.GetKey(),
			})
		case `system`:
			ter.deletePropertyInherited(&PropertySystem{
				SourceID:  srcUUID,
				View:      prop.GetView(),
				Inherited: true,
				Key:       prop.GetKey(),
				Value:     prop.GetValue(),
			})
		case `oncall`:
			oncUUID, _ := uuid.FromString(prop.GetKey())
			ter.deletePropertyInherited(&PropertyOncall{
				SourceID:  srcUUID,
				View:      prop.GetView(),
				Inherited: true,
				OncallID:  oncUUID,
				Name:      prop.(*PropertyOncall).GetName(),
				Number:    prop.(*PropertyOncall).GetNumber(),
			})
		}
	}
	p.SetID(p.GetInstanceID(ter.Type, ter.ID, ter.log))
	if p.Equal(uuid.Nil) {
		p.SetID(uuid.Must(uuid.NewV4()))
	}
	// this property is the source instance
	p.SetInheritedFrom(ter.ID)
	p.SetInherited(false)
	p.SetSourceType(ter.Type)
	if i, e := uuid.FromString(p.GetID()); e == nil {
		p.SetSourceID(i)
	}
	// send a scrubbed copy down
	f := p.Clone()
	f.SetInherited(true)
	f.SetID(uuid.UUID{})
	if f.hasInheritance() {
		ter.setPropertyOnChildren(f)
	}
	// scrub instance startup information prior to storing
	p.clearInstances()
	ter.addProperty(p)
	ter.actionPropertyNew(p.MakeAction())
}

func (ter *Repository) setPropertyInherited(p Property) {
	f := p.Clone()
	f.SetID(f.GetInstanceID(ter.Type, ter.ID, ter.log))
	if f.Equal(uuid.Nil) {
		f.SetID(uuid.Must(uuid.NewV4()))
	}
	f.clearInstances()

	if !f.GetIsInherited() {
		ter.Fault.Error <- &Error{
			Action: `repository.setPropertyInherited on inherited=false`}
		return
	}
	if dupe, deleteOK, _ := ter.checkDuplicate(p); dupe && deleteOK {
		// we received an inherited SetProperty from above us in the
		// tree for a property that is duplicate, but we are not the
		// source of the duplicate -> corrupt tree
		ter.Fault.Error <- &Error{
			Action: `repository.setPropertyInherited corruption detected`}
		return
	} else if dupe && !deleteOK {
		// we received an inherited SetProperty from above us in the
		// tree for a property that is duplicate; we have a locally
		// set property -> stop inheritance, no error
		return
	}
	ter.addProperty(f)
	p.SetID(uuid.UUID{})
	ter.setPropertyOnChildren(p)
	ter.actionPropertyNew(f.MakeAction())
}

func (ter *Repository) setPropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	for child := range ter.Children {
		wg.Add(1)
		go func(stp Property, c string) {
			defer wg.Done()
			ter.Children[c].setPropertyInherited(stp)
		}(p, child)
	}
	wg.Wait()
}

func (ter *Repository) addProperty(p Property) {
	switch p.GetType() {
	case `custom`:
		ter.PropertyCustom[p.GetID()] = p
	case `system`:
		ter.PropertySystem[p.GetID()] = p
	case `service`:
		ter.PropertyService[p.GetID()] = p
	case `oncall`:
		ter.PropertyOncall[p.GetID()] = p
	default:
		ter.Fault.Error <- &Error{Action: `repository.addProperty unknown type`}
	}
}

//
// Propertier:> Update Property

func (ter *Repository) UpdateProperty(p Property) {
	if !ter.verifySourceInstance(
		p.GetSourceInstance(),
		p.GetType(),
	) {
		ter.Fault.Error <- &Error{Action: `update_property_on_non_source`}
		return
	}

	// keep a copy for ourselves, no shared pointers
	p.SetInheritedFrom(ter.ID)
	p.SetSourceType(ter.Type)
	p.SetInherited(true)
	f := p.Clone()
	f.SetInherited(false)
	if ter.switchProperty(f) {
		ter.updatePropertyOnChildren(p)
	}
}

func (ter *Repository) updatePropertyInherited(p Property) {
	// keep a copy for ourselves, no shared pointers
	f := p.Clone()
	if !f.GetIsInherited() {
		ter.Fault.Error <- &Error{
			Action: `repository.updatePropertyInherited on inherited=false`}
		return
	}
	if ter.switchProperty(f) {
		ter.updatePropertyOnChildren(p)
	}
}

func (ter *Repository) updatePropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	for child := range ter.Children {
		wg.Add(1)
		go func(stp Property, c string) {
			defer wg.Done()
			ter.Children[c].updatePropertyInherited(stp)
		}(p, child)
	}
	wg.Wait()
}

func (ter *Repository) switchProperty(p Property) bool {
	uid := ter.findIDForSource(
		p.GetSourceInstance(),
		p.GetType(),
	)
	if uid == `` {
		// we do not have the property for which we received an update
		if dupe, deleteOK, _ := ter.checkDuplicate(p); dupe && !deleteOK {
			// the update is duplicate to an property for which we
			// have the source instance, ie we just received an update
			// for which we have an overwrite. Ignore it and do not
			// inherit it further down
			return false
		}
		ter.Fault.Error <- &Error{
			Action: `repository.switchProperty property not found`}
		return false
	}
	updID, _ := uuid.FromString(uid)
	p.SetID(updID)
	curr := ter.getCurrentProperty(p)
	if curr == nil {
		return false
	}
	ter.addProperty(p)
	ter.actionPropertyUpdate(p.MakeAction())

	if !p.hasInheritance() && curr.hasInheritance() {
		// replacing inheritance with !inheritance:
		// call deletePropertyOnChildren(curr) to clean up
		srcUUID, _ := uuid.FromString(curr.GetSourceInstance())
		switch curr.GetType() {
		case `custom`:
			cstUUID, _ := uuid.FromString(curr.GetKey())
			ter.deletePropertyOnChildren(&PropertyCustom{
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
			ter.deletePropertyOnChildren(&PropertyService{
				SourceID:    srcUUID,
				View:        curr.GetView(),
				Inherited:   true,
				Service:     curr.GetKey(),
				Inheritance: true,
			})
		case `system`:
			ter.deletePropertyOnChildren(&PropertySystem{
				SourceID:    srcUUID,
				View:        curr.GetView(),
				Inherited:   true,
				Key:         curr.GetKey(),
				Value:       curr.GetValue(),
				Inheritance: true,
			})
		case `oncall`:
			oncUUID, _ := uuid.FromString(curr.GetKey())
			ter.deletePropertyOnChildren(&PropertyOncall{
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
		ter.setPropertyOnChildren(t)
	}
	return p.hasInheritance() && curr.hasInheritance()
}

func (ter *Repository) getCurrentProperty(p Property) Property {
	switch p.GetType() {
	case `custom`:
		return ter.PropertyCustom[p.GetID()].Clone()
	case `system`:
		return ter.PropertySystem[p.GetID()].Clone()
	case `service`:
		return ter.PropertyService[p.GetID()].Clone()
	case `oncall`:
		return ter.PropertyOncall[p.GetID()].Clone()
	}
	ter.Fault.Error <- &Error{
		Action: `repository.getCurrentProperty unknown type`}
	return nil
}

//
// Propertier:> Delete Property

func (ter *Repository) DeleteProperty(p Property) {
	if !ter.verifySourceInstance(
		p.GetSourceInstance(),
		p.GetType(),
	) {
		ter.Fault.Error <- &Error{Action: `repository.DeleteProperty on !source`}
		return
	}

	p.SetInherited(false)
	if ter.rmProperty(p) {
		p.SetInherited(true)
		ter.deletePropertyOnChildren(p)
	}
}

func (ter *Repository) deletePropertyInherited(p Property) {
	if ter.rmProperty(p) {
		ter.deletePropertyOnChildren(p)
	}
}

func (ter *Repository) deletePropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	for child := range ter.Children {
		wg.Add(1)
		go func(stp Property, c string) {
			defer wg.Done()
			ter.Children[c].deletePropertyInherited(stp)
		}(p, child)
	}
	wg.Wait()
}

func (ter *Repository) deletePropertyAllInherited() {
	for _, p := range ter.PropertyCustom {
		if !p.GetIsInherited() {
			continue
		}
		ter.deletePropertyInherited(p.Clone())
	}
	for _, p := range ter.PropertySystem {
		if !p.GetIsInherited() {
			continue
		}
		ter.deletePropertyInherited(p.Clone())
	}
	for _, p := range ter.PropertyService {
		if !p.GetIsInherited() {
			continue
		}
		ter.deletePropertyInherited(p.Clone())
	}
	for _, p := range ter.PropertyOncall {
		if !p.GetIsInherited() {
			continue
		}
		ter.deletePropertyInherited(p.Clone())
	}
}

func (ter *Repository) deletePropertyAllLocal() {
	for _, p := range ter.PropertyCustom {
		if p.GetIsInherited() {
			continue
		}
		ter.DeleteProperty(p.Clone())
	}
	for _, p := range ter.PropertySystem {
		if p.GetIsInherited() {
			continue
		}
		ter.DeleteProperty(p.Clone())
	}
	for _, p := range ter.PropertyService {
		if p.GetIsInherited() {
			continue
		}
		ter.DeleteProperty(p.Clone())
	}
	for _, p := range ter.PropertyOncall {
		if p.GetIsInherited() {
			continue
		}
		ter.DeleteProperty(p.Clone())
	}
}

func (ter *Repository) rmProperty(p Property) bool {
	delID := ter.findIDForSource(
		p.GetSourceInstance(),
		p.GetType(),
	)
	if delID == `` {
		// we do not have the property for which we received a delete
		if dupe, deleteOK, _ := ter.checkDuplicate(p); dupe && !deleteOK {
			// the delete is duplicate to a property for which we
			// have the source instance, ie we just received a delete
			// for which we have an overwrite. Ignore it and do not
			// inherit it further down
			return false
		}

		ter.Fault.Error <- &Error{
			Action: `repository.rmProperty property not found`}
		return false
	}

	hasInheritance := false
	switch p.GetType() {
	case `custom`:
		ter.actionPropertyDelete(
			ter.PropertyCustom[delID].MakeAction(),
		)
		hasInheritance = ter.PropertyCustom[delID].hasInheritance()
		delete(ter.PropertyCustom, delID)
	case `service`:
		ter.actionPropertyDelete(
			ter.PropertyService[delID].MakeAction(),
		)
		hasInheritance = ter.PropertyService[delID].hasInheritance()
		delete(ter.PropertyService, delID)
	case `system`:
		ter.actionPropertyDelete(
			ter.PropertySystem[delID].MakeAction(),
		)
		hasInheritance = ter.PropertySystem[delID].hasInheritance()
		delete(ter.PropertySystem, delID)
	case `oncall`:
		ter.actionPropertyDelete(
			ter.PropertyOncall[delID].MakeAction(),
		)
		hasInheritance = ter.PropertyOncall[delID].hasInheritance()
		delete(ter.PropertyOncall, delID)
	default:
		ter.Fault.Error <- &Error{Action: `repository.rmProperty unknown type`}
		return false
	}
	return hasInheritance
}

//
// Propertier:> Utility

//
func (ter *Repository) verifySourceInstance(id, prop string) bool {
	switch prop {
	case `custom`:
		if _, ok := ter.PropertyCustom[id]; !ok {
			goto bailout
		}
		return ter.PropertyCustom[id].GetSourceInstance() == id
	case `service`:
		if _, ok := ter.PropertyService[id]; !ok {
			goto bailout
		}
		return ter.PropertyService[id].GetSourceInstance() == id
	case `system`:
		if _, ok := ter.PropertySystem[id]; !ok {
			goto bailout
		}
		return ter.PropertySystem[id].GetSourceInstance() == id
	case `oncall`:
		if _, ok := ter.PropertyOncall[id]; !ok {
			goto bailout
		}
		return ter.PropertyOncall[id].GetSourceInstance() == id
	}

bailout:
	ter.Fault.Error <- &Error{
		Action: `repository.verifySourceInstance not found`}
	return false
}

func (ter *Repository) findIDForSource(source, prop string) string {
	switch prop {
	case `custom`:
		for id := range ter.PropertyCustom {
			if ter.PropertyCustom[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	case `system`:
		for id := range ter.PropertySystem {
			if ter.PropertySystem[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	case `service`:
		for id := range ter.PropertyService {
			if ter.PropertyService[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	case `oncall`:
		for id := range ter.PropertyOncall {
			if ter.PropertyOncall[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	}
	return ``
}

//
func (ter *Repository) resyncProperty(srcID, pType, childID string) {
	pID := ter.findIDForSource(srcID, pType)
	if pID == `` {
		return
	}

	var f Property
	switch pType {
	case `custom`:
		f = ter.PropertyCustom[pID].(*PropertyCustom).Clone()
	case `oncall`:
		f = ter.PropertyOncall[pID].(*PropertyOncall).Clone()
	case `service`:
		f = ter.PropertyService[pID].(*PropertyService).Clone()
	case `system`:
		f = ter.PropertySystem[pID].(*PropertySystem).Clone()
	}
	if !f.hasInheritance() {
		return
	}
	f.SetInherited(true)
	f.SetID(uuid.UUID{})
	f.clearInstances()
	ter.Children[childID].setPropertyInherited(f)
}

// when a child attaches, it calls self.Parent.syncProperty(self.ID)
// to get get all properties of that part of the tree
func (ter *Repository) syncProperty(childID string) {
customloop:
	for prop := range ter.PropertyCustom {
		if !ter.PropertyCustom[prop].hasInheritance() {
			continue customloop
		}
		f := ter.PropertyCustom[prop].(*PropertyCustom).Clone()
		f.SetInherited(true)
		f.SetID(uuid.UUID{})
		f.clearInstances()
		ter.Children[childID].setPropertyInherited(f)
	}
oncallloop:
	for prop := range ter.PropertyOncall {
		if !ter.PropertyOncall[prop].hasInheritance() {
			continue oncallloop
		}
		f := ter.PropertyOncall[prop].(*PropertyOncall).Clone()
		f.SetInherited(true)
		f.SetID(uuid.UUID{})
		f.clearInstances()
		ter.Children[childID].setPropertyInherited(f)
	}
serviceloop:
	for prop := range ter.PropertyService {
		if !ter.PropertyService[prop].hasInheritance() {
			continue serviceloop
		}
		f := ter.PropertyService[prop].(*PropertyService).Clone()
		f.SetInherited(true)
		f.SetID(uuid.UUID{})
		f.clearInstances()
		ter.Children[childID].setPropertyInherited(f)
	}
systemloop:
	for prop := range ter.PropertySystem {
		if !ter.PropertySystem[prop].hasInheritance() {
			continue systemloop
		}
		f := ter.PropertySystem[prop].(*PropertySystem).Clone()
		f.SetInherited(true)
		f.SetID(uuid.UUID{})
		f.clearInstances()
		ter.Children[childID].setPropertyInherited(f)
	}
}

// function to be used by a child to check if the parent has a
// specific Property
func (ter *Repository) checkProperty(propType string, propID string) bool {
	switch propType {
	case "custom":
		if _, ok := ter.PropertyCustom[propID]; ok {
			return true
		}
	case "service":
		if _, ok := ter.PropertyService[propID]; ok {
			return true
		}
	case "system":
		if _, ok := ter.PropertySystem[propID]; ok {
			return true
		}
	case "oncall":
		if _, ok := ter.PropertyOncall[propID]; ok {
			return true
		}
	}
	return false
}

// Checks if this property is already defined on this node, and
// whether it was inherited, ie. can be deleted so it can be
// overwritten
func (ter *Repository) checkDuplicate(p Property) (bool, bool, Property) {
	var dupe, deleteOK bool
	var prop Property

propswitch:
	switch p.GetType() {
	case "custom":
		for _, pVal := range ter.PropertyCustom {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "service":
		for _, pVal := range ter.PropertyService {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "oncall":
		for _, pVal := range ter.PropertyOncall {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case msg.PropertySystem:
		for _, pVal := range ter.PropertySystem {
			switch p.GetKey() {
			case msg.SystemPropertyTag:
				// tags are only dupes if the value is the same as well
				fallthrough
			case msg.SystemPropertyDisableCheckConfiguration:
				// disable_check_configuration checks values as well
				if p.GetValue() == pVal.GetValue() {
					// same value, can be a dupe
					dupe, deleteOK, prop = isDupe(pVal, p)
					if dupe {
						break propswitch
					}
				}
			default:
				dupe, deleteOK, prop = isDupe(pVal, p)
				if dupe {
					break propswitch
				}
			}
		}
	default:
		// trigger error path
		ter.Fault.Error <- &Error{Action: `repository.checkDuplicate unknown type`}
		dupe = true
		deleteOK = false
	}
	return dupe, deleteOK, prop
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
