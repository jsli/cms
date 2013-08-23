package controllers

import (
	"fmt"
	"github.com/jgraham909/revmgo"
	"github.com/jsli/cms/Account/app/models"
	"github.com/robfig/revel"

	"labix.org/v2/mgo/bson"
)

type Account struct {
	*revel.Controller
	revmgo.MongoController
}

func init() {
	revmgo.ControllerInit()
	revel.OnAppStart(revmgo.AppInit)
}

func (c Account) Index() revel.Result {
	models.SuperUser.ListUsers(c.MongoSession, 1, 10)
	models.SuperUser.GetUserById(c.MongoSession, bson.NewObjectId())
	models.SuperUser.GetUserById(c.MongoSession, bson.ObjectIdHex("5217211f8223a732f1000002"))
	models.SuperUser.GetUserByName(c.MongoSession, "manson")
	user, _ := models.SuperUser.GetUserByEmail(c.MongoSession, "lijinsong@163.com")
	user.GetUserByName(c.MongoSession, "manson")
	return c.Render()
}

func (c Account) GetLogin() revel.Result {
	return c.Render()
}

func (c Account) PostLogin(loginUser *models.LoginUser) revel.Result {
	//workflow is the same as PostRegister
	//step 0: check user is exist or not
	loginUser.Validate(c.Validation, c.MongoSession)

	//step 1: validation
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Account.GetLogin)
	}

	//step 3: save cookie, flash or session
	c.Session["user"] = loginUser.UserName
	c.Flash.Success("Welcome, login " + loginUser.UserName)
	fmt.Println("Welcome, login ", loginUser.UserName)

	//step 4: rediret
	return c.Redirect(Account.Index)
}

func (c Account) Logout() revel.Result {
	//logout status
	user := models.GetUserByName(c.MongoSession, c.Session["user"])
	if user != nil {
		user.IsLogined = false
		models.UpdateUser(c.MongoSession, *user)
	}
	for k := range c.Session {
		delete(c.Session, k)
	}
	c.Flash.Success("Welcome, logout ")
	return c.Redirect(App.Index)
}

func (c Account) GetRegister() revel.Result {
	return c.Render()
}

/*
 * regUser is a struct's name in the template
 * see {{with $field := field "regUser.Field" .}} in template
 */
func (c Account) PostRegister(regUser *models.RegUser) revel.Result {
	//step 0: check user is exist or not
	regUser.Validate(c.Validation, c.MongoSession)

	//step 1: validation
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Account.GetRegister)
	}

	//step 2: save user
	err := regUser.SaveUser(c.MongoSession)
	if err != nil {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Account.GetRegister)
	}

	//step 3: save cookie, flash or session
	c.Session["user"] = regUser.UserName
	c.Flash.Success("Welcome, register " + regUser.UserName)
	fmt.Println("welcome register", regUser.UserName)

	//step 4: rediret
	return c.Redirect(Account.Index)
}
