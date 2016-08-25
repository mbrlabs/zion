package hodor

import (
	"time"
)

const (
	sessionExpire     = 24 * 3 * time.Hour
	sessionLength     = 64
	sessionCookieName = "hsession"
)

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
// 					interface SessionStore & MemorySessionStore
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
