/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type CheckConfig struct {
	ID           string                  `json:"ID,omitempty"`
	Name         string                  `json:"name,omitempty"`
	Interval     uint64                  `json:"interval,omitempty"`
	RepositoryID string                  `json:"repositoryID,omitempty"`
	BucketID     string                  `json:"bucketID,omitempty"`
	CapabilityID string                  `json:"capabilityID,omitempty"`
	ObjectID     string                  `json:"objectID,omitempty"`
	ObjectType   string                  `json:"objectType,omitempty"`
	IsActive     bool                    `json:"isActive,omitempty"`
	IsEnabled    bool                    `json:"isEnabled,omitempty"`
	Inheritance  bool                    `json:"inheritance,omitempty"`
	ChildrenOnly bool                    `json:"childrenOnly,omitempty"`
	ExternalID   string                  `json:"externalID,omitempty"`
	Constraints  []CheckConfigConstraint `json:"constraints,omitempty"`
	Thresholds   []CheckConfigThreshold  `json:"thresholds,omitempty"`
	Details      *CheckConfigDetails     `json:"details,omitempty"`
}

func (c *CheckConfig) Clone() CheckConfig {
	clone := CheckConfig{
		ID:           c.ID,
		Name:         c.Name,
		Interval:     c.Interval,
		RepositoryID: c.RepositoryID,
		BucketID:     c.BucketID,
		CapabilityID: c.CapabilityID,
		ObjectID:     c.ObjectID,
		ObjectType:   c.ObjectType,
		IsActive:     c.IsActive,
		IsEnabled:    c.IsEnabled,
		Inheritance:  c.Inheritance,
		ChildrenOnly: c.ChildrenOnly,
		ExternalID:   c.ExternalID,
	}
	clone.Constraints = make([]CheckConfigConstraint, len(c.Constraints))
	for i := range c.Constraints {
		clone.Constraints[i] = c.Constraints[i].Clone()
	}
	clone.Thresholds = make([]CheckConfigThreshold, len(c.Thresholds))
	for i := range c.Thresholds {
		clone.Thresholds[i] = c.Thresholds[i].Clone()
	}
	if c.Details != nil {
		clone.Details = c.Details.Clone()
	}
	return clone
}

type CheckConfigConstraint struct {
	ConstraintType string            `json:"constraintType,omitempty"`
	Native         *PropertyNative   `json:"native,omitempty"`
	Oncall         *PropertyOncall   `json:"oncall,omitempty"`
	Custom         *PropertyCustom   `json:"custom,omitempty"`
	System         *PropertySystem   `json:"system,omitempty"`
	Service        *PropertyService  `json:"service,omitempty"`
	Attribute      *ServiceAttribute `json:"attribute,omitempty"`
}

func (c *CheckConfigConstraint) Clone() CheckConfigConstraint {
	clone := CheckConfigConstraint{
		ConstraintType: c.ConstraintType,
	}
	if c.Native != nil {
		clone.Native = c.Native.Clone()
	}
	if c.Oncall != nil {
		clone.Oncall = c.Oncall.Clone()
	}
	if c.Custom != nil {
		clone.Custom = c.Custom.Clone()
	}
	if c.System != nil {
		clone.System = c.System.Clone()
	}
	if c.Service != nil {
		clone.Service = c.Service.Clone()
	}
	if c.Attribute != nil {
		ac := c.Attribute.Clone()
		clone.Attribute = &ac
	}
	return clone
}

func (c *CheckConfigConstraint) DeepCompare(a *CheckConfigConstraint) bool {
	if c.ConstraintType != a.ConstraintType {
		return false
	}
	switch c.ConstraintType {
	case "native":
		if c.Native.DeepCompare(a.Native) {
			return true
		}
	case "oncall":
		if c.Oncall.DeepCompare(a.Oncall) {
			return true
		}
	case "custom":
		if c.Custom.DeepCompare(a.Custom) {
			return true
		}
	case "system":
		if c.System.DeepCompare(a.System) {
			return true
		}
	case "service":
		if c.Service.DeepCompare(a.Service) {
			return true
		}
	case "attribute":
		if c.Attribute.DeepCompare(a.Attribute) {
			return true
		}
	}
	return false
}

func (c *CheckConfigConstraint) DeepCompareSlice(a []CheckConfigConstraint) bool {
	if a == nil {
		return false
	}
	for _, constr := range a {
		if c.DeepCompare(&constr) {
			return true
		}
	}
	return false
}

type CheckConfigThreshold struct {
	Predicate Predicate
	Level     Level
	Value     int64
}

func (c *CheckConfigThreshold) Clone() CheckConfigThreshold {
	return CheckConfigThreshold{
		Predicate: c.Predicate,
		Level:     c.Level,
		Value:     c.Value,
	}
}

func (c *CheckConfigThreshold) DeepCompareSlice(a []CheckConfigThreshold) bool {
	if a == nil {
		return false
	}
	for _, thr := range a {
		if c.DeepCompare(&thr) {
			return true
		}
	}
	return false
}

func (c *CheckConfigThreshold) DeepCompare(a *CheckConfigThreshold) bool {
	if c.Value != a.Value || c.Level.Name != a.Level.Name ||
		c.Predicate.Symbol != a.Predicate.Symbol {
		return false
	}
	return true
}

type CheckConfigDetails struct {
	Creation  *DetailsCreation    `json:"creation,omitempty"`
	Instances []CheckInstanceInfo `json:"instances,omitempty"`
}

func (c *CheckConfigDetails) Clone() *CheckConfigDetails {
	clone := &CheckConfigDetails{}
	if c.Creation != nil {
		clone.Creation = c.Creation.Clone()
	}
	clone.Instances = make([]CheckInstanceInfo, len(c.Instances))
	for i := range c.Instances {
		clone.Instances[i] = c.Instances[i].Clone()
	}
	return clone
}

type CheckConfigFilter struct {
	ID           string `json:"ID,omitempty"`
	Name         string `json:"name,omitempty"`
	CapabilityID string `json:"capabilityID,omitempty"`
}

type CheckInstanceInfo struct {
	ID            string `json:"ID,omitempty"`
	ObjectID      string `json:"objectID,omitempty"`
	ObjectType    string `json:"objectType,omitempty"`
	CurrentStatus string `json:"currentStatus,omitempty"`
	NextStatus    string `json:"nextStatus,omitempty"`
}

func (c *CheckInstanceInfo) Clone() CheckInstanceInfo {
	return CheckInstanceInfo{
		ID:            c.ID,
		ObjectID:      c.ObjectID,
		ObjectType:    c.ObjectType,
		CurrentStatus: c.CurrentStatus,
		NextStatus:    c.NextStatus,
	}
}

func (c *CheckConfig) DeepCompare(a *CheckConfig) bool {
	if a == nil {
		return false
	}
	if c.ID != a.ID || c.Name != a.Name || c.Interval != a.Interval ||
		c.RepositoryID != a.RepositoryID || c.BucketID != a.BucketID ||
		c.CapabilityID != a.CapabilityID || c.ObjectID != a.ObjectID ||
		c.ObjectType != a.ObjectType || c.IsActive != a.IsActive ||
		c.IsEnabled != a.IsEnabled || c.Inheritance != a.Inheritance ||
		c.ChildrenOnly != a.ChildrenOnly || c.ExternalID != a.ExternalID {
		return false
	}
threshloop:
	for _, thr := range c.Thresholds {
		if thr.DeepCompareSlice(a.Thresholds) {
			continue threshloop
		}
		return false
	}
revthreshloop:
	for _, thr := range a.Thresholds {
		if thr.DeepCompareSlice(c.Thresholds) {
			continue revthreshloop
		}
		return false
	}
constrloop:
	for _, constr := range c.Constraints {
		if constr.DeepCompareSlice(a.Constraints) {
			continue constrloop
		}
		return false
	}
revconstrloop:
	for _, constr := range a.Constraints {
		if constr.DeepCompareSlice(c.Constraints) {
			continue revconstrloop
		}
		return false
	}
	return true
}

func NewCheckConfigRequest() Request {
	return Request{
		Flags:       &Flags{},
		CheckConfig: &CheckConfig{},
	}
}

func NewCheckConfigFilter() Request {
	return Request{
		Filter: &Filter{
			CheckConfig: &CheckConfigFilter{},
		},
	}
}

func NewCheckConfigResult() Result {
	return Result{
		Errors:       &[]string{},
		CheckConfigs: &[]CheckConfig{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
