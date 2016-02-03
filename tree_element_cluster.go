package somatree

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/satori/go.uuid"
)

type SomaTreeElemCluster struct {
	Id       uuid.UUID
	Name     string
	State    string
	Team     uuid.UUID
	Type     string
	Parent   SomaTreeClusterReceiver            `json:"-"`
	Fault    *SomaTreeElemFault                 `json:"-"`
	Children map[string]SomaTreeClusterAttacher `json:"-"`
	//PropertyOncall  map[string]*SomaTreePropertyOncall
	//PropertyService map[string]*SomaTreePropertyService
	//PropertySystem  map[string]*SomaTreePropertySystem
	//PropertyCustom  map[string]*SomaTreePropertyCustom
	//Checks          map[string]*SomaTreeCheck
}

type ClusterSpec struct {
	Id   uuid.UUID
	Name string
	Team uuid.UUID
}

//
// NEW
func NewCluster(name string) *SomaTreeElemCluster {
	tec := new(SomaTreeElemCluster)
	tec.Id = uuid.NewV4()
	tec.Name = name
	tec.Type = "cluster"
	tec.State = "floating"
	tec.Children = make(map[string]SomaTreeClusterAttacher)
	//tec.PropertyOncall = make(map[string]*SomaTreePropertyOncall)
	//tec.PropertyService = make(map[string]*SomaTreePropertyService)
	//tec.PropertySystem = make(map[string]*SomaTreePropertySystem)
	//tec.PropertyCustom = make(map[string]*SomaTreePropertyCustom)
	//tec.Checks = make(map[string]*SomaTreeCheck)

	return tec
}

func (tec SomaTreeElemCluster) CloneBucket() SomaTreeBucketAttacher {
	for k, child := range tec.Children {
		tec.Children[k] = child.CloneCluster()
	}
	return &tec
}

func (tec SomaTreeElemCluster) CloneGroup() SomaTreeGroupAttacher {
	f := make(map[string]SomaTreeClusterAttacher)
	for k, child := range tec.Children {
		f[k] = child.CloneCluster()
	}
	tec.Children = f
	return &tec
}

//
// Interface: SomaTreeBuilder
func (tec *SomaTreeElemCluster) GetID() string {
	return tec.Id.String()
}

func (tec *SomaTreeElemCluster) GetName() string {
	return tec.Name
}

func (tec *SomaTreeElemCluster) GetType() string {
	return tec.Type
}

//
// Interface: SomaTreeAttacher
func (tec *SomaTreeElemCluster) Attach(a AttachRequest) {
	switch {
	case a.ParentType == "bucket":
		tec.attachToBucket(a)
	case a.ParentType == "group":
		tec.attachToGroup(a)
	default:
		panic(`SomaTreeElemCluster.Attach`)
	}
}

func (tec *SomaTreeElemCluster) ReAttach(a AttachRequest) {
	if tec.Parent == nil {
		panic(`SomaTreeElemGroup.ReAttach: not attached`)
	}
	tec.Parent.Unlink(UnlinkRequest{
		ParentType: tec.Parent.(SomaTreeBuilder).GetType(),
		ParentName: tec.Parent.(SomaTreeBuilder).GetName(),
		ParentId:   tec.Parent.(SomaTreeBuilder).GetID(),
		ChildType:  tec.GetType(),
		ChildName:  tec.GetName(),
		ChildId:    tec.GetID(),
	},
	)

	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  tec.GetType(),
		Cluster:    tec,
	},
	)
}

func (tec *SomaTreeElemCluster) setParent(p SomaTreeReceiver) {
	switch p.(type) {
	case *SomaTreeElemBucket:
		tec.setClusterParent(p.(SomaTreeClusterReceiver))
		tec.State = "standalone"
	case *SomaTreeElemGroup:
		tec.setClusterParent(p.(SomaTreeClusterReceiver))
		tec.State = "grouped"
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`SomaTreeElemCluster.setParent`)
	}
}

func (tec *SomaTreeElemCluster) updateParentRecursive(p SomaTreeReceiver) {
	tec.setParent(p)
	var wg sync.WaitGroup
	for child, _ := range tec.Children {
		wg.Add(1)
		c := child
		go func(str SomaTreeReceiver) {
			defer wg.Done()
			tec.Children[c].updateParentRecursive(str)
		}(tec)
	}
	wg.Wait()
}

// SomaTreeClusterReceiver == can receive Clusters as children
func (tec *SomaTreeElemCluster) setClusterParent(p SomaTreeClusterReceiver) {
	tec.Parent = p
}

func (tec *SomaTreeElemCluster) clearParent() {
	tec.Parent = nil
	tec.State = "floating"
}

func (tec *SomaTreeElemCluster) setFault(f *SomaTreeElemFault) {
	tec.Fault = f
}

func (tec *SomaTreeElemCluster) updateFaultRecursive(f *SomaTreeElemFault) {
	tec.setFault(f)
	var wg sync.WaitGroup
	for child, _ := range tec.Children {
		wg.Add(1)
		c := child
		go func(ptr *SomaTreeElemFault) {
			defer wg.Done()
			tec.Children[c].updateFaultRecursive(ptr)
		}(f)
	}
	wg.Wait()
}

func (tec *SomaTreeElemCluster) Destroy() {
	if tec.Parent == nil {
		panic(`SomaTreeElemCluster.Destroy called without Parent to unlink from`)
	}

	tec.Parent.Unlink(UnlinkRequest{
		ParentType: tec.Parent.(SomaTreeBuilder).GetType(),
		ParentId:   tec.Parent.(SomaTreeBuilder).GetID(),
		ParentName: tec.Parent.(SomaTreeBuilder).GetName(),
		ChildType:  tec.GetType(),
		ChildName:  tec.GetName(),
		ChildId:    tec.GetID(),
	},
	)

	tec.setFault(nil)
}

func (tec *SomaTreeElemCluster) Detach() {
	if tec.Parent == nil {
		panic(`SomaTreeElemCluster.Detach called without Parent to detach from`)
	}
	bucket := tec.Parent.(SomaTreeBucketeer).GetBucket()

	tec.Parent.Unlink(UnlinkRequest{
		ParentType: tec.Parent.(SomaTreeBuilder).GetType(),
		ParentId:   tec.Parent.(SomaTreeBuilder).GetID(),
		ParentName: tec.Parent.(SomaTreeBuilder).GetName(),
		ChildType:  tec.GetType(),
		ChildName:  tec.GetName(),
		ChildId:    tec.GetID(),
	},
	)

	bucket.Receive(ReceiveRequest{
		ParentType: bucket.(SomaTreeBuilder).GetType(),
		ParentId:   bucket.(SomaTreeBuilder).GetID(),
		ParentName: bucket.(SomaTreeBuilder).GetName(),
		ChildType:  tec.Type,
		Cluster:    tec,
	},
	)
}

//
// Interface: SomaTreeBucketAttacher
func (tec *SomaTreeElemCluster) attachToBucket(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  tec.Type,
		Cluster:    tec,
	})
}

//
// Interface: SomaTreeGroupAttacher
func (tec *SomaTreeElemCluster) attachToGroup(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  tec.Type,
		Cluster:    tec,
	})
}

//
// Interface: SomaTreeReceiver
func (tec *SomaTreeElemCluster) Receive(r ReceiveRequest) {
	if receiveRequestCheck(r, tec) {
		switch r.ChildType {
		case "node":
			tec.receiveNode(r)
		default:
			panic(`SomaTreeElemCluster.Receive`)
		}
	}
	// no passing along since only nodes are a SomeTreeClusterAttacher
	// and nodes can have no children
	return
}

//
// Interface: SomaTreeBucketeer
func (tec *SomaTreeElemCluster) GetBucket() SomaTreeReceiver {
	if tec.Parent == nil {
		if tec.Fault == nil {
			panic(`SomaTreeElemCluster.GetBucket called without Parent`)
		} else {
			return tec.Fault
		}
	}
	return tec.Parent.(SomaTreeBucketeer).GetBucket()
}

//
// Interface: SomaTreeUnlinker
func (tec *SomaTreeElemCluster) Unlink(u UnlinkRequest) {
	if unlinkRequestCheck(u, tec) {
		switch u.ChildType {
		case "node":
			tec.unlinkNode(u)
		default:
			panic(`SomaTreeElemCluster.Unlink`)
		}
	}
	// no passing along since only nodes are a SomeTreeClusterAttacher
	// and nodes can have no children
	return
}

//
// Interface: SomaTreeNodeReceiver
func (tec *SomaTreeElemCluster) receiveNode(r ReceiveRequest) {
	if receiveRequestCheck(r, tec) {
		switch r.ChildType {
		case "node":
			tec.Children[r.Node.GetID()] = r.Node
			r.Node.setParent(tec)
			r.Node.setFault(tec.Fault)
		default:
			panic(`SomaTreeElemCluster.receiveNode`)
		}
		return
	}
	panic(`SomaTreeElemCluster.receiveNode`)
}

//
// Interface: SomaTreeNodeUnlinker
func (tec *SomaTreeElemCluster) unlinkNode(u UnlinkRequest) {
	if unlinkRequestCheck(u, tec) {
		switch u.ChildType {
		case "node":
			if _, ok := tec.Children[u.ChildId]; ok {
				if u.ChildName == tec.Children[u.ChildId].GetName() {
					tec.Children[u.ChildId].clearParent()
					delete(tec.Children, u.ChildId)
				}
			}
		default:
			panic(`SomaTreeElemCluster.unlinkNode`)
		}
		return
	}
	panic(`SomaTreeElemCluster.unlinkNode`)
}

//
// Interface: SomaTreeFinder
func (tec *SomaTreeElemCluster) Find(f FindRequest, b bool) SomaTreeAttacher {
	if findRequestCheck(f, tec) {
		return tec
	}
	var (
		wg             sync.WaitGroup
		rawResult, res chan SomaTreeAttacher
	)
	if len(tec.Children) == 0 {
		goto skip
	}
	if f.ElementId != "" {
		if _, ok := tec.Children[f.ElementId]; ok {
			return tec.Children[f.ElementId]
		} else {
			// f.ElementId is not a child of ours
			goto skip
		}
	}
	rawResult = make(chan SomaTreeAttacher, len(tec.Children))
	for child, _ := range tec.Children {
		wg.Add(1)
		c := child
		go func(fr FindRequest, bl bool) {
			defer wg.Done()
			rawResult <- tec.Children[c].(SomaTreeFinder).Find(fr, bl)
		}(f, false)
	}
	wg.Wait()
	close(rawResult)

	res = make(chan SomaTreeAttacher, len(rawResult))
	for sta := range rawResult {
		if sta != nil {
			res <- sta
		}
	}
	close(res)
skip:
	switch {
	case len(res) == 0:
		if b {
			return tec.Fault
		} else {
			return nil
		}
	case len(res) > 1:
		return tec.Fault
	}
	return <-res
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
