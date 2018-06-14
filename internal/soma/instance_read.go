/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/lib/pq"
	"github.com/mjolnir42/soma/internal/handler"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/proto"
)

// InstanceRead handles read requests for check instances
type InstanceRead struct {
	Input        chan msg.Request
	Shutdown     chan struct{}
	handlerName  string
	conn         *sql.DB
	stmtList     *sql.Stmt
	stmtShow     *sql.Stmt
	stmtVersions *sql.Stmt
	appLog       *logrus.Logger
	reqLog       *logrus.Logger
	errLog       *logrus.Logger
}

// newInstanceRead return a new InstanceRead handler with
// input buffer of length
func newInstanceRead(length int) (string, *InstanceRead) {
	r := &InstanceRead{}
	r.handlerName = generateHandlerName() + `_r`
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return r.handlerName, r
}

// Register initializes resources provided by the Soma app
func (r *InstanceRead) Register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (r *InstanceRead) RegisterRequests(hmap *handler.Map) {
	for _, action := range []string{
		msg.ActionAll,
		msg.ActionList,
		msg.ActionShow,
		msg.ActionVersions,
	} {
		hmap.Request(msg.SectionInstance, action, r.handlerName)
	}
}

// Intake exposes the Input channel as part of the handler interface
func (r *InstanceRead) Intake() chan msg.Request {
	return r.Input
}

// Run is the event loop for InstanceRead
func (r *InstanceRead) Run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.InstanceScopedList: r.stmtList,
		stmt.InstanceShow:       r.stmtShow,
		stmt.InstanceVersions:   r.stmtVersions,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`instance`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-r.Shutdown:
			break runloop
		case req := <-r.Input:
			go func() {
				r.process(&req)
			}()
		}
	}
}

// process is the request dispatcher
func (r *InstanceRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(r.reqLog, q)

	switch q.Action {
	case msg.ActionAll, msg.ActionList:
		r.list(q, &result)
	case msg.ActionShow:
		r.show(q, &result)
	case msg.ActionVersions:
		r.versions(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// show returns the details of a specific instance
func (r *InstanceRead) show(q *msg.Request, mr *msg.Result) {
	var (
		err                                      error
		version                                  int64
		isInherited                              bool
		instanceID, checkID, configID, details   string
		objectID, objectType, status, nextStatus string
		repositoryID, bucketID, instanceConfigID string
	)

	if err = r.stmtShow.QueryRow(
		q.Instance.ID,
	).Scan(
		&instanceID,
		&version,
		&checkID,
		&configID,
		&instanceConfigID,
		&repositoryID,
		&bucketID,
		&objectID,
		&objectType,
		&status,
		&nextStatus,
		&isInherited,
		&details,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	// unmarhal JSONB deployment details
	depl := proto.Deployment{}
	if err = json.Unmarshal([]byte(details), &depl); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	mr.Instance = append(mr.Instance, proto.Instance{
		ID:               instanceID,
		Version:          uint64(version),
		CheckID:          checkID,
		ConfigID:         configID,
		InstanceConfigID: instanceConfigID,
		RepositoryID:     repositoryID,
		BucketID:         bucketID,
		ObjectID:         objectID,
		ObjectType:       objectType,
		CurrentStatus:    status,
		NextStatus:       nextStatus,
		IsInherited:      isInherited,
		Deployment:       &depl,
	})
	mr.OK()
}

// list returns all instances either globally, within a repository
// or within a bucket
func (r *InstanceRead) list(q *msg.Request, mr *msg.Result) {
	var (
		err                                      error
		version                                  int64
		isInherited                              bool
		rows                                     *sql.Rows
		nullRepositoryID, nullBucketID           *sql.NullString
		instanceID, checkID, configID            string
		objectID, objectType, status, nextStatus string
		repositoryID, bucketID, instanceConfigID string
	)

	switch q.Instance.ObjectType {
	case msg.EntityRepository:
		nullRepositoryID.String = q.Instance.ObjectID
		nullRepositoryID.Valid = true
	case msg.EntityBucket:
		nullBucketID.String = q.Instance.ObjectID
		nullBucketID.Valid = true
	default:
		// only run an unscoped query if the flag has been explicitly
		// set
		if !(q.Flag.Unscoped && q.Action == msg.ActionAll) {
			mr.NotImplemented(
				fmt.Errorf("Instance listing for entity"+
					" type %s is currently not implemented",
					q.Instance.ObjectType,
				),
			)
			return
		}
	}

	if rows, err = r.stmtList.Query(
		nullRepositoryID,
		nullBucketID,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&instanceID,
			&version,
			&checkID,
			&configID,
			&instanceConfigID,
			&nullRepositoryID,
			&nullBucketID,
			&objectID,
			&objectType,
			&status,
			&nextStatus,
			&isInherited,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}

		if nullRepositoryID.Valid {
			repositoryID = nullRepositoryID.String
		}
		if nullBucketID.Valid {
			bucketID = nullBucketID.String
		}

		mr.Instance = append(mr.Instance, proto.Instance{
			ID:               instanceID,
			Version:          uint64(version),
			CheckID:          checkID,
			ConfigID:         configID,
			InstanceConfigID: instanceConfigID,
			RepositoryID:     repositoryID,
			BucketID:         bucketID,
			ObjectID:         objectID,
			ObjectType:       objectType,
			CurrentStatus:    status,
			NextStatus:       nextStatus,
			IsInherited:      isInherited,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// versions returns all versions of a specific instance
func (r *InstanceRead) versions(q *msg.Request, mr *msg.Result) {
	var (
		err                                              error
		version                                          int64
		isInherited                                      bool
		rows                                             *sql.Rows
		instanceID, status, nextStatus, instanceConfigID string
		createdNull, activatedNull, deprovisionedNull    pq.NullTime
		updatedNull, notifiedNull                        pq.NullTime
	)

	if rows, err = r.stmtVersions.Query(
		q.Instance.ID,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&instanceConfigID,
			&version,
			&instanceID,
			&createdNull,
			&activatedNull,
			&deprovisionedNull,
			&updatedNull,
			&notifiedNull,
			&status,
			&nextStatus,
			&isInherited,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		inst := proto.Instance{
			InstanceConfigID: instanceConfigID,
			Version:          uint64(version),
			ID:               instanceID,
			CurrentStatus:    status,
			NextStatus:       nextStatus,
			IsInherited:      isInherited,
			Info: &proto.InstanceVersionInfo{
				// created timestamp is a not null column
				CreatedAt: createdNull.Time.Format(msg.RFC3339Milli),
			},
		}
		if activatedNull.Valid {
			inst.Info.ActivatedAt = activatedNull.Time.Format(
				msg.RFC3339Milli)
		}
		if deprovisionedNull.Valid {
			inst.Info.DeprovisionedAt = deprovisionedNull.Time.
				Format(msg.RFC3339Milli)
		}
		if updatedNull.Valid {
			inst.Info.StatusLastUpdatedAt = updatedNull.Time.
				Format(msg.RFC3339Milli)
		}
		if notifiedNull.Valid {
			inst.Info.NotifiedAt = notifiedNull.Time.Format(
				msg.RFC3339Milli)
		}
		mr.Instance = append(mr.Instance, inst)
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// ShutdownNow signals the handler to shut down
func (r *InstanceRead) ShutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
