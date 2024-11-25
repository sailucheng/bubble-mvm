package mvm

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestModeInit(t *testing.T) {
	m := CreateModel(nil)
	assert.Nil(t, m.Init())
	assert.IsType(t, nopViewer{}, m.viwerer)
	m = m.WithInitView(&testViewer{})
	var tw testViewer
	assert.IsType(t, &tw, m.viwerer)
	assert.Equal(t, tw.Init()(), m.Init()())
}

func TestModelUpdate(t *testing.T) {
	expected := "for test viewer out"
	var hint bool
	controller := modelTestController{
		t:               t,
		expectedViewout: expected,
		hints: func() {
			hint = true
		},
	}
	pipe := CreatePipe(withPipeControllers(controller))
	m := CreateModel(nil)
	m.pipe = pipe

	model, _ := m.Update(nil)

	assert.True(t, hint)
	assert.Equal(t, expected, model.View())
}

func TestResultWithoutViewWillUseInitialView(t *testing.T) {
	v := &testViewer{
		t:        t,
		expected: "foo zoo bar",
	}

	pip := CreatePipe()
	model := CreateModel(nil).WithInitView(v).WithPipe(pip)
	model.Update(nil)

	assert.Same(t, v, model.viwerer)
}

type modelTestController struct {
	t               *testing.T
	expectedViewout string
	hints           func()
}

func (controller modelTestController) Filter(*Context) bool {
	return true
}

func (controller modelTestController) Handle(c *Context) Result {
	controller.hints()

	v := testViewer{
		t:        controller.t,
		expected: controller.expectedViewout,
	}

	return c.View(&v)
}

type testViewer struct {
	t        *testing.T
	expected string
}

func (tw *testViewer) Render(model any) string {
	return tw.expected
}

func (*testViewer) Init() tea.Cmd {
	return func() tea.Msg {
		return "test"
	}
}
