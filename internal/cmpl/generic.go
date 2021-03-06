package cmpl

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func Generic(c *cli.Context, keywords []string) {
	switch {
	case c.NArg() == 0:
		return
	case c.NArg() == 1:
		for _, t := range keywords {
			fmt.Println(t)
		}
		return
	}

	skip := 0
	match := make(map[string]bool)

	for _, t := range c.Args().Tail() {
		if skip > 0 {
			skip--
			continue
		}
		skip = 1
		match[t] = true
		continue
	}
	// do not complete in positions where arguments are expected
	if skip > 0 {
		return
	}
	for _, t := range keywords {
		if !match[t] {
			fmt.Println(t)
		}
	}
}

func GenericMulti(c *cli.Context, singlewords, multiwords []string) {
	keywords := append(singlewords, multiwords...)

	switch {
	case c.NArg() == 0:
		return
	case c.NArg() == 1:
		for _, t := range keywords {
			fmt.Println(t)
		}
		return
	}

	skip := 0
	match := make(map[string]bool)

	for _, t := range c.Args().Tail() {
		if skip > 0 {
			skip--
			continue
		}
		skip = 1
		match[t] = true
		continue
	}
	// do not complete in positions where arguments are expected
	if skip > 0 {
		return
	}
	for _, t := range singlewords {
		if !match[t] {
			fmt.Println(t)
		}
	}
	for _, t := range multiwords {
		fmt.Println(t)
	}
}

func GenericDirect(c *cli.Context, keywords []string) {
	switch {
	case c.NArg() == 0:
		for _, t := range keywords {
			fmt.Println(t)
		}
		return
	case c.NArg() == 1:
		return
	}

	skip := 0
	match := make(map[string]bool)

	fullArgs := append([]string{c.Args().First()}, c.Args().Tail()...)

	for _, t := range fullArgs {
		if skip > 0 {
			skip--
			continue
		}
		skip = 1
		match[t] = true
		continue
	}
	// do not complete in positions where arguments are expected
	if skip > 0 {
		return
	}
	for _, t := range keywords {
		if !match[t] {
			fmt.Println(t)
		}
	}
}

func GenericTriple(c *cli.Context, keywords []string) {
	switch {
	case c.NArg() == 0:
		return
	case c.NArg() == 1:
		for _, t := range keywords {
			fmt.Println(t)
		}
		return
	}

	skip := 0
	match := make(map[string]bool)

	for _, t := range c.Args().Tail() {
		if skip > 0 {
			skip--
			continue
		}
		skip = 2
		match[t] = true
		continue
	}
	// do not complete in positions where arguments are expected
	if skip > 0 {
		return
	}
	for _, t := range keywords {
		if !match[t] {
			fmt.Println(t)
		}
	}
}

func None(c *cli.Context) {
}

func GenericDataOnly(c *cli.Context, data []string) {
	GenericDataFirst(c, data, []string{})
}

func GenericDataFirst(c *cli.Context, data, keywords []string) {
	switch {
	case c.NArg() == 0:
		for _, entry := range data {
			fmt.Println(entry)
		}
		return
	case c.NArg() == 1:
		for _, t := range keywords {
			fmt.Println(t)
		}
		return
	}

	skip := 0
	match := make(map[string]bool)

	for _, t := range c.Args().Tail() {
		if skip > 0 {
			skip--
			continue
		}
		skip = 1
		match[t] = true
		continue
	}
	// do not complete in positions where arguments are expected
	if skip > 0 {
		return
	}
	for _, t := range keywords {
		if !match[t] {
			fmt.Println(t)
		}
	}
}

func GenericDirectTriple(c *cli.Context, keywords []string) {
	switch {
	case c.NArg() == 0:
		for _, t := range keywords {
			fmt.Println(t)
		}
		return
	case c.NArg() == 1:
		return
	case c.NArg() == 2:
		return
	}

	skip := 0
	match := make(map[string]bool)

	fullArgs := append([]string{c.Args().First()}, c.Args().Tail()...)

	for _, t := range fullArgs {
		if skip > 0 {
			skip--
			continue
		}
		skip = 2
		match[t] = true
		continue
	}
	// do not complete in positions where arguments are expected
	if skip > 0 {
		return
	}
	for _, t := range keywords {
		if !match[t] {
			fmt.Println(t)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
