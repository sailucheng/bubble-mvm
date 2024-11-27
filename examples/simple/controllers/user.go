package controllers

import (
	"github.com/sailucheng/bubble-mvm/examples/simple/models"
	"github.com/sailucheng/bubble-mvm/examples/simple/views"
	"github.com/sailucheng/bubble-mvm/mvm"
)

type UserController struct {
	mvm.ControllerBase
}

func (c UserController) Filter(*mvm.Context) bool {
	return true
}

func (c UserController) Handle(ctx *mvm.Context) mvm.Result {
	v := ctx.Viewer.(*views.UserView)
	result := v.Update(ctx)
	if ctx.IsAbort() {
		return result
	}
	user, ok := ctx.Model.(*models.User)
	if !ok {
		panic("no, it's impossible")
	}

	user.FirstName = "Bob"
	return ctx.NoAction()
}
