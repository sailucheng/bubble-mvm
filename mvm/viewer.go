package mvm

import (
	tea "github.com/charmbracelet/bubbletea"
)

var (
	NopViewer = nopViewer{}
)

type Viewer interface {
	Init() tea.Cmd
	Render(model any) string
}

type nopViewer struct{}

func (nop nopViewer) Init() tea.Cmd {
	return nil
}

func (nop nopViewer) Render(model any) string {
	return ""
}
