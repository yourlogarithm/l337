package tools

import (
	"fmt"
	"reflect"

	"github.com/invopop/jsonschema"
)

type Params map[string]any

type Parameter struct {
	Name        string
	Description string
	Schema      jsonschema.Schema
}

func GetParameter[T any](params Params, name string) (T, error) {
	var zero T
	value, exists := params[name]
	if !exists {
		return zero, fmt.Errorf("parameter `%s` not found", name)
	}

	if v, ok := value.(T); ok {
		return v, nil
	}

	valueVal := reflect.ValueOf(value)
	zeroType := reflect.TypeOf(zero)

	if valueVal.Kind() == reflect.Float64 && isNumericKind(zeroType.Kind()) {
		converted := reflect.New(zeroType).Elem()
		switch zeroType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			converted.SetInt(int64(valueVal.Float()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			converted.SetUint(uint64(valueVal.Float()))
		default:
			return zero, fmt.Errorf("parameter `%s` type mismatch: cannot convert float64 to `%T`", name, zero)
		}
		return converted.Interface().(T), nil
	}

	return zero, fmt.Errorf("parameter `%s` expected `%T` but got `%T`", name, zero, value)
}

func isNumericKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}
