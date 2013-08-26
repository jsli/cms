package models

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"fmt"
	"labix.org/v2/mgo/bson"
)

const (
	DEBUG = true
)

const (
	ROLE_SUPERUSER = 1 // superuser is the KING. It should be created manually, and it cannot be deleted.
	ROLE_ADMIN     = 2
	ROLE_NORMAL    = 3
	ROLE_ANONYMOUS = 4
)

const (
	DbName         = "account"
	CollectionName = "user"
)

const (
	SuperUserName  = "SuperUser"
	SuperUserPwd   = "lijinsong"
	SuperUserEmail = "manson.li3307@gmail.com"
)

const (
	POWER_EDIT_ADMIN_USER  = "edit_admin_user"
	POWER_EDIT_NORMAL_USER = "edit_normal_user"
	POWER_LIST_USERS = "list_users"
)

/*
 *key indicate power's name,
 *value indicate power's description
 *check key only for accessibility
 */
type Power map[string]string

/*
 * TODO: Index and Primary Key???
 */
type User struct {
	Id           bson.ObjectId `bson:"_id"`
	UserName     string        `bson:"user_name"`
	HashPassword []byte        `bson:"password"`
	Role         int           `bson:"role"`
	Email        string        `bson:"email"`
	PowerMap     Power         `bson:"power"`
	IsLogined    bool          `bson:"is_logined"`
}

func (user *User) String() string {
	return fmt.Sprintf("User(id = %s, username = %s, role = %d)\n",
		user.Id.Hex(), user.UserName, user.Role)
}

/*
 *permission check!
 * TODO:
 * 		check power map!!!
 */
func (user *User) IsAdmin() bool {
	return user.Role == ROLE_ADMIN || user.Role == ROLE_SUPERUSER
}

/*
 *TODO: check map!!!
 */
func CheckPermission(user *User, perm_key string) (bool, error) {
	var msg string
	power := user.PowerMap
	access := (power[perm_key] == perm_key)
	if !access {
		msg = fmt.Sprintf("Permission Denied : %s miss permission %s", user, perm_key)
		fmt.Println(msg)
		return access, errors.New(msg)
	} else {
		msg = fmt.Sprintf("Permission Passed : %s has permission %s", user, perm_key)
		fmt.Println(msg)
		return access, nil
	}
}

func GeneratePwdByte(pwd string) []byte {
	pwdByte, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return pwdByte
}
