package views

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/sailucheng/bubble-mvm/examples/huh-demo/models"
)

var (
	normal = lipgloss.NewStyle().Render
	help   = lipgloss.NewStyle().Foreground(lipgloss.Color("213")).MarginTop(1).Render
)

type HuhDemoView struct {
	Model    *models.HuhDemoModel
	Quitting bool
	Form     *huh.Form
}

func InitViewModel(model *models.HuhDemoModel) *HuhDemoView {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("username").Value(&model.Username).Key("username"),
			huh.NewInput().Title("password").Value(&model.Password).Key("password").
				EchoMode(huh.EchoModePassword),
		),
		huh.NewGroup(
			huh.NewConfirm().Title("Are you sure?").Affirmative("ok").Negative("cancel").Validate(func(b bool) error {
				if b == false {
					return fmt.Errorf("please select ok button")
				}
				return nil
			})),
	).WithShowHelp(false)

	return &HuhDemoView{
		Model: model,
		Form:  f,
	}
}

func (v *HuhDemoView) Init() tea.Cmd {
	return v.Form.Init()
}

func (v *HuhDemoView) Render(model any) string {
	if v.Quitting {
		return v.Model.String() + "\n"
	}

	return normal(v.Form.View() + "\n" + help("press ctrl+c to exit"+"\n"))
}
