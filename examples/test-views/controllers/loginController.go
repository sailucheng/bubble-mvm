package controllers

import (
  "github.com/sailucheng/bubble-mvm/examples/test-views/models"
  "github.com/sailucheng/bubble-mvm/examples/test-views/views"
  "github.com/sailucheng/bubble-mvm/mvm"
)

type LoginController struct {}
func (controller LoginController) Filter(ctx *mvm.Context) bool {
	_ , ok := ctx.Model.(*models.LoginModel)
	return ok
}

func (controller LoginController) Handle(ctx *mvm.Context) mvm.Result {
	v := ctx.Viewer.(*views.LoginView)
	result := v.Update(ctx)
	return result
}