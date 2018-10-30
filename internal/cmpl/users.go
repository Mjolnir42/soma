package cmpl

import "github.com/codegangsta/cli"

func UserMgmtAdd(c *cli.Context) {
	Generic(c, []string{`firstname`, `lastname`, `employeenr`, `mailaddr`, `team`, `system`})
}

func UserMgmtUpdate(c *cli.Context) {
	Generic(c, []string{`username`, `firstname`, `lastname`, `employeenr`, `mailaddr`, `team`, `deleted`})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
