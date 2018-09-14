/*-
 * Copyright (c) 2017-2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package handler // import "github.com/mjolnir42/soma/internal/handler"

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/msg"
)

// Map is a concurrent map that is used to look up input
// channels for application handlers
type Map struct {
	hmap        map[string]Handler
	requestPath map[string]map[string]Handler
	sync.RWMutex
}

// NewMap returns a new HandlerMap
func NewMap() *Map {
	h := Map{}
	h.hmap = make(map[string]Handler)
	h.requestPath = make(map[string]map[string]Handler)
	return &h
}

// Add registers a new handler
func (h *Map) Add(key string, value Handler) {
	h.Lock()
	defer h.Unlock()
	h.hmap[key] = value
}

// Get retrieves a handler by name
func (h *Map) Get(key string) Handler {
	h.RLock()
	defer h.RUnlock()
	return h.hmap[key]
}

// MustLookup retrieves a handler by request and panics if no handler
// can be found
func (h *Map) MustLookup(q *msg.Request) Handler {
	h.RLock()
	defer h.RUnlock()
	if _, ok := h.requestPath[q.Section]; !ok {
		panic(`Section not found in requestPath map`)
	}
	if _, ok := h.requestPath[q.Section][q.Action]; !ok {
		panic(`Action not found in requestPath map`)
	}
	return h.requestPath[q.Section][q.Action]
}

// Exists checks if a handler exists. This function is only safe to
// call if it is certain that the calling function is the only one
// that adds or removes the searched handler
func (h *Map) Exists(key string) bool {
	h.RLock()
	defer h.RUnlock()
	if _, ok := h.hmap[key]; ok {
		return true
	}
	return false
}

// Del removes a handler
func (h *Map) Del(key string) {
	h.Lock()
	defer h.Unlock()
	delete(h.hmap, key)
}

// Range returns all handlers
func (h *Map) Range() map[string]Handler {
	h.RLock()
	defer h.RUnlock()
	return h.hmap
}

// Register calls register() for each handler
func (h *Map) Register(n string, c *sql.DB, l []*logrus.Logger) {
	h.Lock()
	defer h.Unlock()
	h.hmap[n].Register(c, l...)
	h.hmap[n].RegisterRequests(h)
}

// Request registers a request to a handler registered as name
func (h *Map) Request(section, action, name string) {
	if _, ok := h.requestPath[section]; !ok {
		h.requestPath[section] = make(map[string]Handler)
	}
	if _, ok := h.hmap[name]; !ok {
		panic(fmt.Sprintf("RequestPath attmpted register for unknown handler, %s::%s by %s",
			section, action, name))
	}
	h.requestPath[section][action] = h.hmap[name]
}

// Run starts the handler n
func (h *Map) Run(n string) {
	h.Lock()
	defer h.Unlock()
	go h.hmap[n].Run()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
