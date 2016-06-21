package staert

import (
	"errors"
	"fmt"
	"github.com/containous/flaeg"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// KvSource implements Source
// It handles all mapstructure features(Squashed Embeded Sub-Structures, Maps, Pointers)
// It supports Slices (and maybe Arraies). They must be sorted in the KvStore like this :
// Key : ".../[sliceIndex]" -> Value
type KvSource struct {
	store.Store
	Prefix string // like this "prefix" (without the /)
}

// NewKvSource creates a new KvSource
func NewKvSource(backend store.Backend, addrs []string, options *store.Config, prefix string) (*KvSource, error) {
	store, err := libkv.NewStore(backend, addrs, options)
	return &KvSource{Store: store, Prefix: prefix}, err
}

// Parse uses libkv and mapstructure to fill the structure
func (kv *KvSource) Parse(cmd *flaeg.Command) (*flaeg.Command, error) {
	err := kv.LoadConfig(cmd.Config)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

// LoadConfig loads data from the KV Store into the config structure (given by reference)
func (kv *KvSource) LoadConfig(config interface{}) error {
	pairs, err := kv.List(kv.Prefix)
	if err != nil {
		return err
	}
	mapstruct, err := generateMapstructure(pairs, kv.Prefix)
	if err != nil {
		return err
	}
	configDecoder := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           config,
		WeaklyTypedInput: true,
		DecodeHook:       decodeHook,
	}
	decoder, err := mapstructure.NewDecoder(configDecoder)
	if err != nil {
		return err
	}
	if err := decoder.Decode(mapstruct); err != nil {
		return err
	}
	return nil
}

func generateMapstructure(pairs []*store.KVPair, prefix string) (map[string]interface{}, error) {
	raw := make(map[string]interface{})
	for _, p := range pairs {
		// Trim the prefix off our key first
		key := strings.TrimPrefix(p.Key, prefix+"/")
		raw, err := processKV(key, string(p.Value), raw)
		if err != nil {
			return raw, err
		}

	}
	return raw, nil
}

func processKV(key string, v string, raw map[string]interface{}) (map[string]interface{}, error) {
	// Determine which map we're writing the value to. We split by '/'
	// to determine any sub-maps that need to be created.
	m := raw
	children := strings.Split(key, "/")
	if len(children) > 0 {
		key = children[len(children)-1]
		children = children[:len(children)-1]
		for _, child := range children {
			if m[child] == nil {
				m[child] = make(map[string]interface{})
			}
			subm, ok := m[child].(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("child is both a data item and dir: %s", child)
			}
			m = subm
		}
	}
	m[key] = string(v)
	return raw, nil
}

func decodeHook(fromType reflect.Type, toType reflect.Type, data interface{}) (interface{}, error) {
	// TODO : Array support
	switch toType.Kind() {
	case reflect.Slice:
		if fromType.Kind() == reflect.Map {
			// Type assertion
			dataMap, ok := data.(map[string]interface{})
			if !ok {
				return data, fmt.Errorf("input data is not a map : %#v", data)
			}
			// Sorting map
			indexes := make([]int, len(dataMap))
			i := 0
			for k := range dataMap {
				ind, err := strconv.Atoi(k)
				if err != nil {
					return dataMap, err
				}
				indexes[i] = ind
				i++
			}
			sort.Ints(indexes)
			// Building slice
			dataOutput := make([]interface{}, i)
			i = 0
			for _, k := range indexes {
				dataOutput[i] = dataMap[strconv.Itoa(k)]
				i++
			}

			return dataOutput, nil
		}
	}
	return data, nil
}

// StoreConfig stores the config into the KV Store
func (kv *KvSource) StoreConfig(config interface{}) error {
	kvMap := map[string]string{}
	if err := collateKvRecursive(reflect.ValueOf(config), kvMap, kv.Prefix); err != nil {
		return err
	}
	for k, v := range kvMap {
		if err := kv.Put(k, []byte(v), nil); err != nil {
			return err
		}
	}
	return nil
}

func collateKvRecursive(objValue reflect.Value, kv map[string]string, key string) error {
	name := key
	kind := objValue.Kind()
	switch kind {
	case reflect.Struct:
		for i := 0; i < objValue.NumField(); i++ {
			objType := objValue.Type()
			if objType.Field(i).Name[:1] != strings.ToUpper(objType.Field(i).Name[:1]) {
				//if unexported field
				continue
			}
			squashed := false
			if objType.Field(i).Anonymous {
				if objValue.Field(i).Kind() == reflect.Struct {
					tags := objType.Field(i).Tag
					if strings.Contains(string(tags), "squash") {
						squashed = true
					}
				}
			}
			if squashed {
				if err := collateKvRecursive(objValue.Field(i), kv, name); err != nil {
					return err
				}
			} else {
				fieldName := objType.Field(i).Name
				//useless if not empty Prefix is required ?
				if len(key) == 0 {
					name = strings.ToLower(fieldName)
				} else {
					name = key + "/" + strings.ToLower(fieldName)
				}

				if err := collateKvRecursive(objValue.Field(i), kv, name); err != nil {
					return err
				}
			}
		}

	case reflect.Ptr:
		if !objValue.IsNil() {
			if err := collateKvRecursive(objValue.Elem(), kv, name); err != nil {
				return err
			}
		}

	case reflect.Map:
		for _, k := range objValue.MapKeys() {
			if k.Kind() == reflect.Struct {
				return errors.New("Struct as key not supported")
			}
			name = key + "/" + fmt.Sprint(k)
			if err := collateKvRecursive(objValue.MapIndex(k), kv, name); err != nil {
				return err
			}
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < objValue.Len(); i++ {
			name = key + "/" + strconv.Itoa(i)
			if err := collateKvRecursive(objValue.Index(i), kv, name); err != nil {
				return err
			}
		}
	case reflect.Interface, reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
		if _, ok := kv[name]; ok {
			return errors.New("key already exists: " + name)
		}
		kv[name] = fmt.Sprint(objValue)

	default:
		return fmt.Errorf("Kind %s not supported", kind.String())
	}
	return nil
}
