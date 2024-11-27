package views

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/sailucheng/bubble-mvm/examples/redirect-views/models"
	"github.com/sailucheng/bubble-mvm/mvm"
)

var (
	help = lipgloss.NewStyle().Foreground(huh.ThemeBase16().Help.Ellipsis.GetForeground()).Render
)

type ContentView struct {
	Model    *models.LoginModel
	quitting bool
}

func InitContentView(m *models.LoginModel) *ContentView {
	v := ContentView{
		Model: m,
	}
	return &v
}

func (v *ContentView) Init() tea.Cmd {
	return nil
}

func (v *ContentView) Update(ctx *mvm.Context) mvm.Result {
	switch ctx.Msg.(type) {
	case tea.KeyMsg:
		v.quitting = true
		return ctx.Quit()
	}
	return ctx.NoAction()
}

func (v *ContentView) Render(model any) string {
	if v.quitting {
		return "logout" + "\n"
	}
	loginModel := model.(*models.LoginModel)
	return fmt.Sprintf("hello %s\n", loginModel.UserName) + help("press any key to exit")
}
