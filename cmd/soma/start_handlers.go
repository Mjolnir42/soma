package main

import (
	"encoding/hex"

	log "github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/msg"
)

func startHandlers(appLog, reqLog, errLog *log.Logger) {
	spawnSupervisorHandler(appLog, reqLog, errLog)

	spawnAttributeRead(appLog, reqLog, errLog)
	spawnEntityRead(appLog, reqLog, errLog)
	spawnEnvironmentReadHandler(appLog, reqLog, errLog)
	spawnGroupReadHandler(appLog, reqLog, errLog)
	spawnInstanceReadHandler(appLog, reqLog, errLog)
	spawnJobReadHandler(appLog, reqLog, errLog)
	spawnMonitoringRead(appLog, reqLog, errLog)
	spawnObjectStateReadHandler(appLog, reqLog, errLog)
	spawnOutputTreeHandler(appLog, reqLog, errLog)
	spawnWorkflowReadHandler(appLog, reqLog, errLog)

	if !SomaCfg.ReadOnly {
		if !SomaCfg.Observer {
			spawnAttributeWrite(appLog, reqLog, errLog)
			spawnEntityWrite(appLog, reqLog, errLog)
			spawnEnvironmentWriteHandler(appLog, reqLog, errLog)
			spawnMonitoringWrite(appLog, reqLog, errLog)
			spawnObjectStateWriteHandler(appLog, reqLog, errLog)
			spawnWorkflowWriteHandler(appLog, reqLog, errLog)
		}
	}
}

func spawnEnvironmentReadHandler(appLog, reqLog, errLog *log.Logger) {
	var environmentReadHandler environmentRead
	environmentReadHandler.input = make(chan msg.Request)
	environmentReadHandler.shutdown = make(chan bool)
	environmentReadHandler.conn = conn
	environmentReadHandler.appLog = appLog
	environmentReadHandler.reqLog = reqLog
	environmentReadHandler.errLog = errLog
	handlerMap[`environment_r`] = &environmentReadHandler
	go environmentReadHandler.run()
}

func spawnEnvironmentWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var environmentWriteHandler environmentWrite
	environmentWriteHandler.input = make(chan msg.Request, 64)
	environmentWriteHandler.shutdown = make(chan bool)
	environmentWriteHandler.conn = conn
	environmentWriteHandler.appLog = appLog
	environmentWriteHandler.reqLog = reqLog
	environmentWriteHandler.errLog = errLog
	handlerMap[`environment_w`] = &environmentWriteHandler
	go environmentWriteHandler.run()
}

func spawnObjectStateReadHandler(appLog, reqLog, errLog *log.Logger) {
	var handler stateRead
	handler.input = make(chan msg.Request)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`state_r`] = &handler
	go handler.run()
}

func spawnObjectStateWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var handler stateWrite
	handler.input = make(chan msg.Request, 64)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`state_w`] = &handler
	go handler.run()
}

func spawnEntityRead(appLog, reqLog, errLog *log.Logger) {
	var handler entityRead
	handler.input = make(chan msg.Request)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`entity_r`] = &handler
	go handler.run()
}

func spawnEntityWrite(appLog, reqLog, errLog *log.Logger) {
	var handler entityWrite
	handler.input = make(chan msg.Request, 64)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`entity_w`] = &handler
	go handler.run()
}

func spawnMonitoringRead(appLog, reqLog, errLog *log.Logger) {
	var handler monitoringRead
	handler.input = make(chan msg.Request, 64)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`monitoring_r`] = &handler
	go handler.run()
}

func spawnMonitoringWrite(appLog, reqLog, errLog *log.Logger) {
	var handler monitoringWrite
	handler.input = make(chan msg.Request, 64)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`monitoring_w`] = &handler
	go handler.run()
}

func spawnAttributeRead(appLog, reqLog, errLog *log.Logger) {
	var handler attributeRead
	handler.input = make(chan msg.Request, 64)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`attribute_r`] = &handler
	go handler.run()
}

func spawnAttributeWrite(appLog, reqLog, errLog *log.Logger) {
	var handler attributeWrite
	handler.input = make(chan msg.Request, 64)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`attribute_w`] = &handler
	go handler.run()
}

func spawnGroupReadHandler(appLog, reqLog, errLog *log.Logger) {
	var groupReadHandler somaGroupReadHandler
	groupReadHandler.input = make(chan somaGroupRequest, 64)
	groupReadHandler.shutdown = make(chan bool)
	groupReadHandler.conn = conn
	groupReadHandler.appLog = appLog
	groupReadHandler.reqLog = reqLog
	groupReadHandler.errLog = errLog
	handlerMap["groupReadHandler"] = &groupReadHandler
	go groupReadHandler.run()
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

func spawnJobReadHandler(appLog, reqLog, errLog *log.Logger) {
	var handler jobsRead
	handler.input = make(chan msg.Request, 256)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`jobs_r`] = &handler
	go handler.run()
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

func spawnInstanceReadHandler(appLog, reqLog, errLog *log.Logger) {
	var handler instance
	handler.input = make(chan msg.Request, 128)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`instance_r`] = &handler
	go handler.run()
}

func spawnWorkflowReadHandler(appLog, reqLog, errLog *log.Logger) {
	var handler workflowRead
	handler.input = make(chan msg.Request, 128)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`workflow_r`] = &handler
	go handler.run()
}

func spawnWorkflowWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var handler workflowWrite
	handler.input = make(chan msg.Request, 128)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`workflow_w`] = &handler
	go handler.run()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
