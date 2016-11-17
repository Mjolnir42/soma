/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

func (s *supervisor) category(q *msg.Request) {
	result := msg.FromRequest(q)

	s.requestLog(q)

	switch q.Action {
	case `list`, `show`:
		go func() { s.category_read(q) }()

	case `add`, `remove`:
		if s.readonly {
			result.Conflict(fmt.Errorf(`Readonly instance`))
			goto abort
		}
		s.category_write(q)

	default:
		result.UnknownRequest(q)
		goto abort
	}
	return

abort:
	q.Reply <- result
}

func (s *supervisor) category_read(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case `list`:
		s.category_list(q, &result)
	case `show`:
		s.category_show(q, &result)
	}
	q.Reply <- result
}

func (s *supervisor) category_write(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case `add`:
		s.category_add(q, &result)
	case `remove`:
		s.category_remove(q, &result)
		return
	}

	q.Reply <- result
}

func (s *supervisor) category_list(q *msg.Request, r *msg.Result) {
	var (
		err      error
		rows     *sql.Rows
		category string
	)
	if rows, err = s.stmt_ListCategory.Query(); err != nil {
		r.ServerError(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&category,
		); err != nil {
			r.ServerError(err)
			r.Clear(q.Section)
			return
		}
		r.Category = append(r.Category,
			proto.Category{Name: category})
	}
	if err = rows.Err(); err != nil {
		r.ServerError(err)
		r.Clear(q.Section)
	}
	r.OK()
}

func (s *supervisor) category_show(q *msg.Request, r *msg.Result) {
	var (
		err            error
		category, user string
		ts             time.Time
	)
	if err = s.stmt_ShowCategory.QueryRow(q.Category.Name).Scan(
		&category,
		&user,
		&ts,
	); err == sql.ErrNoRows {
		r.NotFound(err)
		return
	} else if err != nil {
		r.ServerError(err)
		return
	}
	r.Category = []proto.Category{proto.Category{
		Name: category,
		Details: &proto.CategoryDetails{
			CreatedAt: ts.Format(rfc3339Milli),
			CreatedBy: user,
		},
	}}
	r.OK()
}

func (s *supervisor) category_add(q *msg.Request, r *msg.Result) {
	var (
		err error
		tx  *sql.Tx
		res sql.Result
	)
	txMap := map[string]*sql.Stmt{}

	// open multi-statement transaction
	if tx, err = s.conn.Begin(); err != nil {
		r.ServerError(err)
		return
	}

	// prepare statements for this transaction
	for name, statement := range map[string]string{
		`category_add_tx_cat`:  stmt.CategoryAdd,
		`category_add_tx_perm`: stmt.PermissionAdd,
	} {
		if txMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("s.CategoryTx.Prepare(%s) error: %s",
				name, err.Error())
			r.ServerError(err)
			tx.Rollback()
			return
		}
	}

	if res, err = s.category_add_tx(q, txMap); err != nil {
		r.ServerError(err)
		tx.Rollback()
		return
	}
	// sets r.OK()
	if !r.RowCnt(res.RowsAffected()) {
		tx.Rollback()
		return
	}

	// close transaction
	if err = tx.Commit(); err != nil {
		r.ServerError(err)
		return
	}

	r.Category = []proto.Category{q.Category}
}

func (s *supervisor) category_add_tx(q *msg.Request,
	txMap map[string]*sql.Stmt) (sql.Result, error) {
	var (
		err    error
		res    sql.Result
		permId string
	)

	// create requested category
	if res, err = txMap[`category_add_tx_cat`].Exec(
		q.Category.Name,
		q.User,
	); err != nil {
		return res, err
	}

	// create grant category for requested category
	if res, err = txMap[`category_add_tx_cat`].Exec(
		fmt.Sprintf("%s:grant", q.Category.Name),
		q.User,
	); err != nil {
		return res, err
	}

	// create system permission for category, the category
	// name becomes the permission name in system
	permId = uuid.NewV4().String()
	return txMap[`category_add_tx_perm`].Exec(
		permId,
		q.Category.Name,
		`system`,
		q.User,
	)
}

func (s *supervisor) category_remove(q *msg.Request, r *msg.Result) {
	var (
		err error
		tx  *sql.Tx
		res sql.Result
	)
	txMap := map[string]*sql.Stmt{}

	// open multi-statement transaction
	if tx, err = s.conn.Begin(); err != nil {
		r.ServerError(err)
		return
	}

	// prepare statements for this transaction
	for name, statement := range map[string]string{
		`action_tx_remove`:        stmt.ActionRemove,
		`action_tx_removeMap`:     stmt.ActionRemoveFromMap,
		`section_tx_remove`:       stmt.SectionRemove,
		`section_tx_removeMap`:    stmt.SectionRemoveFromMap,
		`section_tx_actlist`:      stmt.SectionListActions,
		`category_tx_remove`:      stmt.CategoryRemove,
		`category_tx_seclist`:     stmt.CategoryListSections,
		`category_tx_permlist`:    stmt.CategoryListPermissions,
		`grant_tx_rm_system`:      stmt.GrantRemoveSystem,
		`permission_rm_tx_byname`: stmt.PermissionRemoveByName,
		`permission_rm_tx_lookup`: stmt.PermissionLookupGrantId,
		`permission_rm_tx_remove`: stmt.PermissionRemove,
		`permission_rm_tx_rev_gl`: stmt.PermissionRevokeGlobal,
		`permission_rm_tx_rev_mn`: stmt.PermissionRevokeMonitoring,
		`permission_rm_tx_rev_rp`: stmt.PermissionRevokeRepository,
		`permission_rm_tx_rev_tm`: stmt.PermissionRevokeTeam,
		`permission_rm_tx_unlink`: stmt.PermissionRemoveLink,
		`permission_rm_tx_unmapa`: stmt.PermissionUnmapAll,
	} {
		if txMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("s.CategoryTx.Prepare(%s) error: %s",
				name, err.Error())
			r.ServerError(err)
			tx.Rollback()
			return
		}
	}

	if res, err = s.category_remove_tx(q, txMap); err != nil {
		r.ServerError(err)
		tx.Rollback()
		return
	}

	// sets r.OK()
	if !r.RowCnt(res.RowsAffected()) {
		tx.Rollback()
		return
	}

	// close transaction
	if err = tx.Commit(); err != nil {
		r.ServerError(err)
		return
	}

}

func (s *supervisor) category_remove_tx(q *msg.Request,
	txMap map[string]*sql.Stmt) (sql.Result, error) {
	var (
		err                     error
		res                     sql.Result
		rows                    *sql.Rows
		sectionId, permissionId string
		affected                int64
	)

	// remove all sections from category
	if rows, err = txMap[`category_tx_seclist`].Query(
		q.Category.Name); err != nil {
		return res, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&sectionId,
		); err != nil {
			rows.Close()
			return res, err
		}
		if res, err = s.section_remove_tx(sectionId,
			txMap); err != nil {
			rows.Close()
			return res, err
		}
		if affected, err = res.RowsAffected(); err != nil {
			rows.Close()
			return res, err
		} else if affected != 1 {
			rows.Close()
			return res, fmt.Errorf("Delete statement caught %d"+
				" rows of sections instead of 1 (sectionId=%s)",
				affected, sectionId)
		}
	}
	if err = rows.Err(); err != nil {
		return res, err
	}

	// remove all permissions from category
	if rows, err = txMap[`category_tx_permlist`].Query(
		q.Category.Name); err != nil {
		return res, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&permissionId,
		); err != nil {
			rows.Close()
			return res, err
		}
		if res, err = s.permission_remove_tx(&msg.Request{
			Permission: proto.Permission{
				Id:       permissionId,
				Category: q.Category.Name,
			}}, txMap); err != nil {
			rows.Close()
			return res, err
		}
		if affected, err = res.RowsAffected(); err != nil {
			rows.Close()
			return res, err
		} else if affected != 1 {
			rows.Close()
			return res, fmt.Errorf("Delete statement caught %d"+
				" rows of permissions instead of 1 (permissionId=%s)",
				affected, permissionId)
		}
	}
	if err = rows.Err(); err != nil {
		return res, err
	}

	// remove all grants of system permission for category
	// ignore result since there can be any number of grants
	if _, err = txMap[`grant_tx_rm_system`].Exec(
		q.Category.Name); err != nil {
		return res, err
	}

	// remove system permission for category
	if res, err = txMap[`permission_rm_tx_byname`].Exec(
		q.Category.Name,
		`system`); err != nil {
		return res, err
	}
	if affected, err = res.RowsAffected(); err != nil {
		rows.Close()
		return res, err
	} else if affected != 1 {
		rows.Close()
		return res, fmt.Errorf("Delete statement caught %d"+
			" rows of permissions instead of 1 (system/%s)",
			affected, q.Category.Name)
	}

	// remove granting category
	if res, err = txMap[`category_tx_remove`].Exec(
		fmt.Sprintf("%s:grant", q.Category.Name)); err != nil {
		return res, err
	}
	if affected, err = res.RowsAffected(); err != nil {
		rows.Close()
		return res, err
	} else if affected != 1 {
		rows.Close()
		return res, fmt.Errorf("Delete statement caught %d"+
			" rows of categories instead of 1 (%s:grant)",
			affected, q.Category.Name)
	}

	// remove actual category
	return txMap[`category_tx_remove`].Exec(q.Category.Name)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
