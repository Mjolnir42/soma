package main

import (
	"encoding/hex"

	log "github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/msg"
)

func startHandlers(appLog, reqLog, errLog *log.Logger) {
	spawnSupervisorHandler(appLog, reqLog, errLog)

	spawnOutputTreeHandler(appLog, reqLog, errLog)
}

func spawnSupervisorHandler(appLog, reqLog, errLog *log.Logger) {
	var supervisorHandler supervisor
	var err error
	supervisorHandler.input = make(chan msg.Request, 1024)
	supervisorHandler.update = make(chan msg.Request, 1024)
	supervisorHandler.shutdown = make(chan bool)
	supervisorHandler.conn = conn
	supervisorHandler.appLog = appLog
	supervisorHandler.reqLog = reqLog
	supervisorHandler.errLog = errLog
	supervisorHandler.readonly = SomaCfg.ReadOnly
	if supervisorHandler.seed, err = hex.DecodeString(SomaCfg.Auth.TokenSeed); err != nil {
		panic(err)
	}
	if len(supervisorHandler.seed) == 0 {
		panic(`token.seed has length 0`)
	}
	if supervisorHandler.key, err = hex.DecodeString(SomaCfg.Auth.TokenKey); err != nil {
		panic(err)
	}
	if len(supervisorHandler.key) == 0 {
		panic(`token.key has length 0`)
	}
	supervisorHandler.tokenExpiry = SomaCfg.Auth.TokenExpirySeconds
	supervisorHandler.kexExpiry = SomaCfg.Auth.KexExpirySeconds
	supervisorHandler.credExpiry = SomaCfg.Auth.CredentialExpiryDays
	supervisorHandler.activation = SomaCfg.Auth.Activation
	handlerMap[`supervisor`] = &supervisorHandler
	go supervisorHandler.run()
}

func spawnOutputTreeHandler(appLog, reqLog, errLog *log.Logger) {
	var handler outputTree
	handler.input = make(chan msg.Request, 128)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`tree_r`] = &handler
	go handler.run()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
