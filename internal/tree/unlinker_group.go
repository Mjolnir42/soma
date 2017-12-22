/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

//
// Interface: Unlinker
func (teg *Group) Unlink(u UnlinkRequest) {
	if unlinkRequestCheck(u, teg) {
		switch u.ChildType {
		case "group":
			teg.unlinkGroup(u)
		case "cluster":
			teg.unlinkCluster(u)
		case "node":
			teg.unlinkNode(u)
		default:
			panic(`Group.Unlink`)
		}
		return
	}
loop:
	for child := range teg.Children {
		if teg.Children[child].(Builder).GetType() == "node" {
			continue loop
		}
		teg.Children[child].(Unlinker).Unlink(u)
	}
}

//
// Interface: GroupUnlinker
func (teg *Group) unlinkGroup(u UnlinkRequest) {
	if unlinkRequestCheck(u, teg) {
		switch u.ChildType {
		case "group":
			if _, ok := teg.Children[u.ChildID]; ok {
				if u.ChildName == teg.Children[u.ChildID].GetName() {
					a := Action{
						ChildType:  "group",
						ChildGroup: teg.Children[u.ChildID].(*Group).export(),
					}

					teg.Children[u.ChildID].clearParent()
					delete(teg.Children, u.ChildID)
					for i, grp := range teg.ordChildrenGrp {
						if grp == u.ChildID {
							delete(teg.ordChildrenGrp, i)
						}
					}

					teg.actionMemberRemoved(a)
				}
			}
		default:
			panic(`Group.unlinkGroup`)
		}
		return
	}
	panic(`Group.unlinkGroup`)
}

//
// Interface: ClusterUnlinker
func (teg *Group) unlinkCluster(u UnlinkRequest) {
	if unlinkRequestCheck(u, teg) {
		switch u.ChildType {
		case "cluster":
			if _, ok := teg.Children[u.ChildID]; ok {
				if u.ChildName == teg.Children[u.ChildID].GetName() {
					a := Action{
						ChildType:    "cluster",
						ChildCluster: teg.Children[u.ChildID].(*Cluster).export(),
					}

					teg.Children[u.ChildID].clearParent()
					delete(teg.Children, u.ChildID)
					for i, clr := range teg.ordChildrenClr {
						if clr == u.ChildID {
							delete(teg.ordChildrenClr, i)
						}
					}

					teg.actionMemberRemoved(a)
				}
			}
		default:
			panic(`Group.unlinkCluster`)
		}
		return
	}
	panic(`Group.unlinkCluster`)
}

//
// Interface: NodeUnlinker
func (teg *Group) unlinkNode(u UnlinkRequest) {
	if unlinkRequestCheck(u, teg) {
		switch u.ChildType {
		case "node":
			if _, ok := teg.Children[u.ChildID]; ok {
				if u.ChildName == teg.Children[u.ChildID].GetName() {
					a := Action{
						ChildType: "node",
						ChildNode: teg.Children[u.ChildID].(*Node).export(),
					}

					teg.Children[u.ChildID].clearParent()
					delete(teg.Children, u.ChildID)
					for i, nod := range teg.ordChildrenNod {
						if nod == u.ChildID {
							delete(teg.ordChildrenNod, i)
						}
					}

					teg.actionMemberRemoved(a)
				}
			}
		default:
			panic(`Group.unlinkNode`)
		}
		return
	}
	panic(`Group.unlinkNode`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
