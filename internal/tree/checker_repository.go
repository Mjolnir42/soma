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
func (ter *Repository) SetCheck(c Check) {
	c.ID = c.GetItemID(ter.Type, ter.ID)
	if uuid.Equal(c.ID, uuid.Nil) {
		c.ID = uuid.Must(uuid.NewV4())
	}
	// this check is the source check
	c.InheritedFrom = ter.ID
	c.Inherited = false
	c.SourceID, _ = uuid.FromString(c.ID.String())
	c.SourceType = ter.Type
	// send a scrubbed copy downward
	f := c.Clone()
	f.Inherited = true
	f.ID = uuid.Nil
	ter.setCheckOnChildren(f)
	// scrub checkitem startup information prior to storing
	c.Items = nil
	ter.addCheck(c)
}

func (ter *Repository) setCheckInherited(c Check) {
	// we keep a local copy, that way we know it is ours....
	f := c.Clone()
	f.ID = f.GetItemID(ter.Type, ter.ID)
	if uuid.Equal(f.ID, uuid.Nil) {
		f.ID = uuid.Must(uuid.NewV4())
	}
	// send original check downwards
	c.ID = uuid.Nil
	ter.setCheckOnChildren(c)
	f.Items = nil
	ter.addCheck(f)
}

func (ter *Repository) setCheckOnChildren(c Check) {
	switch deterministicInheritanceOrder {
	case true:
		// buckets
		for i := 0; i < ter.ordNumChildBck; i++ {
			if child, ok := ter.ordChildrenBck[i]; ok {
				ter.Children[child].(Checker).setCheckInherited(c)
			}
		}
	default:
		var wg sync.WaitGroup
		for child, _ := range ter.Children {
			wg.Add(1)
			go func(stc Check, ch string) {
				defer wg.Done()
				ter.Children[ch].(Checker).setCheckInherited(stc)
			}(c, child)
		}
		wg.Wait()
	}
}

func (ter *Repository) addCheck(c Check) {
	ter.Checks[c.ID.String()] = c
	ter.actionCheckNew(c.MakeAction())
}

//
// Checker:> Remove Check

func (ter *Repository) DeleteCheck(c Check) {
	ter.deleteCheckOnChildren(c)
	ter.rmCheck(c)
}

func (ter *Repository) deleteCheckInherited(c Check) {
	ter.deleteCheckOnChildren(c)
	ter.rmCheck(c)
}

func (ter *Repository) deleteCheckOnChildren(c Check) {
	switch deterministicInheritanceOrder {
	case true:
		// buckets
		for i := 0; i < ter.ordNumChildBck; i++ {
			if child, ok := ter.ordChildrenBck[i]; ok {
				ter.Children[child].(Checker).deleteCheckInherited(c)
			}
		}
	default:
		var wg sync.WaitGroup
		for child, _ := range ter.Children {
			wg.Add(1)
			go func(stc Check, ch string) {
				defer wg.Done()
				ter.Children[ch].(Checker).deleteCheckInherited(stc)
			}(c, child)
		}
		wg.Wait()
	}
}

func (ter *Repository) deleteCheckLocalAll() {
	localChecks := make(chan *Check, len(ter.Checks)+1)

	for _, check := range ter.Checks {
		if check.GetIsInherited() {
			// not a locally configured check
			continue
		}
		localChecks <- &check
	}
	close(localChecks)

	for check := range localChecks {
		ter.DeleteCheck(check.Clone())
	}
}

func (ter *Repository) rmCheck(c Check) {
	for id := range ter.Checks {
		if uuid.Equal(ter.Checks[id].SourceID, c.SourceID) {
			ter.actionCheckRemoved(ter.setupCheckAction(ter.Checks[id]))
			delete(ter.Checks, id)
			return
		}
	}
}

//
// Checker:> Meta

func (ter *Repository) syncCheck(childID string) {
	for check := range ter.Checks {
		if !ter.Checks[check].Inheritance {
			continue
		}
		// build a pristine version for inheritance
		f := ter.Checks[check]
		c := f.Clone()
		c.Inherited = true
		c.ID = uuid.Nil
		c.Items = nil
		ter.Children[childID].(Checker).setCheckInherited(c)
	}
}

func (ter *Repository) checkCheck(checkID string) bool {
	if _, ok := ter.Checks[checkID]; ok {
		return true
	}
	return false
}

// XXX
func (ter *Repository) LoadInstance(i CheckInstance) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
