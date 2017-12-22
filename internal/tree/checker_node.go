/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import "github.com/satori/go.uuid"

// Implementation of the `Checker` interface

//
// Checker:> Add Check

func (ten *Node) SetCheck(c Check) {
	c.ID = c.GetItemID(ten.Type, ten.ID)
	if uuid.Equal(c.ID, uuid.Nil) {
		c.ID = uuid.NewV4()
	}
	// this check is the source check
	c.InheritedFrom = ten.ID
	c.Inherited = false
	c.SourceID, _ = uuid.FromString(c.ID.String())
	c.SourceType = ten.Type
	// scrub checkitem startup information prior to storing
	c.Items = nil
	ten.addCheck(c)
}

func (ten *Node) setCheckInherited(c Check) {
	// we keep a local copy, that way we know it is ours....
	f := c.Clone()
	f.ID = f.GetItemID(ten.Type, ten.ID)
	if uuid.Equal(f.ID, uuid.Nil) {
		f.ID = uuid.NewV4()
	}
	f.Items = nil
	ten.addCheck(f)
}

func (ten *Node) setCheckOnChildren(c Check) {
}

func (ten *Node) addCheck(c Check) {
	ten.hasUpdate = true
	ten.Checks[c.ID.String()] = c
	ten.actionCheckNew(ten.setupCheckAction(c))
}

//
// Checker:> Remove Check

func (ten *Node) DeleteCheck(c Check) {
	ten.rmCheck(c)
}

func (ten *Node) deleteCheckInherited(c Check) {
	ten.rmCheck(c)
}

func (ten *Node) deleteCheckOnChildren(c Check) {
}

func (ten *Node) rmCheck(c Check) {
	for id := range ten.Checks {
		if uuid.Equal(ten.Checks[id].SourceID, c.SourceID) {
			ten.hasUpdate = true
			ten.actionCheckRemoved(ten.setupCheckAction(ten.Checks[id]))
			delete(ten.Checks, id)
			return
		}
	}
}

// noop, satisfy interface
func (ten *Node) syncCheck(childID string) {
}

func (ten *Node) checkCheck(checkID string) bool {
	if _, ok := ten.Checks[checkID]; ok {
		return true
	}
	return false
}

//
func (ten *Node) LoadInstance(i CheckInstance) {
	ckID := i.CheckID.String()
	ckInstID := i.InstanceID.String()
	if ten.loadedInstances[ckID] == nil {
		ten.loadedInstances[ckID] = map[string]CheckInstance{}
	}
	ten.loadedInstances[ckID][ckInstID] = i
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
