package models

import (
	"fmt"
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/robfig/revel"
	"regexp"
	"labix.org/v2/mgo/bson"
)

const (
	DEBUG = true
	DEBUG_PWD = true
)

var USERNAME_REX, PWD_REX, NICKNAME_REX *regexp.Regexp

func init() {
	USERNAME_REX = regexp.MustCompile(`^[a-z0-9_]{6,16}$`)
	PWD_REX = regexp.MustCompile(`^[\x01-\xfe]{8,20}$`)

	//NICKNAME_REX = regexp.MustCompile(`^[a-zA-Z\xa0-\xff_][0-9a-zA-Z\xa0-\xff_]{3,15}$`)
	NICKNAME_REX = USERNAME_REX
}

func generatePwdByte(pwd string) []byte {
	pwdByte, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return pwdByte
}

/*
 * real struct which was persisted in database
 */
type User struct {
	Id bson.ObjectId "_id"
	UserName string
	Email    string
	NickName string
	HashPassword []byte
}

func (user *User) String() string {
	if !DEBUG_PWD {
		return fmt.Sprintf("User(username = %s, email = %s, nick name = %s)",
			user.UserName, user.Email, user.NickName)
	} else {
		return fmt.Sprintf("User(username = %s, email = %s, nick name = %s), pwd = %s",
			user.UserName, user.Email, user.NickName, user.HashPassword)
	}
}

func GetAllUsers() error {
	manager, err := NewDbManager()
	if err != nil {
		fmt.Println("New db manager error")
		return err
	}
	defer manager.Close()
	manager.GetAllUsers()
	return nil
}

func (user *User) SaveUser() error {
	if DEBUG {
		fmt.Println("SaveUser in ------> User")
	}

	manager, err := NewDbManager()
	if err != nil {
		fmt.Println("New db manager error")
		return err
	}
	defer manager.Close()

	err = manager.SaveUser(user)
	if err != nil {
		fmt.Println("save user failed")
		return err
	} else {
		fmt.Println("Save User success: ", user)
	}
	return nil
}

/*
 * used for login
 */
type LoginUser struct {
	User
	PasswordStr string
}

func (loginUser *LoginUser) Validate(v *revel.Validation) {
	v.Check(loginUser.UserName,
		revel.Required{},
		revel.Match{USERNAME_REX},
	).Message("UserName or password is wrong")

	v.Check(loginUser.PasswordStr,
		revel.Required{},
		revel.Match{PWD_REX},
	).Message("UserName or password is wrong")

	//0: generate passing str
	//1: get pwd bytes from database
	//2: compare them
	//test here
	pwd := "testtest"
	//rPwd := "testtest"
	rPwd := "testtest"
	v.Required(pwd == rPwd).Message("user name or password is wrong!!!")
}

/*
 * used for register
 */
type RegUser struct {
	User
	PasswordStr string
	ConfirmPwdStr string
}

func (regUser *RegUser) SaveUser() error {
	if DEBUG {
		fmt.Println("SaveUser in ------> RegUser")
	}
	regUser.HashPassword = generatePwdByte(regUser.PasswordStr)
	err := regUser.User.SaveUser()
	return err
}

func (regUser *RegUser) Validate(v *revel.Validation) {
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
		DuplicatedUser{},
	)

	v.Check(regUser.NickName,
		revel.Required{},
		revel.Match{NICKNAME_REX},
	)

	//validation provide an convenient method for checking Email.
	//revel has a const for email rexgep, Email will use the rex to check email string.
	v.Email(regUser.Email)
	v.Check(regUser.Email,
		DuplicatedEmail{},
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
 *used for updating user
 */
type UpdateUser RegUser

/*
 * a validator for checking duplicated user
 */
type DuplicatedUser struct{}

func (dup DuplicatedUser) IsSatisfied(obj interface{}) bool {
	manager, err := NewDbManager()
	if err != nil {
		fmt.Println("New db manager error")
		return false
	}
	defer manager.Close()
	
	registed := manager.IsUserRegistedByName(obj.(string))
	return !registed
}

func (dup DuplicatedUser) DefaultMessage() string {
	return "Duplicated User"
}

/*
 * a validator for checking duplicated email
 */
type DuplicatedEmail struct{}

func (dup DuplicatedEmail) IsSatisfied(obj interface{}) bool {
	manager, err := NewDbManager()
	if err != nil {
		fmt.Println("New db manager error")
		return false
	}
	defer manager.Close()
	
	registed := manager.IsUserRegistedByEmail(obj.(string))
	return !registed
}

func (dup DuplicatedEmail) DefaultMessage() string {
	return "Duplicated Email"
}
