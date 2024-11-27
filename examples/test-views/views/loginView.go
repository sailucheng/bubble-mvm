package views
import (
  "fmt"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/sailucheng/bubble-mvm/examples/test-views/models"
  "github.com/sailucheng/bubble-mvm/mvm"
)

type LoginView struct {
	Model *models.LoginModel
	quitting bool
}

func InitLoginView(m *models.LoginModel) *LoginView {
	v := LoginView{
		Model: m,
	}
	return &v
}

func (v *LoginView) Init() tea.Cmd {
	return nil
}

// This is for updating the view logic. You can write the view update logic here. 
// Don't forget to call it in the relevant controller's handle method.
func (v *LoginView) Update(ctx *mvm.Context) mvm.Result {
	switch ctx.Msg.(type) {
	case tea.KeyMsg:
		switch ctx.Msg.(tea.KeyMsg).String() {
		case "ctrl+c":
			v.quitting = true
			return ctx.Quit()
		}
	}
	return ctx.NoAction()
}

func (v *LoginView) Render(model any) string {
	if v.quitting {
		return "bye" + "\n"
	}
	return fmt.Sprintf("%T\n Press ctrl+c to exit.\n",v) 
}