package reflection

import (
	"fmt"
)

// New creates and initializes an empty environment.
func New() *Env {
	return &Env{
		entities: map[string]Value{},
	}
}

// Env maintains entities available in the environment.
type Env struct {
	entities map[string]Value
}

// Bind a value to the given name. If the name already exists,
// the value will be re-bound.
func (env *Env) Bind(name string, v interface{}) error {
	entry := NewValue(v)
	env.entities[name] = entry
	return nil
}

func (env *Env) String() string {
	return fmt.Sprintf("Env[size=%d]", len(env.entities))
}

// Get returns the value bound for the name.
func (env *Env) Get(name string) (interface{}, error) {
	entry, found := env.entities[name]
	if !found {
		return nil, ErrNameNotFound
	}

	return entry.Value.Interface(), nil
}

// GetString attempts resolving the name into a string value.
func (env *Env) GetString(name string) (string, error) {
	val, err := env.Get(name)
	if err != nil {
		return "", err
	}

	stringVal, ok := val.(string)
	if !ok {
		return "", ErrConversionImpossible
	}

	return stringVal, nil
}

// GetInt attempts resolving the name into a integer value.
func (env *Env) GetInt(name string) (int, error) {
	val, err := env.Get(name)
	if err != nil {
		return 0, err
	}

	intVal, ok := val.(int)
	if !ok {
		return 0, ErrConversionImpossible
	}

	return intVal, nil
}

// GetFloat attempts resolving the name into a float64 value.
func (env *Env) GetFloat(name string) (float64, error) {
	val, err := env.Get(name)
	if err != nil {
		return 0, err
	}

	floatVal, ok := val.(float64)
	if !ok {
		return 0, ErrConversionImpossible
	}

	return floatVal, nil
}

// GetBool attempts resolving the name into a bool value.
func (env *Env) GetBool(name string) (bool, error) {
	val, err := env.Get(name)
	if err != nil {
		return false, err
	}

	boolVal, ok := val.(bool)
	if !ok {
		return false, ErrConversionImpossible
	}

	return boolVal, nil
}
