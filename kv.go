package staert

import (
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

// KvSource impement Source
// It handels all mapstructure features(Squashed Embeded Sub-Structures, Maps, Pointers)
// It support Slices (and maybe Arraies). They must be sort in the KvStore like this :
// Key : ".../[sliceIndex]" -> Value
type KvSource struct {
	store.Store
	Prefix string
}

// NewKvSource creats a new KvSource
func NewKvSource(backend store.Backend, addrs []string, options *store.Config, prefix string) (*KvSource, error) {
	store, err := libkv.NewStore(backend, addrs, options)
	return &KvSource{Store: store, Prefix: prefix}, err
}

// Parse use libkv and mapstructure to fill the structure
func (kv *KvSource) Parse(cmd *flaeg.Command) (*flaeg.Command, error) {
	pairs, err := kv.List(kv.Prefix)
	if err != nil {
		return nil, err
	}
	mapstruct, err := generateMapstructure(pairs, kv.Prefix)
	if err != nil {
		return nil, err
	}
	configDecoder := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           cmd.Config,
		WeaklyTypedInput: true,
		DecodeHook:       decodeHook,
	}
	decoder, err := mapstructure.NewDecoder(configDecoder)
	if err != nil {
		return nil, err
	}
	if err := decoder.Decode(mapstruct); err != nil {
		return nil, err
	}
	return cmd, nil
}

func generateMapstructure(pairs []*store.KVPair, prefix string) (map[string]interface{}, error) {
	raw := make(map[string]interface{})
	for _, p := range pairs {
		// Trim the prefix off our key first
		key := strings.TrimPrefix(p.Key, prefix)
		raw, err := processKV(key, string(p.Value), raw)
		if err != nil {
			return raw, err
		}

	}
	return raw, nil
}

func processKV(key string, v string, raw map[string]interface{}) (map[string]interface{}, error) {
	// Determine what map we're writing the value to. We split by '/'
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
