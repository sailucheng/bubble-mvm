package controllers

import (
	"errors"

	"github.com/charmbracelet/huh"
	"github.com/sailucheng/bubble-mvm/examples/redirect-views/models"
	"github.com/sailucheng/bubble-mvm/examples/redirect-views/views"
	"github.com/sailucheng/bubble-mvm/mvm"
)

type LoginController struct{}

func (controller LoginController) Filter(ctx *mvm.Context) bool {
	loginModel, ok := ctx.Model.(*models.LoginModel)
	return ok && !loginModel.Logged
}

func (controller LoginController) Handle(ctx *mvm.Context) mvm.Result {
	v := ctx.Viewer.(*views.LoginView)
	result := v.Update(ctx)
	if ctx.IsAbort() {
		return result
	}
	if v.Form.State != huh.StateCompleted {
		return result
	}
	model := ctx.Model.(*models.LoginModel)

	if model.Password != "123456" || model.UserName != "jojo" {
		v = views.InitLoginView(model)
		v.SetError(errors.New("login failed, check accounts and retry"))
		//recreate view to reset form
		return mvm.Result{
			Viewer: v,
			Cmd:    v.Init(),
			Model:  model,
		}
	}
	model.Logged = true
	return ctx.Redirect(model, views.InitContentView(model))
}
