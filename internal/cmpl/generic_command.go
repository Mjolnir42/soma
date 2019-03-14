package cmpl

import "github.com/codegangsta/cli"

func Datacenter(c *cli.Context) {
	Generic(c, []string{`datacenter`})
}

func In(c *cli.Context) {
	Generic(c, []string{`in`})
}

func DirectIn(c *cli.Context) {
	GenericDirect(c, []string{`in`})
}

func DirectInOf(c *cli.Context) {
	GenericDirect(c, []string{`in`, `of`})
}

func DirectInFrom(c *cli.Context) {
	GenericDirect(c, []string{`from`, `in`})
}

func InTo(c *cli.Context) {
	Generic(c, []string{`in`, `to`})
}

func InFrom(c *cli.Context) {
	Generic(c, []string{`in`, `from`})
}

func InFromView(c *cli.Context) {
	Generic(c, []string{`in`, `from`, `view`})
}

func From(c *cli.Context) {
	Generic(c, []string{`from`})
}

func FromTo(c *cli.Context) {
	Generic(c, []string{`from`, `to`})
}

func FromView(c *cli.Context) {
	Generic(c, []string{`from`, `view`})
}

func Name(c *cli.Context) {
	Generic(c, []string{`name`})
}

func To(c *cli.Context) {
	Generic(c, []string{`to`})
}

func Augmented(c *cli.Context, dispatcher string, data []string) {
	switch dispatcher {
	case `to`:
		GenericDataFirst(c, data, []string{`to`})
	case `from`:
		GenericDataFirst(c, data, []string{`from`})
	case `in`:
		GenericDataFirst(c, data, []string{`in`})
	}
}

func TripleToOn(c *cli.Context) {
	GenericTriple(c, []string{`to`, `on`})
}

func TripleFromOn(c *cli.Context) {
	GenericTriple(c, []string{`from`, `on`})
}

func ValidityAdd(c *cli.Context) {
	Generic(c, []string{`on`, `direct`, `inherited`})
}

func WorkflowSet(c *cli.Context) {
	Generic(c, []string{`status`, `next`})
}

func DirectIDName(c *cli.Context) {
	GenericDirect(c, []string{`id`, `name`})
}

func RepositoryConfigSearch(c *cli.Context) {
	GenericDirect(c, []string{`id`, `name`, `team`, `deleted`, `active`})
}

func CheckConfigList(c *cli.Context) {
	GenericDirectTriple(c, []string{`in`})
}

func CheckConfigDestroy(c *cli.Context) {
	GenericTriple(c, []string{`in`})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
