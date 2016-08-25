package security

import "github.com/mbrlabs/hodor"

// UserStore #
type UserStore interface {
	GetUserByLogin(string) hodor.User
	GetUserByID(string) hodor.User
	Authenticate(hodor.User, string) bool
}
