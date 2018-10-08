/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/proto"
)

func (s *Supervisor) permissionRead(q *msg.Request, mr *msg.Result) {
	switch q.Action {
	case msg.ActionList:
		s.permissionList(q, mr)
	case msg.ActionShow:
		s.permissionShow(q, mr)
	case msg.ActionSearch:
		s.permissionSearch(q, mr)
	}
}

func (s *Supervisor) permissionList(q *msg.Request, mr *msg.Result) {
	var (
		err      error
		rows     *sql.Rows
		id, name string
	)

	if rows, err = s.stmtPermissionList.Query(
		q.Permission.Category,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&id,
			&name,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Permission = append(mr.Permission, proto.Permission{
			ID:       id,
			Name:     name,
			Category: q.Permission.Category,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

func (s *Supervisor) permissionShow(q *msg.Request, mr *msg.Result) {
	var (
		err  error
		tx   *sql.Tx
		perm proto.Permission
	)
	txMap := map[string]*sql.Stmt{}

	// open multi-statement transaction, set it readonly
	if tx, err = s.conn.Begin(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	if _, err = tx.Exec(stmt.ReadOnlyTransaction); err != nil {
		mr.ServerError(err, q.Section)
		tx.Rollback()
		return
	}

	// prepare statements for this transaction
	for name, statement := range map[string]string{
		`permission_show`:     stmt.PermissionShow,
		`permission_actions`:  stmt.PermissionMappedActions,
		`permission_sections`: stmt.PermissionMappedSections,
	} {
		if txMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("s.PermissionTx.Prepare(%s) error: %s",
				name, err.Error())
			mr.ServerError(err, q.Section)
			tx.Rollback()
			return
		}
	}

	if perm, err = s.permissionShowTx(
		q, txMap,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		tx.Rollback()
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		tx.Rollback()
		return
	}

	// close transaction
	if err = tx.Commit(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	mr.Permission = append(mr.Permission, perm)
	mr.OK()
}

func (s *Supervisor) permissionSearch(q *msg.Request, mr *msg.Result) {
	var (
		err      error
		rows     *sql.Rows
		id, name string
	)

	if rows, err = s.stmtPermissionSearch.Query(
		q.Search.Permission.Name,
		q.Search.Permission.Category,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&id,
			&name,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Permission = append(mr.Permission, proto.Permission{
			ID:       id,
			Name:     name,
			Category: q.Permission.Category,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

func (s *Supervisor) permissionShowTx(q *msg.Request,
	txMap map[string]*sql.Stmt) (proto.Permission, error) {
	var (
		err                                          error
		ts                                           time.Time
		id, name, category, user                     string
		perm                                         proto.Permission
		rows                                         *sql.Rows
		actionID, actionName, sectionID, sectionName string
	)

	// query base permission
	if err = txMap[`permission_show`].QueryRow(
		q.Permission.ID,
		q.Permission.Category,
	).Scan(
		&id,
		&name,
		&category,
		&user,
		&ts,
	); err != nil {
		return proto.Permission{}, err
	}

	perm = proto.Permission{
		ID:       id,
		Name:     name,
		Category: category,
		Actions:  &[]proto.Action{},
		Sections: &[]proto.Section{},
		Details: &proto.PermissionDetails{
			Creation: &proto.DetailsCreation{
				CreatedAt: ts.Format(msg.RFC3339Milli),
				CreatedBy: user,
			},
		},
	}

	// query mapped actions
	if rows, err = txMap[`permission_actions`].Query(
		q.Permission.ID,
		q.Permission.Category,
	); err != nil {
		return proto.Permission{}, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&actionID,
			&actionName,
			&sectionID,
			&sectionName,
			&category,
		); err != nil {
			rows.Close()
			return proto.Permission{}, err
		}
		*perm.Actions = append(*perm.Actions, proto.Action{
			ID:          actionID,
			Name:        actionName,
			SectionID:   sectionID,
			SectionName: sectionName,
			Category:    category,
		})
	}
	if err = rows.Err(); err != nil {
		return proto.Permission{}, err
	}

	// query mapped sections
	if rows, err = txMap[`permission_sections`].Query(
		q.Permission.ID,
		q.Permission.Category,
	); err != nil {
		return proto.Permission{}, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&sectionID,
			&sectionName,
			&category,
		); err != nil {
			rows.Close()
			return proto.Permission{}, err
		}
		*perm.Sections = append(*perm.Sections, proto.Section{
			ID:       sectionID,
			Name:     sectionName,
			Category: category,
		})
	}
	if err = rows.Err(); err != nil {
		return proto.Permission{}, err
	}

	if len(*perm.Actions) == 0 {
		perm.Actions = nil
	}
	if len(*perm.Sections) == 0 {
		perm.Sections = nil
	}

	return perm, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
