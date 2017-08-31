/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

// Package soma implements the application handlers of the SOMA
// service.
package soma

import (
	"database/sql"

	"github.com/Sirupsen/logrus"
)

// Soma application struct
type Soma struct {
	handlerMap   *HandlerMap
	logMap       *LogHandleMap
	dbConnection *sql.DB
	conf         *Config
	appLog       *logrus.Logger
	reqLog       *logrus.Logger
	errLog       *logrus.Logger
}

// New returns a new SOMA application
func New(
	appHandlerMap *HandlerMap,
	logHandleMap *LogHandleMap,
	dbConnection *sql.DB,
	conf *Config,
	appLog, reqLog, errLog *logrus.Logger,
) *Soma {
	s := Soma{}
	s.handlerMap = appHandlerMap
	s.logMap = logHandleMap
	s.dbConnection = dbConnection
	s.conf = conf
	s.appLog = appLog
	s.reqLog = reqLog
	s.errLog = errLog
	return &s
}

// exportLogger returns references to the instances loggers
func (s *Soma) exportLogger() []*logrus.Logger {
	return []*logrus.Logger{s.appLog, s.reqLog, s.errLog}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
