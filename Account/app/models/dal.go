package models

import (
//	"labix.org/v2/mgo/bson"
)

type UserCrud interface {
	ListUsers(host *User, page int, count int, role int) ([]*User, error)

	GetUserById(id string) (*User, error)
	GetUserByName(name string) (*User, error)
	GetUserByEmail(email string) (*User, error)

	SaveUser(user *User) error

	DeleteUserById(host *User, id string) error

	UpdateUserById(host *User, user *User) error
}
