package adm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/lib/proto"
	"gopkg.in/resty.v0"
)

// Exported functions

// WRAPPER
func Perform(rqType, path, tmpl string, body interface{}, c *cli.Context) error {
	var (
		err  error
		resp *resty.Response
	)

	if strings.HasSuffix(rqType, `body`) && body == nil {
		goto noattachment
	}

	switch rqType {
	case `get`:
		resp, err = GetReq(path)
	case `head`:
		resp, err = HeadReq(path)
	case `delete`:
		resp, err = DeleteReq(path)
	case `deletebody`:
		resp, err = DeleteReqBody(body, path)
	case `putbody`:
		resp, err = PutReqBody(body, path)
	case `postbody`:
		resp, err = PostReqBody(body, path)
	case `patchbody`:
		resp, err = PatchReqBody(body, path)
	}

	if err != nil {
		return err
	}
	return FormatOut(c, resp, tmpl)

noattachment:
	return fmt.Errorf(`Missing body to client request that requires it.`)
}

func DecodedResponse(resp *resty.Response, res *proto.Result) error {
	if err := decodeResponse(resp, res); err != nil {
		return err
	}
	return checkApplicationError(res)
}

// DELETE
func DeleteReq(p string) (*resty.Response, error) {
	return handleRequestOptions(client.R().Delete(p))
}

func DeleteReqBody(body interface{}, p string) (*resty.Response, error) {
	return handleRequestOptions(
		client.R().SetBody(body).SetContentLength(true).Delete(p))
}

// GET
func GetReq(p string) (*resty.Response, error) {
	return handleRequestOptions(client.R().Get(p))
}

// HEAD
func HeadReq(p string) (*resty.Response, error) {
	return handleRequestOptions(client.R().Head(p))
}

// PATCH
func PatchReqBody(body interface{}, p string) (*resty.Response, error) {
	return handleRequestOptions(
		client.R().SetBody(body).SetContentLength(true).Patch(p))
}

// POST
func PostReqBody(body interface{}, p string) (*resty.Response, error) {
	return handleRequestOptions(
		client.R().SetBody(body).SetContentLength(true).Post(p))
}

// PUT
func PutReq(p string) (*resty.Response, error) {
	return handleRequestOptions(client.R().Put(p))
}

func PutReqBody(body interface{}, p string) (*resty.Response, error) {
	return handleRequestOptions(
		client.R().SetBody(body).SetContentLength(true).Put(p))
}

// Private functions

func handleRequestOptions(resp *resty.Response, err error) (*resty.Response, error) {
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() >= 300 {
		return resp, fmt.Errorf("Request error: %s, %s", resp.Status(), resp.String())
	}

	if !(async || jobSave) {
		return resp, nil
	}

	var result *proto.Result
	if err = decodeResponse(resp, result); err != nil {
		return nil, err
	}

	if jobSave {
		if result.StatusCode == 202 && result.JobID != "" {
			cache.SaveJob(result.JobID, result.JobType)
		}
	}

	if async {
		asyncWait(result)
	}
	return resp, nil
}

func asyncWait(result *proto.Result) {
	if !async {
		return
	}

	if result.StatusCode == 202 && result.JobID != "" {
		fmt.Fprintf(os.Stderr, "Waiting for job: %s\n", result.JobID)
		_, err := PutReq(fmt.Sprintf("/job/%s", result.JobID))
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "Wait error: %s\n", err.Error())
		}
	}
}

func decodeResponse(resp *resty.Response, res *proto.Result) error {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body()))
	return decoder.Decode(res)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
