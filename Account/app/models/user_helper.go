package models

import (
	"labix.org/v2/mgo/bson"
	"errors"
	"fmt"
)

const (
	DbName = "account"
	CollectionName = "user"
)

func (manager *DbManager) IsUserRegistedByName(name string) bool {
	user, err := manager.GetUserByName(name)
	if err == nil && user != nil {
		return true
	}
	return false
}

func (manager *DbManager) IsUserRegistedByEmail(email string) bool {
	user, err := manager.GetUserByEmail(email)
	if err == nil && user != nil {
		return true
	}
	return false
}

func (manager *DbManager) GetUserByName(name string) (user *User, err error) {
	uc := manager.session.DB(DbName).C(CollectionName)
	err = uc.Find(bson.M{"username":name}).One(&user)
	return
}

func (manager *DbManager) GetUserByEmail(email string) (user *User, err error) {
	uc := manager.session.DB(DbName).C(CollectionName)
	err = uc.Find(bson.M{"email":email}).One(&user)
	return
}

func (manager *DbManager) GetAllUsers() ([]User, error) {
	uc := manager.session.DB(DbName).C(CollectionName)
	count, _ := uc.Count()
	fmt.Println("Total user count is ", count)
	allUsers := []User{}
	uc.Find(nil).All(&allUsers)
	for _, user := range allUsers {
		fmt.Println(user)
		fmt.Println("==================")
	}
	return nil, nil
}

func (manager *DbManager) SaveUser(user *User) error {
	uc := manager.session.DB(DbName).C(CollectionName)

	i, _ := uc.Find(bson.M{"username":user.UserName}).Count()
	if i != 0 {
		return errors.New("user name registed!!!")
	}

	i, _ = uc.Find(bson.M{"email":user.Email}).Count()
	if i != 0 {
		return errors.New("email name registed!!!")
	}

	user.Id = bson.NewObjectId()
	err := uc.Insert(user)
	return err
}
