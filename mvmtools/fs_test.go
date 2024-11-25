package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindModule(t *testing.T) {
	pwd, _ := os.Getwd()

	maxDepth := 10
	root, module, err := findGoMod(pwd, maxDepth)
	assert.NoError(t, err)
	assert.NotEmpty(t, module)
	assert.True(t, strings.Index(pwd, root) >= 0)
}

func TestFindModuleInRoot(t *testing.T) {
	dirs := []string{"C:/", "/"}

	maxDepth := 10
	for _, d := range dirs {
		root, module, err := findGoMod(d, maxDepth)
		assert.Errorf(t, err, "cannot find go.mod file, you must initial a golang project first")
		assert.Empty(t, root)
		assert.Empty(t, module)
	}
}
