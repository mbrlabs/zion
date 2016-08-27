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

import "github.com/mbrlabs/hodor"

// UserStore #
type UserStore interface {
	GetUserByLogin(string) (hodor.User, error)
	GetUserByID(string) (hodor.User, error)
	Authenticate(hodor.User, string) bool
}

type MemoryUserStore struct {
	loginToUser map[string]hodor.User
	idToUser    map[string]hodor.User
}

func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{
		loginToUser: make(map[string]hodor.User),
		idToUser:    make(map[string]hodor.User),
	}
}

func (us *MemoryUserStore) GetUserByLogin(login string) (hodor.User, error) {
	return us.loginToUser[login], nil
}

func (us *MemoryUserStore) GetUserByID(id string) (hodor.User, error) {
	return us.idToUser[id], nil
}

func (us *MemoryUserStore) AddUser(user hodor.User) {
	us.loginToUser[user.GetLogin()] = user
	us.idToUser[user.GetID()] = user
}

func (us *MemoryUserStore) Authenticate(user hodor.User, password string) bool {
	return user.GetPassword() == password
}
