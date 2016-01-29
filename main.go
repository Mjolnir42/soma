package main

import (
	"database/sql"
	"github.com/codegangsta/cli"
	"os"
)

var Cfg Config
var db *sql.DB

func main() {
	app := cli.NewApp()
	app.Name = "somadbctl"
	app.Usage = "SOMA Database Control Utility"
	app.Version = "0.1.1"

	app = registerCommands(*app)
	app = registerFlags(*app)

	app.Run(os.Args)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
