/*-
Copyright (c) 2013 Julien Schmidt. All rights reserved.
Copyright (c) 2016-2017 Jörg Pernfuß <joerg.pernfuss@1und1.de>


Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
    * The names of the contributors may not be used to endorse or promote
      products derived from this software without specific prior written
      permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL JULIEN SCHMIDT BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.


The following code is nearly verbatim the example code from the httprouter
distribution. Therefor copyright is set to the license text of that distribution.
*/

package rest // import "github.com/mjolnir42/soma/internal/rest"

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/mjolnir42/soma/internal/handler"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/satori/go.uuid"

	"github.com/julienschmidt/httprouter"
)

// BasicAuth handles HTTP BasicAuth on requests
func (x *Rest) BasicAuth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {
		const basicAuthPrefix string = "Basic "
		var supervisor handler.Handler

		// generate and record the requestID
		requestID := uuid.Must(uuid.NewV4())
		ps = append(ps, httprouter.Param{
			Key:   `RequestID`,
			Value: requestID.String(),
		})

		// record the request URI
		ps = append(ps, httprouter.Param{
			Key:   `RequestURI`,
			Value: r.RequestURI,
		})

		logEntry := x.reqLog.WithField(`RequestID`, requestID.String()).
			WithField(`RequestURI`, r.RequestURI).
			WithField(`Phase`, `basic-auth`)

		// if the supervisor is not available, no requests are accepted
		if supervisor = x.handlerMap.Get(`supervisor`); supervisor == nil {
			http.Error(w, `Authentication supervisor not available`,
				http.StatusServiceUnavailable)
			logEntry.WithField(`Code`, http.StatusServiceUnavailable).
				Warn(`Authentication supervisor not available`)
			return
		}

		// disable authentication much?
		if x.conf.OpenInstance {
			ps = append(ps, httprouter.Param{
				Key:   `AuthenticatedUser`,
				Value: `AnonymousCoward`,
			})
			h(w, r, ps)
			return
		}

		// Get credentials
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, basicAuthPrefix) {
			// Check credentials
			payload, err := base64.StdEncoding.DecodeString(
				auth[len(basicAuthPrefix):])
			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				if len(pair) == 2 {
					request := newRequest(r, ps)
					request.Section = msg.SectionSupervisor
					request.Action = msg.ActionAuthenticate
					request.Super = &msg.Supervisor{
						RestrictedEndpoint: false,
						Task:               msg.TaskBasicAuth,
						BasicAuth: struct {
							User  string
							Token string
						}{
							User:  string(pair[0]),
							Token: string(pair[1]),
						},
					}
					supervisor.Intake() <- request
					result := <-request.Reply
					if result.Error != nil {
						// log authentication errors
						logEntry = logEntry.WithField(`AuthenticationUser`, string(pair[0])).
							WithField(`AuthenticationError`, result.Error.Error())
						goto unauthorized
					}
					if result.Super.Verdict == 200 {
						// record the authenticated user
						ps = append(ps, httprouter.Param{
							Key:   `AuthenticatedUser`,
							Value: string(pair[0]),
						})
						// record the used token for supervisor:token/invalidate
						ps = append(ps, httprouter.Param{
							Key:   `AuthenticatedToken`,
							Value: string(pair[1]),
						})

						// log successful basic auth requests only at debug level
						// since they will also be logged by rest.send()
						logEntry.WithField(`AuthenticationUser`, string(pair[0])).
							WithField(`Code`, int(result.Super.Verdict)).
							Debug(http.StatusText(int(result.Super.Verdict)))

						// Delegate request to given handle
						h(w, r, ps)
						return

					}
				}
			}
		} else {
			logEntry = logEntry.WithField(
				`AuthenticationError`,
				`Invalid BasicAuth authorization header`,
			)
		}

	unauthorized:
		w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		http.Error(w, http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized)
		logEntry.WithField(`Code`, http.StatusUnauthorized).
			Info(http.StatusText(http.StatusUnauthorized))
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
