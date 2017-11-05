/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/auth"
	"gopkg.in/ldap.v2"
)

// authenticateLdap verifies credentials provided in token
// against LDAP
func (s *Supervisor) authenticateLdap(token *auth.Token, mr *msg.Result) bool {
	// check provided password via simple bind
	if ok, err := validateLdapCredentials(
		token.UserName,
		token.Token,
	); err != nil {
		mr.ServerError(err, mr.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return false
	} else if !ok {
		mr.Unauthorized(fmt.Errorf(`Invalid LDAP credentials`),
			mr.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return false
	}

	// fail activation if local password is the same as the
	// upstream password. This error _IS_ sent to the user!
	if token.Token == token.Password {
		mr.Conflict(fmt.Errorf(
			"User %s denied: matching local/upstream passwords",
			token.UserName), mr.Section)
		mr.Super.Audit.
			WithField(`Code`, mr.Code).
			Warningln(mr.Error)
		return false
	}
	return true
}

func validateLdapCredentials(user, password string) (bool, error) {
	var (
		conn *ldap.Conn
		err  error
		pem  []byte
	)

	addr := fmt.Sprintf("%s:%d", cfg.Ldap.Address, cfg.Ldap.Port)
	bindDN := strings.Join(
		[]string{
			strings.Join(
				[]string{
					cfg.Ldap.Attribute,
					user,
				},
				`=`,
			),
			cfg.Ldap.UserDN,
			cfg.Ldap.BaseDN,
		},
		`,`,
	)

	if cfg.Ldap.TLS {
		conf := &tls.Config{
			InsecureSkipVerify: cfg.Ldap.SkipVerify,
			ServerName:         cfg.Ldap.Address,
			MinVersion:         tls.VersionTLS12,
			MaxVersion:         tls.VersionTLS12,
			CipherSuites: []uint16{
				// TODO this should probably be configurable
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			},
		}
		if cfg.Ldap.Cert != "" {
			if pem, err = ioutil.ReadFile(cfg.Ldap.Cert); err != nil {
				return false, err
			}
			conf.RootCAs = x509.NewCertPool()
			conf.RootCAs.AppendCertsFromPEM(pem)
		}
		conn, err = ldap.DialTLS(`tcp`, addr, conf)
	} else {
		log.Println(`REALLY?!! Using unencrypted LDAP connection. Grudgingly.`)
		conn, err = ldap.Dial(`tcp`, addr)
	}
	if err != nil {
		return false, err
	}
	defer conn.Close()

	// attempt bind
	err = conn.Bind(bindDN, password)
	if err != nil && ldap.IsErrorWithCode(err,
		ldap.LDAPResultInvalidCredentials) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
