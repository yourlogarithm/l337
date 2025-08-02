package tools

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"runtime"
	"strings"

	"github.com/invopop/jsonschema"
)

type Tool struct {
	callable any

	Name        string
	Description string
	Parameters  []Parameter
}

type Element struct {
	Type   string
	Nested *Element
	Schema *jsonschema.Schema
}

type Parameter struct {
	Name    string
	Element Element
}

func (t *Tool) Call() {

}

func NewTool(fn any, name, description string, parameters []Parameter) (Tool, error) {
	if err := checkFn(fn); err != nil {
		return Tool{}, fmt.Errorf("NewTool: %w", err)
	}

	return Tool{
		callable:    fn,
		Name:        name,
		Description: description,
		Parameters:  parameters,
	}, nil
}

// CreateToolFromFunc creates a Tool from a function, extracting metadata via reflection and source parsing.
func CreateToolFromFunc(fn any) (Tool, error) {
	if err := checkFn(fn); err != nil {
		return Tool{}, fmt.Errorf("CreateToolFromFunc: %w", err)
	}

	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()

	// Get function pointer, runtime info
	fnPtr := fnValue.Pointer()
	fnInfo := runtime.FuncForPC(fnPtr)
	if fnInfo == nil {
		return Tool{}, fmt.Errorf("CreateToolFromFunc: unable to get function info")
	}
	file, _ := fnInfo.FileLine(fnPtr)
	funcName := fnInfo.Name()
	if idx := strings.LastIndex(funcName, "."); idx != -1 {
		funcName = funcName[idx+1:]
	}

	// Parse source file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return Tool{}, fmt.Errorf("CreateToolFromFunc: failed to parse source: %v", err)
	}

	// Extract doc comment and parameter names
	var doc string
	var paramNames []string
	for _, decl := range node.Decls {
		if fnDecl, ok := decl.(*ast.FuncDecl); ok && fnDecl.Name.Name == funcName {
			if fnDecl.Doc != nil {
				doc = strings.TrimSpace(fnDecl.Doc.Text())
			}
			for _, field := range fnDecl.Type.Params.List {
				if len(field.Names) == 0 {
					return Tool{}, fmt.Errorf("CreateToolFromFunc: function %s has unnamed parameters", funcName)
				}
				for _, ident := range field.Names {
					paramNames = append(paramNames, ident.Name)
				}
			}
			break
		}
	}

	// Match parameter types from reflect
	var parameters []Parameter
	for i := 0; i < fnType.NumIn(); i++ {
		if i >= len(paramNames) {
			return Tool{}, fmt.Errorf("CreateToolFromFunc: not enough parameter names provided for function %s", funcName)
		}
		name := paramNames[i]
		param, err := getParameter(fnType.In(i), name)
		if err != nil {
			return Tool{}, fmt.Errorf("CreateToolFromFunc: failed to get parameter %s for function %s: %v", name, funcName, err)
		}

		parameters = append(parameters, param)
	}

	return Tool{
		callable:    fn,
		Name:        funcName,
		Description: doc,
		Parameters:  parameters,
	}, nil
}

func checkFn(fn any) error {
	fnValue := reflect.ValueOf(fn)
	if fnValue.Kind() != reflect.Func {
		return fmt.Errorf("expected a function, got %s", fnValue.Kind())
	}
	return nil
}

func getParameter(paramType reflect.Type, name string) (Parameter, error) {
	element, err := getElement(paramType)
	if err != nil {
		return Parameter{}, fmt.Errorf("getParameter: failed to get element type for %s: %v", name, err)
	}

	return Parameter{
		Name:    name,
		Element: element,
	}, nil
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
