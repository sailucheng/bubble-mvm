package main

const modelTempl = `package models
type {{ .Name }}Model struct {
}`

const controllerTempl = `package controllers

import (
  "{{.Module}}/models"
  "{{.Module}}/views"
  "github.com/sailucheng/bubble-mvm/mvm"
)

type {{.Name}}Controller struct {}
func (controller {{.Name}}Controller) Filter(ctx *mvm.Context) bool {
	_ , ok := ctx.Model.(*models.{{.Name}}Model)
	return ok
}

func (controller {{.Name}}Controller) Handle(ctx *mvm.Context) mvm.Result {
	v := ctx.Viewer.(*views.{{.Name}}View)
	result := v.Update(ctx)
	return result
}`

const viewTempl = `package views
import (
  "fmt"
  tea "github.com/charmbracelet/bubbletea"
  "{{.Module}}/models"
  "github.com/sailucheng/bubble-mvm/mvm"
)

type {{.Name }}View struct {
	Model *models.{{ .Name }}Model
	quitting bool
}

func Init{{.Name}}View(m *models.{{.Name}}Model) *{{.Name}}View {
	v := {{.Name}}View{
		Model: m,
	}
	return &v
}

func (v *{{.Name}}View) Init() tea.Cmd {
	return nil
}

// This is for updating the view logic. You can write the view update logic here. 
// Don't forget to call it in the relevant controller's handle method.
func (v *{{.Name}}View) Update(ctx *mvm.Context) mvm.Result {
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

func (v *{{.Name}}View) Render(model any) string {
	if v.quitting {
		return "bye" + "\n"
	}
	return fmt.Sprintf("%T\n Press ctrl+c to exit.\n",v) 
}`

const mainTempl = `package main

import (
  "fmt"
  "log"
  "os"

  tea "github.com/charmbracelet/bubbletea"
  "{{.Module}}/controllers"
  "{{.Module}}/models"
  "{{.Module}}/views"
  "github.com/sailucheng/bubble-mvm/mvm"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug:")
	if err != nil {
		log.Fatalf("cannot create open log :%v", err)
	}
	defer f.Close()
	register()
	m := &models.{{.Name}}Model{}
	teaModel := mvm.CreateModel(m).WithInitView(views.Init{{.Name}}View(m))
	if _, err := tea.NewProgram(teaModel).Run(); err != nil {
		fmt.Println("Ops!", err)
		os.Exit(1)
	}
}
func register() {
	mvm.RegisterControllers(controllers.{{.Name}}Controller{})
}
`

type mode int

func (m mode) String() string {
	switch m {
	case ModeMain:
		return "main"
	case ModeModel:
		return "model"
	case ModeController:
		return "controller"
	case ModeView:
		return "View"
	default:
		return "unknown"
	}
}

const (
	ModeMain mode = iota
	ModeModel
	ModeController
	ModeView
)

var templateMap = map[mode]string{
	ModeMain:       mainTempl,
	ModeController: controllerTempl,
	ModeView:       viewTempl,
	ModeModel:      modelTempl,
}
