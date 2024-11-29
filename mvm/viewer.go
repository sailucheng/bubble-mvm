package mvm

import (
	tea "github.com/charmbracelet/bubbletea"
)

var NopeViewer = nopViewer{}

type Viewer interface {
	Init() tea.Cmd
	Update(ctx *Context) Result
	Render(model any) string
}

type nopViewer struct{}

func (nop nopViewer) Init() tea.Cmd {
	return nil
}

func (nop nopViewer) Render(model any) string {
	return ""
}

func (nop nopViewer) Update(ctx *Context) Result {
	return Result{}
}
