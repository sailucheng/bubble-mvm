package controllers

import (
	"github.com/sailucheng/bubble-mvm/examples/huh-demo/views"
	"github.com/sailucheng/bubble-mvm/mvm"
)

type HuhDemoController struct {
	mvm.ControllerBase
}

func (controller HuhDemoController) Filter(ctx *mvm.Context) bool {
	return true
}

func (controller HuhDemoController) Handle(ctx *mvm.Context) mvm.Result {
	v := ctx.Viewer.(*views.HuhDemoView)
	result := v.Update(ctx)
	if ctx.IsAbort() {
		return result
	}
	return result
}
