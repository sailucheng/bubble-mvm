package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sailucheng/bubble-mvm/mvm"
)

type TriggerView struct {
	message  string
	name     string
	quitting bool
}

func (v *TriggerView) Init() tea.Cmd {
	return nil
}

func (v *TriggerView) Update(ctx *mvm.Context) mvm.Result {
	switch msg := ctx.Msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			ctx.Trigger("Quit", "this is quit message", v)
			return ctx.Quit()
		}
	}
	return ctx.Propagate()
}

func (v *TriggerView) Render(model any) string {
	if !v.quitting {
		return "Press Ctrl+C to quit"
	}
	return v.message + "\n"
}
