/*-
 * Copyright (c) 2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import uuid "github.com/satori/go.uuid"

func generateHandlerName() string {
	return uuid.Must(uuid.NewV4()).String()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
