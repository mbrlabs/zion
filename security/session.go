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
