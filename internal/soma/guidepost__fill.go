/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/proto"
	"github.com/satori/go.uuid"
)

//
func (g *GuidePost) fillReqData(q *msg.Request) (bool, error) {
	switch {
	case q.Action == msg.ActionPropertyCreate && q.Property.Type == `service`:
		return g.fillServiceAttributes(q)
	case q.Section == msg.SectionNodeConfig && q.Action == msg.ActionAssign:
		return g.fillNode(q)
	case q.Section == msg.SectionCheckConfig && q.Action == msg.ActionDestroy:
		return g.fillCheckDeleteInfo(q)
	case q.Section == msg.SectionBucket && q.Action == msg.ActionCreate:
		return g.fillBucketID(q)
	case q.Section == msg.SectionGroup && q.Action == msg.ActionCreate:
		return g.fillGroupID(q)
	case q.Section == msg.SectionCluster && q.Action == msg.ActionCreate:
		return g.fillClusterID(q)
	case q.Action == msg.ActionPropertyDestroy:
		return g.fillPropertyDeleteInfo(q)
	case q.Section == msg.SectionCheckConfig && q.Action == msg.ActionCreate:
		return g.fillCheckConfigID(q)
	default:
		return false, nil
	}
}

// generate CheckConfigId
func (g *GuidePost) fillCheckConfigID(q *msg.Request) (bool, error) {
	q.CheckConfig.ID = uuid.Must(uuid.NewV4()).String()
	return false, nil
}

// generate BucketID
func (g *GuidePost) fillBucketID(q *msg.Request) (bool, error) {
	q.Bucket.ID = uuid.Must(uuid.NewV4()).String()
	return false, nil
}

// generate GroupID
func (g *GuidePost) fillGroupID(q *msg.Request) (bool, error) {
	q.Group.ID = uuid.Must(uuid.NewV4()).String()
	return false, nil
}

// generate ClusterID
func (g *GuidePost) fillClusterID(q *msg.Request) (bool, error) {
	q.Cluster.ID = uuid.Must(uuid.NewV4()).String()
	return false, nil
}

// Populate the node structure with data, overwriting the client
// submitted values.
func (g *GuidePost) fillNode(q *msg.Request) (bool, error) {
	var (
		err                      error
		ndName, ndTeam, ndServer string
		ndAsset                  int64
		ndOnline, ndDeleted      bool
	)
	if err = g.stmtNodeDetails.QueryRow(q.Node.ID).Scan(
		&ndAsset,
		&ndName,
		&ndTeam,
		&ndServer,
		&ndOnline,
		&ndDeleted,
	); err != nil {
		if err == sql.ErrNoRows {
			return true, fmt.Errorf("Node not found: %s", q.Node.ID)
		}
		return false, err
	}
	q.Node.AssetID = uint64(ndAsset)
	q.Node.Name = ndName
	q.Node.TeamID = ndTeam
	q.Node.ServerID = ndServer
	q.Node.IsOnline = ndOnline
	q.Node.IsDeleted = ndDeleted
	return false, nil
}

// load authoritative copy of the service attributes from the
// database. Replaces whatever the client sent in.
func (g *GuidePost) fillServiceAttributes(q *msg.Request) (bool, error) {
	var (
		serviceID, attr, val, svName, svTeam, repoID string
		rows                                         *sql.Rows
		err                                          error
		nf                                           bool
	)
	attrs := []proto.ServiceAttribute{}

	switch q.Section {
	case msg.SectionRepositoryConfig:
		// svName may be the ID or the name
		serviceID = (*q.Repository.Properties)[0].Service.ID
		svName = (*q.Repository.Properties)[0].Service.Name
		svTeam = (*q.Repository.Properties)[0].Service.TeamID
	case msg.SectionBucket:
		serviceID = (*q.Bucket.Properties)[0].Service.ID
		svName = (*q.Bucket.Properties)[0].Service.Name
		svTeam = (*q.Bucket.Properties)[0].Service.TeamID
	case msg.SectionGroup:
		serviceID = (*q.Group.Properties)[0].Service.ID
		svName = (*q.Group.Properties)[0].Service.Name
		svTeam = (*q.Group.Properties)[0].Service.TeamID
	case msg.SectionCluster:
		serviceID = (*q.Cluster.Properties)[0].Service.ID
		svName = (*q.Cluster.Properties)[0].Service.Name
		svTeam = (*q.Cluster.Properties)[0].Service.TeamID
	case msg.SectionNodeConfig:
		serviceID = (*q.Node.Properties)[0].Service.ID
		svName = (*q.Node.Properties)[0].Service.Name
		svTeam = (*q.Node.Properties)[0].Service.TeamID
	}

	// ignore error since it would have been caught by GuidePost
	repoID, _, _, _ = g.extractRouting(q)

	// validate the tuple (repo, team, service) is valid.
	// also resolve and disambiguate serviceID and serviceName
	if err = g.stmtServiceLookup.QueryRow(
		repoID, serviceID, svName, svTeam,
	).Scan(
		&serviceID,
		&svName,
	); err != nil {
		if err == sql.ErrNoRows {
			nf = true
			err = fmt.Errorf("Requested service %s not available for team %s",
				svName, svTeam)
		}
		goto abort
	}

	// load attributes
	if rows, err = g.stmtServiceAttributes.Query(
		repoID, serviceID, svTeam,
	); err != nil {
		goto abort
	}
	defer rows.Close()

attrloop:
	for rows.Next() {
		if err = rows.Scan(&attr, &val); err != nil {
			break attrloop
		}
		attrs = append(attrs, proto.ServiceAttribute{
			Name:  attr,
			Value: val,
		})
	}
abort:
	if err != nil {
		return nf, err
	}
	// not aborted: set the loaded attributes
	switch q.Section {
	case msg.SectionRepositoryConfig:
		(*q.Repository.Properties)[0].Service.Attributes = attrs
	case msg.SectionBucket:
		(*q.Bucket.Properties)[0].Service.Attributes = attrs
	case msg.SectionGroup:
		(*q.Group.Properties)[0].Service.Attributes = attrs
	case msg.SectionCluster:
		(*q.Cluster.Properties)[0].Service.Attributes = attrs
	case msg.SectionNodeConfig:
		(*q.Node.Properties)[0].Service.Attributes = attrs
	}
	return false, nil
}

// if the request is a check deletion, populate required IDs
func (g *GuidePost) fillCheckDeleteInfo(q *msg.Request) (bool, error) {
	var delObjID, delObjTyp, delSrcChkID string
	var err error

	if err = g.stmtCheckDetailsForDelete.QueryRow(
		q.CheckConfig.ID,
		q.CheckConfig.RepositoryID,
	).Scan(
		&delObjID,
		&delObjTyp,
		&delSrcChkID,
	); err != nil {
		if err == sql.ErrNoRows {
			return true, fmt.Errorf(
				"Failed to find source check for config %s",
				q.CheckConfig.ID)
		}
		return false, err
	}
	q.CheckConfig.ObjectID = delObjID
	q.CheckConfig.ObjectType = delObjTyp
	q.CheckConfig.ExternalID = delSrcChkID
	return false, nil
}

// if the request is a property deletion, populate required IDs
func (g *GuidePost) fillPropertyDeleteInfo(q *msg.Request) (bool, error) {
	var (
		err                                             error
		row                                             *sql.Row
		queryStmt, view, sysProp, value, cstID, cstProp string
		svcID, oncID, oncName                           string
		oncNumber                                       int
	)

	// select SQL statement
	switch q.Section {
	case msg.SectionRepositoryConfig:
		switch q.Property.Type {
		case msg.PropertySystem:
			queryStmt = stmt.RepoSystemPropertyForDelete
		case msg.PropertyCustom:
			queryStmt = stmt.RepoCustomPropertyForDelete
		case msg.PropertyService:
			queryStmt = stmt.RepoServicePropertyForDelete
		case msg.PropertyOncall:
			queryStmt = stmt.RepoOncallPropertyForDelete
		}
	case msg.SectionBucket:
		switch q.Property.Type {
		case msg.PropertySystem:
			queryStmt = stmt.BucketSystemPropertyForDelete
		case msg.PropertyCustom:
			queryStmt = stmt.BucketCustomPropertyForDelete
		case msg.PropertyService:
			queryStmt = stmt.BucketServicePropertyForDelete
		case msg.PropertyOncall:
			queryStmt = stmt.BucketOncallPropertyForDelete
		}
	case msg.SectionGroup:
		switch q.Property.Type {
		case msg.PropertySystem:
			queryStmt = stmt.GroupSystemPropertyForDelete
		case msg.PropertyCustom:
			queryStmt = stmt.GroupCustomPropertyForDelete
		case msg.PropertyService:
			queryStmt = stmt.GroupServicePropertyForDelete
		case msg.PropertyOncall:
			queryStmt = stmt.GroupOncallPropertyForDelete
		}
	case msg.SectionCluster:
		switch q.Property.Type {
		case msg.PropertySystem:
			queryStmt = stmt.ClusterSystemPropertyForDelete
		case msg.PropertyCustom:
			queryStmt = stmt.ClusterCustomPropertyForDelete
		case msg.PropertyService:
			queryStmt = stmt.ClusterServicePropertyForDelete
		case msg.PropertyOncall:
			queryStmt = stmt.ClusterOncallPropertyForDelete
		}
	case msg.SectionNodeConfig:
		switch q.Property.Type {
		case msg.PropertySystem:
			queryStmt = stmt.NodeSystemPropertyForDelete
		case msg.PropertyCustom:
			queryStmt = stmt.NodeCustomPropertyForDelete
		case msg.PropertyService:
			queryStmt = stmt.NodeServicePropertyForDelete
		case msg.PropertyOncall:
			queryStmt = stmt.NodeOncallPropertyForDelete
		}
	}

	// execute and scan
	switch q.Section {
	case msg.SectionRepository, msg.SectionRepositoryConfig:
		row = g.conn.QueryRow(queryStmt,
			(*q.Repository.Properties)[0].SourceInstanceID)
	case msg.SectionBucket:
		row = g.conn.QueryRow(queryStmt,
			(*q.Bucket.Properties)[0].SourceInstanceID)
	case msg.SectionGroup:
		row = g.conn.QueryRow(queryStmt,
			(*q.Group.Properties)[0].SourceInstanceID)
	case msg.SectionCluster:
		row = g.conn.QueryRow(queryStmt,
			(*q.Cluster.Properties)[0].SourceInstanceID)
	case msg.SectionNode:
		row = g.conn.QueryRow(queryStmt,
			(*q.Node.Properties)[0].SourceInstanceID)
	}
	switch q.Property.Type {
	case msg.PropertySystem:
		err = row.Scan(&view, &sysProp, &value)
	case msg.PropertyCustom:
		err = row.Scan(&view, &cstID, &value, &cstProp)
	case msg.PropertyService:
		err = row.Scan(&view, &svcID)
	case msg.PropertyOncall:
		err = row.Scan(&view, &oncID, &oncName, &oncNumber)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return true, fmt.Errorf(
				"Failed to find source property for %s",
				(*q.Repository.Properties)[0].SourceInstanceID)
		}
		return false, err
	}

	// assemble and set results: property specification
	var (
		pSys *proto.PropertySystem
		pCst *proto.PropertyCustom
		pSvc *proto.PropertyService
		pOnc *proto.PropertyOncall
	)
	switch q.Property.Type {
	case msg.PropertySystem:
		pSys = &proto.PropertySystem{
			Name:  sysProp,
			Value: value,
		}
	case msg.PropertyCustom:
		pCst = &proto.PropertyCustom{
			ID:    cstID,
			Name:  cstProp,
			Value: value,
		}
	case msg.PropertyService:
		pSvc = &proto.PropertyService{
			ID: svcID,
		}
	case msg.PropertyOncall:
		num := strconv.Itoa(oncNumber)
		pOnc = &proto.PropertyOncall{
			ID:     oncID,
			Name:   oncName,
			Number: num,
		}
	}

	// assemble and set results: view
	switch q.Section {
	case msg.SectionRepositoryConfig:
		(*q.Repository.Properties)[0].View = view
	case msg.SectionBucket:
		(*q.Bucket.Properties)[0].View = view
	case msg.SectionGroup:
		(*q.Group.Properties)[0].View = view
	case msg.SectionCluster:
		(*q.Cluster.Properties)[0].View = view
	case msg.SectionNodeConfig:
		(*q.Node.Properties)[0].View = view
	}

	// final assembly step
	switch q.Section {
	case msg.SectionRepositoryConfig:
		switch q.Property.Type {
		case msg.PropertySystem:
			(*q.Repository.Properties)[0].System = pSys
		case msg.PropertyCustom:
			(*q.Repository.Properties)[0].Custom = pCst
		case msg.PropertyService:
			(*q.Repository.Properties)[0].Service = pSvc
		case msg.PropertyOncall:
			(*q.Repository.Properties)[0].Oncall = pOnc
		}
	case msg.SectionBucket:
		switch q.Property.Type {
		case msg.PropertySystem:
			(*q.Bucket.Properties)[0].System = pSys
		case msg.PropertyCustom:
			(*q.Bucket.Properties)[0].Custom = pCst
		case msg.PropertyService:
			(*q.Bucket.Properties)[0].Service = pSvc
		case msg.PropertyOncall:
			(*q.Bucket.Properties)[0].Oncall = pOnc
		}
	case msg.SectionGroup:
		switch q.Property.Type {
		case msg.PropertySystem:
			(*q.Group.Properties)[0].System = pSys
		case msg.PropertyCustom:
			(*q.Group.Properties)[0].Custom = pCst
		case msg.PropertyService:
			(*q.Group.Properties)[0].Service = pSvc
		case msg.PropertyOncall:
			(*q.Group.Properties)[0].Oncall = pOnc
		}
	case msg.SectionCluster:
		switch q.Property.Type {
		case msg.PropertySystem:
			(*q.Cluster.Properties)[0].System = pSys
		case msg.PropertyCustom:
			(*q.Cluster.Properties)[0].Custom = pCst
		case msg.PropertyService:
			(*q.Cluster.Properties)[0].Service = pSvc
		case msg.PropertyOncall:
			(*q.Cluster.Properties)[0].Oncall = pOnc
		}
	case msg.SectionNodeConfig:
		switch q.Property.Type {
		case msg.PropertySystem:
			(*q.Node.Properties)[0].System = pSys
		case msg.PropertyCustom:
			(*q.Node.Properties)[0].Custom = pCst
		case msg.PropertyService:
			(*q.Node.Properties)[0].Service = pSvc
		case msg.PropertyOncall:
			(*q.Node.Properties)[0].Oncall = pOnc
		}
	}
	return false, err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
