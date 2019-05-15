/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

func receiveRequestCheck(r ReceiveRequest, b Builder) bool {
	if r.ParentType == b.GetType() && (r.ParentID == b.GetID() || r.ParentName == b.GetName()) {
		return true
	}
	return false
}

func unlinkRequestCheck(u UnlinkRequest, b Builder) bool {
	if u.ParentType == b.GetType() && (u.ParentID == b.GetID() || u.ParentName == b.GetName()) {
		return true
	}
	return false
}

func findRequestCheck(f FindRequest, b Builder) bool {
	if f.ElementID == b.GetID() || (f.ElementType == b.GetType() && f.ElementName == b.GetName()) {
		return true
	}
	return false
}

func countAttributeConstraints(attributeC map[string][]string) int {
	var count int
	for key := range attributeC {
		count = count + len(attributeC[key])
	}
	return count
}

func removeFromArray(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
