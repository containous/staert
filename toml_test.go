package staert

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
