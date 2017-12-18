package staert

import (
	"reflect"
	"testing"

	"github.com/containous/flaeg/parse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getPtrPtrConfig() *struct {
	StringValue string
	NextPointer **basicAppConfig
} {
	configPtr := &basicAppConfig{
		BoolValue:   true,
		IntValue:    1,
		StringValue: "string",
	}
	expectedPtrPtr := &struct {
		StringValue string
		NextPointer **basicAppConfig
	}{
		StringValue: "FOO",
	}
	expectedPtrPtr.NextPointer = &configPtr
	return expectedPtrPtr
}

func TestFilterEnvVarWithPrefix(t *testing.T) {
	envSource := []*envValue{
		{strValue: "FOOO", path: []string{"Config", "0", "foo", "StringValue"}},
		{strValue: "10", path: []string{"Config", "0", "foo", "IntValue"}},
		{strValue: "10", path: []string{"Config", "IntValue"}},
		{strValue: "10", path: []string{"Config", "0", "0", "IntValue"}},
	}

	result := filterEnvVarWithPrefix(envSource, []string{"Config", "0", "foo"})

	expected := []*envValue{
		{strValue: "FOOO", path: []string{"StringValue"}},
		{strValue: "10", path: []string{"IntValue"}},
	}

	assert.Exactly(t, expected, result)
}

func TestAssignValues(t *testing.T) {
	parsers, _ := parse.LoadParsers(nil)
	subject := &EnvSource{prefix: "", separator: "_", parsers: parsers}

	testCases := []struct {
		desc     string
		source   interface{}
		values   []*envValue
		expected interface{}
	}{
		{
			desc: "with a simple struct",
			source: &struct {
				StringValue      string
				OtherStringValue string
			}{},
			values: []*envValue{
				{strValue: "FOO", path: []string{"StringValue"}},
				{strValue: "BAR", path: []string{"OtherStringValue"}},
			},
			expected: &struct {
				StringValue      string
				OtherStringValue string
			}{"FOO", "BAR"},
		},
		{
			desc: "when it needs an int parser",
			source: &struct {
				StringValue string
				IntValue    int
			}{},
			values: []*envValue{
				{strValue: "FOO", path: []string{"StringValue"}},
				{strValue: "1", path: []string{"IntValue"}},
			},
			expected: &struct {
				StringValue string
				IntValue    int
			}{"FOO", 1},
		},
		{
			desc: "with an embedded struct",
			source: &struct {
				StringValue string
				Next        basicAppConfig
			}{},
			values: []*envValue{
				{strValue: "FOO", path: []string{"StringValue"}},
				{strValue: "1", path: []string{"Next", "IntValue"}},
				{strValue: "true", path: []string{"Next", "BoolValue"}},
				{strValue: "string", path: []string{"Next", "StringValue"}},
			},
			expected: &struct {
				StringValue string
				Next        basicAppConfig
			}{"FOO", basicAppConfig{
				BoolValue:   true,
				IntValue:    1,
				StringValue: "string",
			}},
		},
		{
			desc: "with a pointer to a struct",
			source: &struct {
				StringValue string
				NextPointer *basicAppConfig
			}{},
			values: []*envValue{
				{strValue: "FOO", path: []string{"StringValue"}},
				{strValue: "1", path: []string{"NextPointer", "IntValue"}},
				{strValue: "true", path: []string{"NextPointer", "BoolValue"}},
				{strValue: "string", path: []string{"NextPointer", "StringValue"}},
			},
			expected: &struct {
				StringValue string
				NextPointer *basicAppConfig
			}{"FOO", &basicAppConfig{
				BoolValue:   true,
				IntValue:    1,
				StringValue: "string",
			}},
		},
		{
			desc: "with a pointer to pointer struct",
			source: &struct {
				StringValue string
				NextPointer **basicAppConfig
			}{},
			values: []*envValue{
				{strValue: "FOO", path: []string{"StringValue"}},
				{strValue: "1", path: []string{"NextPointer", "IntValue"}},
				{strValue: "true", path: []string{"NextPointer", "BoolValue"}},
				{strValue: "string", path: []string{"NextPointer", "StringValue"}},
			},
			expected: getPtrPtrConfig(),
		},
		{
			desc:   "when the environment values contain a wrong path",
			source: &delegatorType{},
			values: []*envValue{
				{strValue: "FOO", path: []string{"WrongPath"}},
			},
			expected: &delegatorType{},
		},
		{
			desc:   "with interface delegation",
			source: &delegatorType{},
			values: []*envValue{
				{strValue: "FOO", path: []string{"StringValue"}},
				{strValue: "1", path: []string{"IntValue"}},
			},
			expected: &delegatorType{
				IntValue:    1,
				StringValue: "FOO",
			},
		},
		{
			desc: "with a map[string]string",
			source: &struct {
				Config map[string]string
			}{},
			values: []*envValue{
				{strValue: "FOO", path: []string{"Config", "foo"}},
				{strValue: "MEH", path: []string{"Config", "bar"}},
				{strValue: "BAR", path: []string{"Config", "biz"}},
			},
			expected: &struct {
				Config map[string]string
			}{
				Config: map[string]string{
					"foo": "FOO",
					"bar": "MEH",
					"biz": "BAR",
				},
			},
		},
		{
			desc: "with a map that uses a parser for values",
			source: &struct {
				Config map[string]int
			}{},
			values: []*envValue{
				{strValue: "1", path: []string{"Config", "foo"}},
				{strValue: "2", path: []string{"Config", "bar"}},
				{strValue: "3", path: []string{"Config", "biz"}},
			},
			expected: &struct {
				Config map[string]int
			}{
				Config: map[string]int{
					"foo": 1,
					"bar": 2,
					"biz": 3,
				},
			},
		},
		{
			desc: "with a map that needs a parser for keys",
			source: &struct {
				Config map[int]int
			}{},
			values: []*envValue{
				{strValue: "1", path: []string{"Config", "1"}},
				{strValue: "2", path: []string{"Config", "2"}},
				{strValue: "3", path: []string{"Config", "3"}},
			},
			expected: &struct {
				Config map[int]int
			}{
				Config: map[int]int{
					1: 1,
					2: 2,
					3: 3,
				},
			},
		},
		{
			desc: "with a map[string]struct",
			source: &struct {
				Config map[string]basicAppConfig
			}{},
			values: []*envValue{
				{strValue: "FOOO", path: []string{"Config", "foo", "StringValue"}},
				{strValue: "10", path: []string{"Config", "foo", "IntValue"}},
			},
			expected: &struct {
				Config map[string]basicAppConfig
			}{
				Config: map[string]basicAppConfig{
					"foo": {
						StringValue: "FOOO",
						IntValue:    10,
					},
				},
			},
		},
		{
			desc: "with a map[int]struct",
			source: &struct {
				Config map[int]basicAppConfig
			}{},
			values: []*envValue{
				{strValue: "FOOO", path: []string{"Config", "0", "StringValue"}},
				{strValue: "10", path: []string{"Config", "0", "IntValue"}},
			},
			expected: &struct {
				Config map[int]basicAppConfig
			}{
				Config: map[int]basicAppConfig{
					0: {
						StringValue: "FOOO",
						IntValue:    10,
					},
				},
			},
		},
		{
			desc: "with an array of int",
			source: &struct {
				Config []int
			}{},
			values: []*envValue{
				{strValue: "1", path: []string{"Config", "0"}},
				{strValue: "10", path: []string{"Config", "1"}},
			},
			expected: &struct {
				Config []int
			}{
				Config: []int{1, 10},
			},
		},
		{
			desc: "with an array of struct",
			source: &struct {
				Config []basicAppConfig
			}{},
			values: []*envValue{
				{strValue: "Test", path: []string{"Config", "0", "StringValue"}},
				{strValue: "10", path: []string{"Config", "0", "IntValue"}},
				{strValue: "true", path: []string{"Config", "0", "BoolValue"}},
				{strValue: "Test2", path: []string{"Config", "1", "StringValue"}},
				{strValue: "20", path: []string{"Config", "1", "IntValue"}},
				{strValue: "false", path: []string{"Config", "1", "BoolValue"}},
			},
			expected: &struct {
				Config []basicAppConfig
			}{
				Config: []basicAppConfig{
					{
						BoolValue:   true,
						IntValue:    10,
						StringValue: "Test",
					},
					{
						BoolValue:   false,
						IntValue:    20,
						StringValue: "Test2",
					},
				},
			},
		},
		{
			desc: "with an array of pointer to struct",
			source: &struct {
				Config []*basicAppConfig
			}{},
			values: []*envValue{
				{strValue: "Test", path: []string{"Config", "0", "StringValue"}},
				{strValue: "10", path: []string{"Config", "0", "IntValue"}},
				{strValue: "true", path: []string{"Config", "0", "BoolValue"}},
				{strValue: "Test2", path: []string{"Config", "1", "StringValue"}},
				{strValue: "20", path: []string{"Config", "1", "IntValue"}},
				{strValue: "false", path: []string{"Config", "1", "BoolValue"}},
			},
			expected: &struct {
				Config []*basicAppConfig
			}{
				Config: []*basicAppConfig{
					{
						BoolValue:   true,
						IntValue:    10,
						StringValue: "Test",
					},
					{
						BoolValue:   false,
						IntValue:    20,
						StringValue: "Test2",
					},
				},
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run("should assign values to config "+test.desc, func(t *testing.T) {
			t.Parallel()

			err := subject.assignValues(reflect.ValueOf(test.source).Elem(), test.values, []string{})
			require.NoError(t, err)

			assert.Exactly(t, test.expected, test.source)
		})
	}
}
