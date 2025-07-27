package util

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type (
	BeanScope      int
	BeanDefinition struct {
		Name        string
		Type        reflect.Type
		Scope       BeanScope
		Constructor interface{}

		instance interface{}
		initOnce sync.Once
		mu       sync.Mutex
	}

	BeanContainer struct {
		beans map[string]*BeanDefinition
		mu    sync.RWMutex
	}

	InitializingBean interface {
		AfterPropertiesSet()
	}
)

const (
	BeanScopeSingleton BeanScope = iota
	BeanScopePrototype
)

func NewBeanContainer() *BeanContainer {
	return &BeanContainer{beans: make(map[string]*BeanDefinition)}
}

func (c *BeanContainer) RegisterBean(name string, constructor interface{}, scope BeanScope) error {
	constructorType := reflect.TypeOf(constructor)
	if constructorType.Kind() != reflect.Func {
		return errors.New("constructor must be a function")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.beans[name] = &BeanDefinition{
		Name:        name,
		Type:        constructorType.Out(0),
		Scope:       scope,
		Constructor: constructor,
	}
	return nil
}

func (c *BeanContainer) GetBean(name string) (interface{}, error) {
	c.mu.RLock()
	definition, exists := c.beans[name]
	c.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("bean %s not registered", name)
	}

	if definition.Scope == BeanScopeSingleton {
		definition.mu.Lock()
		defer definition.mu.Unlock()
		var err error = nil
		definition.initOnce.Do(func() {
			definition.instance, err = c.createInstance(definition)
		})
		return definition.instance, err
	} else {
		return c.createInstance(definition)
	}
}

func (c *BeanContainer) createInstance(definition *BeanDefinition) (interface{}, error) {
	constructorValue := reflect.ValueOf(definition.Constructor)
	args := c.prepareArgs(constructorValue)
	instanceValue := constructorValue.Call(args)[0]
	instance := instanceValue.Interface()

	err := c.injectDependencies(instanceValue.Elem())
	if err != nil {
		return nil, err
	}

	if init, ok := instance.(InitializingBean); ok {
		init.AfterPropertiesSet()
	}
	return instance, nil
}

func (c *BeanContainer) prepareArgs(constructor reflect.Value) []reflect.Value {
	constructorType := constructor.Type()
	numArgs := constructorType.NumIn()
	args := make([]reflect.Value, numArgs)

	for i := 0; i < numArgs; i++ {
		argType := constructorType.In(i)
		args[i] = reflect.New(argType).Elem()
	}

	return args
}

func (c *BeanContainer) injectDependencies(instanceValue reflect.Value) error {
	instanceType := instanceValue.Type()
	for i := 0; i < instanceValue.NumField(); i++ {
		field := instanceType.Field(i)
		if tag, ok := field.Tag.Lookup("inject"); ok {
			if tag == "" {
				tag = field.Name
			}
			dep, err := c.GetBean(tag)
			if err != nil {
				return err
			}
			fieldValue := instanceValue.Field(i)
			fieldValue.Set(reflect.ValueOf(dep))
		}
	}
	return nil
}
