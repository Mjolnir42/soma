package util

import (
	"fmt"

	"github.com/1and1/soma/lib/proto"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) tryGetRepositoryByUUIDOrName(c *resty.Client, s string) string {
	if u.isUUID(s) {
		return s
	}
	return u.getRepositoryIdByName(c, s)
}

func (u SomaUtil) getRepositoryIdByName(c *resty.Client, repo string) string {
	req := proto.Request{
		Filter: &proto.Filter{
			Repository: &proto.RepositoryFilter{
				Name: repo,
			},
		},
	}

	resp := u.PostRequestWithBody(c, req, "/filter/repository/")
	repoResult := u.decodeProtoResultRepositoryFromResponse(resp)

	if repo != (*repoResult.Repositories)[0].Name {
		u.abort("Received result set for incorrect repository")
	}
	return (*repoResult.Repositories)[0].Id
}

func (u SomaUtil) GetTeamIdByRepositoryId(c *resty.Client, repo string) string {
	repoId := u.tryGetRepositoryByUUIDOrName(c, repo)

	resp := u.GetRequest(c, fmt.Sprintf("/repository/%s", repoId))
	repoResult := u.DecodeResultFromResponse(resp)
	return (*repoResult.Repositories)[0].TeamId
}

func (u SomaUtil) getRepositoryDetails(c *resty.Client, repoId string) *proto.Repository {
	resp := u.GetRequest(c, fmt.Sprintf("/repository/%s", repoId))
	res := u.DecodeResultFromResponse(resp)
	return &(*res.Repositories)[0]
}

func (u SomaUtil) findSourceForRepoProperty(c *resty.Client, pTyp, pName, view, repoId string) string {
	repo := u.getRepositoryDetails(c, repoId)
	if repo == nil {
		return ``
	}
	for _, prop := range *repo.Properties {
		// wrong type
		if prop.Type != pTyp {
			continue
		}
		// wrong view
		if prop.View != view {
			continue
		}
		// inherited property
		if prop.InstanceId != prop.SourceInstanceId {
			continue
		}
		switch pTyp {
		case `system`:
			if prop.System.Name == pName {
				return prop.SourceInstanceId
			}
		case `oncall`:
			if prop.Oncall.Name == pName {
				return prop.SourceInstanceId
			}
		case `custom`:
			if prop.Custom.Name == pName {
				return prop.SourceInstanceId
			}
		case `service`:
			if prop.Service.Name == pName {
				return prop.SourceInstanceId
			}
		}
	}
	return ``
}

func (u SomaUtil) decodeProtoResultRepositoryFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
