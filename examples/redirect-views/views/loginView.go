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
	wrong = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).MarginTop(1).Render
)

type LoginView struct {
	Model *models.LoginModel
	Form  *huh.Form
	err   error
}

func InitLoginView(m *models.LoginModel) *LoginView {
	h := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("username").Value(&m.UserName).Validate(m.ValidateUserName),
			huh.NewInput().Title("password").Value(&m.Password).Validate(m.ValidatePass).EchoMode(huh.EchoModePassword),
			huh.NewConfirm().Title("login now?").Validate(func(b bool) error {
				if !b {
					return fmt.Errorf("you must select login")
				}
				return nil
			}),
		),
	).WithShowHelp(false)
	v := LoginView{
		Model: m,
		Form:  h,
	}
	return &v
}

func (v *LoginView) Init() tea.Cmd {
	return v.Form.Init()
}
func (v *LoginView) Update(ctx *mvm.Context) mvm.Result {
	switch msg := ctx.Msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return ctx.Quit()
		}
	}
	formModel, cmd := v.Form.Update(ctx.Msg)

	if formModel, ok := formModel.(*huh.Form); ok {
		v.Form = formModel
	}
	return ctx.Cmd(cmd)
}
func (v *LoginView) Render(model any) string {
	if v.err != nil {
		return v.Form.View() + wrong(v.err.Error()) + "\n"
	}

	return v.Form.View() + "\n"
}

func (v *LoginView) SetError(err error) {
	v.err = err
}
