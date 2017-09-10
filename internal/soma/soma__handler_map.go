/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"database/sql"
	"sync"

	"github.com/Sirupsen/logrus"
)

// HandlerMap is a concurrent map that is used to look up input
// channels for application handlers
type HandlerMap struct {
	hmap map[string]Handler
	sync.RWMutex
}

// Add registers a new handler
func (h *HandlerMap) Add(key string, value Handler) {
	h.Lock()
	defer h.Unlock()
	h.hmap[key] = value
}

// Get retrieves a handler
func (h *HandlerMap) Get(key string) Handler {
	h.RLock()
	defer h.RUnlock()
	return h.hmap[key]
}

// Exists checks if a handler exists. This function is only safe to
// call if it is certain that the calling function is the only one
// that adds or removes the searched handler
func (h *HandlerMap) Exists(key string) bool {
	h.RLock()
	defer h.RUnlock()
	if _, ok := h.hmap[key]; ok {
		return true
	}
	return false
}

// Del removes a handler
func (h *HandlerMap) Del(key string) {
	h.Lock()
	defer h.Unlock()
	delete(h.hmap, key)
}

// Range returns all handlers
func (h *HandlerMap) Range() map[string]Handler {
	h.RLock()
	defer h.RUnlock()
	return h.hmap
}

// Register calls register() for each handler
func (h *HandlerMap) Register(n string, c *sql.DB, l []*logrus.Logger) {
	h.Lock()
	defer h.Unlock()
	h.hmap[n].Register(c, l...)
}

// Run starts the handler n
func (h *HandlerMap) Run(n string) {
	h.Lock()
	defer h.Unlock()
	go h.hmap[n].Run()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
