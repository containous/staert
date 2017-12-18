package staert

import (
	"os"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/containous/flaeg/parse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type basicAppConfig struct {
	StringValue string
	IntValue    int
	BoolValue   bool
}

type typeInterface interface {
	Foo() string
}

type delegatorType struct {
	typeInterface
	IntValue    int
	StringValue string
}

type sortableEnvValues []*envValue

func (s sortableEnvValues) Len() int {
	return len(s)
}

func (s sortableEnvValues) Less(i, j int) bool {
	return s[i].strValue < s[j].strValue
}

func (s sortableEnvValues) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Users authentication users
type Users []string

type Basic struct {
	Users     `mapstructure:","`
	UsersFile string
}

func setupEnv(t *testing.T, env map[string]string) {
	t.Helper()

	for k, v := range env {
		err := os.Setenv(k, v)
		require.NoError(t, err)
	}
}

func cleanupEnv(t *testing.T, env map[string]string) {
	t.Helper()

	for k := range env {
		err := os.Unsetenv(k)
		require.NoError(t, err)
	}
}

func analyzeStructShouldSucceed(t *testing.T, expectation, result sortableEnvValues, err error) {
	t.Helper()

	require.NoError(t, err)
	assert.Lenf(t, result, len(expectation), "Unexpected count of values returned")

	// Sort by value, according to strValue (which might not be the best
	// idea ever), in order to ensure index based comparison consistency
	sort.Sort(expectation)
	sort.Sort(result)

	for i, v := range expectation {
		assert.Equal(t, v.strValue, result[i].strValue)
		assert.Exactly(t, v.path, result[i].path)
	}
}

func analyzeStructShouldFail(t *testing.T, expectation, result sortableEnvValues, err error) {
	t.Helper()

	require.Error(t, err)
}

func TestAnalyzeStruct(t *testing.T) {
	subject := NewEnvSource("", "_", map[reflect.Type]parse.Parser{})

	testCases := []struct {
		desc     string
		source   interface{}
		expected []*envValue
		env      map[string]string
		then     func(t *testing.T, expectation, result sortableEnvValues, err error)
	}{
		{
			desc:   "should succeed with basic struct",
			source: &basicAppConfig{},
			expected: []*envValue{
				{strValue: "FOOO", path: []string{"StringValue"}},
				{strValue: "10", path: []string{"IntValue"}},
				{strValue: "true", path: []string{"BoolValue"}},
			},
			env: map[string]string{
				"STRING_VALUE": "FOOO",
				"INT_VALUE":    "10",
				"BOOL_VALUE":   "true",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc: "should succeed with unexported fields",
			source: &struct {
				unexported string
				IntValue   int
			}{},
			expected: []*envValue{
				{strValue: "10", path: []string{"IntValue"}},
			},
			env: map[string]string{
				"UNEXPORTED": "FOOO",
				"INT_VALUE":  "10",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc: "should succeed with embedded struct",
			source: &struct {
				basicAppConfig
				FloatValue float32
			}{},
			expected: []*envValue{
				{strValue: "FOOO", path: []string{"StringValue"}},
				{strValue: "10", path: []string{"IntValue"}},
				{strValue: "true", path: []string{"BoolValue"}},
				{strValue: "42.1", path: []string{"FloatValue"}},
			},
			env: map[string]string{
				"STRING_VALUE": "FOOO",
				"INT_VALUE":    "10",
				"BOOL_VALUE":   "true",
				"FLOAT_VALUE":  "42.1",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc:   "should succeed with nested struct value",
			source: &struct{ Config basicAppConfig }{},
			expected: []*envValue{
				{strValue: "FOOO", path: []string{"Config", "StringValue"}},
				{strValue: "10", path: []string{"Config", "IntValue"}},
				{strValue: "true", path: []string{"Config", "BoolValue"}},
			},
			env: map[string]string{
				"CONFIG_STRING_VALUE": "FOOO",
				"CONFIG_INT_VALUE":    "10",
				"CONFIG_BOOL_VALUE":   "true",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc: "should succeed with double nested struct value",
			source: &struct {
				Nested struct {
					Config basicAppConfig
				}
			}{},
			expected: []*envValue{
				{strValue: "FOOO", path: []string{"Nested", "Config", "StringValue"}},
				{strValue: "10", path: []string{"Nested", "Config", "IntValue"}},
				{strValue: "true", path: []string{"Nested", "Config", "BoolValue"}},
			},
			env: map[string]string{
				"NESTED_CONFIG_STRING_VALUE": "FOOO",
				"NESTED_CONFIG_INT_VALUE":    "10",
				"NESTED_CONFIG_BOOL_VALUE":   "true",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc:   "should succeed with nested struct pointer",
			source: &struct{ Config *basicAppConfig }{},
			expected: []*envValue{
				{strValue: "FOOO", path: []string{"Config", "StringValue"}},
				{strValue: "10", path: []string{"Config", "IntValue"}},
				{strValue: "true", path: []string{"Config", "BoolValue"}},
			},
			env: map[string]string{
				"CONFIG_STRING_VALUE": "FOOO",
				"CONFIG_INT_VALUE":    "10",
				"CONFIG_BOOL_VALUE":   "true",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc: "should succeed with double nested struct pointer",
			source: &struct {
				Nested *struct {
					Config *basicAppConfig
				}
			}{},
			expected: []*envValue{
				{strValue: "FOOO", path: []string{"Nested", "Config", "StringValue"}},
				{strValue: "10", path: []string{"Nested", "Config", "IntValue"}},
				{strValue: "true", path: []string{"Nested", "Config", "BoolValue"}},
			},
			env: map[string]string{
				"NESTED_CONFIG_STRING_VALUE": "FOOO",
				"NESTED_CONFIG_INT_VALUE":    "10",
				"NESTED_CONFIG_BOOL_VALUE":   "true",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc: "should succeed with nested struct pointer to struct",
			source: &struct {
				Nested *struct {
					Config basicAppConfig
				}
			}{},
			expected: []*envValue{
				{strValue: "FOOO", path: []string{"Nested", "Config", "StringValue"}},
				{strValue: "10", path: []string{"Nested", "Config", "IntValue"}},
				{strValue: "true", path: []string{"Nested", "Config", "BoolValue"}},
			},
			env: map[string]string{
				"NESTED_CONFIG_STRING_VALUE": "FOOO",
				"NESTED_CONFIG_INT_VALUE":    "10",
				"NESTED_CONFIG_BOOL_VALUE":   "true",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc:   "should succeed with pointer to value",
			source: &struct{ IntValue *int }{},
			expected: []*envValue{
				{strValue: "10", path: []string{"IntValue"}},
			},
			env: map[string]string{
				"INT_VALUE": "10",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc: "should succeed with nested pointer to value",
			source: &struct {
				Config struct {
					IntValue *int
				}
			}{},
			expected: []*envValue{
				{strValue: "10", path: []string{"Config", "IntValue"}},
			},
			env: map[string]string{
				"CONFIG_INT_VALUE": "10",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc:   "should succeed with double pointer to int value",
			source: &struct{ Config **int }{},
			expected: []*envValue{
				{strValue: "10", path: []string{"Config"}},
			},
			env: map[string]string{
				"CONFIG": "10",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc:   "should succeed with a double pointer to struct",
			source: &struct{ Config **basicAppConfig }{},
			expected: []*envValue{
				{strValue: "FOOO", path: []string{"Config", "StringValue"}},
				{strValue: "10", path: []string{"Config", "IntValue"}},
				{strValue: "true", path: []string{"Config", "BoolValue"}},
			},
			env: map[string]string{
				"CONFIG_STRING_VALUE": "FOOO",
				"CONFIG_INT_VALUE":    "10",
				"CONFIG_BOOL_VALUE":   "true",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc:   "should succeed with interface delegation",
			source: &delegatorType{},
			expected: []*envValue{
				{strValue: "FOOO", path: []string{"StringValue"}},
				{strValue: "10", path: []string{"IntValue"}},
			},
			env: map[string]string{
				"STRING_VALUE": "FOOO",
				"INT_VALUE":    "10",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc:   "should succeed with a map[string]string",
			source: &struct{ Config map[string]string }{},
			expected: []*envValue{
				{strValue: "FOO", path: []string{"Config", "foo"}},
				{strValue: "MEH", path: []string{"Config", "bar"}},
				{strValue: "BAR", path: []string{"Config", "biz"}},
			},
			env: map[string]string{
				"CONFIG_FOO": "FOO",
				"CONFIG_BAR": "MEH",
				"CONFIG_BIZ": "BAR",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc:   "should succeed with a map[string]struct",
			source: &struct{ Config map[string]basicAppConfig }{},
			expected: []*envValue{
				{strValue: "FOO", path: []string{"Config", "foo", "StringValue"}},
				{strValue: "MEH", path: []string{"Config", "bar", "StringValue"}},
				{strValue: "BAR", path: []string{"Config", "biz", "StringValue"}},
			},
			env: map[string]string{
				"CONFIG_FOO_STRING_VALUE": "FOO",
				"CONFIG_BAR_STRING_VALUE": "MEH",
				"CONFIG_BIZ_STRING_VALUE": "BAR",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc:   "should succeed with a map[string]*struct",
			source: &struct{ Config map[string]*basicAppConfig }{},
			expected: []*envValue{
				{strValue: "FOO", path: []string{"Config", "foo", "StringValue"}},
				{strValue: "MEH", path: []string{"Config", "bar", "StringValue"}},
				{strValue: "BAR", path: []string{"Config", "biz", "StringValue"}},
			},
			env: map[string]string{
				"CONFIG_FOO_STRING_VALUE": "FOO",
				"CONFIG_BAR_STRING_VALUE": "MEH",
				"CONFIG_BIZ_STRING_VALUE": "BAR",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc: "should succeed with a map of map",
			source: &struct {
				Config map[int]map[string]*basicAppConfig
			}{},
			expected: []*envValue{
				{strValue: "FOO", path: []string{"Config", "0", "foo", "StringValue"}},
				{strValue: "MEH", path: []string{"Config", "1", "bar", "StringValue"}},
				{strValue: "BAR", path: []string{"Config", "0", "biz", "StringValue"}},
			},
			env: map[string]string{
				"CONFIG_0_FOO_STRING_VALUE": "FOO",
				"CONFIG_1_BAR_STRING_VALUE": "MEH",
				"CONFIG_0_BIZ_STRING_VALUE": "BAR",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc:   "should succeed with a slice of ints",
			source: &struct{ Config []int }{},
			expected: []*envValue{
				{strValue: "FOOO", path: []string{"Config", "0"}},
				{strValue: "10", path: []string{"Config", "1"}},
				{strValue: "true", path: []string{"Config", "2"}},
			},
			env: map[string]string{
				"CONFIG_0": "FOOO",
				"CONFIG_1": "10",
				"CONFIG_2": "true",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc:     "should fail with an invalid key for a slice",
			source:   &struct{ Config []int }{},
			expected: []*envValue{},
			env: map[string]string{
				"CONFIG_0":      "FOOO",
				"CONFIG_1":      "10",
				"CONFIG_PATATE": "true",
			},
			then: analyzeStructShouldFail,
		},
		{
			desc:   "should succeed with an array to value",
			source: &struct{ Config [10]int }{},
			expected: []*envValue{
				{strValue: "FOOO", path: []string{"Config", "0"}},
				{strValue: "10", path: []string{"Config", "1"}},
				{strValue: "true", path: []string{"Config", "2"}},
			},
			env: map[string]string{
				"CONFIG_0": "FOOO",
				"CONFIG_1": "10",
				"CONFIG_2": "true",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc:     "should fail with an array and an out of bound index",
			source:   &struct{ Config [10]int }{},
			expected: []*envValue{},
			env: map[string]string{
				"CONFIG_11": "10",
			},
			then: analyzeStructShouldFail,
		},
		{
			desc:   "should succeed with an array to value",
			source: &struct{ Config [10]int }{},
			expected: []*envValue{
				{strValue: "FOOO", path: []string{"Config", "0"}},
				{strValue: "10", path: []string{"Config", "1"}},
				{strValue: "true", path: []string{"Config", "2"}},
			},
			env: map[string]string{
				"CONFIG_0": "FOOO",
				"CONFIG_1": "10",
				"CONFIG_2": "true",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc:   "should succeed with a slice of struct",
			source: &struct{ Config []basicAppConfig }{},
			expected: []*envValue{
				{strValue: "FOOO", path: []string{"Config", "0", "StringValue"}},
				{strValue: "10", path: []string{"Config", "0", "IntValue"}},
				{strValue: "MIMI", path: []string{"Config", "1", "StringValue"}},
				{strValue: "15", path: []string{"Config", "1", "IntValue"}},
			},
			env: map[string]string{
				"CONFIG_0_STRING_VALUE": "FOOO",
				"CONFIG_0_INT_VALUE":    "10",
				"CONFIG_1_STRING_VALUE": "MIMI",
				"CONFIG_1_INT_VALUE":    "15",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc:   "should succeed with a [][]struct",
			source: &struct{ Config [][]basicAppConfig }{},
			expected: []*envValue{
				{strValue: "FOOO", path: []string{"Config", "0", "0", "StringValue"}},
				{strValue: "10", path: []string{"Config", "0", "0", "IntValue"}},
				{strValue: "MIMI", path: []string{"Config", "1", "1", "StringValue"}},
				{strValue: "15", path: []string{"Config", "1", "1", "IntValue"}},
			},
			env: map[string]string{
				"CONFIG_0_0_STRING_VALUE": "FOOO",
				"CONFIG_0_0_INT_VALUE":    "10",
				"CONFIG_1_1_STRING_VALUE": "MIMI",
				"CONFIG_1_1_INT_VALUE":    "15",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc:   "should succeed with a []map[string]struct",
			source: &struct{ Config []map[string]basicAppConfig }{},
			expected: []*envValue{
				{strValue: "FOOO", path: []string{"Config", "0", "foo", "StringValue"}},
				{strValue: "10", path: []string{"Config", "0", "foo", "IntValue"}},
				{strValue: "MIMI", path: []string{"Config", "1", "bar", "StringValue"}},
				{strValue: "15", path: []string{"Config", "1", "bar", "IntValue"}},
			},
			env: map[string]string{
				"CONFIG_0_FOO_STRING_VALUE": "FOOO",
				"CONFIG_0_FOO_INT_VALUE":    "10",
				"CONFIG_1_BAR_STRING_VALUE": "MIMI",
				"CONFIG_1_BAR_INT_VALUE":    "15",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc: "should succeed when config has exported fields that are of type func",
			source: &struct {
				Config basicAppConfig
				Time   func() time.Time
			}{},
			expected: []*envValue{
				{strValue: "FOOO", path: []string{"Config", "StringValue"}},
				{strValue: "10", path: []string{"Config", "IntValue"}},
				{strValue: "true", path: []string{"Config", "BoolValue"}},
			},
			env: map[string]string{
				"CONFIG_STRING_VALUE": "FOOO",
				"CONFIG_INT_VALUE":    "10",
				"CONFIG_BOOL_VALUE":   "true",
			},
			then: analyzeStructShouldSucceed,
		},
		{
			desc: "should succeed with an type alias to an array",
			source: &struct {
				Basic     *Basic
				UsersFile string
			}{},
			expected: []*envValue{
				{strValue: "UserZero", path: []string{"Basic", "0"}},
				{strValue: "UserOne", path: []string{"Basic", "1"}},
				{strValue: "path/to/file", path: []string{"UsersFile"}},
			},
			env: map[string]string{
				"BASIC_0":    "UserZero",
				"BASIC_1":    "UserOne",
				"USERS_FILE": "path/to/file",
			},
			then: analyzeStructShouldSucceed,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			setupEnv(t, test.env)
			defer cleanupEnv(t, test.env)

			res, err := subject.analyzeStruct(reflect.TypeOf(test.source).Elem(), nil)
			test.then(t, test.expected, res, err)
		})
	}
}

func TestEnvVarFromPath(t *testing.T) {
	testCases := []struct {
		desc        string
		prefix      string
		separator   string
		paths       []string
		expectation string
	}{
		{
			desc:        "BlankPrefix",
			prefix:      "",
			separator:   "_",
			paths:       []string{"Foo"},
			expectation: "FOO",
		},
		{
			desc:        "NonBlankPrefix",
			prefix:      "YOUPI",
			separator:   "_",
			paths:       []string{"Foo"},
			expectation: "YOUPI_FOO",
		},
		{
			desc:        "CamelCasedPathMembers",
			prefix:      "YOUPI",
			separator:   "_",
			paths:       []string{"Foo", "IamGroot", "IAmBatman"},
			expectation: "YOUPI_FOO_IAM_GROOT_I_AM_BATMAN",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			subject := NewEnvSource(test.prefix, test.separator, map[reflect.Type]parse.Parser{})

			result := subject.envVarFromPath(test.paths)
			assert.Exactly(t, test.expectation, result)
		})
	}
}

func TestAnalyzeAndAssignFlowWithArrayConfig(t *testing.T) {
	sourceConfig := struct {
		StringArray []string
	}{}

	test := struct {
		source      interface{}
		expectation []*envValue
		env         map[string]string
	}{
		source: &sourceConfig,
		expectation: []*envValue{
			{strValue: "one", path: []string{"StringArray", "0"}},
			{strValue: "two", path: []string{"StringArray", "1"}},
		},
		env: map[string]string{
			"STRING_ARRAY_0": "one",
			"STRING_ARRAY_1": "two",
		},
	}

	setupEnv(t, test.env)
	defer cleanupEnv(t, test.env)

	parsers, err := parse.LoadParsers(nil)
	require.NoError(t, err)

	subject := NewEnvSource("", "_", parsers)

	res, err := subject.analyzeStruct(reflect.TypeOf(test.source).Elem(), nil)
	require.NoError(t, err)

	err = subject.assignValues(reflect.ValueOf(&sourceConfig), res, nil)
	require.NoError(t, err)

	require.ElementsMatch(t, []string{"one", "two"}, sourceConfig.StringArray)
}

func TestNextLevelKeys(t *testing.T) {
	subject := NewEnvSource("", "_", map[reflect.Type]parse.Parser{})

	testCases := []struct {
		desc     string
		prefix   string
		env      []string
		expected []string
	}{
		{
			desc:   "should strip the key part from the env name",
			prefix: "CONFIG_APP",
			env: []string{
				"CONFIG_APP_BATMAN_FOO",
				"CONFIG_APP_ROBIN_FOO",
				"CONFIG_APP_JOCKER_FOO",
			},
			expected: []string{
				"CONFIG_APP_BATMAN",
				"CONFIG_APP_ROBIN",
				"CONFIG_APP_JOCKER",
			},
		},
		{
			desc:   "should handle multiple equal keys",
			prefix: "CONFIG_APP",
			env: []string{
				"CONFIG_APP_BATMAN_FOO",
				"CONFIG_APP_ROBIN_FOO",
				"CONFIG_APP_JOCKER_FOO",
				"CONFIG_APP_BATMAN_BAR",
			},
			expected: []string{
				"CONFIG_APP_BATMAN",
				"CONFIG_APP_ROBIN",
				"CONFIG_APP_JOCKER",
				"CONFIG_APP_BATMAN",
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			res := subject.nextLevelKeys(test.prefix, test.env)
			assert.ElementsMatch(t, test.expected, res)
		})
	}
}

func TestEnvVarsWithPrefix(t *testing.T) {
	subject := NewEnvSource("", "_", map[reflect.Type]parse.Parser{})

	testCases := []struct {
		desc     string
		prefix   string
		env      map[string]string
		expected []string
	}{
		{
			desc:   "should filter out values that dont start with the prefix",
			prefix: "STAERT_APP",
			env: map[string]string{
				"STRING_VALUE":          "FOOO",
				"INT_VALUE":             "10",
				"BOOL_VALUE":            "true",
				"STAERT_APP_BOOL_VALUE": "true",
				"STAERT_APP_BAR_VALUE":  "true",
			},
			expected: []string{"STAERT_APP_BAR_VALUE", "STAERT_APP_BOOL_VALUE"},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			setupEnv(t, test.env)
			defer cleanupEnv(t, test.env)

			res := subject.envVarsWithPrefix(test.prefix)
			assert.ElementsMatch(t, test.expected, res)
		})
	}
}

func TestUnique(t *testing.T) {
	testCases := []struct {
		desc     string
		in       []string
		expected []string
	}{
		{
			desc:     "WithDuplicates",
			in:       []string{"FOO", "BAR", "BIZ", "FOO", "BIZ"},
			expected: []string{"FOO", "BAR", "BIZ"},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			res := unique(test.in)
			assert.ElementsMatch(t, test.expected, res)
		})
	}
}

func TestKeyFromEnvVar(t *testing.T) {
	subject := NewEnvSource("", "_", map[reflect.Type]parse.Parser{})

	testCases := []struct {
		desc     string
		prefix   string
		envVar   string
		expected string
	}{
		{
			desc:     "WithPrefix",
			prefix:   "CONFIG_APP",
			envVar:   "CONFIG_APP_BATMAN",
			expected: "batman",
		},
		{
			desc:     "WithPrefixAndSuffix",
			prefix:   "CONFIG_APP",
			envVar:   "CONFIG_APP_BATMAN_FOO",
			expected: "batman",
		},
		{
			desc:     "WithoutPrefix",
			prefix:   "",
			envVar:   "BATMAN",
			expected: "batman",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			res := subject.keyFromEnvVar(test.envVar, test.prefix)
			assert.Equal(t, test.expected, res)
		})
	}
}
