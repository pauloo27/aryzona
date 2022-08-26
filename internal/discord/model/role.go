package model

type Role interface {
	ID() string
	Name() string
	Permissions() Permissions
	Position() int
	Color() int
	Mentionable() bool
}
