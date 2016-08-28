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

import "github.com/mbrlabs/zion"

// UserStore #
type UserStore interface {
	GetUserByLogin(string) (zion.User, error)
	GetUserByID(string) (zion.User, error)
	Authenticate(zion.User, string) bool
}

type MemoryUserStore struct {
	loginToUser map[string]zion.User
	idToUser    map[string]zion.User
}

func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{
		loginToUser: make(map[string]zion.User),
		idToUser:    make(map[string]zion.User),
	}
}

func (us *MemoryUserStore) GetUserByLogin(login string) (zion.User, error) {
	return us.loginToUser[login], nil
}

func (us *MemoryUserStore) GetUserByID(id string) (zion.User, error) {
	return us.idToUser[id], nil
}

func (us *MemoryUserStore) AddUser(user zion.User) {
	us.loginToUser[user.GetLogin()] = user
	us.idToUser[user.GetID()] = user
}

func (us *MemoryUserStore) Authenticate(user zion.User, password string) bool {
	return user.GetPassword() == password
}
