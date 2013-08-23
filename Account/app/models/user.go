package models

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"fmt"
	"labix.org/v2/mgo"
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
	return fmt.Sprintf("User(username = %s, role = %d)\n",
		user.UserName, user.Role)
}

/*---------------query-----------------*/
//TODO:
//	support selection query
func (user *User) ListUsers(session *mgo.Session, page int, count int) ([]*User, error) {
	if !user.IsAdmin() {
		fmt.Println("Permission Denied!")
		return nil, errors.New("Permission Denied!")
	}
	uc := session.DB(DbName).C(CollectionName)
	query := uc.Find(nil)
	if query != nil {
		results := []*User{}
		query.Skip((page - 1) * count).Limit(count).All(&results)
		if len(results) > 0 {
			fmt.Println("--->>>dump user list: ", len(results), "<<<---")
			for _, item := range results {
				fmt.Println(item)
			}
			return results, nil
		}
	}
	return nil, errors.New("Cannot find users!")
}

func (user *User) GetUserById(session *mgo.Session, id bson.ObjectId) (*User, error) {
	if !user.IsAdmin() {
		fmt.Println("Permission Denied!")
		return nil, errors.New("Permission Denied!")
	}
	uc := session.DB(DbName).C(CollectionName)
	result := User{}
	err := uc.FindId(id).One(&result)
	if err == nil {
		fmt.Println("find user : ", result)
		return &result, nil
	}
	fmt.Println(fmt.Sprintf("Cannot find user by id: %s", id))
	return nil, errors.New(fmt.Sprintf("Cannot find user by id: %s", id))
}

func (user *User) GetUserByName(session *mgo.Session, name string) (*User, error) {
	if !user.IsAdmin() {
		fmt.Println("Permission Denied!")
		return nil, errors.New("Permission Denied!")
	}
	return getUserByM(session, bson.M{"user_name": name})
}

func (user *User) GetUserByEmail(session *mgo.Session, email string) (*User, error) {
	if !user.IsAdmin() {
		fmt.Println("Permission Denied!")
		return nil, errors.New("Permission Denied!")
	}
	return getUserByM(session, bson.M{"email": email})
}

func getUserByM(session *mgo.Session, m bson.M) (*User, error) {
	uc := session.DB(DbName).C(CollectionName)
	result := User{}
	err := uc.Find(m).One(&result)
	if err == nil {
		fmt.Println("find user : ", result)
		return &result, nil
	}
	return nil, errors.New(fmt.Sprintf("Cannot find user by %s", m))
}

/*------------insert-------------*/


func (user *User) SaveUser(session *mgo.Session) error {
	if DEBUG {
		fmt.Println("SaveUser in ------> User")
	}
	uc := session.DB(DbName).C(CollectionName)

	i, _ := uc.Find(bson.M{"user_name": user.UserName}).Count()
	if i != 0 {
		return errors.New(fmt.Sprintf("Duplicated User: %s\n", user.UserName))
	}

	i, _ = uc.Find(bson.M{"email": user.Email}).Count()
	if i != 0 {
		return errors.New(fmt.Sprintf("Duplicated Email: %s\n", user.Email))
	}

	user.Id = bson.NewObjectId()
	err := uc.Insert(user)
	if err != nil {
		return errors.New(fmt.Sprintf("Register User failed: %s\n", user.UserName))
	} else {
		return nil
	}
}

func (user *User) IsAdmin() bool {
	return user.Role == ROLE_ADMIN || user.Role == ROLE_SUPERUSER
}

func generatePwdByte(pwd string) []byte {
	pwdByte, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return pwdByte
}

func GetUserByName(session *mgo.Session, name string) *User {
	user := User{}
	uc := session.DB(DbName).C(CollectionName)
	err := uc.Find(bson.M{"user_name": name}).One(&user)
	if err == nil {
		return &user
	}
	return nil
}

func GetUserByEmail(session *mgo.Session, email string) *User {
	user := User{}
	uc := session.DB(DbName).C(CollectionName)
	err := uc.Find(bson.M{"email": email}).One(&user)
	if err == nil {
		return &user
	}
	return nil
}

func UpdateUser(session *mgo.Session, user User) error {
	uc := session.DB(DbName).C(CollectionName)
	_, err := uc.Upsert(bson.M{"user_name": user.UserName}, user)
	if err == nil {
		fmt.Println("update ", user.UserName, " ok")
	} else {
		fmt.Println("update ", user.UserName, "failed :", err)
	}
	return err
}
