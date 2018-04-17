/*-
 * Copyright (c) 2016,2018 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <code.jpe@gmail.com>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/soma/lib/proto"

// PushNotification is used to signal monitoring systems that an update
// checkinstance with checkinstanceID UUID is available
type PushNotification struct {
	UUID string `json:"uuid" valid:"uuidv4"`
	Path string `json:"path" valid:"abspath"`
}

// NewPushNotification returns a new push notification
func NewPushNotification() PushNotification {
	return PushNotification{}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
