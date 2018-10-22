package adm

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func FormatOut(c *cli.Context, data []byte, cmd string) error {
	if string(data) == `` {
		return nil
	}

	if c.GlobalBool(`json`) {
		fmt.Println(string(data))
		return nil
	}

	// hardwire JSON output for now
	fmt.Println(string(data))

	/* TODO
	switch cmd {
	case `list`:
	case `show`:
	case `tree`:
	default:
	}
	*/

	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
