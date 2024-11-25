package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateModule(t *testing.T) {
	pwd, _ := os.Getwd()
	root, module, _ := findGoMod(pwd, 10)
	projectName := "jojo"
	module, err := calculateModule(module, projectName, root, pwd)
	assert.NoError(t, err)
}
