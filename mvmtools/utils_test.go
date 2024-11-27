package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateModule(t *testing.T) {
	pwd, _ := os.Getwd()
	root, module, _ := findGoMod(pwd, 10)
	rawName := "jojo"
	module, err := calculateModule(module, rawName, root, pwd)
	assert.NoError(t, err)
}

func TestCalculateModuleCompatibleWindowsSeparator(t *testing.T) {
	root := `D:\code\go\src\bubble-mvm`
	dest := `D:\code\go\src\bubble-mvm\examples\multi-views`
	module := "github.com/sailucheng/bubble-mvm"
	rawName := "jojo"

	calcModule, err := calculateModule(module, rawName, root, dest)
	expected := "github.com/sailucheng/bubble-mvm/examples/multi-views"
	assert.NoError(t, err)
	assert.Equal(t, expected, calcModule)
}
