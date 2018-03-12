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

func TestTomlSource_Parse_Trivial(t *testing.T) {
	src := NewTomlSource("trivial", []string{"./fixtures/", "/any/other/path"})

	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	cmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig(),
		Run: func() error {
			return nil
		},
	}

	command, err := src.Parse(cmd)
	require.NoError(t, err)
	assert.Exactly(t, cmd, command)

	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    28,
			S1String: "S1StringDefaultPointersConfig",
			S1Bool:   true,
		},
		DurationField: parse.Duration(28 * time.Second),
	}

	assert.Exactly(t, expected, command.Config)
}

func TestTomlSource_Parse_FileNotFound(t *testing.T) {
	src := NewTomlSource("nothing", []string{"../path", "/any/other/path"})

	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	cmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig(),
		Run: func() error {
			return nil
		},
	}

	command, err := src.Parse(cmd)
	require.NoError(t, err)

	assert.Exactly(t, cmd, command)
}

func TestTomlSource_Parse_EmptyFile(t *testing.T) {
	src := NewTomlSource("nothing", []string{"./fixtures/", "/any/other/path"})

	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	cmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig(),
		Run: func() error {
			return nil
		},
	}

	command, err := src.Parse(cmd)
	require.NoError(t, err)
	assert.Exactly(t, cmd, command)

	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	assert.Exactly(t, expected, cmd.Config)
}

func TestTomlSource_Parse_FieldUnderPointer(t *testing.T) {
	src := NewTomlSource("fieldUnderPointer", []string{"./fixtures/", "/any/other/path"})

	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	cmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig(),
		Run: func() error {
			return nil
		},
	}

	command, err := src.Parse(cmd)
	require.NoError(t, err)
	assert.Exactly(t, cmd, command)

	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    11,
			S1String: "S1StringDefaultPointersConfig",
			S1Bool:   true,
		},
		DurationField: parse.Duration(42 * time.Second),
	}

	assert.Exactly(t, expected, cmd.Config)
}

func TestTomlSource_Parse_PointerUnderPointer(t *testing.T) {
	src := NewTomlSource("pointerUnderPointer", []string{"./fixtures/", "/any/other/path"})

	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	cmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig(),
		Run: func() error {
			return nil
		},
	}

	command, err := src.Parse(cmd)
	require.NoError(t, err)
	assert.Exactly(t, cmd, command)

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

	assert.Exactly(t, expected, cmd.Config)
}

func TestTomlSource_Parse_FieldUnderPointerUnderPointer(t *testing.T) {
	src := NewTomlSource("fieldUnderPtrUnderPtr", []string{"./fixtures/", "/any/other/path"})

	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	cmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig(),
		Run: func() error {
			return nil
		},
	}

	command, err := src.Parse(cmd)
	require.NoError(t, err)
	assert.Exactly(t, cmd, command)

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

	assert.Exactly(t, expected, cmd.Config)
}

func TestTomlSource_Parse_Pointer(t *testing.T) {
	src := NewTomlSource("pointer", []string{"./fixtures/", "/any/other/path"})

	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	cmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig(),
		Run: func() error {
			return nil
		},
	}

	command, err := src.Parse(cmd)
	require.NoError(t, err)
	assert.Exactly(t, cmd, command)

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

	assert.Exactly(t, expected, cmd.Config)
}

func Test_preProcessDir(t *testing.T) {
	here, err := filepath.Abs(".")
	require.NoError(t, err)

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

			out, err := preProcessDir(test.directory)
			require.NoError(t, err)

			// always check against the absolute path
			expectedPath, _ := filepath.Abs(test.expected)
			assert.Equal(t, expectedPath, out)
		})
	}
}

func Test_preProcessDir_envVariablesExpansions(t *testing.T) {
	err := os.Setenv("TEST_ENV_VARIABLE", "/some/path")
	require.NoError(t, err)

	in := "$TEST_ENV_VARIABLE/my/path"

	expectedPath, _ := filepath.Abs("/some/path/my/path")

	out, err := preProcessDir(in)
	require.NoError(t, err)

	assert.Equal(t, expectedPath, out)
}

func Test_findFile(t *testing.T) {
	here, err := filepath.Abs(".")
	require.NoError(t, err)

	expected := filepath.Join(here, "fixtures", "nothing.toml")

	result := findFile("nothing", []string{"", "$HOME/test", "fixtures"})
	assert.Equal(t, expected, result)
}

func Test_findFile_sliceFileAndDirLastIf(t *testing.T) {
	thisPath, _ := filepath.Abs(".")
	expected := filepath.Join(thisPath, "/fixtures/trivial.toml")

	result := findFile("trivial", []string{"./fixtures/", "/any/other/path"})
	assert.Equal(t, expected, result)
}

func Test_findFile_sliceFileAndDirFirstIf(t *testing.T) {
	inFilename := ""
	inDirNfile := []string{"$PWD/fixtures/nothing.toml"}

	thisPath, _ := filepath.Abs(".")
	expected := filepath.Join(thisPath, "/fixtures/nothing.toml")

	result := findFile(inFilename, inDirNfile)
	assert.Equal(t, expected, result)
}
