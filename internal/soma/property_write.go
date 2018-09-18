/*-
 * Copyright (c) 2016-2018, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma // import "github.com/mjolnir42/soma/internal/soma"

import (
	"database/sql"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/handler"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

// PropertyWrite handles write requests for properties
type PropertyWrite struct {
	Input                  chan msg.Request
	Shutdown               chan struct{}
	handlerName            string
	conn                   *sql.DB
	stmtAddCustom          *sql.Stmt
	stmtAddNative          *sql.Stmt
	stmtAddService         *sql.Stmt
	stmtAddServiceAttr     *sql.Stmt
	stmtAddSystem          *sql.Stmt
	stmtAddTemplate        *sql.Stmt
	stmtAddTemplateAttr    *sql.Stmt
	stmtRemoveCustom       *sql.Stmt
	stmtRemoveNative       *sql.Stmt
	stmtRemoveService      *sql.Stmt
	stmtRemoveServiceAttr  *sql.Stmt
	stmtRemoveSystem       *sql.Stmt
	stmtRemoveTemplate     *sql.Stmt
	stmtRemoveTemplateAttr *sql.Stmt
	appLog                 *logrus.Logger
	reqLog                 *logrus.Logger
	errLog                 *logrus.Logger
}

// newPropertyWrite return a new PropertyWrite handler with input
// buffer of length
func newPropertyWrite(length int) (string, *PropertyWrite) {
	w := &PropertyWrite{}
	w.handlerName = generateHandlerName() + `_w`
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	return w.handlerName, w
}

// Register initializes resources provided by the Soma app
func (w *PropertyWrite) Register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (w *PropertyWrite) RegisterRequests(hmap *handler.Map) {
	for _, section := range []string{
		// XXX INCOMPLETE
		msg.SectionPropertySystem,
		msg.SectionPropertyNative,
		msg.SectionPropertyCustom,
	} {
		for _, action := range []string{
			msg.ActionAdd,
			msg.ActionRemove,
		} {
			hmap.Request(section, action, w.handlerName)
		}
	}
}

// Intake exposes the Input channel as part of the handler interface
func (w *PropertyWrite) Intake() chan msg.Request {
	return w.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (w *PropertyWrite) PriorityIntake() chan msg.Request {
	return w.Intake()
}

// Run is the event loop for PropertyWrite
func (w *PropertyWrite) Run() {
	var err error

	for statement, prepStmt := range map[string]**sql.Stmt{
		stmt.PropertyCustomAdd:            &w.stmtAddCustom,
		stmt.PropertyCustomDel:            &w.stmtRemoveCustom,
		stmt.PropertyNativeAdd:            &w.stmtAddNative,
		stmt.PropertyNativeDel:            &w.stmtRemoveNative,
		stmt.PropertyServiceAdd:           &w.stmtAddService,
		stmt.PropertyServiceAttributeAdd:  &w.stmtAddServiceAttr,
		stmt.PropertyServiceAttributeDel:  &w.stmtRemoveServiceAttr,
		stmt.PropertyServiceDel:           &w.stmtRemoveService,
		stmt.PropertySystemAdd:            &w.stmtAddSystem,
		stmt.PropertySystemDel:            &w.stmtRemoveSystem,
		stmt.PropertyTemplateAdd:          &w.stmtAddTemplate,
		stmt.PropertyTemplateAttributeAdd: &w.stmtAddTemplateAttr,
		stmt.PropertyTemplateAttributeDel: &w.stmtRemoveTemplateAttr,
		stmt.PropertyTemplateDel:          &w.stmtRemoveTemplate,
	} {
		if *prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`property`, err, stmt.Name(statement))
		}
		defer (*prepStmt).Close()
	}

runloop:
	for {
		select {
		case <-w.Shutdown:
			break runloop
		case req := <-w.Input:
			w.process(&req)
		}
	}
}

// process is the request dispatcher
func (w *PropertyWrite) process(q *msg.Request) {
	result := msg.FromRequest(q)
	logRequest(w.reqLog, q)

	switch q.Action {
	case msg.ActionAdd:
		w.add(q, &result)
	case msg.ActionRemove:
		w.remove(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// add inserts a new property
func (w *PropertyWrite) add(q *msg.Request, mr *msg.Result) {
	switch q.Property.Type {
	case `custom`:
		w.addCustom(q, mr)
	case `native`:
		w.addNative(q, mr)
	case `service`, `template`:
		w.addService(q, mr)
	case `system`:
		w.addSystem(q, mr)
	default:
		mr.NotImplemented(fmt.Errorf("Unknown property type: %s",
			q.Property.Type))
	}
}

// addSystem inserts system properties
func (w *PropertyWrite) addSystem(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
	)

	if res, err = w.stmtAddSystem.Exec(
		q.Property.System.Name,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Property = append(mr.Property, q.Property)
	}
}

// addNative inserts native properties
func (w *PropertyWrite) addNative(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
	)

	if res, err = w.stmtAddNative.Exec(
		q.Property.Native.Name,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Property = append(mr.Property, q.Property)
	}
}

// addCustom inserts custom repository properties
func (w *PropertyWrite) addCustom(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
	)

	q.Property.Custom.ID = uuid.Must(uuid.NewV4()).String()
	if res, err = w.stmtAddCustom.Exec(
		q.Property.Custom.ID,
		q.Property.Custom.RepositoryID,
		q.Property.Custom.Name,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Property = append(mr.Property, q.Property)
	}
}

// addService inserts team services or global service templates
func (w *PropertyWrite) addService(q *msg.Request, mr *msg.Result) {
	var (
		res  sql.Result
		err  error
		tx   *sql.Tx
		attr proto.ServiceAttribute
	)

	if tx, err = w.conn.Begin(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	switch q.Property.Type {
	case `service`:
		if res, err = tx.Stmt(w.stmtAddService).Exec(
			q.Property.Service.TeamID,
			q.Property.Service.Name,
		); err != nil {
			mr.ServerError(err, q.Section)
			tx.Rollback()
			return
		}
	case `template`:
		if res, err = tx.Stmt(w.stmtAddTemplate).Exec(
			q.Property.Service.Name,
		); err != nil {
			mr.ServerError(err, q.Section)
			tx.Rollback()
			return
		}
	}
	if !mr.RowCnt(res.RowsAffected()) {
		tx.Rollback()
		return
	}

	for _, attr = range q.Property.Service.Attributes {
		switch q.Property.Type {
		case `service`:
			if res, err = tx.Stmt(w.stmtAddServiceAttr).Exec(
				q.Property.Service.TeamID,
				q.Property.Service.Name,
				attr.Name,
				attr.Value,
			); err != nil {
				mr.ServerError(err, q.Section)
				tx.Rollback()
				return
			}
		case `template`:
			if res, err = tx.Stmt(w.stmtAddTemplateAttr).Exec(
				q.Property.Service.Name,
				attr.Name,
				attr.Value,
			); err != nil {
				mr.ServerError(err, q.Section)
				tx.Rollback()
				return
			}
		}
		if !mr.RowCnt(res.RowsAffected()) {
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Property = append(mr.Property, q.Property)
}

// remove deletes a property
func (w *PropertyWrite) remove(q *msg.Request, mr *msg.Result) {
	switch q.Property.Type {
	case `custom`:
		w.removeCustom(q, mr)
	case `native`:
		w.removeNative(q, mr)
	case `service`, `template`:
		w.removeService(q, mr)
	case `system`:
		w.removeSystem(q, mr)
	default:
		mr.NotImplemented(fmt.Errorf("Unknown property type: %s",
			q.Property.Type))
	}
}

// removeSystem deletes a system property
func (w *PropertyWrite) removeSystem(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
	)

	if res, err = w.stmtRemoveSystem.Exec(
		q.Property.System.Name,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Property = append(mr.Property, q.Property)
	}
}

// removeNative deletes a native property
func (w *PropertyWrite) removeNative(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
	)
	if res, err = w.stmtRemoveNative.Exec(
		q.Property.Native.Name,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Property = append(mr.Property, q.Property)
	}
}

// removeCustom deletes a custom repository property
func (w *PropertyWrite) removeCustom(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
	)
	if res, err = w.stmtRemoveCustom.Exec(
		q.Property.Custom.RepositoryID,
		q.Property.Custom.ID,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Property = append(mr.Property, q.Property)
	}
}

// removeService deletes a team service or service template
func (w *PropertyWrite) removeService(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
		tx  *sql.Tx
	)

	if tx, err = w.conn.Begin(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	switch q.Property.Type {
	case `service`:
		if res, err = tx.Stmt(w.stmtRemoveServiceAttr).Exec(
			q.Property.Service.TeamID,
			q.Property.Service.Name,
		); err != nil {
			mr.ServerError(err, q.Section)
			tx.Rollback()
			return
		}
	case `template`:
		if res, err = tx.Stmt(w.stmtRemoveTemplateAttr).Exec(
			q.Property.Service.TeamID,
			q.Property.Service.Name,
		); err != nil {
			mr.ServerError(err, q.Section)
			tx.Rollback()
			return
		}
	}
	// services can have an arbitrary number of attributes, no
	// rows affected check

	switch q.Property.Type {
	case `service`:
		if res, err = tx.Stmt(w.stmtRemoveService).Exec(
			q.Property.Service.TeamID,
			q.Property.Service.Name,
		); err != nil {
			mr.ServerError(err, q.Section)
			tx.Rollback()
			return
		}
	case `template`:
		if res, err = tx.Stmt(w.stmtRemoveTemplate).Exec(
			q.Property.Service.Name,
		); err != nil {
			mr.ServerError(err, q.Section)
			tx.Rollback()
			return
		}
	}
	if !mr.RowCnt(res.RowsAffected()) {
		tx.Rollback()
		return
	}

	if err = tx.Commit(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Property = append(mr.Property, q.Property)
}

// ShutdownNow signals the handler to shut down
func (w *PropertyWrite) ShutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
