package mvm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTriggerControllerMethodWithoutResult(t *testing.T) {
	expected := "TestTrigger"
	controller := testController{
		callback: func(s string) {
			assert.Equal(t, expected, s)
		},
	}
	ret, err := triggerControllerMethod(&controller, "OnCallback", expected)
	assert.NoError(t, err)
	assert.Empty(t, ret)
}
func TestTriggerControllerWithResult(t *testing.T) {
	expectedMessage := "TestTrigger"
	controller := testController{
		callback2: func(s string) Result {
			return Result{
				Model: s,
			}
		},
	}

	result, err := triggerControllerMethod(&controller, "OnCallback2", expectedMessage)
	assert.NoError(t, err)
	assert.Equal(t, expectedMessage, result.Model)
}

type testController struct {
	callback  func(string)
	callback2 func(string) Result
	nopeController
}

func (c *testController) OnCallback2(s string) Result {
	return c.callback2(s)
}
func (c *testController) OnCallback(s string) {
	c.callback(s)
}
