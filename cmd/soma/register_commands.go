package main

import (
	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/help"
)

func registerCommands(app cli.App) *cli.App {

	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "Initialize local client files",
			Action: cmdClientInit,
		},
		{
			Name:        `login`,
			Usage:       `Authenticate with the SOMA middleware`,
			Description: help.Text(`supervisor::login`),
			Action:      runtime(supervisorLogin),
		},
		{
			Name:        `logout`,
			Usage:       `Revoke currently used password token`,
			Action:      runtime(supervisorLogout),
			Description: help.Text(`supervisor::logout`),
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  `all, a`,
					Usage: `Revoke all active tokens for the account`,
				},
			},
		},
		{
			Name:   "experiment",
			Usage:  "Test cli.Action functionality",
			Action: runtime(cmdExperiment),
			Hidden: true,
		},
	}

	app = *registerAction(app)
	app = *registerAttributes(app)
	app = *registerBuckets(app)
	app = *registerCapability(app)
	app = *registerCategories(app)
	app = *registerChecks(app)
	app = *registerClusters(app)
	app = *registerDatacenters(app)
	app = *registerEntities(app)
	app = *registerEnvironments(app)
	app = *registerGroups(app)
	app = *registerInstanceMgmt(app)
	app = *registerInstances(app)
	app = *registerJobs(app)
	app = *registerLevels(app)
	app = *registerMetrics(app)
	app = *registerModes(app)
	app = *registerMonitoringMgmt(app)
	app = *registerNodes(app)
	app = *registerOncall(app)
	app = *registerPermissions(app)
	app = *registerPredicates(app)
	app = *registerProperty(app)
	app = *registerProviders(app)
	app = *registerRepositoryMgmt(app)
	app = *registerRepository(app)
	app = *registerRights(app)
	app = *registerSection(app)
	app = *registerServers(app)
	app = *registerStates(app)
	app = *registerStatus(app)
	app = *registerTeams(app)
	app = *registerUnits(app)
	app = *registerUserMgmt(app)
	app = *registerValidity(app)
	app = *registerViews(app)
	app = *registerWorkflow(app)
	app = *registerOps(app)

	return &app
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
