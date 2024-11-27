package controllers

import (
	"github.com/sailucheng/bubble-mvm/examples/redirect-views/models"
	"github.com/sailucheng/bubble-mvm/examples/redirect-views/views"
	"github.com/sailucheng/bubble-mvm/mvm"
)

type ContentController struct {
	mvm.ControllerBase
}

func (controller ContentController) Filter(ctx *mvm.Context) bool {
	m, ok := ctx.Model.(*models.LoginModel)
	return ok && m.Logged
}

func (controller ContentController) Handle(ctx *mvm.Context) mvm.Result {
	v := ctx.Viewer.(*views.ContentView)
	return v.Update(ctx)
}
