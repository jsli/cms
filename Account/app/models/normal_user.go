package models

import (
	//	"fmt"
	//	"code.google.com/p/go.crypto/bcrypt"
	//	"github.com/robfig/revel"
	//	"regexp"
	"labix.org/v2/mgo/bson"
)

type NormalUser struct {
	User
	Available    bool              `bson:"available"`
	PersonalInfo map[string]string `bson:"personal_info"`
}

func NewNormalUser() (*NormalUser, error) {
	user := NormalUser{User{Id: bson.NewObjectId(), Role: ROLE_NORMAL}, false, nil}
	return &user, nil
}
