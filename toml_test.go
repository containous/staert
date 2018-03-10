package staert

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/containous/flaeg"
	"github.com/containous/flaeg/parse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTomlSourceTrivial(t *testing.T) {
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	defaultPointersConfig := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    11,
			S1String: "S1StringDefaultPointersConfig",
			S1Bool:   true,
			S1PtrStruct3: &Struct3{
				S3Float64: 11.11,
			},
		},
		PtrStruct2: &Struct2{
			S2Int64:  22,
			S2String: "S2StringDefaultPointersConfig",
			S2Bool:   false,
		},
	}

	// Test
	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig,
		Run: func() error {
			return nil
		},
	}

	s := NewStaert(rootCmd)
	toml := NewTomlSource("trivial", []string{"./fixtures/", "/any/other/path"})
	s.AddSource(toml)

	err := s.parseConfigAllSources(rootCmd)
	require.NoError(t, err)

	// expected
	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    28,
			S1String: "S1StringDefaultPointersConfig",
			S1Bool:   true,
		},
		DurationField: parse.Duration(28 * time.Second),
	}

	assert.Exactly(t, expected, rootCmd.Config)
}

func TestTomlSourceFieldUnderPointer(t *testing.T) {
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	defaultPointersConfig := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    11,
			S1String: "S1StringDefaultPointersConfig",
			S1Bool:   true,
			S1PtrStruct3: &Struct3{
				S3Float64: 11.11,
			},
		},
		PtrStruct2: &Struct2{
			S2Int64:  22,
			S2String: "S2StringDefaultPointersConfig",
			S2Bool:   false,
		},
	}

	// Test
	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig,
		Run: func() error {
			return nil
		},
	}

	s := NewStaert(rootCmd)
	toml := NewTomlSource("fieldUnderPointer", []string{"./fixtures/", "/any/other/path"})
	s.AddSource(toml)

	err := s.parseConfigAllSources(rootCmd)
	require.NoError(t, err)

	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    11,
			S1String: "S1StringDefaultPointersConfig",
			S1Bool:   true,
		},
		DurationField: parse.Duration(42 * time.Second),
	}

	assert.Exactly(t, expected, rootCmd.Config)
}

func TestTomlSourcePointerUnderPointer(t *testing.T) {
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	defaultPointersConfig := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    11,
			S1String: "S1StringDefaultPointersConfig",
			S1Bool:   true,
			S1PtrStruct3: &Struct3{
				S3Float64: 11.11,
			},
		},
		PtrStruct2: &Struct2{
			S2Int64:  22,
			S2String: "S2StringDefaultPointersConfig",
			S2Bool:   false,
		},
	}

	// Test
	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig,
		Run: func() error {
			return nil
		},
	}

	s := NewStaert(rootCmd)
	toml := NewTomlSource("pointerUnderPointer", []string{"./fixtures/", "/any/other/path"})
	s.AddSource(toml)

	err := s.parseConfigAllSources(rootCmd)
	require.NoError(t, err)

	// expected
	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    11,
			S1String: "S1StringDefaultPointersConfig",
			S1Bool:   true,
			S1PtrStruct3: &Struct3{
				S3Float64: 11.11,
			},
		},
		DurationField: parse.Duration(time.Second),
	}

	assert.Exactly(t, expected, rootCmd.Config)
}

func TestTomlSourceFieldUnderPointerUnderPointer(t *testing.T) {
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	defaultPointersConfig := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    11,
			S1String: "S1StringDefaultPointersConfig",
			S1Bool:   true,
			S1PtrStruct3: &Struct3{
				S3Float64: 11.11,
			},
		},
		PtrStruct2: &Struct2{
			S2Int64:  22,
			S2String: "S2StringDefaultPointersConfig",
			S2Bool:   false,
		},
	}

	// Test
	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig,
		Run: func() error {
			return nil
		},
	}

	s := NewStaert(rootCmd)
	toml := NewTomlSource("fieldUnderPtrUnderPtr", []string{"./fixtures/", "/any/other/path"})
	s.AddSource(toml)

	err := s.parseConfigAllSources(rootCmd)
	require.NoError(t, err)

	// expected
	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
			S1PtrStruct3: &Struct3{
				S3Float64: 28.28,
			},
		},
		DurationField: parse.Duration(time.Second),
	}

	assert.Exactly(t, expected, rootCmd.Config)
}

func TestTomlSourcePointer(t *testing.T) {
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	defaultPointersConfig := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    11,
			S1String: "S1StringDefaultPointersConfig",
			S1Bool:   true,
			S1PtrStruct3: &Struct3{
				S3Float64: 11.11,
			},
		},
		PtrStruct2: &Struct2{
			S2Int64:  22,
			S2String: "S2StringDefaultPointersConfig",
			S2Bool:   false,
		},
	}

	// Test
	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig,
		Run: func() error {
			return nil
		},
	}

	s := NewStaert(rootCmd)
	toml := NewTomlSource("pointer", []string{"./fixtures/", "/any/other/path"})
	s.AddSource(toml)

	err := s.parseConfigAllSources(rootCmd)
	require.NoError(t, err)

	// expected
	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		PtrStruct2: &Struct2{
			S2Int64:  22,
			S2String: "S2StringDefaultPointersConfig",
			S2Bool:   false,
		},
		DurationField: parse.Duration(time.Second),
	}

	assert.Exactly(t, expected, rootCmd.Config)
}

func TestTomlSourceNothing(t *testing.T) {
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	defaultPointersConfig := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    11,
			S1String: "S1StringDefaultPointersConfig",
			S1Bool:   true,
			S1PtrStruct3: &Struct3{
				S3Float64: 11.11,
			},
		},
		PtrStruct2: &Struct2{
			S2Int64:  22,
			S2String: "S2StringDefaultPointersConfig",
			S2Bool:   false,
		},
	}

	// Test
	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig,
		Run: func() error {
			return nil
		},
	}

	s := NewStaert(rootCmd)
	toml := NewTomlSource("nothing", []string{"./fixtures/", "/any/other/path"})
	s.AddSource(toml)

	err := s.parseConfigAllSources(rootCmd)
	require.NoError(t, err)

	// expected
	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	assert.Exactly(t, expected, rootCmd.Config)
}

func TestMergeTomlNothingFlaegNoArgs(t *testing.T) {
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	defaultPointersConfig := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    11,
			S1String: "S1StringDefaultPointersConfig",
			S1Bool:   true,
			S1PtrStruct3: &Struct3{
				S3Float64: 11.11,
			},
		},
		PtrStruct2: &Struct2{
			S2Int64:  22,
			S2String: "S2StringDefaultPointersConfig",
			S2Bool:   false,
		},
	}

	var args []string

	// Test
	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig,
		Run: func() error {
			return nil
		},
	}

	s := NewStaert(rootCmd)
	toml := NewTomlSource("nothing", []string{"./fixtures/", "/any/other/path"})
	s.AddSource(toml)
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)

	err := s.parseConfigAllSources(rootCmd)
	require.NoError(t, err)

	// expected
	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	assert.Exactly(t, expected, rootCmd.Config)
}

func TestMergeTomlFieldUnderPointerUnderPointerFlaegNoArgs(t *testing.T) {
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	defaultPointersConfig := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    11,
			S1String: "S1StringDefaultPointersConfig",
			S1Bool:   true,
			S1PtrStruct3: &Struct3{
				S3Float64: 11.11,
			},
		},
		PtrStruct2: &Struct2{
			S2Int64:  22,
			S2String: "S2StringDefaultPointersConfig",
			S2Bool:   false,
		},
	}

	var args []string

	// Test
	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig,
		Run: func() error {
			return nil
		},
	}

	s := NewStaert(rootCmd)
	toml := NewTomlSource("fieldUnderPtrUnderPtr", []string{"./fixtures/", "/any/other/path"})
	s.AddSource(toml)
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)

	err := s.parseConfigAllSources(rootCmd)
	require.NoError(t, err)

	// expected
	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
			S1PtrStruct3: &Struct3{
				S3Float64: 28.28,
			},
		},
		DurationField: parse.Duration(time.Second),
	}

	assert.Exactly(t, expected, rootCmd.Config)
}

func TestMergeTomlTrivialFlaegOverwriteField(t *testing.T) {
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	defaultPointersConfig := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    11,
			S1String: "S1StringDefaultPointersConfig",
			S1Bool:   true,
			S1PtrStruct3: &Struct3{
				S3Float64: 11.11,
			},
		},
		PtrStruct2: &Struct2{
			S2Int64:  22,
			S2String: "S2StringDefaultPointersConfig",
			S2Bool:   false,
		},
	}

	args := []string{"--ptrstruct1.s1int=55"}

	// Test
	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig,
		Run: func() error {
			return nil
		},
	}

	s := NewStaert(rootCmd)
	toml := NewTomlSource("trivial", []string{"./fixtures/", "/any/other/path"})
	s.AddSource(toml)
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)

	err := s.parseConfigAllSources(rootCmd)
	require.NoError(t, err)

	// expected
	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    55,
			S1String: "S1StringDefaultPointersConfig",
			S1Bool:   true,
		},
		DurationField: parse.Duration(28 * time.Second),
	}

	assert.Exactly(t, expected, rootCmd.Config)
}

func TestMergeTomlPointerUnderPointerFlaegManyArgs(t *testing.T) {
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	defaultPointersConfig := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    11,
			S1String: "S1StringDefaultPointersConfig",
			S1Bool:   true,
			S1PtrStruct3: &Struct3{
				S3Float64: 11.11,
			},
		},
		PtrStruct2: &Struct2{
			S2Int64:  22,
			S2String: "S2StringDefaultPointersConfig",
			S2Bool:   false,
		},
	}

	args := []string{
		"--ptrstruct1.s1int=55",
		"--durationfield=55s",
		"--ptrstruct2.s2string=S2StringFlaeg",
	}

	// Test
	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig,
		Run: func() error {
			return nil
		},
	}

	s := NewStaert(rootCmd)
	toml := NewTomlSource("pointerUnderPointer", []string{"./fixtures/", "/any/other/path"})
	s.AddSource(toml)
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)

	err := s.parseConfigAllSources(rootCmd)
	require.NoError(t, err)

	// expected
	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    55,
			S1String: "S1StringDefaultPointersConfig",
			S1Bool:   true,
			S1PtrStruct3: &Struct3{
				S3Float64: 11.11,
			},
		},
		PtrStruct2: &Struct2{
			S2Int64:  22,
			S2String: "S2StringFlaeg",
			S2Bool:   false,
		},
		DurationField: parse.Duration(55 * time.Second),
	}

	assert.Exactly(t, expected, rootCmd.Config)
}

func TestPreprocessDir(t *testing.T) {
	here, err := filepath.Abs(".")
	if err != nil {
		require.NoError(t, err)
	}

	testCases := []struct {
		directory string
		expected  string
	}{
		{
			directory: ".",
			expected:  here,
		},
		{
			directory: "dir1/dir2",
			expected:  here + "/dir1/dir2",
		},
		{
			directory: "/etc/test",
			expected:  "/etc/test",
		},
		{
			directory: "/etc/dir1/file1.ext",
			expected:  "/etc/dir1/file1.ext",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.directory, func(t *testing.T) {
			t.Parallel()

			out, err := preprocessDir(test.directory)
			require.NoError(t, err)

			// always check against the absolute path
			expectedPath, _ := filepath.Abs(test.expected)
			assert.Equal(t, expectedPath, out)
		})
	}
}

func TestPreprocessDirEnvVariablesExpansions(t *testing.T) {
	err := os.Setenv("TEST_ENV_VARIABLE", "/some/path")
	require.NoError(t, err)

	in := "$TEST_ENV_VARIABLE/my/path"

	expectedPath, _ := filepath.Abs("/some/path/my/path")

	out, err := preprocessDir(in)
	require.NoError(t, err)

	assert.Equal(t, expectedPath, out)
}

func TestFindFile(t *testing.T) {
	result := findFile("nothing", []string{"", "$HOME/test", "fixtures"})

	// expected
	here, err := filepath.Abs(".")
	require.NoError(t, err)

	expected := filepath.Join(here, "fixtures", "nothing.toml")
	assert.Equal(t, expected, result)
}

func TestFindFileSliceFileAndDirLastIf(t *testing.T) {
	thisPath, _ := filepath.Abs(".")

	expected := filepath.Join(thisPath, "/fixtures/trivial.toml")
	result := findFile("trivial", []string{"./fixtures/", "/any/other/path"})

	assert.Equal(t, expected, result)
}

func TestFindFileSliceFileAndDirFirstIf(t *testing.T) {
	inFilename := ""
	inDirNfile := []string{"$PWD/fixtures/nothing.toml"}

	thisPath, _ := filepath.Abs(".")
	expected := filepath.Join(thisPath, "/fixtures/nothing.toml")
	result := findFile(inFilename, inDirNfile)

	assert.Equal(t, expected, result)
}
