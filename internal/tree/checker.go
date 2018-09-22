/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"io"
	"sort"

	"github.com/mjolnir42/soma/lib/proto"

	"github.com/satori/go.uuid"
)

type Checker interface {
	SetCheck(c Check)
	LoadInstance(i CheckInstance)
	DeleteCheck(c Check)

	setCheckInherited(c Check)
	setCheckOnChildren(c Check)
	addCheck(c Check)

	deleteCheckInherited(c Check)
	deleteCheckOnChildren(c Check)
	rmCheck(c Check)

	syncCheck(childID string)
	checkCheck(checkID string) bool
}

type CheckGetter interface {
	GetCheckID() string
	GetSourceCheckID() string
	GetCheckConfigID() string
	GetSourceType() string
	GetIsInherited() bool
	GetInheritedFrom() string
	GetInheritance() bool
	GetChildrenOnly() bool
	GetView() string
	GetCapabilityID() string
	GetInterval() uint64
	GetItemID(objType string, objID uuid.UUID) uuid.UUID
}

type Check struct {
	ID            uuid.UUID
	SourceID      uuid.UUID
	SourceType    string
	Inherited     bool
	InheritedFrom uuid.UUID
	CapabilityID  uuid.UUID
	ConfigID      uuid.UUID
	Inheritance   bool
	ChildrenOnly  bool
	View          string
	Interval      uint64
	Thresholds    []CheckThreshold
	Constraints   []CheckConstraint
	Items         []CheckItem
}

func (c *Check) Clone() Check {
	ng := Check{
		SourceType:   c.SourceType,
		Inherited:    c.Inherited,
		Inheritance:  c.Inheritance,
		ChildrenOnly: c.ChildrenOnly,
		View:         c.View,
		Interval:     c.Interval,
	}
	ng.ID, _ = uuid.FromString(c.ID.String())
	ng.SourceID, _ = uuid.FromString(c.SourceID.String())
	ng.InheritedFrom, _ = uuid.FromString(c.InheritedFrom.String())
	ng.CapabilityID, _ = uuid.FromString(c.CapabilityID.String())
	ng.ConfigID, _ = uuid.FromString(c.ConfigID.String())

	ng.Thresholds = make([]CheckThreshold, len(c.Thresholds))
	for i := range c.Thresholds {
		ng.Thresholds[i] = c.Thresholds[i].Clone()
	}

	ng.Constraints = make([]CheckConstraint, len(c.Constraints))
	for i := range c.Constraints {
		ng.Constraints[i] = c.Constraints[i].Clone()
	}

	ng.Items = make([]CheckItem, len(c.Items))
	for i := range c.Items {
		ng.Items[i] = c.Items[i].Clone()
	}

	return ng
}

type CheckItem struct {
	ObjectID   uuid.UUID
	ObjectType string
	ItemID     uuid.UUID
}

func (ci *CheckItem) Clone() CheckItem {
	oid, _ := uuid.FromString(ci.ObjectID.String())
	iid, _ := uuid.FromString(ci.ItemID.String())
	return CheckItem{
		ObjectID:   oid,
		ObjectType: ci.ObjectType,
		ItemID:     iid,
	}
}

type CheckThreshold struct {
	Predicate string
	Level     uint8
	Value     int64
}

func (ct *CheckThreshold) Clone() CheckThreshold {
	return CheckThreshold{
		Predicate: ct.Predicate,
		Level:     ct.Level,
		Value:     ct.Value,
	}
}

type CheckConstraint struct {
	Type  string
	Key   string
	Value string
}

func (cc *CheckConstraint) Clone() CheckConstraint {
	return CheckConstraint{
		Type:  cc.Type,
		Key:   cc.Key,
		Value: cc.Value,
	}
}

type CheckInstance struct {
	InstanceID            uuid.UUID
	CheckID               uuid.UUID
	ConfigID              uuid.UUID
	InstanceConfigID      uuid.UUID
	Version               uint64
	ConstraintHash        string
	ConstraintValHash     string
	ConstraintOncall      string                         // Ids
	ConstraintService     map[string]string              // svcName->value
	ConstraintSystem      map[string]string              // Id->value
	ConstraintCustom      map[string]string              // Id->value
	ConstraintNative      map[string]string              // prop->value
	ConstraintAttribute   map[string]map[string][]string // svcID->attr->[ value, value, ... ]
	InstanceServiceConfig map[string]string              // attr->value
	InstanceService       string
	InstanceSvcCfgHash    string
	oldConstraintHash     string
	oldConstraintValHash  string
	oldInstanceSvcCfgHash string
}

func (c *Check) GetItemID(objType string, objID uuid.UUID) uuid.UUID {
	if !uuid.Equal(c.ID, uuid.Nil) {
		return c.ID
	}
	for _, item := range c.Items {
		if objType == item.ObjectType && uuid.Equal(item.ObjectID, objID) {
			return item.ItemID
		}
	}
	return uuid.Nil
}

func (c *Check) GetCheckID() string {
	return c.ID.String()
}

func (c *Check) GetSourceCheckID() string {
	return c.SourceID.String()
}

func (c *Check) GetCheckConfigID() string {
	return c.ConfigID.String()
}

func (c *Check) GetSourceType() string {
	return c.SourceType
}

func (c *Check) GetIsInherited() bool {
	return c.Inherited
}

func (c *Check) GetInheritedFrom() string {
	return c.InheritedFrom.String()
}

func (c *Check) GetInheritance() bool {
	return c.Inheritance
}

func (c *Check) GetChildrenOnly() bool {
	return c.ChildrenOnly
}

func (c *Check) GetView() string {
	return c.View
}

func (c *Check) GetCapabilityID() string {
	return c.CapabilityID.String()
}

func (c *Check) GetInterval() uint64 {
	return c.Interval
}

func (c *Check) MakeAction() Action {
	return Action{
		Check: proto.Check{
			CheckID:       c.GetCheckID(),
			SourceCheckID: c.GetSourceCheckID(),
			CheckConfigID: c.GetCheckConfigID(),
			SourceType:    c.GetSourceType(),
			IsInherited:   c.GetIsInherited(),
			InheritedFrom: c.GetInheritedFrom(),
			Inheritance:   c.GetInheritance(),
			ChildrenOnly:  c.GetChildrenOnly(),
			CapabilityID:  c.GetCapabilityID(),
		},
	}
}

func (tci *CheckInstance) Clone() CheckInstance {
	cl := CheckInstance{
		Version:            tci.Version,
		ConstraintHash:     tci.ConstraintHash,
		ConstraintValHash:  tci.ConstraintValHash,
		ConstraintOncall:   tci.ConstraintOncall,
		InstanceSvcCfgHash: tci.InstanceSvcCfgHash,
		InstanceService:    tci.InstanceService,
	}
	cl.InstanceConfigID, _ = uuid.FromString(tci.InstanceConfigID.String())
	cl.InstanceID, _ = uuid.FromString(tci.InstanceID.String())
	cl.CheckID, _ = uuid.FromString(tci.CheckID.String())
	cl.ConfigID, _ = uuid.FromString(tci.ConfigID.String())
	cl.ConstraintService = make(map[string]string)
	for k, v := range tci.ConstraintService {
		t := v
		cl.ConstraintService[k] = t
	}
	cl.ConstraintSystem = make(map[string]string)
	for k, v := range tci.ConstraintSystem {
		t := v
		cl.ConstraintSystem[k] = t
	}
	cl.ConstraintCustom = make(map[string]string)
	for k, v := range tci.ConstraintCustom {
		t := v
		cl.ConstraintCustom[k] = t
	}
	cl.ConstraintNative = make(map[string]string)
	for k, v := range tci.ConstraintNative {
		t := v
		cl.ConstraintNative[k] = t
	}
	cl.InstanceServiceConfig = make(map[string]string)
	for k, v := range tci.InstanceServiceConfig {
		t := v
		cl.InstanceServiceConfig[k] = t
	}
	cl.ConstraintAttribute = make(map[string]map[string][]string, 0)
	for k := range tci.ConstraintAttribute {
		cl.ConstraintAttribute[k] = make(map[string][]string)
		for k2, aVal := range tci.ConstraintAttribute[k] {
			for _, val := range aVal {
				t := val
				cl.ConstraintAttribute[k][k2] = append(cl.ConstraintAttribute[k][k2], t)
			}
		}
	}

	return cl
}

func (tci *CheckInstance) calcConstraintHash() {
	h := sha512.New()
	io.WriteString(h, tci.ConstraintOncall)

	services := []string{}
	for i := range tci.ConstraintService {
		j := i
		services = append(services, j)
	}
	sort.Strings(services)
	for _, i := range services {
		io.WriteString(h, i)
	}

	systems := []string{}
	for i := range tci.ConstraintSystem {
		j := i
		systems = append(systems, j)
	}
	sort.Strings(systems)
	for _, i := range systems {
		io.WriteString(h, i)
	}

	customs := []string{}
	for i := range tci.ConstraintCustom {
		j := i
		customs = append(customs, j)
	}
	sort.Strings(customs)
	for _, i := range customs {
		io.WriteString(h, i)
	}

	natives := []string{}
	for i := range tci.ConstraintNative {
		j := i
		natives = append(natives, j)
	}
	sort.Strings(natives)
	for _, i := range natives {
		io.WriteString(h, i)
	}

	attributes := []string{}
	for i := range tci.ConstraintAttribute {
		j := i
		attributes = append(attributes, j)
	}
	sort.Strings(attributes)
	for _, i := range attributes {
		svcattr := []string{}
		for j := range tci.ConstraintAttribute[i] {
			k := j
			svcattr = append(svcattr, k)
		}
		sort.Strings(svcattr)
		io.WriteString(h, i)
		for _, l := range svcattr {
			io.WriteString(h, l)
		}
	}
	tci.oldConstraintHash = base64.URLEncoding.EncodeToString(h.Sum(nil))
	io.WriteString(h, tci.ConfigID.String())
	io.WriteString(h, tci.CheckID.String())
	tci.ConstraintHash = base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (tci *CheckInstance) calcConstraintValHash() {
	h := sha512.New()
	io.WriteString(h, tci.ConstraintOncall)

	services := []string{}
	for i := range tci.ConstraintService {
		j := i
		services = append(services, j)
	}
	sort.Strings(services)
	for _, i := range services {
		io.WriteString(h, i)
		io.WriteString(h, tci.ConstraintService[i])
	}

	systems := []string{}
	for i := range tci.ConstraintSystem {
		j := i
		systems = append(systems, j)
	}
	sort.Strings(systems)
	for _, i := range systems {
		io.WriteString(h, i)
		io.WriteString(h, tci.ConstraintSystem[i])
	}

	customs := []string{}
	for i := range tci.ConstraintCustom {
		j := i
		customs = append(customs, j)
	}
	sort.Strings(customs)
	for _, i := range customs {
		io.WriteString(h, i)
		io.WriteString(h, tci.ConstraintCustom[i])
	}

	natives := []string{}
	for i := range tci.ConstraintNative {
		j := i
		natives = append(natives, j)
	}
	sort.Strings(natives)
	for _, i := range natives {
		io.WriteString(h, i)
		io.WriteString(h, tci.ConstraintNative[i])
	}

	attributes := []string{}
	for i := range tci.ConstraintAttribute {
		j := i
		attributes = append(attributes, j)
	}
	sort.Strings(attributes)
	for _, i := range attributes {
		svcattr := []string{}
		for j := range tci.ConstraintAttribute[i] {
			k := j
			svcattr = append(svcattr, k)
		}
		sort.Strings(svcattr)
		io.WriteString(h, i)
		for _, l := range svcattr {
			io.WriteString(h, l)
			vals := make([]string, len(tci.ConstraintAttribute[i][l]))
			copy(vals, tci.ConstraintAttribute[i][l])
			sort.Strings(vals)
			for _, m := range vals {
				io.WriteString(h, m)
			}
		}
	}
	tci.oldConstraintValHash = base64.URLEncoding.EncodeToString(h.Sum(nil))
	io.WriteString(h, tci.ConfigID.String())
	io.WriteString(h, tci.CheckID.String())
	io.WriteString(h, tci.InstanceService)
	tci.ConstraintValHash = base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (tci *CheckInstance) calcInstanceSvcCfgHash() {
	h := sha512.New()

	attributes := []string{}
	for i := range tci.InstanceServiceConfig {
		j := i
		attributes = append(attributes, j)
	}
	sort.Strings(attributes)
	for _, i := range attributes {
		io.WriteString(h, i)
		io.WriteString(h, tci.InstanceServiceConfig[i])
	}
	tci.oldInstanceSvcCfgHash = base64.URLEncoding.EncodeToString(h.Sum(nil))
	io.WriteString(h, tci.ConfigID.String())
	io.WriteString(h, tci.CheckID.String())
	io.WriteString(h, tci.InstanceService)
	tci.InstanceSvcCfgHash = base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (tci CheckInstance) MakeAction() Action {
	serviceCfg, err := json.Marshal(tci.InstanceServiceConfig)
	if err != nil {
		serviceCfg = []byte{}
	}

	return Action{
		CheckInstance: proto.CheckInstance{
			InstanceID:            tci.InstanceID.String(),
			CheckID:               tci.CheckID.String(),
			ConfigID:              tci.ConfigID.String(),
			InstanceConfigID:      tci.InstanceConfigID.String(),
			Version:               tci.Version,
			ConstraintHash:        tci.ConstraintHash,
			ConstraintValHash:     tci.ConstraintValHash,
			InstanceSvcCfgHash:    tci.InstanceSvcCfgHash,
			InstanceService:       tci.InstanceService,
			InstanceServiceConfig: string(serviceCfg),
		},
	}
}

func (tci *CheckInstance) MatchConstraints(target *CheckInstance) bool {
	if tci.matchConstraintHash(target) && tci.matchConstraintValHash(target) {
		return true
	}
	if tci.matchOldInstance(target) {
		return true
	}
	return false
}

func (tci *CheckInstance) MatchServiceConstraints(target *CheckInstance) bool {
	if tci.matchInstanceSvcCfgHash(target) && tci.MatchConstraints(target) {
		return true
	}
	if tci.matchOldService(target) {
		return true
	}
	return false
}

func (tci *CheckInstance) matchConstraintHash(target *CheckInstance) bool {
	if tci.ConstraintHash == target.ConstraintHash {
		return true
	}
	return false
}

func (tci *CheckInstance) matchConstraintValHash(target *CheckInstance) bool {
	if tci.ConstraintValHash == target.ConstraintValHash {
		return true
	}
	return false
}

func (tci *CheckInstance) matchOldInstance(target *CheckInstance) bool {
	if tci.oldConstraintHash == target.oldConstraintHash &&
		uuid.Equal(tci.ConfigID, target.ConfigID) &&
		tci.oldConstraintValHash == target.oldConstraintValHash {
		return true
	}
	return false
}

func (tci *CheckInstance) matchInstanceSvcCfgHash(target *CheckInstance) bool {
	if tci.InstanceSvcCfgHash == target.InstanceSvcCfgHash {
		return true
	}
	return false
}

func (tci *CheckInstance) matchOldService(target *CheckInstance) bool {
	if tci.oldInstanceSvcCfgHash == target.oldInstanceSvcCfgHash &&
		tci.oldConstraintHash == target.oldConstraintHash &&
		tci.oldConstraintValHash == target.oldConstraintValHash &&
		uuid.Equal(tci.ConfigID, target.ConfigID) {
		return true
	}
	return false
}

type checkContext struct {
	uuid                   string
	brokeConstraint        bool
	hasServiceConstraint   bool
	hasAttributeConstraint bool
	view                   string
	attributes             []CheckConstraint
	oncallConstr           string
	systemConstr           map[string]string              // ID -> Value
	nativeConstr           map[string]string              // Property -> Value
	serviceConstr          map[string]string              // ID -> Value
	customConstr           map[string]string              // ID -> Value
	attributeConstr        map[string]map[string][]string // svcID -> attr -> [ value, ... ]
	newCheckInstances      []string
	newInstances           map[string]CheckInstance
}

func newCheckContext(uuid, view string) *checkContext {
	cc := checkContext{
		uuid: uuid,
		view: view,
	}
	cc.attributes = []CheckConstraint{}
	cc.systemConstr = make(map[string]string)
	cc.nativeConstr = make(map[string]string)
	cc.serviceConstr = make(map[string]string)
	cc.customConstr = make(map[string]string)
	cc.attributeConstr = make(map[string]map[string][]string)
	cc.newCheckInstances = []string{}
	cc.newInstances = make(map[string]CheckInstance)
	return &cc
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
