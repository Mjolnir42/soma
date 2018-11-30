/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import (
	"github.com/satori/go.uuid"
	"sync"
)

// Implementation of the `Checker` interface

//
// Checker:> Add Check

func (teb *Bucket) SetCheck(c Check) {
	c.ID = c.GetItemID(teb.Type, teb.ID)
	if uuid.Equal(c.ID, uuid.Nil) {
		c.ID = uuid.Must(uuid.NewV4())
	}
	// this check is the source check
	c.InheritedFrom = teb.ID
	c.Inherited = false
	c.SourceID, _ = uuid.FromString(c.ID.String())
	c.SourceType = teb.Type
	// send a scrubbed copy downward
	f := c.Clone()
	f.Inherited = true
	f.ID = uuid.Nil
	teb.setCheckOnChildren(f)
	// scrub checkitem startup information prior to storing
	c.Items = nil
	teb.addCheck(c)
}

func (teb *Bucket) setCheckInherited(c Check) {
	// we keep a local copy, that way we know it is ours....
	f := c.Clone()
	f.ID = f.GetItemID(teb.Type, teb.ID)
	if uuid.Equal(f.ID, uuid.Nil) {
		f.ID = uuid.Must(uuid.NewV4())
	}
	// send original check downwards
	c.ID = uuid.Nil
	teb.setCheckOnChildren(c)
	f.Items = nil
	teb.addCheck(f)
}

func (teb *Bucket) setCheckOnChildren(c Check) {
	switch deterministicInheritanceOrder {
	case true:
		// groups
		for i := 0; i < teb.ordNumChildGrp; i++ {
			if child, ok := teb.ordChildrenGrp[i]; ok {
				teb.Children[child].(Checker).setCheckInherited(c)
			}
		}
		// clusters
		for i := 0; i < teb.ordNumChildClr; i++ {
			if child, ok := teb.ordChildrenClr[i]; ok {
				teb.Children[child].(Checker).setCheckInherited(c)
			}
		}
		// nodes
		for i := 0; i < teb.ordNumChildNod; i++ {
			if child, ok := teb.ordChildrenNod[i]; ok {
				teb.Children[child].(Checker).setCheckInherited(c)
			}
		}
	default:
		var wg sync.WaitGroup
		for child, _ := range teb.Children {
			wg.Add(1)
			go func(stc Check, ch string) {
				defer wg.Done()
				teb.Children[ch].(Checker).setCheckInherited(stc)
			}(c, child)
		}
		wg.Wait()
	}
}

func (teb *Bucket) addCheck(c Check) {
	teb.Checks[c.ID.String()] = c
	teb.actionCheckNew(teb.setupCheckAction(c))
}

//
// Checker:> Remove Check

func (teb *Bucket) DeleteCheck(c Check) {
	teb.deleteCheckOnChildren(c)
	teb.rmCheck(c)
}

func (teb *Bucket) deleteCheckInherited(c Check) {
	teb.deleteCheckOnChildren(c)
	teb.rmCheck(c)
}

func (teb *Bucket) deleteCheckOnChildren(c Check) {
	switch deterministicInheritanceOrder {
	case true:
		// groups
		for i := 0; i < teb.ordNumChildGrp; i++ {
			if child, ok := teb.ordChildrenGrp[i]; ok {
				teb.Children[child].(Checker).deleteCheckInherited(c)
			}
		}
		// clusters
		for i := 0; i < teb.ordNumChildClr; i++ {
			if child, ok := teb.ordChildrenClr[i]; ok {
				teb.Children[child].(Checker).deleteCheckInherited(c)
			}
		}
		// nodes
		for i := 0; i < teb.ordNumChildNod; i++ {
			if child, ok := teb.ordChildrenNod[i]; ok {
				teb.Children[child].(Checker).deleteCheckInherited(c)
			}
		}
	default:
		var wg sync.WaitGroup
		for child, _ := range teb.Children {
			wg.Add(1)
			go func(stc Check, ch string) {
				defer wg.Done()
				teb.Children[ch].(Checker).deleteCheckInherited(stc)
			}(c, child)
		}
		wg.Wait()

	}
}

func (teb *Bucket) deleteCheckLocalAll() {
	localChecks := make(chan *Check, len(teb.Checks)+1)

	for _, check := range teb.Checks {
		if check.GetIsInherited() {
			// not a locally configured check
			continue
		}
		localChecks <- &check
	}
	close(localChecks)

	for check := range localChecks {
		teb.DeleteCheck(check.Clone())
	}
}

func (teb *Bucket) rmCheck(c Check) {
	for id := range teb.Checks {
		if uuid.Equal(teb.Checks[id].SourceID, c.SourceID) {
			teb.actionCheckRemoved(teb.setupCheckAction(teb.Checks[id]))
			delete(teb.Checks, id)
			return
		}
	}
}

//
// Checker:> Meta

func (teb *Bucket) syncCheck(childID string) {
	for check := range teb.Checks {
		if !teb.Checks[check].Inheritance {
			continue
		}
		// build a pristine version for inheritance
		f := teb.Checks[check]
		c := f.Clone()
		c.Inherited = true
		c.ID = uuid.Nil
		c.Items = nil
		teb.Children[childID].(Checker).setCheckInherited(c)
	}
}

func (teb *Bucket) checkCheck(checkID string) bool {
	if _, ok := teb.Checks[checkID]; ok {
		return true
	}
	return false
}

// XXX
func (teb *Bucket) LoadInstance(i CheckInstance) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
