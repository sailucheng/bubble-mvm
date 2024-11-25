package models

import "fmt"

type HuhDemoModel struct {
	Username string
	Password string
}

func (m HuhDemoModel) String() string {
	return fmt.Sprintf("%#+v", m)
}
