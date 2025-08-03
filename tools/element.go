package tools

import (
	"fmt"
	"reflect"

	"github.com/invopop/jsonschema"
)

type Element struct {
	Type   string
	Nested *Element
	Schema *jsonschema.Schema
}

func getElement(goType reflect.Type) (Element, error) {
	var t string
	var nested *Element
	var schema *jsonschema.Schema

	switch goType.Kind() {
	case reflect.Bool:
		t = "bool"
	case reflect.String:
		t = "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		t = "number"
	case reflect.Array, reflect.Slice:
		t = "array"
		nestedVar, err := getElement(goType.Elem())
		if err != nil {
			return Element{}, fmt.Errorf("getElement: failed to get nested element for %s: %v", goType.Name(), err)
		}
		nested = &nestedVar
	case reflect.Map:
		t = "object"
		if goType.Key().Kind() != reflect.String {
			return Element{}, fmt.Errorf("getElement: map keys must be strings, got %s", goType.Key().Kind())
		}
		nestedVar, err := getElement(goType.Elem())
		if err != nil {
			return Element{}, fmt.Errorf("getElement: failed to get nested element for %s: %v", goType.Name(), err)
		}
		nested = &nestedVar
	case reflect.Struct:
		t = goType.Name()
		schema = jsonschema.Reflect(reflect.New(goType).Interface())
		if schema == nil {
			return Element{}, fmt.Errorf("getElement: failed to reflect schema for struct %s", goType.Name())
		}
	default:
		return Element{}, fmt.Errorf("getElement: unsupported type %s", goType.Kind())
	}

	return Element{Type: t, Nested: nested, Schema: schema}, nil
}
