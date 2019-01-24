package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/lib/proto"
)

func cmdNodeAdd(c *cli.Context) error {
	opts := map[string][]string{}
	uniqKeys := []string{`assetid`, `name`, `team`, `server`, `online`}
	reqKeys := []string{`assetid`, `name`, `team`}

	var err error
	if err = adm.ParseVariadicArguments(
		opts,
		[]string{},
		uniqKeys,
		reqKeys,
		adm.AllArguments(c),
	); err != nil {
		return err
	}
	req := proto.NewNodeRequest()

	if _, ok := opts[`online`]; ok {
		if err = adm.ValidateBool(opts[`online`][0],
			&req.Node.IsOnline); err != nil {
			return err
		}
	} else {
		req.Node.IsOnline = true
	}
	if _, ok := opts[`server`]; ok {
		if req.Node.ServerID, err = adm.LookupServerID(
			opts[`server`][0]); err != nil {
			return err
		}
	}
	if err = adm.ValidateLBoundUint64(opts[`assetid`][0],
		&req.Node.AssetID, 1); err != nil {
		return err
	}
	req.Node.Name = opts[`name`][0]
	if err = adm.LookupTeamID(
		opts[`team`][0],
		&req.Node.TeamID,
	); err != nil {
		return nil
	}

	return adm.Perform(`postbody`, `/node/`, `command`, req, c)
}

func cmdNodeUpdate(c *cli.Context) error {
	unique := []string{`name`, `assetid`, `server`, `team`,
		`online`, `deleted`}
	required := []string{`name`, `assetid`, `server`, `team`,
		`online`, `deleted`}
	opts := map[string][]string{}

	var err error
	if err = adm.ParseVariadicArguments(
		opts,
		[]string{},
		unique,
		required,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	req := proto.NewNodeRequest()
	if !adm.IsUUID(c.Args().First()) {
		return fmt.Errorf(
			`Node/update command requires UUID as first argument`)
	}
	req.Node.ID = c.Args().First()
	req.Node.Name = opts[`name`][0]
	if err = adm.ValidateBool(opts[`online`][0],
		&req.Node.IsOnline); err != nil {
		return err
	}
	if err = adm.ValidateBool(opts[`deleted`][0],
		&req.Node.IsDeleted); err != nil {
		return err
	}
	if req.Node.ServerID, err = adm.LookupServerID(
		opts[`server`][0]); err != nil {
		return err
	}
	if err = adm.ValidateLBoundUint64(opts[`assetid`][0],
		&req.Node.AssetID, 1); err != nil {
		return err
	}
	if err = adm.LookupTeamID(
		opts[`team`][0],
		&req.Node.TeamID,
	); err != nil {
		return err
	}
	path := fmt.Sprintf("/node/%s", req.Node.ID)
	return adm.Perform(`putbody`, path, `command`, req, c)
}

func cmdNodeDel(c *cli.Context) (err error) {
	if err = adm.VerifySingleArgument(c); err != nil {
		return err
	}
	var id, path string
	if id, err = adm.LookupNodeID(c.Args().First()); err != nil {
		return err
	}
	path = fmt.Sprintf("/node/%s", id)

	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdNodePurge(c *cli.Context) (err error) {
	var (
		id, path string
		req      proto.Request
	)
	if c.Bool(`all`) {
		if err = adm.VerifyNoArgument(c); err != nil {
			return err
		}
		path = "/node/"
	} else {
		if err = adm.VerifySingleArgument(c); err != nil {
			return err
		}
		if id, err = adm.LookupNodeID(c.Args().First()); err != nil {
			return err
		}
		path = fmt.Sprintf("/node/%s", id)
	}

	req = proto.Request{
		Flags: &proto.Flags{
			Purge: true,
		},
	}

	return adm.Perform(`deletebody`, path, `command`, req, c)
}

func cmdNodeRestore(c *cli.Context) (err error) {
	var (
		id, path string
		req      proto.Request
	)
	if c.Bool(`all`) {
		if err = adm.VerifyNoArgument(c); err != nil {
			return err
		}
		path = "/node/"
	} else {
		if err = adm.VerifySingleArgument(c); err != nil {
			return err
		}
		if id, err = adm.LookupNodeID(c.Args().First()); err != nil {
			return err
		}
		path = fmt.Sprintf("/node/%s", id)
	}

	req = proto.Request{
		Flags: &proto.Flags{
			Restore: true,
		},
	}

	return adm.Perform(`deletebody`, path, `command`, req, c)
}

func cmdNodeRename(c *cli.Context) (err error) {
	opts := map[string][]string{}
	if err = adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`to`},
		[]string{`to`},
		c.Args().Tail()); err != nil {
		return err
	}
	var id, path string
	if id, err = adm.LookupNodeID(c.Args().First()); err != nil {
		return err
	}
	path = fmt.Sprintf("/node/%s", id)

	req := proto.NewNodeRequest()
	req.Node.Name = opts[`to`][0]

	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdNodeRepo(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`to`},
		[]string{`to`},
		c.Args().Tail()); err != nil {
		return err
	}
	var id, teamID string
	{
		var err error
		if id, err = adm.LookupNodeID(c.Args().First()); err != nil {
			return err
		}
		if err = adm.LookupTeamID(opts[`to`][0], &teamID); err != nil {
			return err
		}
	}
	path := fmt.Sprintf("/node/%s", id)

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.TeamID = teamID

	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdNodeMove(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`to`},
		[]string{`to`},
		c.Args().Tail()); err != nil {
		return err
	}
	var id string
	{
		var err error
		if id, err = adm.LookupNodeID(c.Args().First()); err != nil {
			return err
		}
	}
	server := opts[`to`][0]
	serverID, err := adm.LookupServerID(server)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/node/%s", id)

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.ServerID = serverID

	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdNodeOnline(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupNodeID(c.Args().First())
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/node/%s", id)

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.IsOnline = true

	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdNodeOffline(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupNodeID(c.Args().First())
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/node/%s", id)

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.IsOnline = false

	return adm.Perform(`patchbody`, path, `command`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
