package somaproto

type ProtoRequestProvider struct {
	Provider *ProtoProvider       `json:"provider,omitempty"`
	Filter   *ProtoProviderFilter `json:"filter,omitempty"`
}

type ProtoResultProvider struct {
	Code      uint16          `json:"code,omitempty"`
	Status    string          `json:"status,omitempty"`
	Text      []string        `json:"text,omitempty"`
	Providers []ProtoProvider `json:"providers,omitempty"`
	JobId     string          `json:"jobid,omitempty"`
}

type ProtoProvider struct {
	Provider string                `json:"provider,omitempty"`
	Details  *ProtoProviderDetails `json:"details,omitempty"`
}

type ProtoProviderFilter struct {
	Provider string `json:"provider,omitempty"`
}

type ProtoProviderDetails struct {
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

//
func (p *ProtoResultProvider) ErrorMark(err error, imp bool, found bool,
	length int) bool {
	if p.markError(err) {
		return true
	}
	if p.markImplemented(imp) {
		return true
	}
	if p.markFound(found, length) {
		return true
	}
	return false
}

func (p *ProtoResultProvider) markError(err error) bool {
	if err != nil {
		p.Code = 500
		p.Status = "ERROR"
		p.Text = []string{err.Error()}
		return true
	}
	return false
}

func (p *ProtoResultProvider) markImplemented(f bool) bool {
	if f {
		p.Code = 501
		p.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (p *ProtoResultProvider) markFound(f bool, i int) bool {
	if f || i == 0 {
		p.Code = 404
		p.Status = "NOT FOUND"
		return true
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
