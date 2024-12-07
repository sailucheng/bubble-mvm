package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sailucheng/bubble-mvm/examples/redirect-views/controllers"
	"github.com/sailucheng/bubble-mvm/examples/redirect-views/models"
	"github.com/sailucheng/bubble-mvm/examples/redirect-views/views"
	"github.com/sailucheng/bubble-mvm/mvm"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug:")
	if err != nil {
		log.Fatalf("cannot create open log :%v", err)
	}
	defer f.Close()
	register()
	m := &models.LoginModel{}
	teaModel := mvm.CreateModel(m).WithInitView(views.InitLoginView(m))
	if _, err := tea.NewProgram(teaModel).Run(); err != nil {
		fmt.Println("Ops!", err)
		os.Exit(1)
	}
}

func register() {
	mvm.RegisterControllers(
		&controllers.LoginController{},
		&controllers.ContentController{},
	)
}
