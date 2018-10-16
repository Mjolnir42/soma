/*-
 * Copyright (c) 2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main // import "github.com/mjolnir42/soma/cmd/soma"

// The error stack is used to defer returning errors from functions
// where the function signature is cleaner to use without having an
// error in the return

var errors []error

func init() {
	errors = make([]error, 0)
}

func pushError(err error) {
	errors = append(errors, err)
}

func popError() error {
	var ret error
	switch len(errors) {
	case 0:
		ret = nil
		errors = make([]error, 0)
	default:
		ret = errors[len(errors)-1]
		errors = errors[:len(errors)-1]
	}
	return ret
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
