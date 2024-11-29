package mvm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTriggerControllerMethod(t *testing.T) {
	expected := "TestTrigger"
	controller := testController{
		callback: func(s string) {
			assert.Equal(t, expected, s)
		},
	}
	err := triggerControllerMethod(&controller, "OnSuccess", expected)
	assert.NoError(t, err)
}

type testController struct {
	callback func(string)
	nopeController
}

func (c *testController) OnSuccess(s string) {
	c.callback(s)
}
