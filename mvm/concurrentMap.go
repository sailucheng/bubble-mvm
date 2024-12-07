package mvm

import (
	"fmt"
	"reflect"
)

var resultType = reflect.TypeOf(nopeResult)

type MethodInvokerCache struct {
	cache map[string]map[string]reflect.Value
	mu    ReentrantMutex
}

func CreateMethodInvokerCache() *MethodInvokerCache {
	return &MethodInvokerCache{
		cache: make(map[string]map[string]reflect.Value),
	}
}
func (c *MethodInvokerCache) GetOrAdd(controller Controller, method string) (reflect.Value, bool, error) {
	var (
		value reflect.Value
		err   error
		ok    bool
	)

	c.mu.Lock()
	defer c.mu.Unlock()

	value, ok, err = c.Get(controller, method)
	if err != nil {
		return value, false, err
	}

	if ok {
		return value, false, nil
	}

	ok, err = c.Add(controller, method)
	if err != nil {
		return value, false, err
	}

	key, _ := getControllerKey(controller)
	return c.cache[key][method], ok, nil
}
func (c *MethodInvokerCache) Add(controller Controller, method string) (bool, error) {
	key, err := getControllerKey(controller)
	if err != nil {
		return false, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	methods, ok := c.cache[key]
	if !ok {
		methods = make(map[string]reflect.Value)
		c.cache[key] = methods
	}
	if _, ok = methods[method]; ok {
		return false, nil
	}

	controllerValue := reflect.ValueOf(controller)
	methodValue := controllerValue.MethodByName(method)
	if !methodValue.IsValid() {
		return false, fmt.Errorf("method %s is not found in the controller", method)
	}
	methodType := methodValue.Type()
	// Check the method parameter type invalid.
	if methodType.NumOut() > 0 {
		firstArgType := methodType.Out(0)
		if firstArgType != resultType {
			return false, fmt.Errorf("the first parameter of the 'trigger' must be of type 'mvm.Result', but found %s", firstArgType.Name())
		}
	}

	methods[method] = methodValue
	return true, nil
}

func (c *MethodInvokerCache) Get(controller Controller, method string) (reflect.Value, bool, error) {
	key, err := getControllerKey(controller)
	var v reflect.Value
	if err != nil {
		return v, false, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	methods, ok := c.cache[key]
	if !ok {
		return v, false, nil
	}

	v, ok = methods[method]
	return v, ok, nil
}
