package staert

import (
	"github.com/docker/libkv/store"
	"reflect"
	"testing"
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
	type BasicStruct struct {
		Bar1 string
		Bar2 string
	}
	type SliceStruct []BasicStruct
	type Test struct {
		Vfoo   string
		Vother SliceStruct
	}
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
	output, err := decodeHookSlice(reflect.TypeOf(data), reflect.TypeOf([]string{}), data)
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
