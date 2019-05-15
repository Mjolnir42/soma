/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import (
	uuid "github.com/satori/go.uuid"
)

//
// Interface: Attacher
func (ten *Node) Attach(a AttachRequest) {
	if ten.Parent != nil {
		panic(`Node.Attach: already attached`)
	}
	switch {
	case a.ParentType == "bucket":
		ten.attachToBucket(a)
	case a.ParentType == "group":
		ten.attachToGroup(a)
	case a.ParentType == "cluster":
		ten.attachToCluster(a)
	default:
		panic(`Node.Attach`)
	}

	ten.Parent.(Propertier).syncProperty(ten.ID.String())
	ten.Parent.(Checker).syncCheck(ten.ID.String())
}

func (ten *Node) ReAttach(a AttachRequest) {
	if ten.Parent == nil {
		panic(`Node.ReAttach: not attached`)
	}
	ten.deletePropertyAllInherited()
	ten.deleteCheckAllInherited()
	ten.updateCheckInstances()

	ten.Parent.Unlink(UnlinkRequest{
		ParentType: ten.Parent.(Builder).GetType(),
		ParentName: ten.Parent.(Builder).GetName(),
		ParentID:   ten.Parent.(Builder).GetID(),
		ChildType:  ten.GetType(),
		ChildName:  ten.GetName(),
		ChildID:    ten.GetID(),
	},
	)

	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentID:   a.ParentID,
		ParentName: a.ParentName,
		ChildType:  ten.GetType(),
		Node:       ten,
	},
	)

	if ten.Parent == nil {
		panic(`Node.ReAttach: not reattached`)
	}
	ten.actionUpdate()
	ten.Parent.(Propertier).syncProperty(ten.ID.String())
	ten.Parent.(Checker).syncCheck(ten.ID.String())
}

func (ten *Node) Destroy() {
	if ten.Parent == nil {
		panic(`Node.Destroy called without Parent to unlink from`)
	}
	// call before unlink since it requires tec.Parent.*
	ten.State = `unassigned`
	ten.actionDelete()
	ten.deletePropertyAllLocal()
	ten.deletePropertyAllInherited()
	ten.deleteCheckLocalAll()
	ten.deleteCheckAllInherited()
	ten.updateCheckInstances()

	ten.Parent.Unlink(UnlinkRequest{
		ParentType: ten.Parent.(Builder).GetType(),
		ParentID:   ten.Parent.(Builder).GetID(),
		ParentName: ten.Parent.(Builder).GetName(),
		ChildType:  ten.GetType(),
		ChildName:  ten.GetName(),
		ChildID:    ten.GetID(),
	},
	)

	ten.setFault(nil)
	ten.setAction(nil)
}

func (ten *Node) Detach() {
	if ten.Parent == nil {
		panic(`Node.Detach called without Parent to detach from`)
	}
	bucket := ten.Parent.(Bucketeer).GetBucket()

	ten.deletePropertyAllInherited()
	ten.deleteCheckAllInherited()
	ten.updateCheckInstances()

	ten.Parent.Unlink(UnlinkRequest{
		ParentType: ten.Parent.(Builder).GetType(),
		ParentID:   ten.Parent.(Builder).GetID(),
		ParentName: ten.Parent.(Builder).GetName(),
		ChildType:  ten.GetType(),
		ChildName:  ten.GetName(),
		ChildID:    ten.GetID(),
	},
	)

	bucket.Receive(ReceiveRequest{
		ParentType: bucket.(Builder).GetType(),
		ParentID:   bucket.(Builder).GetID(),
		ParentName: bucket.(Builder).GetName(),
		ChildType:  ten.Type,
		Node:       ten,
	},
	)

	ten.actionUpdate()
	ten.Parent.(Propertier).syncProperty(ten.ID.String())
}

func (ten *Node) SetName(s string) {
	ten.Name = s
	ten.actionRename()
}

func (ten *Node) inheritTeamID(newTeamID string) {
	ten.Team, _ = uuid.FromString(newTeamID)
	ten.actionRepossess()
}

//
// Interface: BucketAttacher
func (ten *Node) attachToBucket(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentID:   a.ParentID,
		ParentName: a.ParentName,
		ChildType:  ten.Type,
		Node:       ten,
	})

	if ten.Parent == nil {
		a.Root.(*Tree).AttachError(Error{Action: `attach_node`})
		return
	}
	ten.actionUpdate()
}

//
// Interface: GroupAttacher
func (ten *Node) attachToGroup(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentID:   a.ParentID,
		ParentName: a.ParentName,
		ChildType:  ten.Type,
		Node:       ten,
	})

	if ten.Parent == nil {
		a.Root.(*Tree).AttachError(Error{Action: `attach_node`})
		return
	}
	ten.actionUpdate()
}

//
// Interface: ClusterAttacher
func (ten *Node) attachToCluster(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentID:   a.ParentID,
		ParentName: a.ParentName,
		ChildType:  ten.Type,
		Node:       ten,
	})

	if ten.Parent == nil {
		a.Root.(*Tree).AttachError(Error{Action: `attach_node`})
		return
	}
	ten.actionUpdate()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
