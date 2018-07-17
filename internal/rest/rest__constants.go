/*-
 * Copyright (c) 2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/soma/internal/rest"

const (
	rtNode                 = `/node/`
	rtNodeID               = `/node/:nodeID`
	rtNodeConfig           = `/node/:nodeID/config`
	rtNodeUnassign         = `/repository/:repositoryID/bucket/:bucketID/node/:nodeID/config`
	rtNodeInstance         = `/repository/:repositoryID/bucket/:bucketID/node/:nodeID/instance/`
	rtNodeInstanceID       = `/repository/:repositoryID/bucket/:bucketID/node/:nodeID/instance/:instanceID`
	rtNodeInstanceVersions = `/repository/:repositoryID/bucket/:bucketID/node/:nodeID/instance/:instanceID/versions`
	rtNodeProperty         = `/repository/:repositoryID/bucket/:bucketID/node/:nodeID/property/:propertyType/`
	rtNodePropertyID       = `/repository/:repositoryID/bucket/:bucketID/node/:nodeID/property/:propertyType/:sourceID`
	rtNodeTree             = `/repository/:repositoryID/bucket/:bucketID/node/:nodeID/tree`
	rtSearchNode           = `/search/node/`
	rtSyncNode             = `/sync/node/`
)

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
