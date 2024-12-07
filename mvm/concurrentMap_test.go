package mvm

import (
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOrAdd(t *testing.T) {
	cache := CreateMethodInvokerCache()

	validNames := []string{
		"Method1",
		"Method2",
		"CorrectReturnType",
	}

	for _, name := range validNames {
		v, ok, err := cache.GetOrAdd(concurrentMapTestController{}, name)
		assert.NoError(t, err)
		assert.NotEmpty(t, v)
		assert.True(t, ok)
	}
	for _, name := range validNames {
		v, ok, err := cache.GetOrAdd(concurrentMapTestController{}, name)
		assert.NoError(t, err)
		assert.NotEmpty(t, v)
		assert.False(t, ok)
	}
	invalidNames := []string{"NonExistsMethod", "WrongReturnType"}
	for _, name := range invalidNames {
		v, ok, err := cache.GetOrAdd(concurrentMapTestController{}, name)
		assert.Error(t, err)
		assert.Empty(t, v)
		assert.False(t, ok)
	}
}
func TestParallelAdd(t *testing.T) {
	cache := CreateMethodInvokerCache()
	var wg sync.WaitGroup
	validNames := []string{
		"Method1",
		"Method2",
		"CorrectReturnType",
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			// 这里不管新建多少个controller，应该没问题。因为key是幂等的.
			controller := concurrentMapTestController{}
			for _, name := range validNames {
				_, err := cache.Add(controller, name)
				assert.NoError(t, err)
			}

			added, err := cache.Add(controller, "WrongReturnType")
			assert.False(t, added)
			assert.Error(t, err)
		}(i)
	}
	wg.Wait()
	controller := concurrentMapTestController{}
	key, _ := getControllerKey(controller)
	assert.Equal(t, len(cache.cache[key]), 3)
}

func TestGet(t *testing.T) {
	cache := CreateMethodInvokerCache()
	controller := concurrentMapTestController{}
	methodNames := []string{"Method1", "Method2", "CorrectReturnType"}

	for _, name := range methodNames {
		_, _ = cache.Add(controller, name)
	}

	_, _ = cache.Add(controller, "WrongReturnType")
	for _, name := range methodNames {
		invoker, ok, err := cache.Get(controller, name)
		assert.IsType(t, reflect.Value{}, invoker)
		assert.NotEmpty(t, invoker)
		assert.Equal(t, true, ok)
		assert.NoError(t, err)
	}
	for _, name := range []string{"WrongReturnType", "NonExistsMethod"} {
		invoker, ok, err := cache.Get(controller, name)
		assert.IsType(t, reflect.Value{}, invoker)
		assert.Equal(t, false, ok)
		assert.NoError(t, err)
	}

}

type concurrentMapTestController struct {
	ControllerBase
}

func (c concurrentMapTestController) Method1()           {}
func (c concurrentMapTestController) Method2(arg string) {}
func (c concurrentMapTestController) WrongReturnType() string {
	return ""
}

func (c concurrentMapTestController) CorrectReturnType() Result {
	return nopeResult
}
