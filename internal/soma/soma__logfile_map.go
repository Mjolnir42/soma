/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/client9/reopen"
)

// LogHandleMap is a concurrent map that is used to look up
// filehandles of active logfiles
type LogHandleMap struct {
	hmap map[string]*reopen.FileWriter
	lmap map[string]*logrus.Logger
	sync.RWMutex
}

// New returns an initialized LogHandleMap
func NewLogHandleMap() *LogHandleMap {
	lm := &LogHandleMap{}
	lm.hmap = make(map[string]*reopen.FileWriter)
	lm.lmap = make(map[string]*logrus.Logger)
	return lm
}

// Add registers a new filehandle
func (l *LogHandleMap) Add(key string, fh *reopen.FileWriter, lg *logrus.Logger) {
	l.Lock()
	defer l.Unlock()
	l.hmap[key] = fh
	l.lmap[key] = lg
}

// Get retrieves a filehandle
func (l *LogHandleMap) Get(key string) *reopen.FileWriter {
	l.RLock()
	defer l.RUnlock()
	return l.hmap[key]
}

// Del removes a filehandle
func (l *LogHandleMap) Del(key string) {
	l.Lock()
	defer l.Unlock()
	delete(l.hmap, key)
}

// GetLogger retrieves a Logger without locking. This should be used
// inside an active Range() lock.
func (l *LogHandleMap) GetLogger(key string) *logrus.Logger {
	return l.lmap[key]
}

// Range locks l and returns the embedded map. Unlocking must
// be done by the caller via RangeUnlock()
func (l *LogHandleMap) Range() map[string]*reopen.FileWriter {
	l.Lock()
	return l.hmap
}

// RangeUnlock unlocks l. It is required to be called after Range() once
// the caller is finished with the map.
func (l *LogHandleMap) RangeUnlock() {
	l.Unlock()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
