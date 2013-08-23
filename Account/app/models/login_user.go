package models

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/robfig/revel"
	"labix.org/v2/mgo"
)

type LoginUser struct {
	User
	PasswordStr string
}

func (loginUser *LoginUser) Validate(v *revel.Validation, session *mgo.Session) {
	v.Check(loginUser.UserName,
		revel.Required{},
		revel.Match{USERNAME_REX},
	).Message("UserName or password is wrong")

	v.Check(loginUser.PasswordStr,
		revel.Required{},
		revel.Match{PWD_REX},
		LegalUserValidator{session, loginUser.UserName},
	).Message("UserName or password is wrong")

	if !v.HasErrors() {
		user := GetUserByName(session, loginUser.UserName)
		user.IsLogined = true
		UpdateUser(session, *user)
	}
}

/*
 * a validator for checking duplicated user
 */
type LegalUserValidator struct {
	session *mgo.Session
	name    string
}

func (legal LegalUserValidator) IsSatisfied(obj interface{}) bool {
	user := GetUserByName(legal.session, legal.name)
	if user != nil {
		err := bcrypt.CompareHashAndPassword(user.HashPassword, []byte(obj.(string)))
		if err == nil {
			return true
		} else {
			return false
		}
	}
	return false
}

func (legal LegalUserValidator) DefaultMessage() string {
	return "Illegal user, user name or password is wrong"
}