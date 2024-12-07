package mvm

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"reflect"
)

var nopeResult = Result{}

type Controller interface {
	Filter(*Context) bool
	Handle(*Context) Result
}

type Result struct {
	Err    error
	Model  any
	Cmd    tea.Cmd
	Viewer Viewer
}

func (r *Result) Composite(other Result) {
	if other.Err != nil {
		if r.Err == nil {
			r.Err = other.Err
		}
		r.Err = fmt.Errorf("composite errors: %w and %w", r.Err, other.Err)
	}
	r.Model = other.Model
	r.Cmd = other.Cmd
	r.Viewer = other.Viewer
}

type nopeController struct{}

func (nopeController) Filter(*Context) bool {
	return true
}

func (nopeController) Handle(ctx *Context) Result {
	return ctx.None()
}

type ControllerBase struct {
	nopeController
}

func getControllerKey(c Controller) (string, error) {
	v := reflect.ValueOf(c)
	if !v.IsValid() {
		return "", fmt.Errorf("invalid controller: nil or uninitialized")
	}
	t := v.Type()
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	pkgName := t.PkgPath()
	typeName := t.Name()
	return fmt.Sprintf("%s.%s", pkgName, typeName), nil
}
