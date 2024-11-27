package environment

import (
	"fmt"
)

type Environment struct {
	values map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{values: make(map[string]interface{})}
}

func (e *Environment) Get(name string) (interface{}, error) {
	if val, ok := e.values[name]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("Undefined variable '%s'.", name)
}

func (e *Environment) Assign(name string, value interface{}) error {
	if _, ok := e.values[name]; ok {
		e.Define(name, value)
		return nil
	}

	return fmt.Errorf("Undefined variable '%s'.", name)
}

func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}
