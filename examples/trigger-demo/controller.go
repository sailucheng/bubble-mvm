package main

import (
	"fmt"

	"github.com/sailucheng/bubble-mvm/mvm"
)

type TriggerController struct {
}

func (c *TriggerController) Filter(ctx *mvm.Context) bool {
	return true
}
func (c *TriggerController) Handle(ctx *mvm.Context) mvm.Result {
	return ctx.Propagate()
}

func (c *TriggerController) Quit(message string, v *TriggerView) {
	v.quitting = true
	v.message = fmt.Sprintf("message :%s\nview name:%s\n", message, v.name)
}
