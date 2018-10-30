/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main // import "github.com/mjolnir42/soma/cmd/soma"

import (
	"fmt"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/lib/auth"
)

// supervisorLogin function
// soma login
func supervisorLogin(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	path := `/authenticate/validate`
	return adm.Perform(`head`, path, `command`, nil, c)
}

// supervisorLogout function
// soma logout
func supervisorLogout(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	if c.GlobalBool(`doublelogout`) {
		return adm.MockOK(`command`, c)
	}

	var path string
	switch c.Bool(`all`) {
	case true:
		path = `/tokens/self/all`
	case false:
		path = `/tokens/self/active`
	}
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func supervisorActivate(c *cli.Context) error {
	// administrative use, full runtime is available
	if c.GlobalIsSet(`admin`) {
		if err := adm.VerifySingleArgument(c); err != nil {
			return err
		}
		return runtime(supervisorActivateAdmin)(c)
	}
	// user trying to activate the account for the first
	// time, reduced runtime
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}
	return boottime(supervisorActivateUser)(c)
}

func supervisorActivateUser(c *cli.Context) error {
	var err error
	var password string
	var passKey string
	var happy bool
	var cred *auth.Token

	if Cfg.Auth.User == `` {
		fmt.Println(`Please specify which account to activate.`)
		if Cfg.Auth.User, err = adm.Read(`user`); err != nil {
			return err
		}
	} else {
		fmt.Printf("Starting with activation of account '%s' in"+
			" 2 seconds.\n", Cfg.Auth.User)
		fmt.Printf(`Use --user flag to activate a different account.`)
		time.Sleep(2 * time.Second)
	}
	if strings.Contains(Cfg.Auth.User, `:`) {
		return fmt.Errorf(`Usernames must not contain : character.`)
	}

	fmt.Printf("\nPlease provide the password you want to use.\n")
password_read:
	password = adm.ReadVerified(`password`)

	if happy, err = adm.EvaluatePassword(3, password,
		Cfg.Auth.User, `soma`); err != nil {
		return err
	} else if !happy {
		goto password_read
	}

	fmt.Printf("\nTo confirm that this is your account, an" +
		" additional credential is required" +
		" this once.\n")

	switch Cfg.Activation {
	case `ldap`:
		fmt.Println(`Please provide your LDAP password to"+
		" establish ownership.`)
		passKey = adm.ReadVerified(`password`)
	case `mailtoken`:
		fmt.Println(`Please provide the token you received via email.`)
		passKey = adm.ReadVerified(`token`)
	default:
		return fmt.Errorf(`Unknown activation mode`)
	}

	if cred, err = adm.ActivateAccount(Client, &auth.Token{
		UserName: Cfg.Auth.User,
		Password: password,
		Token:    passKey,
	}); err != nil {
		return err
	}

	// validate received token
	if err = adm.ValidateToken(Client, Cfg.Auth.User,
		cred.Token); err != nil {
		return err
	}
	// save received token
	if err = store.SaveToken(
		Cfg.Auth.User,
		cred.ValidFrom,
		cred.ExpiresAt,
		cred.Token,
	); err != nil {
		return err
	}
	return nil
}

func supervisorActivateAdmin(c *cli.Context) error {
	return nil
}

func supervisorPasswordUpdate(c *cli.Context) error {
	var (
		err               error
		password, passKey string
		happy             bool
		cred              *auth.Token
	)

	if Cfg.Auth.User == `` {
		fmt.Println(`Please specify for which  account the"+
		" password should be changed.`)
		if Cfg.Auth.User, err = adm.Read(`user`); err != nil {
			return err
		}
	} else {
		fmt.Printf("Starting with password update of account '%s'"+
			" in 2 seconds.\n", Cfg.Auth.User)
		fmt.Printf(`Use --user flag to switch account account.`)
		time.Sleep(2 * time.Second)
	}
	if strings.Contains(Cfg.Auth.User, `:`) {
		return fmt.Errorf(`Usernames must not contain : character.`)
	}

	fmt.Printf("\nPlease provide the new password you want to set.\n")
password_read:
	password = adm.ReadVerified(`password`)

	if happy, err = adm.EvaluatePassword(3, password,
		Cfg.Auth.User, `soma`); err != nil {
		return err
	} else if !happy {
		goto password_read
	}

	if c.Bool(`reset`) {
		fmt.Printf("\nTo confirm that you are allowed to reset" +
			" this account, an additional" +
			"credential is required.\n")

		switch Cfg.Activation {
		case `ldap`:
			fmt.Println(`Please provide your LDAP password to` +
				`establish ownership.`)
			passKey = adm.ReadVerified(`password`)
		case `mailtoken`:
			fmt.Println(`Please provide the token you received` +
				`via email.`)
			passKey = adm.ReadVerified(`token`)
		default:
			return fmt.Errorf(`Unknown activation mode`)
		}
	} else {
		fmt.Printf("\nPlease provide your currently active/old" +
			" password.\n")
		passKey = adm.ReadVerified(`password`)
	}

	if cred, err = adm.ChangeAccountPassword(
		Client,
		c.Bool(`reset`),
		&auth.Token{
			UserName: Cfg.Auth.User,
			Password: password,
			Token:    passKey,
		},
	); err != nil {
		return err
	}

	// validate received token
	if err = adm.ValidateToken(
		Client,
		Cfg.Auth.User,
		cred.Token,
	); err != nil {
		return err
	}
	// save received token
	if err = store.SaveToken(
		Cfg.Auth.User,
		cred.ValidFrom,
		cred.ExpiresAt,
		cred.Token,
	); err != nil {
		return err
	}
	return nil
}

func supervisorPasswordReset(c *cli.Context) error {
	var (
		err               error
		password, passKey string
		happy             bool
		cred              *auth.Token
	)

	if Cfg.Auth.User == `` {
		fmt.Println(`Please specify for which  account the"+
		" password should be changed.`)
		if Cfg.Auth.User, err = adm.Read(`user`); err != nil {
			return err
		}
	} else {
		fmt.Printf("Starting with password update of account '%s'"+
			" in 2 seconds.\n", Cfg.Auth.User)
		fmt.Printf(`Use --user flag to switch account account.`)
		time.Sleep(2 * time.Second)
	}
	if strings.Contains(Cfg.Auth.User, `:`) {
		return fmt.Errorf(`Usernames must not contain : character.`)
	}

	fmt.Printf("\nPlease provide the new password you want to set.\n")
password_read:
	password = adm.ReadVerified(`password`)

	if happy, err = adm.EvaluatePassword(3, password,
		Cfg.Auth.User, `soma`); err != nil {
		return err
	} else if !happy {
		goto password_read
	}

	if c.Bool(`reset`) {
		fmt.Printf("\nTo confirm that you are allowed to reset" +
			" this account, an additional" +
			"credential is required.\n")

		switch Cfg.Activation {
		case `ldap`:
			fmt.Println(`Please provide your LDAP password to` +
				`establish ownership.`)
			passKey = adm.ReadVerified(`password`)
		case `mailtoken`:
			fmt.Println(`Please provide the token you received` +
				`via email.`)
			passKey = adm.ReadVerified(`token`)
		default:
			return fmt.Errorf(`Unknown activation mode`)
		}
	} else {
		fmt.Printf("\nPlease provide your currently active/old" +
			" password.\n")
		passKey = adm.ReadVerified(`password`)
	}

	if cred, err = adm.ChangeAccountPassword(
		Client,
		c.Bool(`reset`),
		&auth.Token{
			UserName: Cfg.Auth.User,
			Password: password,
			Token:    passKey,
		},
	); err != nil {
		return err
	}

	// validate received token
	if err = adm.ValidateToken(
		Client,
		Cfg.Auth.User,
		cred.Token,
	); err != nil {
		return err
	}
	// save received token
	if err = store.SaveToken(
		Cfg.Auth.User,
		cred.ValidFrom,
		cred.ExpiresAt,
		cred.Token,
	); err != nil {
		return err
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
