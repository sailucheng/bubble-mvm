package mvm

import (
	tea "github.com/charmbracelet/bubbletea"
)

var nopeResult = Result{}

type Controller interface {
	Filter(*Context) bool
	Handle(*Context) Result
}

type Result struct {
	Err    error
	Model  any
	Cmd    tea.Cmd
	Viewer Viewer
}

type nopeController struct{}

func (nopeController) Filter() bool {
	return true
}

func (nopeController) Handle(ctx *Context) Result {
	switch msg := ctx.Msg.(type) {
	case tea.KeyMsg:
		{
			switch msg.String() {
			case "ctrl+c":
				return ctx.Quit()
			}
		}
	}
	return ctx.None()
}

type ControllerBase struct {
	nopeController
}
