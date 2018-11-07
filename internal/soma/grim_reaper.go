/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"database/sql"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/handler"
	"github.com/mjolnir42/soma/internal/msg"
)

// XXX BUG shutdown information must be transported to the
// REST packages

// ShutdownInProgress signals activated system shutdown
var ShutdownInProgress bool

// GrimReaper handles requests for controlled shutdown
type GrimReaper struct {
	Input    chan msg.Request
	Shutdown chan struct{}
	conn     *sql.DB
	appLog   *logrus.Logger
	reqLog   *logrus.Logger
	errLog   *logrus.Logger
	soma     *Soma
}

// newGrimReaper returns a new GrimReaper handler with input
// buffer of length
func newGrimReaper(length int, s *Soma) (grim *GrimReaper) {
	grim = &GrimReaper{}
	grim.Input = make(chan msg.Request, length)
	grim.Shutdown = make(chan struct{})
	grim.soma = s
	return
}

// Register initializes resources provided by the Soma app
func (grim *GrimReaper) Register(c *sql.DB, l ...*logrus.Logger) {
	grim.conn = c
	grim.appLog = l[0]
	grim.reqLog = l[1]
	grim.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (grim *GrimReaper) RegisterRequests(hmap *handler.Map) {
	for _, action := range []string{
		msg.ActionShutdown,
	} {
		hmap.Request(msg.SectionSystem, action, `grimreaper`)
	}
}

// Intake exposes the Input channel as part of the handler interface
func (grim *GrimReaper) Intake() chan msg.Request {
	return grim.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (grim *GrimReaper) PriorityIntake() chan msg.Request {
	return grim.Intake()
}

// Run is the event loop for GrimReaper
func (grim *GrimReaper) Run() {
	// defer calls stack in LIFO order
	defer os.Exit(0)
	defer grim.conn.Close()

	var res bool
	lock := sync.Mutex{}

runloop:
	for {
		select {
		case <-grim.Shutdown:
			lock.Lock()
			go func() {
				req := msg.Request{
					Section: msg.SectionSystem,
					Action:  msg.ActionShutdown,
				}
				req.Reply = make(chan msg.Result, 2)
				res = grim.process(&req)
				lock.Unlock()
			}()
		case req := <-grim.Input:
			// this is mainly so the go runtime does not optimize
			// away waiting for the shutdown routine
			lock.Lock()
			go func() {
				res = grim.process(&req)
				lock.Unlock()
			}()
		}
		break runloop
	}
	// blocks until the go routine has unlocked the mutex
	lock.Lock()
	if !res {
		lock.Unlock()
		goto runloop
	}

	time.Sleep(time.Duration(grim.soma.conf.ShutdownDelay) * time.Second)
	grim.appLog.Println("GrimReaper: shutdown complete")
}

// process is the request dispatcher
func (grim *GrimReaper) process(q *msg.Request) bool {
	result := msg.FromRequest(q)

	switch q.Action {
	case msg.ActionShutdown:
	default:
		result.UnknownRequest(q)
		q.Reply <- result
		return false
	}

	// tell HTTP handlers to start turning people away
	ShutdownInProgress = true

	// answer shutdown request
	result.OK()
	q.Reply <- result

	// give HTTP handlers time to turn people away
	time.Sleep(time.Duration(grim.soma.conf.ShutdownDelay) * time.Second)

	// I have awoken.
	grim.appLog.Println(`GRIM REAPER ACTIVATED. SYSTEM SHUTDOWN INITIATED`)

	// stop + shutdown all treeKeeper   : /^treekeeper_/
	for handler := range grim.soma.handlerMap.Range() {
		if strings.HasPrefix(handler, `treekeeper_`) {
			grim.soma.handlerMap.Get(handler).ShutdownNow()
			grim.soma.handlerMap.Del(handler)
			grim.appLog.Printf("GrimReaper: shut down %s", handler)
		}
	}
	// shutdown all write handler: /_w$/
	for handler := range grim.soma.handlerMap.Range() {
		if !strings.HasSuffix(handler, `_w`) {
			continue
		}
		grim.soma.handlerMap.Get(handler).ShutdownNow()
		grim.soma.handlerMap.Del(handler)
		grim.appLog.Printf("GrimReaper: shut down %s", handler)
	}
	// shutdown all read handler : /_r$/
	for handler := range grim.soma.handlerMap.Range() {
		if !strings.HasSuffix(handler, `_r`) {
			continue
		}
		grim.soma.handlerMap.Get(handler).ShutdownNow()
		grim.soma.handlerMap.Del(handler)
		grim.appLog.Printf("GrimReaper: shut down %s", handler)
	}
	// shutdown special handlers
	for _, h := range []string{
		`job_block`,
		`forest_custodian`,
		`guidepost`,
		`lifecycle`,
		`deployment`,
	} {
		grim.soma.handlerMap.Get(h).ShutdownNow()
		grim.soma.handlerMap.Del(h)
		grim.appLog.Printf("GrimReaper: shut down %s", h)
	}

	// shutdown supervisor -- needs handling in BasicAuth()
	grim.soma.handlerMap.Get(`supervisor`).ShutdownNow()
	grim.soma.handlerMap.Del(`supervisor`)
	grim.appLog.Println(`GrimReaper: shut down the supervisor`)

	// log what we have missed
	grim.appLog.Println(`GrimReaper: checking for still running handlers`)
	for name := range grim.soma.handlerMap.Range() {
		if name == `grimreaper` {
			continue
		}
		grim.appLog.Printf("GrimReaper: %s is still running", name)
	}

	return true
}

// ShutdownNow signals the handler to shut down. In the case of the
// GrimReaper, this will shut down SOMA
func (grim *GrimReaper) ShutdownNow() {
	close(grim.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
