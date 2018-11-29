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

// Implementation of the `Checker` interface

//
// Checker:> Add Check

func (teg *Group) SetCheck(c Check) {
	c.ID = c.GetItemID(teg.Type, teg.ID)
	if uuid.Equal(c.ID, uuid.Nil) {
		c.ID = uuid.Must(uuid.NewV4())
	}
	// this check is the source check
	c.InheritedFrom = teg.ID
	c.Inherited = false
	c.SourceID, _ = uuid.FromString(c.ID.String())
	c.SourceType = teg.Type
	// send a scrubbed copy downward
	f := c.Clone()
	f.Inherited = true
	f.ID = uuid.Nil
	teg.setCheckOnChildren(f)
	// scrub checkitem startup information prior to storing
	c.Items = nil
	teg.addCheck(c)
}

func (teg *Group) setCheckInherited(c Check) {
	// we keep a local copy, that way we know it is ours....
	f := c.Clone()
	f.ID = f.GetItemID(teg.Type, teg.ID)
	if uuid.Equal(f.ID, uuid.Nil) {
		f.ID = uuid.Must(uuid.NewV4())
	}
	// send original check downwards
	c.ID = uuid.Nil
	teg.setCheckOnChildren(c)
	f.Items = nil
	teg.addCheck(f)
}

func (teg *Group) setCheckOnChildren(c Check) {
	switch deterministicInheritanceOrder {
	case true:
		// groups
		for i := 0; i < teg.ordNumChildGrp; i++ {
			if child, ok := teg.ordChildrenGrp[i]; ok {
				teg.Children[child].(Checker).setCheckInherited(c)
			}
		}
		// clusters
		for i := 0; i < teg.ordNumChildClr; i++ {
			if child, ok := teg.ordChildrenClr[i]; ok {
				teg.Children[child].(Checker).setCheckInherited(c)
			}
		}
		// nodes
		for i := 0; i < teg.ordNumChildNod; i++ {
			if child, ok := teg.ordChildrenNod[i]; ok {
				teg.Children[child].(Checker).setCheckInherited(c)
			}
		}
	default:
		var wg sync.WaitGroup
		for child, _ := range teg.Children {
			wg.Add(1)
			go func(stc Check, ch string) {
				defer wg.Done()
				teg.Children[ch].(Checker).setCheckInherited(stc)
			}(c, child)
		}
		wg.Wait()
	}
}

func (teg *Group) addCheck(c Check) {
	teg.hasUpdate = true
	teg.Checks[c.ID.String()] = c
	teg.actionCheckNew(c.MakeAction())
}

//
// Checker:> Remove Check

func (teg *Group) DeleteCheck(c Check) {
	teg.deleteCheckOnChildren(c)
	teg.rmCheck(c)
}

func (teg *Group) deleteCheckInherited(c Check) {
	teg.deleteCheckOnChildren(c)
	teg.rmCheck(c)
}

func (teg *Group) deleteCheckOnChildren(c Check) {
	switch deterministicInheritanceOrder {
	case true:
		// groups
		for i := 0; i < teg.ordNumChildGrp; i++ {
			if child, ok := teg.ordChildrenGrp[i]; ok {
				teg.Children[child].(Checker).deleteCheckInherited(c)
			}
		}
		// clusters
		for i := 0; i < teg.ordNumChildClr; i++ {
			if child, ok := teg.ordChildrenClr[i]; ok {
				teg.Children[child].(Checker).deleteCheckInherited(c)
			}
		}
		// nodes
		for i := 0; i < teg.ordNumChildNod; i++ {
			if child, ok := teg.ordChildrenNod[i]; ok {
				teg.Children[child].(Checker).deleteCheckInherited(c)
			}
		}
	default:
		var wg sync.WaitGroup
		for child, _ := range teg.Children {
			wg.Add(1)
			go func(stc Check, ch string) {
				defer wg.Done()
				teg.Children[ch].(Checker).deleteCheckInherited(stc)
			}(c, child)
		}
		wg.Wait()
	}
}

func (teg *Group) rmCheck(c Check) {
	for id := range teg.Checks {
		if uuid.Equal(teg.Checks[id].SourceID, c.SourceID) {
			teg.hasUpdate = true
			teg.actionCheckRemoved(teg.setupCheckAction(teg.Checks[id]))
			delete(teg.Checks, id)
			return
		}
	}
}

//
// Checker:> Meta

func (teg *Group) syncCheck(childID string) {
	for check := range teg.Checks {
		if !teg.Checks[check].Inheritance {
			continue
		}
		// build a pristine version for inheritance
		f := teg.Checks[check]
		c := f.Clone()
		c.Inherited = true
		c.ID = uuid.Nil
		c.Items = nil
		teg.Children[childID].(Checker).setCheckInherited(c)
	}
}

func (teg *Group) checkCheck(checkID string) bool {
	if _, ok := teg.Checks[checkID]; ok {
		return true
	}
	return false
}

//
func (teg *Group) LoadInstance(i CheckInstance) {
	ckID := i.CheckID.String()
	ckInstID := i.InstanceID.String()
	if teg.loadedInstances[ckID] == nil {
		teg.loadedInstances[ckID] = map[string]CheckInstance{}
	}
	teg.loadedInstances[ckID][ckInstID] = i
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
