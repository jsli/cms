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

	//create an global user object for testing
	session, err := mgo.Dial("localhost")
	if err == nil {
		dal := NewDalMgo(session)
		SuperUser, _ = dal.GetUserByName(SuperUserName)
		if SuperUser != nil {
			fmt.Println("Get SuperUser")
		} else {
			fmt.Println("Cannot get SuperUser")
			power := Power{POWER_EDIT_ADMIN_USER: POWER_EDIT_ADMIN_USER, POWER_EDIT_NORMAL_USER: POWER_EDIT_NORMAL_USER}
			SuperUser = &User{
				Id:           bson.NewObjectId(),
				UserName:     SuperUserName,
				Role:         ROLE_SUPERUSER,
				HashPassword: GeneratePwdByte(SuperUserPwd),
				Email:        SuperUserEmail,
				PowerMap:     power,
				IsLogined:    false,
			}
			err := dal.SaveUser(SuperUser)
			if err == nil {
				fmt.Println("Create the super user!!!")
			} else {
				fmt.Println("Can't create the super user!!!")
			}
		}
		fmt.Println(SuperUser)
	} else {
		fmt.Println("User init failed!")
	}
}