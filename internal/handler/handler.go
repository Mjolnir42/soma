/*-
 * Copyright (c) 2017-2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package handler // import "github.com/mjolnir42/soma/internal/handler"

import (
	"database/sql"

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/msg"
)

// Handler process a specific request type
type Handler interface {
	Register(*sql.DB, ...*logrus.Logger)
	Run()
	ShutdownNow()
	Intake() chan msg.Request
	PriorityIntake() chan msg.Request
	RegisterRequests(*Map)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
