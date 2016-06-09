package staert

import (
	"reflect"
	"testing"
)

func TestGenerateMapstructureBasic(t *testing.T) {
	input := map[string]string{
		"test/addr":       "foo",
		"test/child/data": "bar",
	}
	prefix := "test/"

	output, err := generateMapstructure(input, prefix)
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
	input := map[string]string{
		"test/vfoo":       "foo",
		"test/vother/foo": "foo",
		"test/vother/bar": "bar",
	}
	prefix := "test/"

	output, err := generateMapstructure(input, prefix)
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
	input := map[string]string{
		"test/vfoo":     "foo",
		"test/vother/0": "foo",
		"test/vother/1": "bar",
		"test/vother/2": "bar",
	}
	prefix := "test/"

	output, err := generateMapstructure(input, prefix)
	if err != nil {
		t.Fatalf("Error :%s", err)
	}
	//check
	check := map[string]interface{}{
		"vfoo": "foo",
		"vother": map[string]interface{}{
			"0": "foo",
			"1": "bar",
			"2": "bar",
		},
	}
	if !reflect.DeepEqual(check, output) {
		t.Fatalf("Expected %#v\nGot %#v", check, output)
	}
}

func TestGenerateMapstructureNotTrivialSlice(t *testing.T) {
	input := map[string]string{
		"test/vfoo":          "foo",
		"test/vother/0/foo1": "bar",
		"test/vother/0/foo2": "bar",
		"test/vother/1/bar1": "foo",
		"test/vother/1/bar2": "foo",
	}
	prefix := "test/"

	output, err := generateMapstructure(input, prefix)
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
