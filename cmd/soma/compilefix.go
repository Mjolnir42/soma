package main

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/tree"
	"github.com/mjolnir42/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

type treeRequest struct {
	RequestType string
	Action      string
	User        string
	AuthUser    string
	JobId       uuid.UUID
	reply       chan somaResult
	Repository  somaRepositoryRequest
	Bucket      somaBucketRequest
	Group       somaGroupRequest
	Cluster     somaClusterRequest
	Node        somaNodeRequest
	CheckConfig somaCheckConfigRequest
}

type somaNodeRequest struct {
	action string
	user   string
	Node   proto.Node
	reply  chan somaResult
}

type somaLevelResult struct {
	ResultError error
	Level       proto.Level
}

type somaCapabilityResult struct {
	ResultError error
	Capability  proto.Capability
}

type somaDatacenterResult struct {
	ResultError error
	Datacenter  proto.Datacenter
}

type somaMetricResult struct {
	ResultError error
	Metric      proto.Metric
}

type somaModeResult struct {
	ResultError error
	Mode        proto.Mode
}

type somaNodeResult struct {
	ResultError error
	Node        proto.Node
}

type somaOncallResult struct {
	ResultError error
	Oncall      proto.Oncall
}

type somaPredicateResult struct {
	ResultError error
	Predicate   proto.Predicate
}

type somaPropertyResult struct {
	ResultError error
	prType      string
	System      proto.PropertySystem
	Native      proto.PropertyNative
	Service     proto.PropertyService
	Custom      proto.PropertyCustom
}

type somaProviderResult struct {
	ResultError error
	Provider    proto.Provider
}

type somaServerResult struct {
	ResultError error
	Server      proto.Server
}

type somaStatusResult struct {
	ResultError error
	Status      proto.Status
}

type somaTeamResult struct {
	ResultError error
	Team        proto.Team
}

type somaUnitResult struct {
	ResultError error
	Unit        proto.Unit
}

type somaUserResult struct {
	ResultError error
	User        proto.User
}

type somaValidityResult struct {
	ResultError error
	Validity    proto.Validity
}

type somaViewResult struct {
	ResultError error
	View        proto.View
}

type somaCapabilityRequest struct {
	action     string
	Capability proto.Capability
	reply      chan somaResult
}

type somaCapabilityReadHandler struct {
	input     chan somaCapabilityRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

type somaCapabilityWriteHandler struct {
	input    chan somaCapabilityRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

type somaDatacenterRequest struct {
	action     string
	Datacenter proto.Datacenter
	rename     string
	reply      chan somaResult
}

type somaDatacenterWriteHandler struct {
	input    chan somaDatacenterRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	ren_stmt *sql.Stmt
	grp_add  *sql.Stmt
	grp_del  *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

type somaDatacenterReadHandler struct {
	input     chan somaDatacenterRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	grp_list  *sql.Stmt
	grp_show  *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

type somaLevelRequest struct {
	action string
	Level  proto.Level
	reply  chan somaResult
}

type somaLevelWriteHandler struct {
	input    chan somaLevelRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

type somaLevelReadHandler struct {
	input     chan somaLevelRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

type jobDelay struct {
	input    chan waitSpec
	shutdown chan bool
	notify   chan string
	waitList map[string][]waitSpec
	jobDone  map[string]time.Time
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

type waitSpec struct {
	JobId string
	RecvT time.Time
	Reply chan bool
}

type somaMetricWriteHandler struct {
	input        chan somaMetricRequest
	shutdown     chan bool
	conn         *sql.DB
	add_stmt     *sql.Stmt
	del_stmt     *sql.Stmt
	pkg_add_stmt *sql.Stmt
	pkg_del_stmt *sql.Stmt
	appLog       *log.Logger
	reqLog       *log.Logger
	errLog       *log.Logger
}

type somaMetricReadHandler struct {
	input     chan somaMetricRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

type somaMetricRequest struct {
	action string
	Metric proto.Metric
	reply  chan somaResult
}

type somaModeWriteHandler struct {
	input    chan somaModeRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

type somaModeReadHandler struct {
	input     chan somaModeRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

type somaModeRequest struct {
	action string
	Mode   proto.Mode
	reply  chan somaResult
}

type somaNodeReadHandler struct {
	input     chan somaNodeRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	conf_stmt *sql.Stmt
	sync_stmt *sql.Stmt
	ponc_stmt *sql.Stmt
	psvc_stmt *sql.Stmt
	psys_stmt *sql.Stmt
	pcst_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

type somaOncallRequest struct {
	action string
	Oncall proto.Oncall
	reply  chan somaResult
}

type somaOncallReadHandler struct {
	input     chan somaOncallRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

type somaOncallWriteHandler struct {
	input    chan somaOncallRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	upd_stmt *sql.Stmt
	del_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

type somaPredicateRequest struct {
	action    string
	Predicate proto.Predicate
	reply     chan somaResult
}

type somaPredicateReadHandler struct {
	input     chan somaPredicateRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

type somaPredicateWriteHandler struct {
	input    chan somaPredicateRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

type somaPropertyRequest struct {
	action  string
	prType  string
	System  proto.PropertySystem
	Native  proto.PropertyNative
	Service proto.PropertyService
	Custom  proto.PropertyCustom
	reply   chan somaResult
}

type somaPropertyReadHandler struct {
	input         chan somaPropertyRequest
	shutdown      chan bool
	conn          *sql.DB
	list_sys_stmt *sql.Stmt
	list_srv_stmt *sql.Stmt
	list_nat_stmt *sql.Stmt
	list_tpl_stmt *sql.Stmt
	list_cst_stmt *sql.Stmt
	show_sys_stmt *sql.Stmt
	show_srv_stmt *sql.Stmt
	show_nat_stmt *sql.Stmt
	show_tpl_stmt *sql.Stmt
	show_cst_stmt *sql.Stmt
	appLog        *log.Logger
	reqLog        *log.Logger
	errLog        *log.Logger
}

type somaPropertyWriteHandler struct {
	input             chan somaPropertyRequest
	shutdown          chan bool
	conn              *sql.DB
	add_sys_stmt      *sql.Stmt
	add_nat_stmt      *sql.Stmt
	add_cst_stmt      *sql.Stmt
	add_srv_stmt      *sql.Stmt
	add_tpl_stmt      *sql.Stmt
	add_srv_attr_stmt *sql.Stmt
	add_tpl_attr_stmt *sql.Stmt
	del_sys_stmt      *sql.Stmt
	del_nat_stmt      *sql.Stmt
	del_cst_stmt      *sql.Stmt
	del_srv_stmt      *sql.Stmt
	del_tpl_stmt      *sql.Stmt
	del_srv_attr_stmt *sql.Stmt
	del_tpl_attr_stmt *sql.Stmt
	appLog            *log.Logger
	reqLog            *log.Logger
	errLog            *log.Logger
}

type somaProviderRequest struct {
	action   string
	Provider proto.Provider
	reply    chan somaResult
}

type somaProviderReadHandler struct {
	input     chan somaProviderRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

type somaProviderWriteHandler struct {
	input    chan somaProviderRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

type somaServerRequest struct {
	action string
	Server proto.Server
	Filter proto.Filter
	reply  chan somaResult
}

type somaServerReadHandler struct {
	input     chan somaServerRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	sync_stmt *sql.Stmt
	snam_stmt *sql.Stmt
	sass_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

type somaServerWriteHandler struct {
	input    chan somaServerRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	prg_stmt *sql.Stmt
	upd_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

type somaStatusRequest struct {
	action string
	Status proto.Status
	reply  chan somaResult
}

type somaStatusReadHandler struct {
	input     chan somaStatusRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

type somaStatusWriteHandler struct {
	input    chan somaStatusRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

type somaTeamRequest struct {
	action string
	Team   proto.Team
	reply  chan somaResult
}

type somaTeamReadHandler struct {
	input     chan somaTeamRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	sync_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

type somaTeamWriteHandler struct {
	input    chan somaTeamRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	upd_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

type somaUnitRequest struct {
	action string
	Unit   proto.Unit
	reply  chan somaResult
}

type somaUnitReadHandler struct {
	input     chan somaUnitRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

type somaUnitWriteHandler struct {
	input    chan somaUnitRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

type somaUserRequest struct {
	action string
	User   proto.User
	reply  chan somaResult
}

type somaUserReadHandler struct {
	input     chan somaUserRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	sync_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

type somaUserWriteHandler struct {
	input    chan somaUserRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	prg_stmt *sql.Stmt
	upd_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

type somaValidityRequest struct {
	action   string
	Validity proto.Validity
	reply    chan somaResult
}

type somaValidityReadHandler struct {
	input     chan somaValidityRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

type somaValidityWriteHandler struct {
	input    chan somaValidityRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

type somaViewRequest struct {
	action string
	name   string
	View   proto.View
	reply  chan somaResult
}

type somaViewReadHandler struct {
	input     chan somaViewRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

type somaViewWriteHandler struct {
	input    chan somaViewRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	ren_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

func (a *somaNodeResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Nodes = append(r.Nodes, somaNodeResult{ResultError: err})
	}
}

func (a *somaNodeResult) SomaAppendResult(r *somaResult) {
	r.Nodes = append(r.Nodes, *a)
}

type treeKeeper struct {
	repoId               string
	repoName             string
	team                 string
	rbLevel              string
	broken               bool
	ready                bool
	stopped              bool
	frozen               bool
	rebuild              bool
	input                chan treeRequest
	shutdown             chan bool
	stopchan             chan bool
	errChan              chan *tree.Error
	actionChan           chan *tree.Action
	conn                 *sql.DB
	tree                 *tree.Tree
	get_view             *sql.Stmt
	start_job            *sql.Stmt
	stmt_CapMonMetric    *sql.Stmt
	stmt_Check           *sql.Stmt
	stmt_CheckConfig     *sql.Stmt
	stmt_CheckInstance   *sql.Stmt
	stmt_Cluster         *sql.Stmt
	stmt_ClusterCustProp *sql.Stmt
	stmt_ClusterOncall   *sql.Stmt
	stmt_ClusterService  *sql.Stmt
	stmt_ClusterSysProp  *sql.Stmt
	stmt_DefaultDC       *sql.Stmt
	stmt_DelDuplicate    *sql.Stmt
	stmt_GetComputed     *sql.Stmt
	stmt_GetPrevious     *sql.Stmt
	stmt_Group           *sql.Stmt
	stmt_GroupCustProp   *sql.Stmt
	stmt_GroupOncall     *sql.Stmt
	stmt_GroupService    *sql.Stmt
	stmt_GroupSysProp    *sql.Stmt
	stmt_List            *sql.Stmt
	stmt_Node            *sql.Stmt
	stmt_NodeCustProp    *sql.Stmt
	stmt_NodeOncall      *sql.Stmt
	stmt_NodeService     *sql.Stmt
	stmt_NodeSysProp     *sql.Stmt
	stmt_Pkgs            *sql.Stmt
	stmt_Team            *sql.Stmt
	stmt_Threshold       *sql.Stmt
	stmt_Update          *sql.Stmt
	appLog               *log.Logger
	log                  *log.Logger
	startLog             *log.Logger
}

func (tk *treeKeeper) isReady() bool {
	return tk.ready
}

func (tk *treeKeeper) isBroken() bool {
	return tk.broken
}

func (tk *treeKeeper) isStopped() bool {
	return tk.stopped
}

func (tk *treeKeeper) run() {
}

type somaBucketRequest struct {
	action string
	Bucket proto.Bucket
	reply  chan somaResult
}

type somaBucketResult struct {
	ResultError error
	Bucket      proto.Bucket
}

type somaCheckConfigRequest struct {
	action      string
	CheckConfig proto.CheckConfig
	reply       chan somaResult
}

type somaCheckConfigResult struct {
	ResultError error
	CheckConfig proto.CheckConfig
}

type somaClusterRequest struct {
	action  string
	Cluster proto.Cluster
	reply   chan somaResult
}

type somaClusterResult struct {
	ResultError error
	Cluster     proto.Cluster
}

type somaGroupRequest struct {
	action string
	Group  proto.Group
	reply  chan somaResult
}

type somaGroupResult struct {
	ResultError error
	Group       proto.Group
}

type somaRepositoryRequest struct {
	action     string
	remoteAddr string
	user       string
	rbLevel    string
	rebuild    bool
	Repository proto.Repository
	reply      chan somaResult
}

type somaRepositoryResult struct {
	ResultError error
	Repository  proto.Repository
}

type somaRepositoryReadHandler struct {
	input     chan somaRepositoryRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	ponc_stmt *sql.Stmt
	psvc_stmt *sql.Stmt
	psys_stmt *sql.Stmt
	pcst_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaRepositoryReadHandler) run() {
}

type somaBucketReadHandler struct {
	input     chan somaBucketRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	ponc_stmt *sql.Stmt
	psvc_stmt *sql.Stmt
	psys_stmt *sql.Stmt
	pcst_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaBucketReadHandler) run() {
}

type somaClusterReadHandler struct {
	input     chan somaClusterRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	mbnl_stmt *sql.Stmt
	ponc_stmt *sql.Stmt
	psvc_stmt *sql.Stmt
	psys_stmt *sql.Stmt
	pcst_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaClusterReadHandler) run() {
}

type somaCheckConfigurationReadHandler struct {
	input                 chan somaCheckConfigRequest
	shutdown              chan bool
	conn                  *sql.DB
	list_stmt             *sql.Stmt
	show_base             *sql.Stmt
	show_threshold        *sql.Stmt
	show_constr_custom    *sql.Stmt
	show_constr_system    *sql.Stmt
	show_constr_native    *sql.Stmt
	show_constr_service   *sql.Stmt
	show_constr_attribute *sql.Stmt
	show_constr_oncall    *sql.Stmt
	show_instance_info    *sql.Stmt
	appLog                *log.Logger
	reqLog                *log.Logger
	errLog                *log.Logger
}

func (r *somaCheckConfigurationReadHandler) run() {
}

type somaGroupReadHandler struct {
	input     chan somaGroupRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	mbgl_stmt *sql.Stmt
	mbcl_stmt *sql.Stmt
	mbnl_stmt *sql.Stmt
	ponc_stmt *sql.Stmt
	psvc_stmt *sql.Stmt
	psys_stmt *sql.Stmt
	pcst_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaGroupReadHandler) run() {
}

type guidePost struct {
	input              chan treeRequest
	system             chan msg.Request
	shutdown           chan bool
	conn               *sql.DB
	jbsv_stmt          *sql.Stmt
	repo_stmt          *sql.Stmt
	name_stmt          *sql.Stmt
	node_stmt          *sql.Stmt
	serv_stmt          *sql.Stmt
	attr_stmt          *sql.Stmt
	cthr_stmt          *sql.Stmt
	cdel_stmt          *sql.Stmt
	bucket_for_node    *sql.Stmt
	bucket_for_cluster *sql.Stmt
	bucket_for_group   *sql.Stmt
	appLog             *log.Logger
	reqLog             *log.Logger
	errLog             *log.Logger
}

func (g *guidePost) run() {
}

type forestCustodian struct {
	input     chan somaRepositoryRequest
	system    chan msg.Request
	shutdown  chan bool
	conn      *sql.DB
	add_stmt  *sql.Stmt
	load_stmt *sql.Stmt
	name_stmt *sql.Stmt
	rbck_stmt *sql.Stmt
	rbci_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (f *forestCustodian) run() {
}

type grimReaper struct {
	system chan msg.Request
	conn   *sql.DB
	appLog *log.Logger
	reqLog *log.Logger
	errLog *log.Logger
}

func (grim *grimReaper) run() {
}

type somaDeploymentRequest struct {
	action     string
	Deployment string
	reply      chan somaResult
}

type somaDeploymentResult struct {
	ResultError error
	ListEntry   string
	Deployment  proto.Deployment
}

type somaDeploymentHandler struct {
	input    chan somaDeploymentRequest
	shutdown chan bool
	conn     *sql.DB
	get_stmt *sql.Stmt
	upd_stmt *sql.Stmt
	sta_stmt *sql.Stmt
	act_stmt *sql.Stmt
	lst_stmt *sql.Stmt
	all_stmt *sql.Stmt
	clr_stmt *sql.Stmt
	dpr_stmt *sql.Stmt
	sty_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

func (w *somaDeploymentHandler) run() {
}

type somaHostDeploymentRequest struct {
	action  string
	system  string
	assetid int64
	idlist  []string
	reply   chan somaResult
}

type somaHostDeploymentResult struct {
	ResultError error
	Delete      bool
	DeleteId    string
	Deployment  proto.Deployment
}

type somaHostDeploymentHandler struct {
	input     chan somaHostDeploymentRequest
	shutdown  chan bool
	conn      *sql.DB
	geti_stmt *sql.Stmt
	last_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (self *somaHostDeploymentHandler) run() {
}

type workflowWrite struct {
	input                      chan msg.Request
	shutdown                   chan bool
	conn                       *sql.DB
	stmtRetryDeployment        *sql.Stmt
	stmtTriggerAvailableUpdate *sql.Stmt
	stmtSet                    *sql.Stmt
	appLog                     *log.Logger
	reqLog                     *log.Logger
	errLog                     *log.Logger
}

type workflowRead struct {
	input       chan msg.Request
	shutdown    chan bool
	conn        *sql.DB
	stmtSummary *sql.Stmt
	stmtList    *sql.Stmt
	appLog      *log.Logger
	reqLog      *log.Logger
	errLog      *log.Logger
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
