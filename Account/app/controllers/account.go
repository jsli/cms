package controllers

import (
	"fmt"
	"github.com/jgraham909/revmgo"
	"github.com/jsli/cms/Account/app/models"
	"github.com/robfig/revel"
)

type Account struct {
	*revel.Controller
	revmgo.MongoController
}

func init() {
	revmgo.ControllerInit()
	revel.OnAppStart(revmgo.AppInit)
}

func (c Account)testDal() {
//		ops := models.NewDalMgo(c.MongoSession)
	//		ops.ListUsers(models.SuperUser,  1, 10, models.ROLE_NORMAL)

	//		user, _ := ops.GetUserByName("wangying")
	//		fmt.Println(user)
	//		user, _ := ops.GetUserByEmail("admintest@gmail.com")
	//		fmt.Println(user)
	//		user, _ := ops.GetUserById("52184b318223a72374000002")
	//		fmt.Println(user)

	//	user := models.User{
	//		UserName:     "dalUser",
	//		Role:         models.ROLE_NORMAL,
	//		HashPassword: models.GeneratePwdByte("12345678"),
	//		Email:        "testuser@mail.com111",
	//		IsLogined:    false,
	//	}
	//	ops.SaveUser(user)

	//ops.DeleteUserById(models.SuperUser, "5218bf858223a70fe0000002")

//	user, _ := ops.GetUserById("521b62448223a71e3a000002")
//	user.UserName = "adminadmin####"
//	user.Email = "admintest@gmail.com11!!!!!!!!"
//	ops.UpdateUserById(models.SuperUser, user)
}

func (c Account) Index() revel.Result {
	c.testDal()
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

	//update login status
	dal := models.NewDalMgo(c.MongoSession)
	user, err := dal.GetUserByName(loginUser.UserName)
	if user != nil && err == nil {
		user.IsLogined = true
		dal.UpdateUserById(user, user)
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
	dal := models.NewDalMgo(c.MongoSession)
	user, err := dal.GetUserByName(c.Session["user"])
	if user != nil && err == nil {
		user.IsLogined = false
		dal.UpdateUserById(user, user)
	}
	for k := range c.Session {
		delete(c.Session, k)
	}
	c.Flash.Success(fmt.Sprintf("Welcome, logout %s", user.UserName))
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
	regUser.IsLogined = true
	regUser.Role = models.ROLE_NORMAL
	err := regUser.SaveUser(c.MongoSession)
	if err != nil {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Account.GetRegister)
	}

	//step 3: save cookie, flash or session
	c.Session["user"] = regUser.UserName
	c.Flash.Success("Welcome, register normal " + regUser.UserName)
	fmt.Println("welcome register normal", regUser.UserName)

	//step 4: rediret
	return c.Redirect(revel.MainRouter.Reverse("account.index", make(map[string]string)).Url)
}

func (c Account) GetCreate() revel.Result {
	return c.Render()
}

func (c Account) PostCreate(regUser *models.RegUser) revel.Result {
	//step 0: check user is exist or not
	regUser.Validate(c.Validation, c.MongoSession)

	//step 1: validation
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Account.GetRegister)
	}

	//step 2: save user
	regUser.IsLogined = false
	regUser.Role = models.ROLE_ADMIN
	err := regUser.SaveUser(c.MongoSession)
	if err != nil {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Account.GetRegister)
	}

	//step 3: save cookie, flash or session
	c.Session["user"] = regUser.UserName
	c.Flash.Success("Welcome, register admin " + regUser.UserName)
	fmt.Println("welcome register admin ", regUser.UserName)

	//step 4: rediret
	return c.Redirect(revel.MainRouter.Reverse("account.index", make(map[string]string)).Url)
}
