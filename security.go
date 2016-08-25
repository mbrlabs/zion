package hodor

import (
	"net/http"
)

type HodorSecurity struct {
	Strategy   SecurityStrategy
	Middleware Middleware
}

func (s *HodorSecurity) Authenticate() HandlerFunc {
	return s.Strategy.Authenticate()
}

type SecurityStrategy interface {
	Authenticate() HandlerFunc
}

// ============================================================================
// 							Local strategy + middleware
// ============================================================================

type LocalSecurityStrategy struct {
	userStore       UserStore
	sessionStore    SessionStore
	successRedirect string
	failureRedirect string
	loginNameField  string
	passwordField   string
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

func (ls *LocalSecurityStrategy) SetRedirects(successRedirect string, failureRedirect string) {
	ls.failureRedirect = failureRedirect
	ls.successRedirect = successRedirect
}

func (ls *LocalSecurityStrategy) SetPostParameterFields(loginNameField string, passwordField string) {
	ls.loginNameField = loginNameField
	ls.passwordField = passwordField
}

func (ls *LocalSecurityStrategy) Authenticate() HandlerFunc {
	return func(ctx *Context) {
		login := ctx.Request.FormValue(ls.loginNameField)
		password := ctx.Request.FormValue(ls.passwordField)

		// handle empty input
		if len(login) == 0 || len(password) == 0 {
			http.Redirect(ctx.Writer, ctx.Request, ls.failureRedirect, http.StatusOK)
			return
		}

		// get user
		user := ls.userStore.GetUserByLogin(login)
		if user == nil {
			http.Redirect(ctx.Writer, ctx.Request, ls.failureRedirect, http.StatusOK)
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
				http.SetCookie(ctx.Writer, cookie)
				// redirect to succcess page
				http.Redirect(ctx.Writer, ctx.Request, ls.successRedirect, http.StatusOK)
				return
			}
		}

		// redirect to error page
		http.Redirect(ctx.Writer, ctx.Request, ls.failureRedirect, http.StatusOK)
	}
}

// LocalSecurityMiddleware #
type LocalSecurityMiddleware struct {
	userStore    UserStore
	sessionStore SessionStore
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

func (sm *LocalSecurityMiddleware) Execute(ctx *Context) bool {
	cookie, err := ctx.Request.Cookie(sessionCookieName)
	if err != nil {
		http.NotFound(ctx.Writer, ctx.Request)
		return false
	}

	session := sm.sessionStore.Find(cookie.Value)
	if session == nil {
		http.NotFound(ctx.Writer, ctx.Request)
		return false
	}

	user := sm.userStore.GetUserByID(session.UserID)
	if user == nil {
		http.NotFound(ctx.Writer, ctx.Request)
		return false
	}

	ctx.User = user
	return true
}

func (sm *LocalSecurityMiddleware) Name() string {
	return "LocalSecurityMiddleware"
}
