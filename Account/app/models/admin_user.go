package models

import (
	"labix.org/v2/mgo/bson"
)

/*
 * include SuperUser and Admin
 */
type AdminUser struct {
	User
}

func NewAdminUser() (*AdminUser, error) {
	power := Power{POWER_EDIT_NORMAL_USER: POWER_EDIT_NORMAL_USER}
	user := AdminUser{User{Id: bson.NewObjectId(), Role: ROLE_ADMIN, PowerMap: power}}
	return &user, nil
}
