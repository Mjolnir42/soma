/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/handler"
	"github.com/mjolnir42/soma/internal/msg"
	uuid "github.com/satori/go.uuid"
)

func logrotate(sigChan chan os.Signal, handlerMap *handler.Map) {
	for {
		select {
		case <-sigChan:
			locked := true
		fileloop:
			for name, lfHandle := range logFileMap.Range() {
				// treekeeper startup logs do not get rotated
				if strings.HasPrefix(name, `startup_`) {
					continue
				}

				// reopen logfile handle
				err := lfHandle.Reopen()

				if err != nil {
					logFileMap.GetLogger(`error`).Errorf("Error rotating logfile %s: %s\n", name, err)
					logFileMap.GetLogger(`application`).Infoln(`Shutting down system`)

					logFileMap.RangeUnlock()
					locked = false

					returnChannel := make(chan msg.Result, 1)
					request := msg.Request{
						ID:         uuid.Must(uuid.FromString(`e0000000-e000-4000-e000-e00000000000`)),
						RemoteAddr: `::1`,
						AuthUser:   `root`,
						Reply:      returnChannel,
						Section:    msg.SectionSystem,
						Action:     msg.ActionShutdown,
					}
					handlerMap.MustLookup(&request).PriorityIntake() <- request
					<-returnChannel
					break fileloop
				}
				lg := logFileMap.GetLogger(name)
				lvl := lg.Level
				lg.SetLevel(logrus.InfoLevel)
				lg.Infoln(fmt.Sprintf("Reopened logfile `%s` for logrotate at %s",
					name,
					time.Now().UTC().Format(time.RFC3339),
				))
				lg.SetLevel(lvl)
			}
			if locked {
				logFileMap.RangeUnlock()
			}
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
