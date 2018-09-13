/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/handler"
	"github.com/mjolnir42/soma/internal/msg"
	uuid "github.com/satori/go.uuid"
)

func logrotate(sigChan chan os.Signal, handlerMap *handler.Map) {
	for {
		select {
		case <-sigChan:
		fileloop:
			for name, lfHandle := range logFileMap {
				// treekeeper startup logs do not get rotated
				if strings.HasPrefix(name, `startup_`) {
					continue
				}
				err := lfHandle.Reopen()
				if err != nil {
					log.Printf("Error rotating logfile %s: %s\n", name, err)
					log.Println(`Shutting down system`)

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
			}
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
