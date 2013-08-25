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
	return fmt.Sprintf("User(id = %s, username = %s, role = %d)\n",
		user.Id.Hex(), user.UserName, user.Role)
}

/*---------------query-----------------*/
//TODO:
//	support selection query
func (user *User) ListUsers(session *mgo.Session, page int, count int, role int) ([]*User, error) {
	var perm_key string
	switch role {
		case ROLE_NORMAL:
			perm_key = POWER_EDIT_NORMAL_USER
		case ROLE_ADMIN:
			perm_key = POWER_EDIT_ADMIN_USER
		default:
			perm_key = POWER_EDIT_ADMIN_USER
	}
	check := CheckPermission(user, perm_key)
	if !check {
		return nil, errors.New("Permission Denied")
	}

	fmt.Println("list users: role = ", role, "count = ", count, "page = ", page)
	uc := session.DB(DbName).C(CollectionName)
	query := uc.Find(bson.M{"role":role})
	if query != nil {
		results := []*User{}
		query.Skip((page - 1) * count).Limit(count).All(&results)
		if len(results) > 0 {
			fmt.Println("--->>>dump user list: ", len(results), "<<<---")
			for _, item := range results {
				fmt.Println(item)
			}
			return results, nil
		} else {
			fmt.Println("empty list")
		}
	}
	return nil, errors.New(fmt.Sprintf("Cannot find users! from %d to %d", (page-1)*count+1, page*count))
}

func (user *User) LoadSelf(session *mgo.Session) error {
	dbUser, err := getUserByM(session, bson.M{"user_name": user.UserName})
	if dbUser != nil && err == nil {
		user.Id = dbUser.Id
		user.Email = dbUser.Email
		user.HashPassword = dbUser.HashPassword
		user.IsLogined = dbUser.IsLogined
		user.PowerMap = dbUser.PowerMap
		user.Role = dbUser.Role
		return nil
	}
	return err
}

func GetUserById(session *mgo.Session, id string) (*User, error) {
	return getUserByM(session, bson.M{"_id": bson.ObjectIdHex(id)})
}

func GetUserByName(session *mgo.Session, name string) (*User, error) {
	return getUserByM(session, bson.M{"user_name": name})
}

func GetUserByEmail(session *mgo.Session, email string) (*User, error) {
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
	fmt.Println(fmt.Sprintf("Cannot find user by %s", m))
	return nil, errors.New(fmt.Sprintf("Cannot find user by %s", m))
}

/*------------insert-------------*/
func (user *User) SaveUser(session *mgo.Session) error {
	//dont check permission, all kinds of user can register
	fmt.Println("SaveUser in ------> User")

	if can, err := canBeSaved(session, user); !can {
		return err
	}
	user.Id = bson.NewObjectId()
	return saveUser(session, user)
}

func canBeSaved(session *mgo.Session, user *User) (bool, error) {
	var err error = nil
	var can bool = true
	uc := session.DB(DbName).C(CollectionName)
	//duplicated user check again!!!
	i, _ := uc.Find(bson.M{"user_name": user.UserName}).Count()
	if i != 0 {
		err = errors.New(fmt.Sprintf("Duplicated User: %s\n", user.UserName))
		can = false
	} else {
		//duplicated email check again!!!
		i, _ = uc.Find(bson.M{"email": user.Email}).Count()
		if i != 0 {
			err = errors.New(fmt.Sprintf("Duplicated Email: %s\n", user.Email))
			can = false
		}
	}
	if !can {
		fmt.Println("cannot save user ", user)
		fmt.Println("err_msg: ", err)
	}
	return can, err
}

func saveUser(session *mgo.Session, obj interface{}) error {
	uc := session.DB(DbName).C(CollectionName)
	err := uc.Insert(obj)
	if err != nil {
		fmt.Println("cannot insert ", obj.(string), ", cause of ", err)
		return errors.New(fmt.Sprintf("Save User failed: %s\n", obj.(string)))
	} else {
		return nil
	}
}

/*------------delete--------------*/
func (user *User) DeleteUserById(session *mgo.Session) error {
	return deleteUserByM(session, bson.M{"_id": user.Id})
}

func deleteUserByM(session *mgo.Session, m bson.M) error {
	uc := session.DB(DbName).C(CollectionName)
	err := uc.Remove(m)
	if err != nil {
		fmt.Println(fmt.Sprintf("Cannot delete user by %s", m))
	} else {
		fmt.Println(fmt.Sprintf("delete user by %s", m))
	}
	return err
}

/*----------update------------*/
func (user *User) UpdateUser(session *mgo.Session) error {
	//need check permission here
	return updateUserByM(session, bson.M{"user_name": user.UserName}, user)
}

func updateUserByM(session *mgo.Session, m bson.M, obj interface{}) error {
	uc := session.DB(DbName).C(CollectionName)
	err := uc.Update(m, obj)
	if err == nil {
		fmt.Println("update ", obj, " ok")
	} else {
		fmt.Println("update ", obj, "failed :", err)
	}
	return err
}

/*
 *permission check!
 * TODO:
 * 		check power map!!!
 */
func (user *User) IsAdmin() bool {
	return user.Role == ROLE_ADMIN || user.Role == ROLE_SUPERUSER
}

func CheckPermission(user *User, perm_key string) bool {
	power := user.PowerMap
	access := (power[perm_key] == perm_key)
	if !access {
		fmt.Println("Permission Denied : ", user, " miss permission ", perm_key)
	} else {
		fmt.Println("Permission Passed : ", user, " has permission ", perm_key)
	}
	return access
}

func GeneratePwdByte(pwd string) []byte {
	pwdByte, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return pwdByte
}
