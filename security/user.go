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

// User is used for authentication & authorization. See the security package.
type User interface {
	GetID() string
	GetLogin() string
	SetLogin(string)
	GetEmail() string
	SetEmail(string)
	GetPassword() string
	SetPassword(string)

	GetRoles() []string
	AddRole(string)
}

// UserStore
//------------------------------------------------------------------------------------

// UserStore #
type UserStore interface {
	GetUserByLogin(string) (User, error)
	GetUserByID(string) (User, error)
	Authenticate(User, string) bool
}

// MemoryUserStore
//------------------------------------------------------------------------------------

type MemoryUserStore struct {
	loginToUser map[string]User
	idToUser    map[string]User
}

func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{
		loginToUser: make(map[string]User),
		idToUser:    make(map[string]User),
	}
}

func (us *MemoryUserStore) GetUserByLogin(login string) (User, error) {
	return us.loginToUser[login], nil
}

func (us *MemoryUserStore) GetUserByID(id string) (User, error) {
	return us.idToUser[id], nil
}

func (us *MemoryUserStore) AddUser(user User) {
	us.loginToUser[user.GetLogin()] = user
	us.idToUser[user.GetID()] = user
}

func (us *MemoryUserStore) Authenticate(user User, password string) bool {
	return user.GetPassword() == password
}
