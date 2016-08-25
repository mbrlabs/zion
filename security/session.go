package security

import (
	"github.com/mbrlabs/hodor"
	"time"
)

const (
	sessionExpire     = 24 * 3 * time.Hour
	sessionLength     = 64
	sessionCookieName = "hsession"
	sessionAlphabet   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// Session #TODO
type Session struct {
	ID     string
	UserID string
	Expire time.Time
}

func NewSession(user hodor.User) *Session {
	return &Session{
		ID:     hodor.GenerateRandomString(sessionLength, sessionAlphabet),
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
	return s.sessions[id]
}

func (s MemorySessionStore) Save(session *Session) error {
	s.sessions[session.ID] = session
	return nil
}

func (s MemorySessionStore) Delete(session *Session) error {
	delete(s.sessions, session.ID)
	return nil
}
