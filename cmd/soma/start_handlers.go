package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/super"
)

func startHandlers(appLog, reqLog, errLog *log.Logger) {
	spawnSupervisorHandler(appLog, reqLog, errLog)
}

func spawnSupervisorHandler(appLog, reqLog, errLog *log.Logger) {
	var supervisorHandler super.Supervisor
	//var err error
	supervisorHandler.Input = make(chan msg.Request, 1024)
	supervisorHandler.Update = make(chan msg.Request, 1024)
	supervisorHandler.Shutdown = make(chan struct{})
	supervisorHandler.Register(conn, appLog, reqLog, errLog)
	/* XXX move to NewSupervisor function
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
	*/
	handlerMap[`supervisor`] = &supervisorHandler
	go supervisorHandler.Run()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
