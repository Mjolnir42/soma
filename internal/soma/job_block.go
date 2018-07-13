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
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/handler"
	"github.com/mjolnir42/soma/internal/msg"
)

// XXX BUG self|job|wait permission can currently block on every known job UUID

// JobBlock handles requests to block a client until an asynchronous job
// has finished
type JobBlock struct {
	Input     chan msg.Request
	Shutdown  chan struct{}
	Notify    chan string
	blockList map[string][]blockSpec
	jobDone   map[string]time.Time
	appLog    *logrus.Logger
	reqLog    *logrus.Logger
	errLog    *logrus.Logger
}

// blockSpec identifies a job that a client would like to block on
type blockSpec struct {
	JobID string
	RecvT time.Time
	Reply chan msg.Result
}

// newJobBlock returns a new JobBlock handler with input and notify
// buffers of length
func newJobBlock(length int) (j *JobBlock) {
	j = &JobBlock{}
	j.Input = make(chan msg.Request, length)
	j.Notify = make(chan string, length)
	j.Shutdown = make(chan struct{})
	j.blockList = make(map[string][]blockSpec)
	j.jobDone = make(map[string]time.Time)
	return
}

// Register initializes resources provided by the Soma app. This
// handler does not use the database connection, but accepts it
// to implement the interface
func (j *JobBlock) Register(c *sql.DB, l ...*logrus.Logger) {
	j.appLog = l[0]
	j.reqLog = l[1]
	j.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (j *JobBlock) RegisterRequests(hmap *handler.Map) {
	for _, section := range []string{
		msg.SectionJobMgmt,
		msg.SectionJob,
	} {
		hmap.Request(section, msg.ActionWait, `job_block`)
	}
}

// Intake exposes a dummy channel required to fulfull the Handler
// interface
func (j *JobBlock) Intake() chan msg.Request {
	return j.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (j *JobBlock) PriorityIntake() chan msg.Request {
	return j.Intake()
}

// Run is the event loop for JobBlock
func (j *JobBlock) Run() {
	tock := time.Tick(1 * time.Minute)
	j.jobDone = make(map[string]time.Time)
	j.blockList = make(map[string][]blockSpec)

runloop:
	for {
		select {
		case <-j.Shutdown:
			// a shutdown request was received, cleanup and disconnect
			// connected clients

			// clean jobDone cache
			for jID := range j.jobDone {
				delete(j.jobDone, jID)
			}
			// clean all block specifications
			for jID := range j.blockList {
				// disconnect all clients waiting on that job
				for i := range j.blockList[jID] {
					close(j.blockList[jID][i].Reply)
				}
				delete(j.blockList, jID)
			}
			break runloop
		case jID := <-j.Notify:
			// a job completion notification was received

			j.jobDone[jID] = time.Now().UTC()
			// unblock all clients blocking on this Job
			for _, bs := range j.blockList[jID] {
				close(bs.Reply)
			}
			delete(j.blockList, jID)
		case rq := <-j.Input:
			// a new block request was received
			bs := blockSpec{
				JobID: rq.Job.ID,
				Reply: rq.Reply,
			}

			// unblock immediate if the Job has already finished
			if _, ok := j.jobDone[bs.JobID]; ok {
				close(bs.Reply)
			} else {
				// register request to wait for this Job
				bs.RecvT = time.Now().UTC()
				if _, ok := j.blockList[bs.JobID]; !ok {
					j.blockList[bs.JobID] = []blockSpec{}
				}
				j.blockList[bs.JobID] = append(j.blockList[bs.JobID], bs)
			}
		case <-tock:
			// time for a periodic cleanup

			// clean jobDone cache - the notification about completed
			// jobs is only kept for 2 hours
			for jID := range j.jobDone {
				if time.Since(j.jobDone[jID]) > (2 * time.Hour) {
					delete(j.jobDone, jID)
				}
			}
			// disconnect active blocks that have waited for more than 5
			// minutes
			for jID := range j.blockList {
				newList := []blockSpec{}
				for i := range j.blockList[jID] {
					if time.Since(j.blockList[jID][i].RecvT) < (5 * time.Minute) {
						newList = append(newList, j.blockList[jID][i])
					} else {
						close(j.blockList[jID][i].Reply)
					}
				}
				if len(newList) > 0 {
					j.blockList[jID] = newList
				} else {
					delete(j.blockList, jID)
				}
			}
		}
	}
}

// ShutdownNow signals the handler to shutdown
func (j *JobBlock) ShutdownNow() {
	close(j.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
