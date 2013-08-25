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

func (c Account) Index() revel.Result {
	//	models.SuperUser.ListUsers(c.MongoSession, 1, 10)
	//	models.SuperUser.SaveUser(c.MongoSession)
	//	user := models.User{
	//		UserName:     "testuser",
	//		Role:         models.ROLE_NORMAL,
	//		HashPassword: models.GeneratePwdByte("12345678"),
	//		Email:        "testuser@mail.com",
	//		IsLogined:    false,
	//	}
	//	user.SaveUser(c.MongoSession)

	//	models.CheckPermission(models.SuperUser, models.POWER_EDIT_ADMIN_USER)
	//	models.CheckPermission(models.SuperUser, models.POWER_EDIT_NORMAL_USER)
	//	models.CheckPermission(models.SuperUser, "fack_permission")
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
	err := loginUser.LoadSelf(c.MongoSession)
	if err == nil {
		loginUser.IsLogined = true
		loginUser.UpdateUser(c.MongoSession)
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
	user, err := models.GetUserByName(c.MongoSession, c.Session["user"])
	if user != nil && err == nil {
		user.IsLogined = false
		user.UpdateUser(c.MongoSession)
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

func (c Account) ListUsers() revel.Result {
	var page, count, role int
	c.Params.Bind(&role, "role")
	c.Params.Bind(&count, "count")
	c.Params.Bind(&page, "page")
	if role <= 0 || role > 3{
		role = models.ROLE_NORMAL // normal user default
	}
	if count <= 0{
		count = 10 // 10 per-page default
	}
	if page <= 0 {
		page = 1 // first page default
	}
	
	user, err := models.GetUserByName(c.MongoSession, c.Session["user"])
	if user == nil || err != nil {
		c.Validation.Error("Please login first")
		c.Validation.Keep();
		c.FlashParams();
		return c.Redirect(Account.GetLogin)
	}
	
	user.ListUsers(c.MongoSession, page, count, role)
	return c.Render()
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