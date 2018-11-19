/*-
 * Copyright (c) 2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/soma/internal/rest"

const (
	rtRepository                 = `/repository/`
	rtRepositoryID               = `/repository/:repositoryID`
	rtRepositoryInstance         = `/repository/:repositoryID/instance/`
	rtRepositoryInstanceID       = `/repository/:repositoryID/instance/:instanceID`
	rtRepositoryInstanceVersions = `/repository/:repositoryID/instance/:instanceID/versions`
	rtRepositoryMember           = `/repository/:repositoryID/member/`
	rtRepositoryMemberID         = `/repository/:repositoryID/member/:memberType/:memberID`
	rtRepositoryProperty         = `/repository/:repositoryID/property/`
	rtRepositoryPropertyID       = `/repository/:repositoryID/property/:propertyType/:sourceID`
	rtRepositoryPropertyMgmt     = `/repository/:repositoryID/property-mgmt/:propertyType/`
	rtRepositoryPropertyMgmtID   = `/repository/:repositoryID/property-mgmt/:propertyType/:propertyID`
	rtRepositoryTree             = `/repository/:repositoryID/tree`
	rtGlobalBucket               = `/bucket/`
	rtBucket                     = `/repository/:repositoryID/bucket/`
	rtBucketID                   = `/repository/:repositoryID/bucket/:bucketID`
	rtBucketInstance             = `/repository/:repositoryID/bucket/:bucketID/instance/`
	rtBucketInstanceID           = `/repository/:repositoryID/bucket/:bucketID/instance/:instanceID`
	rtBucketInstanceVersions     = `/repository/:repositoryID/bucket/:bucketID/instance/:instanceID/versions`
	rtBucketMember               = `/repository/:repositoryID/bucket/:bucketID/member/`
	rtBucketMemberID             = `/repository/:repositoryID/bucket/:bucketID/member/:memberType/:memberID`
	rtBucketProperty             = `/repository/:repositoryID/bucket/:bucketID/property/`
	rtBucketPropertyID           = `/repository/:repositoryID/bucket/:bucketID/property/:propertyType/:sourceID`
	rtBucketTree                 = `/repository/:repositoryID/bucket/:bucketID/tree`
	rtCluster                    = `/repository/:repositoryID/bucket/:bucketID/cluster/`
	rtClusterID                  = `/repository/:repositoryID/bucket/:bucketID/cluster/:clusterID`
	rtClusterInstance            = `/repository/:repositoryID/bucket/:bucketID/cluster/:clusterID/instance/`
	rtClusterInstanceID          = `/repository/:repositoryID/bucket/:bucketID/cluster/:clusterID/instance/:instanceID`
	rtClusterInstanceVersions    = `/repository/:repositoryID/bucket/:bucketID/cluster/:clusterID/instance/:instanceID/versions`
	rtClusterMember              = `/repository/:repositoryID/bucket/:bucketID/cluster/:clusterID/member/`
	rtClusterMemberID            = `/repository/:repositoryID/bucket/:bucketID/cluster/:clusterID/member/:memberType/:memberID`
	rtClusterProperty            = `/repository/:repositoryID/bucket/:bucketID/cluster/:clusterID/property/`
	rtClusterPropertyID          = `/repository/:repositoryID/bucket/:bucketID/cluster/:clusterID/property/:propertyType/:sourceID`
	rtClusterTree                = `/repository/:repositoryID/bucket/:bucketID/cluster/:clusterID/tree`
	rtGroup                      = `/repository/:repositoryID/bucket/:bucketID/group/`
	rtGroupID                    = `/repository/:repositoryID/bucket/:bucketID/group/:groupID`
	rtGroupInstance              = `/repository/:repositoryID/bucket/:bucketID/group/:groupID/instance/`
	rtGroupInstanceID            = `/repository/:repositoryID/bucket/:bucketID/group/:groupID/instance/:instanceID`
	rtGroupInstanceVersions      = `/repository/:repositoryID/bucket/:bucketID/group/:groupID/instance/:instanceID/versions`
	rtGroupMember                = `/repository/:repositoryID/bucket/:bucketID/group/:groupID/member/`
	rtGroupMemberID              = `/repository/:repositoryID/bucket/:bucketID/group/:groupID/member/:memberType/:memberID`
	rtGroupProperty              = `/repository/:repositoryID/bucket/:bucketID/group/:groupID/property/`
	rtGroupPropertyID            = `/repository/:repositoryID/bucket/:bucketID/group/:groupID/property/:propertyType/:sourceID`
	rtGroupTree                  = `/repository/:repositoryID/bucket/:bucketID/group/:groupID/tree`
	rtNode                       = `/node/`
	rtNodeID                     = `/node/:nodeID`
	rtNodeConfig                 = `/node/:nodeID/config`
	rtNodeUnassign               = `/repository/:repositoryID/bucket/:bucketID/node/:nodeID/config`
	rtNodeInstance               = `/repository/:repositoryID/bucket/:bucketID/node/:nodeID/instance/`
	rtNodeInstanceID             = `/repository/:repositoryID/bucket/:bucketID/node/:nodeID/instance/:instanceID`
	rtNodeInstanceVersions       = `/repository/:repositoryID/bucket/:bucketID/node/:nodeID/instance/:instanceID/versions`
	rtNodeProperty               = `/repository/:repositoryID/bucket/:bucketID/node/:nodeID/property/`
	rtNodePropertyID             = `/repository/:repositoryID/bucket/:bucketID/node/:nodeID/property/:propertyType/:sourceID`
	rtNodeTree                   = `/repository/:repositoryID/bucket/:bucketID/node/:nodeID/tree`
	rtPermission                 = `/category/:category/permission/`
	rtPermissionID               = `/category/:category/permission/:permissionID`
	rtRight                      = `/category/:category/permission/:permissionID/grant/`
	rtRightID                    = `/category/:category/permission/:permissionID/grant/:grantID`
	rtPropertyMgmt               = `/property-mgmt/:propertyType/`
	rtPropertyMgmtID             = `/property-mgmt/:propertyType/:propertyID`
	rtTeamMember                 = `/team/:teamID/member/`
	rtTeamRepositoryID           = `/team/:teamID/repository/:repositoryID`
	rtTeamRepositoryIDAudit      = `/team/:teamID/repository/:repositoryID/audit`
	rtTeamRepositoryIDName       = `/team/:teamID/repository/:repositoryID/name`
	rtTeamRepositoryIDOwner      = `/team/:teamID/repository/:repositoryID/owner`
	rtTeamPropertyMgmt           = `/team/:teamID/property-mgmt/:propertyType/`
	rtTeamPropertyMgmtID         = `/team/:teamID/property-mgmt/:propertyType/:propertyID`
	rtSearchRepository           = `/search/repository/`
	rtSearchBucket               = `/search/bucket/`
	rtSearchGroup                = `/search/repository/:repositoryID/bucket/:bucketID/group/`
	rtSearchCluster              = `/search/repository/:repositoryID/bucket/:bucketID/cluster/`
	rtSearchCustomProperty       = `/search/repository/:repositoryID/property-mgmt/:propertyType/`
	rtSearchNode                 = `/search/node/`
	rtSearchPermission           = `/search/permission/`
	rtSearchGlobalProperty       = `/search/property-mgmt/:propertyType/`
	rtSearchRight                = `/search/right/`
	rtSearchServiceProperty      = `/search/team/:teamID/property-mgmt/:propertyType/`
	rtSearchJob                  = `/search/job/`
	rtSearchJobType              = `/search/jobType/`
	rtSearchJobResult            = `/search/jobResult/`
	rtSearchJobStatus            = `/search/jobStatus/`
	rtSyncNode                   = `/sync/node/`
	rtDeployment                 = `/monitoringsystem/:monitoringID/deployment/`
	rtDeploymentID               = `/monitoringsystem/:monitoringID/deployment/id/:deploymentID`
	rtDeploymentIDAction         = `/monitoringsystem/:monitoringID/deployment/id/:deploymentID/:action`
	rtDeploymentState            = `/monitoringsystem/:monitoringID/deployment/state/`
	rtDeploymentStateID          = `/monitoringsystem/:monitoringID/deployment/state/:state`
	rtAliasDeploymentID          = `/deployment/id/:deploymentID`
	rtAliasDeploymentIDAction    = `/deployment/id/:deploymentID/:action`
	rtOncallMember               = `/oncall/:oncallID/member/`
	rtOncallMemberID             = `/oncall/:oncallID/member/:userID`
	rtJob                        = `/job/`
	rtJobEntry                   = `/job/byID/`
	rtJobEntryID                 = `/job/byID/:jobID`
	rtJobEntryWaitID             = `/job/byID/:jobID/_processed`
	rtJobTypeMgmt                = `/job/type-mgmt/`
	rtJobTypeMgmtID              = `/job/type-mgmt/:typeID`
	rtJobStatusMgmt              = `/job/status-mgmt/`
	rtJobStatusMgmtID            = `/job/status-mgmt/:statusID`
	rtJobResultMgmt              = `/job/result-mgmt/`
	rtJobResultMgmtID            = `/job/result-mgmt/:resultID`
)

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
