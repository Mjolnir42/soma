/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package msg

// Privileged access permission categories
const (
	CategoryOmnipotence = `omnipotence`
	CategorySystem      = `system`
)

// Sections in category global are for actions with a global
// scope
const (
	CategoryGlobal               = `global`
	SectionAttribute             = `attribute`
	SectionDatacenter            = `datacenter`
	SectionEntity                = `entity`
	SectionEnvironment           = `environment`
	SectionHostDeployment        = `hostdeployment`
	SectionInstanceMgmt          = `instance-mgmt`
	SectionJob                   = `job`
	SectionJobMgmt               = `job-mgmt`
	SectionLevel                 = `level`
	SectionMetric                = `metric`
	SectionMode                  = `mode`
	SectionMonitoringMgmt        = `monitoringsystem-mgmt`
	SectionNodeMgmt              = `node-mgmt`
	SectionOncall                = `oncall`
	SectionPredicate             = `predicate`
	SectionPropertyNative        = `property-native`
	SectionPropertyServiceGlobal = `property-service-global`
	SectionPropertySystem        = `property-system`
	SectionProvider              = `provider`
	SectionRepositoryMgmt        = `repository-mgmt`
	SectionServer                = `server`
	SectionState                 = `state`
	SectionStatus                = `status`
	SectionTeam                  = `team`
	SectionUnit                  = `unit`
	SectionUser                  = `user`
	SectionValidity              = `validity`
	SectionView                  = `view`
)

// Sections in category operation are special global sections
// for actions to run the SOMA system
const (
	CategoryOperation      = `operation`
	SectionSystemOperation = `system-operation`
	SectionWorkflow        = `workflow`
)

// Sections in category permission are special global sections
// for actions on the permission system
const (
	CategoryPermission = `permission`
	SectionCategory    = `category`
	SectionPermission  = `permission`
)

// Sections in category repository are for actions with
// a per-repository scope
const (
	CategoryRepository    = `repository`
	SectionBucket         = `bucket`
	SectionCheckConfig    = `check-config`
	SectionCluster        = `cluster`
	SectionGroup          = `group`
	SectionInstance       = `instance`
	SectionNodeConfig     = `node-config`
	SectionPropertyCustom = `property-custom`
	SectionRepository     = `repository`
)

// Sections in category team are for actions with a per-team
// scope
const (
	CategoryTeam               = `team`
	SectionNode                = `node`
	SectionPropertyServiceTeam = `property-service-team`
)

// Sections in category monitoring are for actions with a
// per-monitoringsystem scope
const (
	CategoryMonitoring = `monitoring`
	SectionCapability  = `capability`
	SectionDeployment  = `deployment`
	SectionMonitoring  = `monitoringsystem`
)

// Actions for the various permission sections
const (
	ActionAll            = `all`
	ActionAssemble       = `assemble`
	ActionAssign         = `assign`
	ActionAudit          = `audit`
	ActionCreate         = `create`
	ActionDeclare        = `declare`
	ActionDelete         = `delete`
	ActionGet            = `get`
	ActionGrant          = `grant`
	ActionInsertNullID   = `insert-null`
	ActionFailed         = `failed`
	ActionList           = `list`
	ActionMemberAdd      = `member-add`
	ActionMemberList     = `member-list`
	ActionMemberRemove   = `member-remove`
	ActionPropertyAdd    = `property-add`
	ActionPropertyRemove = `property-remove`
	ActionPurge          = `purge`
	ActionRename         = `rename`
	ActionRepoRebuild    = `rebuild-repository`
	ActionRepoRestart    = `restart-repository`
	ActionRepoStop       = `stop-repository`
	ActionRetry          = `retry`
	ActionRevoke         = `revoke`
	ActionSearch         = `search`
	ActionSearchAll      = `search/all`
	ActionSearchByAsset  = `search/asset`
	ActionSearchByList   = `search/list`
	ActionSearchByName   = `search/name`
	ActionSet            = `set`
	ActionShow           = `show`
	ActionShowConfig     = `show-config`
	ActionShutdown       = `shutdown`
	ActionSuccess        = `success`
	ActionSummary        = `summary`
	ActionSync           = `sync`
	ActionUnassign       = `unassign`
	ActionUpdate         = `update`
	ActionUse            = `use`
	ActionVersions       = `versions`
)

// Entity types
const (
	EntityRepository = `repository`
	EntityBucket     = `bucket`
	EntityGroup      = `group`
	EntityCluster    = `cluster`
	EntityNode       = `node`
)

const (
	// RFC3339Milli is a format string for millisecond precision RFC3339
	RFC3339Milli string = "2006-01-02T15:04:05.000Z07:00"
	// LogStrReq is a format string for logging requests (deprecated)
	LogStrReq = `Subsystem=%s, Request=%s, User=%s, Addr=%s`
	// LogStrSRq is a format string for logging requests
	LogStrSRq = `Section=%s, Action=%s, User=%s, Addr=%s`
	// LogStrArg is a format string for logging scoped requests
	LogStrArg = `Subsystem=%s, Request=%s, User=%s, Addr=%s, Arg=%s`
	// LogStrOK is a format string for logging OK results
	LogStrOK = `Section=%s, Action=%s, InternalCode=%d, ExternalCode=%d`
	// LogStrErr is a format string for logging ERROR results
	LogStrErr = `Section=%s, Action=%s, InternalCode=%d, Error=%s`
)

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
