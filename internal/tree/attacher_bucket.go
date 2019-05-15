/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import (
	"fmt"
	"reflect"
	"sync"

	uuid "github.com/satori/go.uuid"
)

//
// Interface: Attacher
func (teb *Bucket) Attach(a AttachRequest) {
	if teb.Parent != nil {
		panic(`Bucket.Attach: already attached`)
	}
	switch {
	case a.ParentType == "repository":
		teb.attachToRepository(a)
	default:
		panic(`Bucket.Attach`)
	}

	teb.Parent.(Propertier).syncProperty(teb.ID.String())
	teb.Parent.(Checker).syncCheck(teb.ID.String())
}

func (teb *Bucket) Destroy() {
	if teb.Parent == nil {
		panic(`Bucket.Destroy called without Parent to unlink from`)
	}
	teb.deletePropertyAllLocal()
	teb.deletePropertyAllInherited()
	teb.deleteCheckLocalAll()
	teb.deleteCheckAllInherited()

	wg := new(sync.WaitGroup)
	for child := range teb.Children {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			teb.Children[c].Destroy()
		}(child)
	}
	wg.Wait()

	teb.Parent.Unlink(UnlinkRequest{
		ParentType: teb.Parent.(Builder).GetType(),
		ParentID:   teb.Parent.(Builder).GetID(),
		ParentName: teb.Parent.(Builder).GetName(),
		ChildType:  teb.GetType(),
		ChildName:  teb.GetName(),
		ChildID:    teb.GetID(),
	},
	)

	teb.setFault(nil)
	teb.actionDelete()
	teb.setAction(nil)
}

func (teb *Bucket) Detach() {
	teb.Destroy()
}

func (teb *Bucket) SetName(s string) {
	teb.Name = s
	teb.actionRename()
}

func (teb *Bucket) clearParent() {
	teb.Parent = nil
	teb.State = "floating"
}

func (teb *Bucket) setFault(f *Fault) {
	teb.Fault = f
}

func (teb *Bucket) updateParentRecursive(p Receiver) {
	teb.setParent(p)
	var wg sync.WaitGroup
	for child := range teb.Children {
		wg.Add(1)
		c := child
		go func(str Receiver) {
			defer wg.Done()
			teb.Children[c].updateParentRecursive(str)
		}(teb)
	}
	wg.Wait()
}

func (teb *Bucket) updateFaultRecursive(f *Fault) {
	teb.setFault(f)
	var wg sync.WaitGroup
	for child := range teb.Children {
		wg.Add(1)
		c := child
		go func(ptr *Fault) {
			defer wg.Done()
			teb.Children[c].updateFaultRecursive(ptr)
		}(f)
	}
	wg.Wait()
}

func (teb *Bucket) setParent(p Receiver) {
	switch p.(type) {
	case BucketReceiver:
		teb.setBucketParent(p.(BucketReceiver))
		teb.State = "attached"
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`Bucket.setParent`)
	}
}

func (teb *Bucket) setBucketParent(p BucketReceiver) {
	teb.Parent = p
}

func (teb *Bucket) inheritTeamID(newTeamID string) {
	wg := sync.WaitGroup{}
	switch deterministicInheritanceOrder {
	case true:
		// groups
		for i := 0; i < teb.ordNumChildGrp; i++ {
			if child, ok := teb.ordChildrenGrp[i]; ok {
				teb.Children[child].inheritTeamID(newTeamID)
			}
		}
		// clusters
		for i := 0; i < teb.ordNumChildClr; i++ {
			if child, ok := teb.ordChildrenClr[i]; ok {
				teb.Children[child].inheritTeamID(newTeamID)
			}
		}
		// nodes
		for i := 0; i < teb.ordNumChildNod; i++ {
			if child, ok := teb.ordChildrenNod[i]; ok {
				teb.Children[child].inheritTeamID(newTeamID)
			}
		}
	default:
		for child := range teb.Children {
			wg.Add(1)
			go func(name, teamID string) {
				defer wg.Done()
				teb.Children[name].inheritTeamID(teamID)
			}(child, newTeamID)
		}
	}
	teb.Team, _ = uuid.FromString(newTeamID)
	teb.actionRepossess()
	wg.Wait()
}

//
// Interface: RepositoryAttacher
func (teb *Bucket) attachToRepository(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentID:   a.ParentID,
		ParentName: a.ParentName,
		ChildType:  teb.Type,
		Bucket:     teb,
	})

	if teb.Parent == nil {
		a.Root.(*Tree).AttachError(Error{Action: `attach_bucket`})
		return
	}
	teb.actionCreate()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
