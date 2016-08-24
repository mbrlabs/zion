package hodor

// ============================================================================
// 						interface User & struct HodorUser
// ============================================================================

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

	GetUserRoles() []*UserRole
	GetIndividualUserRights() []*UserRight
}

type HodorUser struct {
	id       string
	login    string
	email    string
	password string
	roles    []*UserRole
	rights   []*UserRight
}

func NewHodorUser(id string, email string, login string, password string) User {
	return &HodorUser{
		id:       id,
		login:    login,
		email:    email,
		password: password,
	}
}

func (u *HodorUser) GetID() string {
	return u.id
}

func (u *HodorUser) SetID(id string) {
	u.id = id
}

func (u *HodorUser) GetPassword() string {
	return u.password
}

func (u *HodorUser) SetPassword(pwd string) {
	u.password = pwd
}

func (u *HodorUser) GetLogin() string {
	return u.login
}

func (u *HodorUser) SetLogin(login string) {
	u.login = login
}

func (u *HodorUser) GetEmail() string {
	return u.email
}

func (u *HodorUser) SetEmail(email string) {
	u.email = email
}

func (u *HodorUser) GetUserRoles() []*UserRole {
	return u.roles
}

func (u *HodorUser) GetIndividualUserRights() []*UserRight {
	return u.rights
}

// ============================================================================
// 						structs UserRole & UserRight
// ============================================================================

// UserRole #
type UserRole struct {
	Name   string
	Rights []*UserRight
}

// UserRight #
type UserRight struct {
	Name string
}
