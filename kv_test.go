package staert

import (
	"encoding/json"
	"errors"
	"github.com/containous/flaeg"
	"github.com/docker/libkv/store"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestGenerateMapstructureBasic(t *testing.T) {
	moke := []*store.KVPair{
		&store.KVPair{
			Key:   "test/addr",
			Value: []byte("foo"),
		},
		&store.KVPair{
			Key:   "test/child/data",
			Value: []byte("bar"),
		},
	}
	prefix := "test/"

	output, err := generateMapstructure(moke, prefix)
	if err != nil {
		t.Fatalf("Error :%s", err)
	}
	//check
	check := map[string]interface{}{
		"addr": "foo",
		"child": map[string]interface{}{
			"data": "bar",
		},
	}
	if !reflect.DeepEqual(check, output) {
		t.Fatalf("Expected %+v\nGot %+v", check, output)
	}
}

func TestGenerateMapstructureTrivialMap(t *testing.T) {
	moke := []*store.KVPair{
		&store.KVPair{
			Key:   "test/vfoo",
			Value: []byte("foo"),
		},
		&store.KVPair{
			Key:   "test/vother/foo",
			Value: []byte("foo"),
		},
		&store.KVPair{
			Key:   "test/vother/bar",
			Value: []byte("bar"),
		},
	}
	prefix := "test/"

	output, err := generateMapstructure(moke, prefix)
	if err != nil {
		t.Fatalf("Error :%s", err)
	}
	//check
	check := map[string]interface{}{
		"vfoo": "foo",
		"vother": map[string]interface{}{
			"foo": "foo",
			"bar": "bar",
		},
	}
	if !reflect.DeepEqual(check, output) {
		t.Fatalf("Expected %#v\nGot %#v", check, output)
	}
}

func TestGenerateMapstructureTrivialSlice(t *testing.T) {
	moke := []*store.KVPair{
		&store.KVPair{
			Key:   "test/vfoo",
			Value: []byte("foo"),
		},
		&store.KVPair{
			Key:   "test/vother/0",
			Value: []byte("foo"),
		},
		&store.KVPair{
			Key:   "test/vother/1",
			Value: []byte("bar1"),
		},
		&store.KVPair{
			Key:   "test/vother/2",
			Value: []byte("bar2"),
		},
	}
	prefix := "test/"

	output, err := generateMapstructure(moke, prefix)
	if err != nil {
		t.Fatalf("Error :%s", err)
	}
	//check
	check := map[string]interface{}{
		"vfoo": "foo",
		"vother": map[string]interface{}{
			"0": "foo",
			"1": "bar1",
			"2": "bar2",
		},
	}
	if !reflect.DeepEqual(check, output) {
		t.Fatalf("Expected %#v\nGot %#v", check, output)
	}
}

func TestGenerateMapstructureNotTrivialSlice(t *testing.T) {
	moke := []*store.KVPair{
		&store.KVPair{
			Key:   "test/vfoo",
			Value: []byte("foo"),
		},
		&store.KVPair{
			Key:   "test/vother/0/foo1",
			Value: []byte("bar"),
		},
		&store.KVPair{
			Key:   "test/vother/0/foo2",
			Value: []byte("bar"),
		},
		&store.KVPair{
			Key:   "test/vother/1/bar1",
			Value: []byte("foo"),
		},
		&store.KVPair{
			Key:   "test/vother/1/bar2",
			Value: []byte("foo"),
		},
	}
	prefix := "test/"

	output, err := generateMapstructure(moke, prefix)
	if err != nil {
		t.Fatalf("Error :%s", err)
	}
	//check
	check := map[string]interface{}{
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
	if !reflect.DeepEqual(check, output) {
		t.Fatalf("Expected %#v\nGot %#v", check, output)
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

	check := []interface{}{
		map[string]interface{}{
			"bar1": "foo1",
			"bar2": "foo2",
		},
		map[string]interface{}{
			"bar1": "bar1",
			"bar2": "bar2",
		},
	}
	if !reflect.DeepEqual(check, output) {
		t.Fatalf("Expected %#v\nGot %#v", check, output)
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
		t.Fatalf("got an err: %s", err.Error())
	}
	if err := decoder.Decode(input); err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	//check
	check := Test{
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

	if !reflect.DeepEqual(check, config) {
		t.Fatalf("Expected %#v\nGot %#v", check, config)
	}
}

// Extremely limited mock store so we can test initialization
type Mock struct {
	Error           bool
	KVPairs         []*store.KVPair
	WatchTreeMethod func() <-chan []*store.KVPair
}

func (s *Mock) Put(key string, value []byte, opts *store.WriteOptions) error {
	return errors.New("Put not supported")
}

func (s *Mock) Get(key string) (*store.KVPair, error) {
	if s.Error {
		return nil, errors.New("Error")
	}
	for _, kvPair := range s.KVPairs {
		if kvPair.Key == key {
			return kvPair, nil
		}
	}
	return nil, nil
}

func (s *Mock) Delete(key string) error {
	return errors.New("Delete not supported")
}

// Exists mock
func (s *Mock) Exists(key string) (bool, error) {
	return false, errors.New("Exists not supported")
}

// Watch mock
func (s *Mock) Watch(key string, stopCh <-chan struct{}) (<-chan *store.KVPair, error) {
	return nil, errors.New("Watch not supported")
}

// WatchTree mock
func (s *Mock) WatchTree(prefix string, stopCh <-chan struct{}) (<-chan []*store.KVPair, error) {
	return s.WatchTreeMethod(), nil
}

// NewLock mock
func (s *Mock) NewLock(key string, options *store.LockOptions) (store.Locker, error) {
	return nil, errors.New("NewLock not supported")
}

// List mock
func (s *Mock) List(prefix string) ([]*store.KVPair, error) {
	if s.Error {
		return nil, errors.New("Error")
	}
	kv := []*store.KVPair{}
	for _, kvPair := range s.KVPairs {
		if strings.HasPrefix(kvPair.Key, prefix) {
			kv = append(kv, kvPair)
		}
	}
	return kv, nil
}

// DeleteTree mock
func (s *Mock) DeleteTree(prefix string) error {
	return errors.New("DeleteTree not supported")
}

// AtomicPut mock
func (s *Mock) AtomicPut(key string, value []byte, previous *store.KVPair, opts *store.WriteOptions) (bool, *store.KVPair, error) {
	return false, nil, errors.New("AtomicPut not supported")
}

// AtomicDelete mock
func (s *Mock) AtomicDelete(key string, previous *store.KVPair) (bool, error) {
	return false, errors.New("AtomicDelete not supported")
}

// Close mock
func (s *Mock) Close() {
	return
}

func TestKvSourceEmpty(t *testing.T) {
	//Init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
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

	//Check
	check := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
	}

	if !reflect.DeepEqual(check, rootCmd.Config) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, rootCmd.Config)
	}
}
func TestGenerateMapstructureTrivial(t *testing.T) {
	input := []*store.KVPair{
		{
			Key:   "test/ptrstruct1/s1int",
			Value: []byte("28"),
		},
		{
			Key:   "test/durationfield",
			Value: []byte("28"),
		},
	}
	prefix := "test/"
	output, err := generateMapstructure(input, prefix)
	if err != nil {
		t.Fatalf("Error :%s", err)
	}
	//check
	check := map[string]interface{}{
		"durationfield": "28",
		"ptrstruct1": map[string]interface{}{
			"s1int": "28",
		},
	}
	if !reflect.DeepEqual(check, output) {
		printResult, err := json.Marshal(output)
		if err != nil {
			t.Fatalf("error: %s", err)
		}
		printCheck, err := json.Marshal(check)
		if err != nil {
			t.Fatalf("error: %s", err)
		}
		t.Fatalf("\nexpected\t: %s\ngot\t\t\t: %s\n", printCheck, printResult)
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
		t.Fatalf("got an err: %s", err.Error())
	}
	if err := decoder.Decode(mapstruct); err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	//check
	check := StructPtr{
		PtrStruct1: &Struct1{
			S1Int: 28,
		},
		DurationField: time.Nanosecond * 28,
	}

	if !reflect.DeepEqual(check, config) {
		printResult, err := json.Marshal(config)
		if err != nil {
			t.Fatalf("error: %s", err)
		}
		printCheck, err := json.Marshal(check)
		if err != nil {
			t.Fatalf("error: %s", err)
		}
		t.Fatalf("\nexpected\t: %s\ngot\t\t\t: %s\n", printCheck, printResult)
	}
}
func TestIntegrationMapstructureInitedPtrReset(t *testing.T) {
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
		DurationField: time.Nanosecond * 28,
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
		t.Fatalf("got an err: %s", err.Error())
	}
	if err := decoder.Decode(mapstruct); err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	//check
	check := StructPtr{
		PtrStruct1: &Struct1{
			S1Int: 24,
		},
		DurationField: time.Nanosecond * 28,
	}

	if !reflect.DeepEqual(check, config) {
		printResult, err := json.Marshal(config)
		if err != nil {
			t.Fatalf("error: %s", err)
		}
		printCheck, err := json.Marshal(check)
		if err != nil {
			t.Fatalf("error: %s", err)
		}
		t.Fatalf("\nexpected\t: %s\ngot\t\t\t: %s\n", printCheck, printResult)
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
				{
					Key:   "test/ptrstruct1/s1int",
					Value: []byte("28"),
				},
				{
					Key:   "test/durationfield",
					Value: []byte("28"),
				},
			},
		},
		"test/",
	}
	if _, err := kv.Parse(rootCmd); err != nil {
		t.Fatalf("Error %s", err)
	}

	//Check
	check := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int: 28,
		},
		DurationField: time.Nanosecond * 28,
	}

	if !reflect.DeepEqual(check, rootCmd.Config) {
		printResult, err := json.Marshal(rootCmd.Config)
		if err != nil {
			t.Fatalf("error: %s", err)
		}
		printCheck, err := json.Marshal(check)
		if err != nil {
			t.Fatalf("error: %s", err)
		}
		t.Fatalf("\nexpected\t: %s\ngot\t\t\t: %s\n", printCheck, printResult)
	}
}
func TestIntegrationMockList(t *testing.T) {
	kv := &Mock{
		KVPairs: []*store.KVPair{
			{
				Key:   "test/ptrstruct1/s1int",
				Value: []byte("28"),
			},
			{
				Key:   "test/durationfield",
				Value: []byte("28"),
			},
		},
	}
	pairs, err := kv.List("test/")
	if err != nil {
		t.Fatalf("Error : %s", err)
	}
	//check
	if len(pairs) != 2 {
		t.Fatalf("Expected 2 pairs got %d", len(pairs))
	}
	check := map[string][]byte{
		"test/ptrstruct1/s1int": []byte("28"),
		"test/durationfield":    []byte("28"),
	}
	for _, p := range pairs {
		if !reflect.DeepEqual(p.Value, check[p.Key]) {
			t.Fatalf("key %s expected value %s got %s", p.Key, check[p.Key], p.Value)
		}
	}

}
