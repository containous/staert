package staert

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/containous/flaeg"
	"github.com/containous/flaeg/parse"
	"github.com/fatih/camelcase"
)

// Loader interface
// A Loader is an object that can be used to load a configuration from a
// configuration structure
type Loader interface {
	LoadConfig(config interface{}) error
}

// SourceLoader : can be both a staert.Source and a staert.Loader
type SourceLoader interface {
	Loader
	Source
}

// envSource implements SourceLoader
// Enables to populate configuration struct with informations extracted from
// process's environment variables. Variables names are like %PREFIX%%SEP%%FIELD_NAME%
// It supports pointer to values and struct, however not slices and arrays..
type envSource struct {
	prefix    string
	separator string
	parsers   map[reflect.Type]parse.Parser
}

// NewEnvSource constructs a new instance of envSource
func NewEnvSource(prefix, separator string, parsers map[reflect.Type]parse.Parser) SourceLoader {
	return &envSource{prefix, separator, parsers}
}

// Parse parse and load config structure
func (e *envSource) Parse(cmd *flaeg.Command) (*flaeg.Command, error) {
	return cmd, e.LoadConfig(cmd.Config)
}

func (e *envSource) LoadConfig(config interface{}) error {
	configVal := reflect.ValueOf(config).Elem()

	values, err := e.analyzeStruct(configVal.Type(), []string{})

	if err != nil {
		return err
	}

	return e.assignValues(configVal, values)
}

type envValue struct {
	StrValue string
	Path     path
}

type path []string

func (p path) clone() []string {
	res := make([]string, len(p))
	copy(res, p)
	return res
}

// Recursively scan the given config structure type information
// and look for defined environment variables.
func (e *envSource) analyzeStruct(configType reflect.Type, currentPath path) ([]*envValue, error) {
	res := []*envValue{}

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)

		//TODO: Handle this case;
		//find the find the underlying struct and process it.
		if field.Type.Kind() == reflect.Interface {
			continue //skip fields of kind interface
		}
		// If we're facing an embedded struct
		if field.Anonymous {
			values, err := e.analyzeStruct(field.Type, currentPath)

			if err != nil {
				return []*envValue{}, err
			}

			res = append(res, values...)
			continue
		}

		// unexported fields must be handled after embdedded structs (field.Anonymous)
		// because the PkgPath is also null for them.
		// ref: https://github.com/golang/go/issues/21122
		if field.PkgPath != "" { //field is unexported
			continue
		}

		values, err := e.analyzeValue(field.Type, append(currentPath, field.Name))

		if err != nil {
			return []*envValue{}, err
		}

		res = append(res, values...)
	}

	return res, nil
}

func (e *envSource) analyzeValue(valType reflect.Type, fieldPath path) ([]*envValue, error) {
	var (
		res []*envValue
		err error
	)
	switch valType.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		res, err = e.analyzeIndexedType(valType, fieldPath)
	case reflect.Ptr:
		res, err = e.analyzeValue(valType.Elem(), fieldPath)
	case reflect.Struct:
		res, err = e.analyzeStruct(valType, fieldPath)
	case reflect.Invalid, reflect.Chan, reflect.Func, reflect.Interface, reflect.UnsafePointer:
		err = fmt.Errorf("type %s is not supported by EnvSource", valType.Name())

	default:
		res = e.loadValue(fieldPath)
	}

	return res, err
}

func (e *envSource) analyzeIndexedType(valType reflect.Type, fieldPath path) ([]*envValue, error) {
	var (
		res []*envValue
	)

	prefix := e.envVarFromPath(fieldPath)
	vars := e.envVarsWithPrefix(prefix)
	nextKeys := unique(e.nextLevelKeys(prefix, vars))

	for _, varName := range nextKeys {
		key := e.keyFromEnvVar(varName, prefix)

		// If we're on an Int based key, we need to be able to convert
		// detected key to an int
		if valType.Kind() == reflect.Array ||
			valType.Kind() == reflect.Slice {
			index, err := strconv.Atoi(key)

			if err != nil {
				return res, fmt.Errorf(
					key,
					varName,
				)
			}

			if valType.Kind() == reflect.Array &&
				index >= valType.Len() {
				return res, fmt.Errorf(
					"Detected key (%s) from variable %s is >= to array length %d",
					key,
					varName,
					valType.Len(),
				)
			}
		}

		valPath := append(fieldPath, key)
		keyValues, err := e.analyzeValue(valType.Elem(), valPath)

		if err != nil {
			return res, err
		}

		res = append(res, keyValues...)
	}

	return res, nil
}

func (e *envSource) loadValue(fieldPath path) []*envValue {
	variableName := e.envVarFromPath(fieldPath)

	value, ok := os.LookupEnv(variableName)

	if !ok {
		return []*envValue{}
	}

	return []*envValue{&envValue{value, fieldPath.clone()}}
}

func (e *envSource) assignValues(configVal reflect.Value, values []*envValue) error {
	for _, v := range values {
		currentValue := configVal
		for _, p := range v.Path {
			currentValue = currentValue.FieldByName(p)

			if e.needsAllocation(currentValue) {
				var err error
				currentValue, err = e.allocate(currentValue)
				if err != nil {
					return err
				}
			}
		}

		if err := e.setValue(currentValue, v.StrValue); err != nil {
			return err
		}
	}
	return nil
}

func (e *envSource) needsAllocation(value reflect.Value) bool {
	// TODO
	return false
}

func (e *envSource) allocate(value reflect.Value) (reflect.Value, error) {
	// TODO
	return value, nil
}

func (e *envSource) setValue(value reflect.Value, strValue string) error {

	if !value.CanSet() {
		return fmt.Errorf(
			"Value [%v] cannot be set",
			value,
		)
	}

	parser, ok := e.parsers[value.Type()]

	if !ok {
		return fmt.Errorf(
			"Unsupported type [%s], please consider adding custom parser",
			value.Type().Name(),
		)
	}

	err := parser.Set(strValue)

	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(parser).Elem().Convert(value.Type()))

	return nil
}

func (e *envSource) nextLevelKeys(prefix string, envVars []string) []string {
	res := make([]string, 0, len(envVars))

	for _, envVar := range envVars {
		nextKey := strings.Split(
			strings.TrimPrefix(envVar, prefix+e.separator),
			e.separator,
		)[0]
		res = append(res, prefix+e.separator+nextKey)

	}

	return res
}

func (e *envSource) envVarsWithPrefix(prefix string) []string {
	res := []string{}

	for _, rawVar := range os.Environ() {
		varName := strings.Split(rawVar, "=")[0]
		if strings.HasPrefix(varName, prefix) {
			res = append(res, varName)
		}
	}

	return res
}

func (e *envSource) keyFromEnvVar(fullVar, prefix string) string {
	return strings.ToLower(
		strings.Split(
			strings.TrimPrefix(fullVar, prefix+e.separator),
			e.separator,
		)[0],
	)
}

func (e *envSource) envVarFromPath(currentPath []string) string {
	if e.prefix != "" {
		currentPath = append([]string{e.prefix}, currentPath...)
	}
	s := make([]string, 0, len(currentPath))

	for _, word := range currentPath {
		s = append(s, camelcase.Split(word)...)
	}

	return strings.ToUpper(strings.Join(s, e.separator))
}

func unique(in []string) []string {
	collector := map[string]struct{}{}
	res := []string{}

	for _, v := range in {
		if _, ok := collector[v]; ok {
			continue
		}

		collector[v] = struct{}{}
		res = append(res, v)
	}

	return res
}
