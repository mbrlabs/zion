// Copyright (c) 2016. See AUTHORS file.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package security

import (
	"fmt"
	"github.com/mbrlabs/zion"
	"net/http"
)

// LocalSecurityStrategy
//------------------------------------------------------------------------------------

type LocalSecurityStrategy struct {
	userStore      UserStore
	sessionStore   SessionStore
	loginNameField string
	passwordField  string

	logoutRedirect  string
	failureRedirect string
	successRedirect string

	logoutHandler  zion.HandlerFunc
	failureHandler zion.HandlerFunc
	successHandler zion.HandlerFunc
}

func NewLocalSecurityStrategy(userStore UserStore, sessionStore SessionStore) *LocalSecurityStrategy {
	return &LocalSecurityStrategy{
		userStore:    userStore,
		sessionStore: sessionStore,
	}
}

func (ls *LocalSecurityStrategy) SetUserStore(store UserStore) {
	ls.userStore = store
}

func (ls *LocalSecurityStrategy) SetSessionStore(store SessionStore) {
	ls.sessionStore = store
}

func (ls *LocalSecurityStrategy) SetRedirects(successRedirect string, failureRedirect string, logoutRedirect string) {
	ls.failureRedirect = failureRedirect
	ls.successRedirect = successRedirect
	ls.logoutRedirect = logoutRedirect
}

func (ls *LocalSecurityStrategy) SetSuccessHandler(onSuccess zion.HandlerFunc) {
	ls.successHandler = onSuccess
}

func (ls *LocalSecurityStrategy) SetFailureHandler(onFailure zion.HandlerFunc) {
	ls.failureHandler = onFailure
}

func (ls *LocalSecurityStrategy) SetLogoutHandler(onLogout zion.HandlerFunc) {
	ls.logoutHandler = onLogout
}

func (ls *LocalSecurityStrategy) SetPostParameterFields(loginNameField string, passwordField string) {
	ls.loginNameField = loginNameField
	ls.passwordField = passwordField
}

func (ls *LocalSecurityStrategy) handleLoginError(ctx *zion.Context) {
	if ls.failureHandler != nil {
		ls.failureHandler(ctx)
	} else if len(ls.failureRedirect) > 0 {
		ctx.Redirect(ls.failureRedirect)
	} else {
		ctx.SendStatus(http.StatusBadRequest)
	}
}

func (ls *LocalSecurityStrategy) handleLoginSuccess(ctx *zion.Context) {
	if ls.successHandler != nil {
		ls.successHandler(ctx)
	} else if len(ls.successRedirect) > 0 {
		ctx.Redirect(ls.successRedirect)
	} else {
		ctx.SendStatus(http.StatusOK)
	}
}

func (ls *LocalSecurityStrategy) Authenticate() zion.HandlerFunc {
	return func(ctx *zion.Context) {
		login := ctx.FormValue(ls.loginNameField)
		password := ctx.FormValue(ls.passwordField)

		// handle empty input
		if len(login) == 0 || len(password) == 0 {
			ls.handleLoginError(ctx)
			return
		}

		// get user
		user, err := ls.userStore.GetUserByLogin(login)
		if user == nil || err != nil {
			ls.handleLoginError(ctx)
			return
		}

		// authenticate user
		if ls.userStore.Authenticate(user, password) {
			// create new session
			session := NewSession(user)
			err := ls.sessionStore.Save(session)
			if err == nil {
				// set cockie
				cookie := &http.Cookie{
					Name:    sessionCookieName,
					Value:   session.ID,
					Expires: session.Expire,
				}
				ctx.SetCookie(cookie)
				ctx.AddExtra(zion.ExtraUser, user)
				ls.handleLoginSuccess(ctx)
				return
			}
		}

		ls.handleLoginError(ctx)
	}
}

func (ls *LocalSecurityStrategy) Logout() zion.HandlerFunc {
	return func(ctx *zion.Context) {
		cookie, err := ctx.Cookie(sessionCookieName)
		if err == nil {
			session := ls.sessionStore.Find(cookie.Value)
			if session != nil {
				ls.sessionStore.Delete(session)
				ctx.Redirect(ls.logoutRedirect)
			}
		}
	}
}

// LocalSecurityMiddleware
//------------------------------------------------------------------------------------

// LocalSecurityMiddleware #
type LocalSecurityMiddleware struct {
	userStore    UserStore
	sessionStore SessionStore
	rules        SecurityRules
	redirect     string
}

func NewLocalSecurityMiddleware(userStore UserStore, sessionStore SessionStore) *LocalSecurityMiddleware {
	return &LocalSecurityMiddleware{
		userStore:    userStore,
		sessionStore: sessionStore,
	}
}

func (ls *LocalSecurityMiddleware) SetUserStore(store UserStore) {
	ls.userStore = store
}

func (ls *LocalSecurityMiddleware) SetSessionStore(store SessionStore) {
	ls.sessionStore = store
}

func (ls *LocalSecurityMiddleware) SetRedirect(path string) {
	ls.redirect = path
}

func (ls *LocalSecurityMiddleware) AddRule(rule SecurityRule) {
	ls.rules = append(ls.rules, rule)
}

func (sm *LocalSecurityMiddleware) redirectOnAuthFailed(ctx *zion.Context) {
	if len(sm.redirect) == 0 {
		ctx.SendStatus(http.StatusForbidden)
	} else {
		ctx.Redirect(sm.redirect)
	}
}

func (sm *LocalSecurityMiddleware) Execute(ctx *zion.Context) bool {
	// get cookie from request header
	cookie, err := ctx.Cookie(sessionCookieName)
	if err != nil {
		if sm.rules.IsAllowed(nil, ctx) {
			return true
		}
		fmt.Println("User with no cookie set tries to acacess restricted page")
		sm.redirectOnAuthFailed(ctx)
		return false
	}

	// get session based on session key in cookie
	session := sm.sessionStore.Find(cookie.Value)
	var user User

	// get user by userID stored in session
	if session != nil {
		var err error
		user, err = sm.userStore.GetUserByID(session.UserID)
		if err != nil {
			user = nil
		}
	}

	// go through all security rules
	if sm.rules.IsAllowed(user, ctx) {
		ctx.AddExtra(zion.ExtraUser, user)
		return true
	}

	sm.redirectOnAuthFailed(ctx)
	return false
}

func (sm *LocalSecurityMiddleware) Name() string {
	return "LocalSecurityMiddleware"
}
