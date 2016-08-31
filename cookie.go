package zion

import "time"

type Cookie struct {
	Name     string
	Value    string
	Path     string
	Domain   string
	Expires  time.Time
	Secure   bool
	HTTPOnly bool
}
