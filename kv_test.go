package staert

import (
	"bytes"
	"compress/gzip"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/abronan/valkeyrie/store"
	"github.com/containous/flaeg"
	"github.com/containous/flaeg/parse"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateMapstructureBasic(t *testing.T) {
	mock := []*store.KVPair{
		{Key: "test/addr", Value: []byte("foo")},
		{Key: "test/child/data", Value: []byte("bar")}}
	prefix := "test"

	output, err := generateMapstructure(mock, prefix)
	require.NoError(t, err)

	expected := map[string]interface{}{
		"addr": "foo",
		"child": map[string]interface{}{
			"data": "bar",
		},
	}

	assert.Exactly(t, expected, output)
}

func TestGenerateMapstructureTrivialMap(t *testing.T) {
	mock := []*store.KVPair{
		{Key: "test/vfoo", Value: []byte("foo")},
		{Key: "test/vother/foo", Value: []byte("foo")},
		{Key: "test/vother/bar", Value: []byte("bar")},
	}
	prefix := "test"

	output, err := generateMapstructure(mock, prefix)
	require.NoError(t, err)

	expected := map[string]interface{}{
		"vfoo": "foo",
		"vother": map[string]interface{}{
			"foo": "foo",
			"bar": "bar",
		},
	}

	assert.Exactly(t, expected, output)
}

func TestGenerateMapstructureTrivialSlice(t *testing.T) {
	mock := []*store.KVPair{
		{Key: "test/vfoo", Value: []byte("foo")},
		{Key: "test/vother/0", Value: []byte("foo")},
		{Key: "test/vother/1", Value: []byte("bar1")},
		{Key: "test/vother/2", Value: []byte("bar2")},
	}
	prefix := "test"

	output, err := generateMapstructure(mock, prefix)
	require.NoError(t, err)

	expected := map[string]interface{}{
		"vfoo": "foo",
		"vother": map[string]interface{}{
			"0": "foo",
			"1": "bar1",
			"2": "bar2",
		},
	}

	assert.Exactly(t, expected, output)
}

func TestGenerateMapstructureNotTrivialSlice(t *testing.T) {
	mock := []*store.KVPair{
		{Key: "test/vfoo", Value: []byte("foo")},
		{Key: "test/vother/0/foo1", Value: []byte("bar")},
		{Key: "test/vother/0/foo2", Value: []byte("bar")},
		{Key: "test/vother/1/bar1", Value: []byte("foo")},
		{Key: "test/vother/1/bar2", Value: []byte("foo")},
	}
	prefix := "test"

	output, err := generateMapstructure(mock, prefix)
	require.NoError(t, err)

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

	assert.Exactly(t, expected, output)
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
	require.NoError(t, err)

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

	assert.Exactly(t, expected, output)
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
	require.NoError(t, err)

	err = decoder.Decode(input)
	require.NoError(t, err)

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

	assert.Exactly(t, expected, config)
}

func TestKvSourceEmpty(t *testing.T) {
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: config,
		Run:                   func() error { return nil },
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
	require.NoError(t, err)

	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	assert.Exactly(t, expected, rootCmd.Config)
}

func TestGenerateMapstructureTrivial(t *testing.T) {
	input := []*store.KVPair{
		{Key: "test/ptrstruct1/s1int", Value: []byte("28")},
		{Key: "test/durationfield", Value: []byte("28")},
	}

	prefix := "test"

	output, err := generateMapstructure(input, prefix)
	require.NoError(t, err)

	expected := map[string]interface{}{
		"durationfield": "28",
		"ptrstruct1": map[string]interface{}{
			"s1int": "28",
		},
	}

	assert.Exactly(t, expected, output)
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
	require.NoError(t, err)

	err = decoder.Decode(mapstruct)
	require.NoError(t, err)

	expected := StructPtr{
		PtrStruct1: &Struct1{
			S1Int: 28,
		},
		DurationField: parse.Duration(28 * time.Nanosecond),
	}

	assert.Exactly(t, expected, config)
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
		DurationField: parse.Duration(28 * time.Second),
	}

	configDecoder := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           &config,
		WeaklyTypedInput: true,
		DecodeHook:       decodeHook,
	}
	decoder, err := mapstructure.NewDecoder(configDecoder)
	require.NoError(t, err)

	err = decoder.Decode(mapstruct)
	require.NoError(t, err)

	expected := StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    24,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(28 * time.Second),
	}

	assert.Exactly(t, expected, config)
}

func TestParseKvSourceTrivial(t *testing.T) {
	config := StructPtr{}

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                &config,
		DefaultPointersConfig: &config,
		Run:                   func() error { return nil },
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

	_, err := kv.Parse(rootCmd)
	require.NoError(t, err)

	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int: 28,
		},
		DurationField: parse.Duration(28 * time.Nanosecond),
	}

	assert.Exactly(t, expected, rootCmd.Config)
}

func TestLoadConfigKvSourceNestedPtrsNil(t *testing.T) {
	config := &StructPtr{}

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

	err := kv.LoadConfig(config)
	require.NoError(t, err)

	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:        1,
			S1String:     "S1StringInitConfig",
			S1PtrStruct3: &Struct3{},
		},
		DurationField: parse.Duration(21 * time.Second),
	}

	assert.Exactly(t, expected, config)
}

func TestParseKvSourceNestedPtrsNil(t *testing.T) {
	config := StructPtr{}

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                &config,
		DefaultPointersConfig: &config,
		Run:                   func() error { return nil },
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

	_, err := kv.Parse(rootCmd)
	require.NoError(t, err)

	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:        1,
			S1String:     "S1StringInitConfig",
			S1PtrStruct3: &Struct3{},
		},
		DurationField: parse.Duration(21 * time.Second),
	}

	assert.Exactly(t, expected, rootCmd.Config)
}

func TestParseKvSourceMap(t *testing.T) {
	config := &struct {
		Vmap map[string]int
	}{}

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: config,
		Run:                   func() error { return nil },
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

	_, err := kv.Parse(rootCmd)
	require.NoError(t, err)

	expected := &struct {
		Vmap map[string]int
	}{
		Vmap: map[string]int{
			"toto": 1,
			"tata": 2,
			"titi": 3,
		},
	}

	assert.Exactly(t, expected, rootCmd.Config)
}

func TestCollateKvPairsBasic(t *testing.T) {
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

	kv := map[string]string{}
	err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix")
	require.NoError(t, err)

	expected := map[string]string{
		"prefix/vbool":   "true",
		"prefix/vfloat":  "1.5",
		"prefix/vextra":  "toto",
		"prefix/vdata":   "42",
		"prefix/vstring": "tata",
		"prefix/vint":    "-15",
		"prefix/vuint":   "51",
	}

	assert.Exactly(t, expected, kv)
}

func TestCollateKvPairsNestedPointers(t *testing.T) {
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:        1,
			S1String:     "S1StringInitConfig",
			S1PtrStruct3: &Struct3{},
		},
		DurationField: parse.Duration(21 * time.Second),
	}

	kv := map[string]string{}
	err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix")
	require.NoError(t, err)

	expected := map[string]string{
		"prefix/ptrstruct1/s1int":                  "1",
		"prefix/ptrstruct1/s1string":               "S1StringInitConfig",
		"prefix/ptrstruct1/s1bool":                 "false",
		"prefix/ptrstruct1/s1ptrstruct3/":          "",
		"prefix/ptrstruct1/s1ptrstruct3/s3float64": "0",
		"prefix/durationfield":                     "21000000000",
	}

	assert.Exactly(t, expected, kv)
}

func TestCollateKvPairsMapStringString(t *testing.T) {
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

	kv := map[string]string{}
	err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix")
	require.NoError(t, err)

	expected := map[string]string{
		"prefix/vother/k1": "v1",
		"prefix/vother/k2": "v2",
		"prefix/vfoo":      "toto",
	}

	assert.Exactly(t, expected, kv)
}

func TestCollateKvPairsMapIntString(t *testing.T) {
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

	kv := map[string]string{}
	err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix")
	require.NoError(t, err)

	expected := map[string]string{
		"prefix/vother/51": "v1",
		"prefix/vother/15": "v2",
		"prefix/vfoo":      "toto",
	}

	assert.Exactly(t, expected, kv)
}

func TestCollateKvPairsMapStringStruct(t *testing.T) {
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

	kv := map[string]string{}
	err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix")
	require.NoError(t, err)

	expected := map[string]string{
		"prefix/vother/k1/s1bool":   "true",
		"prefix/vother/k1/s1int":    "51",
		"prefix/vother/k1/s1string": "",
		"prefix/vother/k2/s1bool":   "false",
		"prefix/vother/k2/s1int":    "0",
		"prefix/vother/k2/s1string": "tata",
		"prefix/vfoo":               "toto",
	}

	assert.Exactly(t, expected, kv)
}

func TestCollateKvPairsMapStructStructSouldFail(t *testing.T) {
	config := &struct {
		Vfoo   string
		Vother map[Struct1]Struct1
	}{
		Vfoo: "toto",
		Vother: map[Struct1]Struct1{
			{
				S1Bool: true,
				S1Int:  1,
			}: {
				S1Int: 11,
			},
			{
				S1Bool: true,
				S1Int:  2,
			}: {
				S1Int: 22,
			},
		},
	}

	kv := map[string]string{}

	err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix")
	require.Error(t, err)
	require.Contains(t, err.Error(), "struct as key not supported")
}

func TestCollateKvPairsSliceInt(t *testing.T) {
	config := &struct {
		Vfoo   string
		Vother []int
	}{
		Vfoo:   "toto",
		Vother: []int{51, 15},
	}

	kv := map[string]string{}
	err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix")
	require.NoError(t, err)

	expected := map[string]string{
		"prefix/vother/0": "51",
		"prefix/vother/1": "15",
		"prefix/vfoo":     "toto",
	}

	assert.Exactly(t, expected, kv)
}

func TestCollateKvPairsSlicePtrOnStruct(t *testing.T) {
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

	kv := map[string]string{}
	err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix")
	require.NoError(t, err)

	expected := map[string]string{
		"prefix/vother/0/":     "",
		"prefix/vother/0/bar1": "",
		"prefix/vother/0/bar2": "",
		"prefix/vother/1/":     "",
		"prefix/vother/1/bar1": "tata",
		"prefix/vother/1/bar2": "titi",
		"prefix/vfoo":          "toto",
	}

	assert.Exactly(t, expected, kv)
}

func TestCollateKvPairsEmbedded(t *testing.T) {
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

	kv := map[string]string{}
	err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix")
	require.NoError(t, err)

	expected := map[string]string{
		"prefix/basicstruct/bar1": "tata",
		"prefix/basicstruct/bar2": "titi",
		"prefix/vfoo":             "toto",
	}

	assert.Exactly(t, expected, kv)
}

func TestCollateKvPairsEmbeddedSquash(t *testing.T) {
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

	kv := map[string]string{}
	err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix")
	require.NoError(t, err)

	expected := map[string]string{
		"prefix/bar1": "tata",
		"prefix/bar2": "titi",
		"prefix/vfoo": "toto",
	}

	assert.Exactly(t, expected, kv)
}

func TestCollateKvPairsNotSupportedKindShouldFail(t *testing.T) {
	config := &struct {
		Vchan chan int
	}{
		Vchan: make(chan int),
	}

	kv := map[string]string{}

	err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix")
	require.Error(t, err)
	require.Contains(t, err.Error(), "kind chan not supported")
}

func TestStoreConfigEmbeddedSquash(t *testing.T) {
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

	err := kv.StoreConfig(config)
	require.NoError(t, err)

	expected := map[string]string{
		"prefix/bar1": "tata",
		"prefix/bar2": "titi",
		"prefix/vfoo": "toto",
	}

	result, err := kv.ListValuedPairWithPrefix("prefix")
	require.NoError(t, err)

	assert.Len(t, result, len(expected))
	for k, v := range result {
		assert.Equal(t, expected[k], string(v))
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

	kv := map[string]string{}
	err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix")
	require.NoError(t, err)

	if _, ok := kv["prefix/vsilent"]; ok {
		t.Fatalf("Exported field should not be in the map : %s", kv)
	}

	expected := map[string]string{
		"prefix/vstring": "mustBeInTheMap",
	}

	assert.Exactly(t, expected, kv)
}

func TestCollateKvPairsShortNameUnexported(t *testing.T) {
	config := &struct {
		E string
		u string
	}{
		E: "mustBeInTheMap",
		u: "mustNotBeInTheMap",
	}

	kv := map[string]string{}
	err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix")
	require.NoError(t, err)

	if _, ok := kv["prefix/u"]; ok {
		t.Fatalf("Exported field should not be in the map : %s", kv)
	}

	expected := map[string]string{
		"prefix/e": "mustBeInTheMap",
	}

	assert.Exactly(t, expected, kv)
}

func TestListRecursive(t *testing.T) {
	testCases := []struct {
		desc     string
		store    store.Store
		expected map[string][]byte
	}{
		{
			desc: "With unknown prefix",
			store: &Mock{
				Error: false,
				KVPairs: []*store.KVPair{
					{Key: "prefix/l1", Value: []byte("")},
				},
				WatchTreeMethod: nil,
				ListError:       store.ErrKeyNotFound,
			},
			expected: map[string][]byte{},
		},
		{
			desc: "recursive 5 levels",
			store: &Mock{
				KVPairs: []*store.KVPair{
					{Key: "prefix/l1", Value: []byte("level1")},
					{Key: "prefix/d1/l1", Value: []byte("level2")},
					{Key: "prefix/d1/l2", Value: []byte("level2")},
					{Key: "prefix/d2/d1/l1", Value: []byte("level3")},
					{Key: "prefix/d3/d2/d1/d1/d1", Value: []byte("level5")},
				},
			},
			expected: map[string][]byte{
				"prefix/l1":             []byte("level1"),
				"prefix/d1/l1":          []byte("level2"),
				"prefix/d1/l2":          []byte("level2"),
				"prefix/d2/d1/l1":       []byte("level3"),
				"prefix/d3/d2/d1/d1/d1": []byte("level5"),
			},
		},
		{
			desc: "recursive empty",
			store: &Mock{
				KVPairs: []*store.KVPair{},
			},
			expected: map[string][]byte{},
		},
		{
			desc: "same prefix",
			store: &Mock{
				KVPairs: []*store.KVPair{
					{Key: "prefix", Value: []byte("")},
					{Key: "prefix/tls", Value: []byte("")},
					{Key: "prefix/tls/ca", Value: []byte("v1")},
					{Key: "prefix/tls/caOptional", Value: []byte("v2")},
					{Key: "otherprefix/tls/ca", Value: []byte("other")},
				},
			},
			expected: map[string][]byte{
				"prefix/tls/caOptional": []byte("v2"),
				// missing the "prefix/tls/ca"
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			kv := &KvSource{
				Store:  test.store,
				Prefix: "prefix",
			}

			pairs := map[string][]byte{}
			err := kv.ListRecursive(kv.Prefix, pairs)
			require.NoError(t, err)

			assert.Equal(t, test.expected, pairs)
		})
	}
}

func TestListRecursive_Error(t *testing.T) {
	testCases := []struct {
		desc     string
		store    store.Store
		expected string
	}{
		{
			desc: "get error",
			store: &KvSource{
				Store: &Mock{
					GetError: errors.New("a GET error"),
					KVPairs: []*store.KVPair{
						{Key: "prefix/l1", Value: []byte("")},
					},
					WatchTreeMethod: nil,
				},
			},
			expected: "a GET error",
		},
		{
			desc: "list Error",
			store: &Mock{
				ListError: errors.New("another error"),
				KVPairs: []*store.KVPair{
					{Key: "prefix/l1", Value: []byte("")},
				},
				WatchTreeMethod: nil,
			},
			expected: "another error",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			kv := &KvSource{Store: test.store, Prefix: "prefix"}

			pairs := map[string][]byte{}
			err := kv.ListRecursive(kv.Prefix, pairs)
			assert.EqualError(t, err, test.expected)
			assert.Len(t, pairs, 0)
		})
	}
}

func TestListValuedPairWithPrefix(t *testing.T) {
	testCases := []struct {
		desc     string
		store    store.Store
		expected map[string][]byte
	}{
		{
			desc: "with unknown prefix",
			store: &Mock{
				Error: false,
				KVPairs: []*store.KVPair{
					{Key: "prefix/l1", Value: []byte("")},
				},
				WatchTreeMethod: nil,
				ListError:       store.ErrKeyNotFound,
			},
			expected: map[string][]byte{},
		},
		{
			desc: "with prefix 5 levels",
			store: &Mock{
				KVPairs: []*store.KVPair{
					{Key: "prefix/", Value: []byte("")},
					{Key: "prefix/l1", Value: []byte("level1")},
					{Key: "prefix/d1/l1", Value: []byte("level2")},
					{Key: "prefix/d1/l2", Value: []byte("level2")},
					{Key: "prefix/d2/d1/l1", Value: []byte("level3")},
					{Key: "prefix/d3/d2/d1/d1/d1", Value: []byte("level5")},
				},
			},
			expected: map[string][]byte{
				"prefix/l1":             []byte("level1"),
				"prefix/d1/l1":          []byte("level2"),
				"prefix/d1/l2":          []byte("level2"),
				"prefix/d2/d1/l1":       []byte("level3"),
				"prefix/d3/d2/d1/d1/d1": []byte("level5"),
			},
		},
		{
			desc: "with prefix same prefix",
			store: &Mock{
				KVPairs: []*store.KVPair{
					{Key: "prefix", Value: []byte("")},
					{Key: "prefix/tls", Value: []byte("")},
					{Key: "prefix/tls/ca", Value: []byte("v1")},
					{Key: "prefix/tls/caOptional", Value: []byte("v2")},
					{Key: "otherprefix/tls/ca", Value: []byte("other")},
				},
			},
			expected: map[string][]byte{
				"prefix/tls/ca":         []byte("v1"),
				"prefix/tls/caOptional": []byte("v2"),
			},
		},
		{
			desc: "TestFetchValuedPairWithPrefixEmpty",
			store: &Mock{
				KVPairs: []*store.KVPair{},
			},
			expected: map[string][]byte{},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			kv := &KvSource{
				Store:  test.store,
				Prefix: "prefix",
			}

			pairs, err := kv.ListValuedPairWithPrefix(kv.Prefix)
			require.NoError(t, err)

			assert.Equal(t, test.expected, pairs)
		})
	}
}

func TestListValuedPairWithPrefix_Error(t *testing.T) {
	kv := &KvSource{
		&Mock{
			ListError: errors.New("another error"),
			KVPairs: []*store.KVPair{
				{Key: "prefix/l1", Value: []byte("")},
			},
			WatchTreeMethod: nil,
		},
		"prefix",
	}

	pairs, err := kv.ListValuedPairWithPrefix(kv.Prefix)
	assert.NotNil(t, err)
	assert.Len(t, pairs, 0)
}

func TestConvertPairs5Levels(t *testing.T) {
	input := map[string][]byte{
		"prefix/l1":             []byte("level1"),
		"prefix/d1/l1":          []byte("level2"),
		"prefix/d1/l2":          []byte("level2"),
		"prefix/d2/d1/l1":       []byte("level3"),
		"prefix/d3/d2/d1/d1/d1": []byte("level5"),
	}
	output := convertPairs(input)

	expected := map[string][]byte{
		"prefix/l1":             []byte("level1"),
		"prefix/d1/l1":          []byte("level2"),
		"prefix/d1/l2":          []byte("level2"),
		"prefix/d2/d1/l1":       []byte("level3"),
		"prefix/d3/d2/d1/d1/d1": []byte("level5"),
	}

	assert.Len(t, output, len(expected))

	for _, p := range output {
		assert.Exactlyf(t, expected[p.Key], p.Value, p.Key)
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
	err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix")
	require.NoError(t, err)

	compressedVal := kv["prefix/compresseddatabytes"]
	if len(compressedVal) == 0 {
		t.Fatal("Error : no entry for 'prefix/compresseddatabytes'.")
	}

	data, err := readCompressedData(compressedVal, gzipReader)
	require.NoError(t, err)

	assert.Exactly(t, strToCompress, string(data))
}

type TestCompressedDataStruct struct {
	CompressedDataBytes []byte
}

func TestParseKvSourceCompressedData(t *testing.T) {
	config := TestCompressedDataStruct{}

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                &config,
		DefaultPointersConfig: &config,
		Run:                   func() error { return nil },
	}

	strToCompress := "Testing automatic compressed data if byte array"

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write([]byte(strToCompress))
	require.NoError(t, err)

	err = w.Close()
	require.NoError(t, err)

	kvs := []*KvSource{
		{
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
		{
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
		_, err := kv.Parse(rootCmd)
		require.NoError(t, err)

		expected := &TestCompressedDataStruct{
			CompressedDataBytes: []byte(strToCompress),
		}

		assert.Exactly(t, expected, rootCmd.Config)
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
	err := collateKvRecursive(reflect.ValueOf(config), kv, "prefix")
	require.NoError(t, err)

	expected := map[string]string{
		"prefix": "Bar1,Bar2",
	}

	assert.Exactly(t, expected, kv)
}

func TestDecodeHookCustomMarshaller(t *testing.T) {
	data := &CustomStruct{
		Bar1: "Bar1",
		Bar2: "Bar2",
	}

	output, err := decodeHook(reflect.TypeOf([]string{}), reflect.TypeOf(data), "Bar1,Bar2")
	require.NoError(t, err)

	assert.Exactly(t, data, output)
}
