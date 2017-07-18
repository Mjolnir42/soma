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
	SectionEnvironment           = `environment`
	SectionInstanceMgmt          = `instance-mgmt`
	SectionJob                   = `job`
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
	SectionEntity                = `entity`
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
	SectionCheck          = `check`
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
	SectionMonitoring  = `monitoringsystem`
)

// Actions for the various permission sections
const (
	ActionAll            = `all`
	ActionAssign         = `assign`
	ActionAudit          = `audit`
	ActionCreate         = `create`
	ActionDeclare        = `declare`
	ActionDelete         = `delete`
	ActionGrant          = `grant`
	ActionList           = `list`
	ActionMemberAdd      = `member-add`
	ActionMemberList     = `member-list`
	ActionMemberRemove   = `member-remove`
	ActionPropertyAdd    = `property-add`
	ActionPropertyRemove = `property-remove`
	ActionRename         = `rename`
	ActionRepoRebuild    = `rebuild-repository`
	ActionRepoRestart    = `restart-repository`
	ActionRepoStop       = `stop-repository`
	ActionRetry          = `retry`
	ActionRevoke         = `revoke`
	ActionSearch         = `search`
	ActionSet            = `set`
	ActionShow           = `show`
	ActionShowConfig     = `show-config`
	ActionShutdown       = `shutdown`
	ActionSummary        = `summary`
	ActionSync           = `sync`
	ActionUnassign       = `unassign`
	ActionUpdate         = `update`
	ActionUse            = `use`
)

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
