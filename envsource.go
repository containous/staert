package staert

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/containous/flaeg"
	"github.com/containous/flaeg/parse"
	"github.com/fatih/camelcase"
)

var mutex sync.Mutex

// EnvSource implements SourceLoader
// Enables to populate configuration struct with information extracted from
// process's environment variables. Variables names are like %PREFIX%%SEP%%FIELD_NAME%
// It supports pointer to values and struct, however not slices and arrays..
type EnvSource struct {
	prefix    string
	separator string
	parsers   map[reflect.Type]parse.Parser
}

type envValue struct {
	strValue string
	path     []string
}

// NewEnvSource constructs a new instance of EnvSource
func NewEnvSource(prefix, separator string, parsers map[reflect.Type]parse.Parser) *EnvSource {
	return &EnvSource{prefix: prefix, separator: separator, parsers: parsers}
}

// Parse parse and load config structure
func (e *EnvSource) Parse(cmd *flaeg.Command) (*flaeg.Command, error) {
	return cmd, e.LoadConfig(cmd.Config)
}

// LoadConfig load a configuration from a configuration structure
func (e *EnvSource) LoadConfig(config interface{}) error {
	configVal := reflect.ValueOf(config).Elem()

	values, err := e.analyzeStruct(configVal.Type(), nil)
	if err != nil {
		return err
	}

	return e.assignValues(configVal, values, nil)
}

// Recursively scan the given config structure type information
// and look for defined environment variables.
func (e *EnvSource) analyzeStruct(configType reflect.Type, currentPath []string) ([]*envValue, error) {
	var res []*envValue

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)

		// TODO: Handle this case;
		// find the underlying struct and process it.
		if field.Type.Kind() == reflect.Interface {
			// skip fields of kind interface
			continue
		}

		// If we're facing an embedded struct
		if field.Anonymous {
			values, err := e.analyzeValue(field.Type, currentPath)
			if err != nil {
				return nil, err
			}

			res = append(res, values...)
			continue
		}

		// unexported fields must be handled after embedded structs (field.Anonymous)
		// because the PkgPath is also null for them.
		// ref: https://github.com/golang/go/issues/21122
		if field.PkgPath != "" {
			// field is unexported
			continue
		}

		values, err := e.analyzeValue(field.Type, append(currentPath, field.Name))
		if err != nil {
			return nil, err
		}

		res = append(res, values...)
	}

	return res, nil
}

func (e *EnvSource) analyzeValue(valType reflect.Type, fieldPath []string) ([]*envValue, error) {
	switch valType.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		return e.analyzeIndexedType(valType, fieldPath)
	case reflect.Ptr:
		return e.analyzeValue(valType.Elem(), fieldPath)
	case reflect.Struct:
		return e.analyzeStruct(valType, fieldPath)
	case reflect.Invalid, reflect.Chan, reflect.Func, reflect.Interface, reflect.UnsafePointer:
		return nil, nil
		// Skip these fields, don't throw...
		// TODO : keep track of the fields ignored by the library to be able to list them.
		// err = fmt.Errorf("type %s is not supported by EnvSource. fieldPath : %v", valType.Name(), fieldPath)
	default:
		return e.loadValue(fieldPath), nil
	}
}

func (e *EnvSource) analyzeIndexedType(valType reflect.Type, fieldPath []string) ([]*envValue, error) {
	prefix := e.envVarFromPath(fieldPath)
	vars := e.envVarsWithPrefix(prefix)
	nextKeys := unique(e.nextLevelKeys(prefix, vars))

	var res []*envValue
	for _, varName := range nextKeys {
		key := e.keyFromEnvVar(varName, prefix)

		// If we're on an Int based key, we need to be able to convert
		// detected key to an int
		if valType.Kind() == reflect.Array || valType.Kind() == reflect.Slice {
			index, err := strconv.Atoi(key)
			if err != nil {
				return nil, fmt.Errorf("invalid key %q for variable %q: %v", key, varName, err)
			}

			if valType.Kind() == reflect.Array && index >= valType.Len() {
				return nil, fmt.Errorf("detected key (%s) from variable %s is >= to array length %d",
					key, varName, valType.Len())
			}
		}

		valPath := append(fieldPath, key)
		keyValues, err := e.analyzeValue(valType.Elem(), valPath)
		if err != nil {
			return nil, err
		}

		res = append(res, keyValues...)
	}

	return res, nil
}

func (e *EnvSource) loadValue(fieldPath []string) []*envValue {
	variableName := e.envVarFromPath(fieldPath)

	value, ok := os.LookupEnv(variableName)
	if !ok {
		return nil
	}

	clone := make([]string, len(fieldPath))
	copy(clone, fieldPath)

	return []*envValue{{strValue: value, path: clone}}
}

func (e *EnvSource) assignValues(configVal reflect.Value, envValues []*envValue, currentPath []string) error {
	if len(currentPath) > 0 {
		envValues = filterEnvVarWithPrefix(envValues, currentPath)
	}

	if configVal.Kind() == reflect.Ptr {
		if configVal.IsNil() {
			configVal.Set(reflect.New(configVal.Type().Elem()))
		}
		return e.assignValues(configVal.Elem(), envValues, nil)
	}

	for _, v := range envValues {
		fieldVal := configVal.FieldByName(v.path[0])
		if !fieldVal.IsValid() {
			// skip field that are found invalid
			continue
		}
		switch fieldVal.Kind() {

		case reflect.Ptr, reflect.Struct:
			err := e.assignValues(fieldVal, []*envValue{v}, []string{v.path[0]})
			if err != nil {
				return err
			}
		case reflect.Slice, reflect.Array:
			err := e.assignArrays(fieldVal, envValues, v)
			if err != nil {
				return err
			}
		case reflect.Map:
			key := v.path[1]
			val := v.strValue

			mapType := fieldVal.Type()
			elemType := mapType.Elem()

			if elemType.Kind() == reflect.Struct {
				elem := reflect.New(elemType).Elem()

				err := e.assignValues(elem, envValues, v.path[:2])
				if err != nil {
					return err
				}

				err = e.assignMap(fieldVal, key, elem)
				if err != nil {
					return err
				}
			} else {
				err := e.assignMap(fieldVal, key, reflect.ValueOf(val))
				if err != nil {
					return err
				}
			}
		default:
			value, err := e.getParsedValue(fieldVal.Type(), v.strValue)
			if err != nil {
				return err
			}

			fieldVal.Set(*value)
		}
	}

	return nil
}

func (e *EnvSource) assignMap(fieldVal reflect.Value, key string, val reflect.Value) error {
	mapType := fieldVal.Type()
	if fieldVal.IsNil() {
		fieldVal.Set(reflect.MakeMap(mapType))
	}

	elemType := mapType.Elem()
	keyType := mapType.Key()

	parsedKey, err := e.getParsedValue(keyType, key)
	if err != nil {
		return err
	}

	if val.Kind() == reflect.String {
		parsedVal, err := e.getParsedValue(elemType, val.String())
		if err != nil {
			return err
		}

		fieldVal.SetMapIndex(*parsedKey, *parsedVal)
	} else {
		fieldVal.SetMapIndex(*parsedKey, val)
	}

	return nil
}

func (e *EnvSource) assignArrays(fieldVal reflect.Value, envValues []*envValue, currentEnvValue *envValue) error {
	arrayType := fieldVal.Type()
	slice := reflect.Zero(reflect.SliceOf(arrayType.Elem()))
	if !fieldVal.IsNil() {
		slice = reflect.Indirect(fieldVal)
		fieldVal.Set(slice)
	}
	elemType := arrayType.Elem()

	if elemType.Kind() != reflect.Struct && elemType.Kind() != reflect.Ptr {
		parsedVal, err := e.getParsedValue(elemType, currentEnvValue.strValue)
		if err != nil {
			return err
		}

		slice = reflect.Append(slice, *parsedVal)
		fieldVal.Set(slice)
	} else {
		if index, err := strconv.Atoi(currentEnvValue.path[1]); err == nil {
			// grow the slice if needed
			if slice.Len() <= index {
				newSlice := reflect.MakeSlice(slice.Type(), index+1, index+1)
				reflect.Copy(newSlice, slice)
				slice = newSlice
			}

			// get item at env value specified index.
			existingValue := slice.Index(index)
			err = e.assignValues(existingValue, envValues, currentEnvValue.path[:2])
			if err != nil {
				return err
			}

			fieldVal.Set(slice)
		}
	}
	return nil
}

func (e *EnvSource) getParsedValue(valType reflect.Type, strValue string) (*reflect.Value, error) {
	mutex.Lock()
	defer mutex.Unlock()

	parser, ok := e.parsers[valType]
	if !ok {
		return nil, fmt.Errorf("parser not found for type %v and value %q", valType, strValue)
	}

	err := parser.Set(strValue)
	if err != nil {
		return nil, err
	}

	value := reflect.ValueOf(parser.Get())
	return &value, nil
}

func (e *EnvSource) nextLevelKeys(prefix string, envVars []string) []string {
	res := make([]string, 0, len(envVars))

	for _, envVar := range envVars {
		nextKey := strings.SplitN(
			strings.TrimPrefix(envVar, prefix+e.separator),
			e.separator, 2,
		)[0]
		res = append(res, prefix+e.separator+nextKey)
	}

	return res
}

func (e *EnvSource) envVarsWithPrefix(prefix string) []string {
	var res []string

	for _, rawVar := range os.Environ() {
		if strings.HasPrefix(rawVar, prefix) {
			varName := strings.SplitN(rawVar, "=", 2)[0]
			res = append(res, varName)
		}
	}

	return res
}

func (e *EnvSource) keyFromEnvVar(fullVar, prefix string) string {
	return strings.ToLower(
		strings.SplitN(
			strings.TrimPrefix(fullVar, prefix+e.separator),
			e.separator, 2,
		)[0],
	)
}

func (e *EnvSource) envVarFromPath(currentPath []string) string {
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
	collector := make(map[string]struct{})

	var res []string
	for _, v := range in {
		if _, ok := collector[v]; ok {
			continue
		}

		collector[v] = struct{}{}
		res = append(res, v)
	}

	return res
}

func filterEnvVarWithPrefix(envValues []*envValue, startFilter []string) []*envValue {
	startFilterPath := strings.Join(startFilter, "")

	var res []*envValue
	for _, currentEnvValue := range envValues {
		if len(currentEnvValue.path) >= len(startFilter) {
			currentPath := strings.Join(currentEnvValue.path[0:len(startFilter)], "")
			if startFilterPath == currentPath {
				newEnvValue := &envValue{
					strValue: currentEnvValue.strValue,
					path:     currentEnvValue.path[len(startFilter):],
				}
				res = append(res, newEnvValue)
			}
		}
	}
	return res
}
