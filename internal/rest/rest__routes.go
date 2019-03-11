/*-
 * Copyright (c) 2017-2018, Jörg Pernfuß
 * Copyright (c) 2018, 1&1 IONOS SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/soma/internal/rest"

import (
	"github.com/julienschmidt/httprouter"
)

// setupRouter returns a configured httprouter
func (x *Rest) setupRouter() *httprouter.Router {
	router := httprouter.New()

	router.HEAD(`/`, x.Unauthenticated(x.Ping))

	router.GET(`/attribute/:attribute`, x.Authenticated(x.AttributeShow))
	router.GET(`/attribute/`, x.Authenticated(x.AttributeList))
	router.GET(`/capability/:capabilityID`, x.Authenticated(x.CapabilityShow))
	router.GET(`/capability/`, x.Authenticated(x.CapabilityList))
	router.GET(`/category/:category/section/:sectionID/action/:actionID`, x.Authenticated(x.ActionShow))
	router.GET(`/category/:category/section/:sectionID/action/`, x.Authenticated(x.ActionList))
	router.GET(`/category/:category/section/:sectionID`, x.Authenticated(x.SectionShow))
	router.GET(`/category/:category/section/`, x.Authenticated(x.SectionList))
	router.GET(`/category/:category`, x.Authenticated(x.CategoryShow))
	router.GET(`/category/`, x.Authenticated(x.CategoryList))
	router.GET(`/checkconfig/:repositoryID/:checkID`, x.Authenticated(x.CheckConfigShow))
	router.GET(`/checkconfig/:repositoryID/`, x.Authenticated(x.CheckConfigList))
	router.GET(`/datacenter/:datacenter`, x.Authenticated(x.DatacenterShow))
	router.GET(`/datacenter/`, x.Authenticated(x.DatacenterList))
	router.GET(`/entity/:entity`, x.Authenticated(x.EntityShow))
	router.GET(`/entity/`, x.Authenticated(x.EntityList))
	router.GET(`/environment/:environment`, x.Authenticated(x.EnvironmentShow))
	router.GET(`/environment/`, x.Authenticated(x.EnvironmentList))
	router.GET(`/hostdeployment/:monitoringID/:assetID`, x.Unauthenticated(x.HostDeploymentFetch))
	router.GET(`/instance/:instanceID/versions`, x.Authenticated(x.InstanceVersions))
	router.GET(`/instance/:instanceID`, x.Authenticated(x.ScopeSelectInstanceShow))
	router.GET(`/instance/`, x.Authenticated(x.ScopeSelectInstanceList))
	router.GET(`/level/:level`, x.Authenticated(x.LevelShow))
	router.GET(`/level/`, x.Authenticated(x.LevelList))
	router.GET(`/metric/:metric`, x.Authenticated(x.MetricShow))
	router.GET(`/metric/`, x.Authenticated(x.MetricList))
	router.GET(`/mode/:mode`, x.Authenticated(x.ModeShow))
	router.GET(`/mode/`, x.Authenticated(x.ModeList))
	router.GET(`/monitoringsystem/:monitoringID`, x.Authenticated(x.MonitoringShow))
	router.GET(`/monitoringsystem/`, x.Authenticated(x.ScopeSelectMonitoringList))
	router.GET(`/oncall/:oncallID`, x.Authenticated(x.OncallShow))
	router.GET(`/oncall/`, x.Authenticated(x.OncallList))
	router.GET(`/predicate/:predicate`, x.Authenticated(x.PredicateShow))
	router.GET(`/predicate/`, x.Authenticated(x.PredicateList))
	router.GET(`/provider/:provider`, x.Authenticated(x.ProviderShow))
	router.GET(`/provider/`, x.Authenticated(x.ProviderList))
	router.GET(`/server/:serverID`, x.Authenticated(x.ServerShow))
	router.GET(`/server/`, x.Authenticated(x.ServerList))
	router.GET(`/state/:state`, x.Authenticated(x.StateShow))
	router.GET(`/state/`, x.Authenticated(x.StateList))
	router.GET(`/status/:status`, x.Authenticated(x.StatusShow))
	router.GET(`/status/`, x.Authenticated(x.StatusList))
	router.GET(`/sync/datacenter/`, x.Authenticated(x.DatacenterSync))
	router.GET(`/sync/server/`, x.Authenticated(x.ServerSync))
	router.GET(`/sync/team/`, x.Authenticated(x.TeamMgmtSync))
	router.GET(`/sync/user/`, x.Authenticated(x.UserMgmtSync))
	router.GET(`/team/:teamID`, x.Authenticated(x.ScopeSelectTeamShow))
	router.GET(`/team/`, x.Authenticated(x.TeamMgmtList))
	router.GET(`/unit/:unit`, x.Authenticated(x.UnitShow))
	router.GET(`/unit/`, x.Authenticated(x.UnitList))
	router.GET(`/user/:userID`, x.Authenticated(x.ScopeSelectUserShow))
	router.GET(`/user/`, x.Authenticated(x.UserMgmtList))
	router.GET(`/validity/:property`, x.Authenticated(x.ValidityShow))
	router.GET(`/validity/`, x.Authenticated(x.ValidityList))
	router.GET(`/view/:view`, x.Authenticated(x.ViewShow))
	router.GET(`/view/`, x.Authenticated(x.ViewList))
	router.GET(`/workflow/`, x.Authenticated(x.WorkflowList))
	router.GET(`/workflow/summary`, x.Authenticated(x.WorkflowSummary))
	router.GET(rtBucket, x.Authenticated(x.BucketList))
	router.GET(rtBucketID, x.Authenticated(x.BucketShow))
	router.GET(rtBucketInstance, x.Authenticated(x.InstanceList))
	router.GET(rtBucketInstanceID, x.Authenticated(x.InstanceShow))
	router.GET(rtBucketInstanceVersions, x.Authenticated(x.InstanceVersions))
	router.GET(rtBucketMember, x.Authenticated(x.BucketMemberList))
	router.GET(rtBucketTree, x.Authenticated(x.BucketTree))
	router.GET(rtCluster, x.Authenticated(x.ClusterList))
	router.GET(rtClusterID, x.Authenticated(x.ClusterShow))
	router.GET(rtClusterInstance, x.Authenticated(x.InstanceList))
	router.GET(rtClusterInstanceID, x.Authenticated(x.InstanceShow))
	router.GET(rtClusterInstanceVersions, x.Authenticated(x.InstanceVersions))
	router.GET(rtClusterMember, x.Authenticated(x.ClusterMemberList))
	router.GET(rtClusterTree, x.Authenticated(x.ClusterTree))
	router.GET(rtGroup, x.Authenticated(x.GroupList))
	router.GET(rtGroupID, x.Authenticated(x.GroupShow))
	router.GET(rtGroupInstance, x.Authenticated(x.InstanceList))
	router.GET(rtGroupInstanceID, x.Authenticated(x.InstanceShow))
	router.GET(rtGroupInstanceVersions, x.Authenticated(x.InstanceVersions))
	router.GET(rtGroupMember, x.Authenticated(x.GroupMemberList))
	router.GET(rtGroupTree, x.Authenticated(x.GroupTree))
	router.GET(rtJob, x.Authenticated(x.ScopeSelectJobList))
	router.GET(rtJobEntry, x.Authenticated(x.ScopeSelectJobList))
	router.GET(rtJobEntryID, x.Authenticated(x.JobShow))
	router.GET(rtJobResultMgmt, x.Authenticated(x.JobResultMgmtList))
	router.GET(rtJobResultMgmtID, x.Authenticated(x.JobResultMgmtShow))
	router.GET(rtJobStatusMgmt, x.Authenticated(x.JobStatusMgmtList))
	router.GET(rtJobStatusMgmtID, x.Authenticated(x.JobStatusMgmtShow))
	router.GET(rtJobTypeMgmt, x.Authenticated(x.JobTypeMgmtList))
	router.GET(rtJobTypeMgmtID, x.Authenticated(x.JobTypeMgmtShow))
	router.GET(rtNode, x.Authenticated(x.NodeList))
	router.GET(rtNodeConfig, x.Authenticated(x.NodeShowConfig))
	router.GET(rtNodeID, x.Authenticated(x.NodeShow))
	router.GET(rtNodeInstance, x.Authenticated(x.InstanceList))
	router.GET(rtNodeInstanceID, x.Authenticated(x.InstanceShow))
	router.GET(rtNodeInstanceVersions, x.Authenticated(x.InstanceVersions))
	router.GET(rtNodeTree, x.Authenticated(x.NodeConfigTree))
	router.GET(rtOncallMember, x.Authenticated(x.OncallMemberList))
	router.GET(rtPermission, x.Authenticated(x.PermissionList))
	router.GET(rtPermissionID, x.Authenticated(x.PermissionShow))
	router.GET(rtPropertyMgmt, x.Authenticated(x.PropertyMgmtList))
	router.GET(rtPropertyMgmtID, x.Authenticated(x.PropertyMgmtShow))
	router.GET(rtRepository, x.Authenticated(x.RepositoryConfigList))
	router.GET(rtTeamRepositoryID, x.Authenticated(x.ScopeSelectRepositoryShow))
	router.GET(rtRepositoryInstance, x.Authenticated(x.InstanceList))
	router.GET(rtRepositoryInstanceID, x.Authenticated(x.InstanceShow))
	router.GET(rtRepositoryInstanceVersions, x.Authenticated(x.InstanceVersions))
	router.GET(rtRepositoryPropertyMgmt, x.Authenticated(x.PropertyMgmtList))
	router.GET(rtRepositoryPropertyMgmtID, x.Authenticated(x.PropertyMgmtShow))
	router.GET(rtRepositoryTree, x.Authenticated(x.RepositoryConfigTree))
	router.GET(rtRight, x.Authenticated(x.RightList))
	router.GET(rtRightID, x.Authenticated(x.RightShow))
	router.GET(rtSyncNode, x.Authenticated(x.NodeMgmtSync))
	router.GET(rtTeamMember, x.Authenticated(x.TeamMgmtMemberList))
	router.GET(rtTeamPropertyMgmt, x.Authenticated(x.PropertyMgmtList))
	router.GET(rtTeamPropertyMgmtID, x.Authenticated(x.PropertyMgmtShow))
	router.HEAD(`/authenticate/validate`, x.Authenticated(x.SupervisorValidate))
	router.POST(`/hostdeployment/:monitoringID/:assetID`, x.Unauthenticated(x.HostDeploymentAssemble))
	router.POST(`/search/action/`, x.Authenticated(x.ActionSearch))
	router.POST(`/search/capability/`, x.Authenticated(x.CapabilitySearch))
	router.POST(`/search/checkconfig/:repositoryID/`, x.Authenticated(x.CheckConfigSearch))
	router.POST(`/search/level/`, x.Authenticated(x.LevelSearch))
	router.POST(`/search/monitoringsystem/`, x.Authenticated(x.ScopeSelectMonitoringSearch))
	router.POST(`/search/oncall/`, x.Authenticated(x.OncallSearch))
	router.POST(`/search/section/`, x.Authenticated(x.SectionSearch))
	router.POST(`/search/server/`, x.Authenticated(x.ServerSearch))
	router.POST(`/search/team/`, x.Authenticated(x.ScopeSelectTeamSearch))
	router.POST(`/search/user/`, x.Authenticated(x.ScopeSelectUserSearch))
	router.POST(`/search/workflow/`, x.Authenticated(x.WorkflowSearch))
	router.POST(rtSearchBucket, x.Authenticated(x.BucketSearch))
	router.POST(rtSearchCluster, x.Authenticated(x.ClusterSearch))
	router.POST(rtSearchCustomProperty, x.Authenticated(x.PropertyMgmtSearch))
	router.POST(rtSearchGlobalProperty, x.Authenticated(x.PropertyMgmtSearch))
	router.POST(rtSearchGroup, x.Authenticated(x.GroupSearch))
	router.POST(rtSearchJob, x.Authenticated(x.JobSearch))
	router.POST(rtSearchJobType, x.Authenticated(x.JobTypeMgmtSearch))
	router.POST(rtSearchJobResult, x.Authenticated(x.JobResultMgmtSearch))
	router.POST(rtSearchJobStatus, x.Authenticated(x.JobStatusMgmtSearch))
	router.POST(rtSearchNode, x.Authenticated(x.NodeSearch))
	router.POST(rtSearchPermission, x.Authenticated(x.PermissionSearch))
	router.POST(rtSearchRepository, x.Authenticated(x.ScopeSelectRepositorySearch))
	router.POST(rtSearchRight, x.Authenticated(x.RightSearch))
	router.POST(rtSearchServiceProperty, x.Authenticated(x.PropertyMgmtSearch))

	if !x.conf.ReadOnly {
		if !x.conf.Observer {
			router.DELETE(`/accounts/tokens/:account`, x.Authenticated(x.SupervisorTokenInvalidateAccount))
			router.DELETE(`/attribute/:attribute`, x.Authenticated(x.AttributeRemove))
			router.DELETE(`/capability/:capabilityID`, x.Authenticated(x.CapabilityRevoke))
			router.DELETE(`/category/:category/section/:sectionID/action/:actionID`, x.Authenticated(x.ActionRemove))
			router.DELETE(`/category/:category/section/:sectionID`, x.Authenticated(x.SectionRemove))
			router.DELETE(`/category/:category`, x.Authenticated(x.CategoryRemove))
			router.DELETE(`/checkconfig/:repositoryID/:checkID`, x.Authenticated(x.CheckConfigDestroy))
			router.DELETE(`/datacenter/:datacenter`, x.Authenticated(x.DatacenterRemove))
			router.DELETE(`/entity/:entity`, x.Authenticated(x.EntityRemove))
			router.DELETE(`/environment/:environment`, x.Authenticated(x.EnvironmentRemove))
			router.DELETE(`/level/:level`, x.Authenticated(x.LevelRemove))
			router.DELETE(`/metric/:metric`, x.Authenticated(x.MetricRemove))
			router.DELETE(`/mode/:mode`, x.Authenticated(x.ModeRemove))
			router.DELETE(`/monitoringsystem/:monitoringID`, x.Authenticated(x.MonitoringMgmtRemove))
			router.DELETE(`/oncall/:oncallID`, x.Authenticated(x.OncallRemove))
			router.DELETE(`/predicate/:predicate`, x.Authenticated(x.PredicateRemove))
			router.DELETE(`/provider/:provider`, x.Authenticated(x.ProviderRemove))
			router.DELETE(`/server/:serverID`, x.Authenticated(x.ServerRemove))
			router.DELETE(`/state/:state`, x.Authenticated(x.StateRemove))
			router.DELETE(`/status/:status`, x.Authenticated(x.StatusRemove))
			router.DELETE(`/team/:teamID`, x.Authenticated(x.TeamMgmtRemove))
			router.DELETE(`/tokens/global`, x.Authenticated(x.SupervisorTokenInvalidateGlobal))
			router.DELETE(`/tokens/self/active`, x.Authenticated(x.SupervisorTokenInvalidate))
			router.DELETE(`/tokens/self/all`, x.Authenticated(x.SupervisorTokenInvalidateSelf))
			router.DELETE(`/unit/:unit`, x.Authenticated(x.UnitRemove))
			router.DELETE(`/user/:userID`, x.Authenticated(x.UserMgmtRemove))
			router.DELETE(`/validity/:property`, x.Authenticated(x.ValidityRemove))
			router.DELETE(`/view/:view`, x.Authenticated(x.ViewRemove))
			router.DELETE(rtBucketID, x.Authenticated(x.BucketDestroy))
			router.DELETE(rtBucketMemberID, x.Authenticated(x.BucketMemberUnassign))
			router.DELETE(rtBucketPropertyID, x.Authenticated(x.BucketPropertyDestroy))
			router.DELETE(rtClusterID, x.Authenticated(x.ClusterDestroy))
			router.DELETE(rtClusterMemberID, x.Authenticated(x.ClusterMemberUnassign))
			router.DELETE(rtClusterPropertyID, x.Authenticated(x.ClusterPropertyDestroy))
			router.DELETE(rtGroupID, x.Authenticated(x.GroupDestroy))
			router.DELETE(rtGroupMemberID, x.Authenticated(x.GroupMemberUnassign))
			router.DELETE(rtGroupPropertyID, x.Authenticated(x.GroupPropertyDestroy))
			router.DELETE(rtJobResultMgmtID, x.Authenticated(x.JobResultMgmtRemove))
			router.DELETE(rtJobStatusMgmtID, x.Authenticated(x.JobStatusMgmtRemove))
			router.DELETE(rtJobTypeMgmtID, x.Authenticated(x.JobTypeMgmtRemove))
			router.DELETE(rtNode, x.Authenticated(x.NodeMgmtRemove))
			router.DELETE(rtNodeID, x.Authenticated(x.NodeMgmtRemove))
			router.DELETE(rtNodePropertyID, x.Authenticated(x.NodeConfigPropertyDestroy))
			router.DELETE(rtNodeUnassign, x.Authenticated(x.NodeConfigUnassign))
			router.DELETE(rtOncallMemberID, x.Authenticated(x.OncallMemberUnassign))
			router.DELETE(rtPermissionID, x.Authenticated(x.PermissionRemove))
			router.DELETE(rtPropertyMgmtID, x.Authenticated(x.PropertyMgmtRemove))
			router.DELETE(rtRepositoryPropertyID, x.Authenticated(x.RepositoryConfigPropertyDestroy))
			router.DELETE(rtRepositoryPropertyMgmtID, x.Authenticated(x.PropertyMgmtCustomRemove))
			router.DELETE(rtRightID, x.Authenticated(x.RightRevoke))
			router.DELETE(rtTeamPropertyMgmtID, x.Authenticated(x.PropertyMgmtServiceRemove))
			router.DELETE(rtTeamRepositoryID, x.Authenticated(x.RepositoryDestroy))
			router.GET(rtAliasDeploymentID, x.Unauthenticated(x.DeploymentShow))
			router.GET(rtCompatDeploymentID, x.Unauthenticated(x.DeploymentShow))
			router.GET(rtDeployment, x.Unauthenticated(x.DeploymentList))
			router.GET(rtDeploymentID, x.Unauthenticated(x.DeploymentShow))
			router.GET(rtDeploymentState, x.Unauthenticated(x.DeploymentPending))
			router.GET(rtDeploymentStateID, x.Unauthenticated(x.DeploymentFilter))
			router.GET(rtJobEntryWaitID, x.Authenticated(x.ScopeSelectJobWait))
			router.GET(rtTeamRepositoryIDAudit, x.Authenticated(x.RepositoryAudit))
			router.PATCH(`/accounts/password/:kexID`, x.Unauthenticated(x.SupervisorPasswordChange))
			router.PATCH(`/oncall/:oncallID`, x.Authenticated(x.OncallUpdate))
			router.PATCH(`/workflow/retry`, x.Authenticated(x.WorkflowRetry))
			router.PATCH(`/workflow/set/:instanceconfigID`, x.Authenticated(x.WorkflowSet))
			router.PATCH(rtAliasDeploymentIDAction, x.Unauthenticated(x.DeploymentUpdate))
			router.PATCH(rtCompatDeploymentIDAction, x.Unauthenticated(x.DeploymentUpdate))
			router.PATCH(rtClusterID, x.Authenticated(x.ClusterRename))
			router.PATCH(rtDeploymentIDAction, x.Unauthenticated(x.DeploymentUpdate))
			router.PATCH(rtOncallMember, x.Authenticated(x.OncallMemberAssign))
			router.PATCH(rtPermissionID, x.Authenticated(x.PermissionEdit))
			router.PATCH(rtTeamRepositoryIDName, x.Authenticated(x.RepositoryRename))
			router.PATCH(rtTeamRepositoryIDOwner, x.Authenticated(x.RepositoryRepossess))
			router.POST(`/attribute/`, x.Authenticated(x.AttributeAdd))
			router.POST(`/capability/`, x.Authenticated(x.CapabilityDeclare))
			router.POST(`/category/:category/section/:sectionID/action/`, x.Authenticated(x.ActionAdd))
			router.POST(`/category/:category/section/`, x.Authenticated(x.SectionAdd))
			router.POST(`/category/`, x.Authenticated(x.CategoryAdd))
			router.POST(`/checkconfig/:repositoryID/`, x.Authenticated(x.CheckConfigCreate))
			router.POST(`/datacenter/`, x.Authenticated(x.DatacenterAdd))
			router.POST(`/entity/`, x.Authenticated(x.EntityAdd))
			router.POST(`/environment/`, x.Authenticated(x.EnvironmentAdd))
			router.POST(`/kex/`, x.Unauthenticated(x.SupervisorKex))
			router.POST(`/level/`, x.Authenticated(x.LevelAdd))
			router.POST(`/metric/`, x.Authenticated(x.MetricAdd))
			router.POST(`/mode/`, x.Authenticated(x.ModeAdd))
			router.POST(`/monitoringsystem/`, x.Authenticated(x.MonitoringMgmtAdd))
			router.POST(`/oncall/`, x.Authenticated(x.OncallAdd))
			router.POST(`/predicate/`, x.Authenticated(x.PredicateAdd))
			router.POST(`/provider/`, x.Authenticated(x.ProviderAdd))
			router.POST(`/server/:serverID`, x.Authenticated(x.ServerAddNull))
			router.POST(`/server/`, x.Authenticated(x.ServerAdd))
			router.POST(`/state/`, x.Authenticated(x.StateAdd))
			router.POST(`/status/`, x.Authenticated(x.StatusAdd))
			router.POST(`/system/`, x.Authenticated(x.SystemOperation))
			router.POST(`/team/`, x.Authenticated(x.TeamMgmtAdd))
			router.POST(`/unit/`, x.Authenticated(x.UnitAdd))
			router.POST(`/user/`, x.Authenticated(x.UserMgmtAdd))
			router.POST(`/validity/`, x.Authenticated(x.ValidityAdd))
			router.POST(`/view/`, x.Authenticated(x.ViewAdd))
			router.POST(rtBucket, x.Authenticated(x.BucketCreate))
			router.POST(rtBucketMember, x.Authenticated(x.BucketMemberAssign))
			router.POST(rtBucketProperty, x.Authenticated(x.BucketPropertyCreate))
			router.POST(rtCluster, x.Authenticated(x.ClusterCreate))
			router.POST(rtClusterMember, x.Authenticated(x.ClusterMemberAssign))
			router.POST(rtClusterProperty, x.Authenticated(x.ClusterPropertyCreate))
			router.POST(rtGroup, x.Authenticated(x.GroupCreate))
			router.POST(rtGroupMemberAssign, x.Authenticated(x.GroupMemberAssign))
			router.POST(rtGroupProperty, x.Authenticated(x.GroupPropertyCreate))
			router.POST(rtJobResultMgmt, x.Authenticated(x.JobResultMgmtAdd))
			router.POST(rtJobStatusMgmt, x.Authenticated(x.JobStatusMgmtAdd))
			router.POST(rtJobTypeMgmt, x.Authenticated(x.JobTypeMgmtAdd))
			router.POST(rtNode, x.Authenticated(x.NodeMgmtAdd))
			router.POST(rtNodeProperty, x.Authenticated(x.NodeConfigPropertyCreate))
			router.POST(rtPermission, x.Authenticated(x.PermissionAdd))
			router.POST(rtPropertyMgmt, x.Authenticated(x.PropertyMgmtAdd))
			router.POST(rtRepository, x.Authenticated(x.RepositoryMgmtCreate))
			router.POST(rtRepositoryProperty, x.Authenticated(x.RepositoryConfigPropertyCreate))
			router.POST(rtRepositoryPropertyMgmt, x.Authenticated(x.PropertyMgmtCustomAdd))
			router.POST(rtRight, x.Authenticated(x.RightGrant))
			router.POST(rtTeamPropertyMgmt, x.Authenticated(x.PropertyMgmtServiceAdd))
			router.PUT(`/accounts/activate/root/:kexID`, x.Unauthenticated(x.SupervisorActivateRoot))
			router.PUT(`/accounts/activate/user/:kexID`, x.Unauthenticated(x.SupervisorActivateUser))
			router.PUT(`/accounts/password/:kexID`, x.Unauthenticated(x.SupervisorPasswordReset))
			router.PUT(`/datacenter/:datacenter`, x.Authenticated(x.DatacenterRename))
			router.PUT(`/entity/:entity`, x.Authenticated(x.EntityRename))
			router.PUT(`/environment/:environment`, x.Authenticated(x.EnvironmentRename))
			router.PUT(`/server/:serverID`, x.Authenticated(x.ServerUpdate))
			router.PUT(`/state/:state`, x.Authenticated(x.StateRename))
			router.PUT(`/team/:teamID`, x.Authenticated(x.TeamMgmtUpdate))
			router.PUT(`/tokens/request/:kexID`, x.Unauthenticated(x.SupervisorTokenRequest))
			router.PUT(`/user/:userID`, x.Authenticated(x.UserMgmtUpdate))
			router.PUT(`/view/:view`, x.Authenticated(x.ViewRename))
			router.PUT(rtBucketPropertyID, x.Authenticated(x.BucketPropertyUpdate))
			router.PUT(rtClusterPropertyID, x.Authenticated(x.ClusterPropertyUpdate))
			router.PUT(rtGroupPropertyID, x.Authenticated(x.GroupPropertyUpdate))
			router.PUT(rtNodeConfig, x.Authenticated(x.NodeConfigAssign))
			router.PUT(rtNodeID, x.Authenticated(x.NodeMgmtUpdate))
			router.PUT(rtNodePropertyID, x.Authenticated(x.NodeConfigPropertyUpdate))
			router.PUT(rtRepositoryPropertyID, x.Authenticated(x.RepositoryConfigPropertyUpdate))
		}
	}
	return router
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
