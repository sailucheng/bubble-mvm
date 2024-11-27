package models

import "fmt"

type LoginModel struct {
	Logged   bool
	UserName string
	Password string
}

func (m *LoginModel) ValidateUserName(s string) error {
	n := len(s)
	if n == 0 {
		return fmt.Errorf("UserName required")
	}
	return nil
}

func (m *LoginModel) ValidatePass(s string) error {
	n := len(s)
	if n == 0 {
		return fmt.Errorf("password required")
	}
	if n < 6 {
		return fmt.Errorf("password len at least %d", 6)
	}
	return nil
}
