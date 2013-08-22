package controllers

import (
	"github.com/robfig/revel"
	"github.com/jsli/revel-in-action/Account/app/models"
)

type Account struct {
	*revel.Controller
}

func (c Account) Index() revel.Result {
	return c.Render()
}

func (c Account) GetLogin() revel.Result {
	models.GetAllUsers()
	return c.Render()
}

func (c Account) PostLogin(loginUser *models.LoginUser) revel.Result {
	//workflow is the same as PostRegister
	//step 0: check user is exist or not
	loginUser.Validate(c.Validation)

	//step 1: validation
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Account.GetLogin)
	}

	//step 3: save cookie, flash or session
	c.Session["user"] = loginUser.UserName
	c.Flash.Success("Welcome, login " + loginUser.UserName)

	//step 4: rediret
	return c.Redirect(Account.Index)
}

func (c Account) Logout() revel.Result {
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
	regUser.Validate(c.Validation)

	//step 1: validation
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Account.GetRegister)
	}

	//step 2: save user
	err := regUser.SaveUser()
	if err != nil {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Account.GetRegister)
	}

	//step 3: save cookie, flash or session
	c.Session["user"] = regUser.UserName
	c.Flash.Success("Welcome, register " + regUser.UserName)

	//step 4: rediret
	return c.Redirect(Account.Index)
}
