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

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/proto"
)

// BucketRead handles read requests for buckets
type BucketRead struct {
	Input           chan msg.Request
	Shutdown        chan struct{}
	conn            *sql.DB
	stmtList        *sql.Stmt
	stmtShow        *sql.Stmt
	stmtPropOncall  *sql.Stmt
	stmtPropService *sql.Stmt
	stmtPropSystem  *sql.Stmt
	stmtPropCustom  *sql.Stmt
	appLog          *logrus.Logger
	reqLog          *logrus.Logger
	errLog          *logrus.Logger
}

// newBucketRead returns a new BucketRead handler with input
// buffer of length
func newBucketRead(length int) (r *BucketRead) {
	r = &BucketRead{}
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (r *BucketRead) register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// run is the event loop for BucketRead
func (r *BucketRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.BucketList:     r.stmtList,
		stmt.BucketShow:     r.stmtShow,
		stmt.BucketOncProps: r.stmtPropOncall,
		stmt.BucketSvcProps: r.stmtPropService,
		stmt.BucketSysProps: r.stmtPropSystem,
		stmt.BucketCstProps: r.stmtPropCustom,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`bucket`, err, stmt.Name(statement))
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
func (r *BucketRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(r.reqLog, q)

	switch q.Action {
	case `list`:
		r.list(q, &result)
	case `show`:
		r.show(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// list returns all buckets
func (r *BucketRead) list(q *msg.Request, mr *msg.Result) {
	var (
		bucketID, bucketName string
		rows                 *sql.Rows
		err                  error
	)

	if rows, err = r.stmtList.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&bucketID,
			&bucketName,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Bucket = append(mr.Bucket, proto.Bucket{
			Id:   bucketID,
			Name: bucketName,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// show returns the details of a specific bucket
func (r *BucketRead) show(q *msg.Request, mr *msg.Result) {
	var (
		ID, name, env, repoID, teamID string
		isDeleted, isFrozen           bool
		err                           error
	)

	if err = r.stmtShow.QueryRow(
		q.Bucket.Id,
	).Scan(
		&ID,
		&name,
		&isFrozen,
		&isDeleted,
		&repoID,
		&env,
		&teamID,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	bucket := proto.Bucket{
		Id:           ID,
		Name:         name,
		RepositoryId: repoID,
		TeamId:       teamID,
		Environment:  env,
		IsDeleted:    isDeleted,
		IsFrozen:     isFrozen,
	}

	// add properties
	bucket.Properties = &[]proto.Property{}

	if err = r.propertyOncall(q, &bucket); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if err = r.propertyService(q, &bucket); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if err = r.propertySystem(q, &bucket); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if err = r.propertyCustom(q, &bucket); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	mr.Bucket = append(mr.Bucket, bucket)
	mr.OK()
}

// propertyOncall adds the oncall properties to a bucket
func (r *BucketRead) propertyOncall(q *msg.Request, b *proto.Bucket) error {
	var (
		instanceID, sourceInstanceID string
		view, oncallID, oncallName   string
		rows                         *sql.Rows
		err                          error
	)

	if rows, err = r.stmtPropOncall.Query(
		q.Bucket.Id,
	); err != nil {
		return err
	}

	for rows.Next() {
		if err = rows.Scan(
			&instanceID,
			&sourceInstanceID,
			&view,
			&oncallID,
			&oncallName,
		); err != nil {
			rows.Close()
			return err
		}
		*b.Properties = append(*b.Properties,
			proto.Property{
				Type:             `oncall`,
				RepositoryId:     b.RepositoryId,
				BucketId:         b.Id,
				InstanceId:       instanceID,
				SourceInstanceId: sourceInstanceID,
				View:             view,
				Oncall: &proto.PropertyOncall{
					Id:   oncallID,
					Name: oncallName,
				},
			},
		)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

// propertyService adds the service properties to a bucket
func (r *BucketRead) propertyService(q *msg.Request, b *proto.Bucket) error {
	var (
		instanceID, sourceInstanceID string
		serviceName, view            string
		rows                         *sql.Rows
		err                          error
	)

	if rows, err = r.stmtPropService.Query(
		q.Bucket.Id,
	); err != nil {
		return err
	}

	for rows.Next() {
		if err = rows.Scan(
			&instanceID,
			&sourceInstanceID,
			&view,
			&serviceName,
		); err != nil {
			rows.Close()
			return err
		}
		*b.Properties = append(*b.Properties,
			proto.Property{
				Type:             `service`,
				RepositoryId:     b.RepositoryId,
				BucketId:         b.Id,
				InstanceId:       instanceID,
				SourceInstanceId: sourceInstanceID,
				View:             view,
				Service: &proto.PropertyService{
					Name: serviceName,
				},
			},
		)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

// propertySystem adds the system properties to a bucket
func (r *BucketRead) propertySystem(q *msg.Request, b *proto.Bucket) error {
	var (
		instanceID, sourceInstanceID, view string
		systemProp, value                  string
		rows                               *sql.Rows
		err                                error
	)

	if rows, err = r.stmtPropSystem.Query(
		q.Bucket.Id,
	); err != nil {
		return err
	}

	for rows.Next() {
		if err = rows.Scan(
			&instanceID,
			&sourceInstanceID,
			&view,
			&systemProp,
			&value,
		); err != nil {
			rows.Close()
			return err
		}
		*b.Properties = append(*b.Properties,
			proto.Property{
				Type:             `system`,
				RepositoryId:     b.RepositoryId,
				BucketId:         b.Id,
				InstanceId:       instanceID,
				SourceInstanceId: sourceInstanceID,
				View:             view,
				System: &proto.PropertySystem{
					Name:  systemProp,
					Value: value,
				},
			},
		)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

// propertyCustom adds the custom properties to a bucket
func (r *BucketRead) propertyCustom(q *msg.Request, b *proto.Bucket) error {
	var (
		instanceID, sourceInstanceID, view string
		customID, value, customProp        string
		rows                               *sql.Rows
		err                                error
	)

	if rows, err = r.stmtPropCustom.Query(
		q.Bucket.Id,
	); err != nil {
		return err
	}

	for rows.Next() {
		if err = rows.Scan(
			&instanceID,
			&sourceInstanceID,
			&view,
			&customID,
			&value,
			&customProp,
		); err != nil {
			rows.Close()
			return err
		}
		*b.Properties = append(*b.Properties,
			proto.Property{
				Type:             `custom`,
				RepositoryId:     b.RepositoryId,
				BucketId:         b.Id,
				InstanceId:       instanceID,
				SourceInstanceId: sourceInstanceID,
				View:             view,
				Custom: &proto.PropertyCustom{
					Id:    customID,
					Name:  customProp,
					Value: value,
				},
			},
		)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

// shutdown signals the handler to shut down
func (r *BucketRead) shutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
