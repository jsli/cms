package models

import (
	"github.com/robfig/revel"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"fmt"
)

type RegUser struct {
	User
	PasswordStr   string
	ConfirmPwdStr string
}

func (regUser *RegUser) SaveUser(session *mgo.Session) error {
	if DEBUG {
		fmt.Println("SaveUser in ------> RegUser")
	}
	regUser.HashPassword = GeneratePwdByte(regUser.PasswordStr)
	dal := NewDalMgo(session)
	err := dal.SaveUser(&regUser.User)
	return err
}

func (regUser *RegUser) Validate(v *revel.Validation, session *mgo.Session) {
	//Check workflow:
	//see @validation.go Check(obj interface{}, checks ...Validator)
	//Validator is an interface, v.Check invoke v.Apply for each validator.
	//Further, v.Apply invoke validator.IsSatisfied with passing obj.
	//Checking result is an object of ValidationResult. The field Ok of ValidationResult
	//would be true if checking success. Otherwise, Ok would be false, and another filed
	//Error of ValidationResult would be non-nil, an ValidationError filled with error message
	//should be assigned to Error.
	v.Check(regUser.UserName,
		revel.Required{},
		revel.Match{USERNAME_REX},
		DuplicatedUserValidator{session},
	)

	//validation provide an convenient method for checking Email.
	//revel has a const for email rexgep, Email will use the rex to check email string.
	v.Email(regUser.Email)
	v.Check(regUser.Email,
		DuplicatedEmailValidator{session},
	)

	v.Check(regUser.PasswordStr,
		revel.Required{},
		revel.Match{PWD_REX},
	)
	v.Check(regUser.ConfirmPwdStr,
		revel.Required{},
		revel.Match{PWD_REX},
	)
	//pwd and comfirm_pwd should be equal
	v.Required(regUser.PasswordStr == regUser.ConfirmPwdStr).Message("The passwords do not match.")
}

/*
 * a validator for checking duplicated user
 */
type DuplicatedUserValidator struct {
	session *mgo.Session
}

func (dup DuplicatedUserValidator) IsSatisfied(obj interface{}) bool {
	user := User{}
	uc := dup.session.DB(DbName).C(CollectionName)
	err := uc.Find(bson.M{"user_name": obj.(string)}).One(&user)
	if err != nil {
		return true
	}
	return false
}

func (dup DuplicatedUserValidator) DefaultMessage() string {
	return "Duplicated User"
}

/*
 * a validator for checking duplicated email
 */
type DuplicatedEmailValidator struct {
	session *mgo.Session
}

func (dup DuplicatedEmailValidator) IsSatisfied(obj interface{}) bool {
	user := User{}
	uc := dup.session.DB(DbName).C(CollectionName)
	err := uc.Find(bson.M{"email": obj.(string)}).One(&user)
	if err != nil {
		return true
	}
	return false
}

func (dup DuplicatedEmailValidator) DefaultMessage() string {
	return "Duplicated Email"
}