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

func (tec *Cluster) SetCheck(c Check) {
	c.ID = c.GetItemID(tec.Type, tec.ID)
	if uuid.Equal(c.ID, uuid.Nil) {
		c.ID = uuid.Must(uuid.NewV4())
	}
	// this check is the source check
	c.InheritedFrom = tec.ID
	c.Inherited = false
	c.SourceID, _ = uuid.FromString(c.ID.String())
	c.SourceType = tec.Type
	// send a scrubbed copy downward
	f := c.Clone()
	f.Inherited = true
	f.ID = uuid.Nil
	tec.setCheckOnChildren(f)
	// scrub checkitem startup information prior to storing
	c.Items = nil
	tec.addCheck(c)
}

func (tec *Cluster) setCheckInherited(c Check) {
	// we keep a local copy, that way we know it is ours....
	f := c.Clone()
	f.ID = f.GetItemID(tec.Type, tec.ID)
	if uuid.Equal(f.ID, uuid.Nil) {
		f.ID = uuid.Must(uuid.NewV4())
	}
	// send original check downwards
	c.ID = uuid.Nil
	tec.setCheckOnChildren(c)
	f.Items = nil
	tec.addCheck(f)
}

func (tec *Cluster) setCheckOnChildren(c Check) {
	switch deterministicInheritanceOrder {
	case true:
		// nodes
		for i := 0; i < tec.ordNumChildNod; i++ {
			if child, ok := tec.ordChildrenNod[i]; ok {
				tec.Children[child].(Checker).setCheckInherited(c)
			}
		}
	default:
		var wg sync.WaitGroup
		for child := range tec.Children {
			wg.Add(1)
			go func(stc Check, ch string) {
				defer wg.Done()
				tec.Children[ch].(Checker).setCheckInherited(stc)
			}(c, child)
		}
		wg.Wait()
	}
}

func (tec *Cluster) addCheck(c Check) {
	tec.hasUpdate = true
	tec.Checks[c.ID.String()] = c
	tec.actionCheckNew(tec.setupCheckAction(c))
}

//
// Checker:> Remove Check

func (tec *Cluster) DeleteCheck(c Check) {
	tec.deleteCheckOnChildren(c)
	tec.rmCheck(c)
}

func (tec *Cluster) deleteCheckInherited(c Check) {
	tec.deleteCheckOnChildren(c)
	tec.rmCheck(c)
}

func (tec *Cluster) deleteCheckAllInherited() {
	for _, check := range tec.Checks {
		if check.GetIsInherited() {
			tec.deleteCheckInherited(check.Clone())
		}

	}

}

func (tec *Cluster) deleteCheckOnChildren(c Check) {
	switch deterministicInheritanceOrder {
	case true:
		for i := 0; i < tec.ordNumChildNod; i++ {
			if child, ok := tec.ordChildrenNod[i]; ok {
				tec.Children[child].(Checker).deleteCheckInherited(c)
			}
		}
	default:
		var wg sync.WaitGroup
		for child := range tec.Children {
			wg.Add(1)
			go func(stc Check, ch string) {
				defer wg.Done()
				tec.Children[ch].(Checker).deleteCheckInherited(stc)
			}(c, child)
		}
		wg.Wait()
	}
}

func (tec *Cluster) deleteCheckLocalAll() {
	localChecks := make(chan *Check, len(tec.Checks)+1)

	for _, check := range tec.Checks {
		if check.GetIsInherited() {
			// not a locally configured check
			continue
		}
		localChecks <- &check
	}
	close(localChecks)

	for check := range localChecks {
		tec.DeleteCheck(check.Clone())
	}
}

func (tec *Cluster) rmCheck(c Check) {
	for id := range tec.Checks {
		if uuid.Equal(tec.Checks[id].SourceID, c.SourceID) {
			tec.hasUpdate = true
			tec.actionCheckRemoved(tec.setupCheckAction(tec.Checks[id]))
			delete(tec.Checks, id)
			return
		}
	}
}

//
// Checker:> Meta

func (tec *Cluster) syncCheck(childID string) {
	for check := range tec.Checks {
		if !tec.Checks[check].Inheritance {
			continue
		}
		// build a pristine version for inheritance
		f := tec.Checks[check]
		c := f.Clone()
		c.Inherited = true
		c.ID = uuid.Nil
		c.Items = nil
		tec.Children[childID].(Checker).setCheckInherited(c)
	}
}

func (tec *Cluster) checkCheck(checkID string) bool {
	if _, ok := tec.Checks[checkID]; ok {
		return true
	}
	return false
}

func (tec *Cluster) LoadInstance(i CheckInstance) {
	ckID := i.CheckID.String()
	ckInstID := i.InstanceID.String()
	if tec.loadedInstances[ckID] == nil {
		tec.loadedInstances[ckID] = map[string]CheckInstance{}
	}
	tec.loadedInstances[ckID][ckInstID] = i
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
