package models

import (
	"labix.org/v2/mgo/bson"
)

type AnonymousUser struct {
	User
}

func NewAnonymousUser() (*AnonymousUser, error) {
	user := AnonymousUser{User{Id: bson.NewObjectId(), Role: ROLE_ANONYMOUS}}
	user.UserName = user.Id.Hex()
	return &user, nil
}
