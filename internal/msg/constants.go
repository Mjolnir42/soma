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

// Grant categories mirror the regular categories but allow
// to grant the permission
const (
	CategoryGrantGlobal     = `global:grant`
	CategoryGrantIdentity   = `identity:grant`
	CategoryGrantSelf       = `self:grant`
	CategoryGrantOperation  = `operation:grant`
	CategoryGrantPermission = `permission:grant`
	CategoryGrantRepository = `repository:grant`
	CategoryGrantTeam       = `team:grant`
	CategoryGrantMonitoring = `monitoring:grant`
)

// Sections in category global are for actions with a global
// scope
const (
	CategoryGlobal          = `global`
	SectionAttribute        = `attribute`
	SectionDatacenter       = `datacenter`
	SectionEntity           = `entity`
	SectionEnvironment      = `environment`
	SectionHostDeployment   = `hostdeployment`
	SectionInstanceMgmt     = `instance-mgmt`
	SectionJobMgmt          = `job-mgmt`
	SectionLevel            = `level`
	SectionMetric           = `metric`
	SectionMode             = `mode`
	SectionMonitoringMgmt   = `monitoringsystem-mgmt`
	SectionNodeMgmt         = `node-mgmt`
	SectionOncall           = `oncall`
	SectionPredicate        = `predicate`
	SectionPropertyMgmt     = `property-mgmt`
	SectionPropertyNative   = `property-native`
	SectionPropertySystem   = `property-system`
	SectionPropertyTemplate = `property-template`
	SectionProvider         = `provider`
	SectionRepositoryMgmt   = `repository-mgmt`
	SectionServer           = `server`
	SectionState            = `state`
	SectionStatus           = `status`
	SectionUnit             = `unit`
	SectionValidity         = `validity`
	SectionView             = `view`
)

// Sections in category Identity are special global sections for actions
// related to identity management
const (
	CategoryIdentity = `identity`
	SectionTeamMgmt  = `team-mgmt`
	SectionUserMgmt  = `user-mgmt`
)

// Sections in category self are for actions with a per-user
// scope
const (
	CategorySelf = `self`
	SectionJob   = `job`
	SectionTeam  = `team`
	SectionUser  = `user`
)

// Sections in category operation are special global sections
// for actions to run the SOMA system
const (
	CategoryOperation = `operation`
	SectionSystem     = `system`
	SectionWorkflow   = `workflow`
)

// Sections in category permission are special global sections
// for actions on the permission system
const (
	CategoryPermission = `permission`
	SectionAction      = `action`
	SectionCategory    = `category`
	SectionPermission  = `permission`
	SectionRight       = `right`
	SectionSection     = `section`
)

// Sections in category repository are for actions with
// a per-repository scope
const (
	CategoryRepository      = `repository`
	SectionBucket           = `bucket`
	SectionCheckConfig      = `check-config`
	SectionCluster          = `cluster`
	SectionGroup            = `group`
	SectionInstance         = `instance`
	SectionNodeConfig       = `node-config`
	SectionPropertyCustom   = `property-custom`
	SectionRepositoryConfig = `repository-config`
)

// Sections in category team are for actions with a per-team
// scope
const (
	CategoryTeam           = `team`
	SectionNode            = `node`
	SectionPropertyService = `property-service`
	SectionRepository      = `repository`
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
	ActionAdd             = `add`
	ActionAll             = `all`
	ActionAssemble        = `assemble`
	ActionAssign          = `assign`
	ActionAudit           = `audit`
	ActionCreate          = `create`
	ActionDeclare         = `declare`
	ActionDelete          = `delete`
	ActionDestroy         = `destroy`
	ActionFailed          = `failed`
	ActionFilter          = `filter`
	ActionGet             = `get`
	ActionGrant           = `grant`
	ActionInsertNullID    = `insert-null`
	ActionList            = `list`
	ActionMap             = `map`
	ActionMemberAssign    = `member-assign`
	ActionMemberList      = `member-list`
	ActionMemberUnassign  = `member-unassign`
	ActionPending         = `pending`
	ActionPropertyCreate  = `property-create`
	ActionPropertyDestroy = `property-destroy`
	ActionPropertyUpdate  = `property-update`
	ActionPurge           = `purge`
	ActionRemove          = `remove`
	ActionRename          = `rename`
	ActionRepoRebuild     = `rebuild-repository`
	ActionRepoRestart     = `restart-repository`
	ActionRepoStop        = `stop-repository`
	ActionRetry           = `retry`
	ActionRevoke          = `revoke`
	ActionSearch          = `search`
	ActionSearchAll       = `search/all`
	ActionSearchByList    = `search/list`
	ActionSearchByName    = `search/name`
	ActionSet             = `set`
	ActionShow            = `show`
	ActionShowConfig      = `show-config`
	ActionShutdown        = `shutdown`
	ActionSuccess         = `success`
	ActionSummary         = `summary`
	ActionSync            = `sync`
	ActionTree            = `tree`
	ActionUnassign        = `unassign`
	ActionUnmap           = `unmap`
	ActionUpdate          = `update`
	ActionUse             = `use`
	ActionVersions        = `versions`
	ActionWait            = `wait`
)

// Section supervisor handles AAA requests outside the permission
// model
const (
	SectionSupervisor     = `supervisor`
	ActionActivate        = `activate`
	ActionAuthenticate    = `authenticate`
	ActionAuthorize       = `authorize`
	ActionCacheUpdate     = `cacheupdate`
	ActionDeactivate      = `deactivate`
	ActionGC              = `gc`
	ActionKex             = `kex`
	ActionPassword        = `password`
	ActionToken           = `token`
	TaskBasicAuth         = `basic-auth`
	TaskChange            = `change`
	TaskInvalidate        = `invalidate`
	TaskInvalidateAccount = `invalidate-account`
	TaskInvalidateGlobal  = `invalidate-global`
	TaskNone              = `none`
	TaskRequest           = `request`
	TaskReset             = `reset`
	TaskUser              = `user`
)

// Entity types
const (
	EntityRepository = `repository`
	EntityBucket     = `bucket`
	EntityGroup      = `group`
	EntityCluster    = `cluster`
	EntityNode       = `node`
	EntityMonitoring = `monitoring`
	EntityTeam       = `team`
	InvalidObjectID  = `ffffffff-ffff-3fff-ffff-ffffffffffff`
)

// Subject types
const (
	SubjectRoot  = `root`
	SubjectAdmin = `admin`
	SubjectUser  = `user`
	SubjectTool  = `tool`
	SubjectTeam  = `team`
)

// Property types
const (
	PropertyCustom   = `custom`
	PropertyNative   = `native`
	PropertyService  = `service`
	PropertySystem   = `system`
	PropertyTemplate = `template`
	PropertyOncall   = `oncall`
)

// Constraint Types
const (
	ConstraintCustom    = `custom`
	ConstraintNative    = `native`
	ConstraintService   = `service`
	ConstraintSystem    = `system`
	ConstraintOncall    = `oncall`
	ConstraintAttribute = `attribute`
)

const (
	SystemPropertyDisableAllMonitoring      = `disable_all_monitoring`
	SystemPropertyDisableCheckConfiguration = `disable_check_configuration`
	ViewLocal                               = `local`
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
