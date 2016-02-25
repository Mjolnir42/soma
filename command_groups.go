package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func cmdGroupCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys, // as uniqKeys
		multKeys, // as reqKeys
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])

	var req somaproto.ProtoRequestGroup
	req.Group = &somaproto.ProtoGroup{}
	req.Group.Name = c.Args().First()
	req.Group.BucketId = bucketId

	resp := utl.PostRequestWithBody(req, "/groups/")
	fmt.Println(resp)
}

func cmdGroupDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys, // as uniqKeys
		multKeys, // as reqKeys
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	groupId := utl.TryGetGroupByUUIDOrName(
		c.Args().First(),
		bucketId)
	path := fmt.Sprintf("/groups/%s", groupId)

	_ = utl.DeleteRequest(path)
}

func cmdGroupRename(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"to", "bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys, // as uniqKeys
		multKeys, // as reqKeys
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	groupId := utl.TryGetGroupByUUIDOrName(
		c.Args().First(),
		bucketId)
	path := fmt.Sprintf("/groups/%s", groupId)

	var req somaproto.ProtoRequestGroup
	req.Group = &somaproto.ProtoGroup{}
	req.Group.Name = opts["to"][0]

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdGroupList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 0)
	/*
		multiple := []string{}
		unique := []string{"bucket"}
		required := []string{"bucket"}

		opts := utl.ParseVariadicArguments(
			multiple,
			unique,
			required,
			c.Args())

		req := somaproto.ProtoRequestGroup{}
		req.Group = &somaproto.ProtoGroup{}
		req.Filter = &somaproto.ProtoGroupFilter{}
		req.Filter.BucketId = utl.BucketByUUIDOrName(opts["bucket"][0])
		resp := utl.GetRequestWithBody(req, "/groups/")
	*/
	resp := utl.GetRequest("/groups/")
	fmt.Println(resp)
}

func cmdGroupShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	groupId := utl.TryGetGroupByUUIDOrName(
		c.Args().First(),
		bucketId)
	path := fmt.Sprintf("/groups/%s", groupId)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

func cmdGroupMemberAddGroup(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"to", "bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	mGroupId := utl.TryGetGroupByUUIDOrName(
		c.Args().First(),
		bucketId)
	groupId := utl.TryGetGroupByUUIDOrName(
		opts["to"][0],
		bucketId)

	var req somaproto.ProtoRequestGroup
	var group somaproto.ProtoGroup
	group.Id = mGroupId
	req.Group = &somaproto.ProtoGroup{}
	req.Group.Id = groupId
	req.Group.BucketId = bucketId
	req.Group.MemberGroups = append(req.Group.MemberGroups, group)

	path := fmt.Sprintf("/groups/%s/members/", groupId)

	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdGroupMemberAddCluster(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"to", "bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	mClusterId := utl.TryGetClusterByUUIDOrName(
		c.Args().First(),
		bucketId)
	groupId := utl.TryGetGroupByUUIDOrName(
		opts["to"][0],
		bucketId)

	var req somaproto.ProtoRequestGroup
	var cluster somaproto.ProtoCluster
	cluster.Id = mClusterId
	req.Group = &somaproto.ProtoGroup{}
	req.Group.Id = groupId
	req.Group.BucketId = bucketId
	req.Group.MemberClusters = append(req.Group.MemberClusters, cluster)

	path := fmt.Sprintf("/groups/%s/members/", groupId)

	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdGroupMemberAddNode(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"to", "bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	mNodeId := utl.TryGetNodeByUUIDOrName(c.Args().First())
	groupId := utl.TryGetGroupByUUIDOrName(
		opts["to"][0],
		bucketId)

	var req somaproto.ProtoRequestGroup
	var node somaproto.ProtoNode
	node.Id = mNodeId
	req.Group = &somaproto.ProtoGroup{}
	req.Group.Id = groupId
	req.Group.BucketId = bucketId
	req.Group.MemberNodes = append(req.Group.MemberNodes, node)

	path := fmt.Sprintf("/groups/%s/members/", groupId)

	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdGroupMemberDeleteGroup(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"from", "bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	mGroupId := utl.TryGetGroupByUUIDOrName(
		c.Args().First(),
		bucketId)
	groupId := utl.TryGetGroupByUUIDOrName(
		opts["from"][0],
		bucketId)

	path := fmt.Sprintf("/groups/%s/members/%s", groupId,
		mGroupId)

	_ = utl.DeleteRequest(path)
}

func cmdGroupMemberDeleteCluster(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"from", "bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	mClusterId := utl.TryGetClusterByUUIDOrName(
		c.Args().First(),
		bucketId)
	groupId := utl.TryGetGroupByUUIDOrName(
		opts["from"][0],
		bucketId)

	path := fmt.Sprintf("/groups/%s/members/%s", groupId,
		mClusterId)

	_ = utl.DeleteRequest(path)
}

func cmdGroupMemberDeleteNode(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"from", "bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	mNodeId := utl.TryGetNodeByUUIDOrName(c.Args().First())
	groupId := utl.TryGetGroupByUUIDOrName(
		opts["from"][0],
		bucketId)

	path := fmt.Sprintf("/groups/%s/members/%s", groupId,
		mNodeId)

	_ = utl.DeleteRequest(path)
}

func cmdGroupMemberList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	groupId := utl.TryGetGroupByUUIDOrName(
		c.Args().First(),
		bucketId)

	path := fmt.Sprintf("/groups/%s/members/", groupId)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

func cmdGroupSystemPropertyAdd(c *cli.Context) {
	utl.ValidateCliMinArgumentCount(c, 9)
	multiple := []string{}
	required := []string{"to", "in", "value", "view"}
	unique := []string{"to", "in", "value", "view", "inheritance", "childrenonly"}

	opts := utl.ParseVariadicArguments(multiple, required, unique, c.Args().Tail())
	bucketId := utl.BucketByUUIDOrName(opts["in"][0])
	groupId := utl.TryGetGroupByUUIDOrName(opts["to"][0], bucketId)
	utl.CheckStringIsSystemProperty(c.Args().First())

	sprop := somaproto.TreePropertySystem{
		Name:  c.Args().First(),
		Value: opts["value"][0],
	}

	tprop := somaproto.TreeProperty{
		PropertyType: "system",
		View:         opts["view"][0],
		System:       &sprop,
	}
	if _, ok := opts["inheritance"]; ok {
		tprop.Inheritance = utl.GetValidatedBool(opts["inheritance"][0])
	} else {
		tprop.Inheritance = true
	}
	if _, ok := opts["childrenonly"]; ok {
		tprop.ChildrenOnly = utl.GetValidatedBool(opts["childrenonly"][0])
	} else {
		tprop.ChildrenOnly = false
	}

	propList := []somaproto.TreeProperty{tprop}

	group := somaproto.ProtoGroup{
		Id:         groupId,
		BucketId:   bucketId,
		Properties: &propList,
	}

	req := somaproto.ProtoRequestGroup{
		Group: &group,
	}

	path := fmt.Sprintf("/groups/%s/property/system/", groupId)
	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
