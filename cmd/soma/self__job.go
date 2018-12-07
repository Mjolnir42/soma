/*-
 * Copyright (c) 2016-2018, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main // import "github.com/mjolnir42/soma/cmd/soma"

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/internal/help"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerJobs(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:        `job`,
				Usage:       `SUBCOMMANDS for job information`,
				Description: help.Text(`job::`),
				Subcommands: []cli.Command{
					{
						Name:        `update`,
						Usage:       `Check and update status of outstanding locally cached jobs`,
						Description: help.Text(`job::update`),
						Action:      runtime(clientlocalJobUpdate),
						Flags: []cli.Flag{
							cli.BoolFlag{
								Name:  "verbose, v",
								Usage: "Include full raw job request (admin only)",
							},
						},
					},
					{
						Name:        `show`,
						Usage:       `Show details about a job`,
						Description: help.Text(`job::show`),
						Action:      runtime(jobShow),
						Flags: []cli.Flag{
							cli.BoolFlag{
								Name:  "verbose, v",
								Usage: "Include full raw job request (admin only)",
							},
						},
					},
					{
						Name:        `wait`,
						Usage:       `Block until a job has completed`,
						Description: help.Text(`job::wait`),
						Action:      runtime(jobWait),
					},
					{
						Name:        `list`,
						Usage:       `SUBCOMMANDS for listing job information`,
						Description: help.Text(`job::list`),
						Subcommands: []cli.Command{
							{
								Name:        `outstanding`,
								Usage:       `List outstanding jobs from local cache DB`,
								Description: help.Text(`job::list`),
								Action:      runtime(clientlocalJobListOutstanding),
							},
							{
								Name:        `local`,
								Usage:       `List all jobs from local cache DB`,
								Description: help.Text(`job::list`),
								Action:      runtime(clientlocalJobListLocal),
							},
							{
								Name:        `remote`,
								Usage:       `List all jobs from server`,
								Description: help.Text(`job::list`),
								Action:      runtime(jobList),
							},
						},
					},
					{
						Name:        `prune`,
						Usage:       `Delete completed jobs from local cache`,
						Description: help.Text(`job::prune`),
						Action:      runtime(clientlocalJobPruneDB),
					},
					//					{
					//						Name:   `search`,
					//						Usage:  `Search for a job`,
					//						Action: runtime(jobSearch),
					//					},
					{
						Name:        `type-mgmt`,
						Usage:       `SUBCOMMANDS for job type management`,
						Description: help.Text(`job-type-mgmt::`),
						Subcommands: []cli.Command{
							{
								Name:        `add`,
								Usage:       `Add a job type`,
								Description: help.Text(`job-type-mgmt::add`),
								Action:      runtime(jobTypeMgmtAdd),
							},
							{
								Name:        `remove`,
								Usage:       `Remove an unused job type`,
								Description: help.Text(`job-type-mgmt::remove`),
								Action:      runtime(jobTypeMgmtRemove),
							},
							{
								Name:        `show`,
								Usage:       `Show details about a job type`,
								Description: help.Text(`job-type-mgmt::show`),
								Action:      runtime(jobTypeMgmtShow),
							},
							{
								Name:        `list`,
								Usage:       `List all defined job types`,
								Description: help.Text(`job-type-mgmt::list`),
								Action:      runtime(jobTypeMgmtList),
							},
							{
								Name:         `search`,
								Usage:        `Search for a defined job types`,
								Description:  help.Text(`job-type-mgmt::search`),
								Action:       runtime(jobTypeMgmtSearch),
								BashComplete: cmpl.DirectIdName,
							},
						},
					},
					{
						Name:        `result-mgmt`,
						Usage:       `SUBCOMMANDS for job result management`,
						Description: help.Text(`job-result-mgmt::`),
						Subcommands: []cli.Command{
							{
								Name:        `add`,
								Usage:       `Add a job result`,
								Description: help.Text(`job-result-mgmt::add`),
								Action:      runtime(jobResultMgmtAdd),
							},
							{
								Name:        `remove`,
								Usage:       `Remove an unused job result`,
								Description: help.Text(`job-result-mgmt::remove`),
								Action:      runtime(jobResultMgmtRemove),
							},
							{
								Name:        `show`,
								Usage:       `Show details about a job result`,
								Description: help.Text(`job-result-mgmt::show`),
								Action:      runtime(jobResultMgmtShow),
							},
							{
								Name:        `list`,
								Usage:       `List all defined job results`,
								Description: help.Text(`job-result-mgmt::list`),
								Action:      runtime(jobResultMgmtList),
							},
							{
								Name:         `search`,
								Usage:        `Search for a defined job result`,
								Description:  help.Text(`job-result-mgmt::search`),
								Action:       runtime(jobResultMgmtSearch),
								BashComplete: cmpl.DirectIdName,
							},
						},
					},
					{
						Name:        `status-mgmt`,
						Usage:       `SUBCOMMANDS for job status management`,
						Description: help.Text(`job-status-mgmt::`),
						Subcommands: []cli.Command{
							{
								Name:        `add`,
								Usage:       `Add a job status`,
								Description: help.Text(`job-status-mgmt::add`),
								Action:      runtime(jobStatusMgmtAdd),
							},
							{
								Name:        `remove`,
								Usage:       `Remove an unused job status`,
								Description: help.Text(`job-status-mgmt::remove`),
								Action:      runtime(jobStatusMgmtRemove),
							},
							{
								Name:        `show`,
								Usage:       `Show details about a job status`,
								Description: help.Text(`job-status-mgmt::show`),
								Action:      runtime(jobStatusMgmtShow),
							},
							{
								Name:        `list`,
								Usage:       `List all defined job status`,
								Description: help.Text(`job-status-mgmt::list`),
								Action:      runtime(jobStatusMgmtList),
							},
							{
								Name:         `search`,
								Usage:        `Search for a defined job status`,
								Description:  help.Text(`job-status-mgmt::search`),
								Action:       runtime(jobStatusMgmtSearch),
								BashComplete: cmpl.DirectIdName,
							},
						},
					},
				},
			},
		}...,
	)
	return &app
}

func jobList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/job/`, `list`, nil, c)
}

func jobShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if !adm.IsUUID(c.Args().First()) {
		return fmt.Errorf("Argument is not a UUID: %s",
			c.Args().First())
	}

	path := fmt.Sprintf("/job/%s", c.Args().First())
	return adm.Perform(`get`, path, `show`, nil, c)
}

func jobWait(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if !adm.IsUUID(c.Args().First()) {
		return fmt.Errorf("Argument is not a UUID: %s",
			c.Args().First())
	}

	path := fmt.Sprintf("/job/byID/%s/_processed", c.Args().First())
	return adm.Perform(`get`, path, `wait`, nil, c)
}

func jobSearch(c *cli.Context) error {
	return fmt.Errorf(`Not implemented`)
}

func clientlocalJobListOutstanding(c *cli.Context) error {
	jobs, err := store.ActiveJobs()
	if err != nil && err != bolt.ErrBucketNotFound {
		return err
	}

	pj := []proto.Job{}
	for _, iArray := range jobs {
		pj = append(pj, proto.Job{
			ID:       iArray[1],
			TsQueued: iArray[2],
			Type:     iArray[3],
		})
	}

	enc, err := json.Marshal(&pj)
	if err != nil {
		return err
	}

	fmt.Println(string(enc))
	// XXX adm.FormatOut support missing
	return nil
}

func clientlocalJobUpdate(c *cli.Context) error {
	jobs, err := store.ActiveJobs()
	if err != nil && err != bolt.ErrBucketNotFound {
		return err
	} else if err == bolt.ErrBucketNotFound {
		// nothing found
		return nil
	}

	req := proto.NewJobFilter()
	req.Flags.Detailed = c.Bool(`verbose`)
	jobMap := map[string]string{}
	for _, v := range jobs {
		// jobID -> storeID
		jobMap[v[1]] = v[0]
		req.Filter.Job.IDList = append(req.Filter.Job.IDList, v[1])
	}
	resp, err := adm.PostReqBody(req, `/search/job/`)
	if err != nil {
		return fmt.Errorf("Job update request error: %s", err)
	}
	var res *proto.Result
	if err = adm.DecodedResponse(resp, res); err != nil {
		return err
	}
	if res.Jobs == nil {
		return fmt.Errorf("Result contained no jobs array")
	}
	for _, j := range *res.Jobs {
		if j.Status != `processed` {
			// only finish Jobs in DB that actually finished
			continue
		}
		strID := jobMap[j.ID]
		var storeID uint64
		if err := adm.ValidateLBoundUint64(strID, &storeID,
			0); err != nil {
			return fmt.Errorf("somaadm: Job update cache error: %s",
				err.Error())
		}
		if err := store.FinishJob(storeID, &j); err != nil {
			return fmt.Errorf("somaadm: Job update cache error: %s",
				err.Error())
		}
	}
	return adm.FormatOut(c, resp.Body(), `list`)
}

func clientlocalJobListLocal(c *cli.Context) error {
	active, err := store.ActiveJobs()
	if err != nil && err != bolt.ErrBucketNotFound {
		return err
	}

	jobs := []proto.Job{}
	for _, iArray := range active {
		jobs = append(jobs, proto.Job{
			ID:       iArray[1],
			TsQueued: iArray[2],
			Type:     iArray[3],
		})
	}

	finished, err := store.FinishedJobs()
	if err != nil && err != bolt.ErrBucketNotFound {
		return err
	}

	jobs = append(jobs, finished...)
	enc, err := json.Marshal(&jobs)
	if err != nil {
		return err
	}
	fmt.Println(string(enc))
	// XXX adm.FormatOut support missing
	return nil
}

func clientlocalJobPruneDB(c *cli.Context) error {
	return store.PruneFinishedJobs()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
