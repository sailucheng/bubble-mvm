package controllers

import (
	"github.com/sailucheng/bubble-mvm/mvm"
	"log"
)

type HuhDemoController struct {
	mvm.ControllerBase
}

func (controller *HuhDemoController) Filter(ctx *mvm.Context) bool {
	return true
}

func (controller *HuhDemoController) Handle(ctx *mvm.Context) mvm.Result {
	return ctx.Cmd(ctx.Result.Cmd)
}

func (controller *HuhDemoController) OnExit() {
	log.Println("on exit method called")
}
