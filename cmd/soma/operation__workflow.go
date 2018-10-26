/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main // import "github.com/mjolnir42/soma/cmd/soma"

import (
	"fmt"
	"net/url"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/internal/help"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerWorkflow(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:        `workflow`,
				Usage:       `SUBCOMMANDS for workflow inquiry`,
				Description: help.Text(`workflow::`),
				Subcommands: []cli.Command{
					{
						Name:        `summary`,
						Usage:       `Show summary of workflow status`,
						Description: help.Text(`workflow::summary`),
						Action:      runtime(workflowSummary),
					},
					{
						Name:        `list`,
						Usage:       `List instances in all workflow states`,
						Description: help.Text(`workflow::list`),
						Action:      runtime(workflowList),
					},
					{
						Name:        `search`,
						Usage:       `Search for instances in a specific workflow state`,
						Description: help.Text(`workflow::search`),
						Action:      runtime(workflowSearch),
					},
					{
						Name:        `retry`,
						Usage:       `Reschedule an instance in a failed state`,
						Description: help.Text(`workflow::retry`),
						Action:      runtime(workflowRetry),
					},
					{
						Name:         `set`,
						Usage:        `Hard-set an instance's worflow status`,
						Description:  help.Text(`workflow::set`),
						Action:       runtime(workflowSet),
						BashComplete: cmpl.WorkflowSet,
						Flags: []cli.Flag{
							cli.BoolFlag{
								Name:  `force, f`,
								Usage: `Force is required to break the workflow`,
							},
						},
					},
				},
			},
		}...,
	)
	return &app
}

// workflowSummary function
// soma workflow summary
func workflowSummary(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/workflow/summary`, `list`, nil, c)
}

// workflowList function
// soma workflow list
func workflowList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/workflow/`, `list`, nil, c)
}

// workflowSearch function
// soma workflow search ${status}
func workflowSearch(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	if err := adm.ValidateStatus(c.Args().First()); err != nil {
		return err
	}

	req := proto.NewWorkflowFilter()
	req.Filter.Workflow.Status = c.Args().First()

	return adm.Perform(`postbody`, `/search/workflow/`, `list`, req, c)
}

// workflowRetry function XXX UNTESTED
// soma workflow retry ${instanceID}
func workflowRetry(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	if err := adm.ValidateInstance(c.Args().First()); err != nil {
		return err
	}

	req := proto.NewWorkflowRequest()
	req.Workflow.InstanceID = c.Args().First()

	path := `/workflow/retry`
	return adm.Perform(`patchbody`, path, `command`, req, c)
}

// workflowSet function XXX UNTESTED
// soma workflow set ${instanceConfigID} status ${current} next ${nextStatus}
func workflowSet(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`status`, `next`}
	mandatoryOptions := []string{`status`, `next`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	if !adm.IsUUID(c.Args().First()) {
		return fmt.Errorf("Argument is not a UUID: %s",
			c.Args().First())
	}
	if err := adm.ValidateStatus(opts[`status`][0]); err != nil {
		return err
	}
	if err := adm.ValidateStatus(opts[`next`][0]); err != nil {
		return err
	}
	req := proto.NewWorkflowRequest()
	req.Flags.Forced = c.Bool(`force`)
	req.Workflow.InstanceConfigID = c.Args().First()
	req.Workflow.Status = opts[`status`][0]
	req.Workflow.NextStatus = opts[`next`][0]

	path := fmt.Sprintf(
		"/workflow/set/%s",
		url.QueryEscape(c.Args().First()),
	)
	return adm.Perform(`patchbody`, path, `command`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
