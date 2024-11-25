package mvm

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	pipe    *Pipe
	Model   any
	viwerer Viewer
}

func CreateModel(model any) Model {
	m := Model{
		pipe:    DefaultPipe,
		Model:   model,
		viwerer: nopViewer{},
	}
	return m
}
func (m Model) WithPipe(pip *Pipe) Model {
	m.pipe = pip
	return m
}
func (m Model) WithInitView(vi Viewer) Model {
	m.viwerer = vi
	return m
}

func (m Model) Init() tea.Cmd {
	return m.viwerer.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	ctx := m.buildContext(msg)

	defer func() {
		ContextPool.Put(ctx)
	}()

	m.pipe.Execute(ctx)

	if ctx.Result != nil {
		return ctx.TeaModel, ctx.Result.Cmd
	}
	return ctx.TeaModel, nil
}

func (m Model) View() string {
	return m.viwerer.Render(m.Model)
}

func (m Model) buildContext(msg tea.Msg) *Context {
	c := ContextPool.Get().(*Context)
	c.TeaModel = m
	c.Model = m.Model
	c.Msg = msg
	c.Viewer = m.viwerer
	return c
}
