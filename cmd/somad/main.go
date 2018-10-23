package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	"github.com/client9/reopen"
	"github.com/mjolnir42/soma/internal/config"
	"github.com/mjolnir42/soma/internal/handler"
	"github.com/mjolnir42/soma/internal/rest"
	"github.com/mjolnir42/soma/internal/soma"
	"github.com/mjolnir42/soma/internal/super"
	metrics "github.com/rcrowley/go-metrics"
)

// global variables
var (
	// main database connection pool
	conn *sql.DB
	// config file runtime configuration
	SomaCfg config.Config
	// lookup table of logfile handles for logrotate reopen
	logFileMap *soma.LogHandleMap
	// Global metrics registry
	Metrics = make(map[string]metrics.Registry)
	// version string set at compile time
	somaVersion string
)

func init() {
	logrus.SetOutput(os.Stderr)
	logFileMap = soma.NewLogHandleMap()
}

func main() {
	var (
		configFlag, configFile, obsRepoFlag         string
		noPokeFlag, forcedCorruption, versionFlag   bool
		err                                         error
		appLog, reqLog, errLog, auditLog            *logrus.Logger
		lfhGlobal, lfhApp, lfhReq, lfhErr, lfhAudit *reopen.FileWriter
		app                                         *soma.Soma
		hm                                          *handler.Map
		rst                                         *rest.Rest
		lvl                                         logrus.Level
	)

	// Daemon command line flags
	flag.StringVar(&configFlag, "config", "/srv/soma/huxley/conf/soma.conf", "Configuration file location")
	flag.StringVar(&obsRepoFlag, "repo", "", "Single-repository mode target repository")
	flag.BoolVar(&noPokeFlag, "nopoke", false, "Disable lifecycle pokes")
	flag.BoolVar(&forcedCorruption, `allowdatacorruption`, false, `Allow single-repo mode on production`)
	flag.BoolVar(&versionFlag, `version`, false, `Print version information`)
	flag.Parse()

	if versionFlag {
		version() // exit(0)
	}

	logrus.Printf("Starting runtime config initialization, SOMA v%s", somaVersion)
	/*
	 * Read configuration file
	 */
	if configFile, err = filepath.Abs(configFlag); err != nil {
		logrus.Fatal(err)
	}
	if configFile, err = filepath.EvalSymlinks(configFile); err != nil {
		logrus.Fatal(err)
	}
	err = SomaCfg.ReadConfigFile(configFile)
	if err != nil {
		logrus.Fatal(err)
	}
	SomaCfg.Version = somaVersion

	switch SomaCfg.LogLevel {
	case `debug`:
		lvl = logrus.DebugLevel
	case `info`:
		lvl = logrus.InfoLevel
	case `warn`:
		lvl = logrus.WarnLevel
	case `error`:
		lvl = logrus.ErrorLevel
	case `fatal`:
		lvl = logrus.FatalLevel
	case `panic`:
		lvl = logrus.PanicLevel
	}

	// Open logfiles
	if lfhGlobal, err = reopen.NewFileWriter(
		filepath.Join(SomaCfg.LogPath, `global.log`),
	); err != nil {
		logrus.Fatalf("Unable to open global output log: %s", err)
	}
	logrus.SetOutput(lfhGlobal)
	logrus.SetLevel(lvl)
	logFileMap.Add(`global`, lfhGlobal, logrus.StandardLogger())

	appLog = logrus.New()
	if lfhApp, err = reopen.NewFileWriter(
		filepath.Join(SomaCfg.LogPath, `application.log`),
	); err != nil {
		logrus.Fatalf("Unable to open application log: %s", err)
	}
	appLog.SetLevel(lvl)
	appLog.Out = lfhApp
	appLog.Formatter = &logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	}
	logFileMap.Add(`application`, lfhApp, appLog)

	reqLog = logrus.New()
	if lfhReq, err = reopen.NewFileWriter(
		filepath.Join(SomaCfg.LogPath, `request.log`),
	); err != nil {
		logrus.Fatalf("Unable to open request log: %s", err)
	}
	reqLog.SetLevel(lvl)
	reqLog.Out = lfhReq
	reqLog.Formatter = &logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	}
	logFileMap.Add(`request`, lfhReq, reqLog)

	errLog = logrus.New()
	if lfhErr, err = reopen.NewFileWriter(
		filepath.Join(SomaCfg.LogPath, `error.log`),
	); err != nil {
		logrus.Fatalf("Unable to open error log: %s", err)
	}
	// it is called errorlog for a reason
	errLog.SetLevel(logrus.ErrorLevel)
	errLog.Out = lfhErr
	errLog.Formatter = &logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	}
	logFileMap.Add(`error`, lfhErr, errLog)

	auditLog = logrus.New()
	if lfhAudit, err = reopen.NewFileWriter(
		filepath.Join(SomaCfg.LogPath, `audit.log`),
	); err != nil {
		logrus.Fatalf("Unable to open audit log: %s", err)
	}
	auditLog.SetLevel(lvl)
	auditLog.Out = lfhAudit
	auditLog.Formatter = &logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	}
	logFileMap.Add(`audit`, lfhAudit, auditLog)

	hm = handler.NewMap()

	// signal handler will reopen all logfiles on USR2
	sigChanLogRotate := make(chan os.Signal, 1)
	signal.Notify(sigChanLogRotate, syscall.SIGUSR2)
	go logrotate(sigChanLogRotate, hm)

	// print selected runtime mode
	if SomaCfg.ReadOnly {
		appLog.Println(`Instance has been configured as: read-only mode`)
	} else if SomaCfg.Observer {
		appLog.Println(`Instance has been configured as: observer mode`)
	} else {
		appLog.Println(`Instance has been configured as: normal mode`)
	}

	// single-repo cli argument overwrites config file
	if obsRepoFlag != `` {
		SomaCfg.ObserverRepo = obsRepoFlag
	}
	if SomaCfg.ObserverRepo != `` {
		appLog.Printf("Single-repository mode active for: %s", SomaCfg.ObserverRepo)
	}

	// disallow single-repository mode on production r/w instances
	if !SomaCfg.ReadOnly && !SomaCfg.Observer &&
		SomaCfg.ObserverRepo != `` && SomaCfg.Environment == `production` &&
		!forcedCorruption {
		errLog.Fatal(`Single-repository r/w mode disallowed for production environments. ` +
			`Use the -allowdatacorruption flag if you are sure this will be the only ` +
			`running SOMA instance.`)
	}

	if noPokeFlag {
		SomaCfg.NoPoke = true
		appLog.Println(`Instance has disabled outgoing pokes by lifeCycle manager`)
	}

	/*
	 * Register metrics collections
	 */
	Metrics[`golang`] = metrics.NewPrefixedRegistry(`golang.`)
	metrics.RegisterRuntimeMemStats(Metrics[`golang`])
	go metrics.CaptureRuntimeMemStats(Metrics[`golang`], time.Second*60)

	Metrics[`soma`] = metrics.NewPrefixedRegistry(`soma`)
	Metrics[`soma`].Register(`requests.latency`,
		// TODO NewCustomTimer(Histogram, Meter) so there is access
		// to Histogram.Clear()
		metrics.NewTimer())
	soma.Metrics = Metrics
	rest.Metrics = Metrics

	/*
	 * Construct listen address
	 */
	SomaCfg.Daemon.URL = &url.URL{}
	SomaCfg.Daemon.URL.Host = fmt.Sprintf("%s:%s", SomaCfg.Daemon.Listen, SomaCfg.Daemon.Port)
	if SomaCfg.Daemon.TLS {
		SomaCfg.Daemon.URL.Scheme = "https"
		if ok, pt := govalidator.IsFilePath(SomaCfg.Daemon.Cert); !ok {
			errLog.Fatal("Missing required certificate configuration config/daemon/cert-file")
		} else {
			if pt != govalidator.Unix {
				errLog.Fatal("config/daemon/cert-File: valid Windows paths are not helpful")
			}
		}
		if ok, pt := govalidator.IsFilePath(SomaCfg.Daemon.Key); !ok {
			errLog.Fatal("Missing required key configuration config/daemon/key-file")
		} else {
			if pt != govalidator.Unix {
				errLog.Fatal("config/daemon/key-file: valid Windows paths are not helpful")
			}
		}
	} else {
		SomaCfg.Daemon.URL.Scheme = "http"
	}

	connectToDatabase(appLog, errLog)
	go pingDatabase(errLog)

	app = soma.New(hm, logFileMap, conn, &SomaCfg, appLog, reqLog, errLog, auditLog)
	app.Start()

	rst = rest.New(super.IsAuthorized, hm, &SomaCfg, reqLog, errLog)

	go rst.Run()

	// signal handler for shutdown
	sigChanShutdown := make(chan os.Signal, 1)
	signal.Notify(sigChanShutdown, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	appLog.Println(`somad server process started, waiting for shutdown signal`)
	<-sigChanShutdown
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
