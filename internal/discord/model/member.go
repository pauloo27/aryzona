package model

type Member interface {
	Roles() []Role
	Permissions() Permissions
}
