package controllers

import (
	"github.com/charmbracelet/huh"
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
	result := controller.ControllerBase.Handle(ctx)
	if result.Cmd != nil {
		return result
	}

	v, _ := ctx.Viewer.(*views.HuhDemoView)
	formModel, cmd := v.Form.Update(ctx.Msg)
	if formModel, ok := formModel.(*huh.Form); ok {
		v.Form = formModel
		if formModel.State == huh.StateCompleted {
			v.Quitting = true
			return ctx.Quit()
		}
	}

	return mvm.Result{
		Model: ctx.Model,
		Cmd:   cmd,
	}
}
