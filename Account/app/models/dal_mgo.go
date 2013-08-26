package models

import (
	"errors"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type DalMgo struct {
	Session *mgo.Session
}

func NewDalMgo(session *mgo.Session) *DalMgo {
	return &DalMgo{session}
}

func (dal *DalMgo) ListUsers(host *User, page int, count int, role int) ([]*User, error) {
	var perm_key string
	switch role {
	case ROLE_NORMAL:
		perm_key = POWER_EDIT_NORMAL_USER
	case ROLE_ADMIN:
		perm_key = POWER_EDIT_ADMIN_USER
	default:
		perm_key = POWER_EDIT_ADMIN_USER
	}
	if access, err := CheckPermission(host, perm_key); !access {
		return nil, err
	}

	fmt.Println("list users: role = ", role, "count = ", count, "page = ", page)
	uc := dal.Session.DB(DbName).C(CollectionName)
	query := uc.Find(bson.M{"role": role})
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
	return nil, errors.New(fmt.Sprintf("Cannot list users! from %d to %d", (page-1)*count+1, page*count))
}

func (dal *DalMgo) GetUserById(id string) (*User, error) {
	return GetUserByM(dal.Session, bson.M{"_id": bson.ObjectIdHex(id)})
}

func (dal *DalMgo) GetUserByName(name string) (*User, error) {
	return GetUserByM(dal.Session, bson.M{"user_name": name})
}

func (dal *DalMgo) GetUserByEmail(email string) (*User, error) {
	return GetUserByM(dal.Session, bson.M{"email": email})
}

func GetUserByM(session *mgo.Session, m bson.M) (*User, error) {
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

func (dal *DalMgo) SaveUser(user *User) error {
	//dont check permission, all kinds of user can register
	fmt.Println("SaveUser in ------> User")

	if can, err := dal.CanBeSaved(user); !can {
		return err
	}
	user.Id = bson.NewObjectId()
	return saveUser(dal.Session, user)
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

func (dal *DalMgo) CanBeSaved(user *User) (bool, error) {
	var err error = nil
	var can bool = true
	uc := dal.Session.DB(DbName).C(CollectionName)
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

func (dal *DalMgo) DeleteUserById(host *User, id string) error {
	var perm_key string
	switch host.Role {
	case ROLE_NORMAL:
		perm_key = POWER_EDIT_NORMAL_USER
	case ROLE_ADMIN:
		perm_key = POWER_EDIT_ADMIN_USER
	default:
		perm_key = POWER_EDIT_ADMIN_USER
	}
	if access, err := CheckPermission(host, perm_key); !access {
		return err
	}

	return DeleteUserByM(dal.Session, host, bson.M{"_id": bson.ObjectIdHex(id)})
}

func DeleteUserByM(session *mgo.Session, host *User, m bson.M) error {
	uc := session.DB(DbName).C(CollectionName)
	err := uc.Remove(m)
	if err != nil {
		fmt.Println(fmt.Sprintf("Cannot delete user by %s", m))
	} else {
		fmt.Println(fmt.Sprintf("delete user by %s", m))
	}
	return err
}

func (dal *DalMgo) UpdateUserById(host *User, user *User) error {
	if host.UserName == user.UserName {
	} else {
		var perm_key string
		switch host.Role {
		case ROLE_NORMAL:
			perm_key = POWER_EDIT_NORMAL_USER
		case ROLE_ADMIN:
			perm_key = POWER_EDIT_ADMIN_USER
		default:
			perm_key = POWER_EDIT_ADMIN_USER
		}
		if access, err := CheckPermission(host, perm_key); !access {
			return err
		}
	}

	return UpdateUserByM(dal.Session, bson.M{"_id": user.Id}, user)
}

func UpdateUserByM(session *mgo.Session, m bson.M, obj interface{}) error {
	uc := session.DB(DbName).C(CollectionName)
	err := uc.Update(m, obj)
	if err == nil {
		fmt.Println("update ", obj, " ok")
	} else {
		fmt.Println("update ", obj, "failed :", err)
	}
	return err
}
