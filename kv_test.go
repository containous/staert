package staert

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/containous/flaeg"
	"github.com/docker/libkv/store"
	"github.com/mitchellh/mapstructure"
)

func TestGenerateMapstructureBasic(t *testing.T) {
	moke := []*store.KVPair{
		{Key: "test/addr", Value: []byte("foo")},
		{Key: "test/child/data", Value: []byte("bar")}}
	prefix := "test"

	output, err := generateMapstructure(moke, prefix)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	expected := map[string]interface{}{
		"addr": "foo",
		"child": map[string]interface{}{
			"data": "bar",
		},
	}
	if !reflect.DeepEqual(expected, output) {
		t.Fatalf("Got: %#v\nExpected: %#v", expected, output)
	}
}

func TestGenerateMapstructureTrivialMap(t *testing.T) {
	moke := []*store.KVPair{
		{Key: "test/vfoo", Value: []byte("foo")},
		{Key: "test/vother/foo", Value: []byte("foo")},
		{Key: "test/vother/bar", Value: []byte("bar")},
	}
	prefix := "test"

	output, err := generateMapstructure(moke, prefix)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	expected := map[string]interface{}{
		"vfoo": "foo",
		"vother": map[string]interface{}{
			"foo": "foo",
			"bar": "bar",
		},
	}
	if !reflect.DeepEqual(expected, output) {
		t.Fatalf("Got: %#v\nExpected: %#v", expected, output)
	}
}

func TestGenerateMapstructureTrivialSlice(t *testing.T) {
	moke := []*store.KVPair{
		{Key: "test/vfoo", Value: []byte("foo")},
		{Key: "test/vother/0", Value: []byte("foo")},
		{Key: "test/vother/1", Value: []byte("bar1")},
		{Key: "test/vother/2", Value: []byte("bar2")},
	}
	prefix := "test"

	output, err := generateMapstructure(moke, prefix)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	expected := map[string]interface{}{
		"vfoo": "foo",
		"vother": map[string]interface{}{
			"0": "foo",
			"1": "bar1",
			"2": "bar2",
		},
	}
	if !reflect.DeepEqual(expected, output) {
		t.Fatalf("Got: %#v\nExpected: %#v", expected, output)
	}
}

func TestGenerateMapstructureNotTrivialSlice(t *testing.T) {
	moke := []*store.KVPair{
		{Key: "test/vfoo", Value: []byte("foo")},
		{Key: "test/vother/0/foo1", Value: []byte("bar")},
		{Key: "test/vother/0/foo2", Value: []byte("bar")},
		{Key: "test/vother/1/bar1", Value: []byte("foo")},
		{Key: "test/vother/1/bar2", Value: []byte("foo")},
	}
	prefix := "test"

	output, err := generateMapstructure(moke, prefix)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	expected := map[string]interface{}{
		"vfoo": "foo",
		"vother": map[string]interface{}{
			"0": map[string]interface{}{
				"foo1": "bar",
				"foo2": "bar",
			},
			"1": map[string]interface{}{
				"bar1": "foo",
				"bar2": "foo",
			},
		},
	}
	if !reflect.DeepEqual(expected, output) {
		t.Fatalf("Got: %#v\nExpected: %#v", expected, output)
	}
}

func TestDecodeHookSlice(t *testing.T) {
	data := map[string]interface{}{
		"10": map[string]interface{}{
			"bar1": "bar1",
			"bar2": "bar2",
		},
		"2": map[string]interface{}{
			"bar1": "foo1",
			"bar2": "foo2",
		},
	}
	output, err := decodeHook(reflect.TypeOf(data), reflect.TypeOf([]string{}), data)
	if err != nil {
		t.Fatalf("Error : %s", err)
	}

	expected := []interface{}{
		map[string]interface{}{
			"bar1": "foo1",
			"bar2": "foo2",
		},
		map[string]interface{}{
			"bar1": "bar1",
			"bar2": "bar2",
		},
	}
	if !reflect.DeepEqual(expected, output) {
		t.Fatalf("Got: %#v\nExpected: %#v", expected, output)
	}
}

type BasicStruct struct {
	Bar1 string
	Bar2 string
}

type SliceStruct []BasicStruct

type Test struct {
	Vfoo   string
	Vother SliceStruct
}

func TestIntegrationMapstructureWithDecodeHook(t *testing.T) {
	input := map[string]interface{}{
		"vfoo": "foo",
		"vother": map[string]interface{}{
			"10": map[string]interface{}{
				"bar1": "bar1",
				"bar2": "bar2",
			},
			"2": map[string]interface{}{
				"bar1": "foo1",
				"bar2": "foo2",
			},
		},
	}
	var config Test

	//test
	configDecoder := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           &config,
		WeaklyTypedInput: true,
		DecodeHook:       decodeHook,
	}
	decoder, err := mapstructure.NewDecoder(configDecoder)
	if err != nil {
		t.Fatalf("got an err: %v", err)
	}
	if err := decoder.Decode(input); err != nil {
		t.Fatalf("got an err: %v", err)
	}

	expected := Test{
		Vfoo: "foo",
		Vother: SliceStruct{
			BasicStruct{
				Bar1: "foo1",
				Bar2: "foo2",
			},
			BasicStruct{
				Bar1: "bar1",
				Bar2: "bar2",
			},
		},
	}

	if !reflect.DeepEqual(expected, config) {
		t.Fatalf("Got: %#v\nExpected: %#v", expected, config)
	}
}

func TestKvSourceEmpty(t *testing.T) {
	//Init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: flaeg.Duration(time.Second),
	}

	//Test
	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: config,
		Run: func() error { return nil },
	}
	s := NewStaert(rootCmd)
	kv := &KvSource{
		&Mock{
			KVPairs: []*store.KVPair{},
		},
		"test/",
	}
	s.AddSource(kv)

	_, err := s.LoadConfig()
	if err != nil {
		t.Fatalf("Error %s", err)
	}

	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: flaeg.Duration(time.Second),
	}

	if !reflect.DeepEqual(expected, rootCmd.Config) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected, rootCmd.Config)
	}
}

func TestGenerateMapstructureTrivial(t *testing.T) {
	input := []*store.KVPair{
		{Key: "test/ptrstruct1/s1int", Value: []byte("28")},
		{Key: "test/durationfield", Value: []byte("28")},
	}
	prefix := "test"
	output, err := generateMapstructure(input, prefix)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	expected := map[string]interface{}{
		"durationfield": "28",
		"ptrstruct1": map[string]interface{}{
			"s1int": "28",
		},
	}
	if !reflect.DeepEqual(expected, output) {
		actualJSON, err := json.Marshal(output)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		expectedJSON, err := json.Marshal(expected)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		t.Fatalf("\nexpected\t: %s\ngot\t\t\t: %s\n", expectedJSON, actualJSON)
	}
}

func TestIntegrationMapstructureWithDecodeHookPointer(t *testing.T) {
	mapstruct := map[string]interface{}{
		"durationfield": "28",
		"ptrstruct1": map[string]interface{}{
			"s1int": "28",
		},
	}
	config := StructPtr{}

	//test
	configDecoder := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           &config,
		WeaklyTypedInput: true,
		DecodeHook:       decodeHook,
		//TODO : ZeroFields:       false, doesn't work
	}
	decoder, err := mapstructure.NewDecoder(configDecoder)
	if err != nil {
		t.Fatalf("got an err: %v", err)
	}
	if err := decoder.Decode(mapstruct); err != nil {
		t.Fatalf("got an err: %v", err)
	}

	expected := StructPtr{
		PtrStruct1: &Struct1{
			S1Int: 28,
		},
		DurationField: flaeg.Duration(28 * time.Nanosecond),
	}

	if !reflect.DeepEqual(expected, config) {
		actualJSON, err := json.Marshal(config)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		expectedJSON, err := json.Marshal(expected)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		t.Fatalf("\nexpected\t: %s\ngot\t\t\t: %s\n", expectedJSON, actualJSON)
	}
}

func TestIntegrationMapstructureInitiatedPtrReset(t *testing.T) {
	mapstruct := map[string]interface{}{
		// "durationfield": "28",
		"ptrstruct1": map[string]interface{}{
			"s1int": "24",
		},
	}
	config := StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: flaeg.Duration(28 * time.Second),
	}

	//test
	configDecoder := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           &config,
		WeaklyTypedInput: true,
		DecodeHook:       decodeHook,
	}
	decoder, err := mapstructure.NewDecoder(configDecoder)
	if err != nil {
		t.Fatalf("got an err: %v", err)
	}
	if err := decoder.Decode(mapstruct); err != nil {
		t.Fatalf("got an err: %v", err)
	}

	expected := StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    24,
			S1String: "S1StringInitConfig",
		},
		DurationField: flaeg.Duration(28 * time.Second),
	}

	if !reflect.DeepEqual(expected, config) {
		actualJSON, err := json.Marshal(config)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		expectedJSON, err := json.Marshal(expected)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		t.Fatalf("\nexpected\t: %s\ngot\t\t\t: %s\n", expectedJSON, actualJSON)
	}
}

func TestParseKvSourceTrivial(t *testing.T) {
	//Init
	config := StructPtr{}

	//Test
	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                &config,
		DefaultPointersConfig: &config,
		Run: func() error { return nil },
	}
	kv := &KvSource{
		&Mock{
			KVPairs: []*store.KVPair{
				{Key: "test/ptrstruct1/s1int", Value: []byte("28")},
				{Key: "test/durationfield", Value: []byte("28")},
			},
		},
		"test",
	}
	if _, err := kv.Parse(rootCmd); err != nil {
		t.Fatalf("Error %s", err)
	}

	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int: 28,
		},
		DurationField: flaeg.Duration(28 * time.Nanosecond),
	}

	if !reflect.DeepEqual(expected, rootCmd.Config) {
		actualJSON, err := json.Marshal(rootCmd.Config)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		expectedJSON, err := json.Marshal(expected)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		t.Fatalf("\nexpected\t: %s\ngot\t\t\t: %s\n", expectedJSON, actualJSON)
	}
}

func TestLoadConfigKvSourceNestedPtrsNil(t *testing.T) {
	//Init
	config := &StructPtr{}

	//Test
	kv := &KvSource{
		&Mock{
			KVPairs: []*store.KVPair{
				{Key: "prefix/ptrstruct1/s1int", Value: []byte("1")},
				{Key: "prefix/ptrstruct1/s1string", Value: []byte("S1StringInitConfig")},
				{Key: "prefix/ptrstruct1/s1bool", Value: []byte("false")},
				{Key: "prefix/ptrstruct1/s1ptrstruct3/s3float64", Value: []byte("0")},
				{Key: "prefix/durationfield", Value: []byte("21000000000")},
			},
		},
		"prefix",
	}
	if err := kv.LoadConfig(config); err != nil {
		t.Fatalf("Error %s", err)
	}

	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:        1,
			S1String:     "S1StringInitConfig",
			S1PtrStruct3: &Struct3{},
		},
		DurationField: flaeg.Duration(21 * time.Second),
	}

	if !reflect.DeepEqual(expected, config) {
		actualJSON, err := json.Marshal(config)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		expectedJSON, err := json.Marshal(expected)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		t.Fatalf("\nexpected\t: %s\ngot\t\t\t: %s\n", expectedJSON, actualJSON)
	}
}

func TestParseKvSourceNestedPtrsNil(t *testing.T) {
	//Init
	config := StructPtr{}

	//Test
	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                &config,
		DefaultPointersConfig: &config,
		Run: func() error { return nil },
	}
	kv := &KvSource{
		&Mock{
			KVPairs: []*store.KVPair{
				{Key: "prefix/ptrstruct1/s1int", Value: []byte("1")},
				{Key: "prefix/ptrstruct1/s1string", Value: []byte("S1StringInitConfig")},
				{Key: "prefix/ptrstruct1/s1bool", Value: []byte("false")},
				{Key: "prefix/ptrstruct1/s1ptrstruct3/s3float64", Value: []byte("0")},
				{Key: "prefix/durationfield", Value: []byte("21000000000")},
			},
		},
		"prefix",
	}
	if _, err := kv.Parse(rootCmd); err != nil {
		t.Fatalf("Error %s", err)
	}

	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:        1,
			S1String:     "S1StringInitConfig",
			S1PtrStruct3: &Struct3{},
		},
		DurationField: flaeg.Duration(21 * time.Second),
	}

	if !reflect.DeepEqual(expected, rootCmd.Config) {
		actualJSON, err := json.Marshal(rootCmd.Config)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		expectedJSON, err := json.Marshal(expected)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		t.Fatalf("\nexpected\t: %s\ngot\t\t\t: %s\n", expectedJSON, actualJSON)
	}
}

func TestParseKvSourceMap(t *testing.T) {
	//Init
	config := &struct {
		Vmap map[string]int
	}{}

	//Test
	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: config,
		Run: func() error { return nil },
	}
	kv := &KvSource{
		&Mock{
			KVPairs: []*store.KVPair{
				{Key: "prefix/vmap/toto", Value: []byte("1")},
				{Key: "prefix/vmap/tata", Value: []byte("2")},
				{Key: "prefix/vmap/titi", Value: []byte("3")},
			},
		},
		"prefix",
	}
	if _, err := kv.Parse(rootCmd); err != nil {
		t.Fatalf("Error %v", err)
	}

	expected := &struct {
		Vmap map[string]int
	}{
		Vmap: map[string]int{
			"toto": 1,
			"tata": 2,
			"titi": 3,
		},
	}

	if !reflect.DeepEqual(expected, rootCmd.Config) {
		t.Fatalf("\nexpected\t: %#v\ngot\t\t\t: %#v\n", expected, rootCmd.Config)
	}
}

func TestCollateKvPairsBasic(t *testing.T) {
	//init
	config := &struct {
		Vstring string
		Vint    int
		Vuint   uint
		Vbool   bool
		Vfloat  float64
		Vextra  string
		vsilent bool
		Vdata   interface{}
	}{
		Vstring: "tata",
		Vint:    -15,
		Vuint:   51,
		Vbool:   true,
		Vfloat:  1.5,
		Vextra:  "toto",
		vsilent: true, //Unexported : must not be in the map
		Vdata:   42,
	}

	// test
	kv := map[string]string{}
	if err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix"); err != nil {
		t.Fatalf("Error : %v", err)
	}

	expected := map[string]string{
		"prefix/vbool":   "true",
		"prefix/vfloat":  "1.5",
		"prefix/vextra":  "toto",
		"prefix/vdata":   "42",
		"prefix/vstring": "tata",
		"prefix/vint":    "-15",
		"prefix/vuint":   "51",
	}
	if !reflect.DeepEqual(kv, expected) {
		t.Fatalf("Got: %s\nExpected: %s", kv, expected)
	}
}

func TestCollateKvPairsNestedPointers(t *testing.T) {
	//init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:        1,
			S1String:     "S1StringInitConfig",
			S1PtrStruct3: &Struct3{},
		},
		DurationField: flaeg.Duration(21 * time.Second),
	}

	//test
	kv := map[string]string{}
	if err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix"); err != nil {
		t.Fatalf("Error : %v", err)
	}

	//check
	expected := map[string]string{
		"prefix/ptrstruct1/s1int":                  "1",
		"prefix/ptrstruct1/s1string":               "S1StringInitConfig",
		"prefix/ptrstruct1/s1bool":                 "false",
		"prefix/ptrstruct1/s1ptrstruct3/":          "",
		"prefix/ptrstruct1/s1ptrstruct3/s3float64": "0",
		"prefix/durationfield":                     "21000000000",
	}
	if !reflect.DeepEqual(kv, expected) {
		t.Fatalf("Got: %s\nExpected: %s", kv, expected)
	}
}

func TestCollateKvPairsMapStringString(t *testing.T) {
	//init
	config := &struct {
		Vfoo   string
		Vother map[string]string
	}{
		Vfoo: "toto",
		Vother: map[string]string{
			"k1": "v1",
			"k2": "v2",
		},
	}

	//test
	kv := map[string]string{}
	if err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix"); err != nil {
		t.Fatalf("Error : %v", err)
	}

	//check
	expected := map[string]string{
		"prefix/vother/k1": "v1",
		"prefix/vother/k2": "v2",
		"prefix/vfoo":      "toto",
	}
	if !reflect.DeepEqual(kv, expected) {
		t.Fatalf("Got: %s\nExpected: %s", kv, expected)
	}
}

func TestCollateKvPairsMapIntString(t *testing.T) {
	//init
	config := &struct {
		Vfoo   string
		Vother map[int]string
	}{
		Vfoo: "toto",
		Vother: map[int]string{
			51: "v1",
			15: "v2",
		},
	}

	//test
	kv := map[string]string{}
	if err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix"); err != nil {
		t.Fatalf("Error : %v", err)
	}

	//check
	expected := map[string]string{
		"prefix/vother/51": "v1",
		"prefix/vother/15": "v2",
		"prefix/vfoo":      "toto",
	}
	if !reflect.DeepEqual(kv, expected) {
		t.Fatalf("Got: %s\nExpected: %s", kv, expected)
	}
}

func TestCollateKvPairsMapStringStruct(t *testing.T) {
	//init
	config := &struct {
		Vfoo   string
		Vother map[string]Struct1
	}{
		Vfoo: "toto",
		Vother: map[string]Struct1{
			"k1": {
				S1Bool:       true,
				S1Int:        51,
				S1PtrStruct3: nil,
			},
			"k2": {
				S1String: "tata",
			},
		},
	}

	//test
	kv := map[string]string{}
	if err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix"); err != nil {
		t.Fatalf("Error : %v", err)
	}

	//check
	expected := map[string]string{
		"prefix/vother/k1/s1bool":   "true",
		"prefix/vother/k1/s1int":    "51",
		"prefix/vother/k1/s1string": "",
		"prefix/vother/k2/s1bool":   "false",
		"prefix/vother/k2/s1int":    "0",
		"prefix/vother/k2/s1string": "tata",
		"prefix/vfoo":               "toto",
	}
	if !reflect.DeepEqual(kv, expected) {
		t.Fatalf("Got: %s\nExpected: %s", kv, expected)
	}
}

func TestCollateKvPairsMapStructStructSouldFail(t *testing.T) {
	//init
	config := &struct {
		Vfoo   string
		Vother map[Struct1]Struct1
	}{
		Vfoo: "toto",
		Vother: map[Struct1]Struct1{
			Struct1{
				S1Bool: true,
				S1Int:  1,
			}: {
				S1Int: 11,
			},
			Struct1{
				S1Bool: true,
				S1Int:  2,
			}: {
				S1Int: 22,
			},
		},
	}

	//test
	kv := map[string]string{}
	err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix")
	if err == nil || !strings.Contains(err.Error(), "struct as key not supported") {
		t.Fatalf("Expected error Struct as key not supported\ngot: %v", err)
	}
}

func TestCollateKvPairsSliceInt(t *testing.T) {
	//init
	config := &struct {
		Vfoo   string
		Vother []int
	}{
		Vfoo:   "toto",
		Vother: []int{51, 15},
	}

	//test
	kv := map[string]string{}
	if err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix"); err != nil {
		t.Fatalf("Error: %v", err)
	}

	//check
	expected := map[string]string{
		"prefix/vother/0": "51",
		"prefix/vother/1": "15",
		"prefix/vfoo":     "toto",
	}
	if !reflect.DeepEqual(kv, expected) {
		t.Fatalf("Got: %s\nExpected: %s", kv, expected)
	}
}

func TestCollateKvPairsSlicePtrOnStruct(t *testing.T) {
	//init
	config := &struct {
		Vfoo   string
		Vother []*BasicStruct
	}{
		Vfoo: "toto",
		Vother: []*BasicStruct{
			{},
			{
				Bar1: "tata",
				Bar2: "titi",
			},
		},
	}

	//test
	kv := map[string]string{}
	if err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix"); err != nil {
		t.Fatalf("Error: %v", err)
	}

	//check
	expected := map[string]string{
		"prefix/vother/0/":     "",
		"prefix/vother/0/bar1": "",
		"prefix/vother/0/bar2": "",
		"prefix/vother/1/":     "",
		"prefix/vother/1/bar1": "tata",
		"prefix/vother/1/bar2": "titi",
		"prefix/vfoo":          "toto",
	}
	if !reflect.DeepEqual(kv, expected) {
		t.Fatalf("Got: %s\nExpected: %s", kv, expected)
	}
}

func TestCollateKvPairsEmbedded(t *testing.T) {
	//init
	config := &struct {
		BasicStruct
		Vfoo string
	}{
		BasicStruct: BasicStruct{
			Bar1: "tata",
			Bar2: "titi",
		},
		Vfoo: "toto",
	}

	//test
	kv := map[string]string{}
	if err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix"); err != nil {
		t.Fatalf("Error: %v", err)
	}

	//check
	expected := map[string]string{
		"prefix/basicstruct/bar1": "tata",
		"prefix/basicstruct/bar2": "titi",
		"prefix/vfoo":             "toto",
	}
	if !reflect.DeepEqual(kv, expected) {
		t.Fatalf("Got: %s\nExpected: %s", kv, expected)
	}
}

func TestCollateKvPairsEmbeddedSquash(t *testing.T) {
	//init
	config := &struct {
		BasicStruct `mapstructure:",squash"`
		Vfoo        string
	}{
		BasicStruct: BasicStruct{
			Bar1: "tata",
			Bar2: "titi",
		},
		Vfoo: "toto",
	}

	//test
	kv := map[string]string{}
	if err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix"); err != nil {
		t.Fatalf("Error: %v", err)
	}

	//check
	expected := map[string]string{
		"prefix/bar1": "tata",
		"prefix/bar2": "titi",
		"prefix/vfoo": "toto",
	}
	if !reflect.DeepEqual(kv, expected) {
		t.Fatalf("Got: %s\nExpected: %s", kv, expected)
	}
}

func TestCollateKvPairsNotSupportedKindShouldFail(t *testing.T) {
	//init
	config := &struct {
		Vchan chan int
	}{
		Vchan: make(chan int),
	}

	//test
	kv := map[string]string{}
	err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix")
	if err == nil || !strings.Contains(err.Error(), "kind chan not supported") {
		t.Fatalf("Expected error : Kind chan not supported\nGot : %v", err)
	}
}

func TestStoreConfigEmbeddedSquash(t *testing.T) {
	//init
	config := &struct {
		BasicStruct `mapstructure:",squash"`
		Vfoo        string
	}{
		BasicStruct: BasicStruct{
			Bar1: "tata",
			Bar2: "titi",
		},
		Vfoo: "toto",
	}
	kv := &KvSource{
		&Mock{},
		"prefix",
	}
	//test
	if err := kv.StoreConfig(config); err != nil {
		t.Fatalf("Error: %v", err)
	}

	//check
	expected := map[string]string{
		"prefix/bar1": "tata",
		"prefix/bar2": "titi",
		"prefix/vfoo": "toto",
	}
	result := map[string][]byte{}
	err := kv.ListRecursive("prefix", result)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if len(result) != len(expected) {
		t.Fatalf("length of kv.List is not %d", len(expected))
	}
	for k, v := range result {
		if string(v) != expected[k] {
			t.Fatalf("Key %s\nExpected value %s, got %s", k, v, expected[k])
		}
	}
}

func TestCollateKvPairsUnexported(t *testing.T) {
	config := &struct {
		Vstring string
		vsilent string
	}{
		Vstring: "mustBeInTheMap",
		vsilent: "mustNotBeInTheMap",
	}

	//test
	kv := map[string]string{}
	if err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix"); err != nil {
		t.Fatalf("Error: %v", err)
	}

	//check
	if _, ok := kv["prefix/vsilent"]; ok {
		t.Fatalf("Exported field should not be in the map : %s", kv)
	}

	expected := map[string]string{
		"prefix/vstring": "mustBeInTheMap",
	}
	if !reflect.DeepEqual(kv, expected) {
		t.Fatalf("Got: %s\nExpected: %s", kv, expected)
	}

}

func TestCollateKvPairsShortNameUnexported(t *testing.T) {
	config := &struct {
		E string
		u string
	}{
		E: "mustBeInTheMap",
		u: "mustNotBeInTheMap",
	}

	//test
	kv := map[string]string{}
	if err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix"); err != nil {
		t.Fatalf("Error: %v", err)
	}

	//check
	if _, ok := kv["prefix/u"]; ok {
		t.Fatalf("Exported field should not be in the map : %s", kv)
	}

	expected := map[string]string{
		"prefix/e": "mustBeInTheMap",
	}
	if !reflect.DeepEqual(kv, expected) {
		t.Fatalf("Got: %s\nExpected: %s", kv, expected)
	}
}

func TestListRecursive5Levels(t *testing.T) {
	kv := &KvSource{
		&Mock{
			KVPairs: []*store.KVPair{
				{Key: "prefix/l1", Value: []byte("level1")},
				{Key: "prefix/d1/l1", Value: []byte("level2")},
				{Key: "prefix/d1/l2", Value: []byte("level2")},
				{Key: "prefix/d2/d1/l1", Value: []byte("level3")},
				{Key: "prefix/d3/d2/d1/d1/d1", Value: []byte("level5")},
			},
		},
		"prefix",
	}
	pairs := map[string][]byte{}
	err := kv.ListRecursive(kv.Prefix, pairs)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	//check
	expected := map[string][]byte{
		"prefix/l1":             []byte("level1"),
		"prefix/d1/l1":          []byte("level2"),
		"prefix/d1/l2":          []byte("level2"),
		"prefix/d2/d1/l1":       []byte("level3"),
		"prefix/d3/d2/d1/d1/d1": []byte("level5"),
	}
	if len(pairs) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(pairs))
	}
	for k, v := range pairs {
		if !reflect.DeepEqual(v, expected[k]) {
			t.Fatalf("Key %s\nExpected %s\nGot %s", k, expected[k], v)
		}
	}
}

func TestListRecursiveEmpty(t *testing.T) {
	kv := &KvSource{
		&Mock{
			KVPairs: []*store.KVPair{},
		},
		"prefix",
	}
	pairs := map[string][]byte{}
	err := kv.ListRecursive(kv.Prefix, pairs)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	//check
	expected := map[string][]byte{}
	if len(pairs) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(pairs))
	}
}

func TestConvertPairs5Levels(t *testing.T) {
	input := map[string][]byte{
		"prefix/l1":             []byte("level1"),
		"prefix/d1/l1":          []byte("level2"),
		"prefix/d1/l2":          []byte("level2"),
		"prefix/d2/d1/l1":       []byte("level3"),
		"prefix/d3/d2/d1/d1/d1": []byte("level5"),
	}
	//test
	output := convertPairs(input)

	//check
	expected := map[string][]byte{
		"prefix/l1":             []byte("level1"),
		"prefix/d1/l1":          []byte("level2"),
		"prefix/d1/l2":          []byte("level2"),
		"prefix/d2/d1/l1":       []byte("level3"),
		"prefix/d3/d2/d1/d1/d1": []byte("level5"),
	}

	if len(output) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(output))
	}
	for _, p := range output {
		if !reflect.DeepEqual(p.Value, expected[p.Key]) {
			t.Fatalf("Key : %s\nExpected %s\nGot %s", p.Key, expected[p.Key], p.Value)
		}
	}
}

func TestCollateKvPairsCompressedData(t *testing.T) {

	strToCompress := "Testing automatic compressed data if byte array"
	config := &struct {
		CompressedDataBytes []byte
	}{
		CompressedDataBytes: []byte(strToCompress),
	}

	//test
	kv := map[string]string{}
	if err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix"); err != nil {
		t.Fatalf("Error: %v", err)
	}

	compressedVal := kv["prefix/compresseddatabytes"]
	if len(compressedVal) == 0 {
		t.Fatal("Error : no entry for 'prefix/compresseddatabytes'.")
	}

	data, err := readCompressedData(compressedVal, gzipReader)
	dataStr := string(data)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if strToCompress != dataStr {
		t.Fatalf("Got: %q\nExpected: %q", dataStr, strToCompress)
	}
}

type TestCompressedDataStruct struct {
	CompressedDataBytes []byte
}

func TestParseKvSourceCompressedData(t *testing.T) {
	//Init
	config := TestCompressedDataStruct{}

	//Test
	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                &config,
		DefaultPointersConfig: &config,
		Run: func() error { return nil },
	}

	strToCompress := "Testing automatic compressed data if byte array"

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write([]byte(strToCompress))
	w.Close()

	kvs := []*KvSource{
		&KvSource{
			&Mock{
				KVPairs: []*store.KVPair{
					{
						Key:   "test/compresseddatabytes",
						Value: b.Bytes(),
					},
				},
			},
			"test",
		},
		&KvSource{
			&Mock{
				KVPairs: []*store.KVPair{
					{
						Key:   "test/compresseddatabytes",
						Value: []byte("VGVzdGluZyBhdXRvbWF0aWMgY29tcHJlc3NlZCBkYXRhIGlmIGJ5dGUgYXJyYXk="),
					},
				},
			},
			"test",
		},
	}

	for _, kv := range kvs {
		if _, err := kv.Parse(rootCmd); err != nil {
			t.Fatalf("Error %s", err)
		}

		//Check
		expected := &TestCompressedDataStruct{
			CompressedDataBytes: []byte(strToCompress),
		}

		if !reflect.DeepEqual(expected, rootCmd.Config) {
			actualJSON, err := json.Marshal(rootCmd.Config)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			expectedJSON, err := json.Marshal(expected)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			t.Fatalf("\nexpected\t: %s\ngot\t\t\t: %s\n", expectedJSON, actualJSON)
		}
	}
}

type CustomStruct struct {
	Bar1 string
	Bar2 string
}

// UnmarshalText define how unmarshal in TOML parsing
func (c *CustomStruct) UnmarshalText(text []byte) error {
	res := strings.Split(string(text), ",")
	c.Bar1 = res[0]
	c.Bar2 = res[1]
	return nil
}

// MarshalText encodes the receiver into UTF-8-encoded text and returns the result.
func (c *CustomStruct) MarshalText() (text []byte, err error) {
	return []byte(c.Bar1 + "," + c.Bar2), nil
}

func TestCollateCustomMarshaller(t *testing.T) {
	config := &CustomStruct{
		Bar1: "Bar1",
		Bar2: "Bar2",
	}
	//test
	kv := map[string]string{}
	if err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix"); err != nil {
		t.Fatalf("Error: %v", err)
	}
	//check

	check := map[string]string{
		"prefix": "Bar1,Bar2",
	}
	if !reflect.DeepEqual(kv, check) {
		t.Fatalf("Expected %s\nGot %s", check, kv)
	}
}

func TestDecodeHookCustomMarshaller(t *testing.T) {
	data := &CustomStruct{
		Bar1: "Bar1",
		Bar2: "Bar2",
	}
	output, err := decodeHook(reflect.TypeOf([]string{}), reflect.TypeOf(data), "Bar1,Bar2")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if !reflect.DeepEqual(data, output) {
		t.Fatalf("Got %#v\nExpected %#v", output, data)
	}
}
