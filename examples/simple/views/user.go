package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/sailucheng/bubble-mvm/examples/simple/models"
	"github.com/sailucheng/bubble-mvm/mvm"
)

var help = lipgloss.NewStyle().Foreground(lipgloss.Color("213")).MarginTop(1).Render

type UserView struct {
	Model    *models.User
	Quitting bool
}

func InitUserView(model *models.User) *UserView {
	return &UserView{
		Model: model,
	}
}

func (v *UserView) Init() tea.Cmd {
	return nil
}
func (v *UserView) Update(ctx *mvm.Context) mvm.Result {
	switch ctx.Msg.(type) {
	case tea.KeyMsg:
		switch ctx.Msg.(tea.KeyMsg).String() {
		case "ctrl+c":
			v.Quitting = true
			return ctx.Quit()
		}
	}
	return ctx.NoAction()
}
func (v *UserView) Render(model any) string {
	if v.Quitting {
		return "bye"
	}
	return v.Model.FullName() + "\n" +
		help("press ctrl+c to exit.")
}
