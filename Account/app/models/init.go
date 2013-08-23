package models

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"regexp"
)

var USERNAME_REX, PWD_REX, NICKNAME_REX *regexp.Regexp
var SuperUser *User

func init() {
	USERNAME_REX = regexp.MustCompile(`^[a-zA-Z0-9_]{6,16}$`)
	PWD_REX = regexp.MustCompile(`^[\x01-\xfe]{8,20}$`)

	//NICKNAME_REX = regexp.MustCompile(`^[a-zA-Z\xa0-\xff_][0-9a-zA-Z\xa0-\xff_]{3,15}$`)
	NICKNAME_REX = USERNAME_REX

	session, err := mgo.Dial("localhost")
	if err == nil {
		user := User{}
		uc := session.DB(DbName).C(CollectionName)
		err = uc.Find(bson.M{"user_name": SuperUserName}).One(&user)
		if err == nil {
			SuperUser = &user
			fmt.Println(SuperUser)
		} else {
			fmt.Println("Cannot get SuperUser")
			power := Power{POWER_EDIT_ADMIN_USER: POWER_EDIT_ADMIN_USER, POWER_EDIT_NORMAL_USER: POWER_EDIT_NORMAL_USER}
			user = User{
				Id:           bson.NewObjectId(),
				UserName:     SuperUserName,
				Role:         ROLE_SUPERUSER,
				HashPassword: generatePwdByte(SuperUserPwd),
				Email:        SuperUserEmail,
				PowerMap:     power,
				IsLogined:    false,
			}
			err := uc.Insert(user)
			if err == nil {
				fmt.Println("Create the super user!!!")
				SuperUser = &user
				fmt.Println(SuperUser)
			} else {
				fmt.Println("can't create the super user!!!")
			}
		}
	} else {
		fmt.Println("User init failed!")
	}
}