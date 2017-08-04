/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"sync"

	"github.com/client9/reopen"
)

// LogHandleMap is a concurrent map that is used to look up
// filehandles of active logfiles
type LogHandleMap struct {
	hmap map[string]*reopen.FileWriter
	sync.RWMutex
}

// Add registers a new filehandle
func (l *LogHandleMap) Add(key string, fh *reopen.FileWriter) {
	l.Lock()
	defer l.Unlock()
	l.hmap[key] = fh
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
