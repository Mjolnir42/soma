/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
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
	uuid "github.com/satori/go.uuid"
)

func (s *Supervisor) permission(q *msg.Request) {
	result := msg.FromRequest(q)

	s.requestLog(q)

	switch q.Action {
	case msg.ActionList,
		msg.ActionSearchByName,
		msg.ActionShow:
		go func() { s.permissionRead(q) }()

	case msg.ActionAdd,
		msg.ActionRemove,
		msg.ActionMap,
		msg.ActionUnmap:
		if s.readonly {
			result.Conflict(fmt.Errorf(`Readonly instance`))
			goto abort
		}
		s.permissionWrite(q)

	default:
		result.UnknownRequest(q)
		goto abort
	}
	return

abort:
	q.Reply <- result
}

func (s *Supervisor) permissionWrite(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case msg.ActionAdd:
		switch q.Permission.Category {
		case msg.CategoryGlobal,
			msg.CategoryPermission,
			msg.CategoryOperation,
			msg.CategoryRepository,
			msg.CategoryTeam,
			msg.CategoryMonitoring:
			s.permissionAdd(q, &result)
		default:
			result.ServerError(fmt.Errorf(`Illegal category`))
		}
	case msg.ActionRemove:
		switch q.Permission.Category {
		case msg.CategoryGlobal,
			msg.CategoryPermission,
			msg.CategoryOperation,
			msg.CategoryRepository,
			msg.CategoryTeam,
			msg.CategoryMonitoring:
			s.permissionRemove(q, &result)
		default:
			result.ServerError(fmt.Errorf(`Illegal category`))
		}
	case msg.ActionMap:
		s.permissionMap(q, &result)
	case msg.ActionUnmap:
		s.permissionUnmap(q, &result)
	}

	if result.IsOK() {
		s.Update <- msg.CacheUpdateFromRequest(q)
	}

	q.Reply <- result
}

func (s *Supervisor) permissionAdd(q *msg.Request, r *msg.Result) {
	var (
		err error
		tx  *sql.Tx
		res sql.Result
	)
	txMap := map[string]*sql.Stmt{}

	// open multi-statement transaction
	if tx, err = s.conn.Begin(); err != nil {
		r.ServerError(err, q.Section)
		return
	}

	// prepare statements for this transaction
	for name, statement := range map[string]string{
		`permission_add_tx_perm`: stmt.PermissionAdd,
		`permission_add_tx_link`: stmt.PermissionLinkGrant,
	} {
		if txMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("s.PermissionTx.Prepare(%s) error: %s",
				name, err.Error())
			r.ServerError(err, q.Section)
			tx.Rollback()
			return
		}
	}

	if res, err = s.permissionAddTx(q, txMap); err != nil {
		r.ServerError(err, q.Section)
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
		r.ServerError(err, q.Section)
		return
	}

	r.Permission = []proto.Permission{q.Permission}
}

func (s *Supervisor) permissionAddTx(q *msg.Request,
	txMap map[string]*sql.Stmt) (sql.Result, error) {
	var (
		err                        error
		res                        sql.Result
		grantPermID, grantCategory string
	)
	q.Permission.Id = uuid.NewV4().String()
	grantPermID = uuid.NewV4().String()
	switch q.Permission.Category {
	case msg.CategoryGlobal:
		grantCategory = msg.CategoryGrantGlobal
	case msg.CategoryPermission:
		grantCategory = msg.CategoryGrantPermission
	case msg.CategoryOperation:
		grantCategory = msg.CategoryGrantOperation
	case msg.CategoryRepository:
		grantCategory = msg.CategoryGrantRepository
	case msg.CategoryTeam:
		grantCategory = msg.CategoryGrantTeam
	case msg.CategoryMonitoring:
		grantCategory = msg.CategoryGrantMonitoring
	}

	if res, err = txMap[`permission_add_tx_perm`].Exec(
		q.Permission.Id,
		q.Permission.Name,
		q.Permission.Category,
		q.AuthUser,
	); err != nil {
		return res, err
	}

	if res, err = txMap[`permission_add_tx_perm`].Exec(
		grantPermID,
		q.Permission.Name,
		grantCategory,
		q.AuthUser,
	); err != nil {
		return res, err
	}

	return txMap[`permission_add_tx_link`].Exec(
		grantCategory,
		grantPermID,
		q.Permission.Category,
		q.Permission.Id,
	)
}

func (s *Supervisor) permissionRemove(q *msg.Request, r *msg.Result) {
	var (
		err error
		tx  *sql.Tx
		res sql.Result
	)
	txMap := map[string]*sql.Stmt{}

	// open multi-statement transaction
	if tx, err = s.conn.Begin(); err != nil {
		r.ServerError(err, q.Section)
		return
	}

	// prepare statements for this transaction
	for name, statement := range map[string]string{
		`permission_rm_tx_rev_gl`: stmt.PermissionRevokeGlobal,
		`permission_rm_tx_rev_rp`: stmt.PermissionRevokeRepository,
		`permission_rm_tx_rev_tm`: stmt.PermissionRevokeTeam,
		`permission_rm_tx_rev_mn`: stmt.PermissionRevokeMonitoring,
		`permission_rm_tx_lookup`: stmt.PermissionLookupGrantId,
		`permission_rm_tx_unlink`: stmt.PermissionRemoveLink,
		`permission_rm_tx_remove`: stmt.PermissionRemove,
		`permission_rm_tx_unmapa`: stmt.PermissionUnmapAll,
	} {
		if txMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("s.PermissionTx.Prepare(%s) error: %s",
				name, err.Error())
			r.ServerError(err, q.Section)
			tx.Rollback()
			return
		}
	}

	if res, err = s.permissionRemoveTx(q, txMap); err != nil {
		r.ServerError(err, q.Section)
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
		r.ServerError(err, q.Section)
		return
	}

	r.Permission = []proto.Permission{q.Permission}
}

func (s *Supervisor) permissionRemoveTx(q *msg.Request,
	txMap map[string]*sql.Stmt) (sql.Result, error) {
	var (
		err                  error
		res                  sql.Result
		grantingPermissionID string
		revocation           string
	)

	// select correct revocation statement scope
	switch q.Permission.Category {
	case msg.CategoryGlobal,
		msg.CategoryPermission,
		msg.CategoryOperation:
		revocation = `permission_rm_tx_rev_gl`
	case msg.CategoryRepository:
		revocation = `permission_rm_tx_rev_rp`
	case msg.CategoryTeam:
		revocation = `permission_rm_tx_rev_tm`
	case msg.CategoryMonitoring:
		revocation = `permission_rm_tx_rev_mn`
	}

	// lookup which permission grants this permission
	if err = txMap[`permission_rm_tx_lookup`].QueryRow(
		q.Permission.Id,
	).Scan(
		&grantingPermissionID,
	); err != nil {
		return res, err
	}

	// revoke all grants of the granting permission
	if res, err = txMap[revocation].Exec(
		grantingPermissionID,
	); err != nil {
		return res, err
	}

	// sever the link between permission and granting permission
	if res, err = txMap[`permission_rm_tx_unlink`].Exec(
		q.Permission.Id,
	); err != nil {
		return res, err
	}

	// remove granting permission
	if res, err = txMap[`permission_rm_tx_remove`].Exec(
		grantingPermissionID,
	); err != nil {
		return res, err
	}

	// revoke all grants of the permission
	if res, err = txMap[revocation].Exec(
		q.Permission.Id,
	); err != nil {
		return res, err
	}

	// unmap all sections & actions from the permission
	if res, err = txMap[`permission_rm_tx_unmapa`].Exec(
		q.Permission.Id,
	); err != nil {
		return res, err
	}

	// remove permission
	return txMap[`permission_rm_tx_remove`].Exec(q.Permission.Id)
}

func (s *Supervisor) permissionMap(q *msg.Request, r *msg.Result) {
	var (
		err                 error
		res                 sql.Result
		sectionID, actionID sql.NullString
		mapID               string
	)
	if q.Permission.Actions != nil {
		sectionID.String = (*q.Permission.Actions)[0].SectionId
		sectionID.Valid = true
		actionID.String = (*q.Permission.Actions)[0].Id
		actionID.Valid = true
	} else if q.Permission.Sections != nil {
		sectionID.String = (*q.Permission.Sections)[0].Id
		sectionID.Valid = true
	} else {
		r.ServerError(fmt.Errorf(`Nothing to map`))
		return
	}
	mapID = uuid.NewV4().String()

	if res, err = s.stmtPermissionMapEntry.Exec(
		mapID,
		q.Permission.Category,
		q.Permission.Id,
		sectionID,
		actionID,
	); err != nil {
		r.ServerError(err, q.Section)
		return
	}
	if r.RowCnt(res.RowsAffected()) {
		r.Permission = []proto.Permission{q.Permission}
	}
}

func (s *Supervisor) permissionUnmap(q *msg.Request, r *msg.Result) {
	var (
		err                 error
		res                 sql.Result
		sectionID, actionID sql.NullString
	)
	if q.Permission.Actions != nil {
		sectionID.String = (*q.Permission.Actions)[0].SectionId
		sectionID.Valid = true
		actionID.String = (*q.Permission.Actions)[0].Id
		actionID.Valid = true
	} else if q.Permission.Sections != nil {
		sectionID.String = (*q.Permission.Sections)[0].Id
		sectionID.Valid = true
	} else {
		r.ServerError(fmt.Errorf(`Nothing to map`))
		return
	}

	if res, err = s.stmtPermissionUnmapEntry.Exec(
		q.Permission.Id,
		q.Permission.Category,
		sectionID,
		actionID,
	); err != nil {
		r.ServerError(err, q.Section)
		return
	}
	if r.RowCnt(res.RowsAffected()) {
		r.Permission = []proto.Permission{q.Permission}
	}
}

func (s *Supervisor) permissionRead(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case msg.ActionList:
		s.permissionList(q, &result)
	case msg.ActionShow:
		s.permissionShow(q, &result)
	case msg.ActionSearchByName:
		s.permissionSearch(q, &result)
	}

	q.Reply <- result
}

func (s *Supervisor) permissionList(q *msg.Request, r *msg.Result) {
	var (
		err      error
		rows     *sql.Rows
		id, name string
	)
	if rows, err = s.stmtPermissionList.Query(
		q.Permission.Category,
	); err != nil {
		r.ServerError(err, q.Section)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&id,
			&name,
		); err != nil {
			r.ServerError(err, q.Section)
			return
		}
		r.Permission = append(r.Permission, proto.Permission{
			Id:       id,
			Name:     name,
			Category: q.Permission.Category,
		})
	}
	if err = rows.Err(); err != nil {
		r.ServerError(err, q.Section)
		return
	}
	r.OK()
}

func (s *Supervisor) permissionShow(q *msg.Request, r *msg.Result) {
	var (
		err                                          error
		tx                                           *sql.Tx
		ts                                           time.Time
		id, name, category, user                     string
		perm                                         proto.Permission
		rows                                         *sql.Rows
		actionID, actionName, sectionID, sectionName string
	)
	txMap := map[string]*sql.Stmt{}

	// open multi-statement transaction, set it readonly
	if tx, err = s.conn.Begin(); err != nil {
		r.ServerError(err, q.Section)
		return
	}
	if _, err = tx.Exec(stmt.ReadOnlyTransaction); err != nil {
		r.ServerError(err, q.Section)
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
			r.ServerError(err, q.Section)
			tx.Rollback()
			return
		}
	}

	if err = txMap[`permission_show`].QueryRow(
		q.Permission.Id,
		q.Permission.Category,
	).Scan(
		&id,
		&name,
		&category,
		&user,
		&ts,
	); err == sql.ErrNoRows {
		r.NotFound(err, q.Section)
		tx.Rollback()
		return
	} else if err != nil {
		r.ServerError(err, q.Section)
		tx.Rollback()
		return
	}
	perm = proto.Permission{
		Id:       id,
		Name:     name,
		Category: category,
		Actions:  &[]proto.Action{},
		Sections: &[]proto.Section{},
		Details: &proto.DetailsCreation{
			CreatedAt: ts.Format(msg.RFC3339Milli),
			CreatedBy: user,
		},
	}

	if rows, err = txMap[`permission_actions`].Query(
		q.Permission.Id,
		q.Permission.Category,
	); err != nil {
		r.ServerError(err, q.Section)
		tx.Rollback()
		return
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
			r.ServerError(err, q.Section)
			tx.Rollback()
			return
		}
		*perm.Actions = append(*perm.Actions, proto.Action{
			Id:          actionID,
			Name:        actionName,
			SectionId:   sectionID,
			SectionName: sectionName,
			Category:    category,
		})
	}
	if err = rows.Err(); err != nil {
		r.ServerError(err, q.Section)
		tx.Rollback()
		return
	}

	if rows, err = txMap[`permission_sections`].Query(
		q.Permission.Id,
		q.Permission.Category,
	); err != nil {
		r.ServerError(err, q.Section)
		tx.Rollback()
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&sectionID,
			&sectionName,
			&category,
		); err != nil {
			rows.Close()
			r.ServerError(err, q.Section)
			tx.Rollback()
			return
		}
		*perm.Sections = append(*perm.Sections, proto.Section{
			Id:       sectionID,
			Name:     sectionName,
			Category: category,
		})
	}
	if err = rows.Err(); err != nil {
		r.ServerError(err, q.Section)
		tx.Rollback()
		return
	}

	// close transaction
	if err = tx.Commit(); err != nil {
		r.ServerError(err, q.Section)
		return
	}

	if len(*perm.Actions) == 0 {
		perm.Actions = nil
	}
	if len(*perm.Sections) == 0 {
		perm.Sections = nil
	}
	r.Permission = append(r.Permission, perm)
	r.OK()
}

func (s *Supervisor) permissionSearch(q *msg.Request, r *msg.Result) {
	var (
		err      error
		rows     *sql.Rows
		id, name string
	)
	if rows, err = s.stmtPermissionSearch.Query(
		q.Permission.Name,
		q.Permission.Category,
	); err != nil {
		r.ServerError(err, q.Section)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&id,
			&name,
		); err != nil {
			r.ServerError(err, q.Section)
			return
		}
		r.Permission = append(r.Permission, proto.Permission{
			Id:       id,
			Name:     name,
			Category: q.Permission.Category,
		})
	}
	if err = rows.Err(); err != nil {
		r.ServerError(err, q.Section)
		return
	}
	r.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix