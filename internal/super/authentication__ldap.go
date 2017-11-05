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
	"gopkg.in/ldap.v2"
)

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
