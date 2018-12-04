/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package adm

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/mjolnir42/soma/lib/proto"
	resty "gopkg.in/resty.v0"
)

// LookupOncallID looks up the UUID for an oncall duty on the
// server with name s. Error is set if no such oncall duty was
// found or an error occurred.
// If s is already a UUID, then s is immediately returned.
func LookupOncallID(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	return oncallIDByName(s)
}

// LookupOncallDetails looks up the details for oncall duty s.
func LookupOncallDetails(s string) (string, string, error) {
	var oID, o string
	oID = s
	if !IsUUID(s) {
		var err error
		if o, err = LookupOncallID(s); err != nil {
			return ``, ``, err
		}
		oID = o
	}

	return oncallDetailsByID(oID)
}

// LookupOncallId looks up the UUID for a user on the server
// with username s. Error is set if no such user was found
// or an error occurred.
// If s is already a UUID, then s is immediately returned.
func LookupUserID(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	return userIDByUserName(s)
}

// LookupTeamID looks up the UUID for a team on the server
// with teamname s. Error is set if no such team was found
// or an error occurred.
// If such a team is found, r is set to the UUID of the team.
// If s is already a UUID, then r is immediately set.
func LookupTeamID(s string, r *string) error {
	if IsUUID(s) {
		*r = s
		return nil
	}
	return teamIDByName(s, r)
}

// LookupTeamByRepo looks up the UUID for the team that is the
// owner of a given repository s, which can be the name or UUID
// of the repository.
func LookupTeamByRepo(s string, r *string) error {
	var (
		bID string
		err error
	)

	if !IsUUID(s) {
		if bID, err = LookupRepoID(s); err != nil {
			return err
		}
	} else {
		bID = s
	}

	return teamIDByRepoID(bID, r)
}

// LookupTeamByBucket looks up the UUID for the team that is
// the owner of a given bucket s, which can be the name or
// UUID of the bucket.
func LookupTeamByBucket(s string) (string, error) {
	var (
		bID string
		err error
	)

	if !IsUUID(s) {
		if bID, err = LookupBucketID(s); err != nil {
			return ``, err
		}
	} else {
		bID = s
	}

	return teamIDByBucketID(bID)
}

// LookupTeamByNode looks up the UUID for the team that is
// the owner of a given node s, which can be the name or
// UUID of the node.
func LookupTeamByNode(s string) (string, error) {
	var (
		nID string
		err error
	)

	if !IsUUID(s) {
		if nID, err = LookupBucketID(s); err != nil {
			return ``, err
		}
	} else {
		nID = s
	}

	return teamIDByNodeID(nID)
}

// LookupRepoID looks up the UUID for a repository on the server
// with reponame s. Error is set if no such repository was found
// or an error occurred.
// If s is already a UUID, then s is immediately returned.
func LookupRepoID(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	return repoIDByName(s)
}

//  LookupRepoName looks up the name for a repository on the server
// with repoID id. Error is set if no such repository was found or
// an error occured.
// If id is not an UUID, then id is assumed to be the name. The result
// is set in name.
func LookupRepoName(id string, name *string) error {
	if !IsUUID(id) {
		*name = id
		return nil
	}
	return repoNameByID(id, name)
}

// LookupRepoByBucket looks up the UUI for a repository by either
// the UUID or name of a bucket in that repository.
func LookupRepoByBucket(s string) (string, error) {
	var (
		bID string
		err error
	)

	if !IsUUID(s) {
		if bID, err = LookupBucketID(s); err != nil {
			return ``, err
		}
	} else {
		bID = s
	}

	return repoIDByBucketID(bID)
}

// LookupBucketID looks up the UUID for a bucket on the server
// with bucketname s. Error is set if no such bucket was found
// or an error occurred.
// If s is already a UUID, then s is immediately returned.
func LookupBucketID(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	return bucketIDByName(s)
}

// LookupGroupID looks up the UUID for group group in bucket
// bucket on the server.
// If group is already a UUID, then group is immediately returned.
func LookupGroupID(group, bucket string) (string, error) {
	if IsUUID(group) {
		return group, nil
	}
	var (
		bID string
		err error
	)
	if !IsUUID(bucket) {
		if bID, err = LookupBucketID(bucket); err != nil {
			return ``, err
		}
	} else {
		bID = bucket
	}

	return groupIDByName(group, bID)
}

// LookupClusterID looks up the UUID for cluster cluster in
// bucket bucket on the server.
// If cluster is already a UUID, then cluster is immediately returned.
func LookupClusterID(cluster, bucket string) (string, error) {
	if IsUUID(cluster) {
		return cluster, nil
	}
	var (
		bID string
		err error
	)
	if !IsUUID(bucket) {
		if bID, err = LookupBucketID(bucket); err != nil {
			return ``, err
		}
	} else {
		bID = bucket
	}

	return clusterIDByName(cluster, bID)
}

// LookupServerID looks up the UUID for a server either in the
// local cache or on the server. Error is set if no such server
// was found or an error occurred.
// If s is already a UUID, then s is immediately returned.
// If s is a Uint64 number, then the serverlookup is by AssetID.
// Otherwise s is the server name.
func LookupServerID(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	if ok, num := isUint64(s); ok {
		return serverIDByAsset(s, num)
	}
	return serverIDByName(s)
}

// LookupPermIDRef looks up the UUID for a permission from
// the server. Error is set if no such permission was found or
// an error occurred. The permission must be in category c.
// If s is already a UUID, then is is immediately returned.
func LookupPermIDRef(s, c string, id *string) error {
	if IsUUID(s) {
		*id = s
		return nil
	}
	return permissionIDByName(s, c, id)
}

// LookupGrantIDRef looks up the UUID of a permission grant from
// the server and fills it into the provided id pointer.
// Error is set if no such grant was found or an error occurred.
func LookupGrantIDRef(rcptType, rcptID, permID, cat string,
	id *string) error {
	return grantIDFromServer(rcptType, rcptID, permID, cat, id)
}

// LookupMonitoringID looks up the UUID of the monitoring system
// with the name s. Returns immediately if s is a UUID.
func LookupMonitoringID(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	return monitoringIDByName(s)
}

// LookupNodeID looks up the UUID of the repository the bucket
// given via string s is part of. If s is a UUID, it is used as
// bucketID for the lookup.
func LookupNodeID(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	return nodeIDByName(s)
}

// LookupCapabilityID looks up the UUID of the capability with the
// name s. Returns immediately if s is a UUID.
func LookupCapabilityID(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	return capabilityIDByName(s)
}

// LookupSectionID looks up the UUID of the section with the name
// s. Returns immediately if s is a UUID.
func LookupSectionID(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	return sectionIDByName(s)
}

// LookupActionID looks up the UUID of the action with the name
// a in section s. Return immediately if a is a UUID.
func LookupActionID(a, s string) (string, error) {
	if IsUUID(a) {
		return a, nil
	}
	var sID string
	var err error
	if sID, err = LookupSectionID(s); err != nil {
		return ``, err
	}
	return actionIDByName(a, sID)
}

// LookupCategoryBySection returns the category for section s
func LookupCategoryBySection(s string) (string, error) {
	if IsUUID(s) {
		return categoryBySectionID(s, nil)
	}
	return categoryBySectionID(sectionIDByName(s))
}

// LookupNodeConfig looks up the node repo/bucket configuration
// given the name or UUID s of the node.
func LookupNodeConfig(s string) (*proto.NodeConfig, error) {
	var (
		nID string
		err error
	)

	if !IsUUID(s) {
		if nID, err = LookupNodeID(s); err != nil {
			return nil, err
		}
	} else {
		nID = s
	}

	return nodeConfigByID(nID)
}

// LookupCheckConfigID looks up the UUID of check configuration.
// Lookup requires either the name and repository of the check,
// or the id of a check instance that was created by this check
// configuration.
// When the lookup is performed via name and repository, if the
// name is already a UUID it is returned immediately.
func LookupCheckConfigID(name, repo, instance string) (string, string, error) {
	if name != `` && repo != `` {
		if IsUUID(name) {
			return name, ``, nil
		}
		var repoID, r string
		var err error
		if r, err = LookupRepoID(repo); err != nil {
			return ``, ``, err
		}
		repoID = r
		return checkConfigIDByName(name, repoID)
	} else if instance != `` {
		return checkConfigIDByInstance(instance)
	}
	return ``, ``, fmt.Errorf(`Invalid argument combination for CheckConfigID`)
}

// LookupCustomPropertyID looks up the UUID of a custom property s
// in Repository repo. Returns immediately if s is a UUID.
func LookupCustomPropertyID(s, repo string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	var rID, r string
	var err error
	if r, err = LookupRepoID(repo); err != nil {
		return ``, err
	}
	rID = r
	return propertyIDByName(`custom`, s, rID)
}

// LookupServicePropertyID looks up the id of a service property s
// of team team.
func LookupServicePropertyID(s, team string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	var tID string
	if IsUUID(team) {
		tID = team
	} else {
		if err := teamIDByName(team, &tID); err != nil {
			return ``, err
		}
	}
	return propertyIDByName(`service`, s, tID)
}

// LookupTemplatePropertyID looks up the id of a service template
// property
func LookupTemplatePropertyID(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	return propertyIDByName(proto.PropertyTypeTemplate, s, `none`)
}

// LookupLevelName looks up the long name of a level s, where s
// can be the level's long or short name.
func LookupLevelName(s string, name *string) error {
	return levelByName(s, name)
}

// LookupCheckObjectID looks up the UUID for whichever object a
// check was defined on
func LookupCheckObjectID(oType, oName, buck string) (string, error) {
	switch oType {
	case `repository`:
		return LookupRepoID(oName)

	case `bucket`:
		return LookupBucketID(oName)

	case `group`:
		return LookupGroupID(oName, buck)

	case `cluster`:
		return LookupClusterID(oName, buck)

	case `node`:
		return LookupNodeID(oName)
	}

	return ``, fmt.Errorf("Unknown object type: %s", oType)
}

// LookupJobResultID looks up the UUID for the JobResult with name s
// from the server
func LookupJobResultID(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	return jobMetaIDByName(`result`, s)
}

// LookupJobStatusID looks up the UUID for the JobStatus with name s
// from the server
func LookupJobStatusID(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	return jobMetaIDByName(`status`, s)
}

// LookupJobTypeID looks up the UUID for the JobType with name s from
// the server
func LookupJobTypeID(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	return jobMetaIDByName(`type`, s)
}

// oncallIDByName implements the actual serverside lookup of the
// oncall duty UUID
func oncallIDByName(oncall string) (string, error) {
	req := proto.NewOncallFilter()
	req.Filter.Oncall = &proto.OncallFilter{Name: oncall}

	res, err := fetchFilter(req, `/search/oncall/`)
	if err != nil {
		goto abort
	}

	if res.Oncalls == nil || len(*res.Oncalls) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	// check the received record against the input
	if oncall != (*res.Oncalls)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			oncall, (*res.Oncalls)[0].Name)
		goto abort
	}
	return (*res.Oncalls)[0].ID, nil

abort:
	return ``, fmt.Errorf("OncallId lookup failed: %s", err.Error())
}

// oncallDetailsByID implements the actual serverside lookup of
// the oncall duty details
func oncallDetailsByID(oncall string) (string, string, error) {
	res, err := fetchObjList(fmt.Sprintf("/oncall/%s", oncall))
	if err != nil {
		goto abort
	}

	if res.Oncalls == nil || len(*res.Oncalls) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if oncall != (*res.Oncalls)[0].ID {
		err = fmt.Errorf("OncallId mismatch: %s vs %s",
			oncall, (*res.Oncalls)[0].ID)
		goto abort
	}
	return (*res.Oncalls)[0].Name, (*res.Oncalls)[0].Number, nil

abort:
	return ``, ``, fmt.Errorf("OncallDetails lookup failed: %s",
		err.Error())
}

// userIDByUserName implements the actual serverside lookup of the
// user's UUID
func userIDByUserName(user string) (string, error) {
	req := proto.NewUserFilter()
	req.Filter.User.UserName = user

	res, err := fetchFilter(req, `/search/user/`)
	if err != nil {
		goto abort
	}

	if res.Users == nil || len(*res.Users) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	// check the received record against the input
	if user != (*res.Users)[0].UserName {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			user, (*res.Users)[0].UserName)
		goto abort
	}
	return (*res.Users)[0].ID, nil

abort:
	return ``, fmt.Errorf("UserID lookup failed: %s", err.Error())
}

// teamIDByName implements the actual serverside lookup of the
// team's UUID
func teamIDByName(team string, id *string) error {
	req := proto.NewTeamFilter()
	req.Filter.Team.Name = team

	res, err := fetchFilter(req, `/search/team/`)
	if err != nil {
		goto abort
	}

	if res.Teams == nil || len(*res.Teams) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	// check the received record against the input
	if team != (*res.Teams)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			team, (*res.Teams)[0].Name)
		goto abort
	}
	*id = (*res.Teams)[0].ID
	return nil

abort:
	return fmt.Errorf("TeamID lookup failed for %s: %s", team, err.Error())
}

// teamIDByRepoID implements the actual serverside lookup of
// a repository's TeamID
func teamIDByRepoID(repoID string, team *string) error {
	req := proto.NewRepositoryFilter()
	req.Filter.Repository.ID = repoID

	res, err := fetchFilter(req, `/search/repository/`)
	if err != nil {
		goto abort
	}

	if res.Repositories == nil || len(*res.Repositories) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	// check the received record against the input
	if repoID != (*res.Repositories)[0].ID {
		err = fmt.Errorf("RepositoryID mismatch: %s vs %s",
			repoID, (*res.Repositories)[0].ID)
		goto abort
	}
	*team = (*res.Repositories)[0].TeamID
	return nil

abort:
	return fmt.Errorf("TeamID lookup failed: %s",
		err.Error())
}

// teamIDByBucketID implements the actual serverside lookup of
// a bucket's TeamID
func teamIDByBucketID(bucket string) (string, error) {
	res, err := fetchObjList(fmt.Sprintf("/bucket/%s", bucket))
	if err != nil {
		goto abort
	}

	if res.Buckets == nil || len(*res.Buckets) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	// check the received record against the input
	if bucket != (*res.Buckets)[0].ID {
		err = fmt.Errorf("BucketID mismatch: %s vs %s",
			bucket, (*res.Buckets)[0].ID)
		goto abort
	}
	return (*res.Buckets)[0].TeamID, nil

abort:
	return ``, fmt.Errorf("TeamID lookup failed: %s",
		err.Error())
}

// teamIDByNodeID implements the actual serverside lookup of a
// node's TeamID
func teamIDByNodeID(node string) (string, error) {
	res, err := fetchObjList(fmt.Sprintf("/nodes/%s", node))
	if err != nil {
		goto abort
	}

	if res.Nodes == nil || len(*res.Nodes) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	// check the received record against the input
	if node != (*res.Nodes)[0].ID {
		err = fmt.Errorf("NodeId mismatch: %s vs %s",
			node, (*res.Nodes)[0].ID)
		goto abort
	}
	return (*res.Nodes)[0].TeamID, nil

abort:
	return ``, fmt.Errorf("TeamID lookup failed: %s",
		err.Error())
}

// repoIDByName implements the actual serverside lookup of the
// repo's UUID
func repoIDByName(repo string) (string, error) {
	req := proto.NewRepositoryFilter()
	req.Filter.Repository.Name = repo

	res, err := fetchFilter(req, `/search/repository/`)
	if err != nil {
		goto abort
	}

	if res.Repositories == nil || len(*res.Repositories) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	// check the received record against the input
	if repo != (*res.Repositories)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			repo, (*res.Repositories)[0].Name)
		goto abort
	}
	return (*res.Repositories)[0].ID, nil

abort:
	return ``, fmt.Errorf("RepositoryId lookup failed: %s",
		err.Error())
}

// repoNameById mplements the actual serverside lookup
func repoNameByID(id string, name *string) error {
	var teamID string
	if err := LookupTeamByRepo(id, &teamID); err != nil {
		return err
	}
	res, err := fetchObjList(fmt.Sprintf("/team/%s/repository/%s", teamID, id))
	if err != nil {
		goto abort
	}

	if res.Repositories == nil || len(*res.Repositories) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}
	*name = (*res.Repositories)[0].Name
	return nil

abort:
	*name = ``
	return fmt.Errorf("RepositoryName lookup failed: %s", err.Error())
}

// repoIDByBucketID implements the actual serverside lookup of the
// repo's UUID
func repoIDByBucketID(bucket string) (string, error) {
	res, err := fetchObjList(fmt.Sprintf("/bucket/%s", bucket))
	if err != nil {
		goto abort
	}

	if res.Buckets == nil || len(*res.Buckets) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	// check the received record against the input
	if bucket != (*res.Buckets)[0].ID {
		err = fmt.Errorf("BucketID mismatch: %s vs %s",
			bucket, (*res.Buckets)[0].ID)
		goto abort
	}
	return (*res.Buckets)[0].RepositoryID, nil

abort:
	return ``, fmt.Errorf("RepositoryId lookup failed: %s",
		err.Error())
}

// bucketIDByName implements the actual serverside lookup of the
// bucket's UUID
func bucketIDByName(bucket string) (string, error) {
	req := proto.NewBucketFilter()
	req.Filter.Bucket.Name = bucket

	res, err := fetchFilter(req, `/search/bucket/`)
	if err != nil {
		goto abort
	}

	if res.Buckets == nil || len(*res.Buckets) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	// check the received record against the input
	if bucket != (*res.Buckets)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			bucket, (*res.Buckets)[0].Name)
		goto abort
	}
	return (*res.Buckets)[0].ID, nil

abort:
	return ``, fmt.Errorf("BucketID lookup failed: %s",
		err.Error())
}

//
func groupIDByName(group, bucketID string) (string, error) {
	var (
		res                       *proto.Result
		err                       error
		repositoryID, requestPath string
	)

	req := proto.NewGroupFilter()
	req.Filter.Group.Name = group
	req.Filter.Group.BucketID = bucketID

	repositoryID, err = repoIDByBucketID(bucketID)
	if err != nil {
		goto abort
	}

	requestPath = fmt.Sprintf("/search/repository/%s/bucket/%s/group/",
		repositoryID,
		bucketID,
	)

	res, err = fetchFilter(req, requestPath)
	if err != nil {
		goto abort
	}

	if res.Groups == nil || len(*res.Groups) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if group != (*res.Groups)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			group, (*res.Groups)[0].Name)
	}
	return (*res.Groups)[0].ID, nil

abort:
	return ``, fmt.Errorf("GroupID lookup failed: %s",
		err.Error())
}

//
func clusterIDByName(cluster, bucketID string) (string, error) {
	req := proto.NewClusterFilter()
	req.Filter.Cluster.Name = cluster
	req.Filter.Cluster.BucketID = bucketID

	res, err := fetchFilter(req, `/search/clusters/`)
	if err != nil {
		goto abort
	}

	if res.Clusters == nil || len(*res.Clusters) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if cluster != (*res.Clusters)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			cluster, (*res.Clusters)[0].Name)
	}
	return (*res.Clusters)[0].ID, nil

abort:
	return ``, fmt.Errorf("ClusterID lookup failed: %s",
		err.Error())
}

// serverIDByName implements the actual lookup of the server UUID
// by name
func serverIDByName(s string) (string, error) {
	if m, err := cache.ServerByName(s); err == nil {
		return m[`id`], nil
	}
	req := proto.NewServerFilter()
	req.Filter.Server.Name = s

	res, err := fetchFilter(req, `/search/server/`)
	if err != nil {
		goto abort
	}

	if res.Servers == nil || len(*res.Servers) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if s != (*res.Servers)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			s, (*res.Servers)[0].Name)
		goto abort
	}
	// save server in cacheDB
	cache.Server(
		(*res.Servers)[0].Name,
		(*res.Servers)[0].ID,
		strconv.Itoa(int((*res.Servers)[0].AssetID)),
	)
	return (*res.Servers)[0].ID, nil

abort:
	return ``, fmt.Errorf("ServerID lookup failed: %s",
		err.Error())
}

// serverIDByAsset implements the actual lookup of the server UUID
// by numeric AssetID
func serverIDByAsset(s string, aid uint64) (string, error) {
	if m, err := cache.ServerByAsset(s); err == nil {
		return m[`id`], nil
	}
	req := proto.NewServerFilter()
	req.Filter.Server.AssetID = aid

	res, err := fetchFilter(req, `/search/server/`)
	if err != nil {
		goto abort
	}

	if res.Servers == nil || len(*res.Servers) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if aid != (*res.Servers)[0].AssetID {
		err = fmt.Errorf("AssetId mismatch: %d vs %d",
			aid, (*res.Servers)[0].AssetID)
		goto abort
	}
	// save server in cacheDB
	cache.Server(
		(*res.Servers)[0].Name,
		(*res.Servers)[0].ID,
		strconv.Itoa(int((*res.Servers)[0].AssetID)),
	)
	return (*res.Servers)[0].ID, nil

abort:
	return ``, fmt.Errorf("ServerID lookup failed: %s",
		err.Error())
}

// permissionIDByName implements the actual lookup of the permission
// UUID by name
func permissionIDByName(perm, cat string, id *string) error {
	req := proto.NewPermissionFilter()
	req.Filter.Permission.Name = perm
	req.Filter.Permission.Category = cat

	res, err := fetchFilter(req, `/search/permission/`)
	if err != nil {
		goto abort
	}

	if res.Permissions == nil || len(*res.Permissions) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if perm != (*res.Permissions)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			perm, (*res.Permissions)[0].Name)
		goto abort
	}
	*id = (*res.Permissions)[0].ID
	return nil

abort:
	return fmt.Errorf("PermissionID lookup failed: %s",
		err.Error())
}

// grantIDFromServer implements the actual lookup of the grant UUID
func grantIDFromServer(rcptType, rcptID, permID, cat string,
	id *string) error {
	req := proto.NewGrantFilter()
	req.Filter.Grant.RecipientType = rcptType
	req.Filter.Grant.RecipientID = rcptID
	req.Filter.Grant.PermissionID = permID
	req.Filter.Grant.Category = cat

	res, err := fetchFilter(req, `/filter/grant/`)
	if err != nil {
		goto abort
	}

	if res.Grants == nil || len(*res.Grants) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if permID != (*res.Grants)[0].PermissionID {
		err = fmt.Errorf("PermissionID mismatch: %s vs %s",
			permID, (*res.Grants)[0].PermissionID)
		goto abort
	}
	*id = (*res.Grants)[0].ID

abort:
	return fmt.Errorf("GrantId lookup failed: %s",
		err.Error())
}

// monitoringIDByName implements the actual lookup of the monitoring
// system UUID
func monitoringIDByName(monitoring string) (string, error) {
	req := proto.NewMonitoringFilter()
	req.Filter.Monitoring.Name = monitoring

	res, err := fetchFilter(req, `/search/monitoringsystem/`)
	if err != nil {
		goto abort
	}

	if res.Monitorings == nil || len(*res.Monitorings) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if monitoring != (*res.Monitorings)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			monitoring, (*res.Monitorings)[0].Name)
		goto abort
	}
	return (*res.Monitorings)[0].ID, nil

abort:
	return ``, fmt.Errorf("MonitoringId lookup failed: %s",
		err.Error())
}

// nodeIDByName implements the actual lookup of the node UUID
func nodeIDByName(node string) (string, error) {
	req := proto.NewNodeFilter()
	req.Filter.Node.Name = node

	res, err := fetchFilter(req, `/filter/nodes/`)
	if err != nil {
		goto abort
	}

	if res.Nodes == nil || len(*res.Nodes) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if node != (*res.Nodes)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			node, (*res.Nodes)[0].Name)
		goto abort
	}
	return (*res.Nodes)[0].ID, nil

abort:
	return ``, fmt.Errorf("NodeId lookup failed: %s",
		err.Error())
}

// sectionIDByName implements the actual lookup of the section
// UUID from the server
func sectionIDByName(section string) (string, error) {
	req := proto.NewSectionFilter()
	req.Filter.Section.Name = section

	res, err := fetchFilter(req, `/search/section/`)
	if err != nil {
		goto abort
	}

	if res.Sections == nil || len(*res.Sections) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if section != (*res.Sections)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			section, (*res.Sections)[0].Name)
		goto abort
	}
	return (*res.Sections)[0].ID, nil

abort:
	return ``, fmt.Errorf("SectionID lookup failed: %s",
		err.Error())
}

func categoryBySectionID(section string, e error) (string, error) {
	req := proto.NewSectionFilter()
	req.Filter.Section.ID = section

	res, err := fetchFilter(req, `/search/section/`)
	if err != nil {
		goto abort
	}

	if res.Sections == nil || len(*res.Sections) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if section != (*res.Sections)[0].ID {
		err = fmt.Errorf("ID mismatch: %s vs %s",
			section, (*res.Sections)[0].ID)
		goto abort
	}
	return (*res.Sections)[0].Category, nil

abort:
	return ``, fmt.Errorf("Category by SectionID lookup failed: %s",
		err.Error())
}

// actionIDByName implements the actual lookup of the action
// UUID from the server
func actionIDByName(action, section string) (string, error) {
	req := proto.NewActionFilter()
	req.Filter.Action.Name = action
	req.Filter.Action.SectionID = section

	res, err := fetchFilter(req, `/search/action/`)
	if err != nil {
		goto abort
	}

	if res.Actions == nil || len(*res.Actions) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if action != (*res.Actions)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			action, (*res.Actions)[0].Name)
		goto abort
	}
	return (*res.Actions)[0].ID, nil

abort:
	return ``, fmt.Errorf("ActionID lookup failed: %s",
		err.Error())
}

// nodeConfigByID implements the actual lookup of the node's repo
// and bucket assignment information from the server
func nodeConfigByID(node string) (*proto.NodeConfig, error) {
	path := fmt.Sprintf("/nodes/%s/config", node)
	var (
		err  error
		resp *resty.Response
		res  *proto.Result
	)
	if resp, err = GetReq(path); err != nil {
		goto abort
	}
	if err = decodeResponse(resp, res); err != nil {
		goto abort
	}
	if res.StatusCode == 404 {
		err = fmt.Errorf(`Node is not assigned to a configuration` +
			` repository yet.`)
		goto abort
	}
	if err = checkApplicationError(res); err != nil {
		goto abort
	}

	if res.Nodes == nil || len(*res.Nodes) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	// check the received record against the input
	if node != (*res.Nodes)[0].ID {
		err = fmt.Errorf("NodeId mismatch: %s vs %s",
			node, (*res.Nodes)[0].ID)
		goto abort
	}
	return (*res.Nodes)[0].Config, nil

abort:
	return nil, fmt.Errorf("NodeConfig lookup failed: %s",
		err.Error())
}

// capabilityIDByName implements the actual lookup of the capability
// UUID from the server
func capabilityIDByName(cap string) (string, error) {
	var err error
	var res *proto.Result
	req := proto.NewCapabilityFilter()

	split := strings.SplitN(cap, ".", 3)
	if len(split) != 3 {
		err = fmt.Errorf(`Capability split failed, name invalid`)
		goto abort
	}
	if req.Filter.Capability.MonitoringID, err = LookupMonitoringID(
		split[0]); err != nil {
		goto abort
	}
	req.Filter.Capability.View = split[1]
	req.Filter.Capability.Metric = split[2]

	if res, err = fetchFilter(req, `/search/capability/`); err != nil {
		goto abort
	}

	if res.Capabilities == nil || len(*res.Capabilities) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if cap != (*res.Capabilities)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			cap, (*res.Capabilities)[0].Name)
		goto abort
	}
	return (*res.Capabilities)[0].ID, nil

abort:
	return ``, fmt.Errorf("CapabilityId lookup failed: %s",
		err.Error())
}

// checkConfigIDByName implements the actual lookup of the check
// configuration's UUID from the server by check config name
func checkConfigIDByName(check, repo string) (string, string, error) {
	req := proto.NewCheckConfigFilter()
	req.Filter.CheckConfig.Name = check

	res, err := fetchFilter(req, fmt.Sprintf(
		"/search/checkconfig/%s/", repo))
	if err != nil {
		goto abort
	}

	if res.CheckConfigs == nil || len(*res.CheckConfigs) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if check != (*res.CheckConfigs)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			check, (*res.CheckConfigs)[0].Name)
		goto abort
	}
	return (*res.CheckConfigs)[0].ID, repo, nil

abort:
	return ``, ``, fmt.Errorf("CheckConfigID lookup failed: %s",
		err.Error())
}

// checkConfigIDByInstance implements the actual lookup of the check
// configuration's UUID from the server via an instance ID it created
func checkConfigIDByInstance(instance string) (string, string, error) {
	res, err := fetchObjList(fmt.Sprintf("/instance/%s", instance))
	if err != nil {
		goto abort
	}

	if res.Instances == nil || len(*res.Instances) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if instance != (*res.Instances)[0].ID {
		err = fmt.Errorf("Id mismatch: %s vs %s",
			instance, (*res.Instances)[0].ID)
		goto abort
	}
	return (*res.Instances)[0].ConfigID,
		(*res.Instances)[0].RepositoryID, nil

abort:
	return ``, ``, fmt.Errorf("CheckConfigID lookup failed: %s",
		err.Error())
}

// propertyIDByName implements the actual lookup of property ids
// from the server
func propertyIDByName(pType, pName, refID string) (string, error) {
	req := proto.NewPropertyFilter()
	req.Filter.Property.Type = pType
	req.Filter.Property.Name = pName

	var (
		path string
		err  error
	)
	res := &proto.Result{}

	switch pType {
	case proto.PropertyTypeCustom:
		// custom properties are per-repository
		req.Filter.Property.RepositoryID = refID
		path = fmt.Sprintf("/search/repository/%s/property-mgmt/%s/",
			url.QueryEscape(refID),
			url.QueryEscape(proto.PropertyTypeCustom),
		)
	case proto.PropertyTypeService:
		path = fmt.Sprintf("/search/team/%s/property-mgmt/%s/",
			url.QueryEscape(refID),
			url.QueryEscape(proto.PropertyTypeService),
		)
	case proto.PropertyTypeTemplate:
		path = fmt.Sprintf("/search/property-mgmt/%s/",
			url.QueryEscape(proto.PropertyTypeTemplate),
		)
	default:
		err = fmt.Errorf("Unknown property type: %s", pType)
		goto abort
	}

	if res, err = fetchFilter(req, path); err != nil {
		goto abort
	}

	if res.Properties == nil || len(*res.Properties) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	switch pType {
	case proto.PropertyTypeCustom:
		if pName != (*res.Properties)[0].Custom.Name {
			err = fmt.Errorf("Name mismatch: %s vs %s",
				pName, (*res.Properties)[0].Custom.Name)
			goto abort
		}
		if refID != (*res.Properties)[0].Custom.RepositoryID {
			err = fmt.Errorf("RepositoryId mismatch: %s vs %s",
				refID, (*res.Properties)[0].Custom.RepositoryID)
			goto abort
		}
		return (*res.Properties)[0].Custom.ID, nil
	case proto.PropertyTypeService:
		if refID != (*res.Properties)[0].Service.TeamID {
			err = fmt.Errorf("TeamID mismatch: %s vs %s",
				refID, (*res.Properties)[0].Service.TeamID)
			goto abort
		}
		fallthrough
	case proto.PropertyTypeTemplate:
		if pName != (*res.Properties)[0].Service.Name {
			err = fmt.Errorf("Name mismatch: %s vs %s",
				pName, (*res.Properties)[0].Service.Name)
			goto abort
		}
		return (*res.Properties)[0].Service.ID, nil
	default:
		err = fmt.Errorf("Unknown property type: %s", pType)
	}

abort:
	return ``, fmt.Errorf("PropertyID lookup failed: %s", err.Error())
}

// levelByName implements the actual lookup of the level details
// from the server
func levelByName(lvl string, name *string) error {
	req := proto.NewLevelFilter()
	req.Filter.Level.Name = lvl
	req.Filter.Level.ShortName = lvl

	res, err := fetchFilter(req, `/search/level/`)
	if err != nil {
		goto abort
	}

	if res.Levels == nil || len(*res.Levels) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if lvl != (*res.Levels)[0].Name &&
		lvl != (*res.Levels)[0].ShortName {
		err = fmt.Errorf("Name mismatch: %s vs %s/%s",
			lvl, (*res.Levels)[0].Name, (*res.Levels)[0].ShortName)
		goto abort
	}
	*name = (*res.Levels)[0].Name
	return nil

abort:
	return fmt.Errorf("LevelName lookup failed: %s",
		err.Error())
}

// fetchFilter is a helper used in the ...IDByFoo functions
func fetchFilter(req proto.Request, path string) (*proto.Result, error) {
	var (
		err  error
		resp *resty.Response
	)
	res := &proto.Result{}

	if resp, err = PostReqBody(req, path); err != nil {
		// transport errors
		return nil, err
	}

	if err = decodeResponse(resp, res); err != nil {
		// http code errors
		return nil, err
	}

	if err = checkApplicationError(res); err != nil {
		return nil, err
	}
	return res, nil
}

// checkApplicationError tests the server result for
// application errors
func checkApplicationError(result *proto.Result) error {
	if result.StatusCode >= 300 {
		var s string
		// application errors
		if result.StatusCode == 404 {
			s = fmt.Sprintf("Object lookup error: %d - %s",
				result.StatusCode, result.StatusText)
		} else {
			s = fmt.Sprintf("Application error: %d - %s",
				result.StatusCode, result.StatusText)
		}
		m := []string{s}

		if result.Errors != nil {
			m = append(m, *result.Errors...)
		}

		return fmt.Errorf(combineStrings(m...))
	}
	return nil
}

// jobMetaIDByName implements the actual lookup of job metadata
// definitions from the server
func jobMetaIDByName(meta, name string) (string, error) {
	var (
		req  proto.Request
		path string
		err  error
	)
	res := &proto.Result{}

	switch meta {
	case `result`:
		req = proto.NewJobResultFilter()
		req.Filter.JobResult.Name = name
		path = `/search/jobResult/`
	case `status`:
		req = proto.NewJobStatusFilter()
		req.Filter.JobStatus.Name = name
		path = `/search/jobStatus/`
	case `type`:
		req = proto.NewJobTypeFilter()
		req.Filter.JobType.Name = name
		path = `/search/jobType/`
	default:
		err = fmt.Errorf("unknown meta object: %s", meta)
		goto abort
	}

	if res, err = fetchFilter(req, path); err != nil {
		goto abort
	}

	switch meta {
	case `result`:
		if res.JobResults == nil || len(*res.JobResults) == 0 {
			err = fmt.Errorf("no object returned for %s", meta)
			goto abort
		}
		if name != (*res.JobResults)[0].Name {
			err = fmt.Errorf(
				"name mismatch: %s vs %s",
				name, (*res.JobResults)[0].Name,
			)
			goto abort
		}
		return (*res.JobResults)[0].ID, nil
	case `status`:
		if res.JobStatus == nil || len(*res.JobStatus) == 0 {
			err = fmt.Errorf("no object returned for %s", meta)
			goto abort
		}
		if name != (*res.JobStatus)[0].Name {
			err = fmt.Errorf(
				"name mismatch: %s vs %s",
				name, (*res.JobStatus)[0].Name,
			)
			goto abort
		}
		return (*res.JobStatus)[0].ID, nil
	case `type`:
		if res.JobTypes == nil || len(*res.JobTypes) == 0 {
			err = fmt.Errorf("no object returned for %s", meta)
			goto abort
		}
		if name != (*res.JobTypes)[0].Name {
			err = fmt.Errorf(
				"name mismatch: %s vs %s",
				name, (*res.JobTypes)[0].Name,
			)
			goto abort
		}
		return (*res.JobTypes)[0].ID, nil
	default:
		err = fmt.Errorf("unknown meta object: %s", meta)
	}

abort:
	return ``, fmt.Errorf("JobMetadataID lookup failed: %s", err.Error())
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
