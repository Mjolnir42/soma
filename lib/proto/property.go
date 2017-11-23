/*-
 * Copyright (c) 2015-2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Property struct {
	Type             string           `json:"type"`
	RepositoryID     string           `json:"repositoryID,omitempty"`
	BucketID         string           `json:"bucketID,omitempty"`
	InstanceID       string           `json:"instanceID,omitempty"`
	View             string           `json:"view,omitempty"`
	Inheritance      bool             `json:"inheritance,omitempty"`
	ChildrenOnly     bool             `json:"childrenOnly,omitempty"`
	IsInherited      bool             `json:"isInherited,omitempty"`
	SourceInstanceID string           `json:"sourceInstanceID,omitempty"`
	SourceType       string           `json:"sourceType,omitempty"`
	InheritedFrom    string           `json:"inheritedFrom,omitempty"`
	Custom           *PropertyCustom  `json:"custom,omitempty"`
	System           *PropertySystem  `json:"system,omitempty"`
	Service          *PropertyService `json:"service,omitempty"`
	Native           *PropertyNative  `json:"native,omitempty"`
	Oncall           *PropertyOncall  `json:"oncall,omitempty"`
	Details          *PropertyDetails `json:"details,omitempty"`
}

func (p *Property) Clone() Property {
	clone := Property{
		Type:             p.Type,
		RepositoryID:     p.RepositoryID,
		BucketID:         p.BucketID,
		InstanceID:       p.InstanceID,
		View:             p.View,
		Inheritance:      p.Inheritance,
		ChildrenOnly:     p.ChildrenOnly,
		IsInherited:      p.IsInherited,
		SourceInstanceID: p.SourceInstanceID,
		SourceType:       p.SourceType,
		InheritedFrom:    p.InheritedFrom,
	}
	if p.Custom != nil {
		clone.Custom = p.Custom.Clone()
	}
	if p.System != nil {
		clone.System = p.System.Clone()
	}
	if p.Service != nil {
		clone.Service = p.Service.Clone()
	}
	if p.Native != nil {
		clone.Native = p.Native.Clone()
	}
	if p.Oncall != nil {
		clone.Oncall = p.Oncall.Clone()
	}
	if p.Details != nil {
		clone.Details = p.Details.Clone()
	}
	return clone
}

type PropertyFilter struct {
	Name         string `json:"name,omitempty"`
	Type         string `json:"type,omitempty"`
	RepositoryID string `json:"repositoryID,omitempty"`
}

type PropertyDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
}

func (t *PropertyDetails) Clone() *PropertyDetails {
	clone := &PropertyDetails{}
	if t.Creation != nil {
		clone.Creation = t.Creation.Clone()
	}
	return clone
}

type PropertyCustom struct {
	ID           string `json:"ID,omitempty"`
	Name         string `json:"name,omitempty"`
	RepositoryID string `json:"repositoryID,omitempty"`
	Value        string `json:"value,omitempty"`
}

func (t *PropertyCustom) Clone() *PropertyCustom {
	return &PropertyCustom{
		ID:           t.ID,
		Name:         t.Name,
		RepositoryID: t.RepositoryID,
		Value:        t.Value,
	}
}

func (t *PropertyCustom) DeepCompare(a *PropertyCustom) bool {
	if t.ID != a.ID || t.Name != a.Name || t.RepositoryID != a.RepositoryID || t.Value != a.Value {
		return false
	}
	return true
}

func (t *PropertyCustom) DeepCompareSlice(a *[]PropertyCustom) bool {
	if a == nil || *a == nil {
		return false
	}
	for _, cust := range *a {
		if t.DeepCompare(&cust) {
			return true
		}
	}
	return false
}

type PropertySystem struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

func (t *PropertySystem) Clone() *PropertySystem {
	return &PropertySystem{
		Name:  t.Name,
		Value: t.Value,
	}
}

func (t *PropertySystem) DeepCompare(a *PropertySystem) bool {
	if t.Name != a.Name || t.Value != a.Value {
		return false
	}
	return true
}

func (t *PropertySystem) DeepCompareSlice(a *[]PropertySystem) bool {
	if a == nil || *a == nil {
		return false
	}
	for _, sys := range *a {
		if t.DeepCompare(&sys) {
			return true
		}
	}
	return false
}

type PropertyService struct {
	Name       string             `json:"name,omitempty"`
	TeamID     string             `json:"teamID,omitempty"`
	Attributes []ServiceAttribute `json:"attributes"`
}

func (t *PropertyService) Clone() *PropertyService {
	clone := &PropertyService{
		Name:   t.Name,
		TeamID: t.TeamID,
	}
	clone.Attributes = make([]ServiceAttribute, len(t.Attributes))
	for i := range t.Attributes {
		clone.Attributes[i] = t.Attributes[i].Clone()
	}
	return clone
}

func (t *PropertyService) DeepCompare(a *PropertyService) bool {
	if t.Name != a.Name || t.TeamID != a.TeamID {
		return false
	}
attrloop:
	for _, attr := range t.Attributes {
		if attr.DeepCompareSlice(&a.Attributes) {
			continue attrloop
		}
		return false
	}
revattrloop:
	for _, attr := range a.Attributes {
		if attr.DeepCompareSlice(&t.Attributes) {
			continue revattrloop
		}
		return false
	}
	return true
}

type PropertyNative struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

func (t *PropertyNative) Clone() *PropertyNative {
	return &PropertyNative{
		Name:  t.Name,
		Value: t.Value,
	}
}

func (t *PropertyNative) DeepCompare(a *PropertyNative) bool {
	if t.Name != a.Name || t.Value != a.Value {
		return false
	}
	return true
}

type PropertyOncall struct {
	ID     string `json:"ID,omitempty"`
	Name   string `json:"name,omitempty"`
	Number string `json:"number,omitempty"`
}

func (t *PropertyOncall) Clone() *PropertyOncall {
	return &PropertyOncall{
		ID:     t.ID,
		Name:   t.Name,
		Number: t.Number,
	}
}

func (t *PropertyOncall) DeepCompare(a *PropertyOncall) bool {
	if t.ID != a.ID || t.Name != a.Name || t.Number != a.Number {
		return false
	}
	return true
}

type ServiceAttribute struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

func (t *ServiceAttribute) Clone() ServiceAttribute {
	return ServiceAttribute{
		Name:  t.Name,
		Value: t.Value,
	}
}

func (t *ServiceAttribute) DeepCompare(a *ServiceAttribute) bool {
	if t.Name != a.Name || t.Value != a.Value {
		return false
	}
	return true
}

func (t *ServiceAttribute) DeepCompareSlice(a *[]ServiceAttribute) bool {
	if a == nil || *a == nil {
		return false
	}
	for _, attr := range *a {
		if t.DeepCompare(&attr) {
			return true
		}
	}
	return false
}

func NewPropertyRequest() Request {
	return Request{
		Flags:    &Flags{},
		Property: &Property{},
	}
}

func NewPropertyFilter() Request {
	return Request{
		Filter: &Filter{
			Property: &PropertyFilter{},
		},
	}
}

func NewPropertyResult() Result {
	return Result{
		Errors:     &[]string{},
		Properties: &[]Property{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
