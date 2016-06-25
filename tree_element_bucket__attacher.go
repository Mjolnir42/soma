package somatree

import (
	"fmt"
	"reflect"
	"sync"
)

//
// Interface: SomaTreeAttacher
func (teb *SomaTreeElemBucket) Attach(a AttachRequest) {
	if teb.Parent != nil {
		panic(`SomaTreeElemBucket.Attach: already attached`)
	}
	switch {
	case a.ParentType == "repository":
		teb.attachToRepository(a)
	default:
		panic(`SomaTreeElemBucket.Attach`)
	}

	if teb.Parent == nil {
		panic(`SomaTreeElemBucket.Attach: failed`)
	}
	teb.Parent.(Propertier).syncProperty(teb.Id.String())
}

func (teb *SomaTreeElemBucket) Destroy() {
	if teb.Parent == nil {
		panic(`SomaTreeElemBucket.Destroy called without Parent to unlink from`)
	}
	// XXX: destroy all inherited properties before unlinking
	// teb.(SomaTreePropertier).destroyInheritedProperties()

	teb.Parent.Unlink(UnlinkRequest{
		ParentType: teb.Parent.(Builder).GetType(),
		ParentId:   teb.Parent.(Builder).GetID(),
		ParentName: teb.Parent.(Builder).GetName(),
		ChildType:  teb.GetType(),
		ChildName:  teb.GetName(),
		ChildId:    teb.GetID(),
	},
	)

	teb.setFault(nil)
	teb.actionDelete()
	teb.setAction(nil)
}

func (teb *SomaTreeElemBucket) Detach() {
	teb.Destroy()
}

func (teb *SomaTreeElemBucket) clearParent() {
	teb.Parent = nil
	teb.State = "floating"
}

func (teb *SomaTreeElemBucket) setFault(f *SomaTreeElemFault) {
	teb.Fault = f
}

func (teb *SomaTreeElemBucket) updateParentRecursive(p SomaTreeReceiver) {
	teb.setParent(p)
	var wg sync.WaitGroup
	for child, _ := range teb.Children {
		wg.Add(1)
		c := child
		go func(str SomaTreeReceiver) {
			defer wg.Done()
			teb.Children[c].updateParentRecursive(str)
		}(teb)
	}
	wg.Wait()
}

func (teb *SomaTreeElemBucket) updateFaultRecursive(f *SomaTreeElemFault) {
	teb.setFault(f)
	var wg sync.WaitGroup
	for child, _ := range teb.Children {
		wg.Add(1)
		c := child
		go func(ptr *SomaTreeElemFault) {
			defer wg.Done()
			teb.Children[c].updateFaultRecursive(ptr)
		}(f)
	}
	wg.Wait()
}

func (teb *SomaTreeElemBucket) setParent(p SomaTreeReceiver) {
	switch p.(type) {
	case SomaTreeBucketReceiver:
		teb.setBucketParent(p.(SomaTreeBucketReceiver))
		teb.State = "attached"
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`SomaTreeElemBucket.setParent`)
	}
}

func (teb *SomaTreeElemBucket) setBucketParent(p SomaTreeBucketReceiver) {
	teb.Parent = p
}

//
// Interface: SomaTreeRepositoryAttacher
func (teb *SomaTreeElemBucket) attachToRepository(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  teb.Type,
		Bucket:     teb,
	})

	if teb.Parent == nil {
		a.Root.(*SomaTree).AttachError(Error{Action: `attach_bucket`})
		return
	}
	teb.actionCreate()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
