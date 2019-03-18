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

	uuid "github.com/satori/go.uuid"
)

//
// Interface: Attacher
func (teg *Group) Attach(a AttachRequest) {
	if teg.Parent != nil {
		panic(`Group.Attach: already attached`)
	}
	switch {
	case a.ParentType == "bucket":
		teg.attachToBucket(a)
	case a.ParentType == "group":
		teg.attachToGroup(a)
	default:
		panic(`Group.Attach`)
	}

	teg.Parent.(Propertier).syncProperty(teg.ID.String())
	teg.Parent.(Checker).syncCheck(teg.ID.String())
}

func (teg *Group) ReAttach(a AttachRequest) {
	if teg.Parent == nil {
		panic(`Group.ReAttach: not attached`)
	}
	teg.deletePropertyAllInherited()
	// TODO delete all inherited checks + check instances

	teg.Parent.Unlink(UnlinkRequest{
		ParentType: teg.Parent.(Builder).GetType(),
		ParentName: teg.Parent.(Builder).GetName(),
		ParentID:   teg.Parent.(Builder).GetID(),
		ChildType:  teg.GetType(),
		ChildName:  teg.GetName(),
		ChildID:    teg.GetID(),
	},
	)

	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentID:   a.ParentID,
		ParentName: a.ParentName,
		ChildType:  teg.GetType(),
		Group:      teg,
	},
	)

	if teg.Parent == nil {
		panic(`Group.ReAttach: not reattached`)
	}
	teg.actionUpdate()
	teg.Parent.(Propertier).syncProperty(teg.ID.String())
	teg.Parent.(Checker).syncCheck(teg.ID.String())
}

func (teg *Group) Destroy() {
	if teg.Parent == nil {
		panic(`Group.Destroy called without Parent to unlink from`)
	}

	// call before unlink since it requires teg.Parent.*
	teg.actionDelete()
	teg.deletePropertyAllLocal()
	teg.deletePropertyAllInherited()
	teg.deleteCheckLocalAll()
	teg.updateCheckInstances()

	wg := new(sync.WaitGroup)
	for child := range teg.Children {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			teg.Children[c].Destroy()
		}(child)
	}
	wg.Wait()

	teg.Parent.Unlink(UnlinkRequest{
		ParentType: teg.Parent.(Builder).GetType(),
		ParentID:   teg.Parent.(Builder).GetID(),
		ParentName: teg.Parent.(Builder).GetName(),
		ChildType:  teg.GetType(),
		ChildName:  teg.GetName(),
		ChildID:    teg.GetID(),
	},
	)

	teg.setFault(nil)
	teg.setAction(nil)
}

func (teg *Group) Detach() {
	if teg.Parent == nil {
		panic(`Group.Destroy called without Parent to detach from`)
	}
	bucket := teg.Parent.(Bucketeer).GetBucket()

	teg.deletePropertyAllInherited()
	// TODO delete all inherited checks + check instances

	teg.Parent.Unlink(UnlinkRequest{
		ParentType: teg.Parent.(Builder).GetType(),
		ParentID:   teg.Parent.(Builder).GetID(),
		ParentName: teg.Parent.(Builder).GetName(),
		ChildType:  teg.GetType(),
		ChildName:  teg.GetName(),
		ChildID:    teg.GetID(),
	},
	)

	bucket.Receive(ReceiveRequest{
		ParentType: bucket.(Builder).GetType(),
		ParentID:   bucket.(Builder).GetID(),
		ParentName: bucket.(Builder).GetName(),
		ChildType:  teg.Type,
		Group:      teg,
	},
	)

	teg.actionUpdate()
	teg.Parent.(Propertier).syncProperty(teg.ID.String())
}

func (teg *Group) SetName(s string) {
	teg.Name = s
	teg.actionRename()
}

func (teg *Group) inheritTeamID(newTeamID string) {
	wg := sync.WaitGroup{}
	switch deterministicInheritanceOrder {
	case true:
		// groups
		for i := 0; i < teg.ordNumChildGrp; i++ {
			if child, ok := teg.ordChildrenGrp[i]; ok {
				teg.Children[child].inheritTeamID(newTeamID)
			}
		}
		// clusters
		for i := 0; i < teg.ordNumChildClr; i++ {
			if child, ok := teg.ordChildrenClr[i]; ok {
				teg.Children[child].inheritTeamID(newTeamID)
			}
		}
		// nodes
		for i := 0; i < teg.ordNumChildNod; i++ {
			if child, ok := teg.ordChildrenNod[i]; ok {
				teg.Children[child].inheritTeamID(newTeamID)
			}
		}
	default:
		for child := range teg.Children {
			wg.Add(1)
			go func(name, teamID string) {
				defer wg.Done()
				teg.Children[name].inheritTeamID(teamID)
			}(child, newTeamID)
		}
	}
	teg.Team, _ = uuid.FromString(newTeamID)
	teg.actionRepossess()
	wg.Wait()
}

//
// Interface: BucketAttacher
func (teg *Group) attachToBucket(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentID:   a.ParentID,
		ParentName: a.ParentName,
		ChildType:  teg.Type,
		Group:      teg,
	})

	if teg.Parent == nil {
		a.Root.(*Tree).AttachError(Error{Action: `attach_group`})
		return
	}
	teg.actionCreate()
}

//
// Interface: GroupAttacher
func (teg *Group) attachToGroup(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentID:   a.ParentID,
		ParentName: a.ParentName,
		ChildType:  teg.Type,
		Group:      teg,
	})

	if teg.Parent == nil {
		a.Root.(*Tree).AttachError(Error{Action: `attach_group`})
		return
	}
	teg.actionCreate()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
