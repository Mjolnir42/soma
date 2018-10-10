/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <code.jpe@gmail.com>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// Provider is the definition a statistics provider system
type Provider struct {
	Name    string           `json:"name,omitempty"`
	Details *ProviderDetails `json:"details,omitempty"`
}

// Clone returns a copy of p
func (p *Provider) Clone() Provider {
	clone := Provider{
		Name: p.Name,
	}
	if p.Details != nil {
		clone.Details = p.Details.Clone()
	}
	return clone
}

// ProviderDetails contains metadata about a provider
type ProviderDetails struct {
	Creation *DetailsCreation `json:"creation,omitempty"`
}

// Clone returns a copy of p
func (p *ProviderDetails) Clone() *ProviderDetails {
	clone := &ProviderDetails{}
	if p.Creation != nil {
		clone.Creation = p.Creation.Clone()
	}
	return clone
}

// ProviderFilter represents parts of a Provider that can be searched
type ProviderFilter struct {
	Name string `json:"name,omitempty"`
}

// NewProviderRequest returns a new Request with fields preallocated
// for filling in Provider data, ensuring no nilptr-deref takes place.
func NewProviderRequest() Request {
	return Request{
		Flags:    &Flags{},
		Provider: &Provider{},
	}
}

// NewProviderFilter returns a new Request with fields preallocated
// for filling in a Provider filter, ensuring no nilptr-deref takes place.
func NewProviderFilter() Request {
	return Request{
		Filter: &Filter{
			Provider: &ProviderFilter{},
		},
	}
}

// NewProviderResult returns a new Result with fields preallocated
// for filling in Provider data, ensuring no nilptr-deref takes place.
func NewProviderResult() Result {
	return Result{
		Errors:    &[]string{},
		Providers: &[]Provider{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
