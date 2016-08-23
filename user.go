package hodor

// User #
type User interface {
	GetID() string
	SetID(string)
	GetLogin() string
	SetLogin(string)
	GetEmail() string
	SetEmail(string)
	GetPassword() string
	SetPassword(string)

	GetUserRoles() []UserRole
	GetIndividualUserRights() []UserRight
}

// UserRole #
type UserRole interface {
	GetName() string
	SetName() string
	GetRights() []UserRight
}

// UserRight #
type UserRight interface {
	GetName() string
	SetName(string)
}
