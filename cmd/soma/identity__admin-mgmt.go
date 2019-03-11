/*-
 * Copyright (c) 2019, Jörg Pernfuß
 * Copyright (c) 2019, 1&1 IONOS SE
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
	"github.com/mjolnir42/soma/lib/proto"
)

// adminMgmtAdd function
// soma user-mgmt admin grant ${username}
func adminMgmtAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	var err error
	var userID string
	if userID, err = adm.LookupUserID(c.Args().First()); err != nil {
		return err
	}
	if err := adm.ValidateNotUUID(c.Args().First()); err != nil {
		return err
	}

	req := proto.NewAdminRequest()
	req.Admin.UserName = c.Args().First()

	path := fmt.Sprintf("/user/%s/admin", url.QueryEscape(userID))
	return adm.Perform(`putbody`, path, `admin-mgmt::add`, req, c)
}

// adminMgmtRemove function
// soma user-mgmt admin revoke ${username}
func adminMgmtRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	var err error
	var userID, adminID string
	if userID, err = adm.LookupUserID(c.Args().First()); err != nil {
		return err
	}
	if adminID, err = adm.LookupAdminID(userID); err != nil {
		return err
	}

	path := fmt.Sprintf("/user/%s/admin/%s",
		url.QueryEscape(userID),
		url.QueryEscape(adminID),
	)
	return adm.Perform(`delete`, path, `admin-mgmt::remove`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
