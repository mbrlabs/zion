package hodor

import (
	"net/http"
	"time"
)

const (
	sessionExpire     = 24 * 3 * time.Hour
	sessionLength     = 42
	sessionCookieName = "hsession"
)

// AuthenticateHandlerFunc #TODO
func AuthenticateHandlerFunc(userStore UserStore, sessionStore SessionStore,
	loginFieldName string, passwordFieldName, successPath string, errorPath string) HandlerFunc {

	return func(ctx *Context) {
		login := ctx.Request.FormValue(loginFieldName)
		password := ctx.Request.FormValue(passwordFieldName)

		// handle empty input
		if len(login) == 0 || len(password) == 0 {
			http.Redirect(ctx.Writer, ctx.Request, errorPath, http.StatusOK)
			return
		}

		// get user
		user := userStore.GetUserByLogin(login)
		if user == nil {
			http.Redirect(ctx.Writer, ctx.Request, errorPath, http.StatusOK)
			return
		}

		// authenticate user
		if userStore.Authenticate(user, password) {
			// create new session
			session := NewSession(user)
			err := sessionStore.Save(session)
			if err != nil {
				// set cockie
				cookie := &http.Cookie{
					Name:    sessionCookieName,
					Value:   session.ID,
					Expires: session.Expire,
				}
				http.SetCookie(ctx.Writer, cookie)

				// redirect to succcess page
				http.Redirect(ctx.Writer, ctx.Request, successPath, http.StatusOK)
				return
			}
		}

		// redirect to error page
		http.Redirect(ctx.Writer, ctx.Request, errorPath, http.StatusOK)
	}
}

// ============================================================================
// 								struct Session
// ============================================================================

// Session #TODO
type Session struct {
	ID     string
	UserID string
	Expire time.Time
}

func NewSession(user User) *Session {
	return &Session{
		ID:     generateRandomString(sessionLength, alphabetAlphaNumPlus),
		UserID: user.GetID(),
		Expire: time.Now().Add(sessionExpire),
	}
}

// ============================================================================
// 								UserStore
// ============================================================================

// UserStore #
type UserStore interface {
	GetUserByLogin(string) User
	Authenticate(User, string) bool
}

// ============================================================================
// 								SessionStore
// ============================================================================

// SessionStore #
type SessionStore interface {
	Find(string) *Session
	Save(*Session) error
	Delete(*Session) error
}
