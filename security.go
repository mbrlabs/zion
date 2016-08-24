package hodor

import (
	"fmt"
	"net/http"
	"time"
)

const (
	sessionExpire     = 24 * 3 * time.Hour
	sessionLength     = 64
	sessionCookieName = "hsession"
)

// AuthenticateHandlerFunc #TODO
func authenticateHandlerFunc(hodor *Hodor, loginFieldName string, passwordFieldName,
	successPath string, errorPath string) HandlerFunc {

	return func(ctx *Context) {
		login := ctx.Request.FormValue(loginFieldName)
		password := ctx.Request.FormValue(passwordFieldName)

		// handle empty input
		if len(login) == 0 || len(password) == 0 {
			http.Redirect(ctx.Writer, ctx.Request, errorPath, http.StatusOK)
			return
		}

		// get user
		user := hodor.UserStore.GetUserByLogin(login)
		if user == nil {
			http.Redirect(ctx.Writer, ctx.Request, errorPath, http.StatusOK)
			return
		}

		// authenticate user
		if hodor.UserStore.Authenticate(user, password) {
			// create new session
			session := NewSession(user)
			err := hodor.SessionStore.Save(session)
			if err == nil {
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
		ID:     generateRandomString(sessionLength, alphabetAlphaNum),
		UserID: user.GetID(),
		Expire: time.Now().Add(sessionExpire),
	}
}

// ============================================================================
// 					interface UserStore & struct MemoryUserStore
// ============================================================================

// UserStore #
type UserStore interface {
	GetUserByLogin(string) User
	Authenticate(User, string) bool
}

type MemoryUserStore struct {
	users map[string]User
}

func NewMemoryUserStore() UserStore {
	users := make(map[string]User)
	user := NewHodorUser("123", "root@hodor.com", "root@hodor.com", "hodor")
	users[user.GetLogin()] = user
	return MemoryUserStore{users: users}
}

func (s MemoryUserStore) GetUserByLogin(login string) User {
	return s.users[login]
}

func (s MemoryUserStore) Authenticate(user User, password string) bool {
	fmt.Printf("%s store: %s | user: %s", user.GetLogin(), user.GetPassword(), password)
	return user.GetPassword() == password
}

// ============================================================================
// 					interface SessionStore & Memory UserStore
// ============================================================================

// SessionStore #
type SessionStore interface {
	Find(string) *Session
	Save(*Session) error
	Delete(*Session) error
}

type MemorySessionStore struct {
	sessions map[string]*Session
}

func NewMemorySessionStore() SessionStore {
	return MemorySessionStore{sessions: make(map[string]*Session)}
}

func (s MemorySessionStore) Find(id string) *Session {
	return s.Find(id)
}

func (s MemorySessionStore) Save(session *Session) error {
	s.sessions[session.ID] = session
	return nil
}

func (s MemorySessionStore) Delete(session *Session) error {
	delete(s.sessions, session.ID)
	return nil
}
