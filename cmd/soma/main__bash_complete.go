/*-
 * Copyright (c) 2015-2018, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main // import "github.com/mjolnir42/soma/cmd/soma"

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
)

// in main and not the cmpl lib because the full runtime is required
// to provide the completion options. This means we need access to
// globals that do not fit the function signature
func bashCompServiceAdd(c *cli.Context) {
	multipleAllowed := []string{}
	uniqueOptions := []string{}

	// sort attributes based on their cardinality so we can use them
	// for command line parsing
	for _, attr := range attributeFetch() {
		switch attr.Cardinality {
		case `once`:
			uniqueOptions = append(uniqueOptions, attr.Name)
		case `multi`:
			multipleAllowed = append(multipleAllowed, attr.Name)
		default:
			adm.Abort(fmt.Sprintf("Unknown attribute cardinality: %s",
				attr.Cardinality))
		}
	}
	cmpl.GenericMulti(c, uniqueOptions, multipleAllowed)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
