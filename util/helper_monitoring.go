package util

import (
	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u *SomaUtil) TryGetMonitoringByUUIDOrName(s string) string {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		return u.GetMonitoringIdByName(s)
	}
	return id.String()
}

func (u *SomaUtil) GetMonitoringIdByName(monitoring string) string {
	req := somaproto.ProtoRequestMonitoring{}
	req.Filter = &somaproto.ProtoMonitoringFilter{}
	req.Filter.Name = monitoring

	resp := u.PostRequestWithBody(req, "/filter/monitoring/")
	monitoringResult := u.DecodeProtoResultMonitoringFromResponse(resp)

	if monitoring != monitoringResult.Systems[0].Name {
		u.Abort("Received result set for incorrect monitoring system")
	}
	return monitoringResult.Systems[0].Id
}

func (u *SomaUtil) DecodeProtoResultMonitoringFromResponse(resp *resty.Response) *somaproto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
