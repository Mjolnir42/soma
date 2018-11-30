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
func (tec *Cluster) Attach(a AttachRequest) {
	if tec.Parent != nil {
		panic(`Cluster.Attach: already attached`)
	}
	switch {
	case a.ParentType == "bucket":
		tec.attachToBucket(a)
	case a.ParentType == "group":
		tec.attachToGroup(a)
	default:
		panic(`Cluster.Attach`)
	}

	tec.Parent.(Propertier).syncProperty(tec.ID.String())
	tec.Parent.(Checker).syncCheck(tec.ID.String())
}

func (tec *Cluster) ReAttach(a AttachRequest) {
	if tec.Parent == nil {
		panic(`Cluster.ReAttach: not attached`)
	}
	tec.deletePropertyAllInherited()
	// TODO delete all inherited checks + check instances

	tec.Parent.Unlink(UnlinkRequest{
		ParentType: tec.Parent.(Builder).GetType(),
		ParentName: tec.Parent.(Builder).GetName(),
		ParentID:   tec.Parent.(Builder).GetID(),
		ChildType:  tec.GetType(),
		ChildName:  tec.GetName(),
		ChildID:    tec.GetID(),
	},
	)

	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentID:   a.ParentID,
		ParentName: a.ParentName,
		ChildType:  tec.GetType(),
		Cluster:    tec,
	},
	)

	if tec.Parent == nil {
		panic(`Group.ReAttach: not reattached`)
	}
	tec.actionUpdate()
	tec.Parent.(Propertier).syncProperty(tec.ID.String())
	tec.Parent.(Checker).syncCheck(tec.ID.String())
}

func (tec *Cluster) Destroy() {
	if tec.Parent == nil {
		panic(`Cluster.Destroy called without Parent to unlink from`)
	}

	// call before unlink since it requires tec.Parent.*
	tec.actionDelete()
	tec.deletePropertyAllLocal()
	tec.deletePropertyAllInherited()
	tec.deleteCheckLocalAll()
	tec.updateCheckInstances()

	wg := new(sync.WaitGroup)
	for child := range tec.Children {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			tec.Children[c].Destroy()
		}(child)
	}
	wg.Wait()

	tec.Parent.Unlink(UnlinkRequest{
		ParentType: tec.Parent.(Builder).GetType(),
		ParentID:   tec.Parent.(Builder).GetID(),
		ParentName: tec.Parent.(Builder).GetName(),
		ChildType:  tec.GetType(),
		ChildName:  tec.GetName(),
		ChildID:    tec.GetID(),
	},
	)

	tec.setFault(nil)
	tec.setAction(nil)
}

func (tec *Cluster) Detach() {
	if tec.Parent == nil {
		panic(`Cluster.Detach called without Parent to detach from`)
	}
	bucket := tec.Parent.(Bucketeer).GetBucket()

	tec.deletePropertyAllInherited()
	// TODO delete all inherited checks + check instances

	tec.Parent.Unlink(UnlinkRequest{
		ParentType: tec.Parent.(Builder).GetType(),
		ParentID:   tec.Parent.(Builder).GetID(),
		ParentName: tec.Parent.(Builder).GetName(),
		ChildType:  tec.GetType(),
		ChildName:  tec.GetName(),
		ChildID:    tec.GetID(),
	},
	)

	bucket.Receive(ReceiveRequest{
		ParentType: bucket.(Builder).GetType(),
		ParentID:   bucket.(Builder).GetID(),
		ParentName: bucket.(Builder).GetName(),
		ChildType:  tec.Type,
		Cluster:    tec,
	},
	)

	tec.actionUpdate()
	tec.Parent.(Propertier).syncProperty(tec.ID.String())
}

func (tec *Cluster) SetName(s string) {
	tec.Name = s
	tec.actionRename()
}

func (tec *Cluster) inheritTeamID(newTeamID string) {
	wg := sync.WaitGroup{}
	switch deterministicInheritanceOrder {
	case true:
		// nodes
		for i := 0; i < tec.ordNumChildNod; i++ {
			if child, ok := tec.ordChildrenNod[i]; ok {
				tec.Children[child].inheritTeamID(newTeamID)
			}
		}
	default:
		for child := range tec.Children {
			wg.Add(1)
			go func(name, teamID string) {
				defer wg.Done()
				tec.Children[name].inheritTeamID(teamID)
			}(child, newTeamID)
		}
	}
	tec.Team, _ = uuid.FromString(newTeamID)
	tec.actionRepossess()
	wg.Wait()
}

//
// Interface: BucketAttacher
func (tec *Cluster) attachToBucket(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentID:   a.ParentID,
		ParentName: a.ParentName,
		ChildType:  tec.Type,
		Cluster:    tec,
	})

	if tec.Parent == nil {
		a.Root.(*Tree).AttachError(Error{Action: `attach_cluster`})
		return
	}
	tec.actionCreate()
}

//
// Interface: GroupAttacher
func (tec *Cluster) attachToGroup(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentID:   a.ParentID,
		ParentName: a.ParentName,
		ChildType:  tec.Type,
		Cluster:    tec,
	})

	if tec.Parent == nil {
		a.Root.(*Tree).AttachError(Error{Action: `attach_cluster`})
		return
	}
	tec.actionCreate()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
