/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import (
	"testing"

	"github.com/satori/go.uuid"
)

func TestCheckClone(t *testing.T) {
	check := testSpawnCheck(false, false, false)

	clone := check.Clone()

	if !uuid.Equal(check.ID, clone.ID) {
		t.Errorf(`Illegal clone`)
	}
}

func TestCheckGetter(t *testing.T) {
	check := testSpawnCheck(false, false, true)

	if _, err := uuid.FromString(check.GetSourceCheckID()); err != nil {
		t.Error(`Received error`, err)
	}

	if _, err := uuid.FromString(check.GetCheckConfigID()); err != nil {
		t.Error(`Received error`, err)
	}

	if sourceType := check.GetSourceType(); sourceType == "" {
		t.Error(`Received empty Check.SourceType`)
	}

	if _, err := uuid.FromString(check.GetCapabilityID()); err != nil {
		t.Error(`Received error`, err)
	}

	if view := check.GetView(); view == "" {
		t.Error(`Received empty Check.View`)
	} else {
		switch view {
		case `internal`, `external`, `local`, `any`:
		default:
			t.Error(`Received unknown View`)
		}
	}

	if interval := check.GetInterval(); interval == 0 {
		t.Errorf(`Execution interval is every zero seconds`)
	}

	if child := check.GetChildrenOnly(); child == false {
		t.Errorf(`GetChildren received zero value return`)
	}
}

func TestCheckInherited(t *testing.T) {
	check := testSpawnCheck(true, true, false)

	if inherit := check.GetIsInherited(); inherit != true {
		t.Errorf(`Incorrect inheritance`)
	}

	if inheritance := check.GetInheritance(); inheritance == false {
		t.Errorf(`Inherited check can not have inheritance disabled`)
	}

	var id, idFrom uuid.UUID
	var err error

	if id, err = uuid.FromString(check.GetCheckID()); err != nil {
		t.Error(`Received error`, err)
	}
	if idFrom, err = uuid.FromString(check.GetInheritedFrom()); err != nil {
		t.Error(`Received error`, err)
	}
	if uuid.Equal(id, idFrom) {
		t.Error(`Equal id/sourceId for inherited check`)
	}
}

func TestCheckNotInherited(t *testing.T) {
	check := testSpawnCheck(false, false, false)

	if inherit := check.GetIsInherited(); inherit != false {
		t.Errorf(`Incorrect inheritance`)
	}

	var id, idFrom uuid.UUID
	var err error

	if id, err = uuid.FromString(check.GetCheckID()); err != nil {
		t.Error(`Received error`, err)
	}
	if idFrom, err = uuid.FromString(check.GetInheritedFrom()); err != nil {
		t.Error(`Received error`, err)
	}
	if !uuid.Equal(id, idFrom) {
		t.Error(`Unequal id/sourceId for non-inherited check`)
	}
}

func TestCheckAction(t *testing.T) {
	check := testSpawnCheck(false, false, false)

	action := check.MakeAction()

	if action.Check.CheckID != check.GetCheckID() {
		t.Errorf(`Created action is incorrect`)
	}
}

func TestCheckGetItemNotNil(t *testing.T) {
	check := testSpawnCheck(false, false, false)

	if check.GetCheckID() != check.GetItemID(`node`, uuid.Nil).String() {
		t.Errorf(`GetItemID did not return already set ID`)
	}
}

func TestCheckGetItemNoMatch(t *testing.T) {
	check := testSpawnCheck(false, false, false)
	check.ID = uuid.UUID{}

	if !uuid.Equal(uuid.Nil, check.GetItemID(`node`, uuid.Nil)) {
		t.Errorf(`GetItemID did not return uuid.Nil in non-match case`)
	}
}

func TestCheckGetItem(t *testing.T) {
	check := testSpawnCheck(false, false, false)
	check.ID = uuid.UUID{}

	itemID := uuid.NewV4()
	objID := uuid.NewV4()
	check.Items = append(check.Items, CheckItem{
		ObjectID: func() uuid.UUID {
			ui, _ := uuid.FromString(objID.String())
			return ui
		}(),
		ObjectType: `node`,
		ItemID: func() uuid.UUID {
			ui, _ := uuid.FromString(itemID.String())
			return ui
		}(),
	})

	if !uuid.Equal(itemID, check.GetItemID(`node`, objID)) {
		t.Errorf(`GetItemID did not correctly match objects`)
	}
}

func TestCheckInstanceClone(t *testing.T) {
	check := testSpawnCheck(false, false, false)
	instance := testSpawnCheckInstance(check)

	clone := instance.Clone()

	if !uuid.Equal(instance.InstanceID, clone.InstanceID) {
		t.Errorf(`Faulty checkinstance clone`)
	}
	if !uuid.Equal(instance.CheckID, clone.CheckID) {
		t.Errorf(`Faulty checkinstance clone - CheckID`)
	}
	if !uuid.Equal(instance.ConfigID, clone.ConfigID) {
		t.Errorf(`Faulty checkinstance clone - ConfigID`)
	}
}

func TestCheckInstanceAction(t *testing.T) {
	check := testSpawnCheck(false, false, false)
	instance := testSpawnCheckInstance(check)

	action := instance.MakeAction()

	if action.CheckInstance.InstanceID != instance.InstanceID.String() {
		t.Errorf(`Created instance action is incorrect`)
	}
}

func testSpawnCheckInstance(chk Check) CheckInstance {
	ci := CheckInstance{
		InstanceID: uuid.NewV4(),
		CheckID: func(id string) uuid.UUID {
			f, _ := uuid.FromString(id)
			return f
		}(chk.GetCheckID()),
		ConfigID: func(id string) uuid.UUID {
			f, _ := uuid.FromString(id)
			return f
		}(chk.GetCheckConfigID()),
		InstanceConfigID:    uuid.NewV4(),
		ConstraintOncall:    ``,
		ConstraintService:   map[string]string{},
		ConstraintSystem:    map[string]string{},
		ConstraintCustom:    map[string]string{},
		ConstraintNative:    map[string]string{},
		ConstraintAttribute: map[string]map[string][]string{},
		InstanceService:     `Important Enterprise Business Application`,
		InstanceServiceConfig: map[string]string{
			`port`:            `443`,
			`transport_proto`: `tcp`,
		},
	}

	for _, c := range chk.Constraints {
		switch c.Type {
		case `native`:
			ci.ConstraintNative[c.Key] = c.Value
		case `system`:
			ci.ConstraintSystem[c.Key] = c.Value
		case `custom`:
			ci.ConstraintCustom[c.Key] = c.Value
		case `service`:
			ci.ConstraintService[c.Key] = c.Value
		case `oncall`:
			ci.ConstraintOncall = c.Key
		case `attribute`:
			ci.ConstraintAttribute[`Important Enterprise Business Application`] =
				map[string][]string{
					`port`: []string{`80`, `443`},
				}
		}
	}
	ci.calcConstraintHash()
	ci.calcConstraintValHash()
	ci.calcInstanceSvcCfgHash()

	return ci
}

func testSpawnCheck(inherited, inheritance, childrenOnly bool) Check {
	id := uuid.NewV4()
	var idFrom uuid.UUID
	if inherited {
		idFrom = uuid.NewV4()
	} else {
		idFrom, _ = uuid.FromString(id.String())
	}

	return Check{
		Id:            id,
		SourceID:      uuid.NewV4(),
		SourceType:    `sourceType`,
		Inherited:     inherited,
		InheritedFrom: idFrom,
		CapabilityID:  uuid.NewV4(),
		ConfigID:      uuid.NewV4(),
		Inheritance:   inheritance,
		ChildrenOnly:  childrenOnly,
		View:          `any`,
		Interval:      1,
		Thresholds: []CheckThreshold{
			{
				Predicate: `>=`,
				Level:     1,
				Value:     1,
			},
			{
				Predicate: `>=`,
				Level:     3,
				Value:     5,
			},
		},
		Constraints: []CheckConstraint{
			{
				Type:  `native`,
				Key:   `object_type`,
				Value: `node`,
			},
			{
				Type:  `system`,
				Key:   `fqdn`,
				Value: `host.example.org`,
			},
			{
				Type:  `custom`,
				Key:   `47db9a82-83f2-11e6-bc90-f8a963a55ba6`,
				Value: `foobar`,
			},
			{
				Type:  `service`,
				Key:   `name`,
				Value: `Important Enterprise Business Application`,
			},
			{
				Type:  `attribute`,
				Key:   `port`,
				Value: `@defined`,
			},
			{
				Type:  `oncall`,
				Key:   `70c1570b-83f4-11e6-bc90-f8a963a55ba6`,
				Value: `Heroes of Oncall Duty`,
			},
		},
		Items: []CheckItem{
			{
				ObjectID:   uuid.NewV4(),
				ObjectType: `objectType`,
				ItemID:     uuid.NewV4(),
			},
		},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
