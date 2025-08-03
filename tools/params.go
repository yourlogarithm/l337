package tools

import "fmt"

type Params map[string]any

type Parameter struct {
	Name    string
	Element Element
}

func GetParameter[T any](params Params, name string) (T, error) {
	var zero T
	value, exists := params[name]
	if !exists {
		return zero, fmt.Errorf("parameter %s not found", name)
	}
	if v, ok := value.(T); ok {
		return v, nil
	}
	return zero, fmt.Errorf("parameter %s expected %T but got %s", name, zero, value)
}
