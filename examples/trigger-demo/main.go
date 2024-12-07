package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sailucheng/bubble-mvm/mvm"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug:")
	if err != nil {
		log.Fatalf("cannot create open log :%v", err)
	}
	defer f.Close()
	mvm.RegisterControllers(&TriggerController{})
	v := TriggerView{
		name: "trigger view",
	}
	model := mvm.CreateModel(nil).WithInitView(&v)
	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Println("Ops!", err)
		os.Exit(1)
	}
}
