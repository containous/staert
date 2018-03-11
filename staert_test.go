package staert

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/containous/flaeg"
	"github.com/containous/flaeg/parse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// StructPtr : Struct with pointers
type StructPtr struct {
	PtrStruct1    *Struct1       `description:"Enable Struct1"`
	PtrStruct2    *Struct2       `description:"Enable Struct1"`
	DurationField parse.Duration `description:"Duration Field"`
}

// Struct1 : Struct with pointer
type Struct1 struct {
	S1Int        int      `description:"Struct 1 Int"`
	S1String     string   `description:"Struct 1 String"`
	S1Bool       bool     `description:"Struct 1 Bool"`
	S1PtrStruct3 *Struct3 `description:"Enable Struct3"`
}

// Struct2 : trivial Struct
type Struct2 struct {
	S2Int64  int64  `description:"Struct 2 Int64"`
	S2String string `description:"Struct 2 String"`
	S2Bool   bool   `description:"Struct 2 Bool"`
}

// Struct3 : trivial Struct
type Struct3 struct {
	S3Float64 float64 `description:"Struct 3 float64"`
}

// Version Config
type VersionConfig struct {
	Version string `short:"v" description:"Version"`
}

type SliceStr []string

type StructPtrCustom struct {
	PtrCustom *StructCustomParser `description:"Ptr on StructCustomParser"`
}

type StructCustomParser struct {
	CustomField SliceStr `description:"CustomField which requires custom parser"`
}

func defaultPointersConfig() *StructPtr {
	return &StructPtr{
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
}

func Test_parseConfigAllSources_flaegSourceNoArgs(t *testing.T) {
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
		DefaultPointersConfig: defaultPointersConfig(),
		Run: func() error {
			return nil
		},
	}

	s := NewStaert(rootCmd)

	var args []string
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)

	err := s.parseConfigAllSources(rootCmd)
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

func Test_parseConfigAllSources_flaegSourcePtrUnderPtrArgs(t *testing.T) {
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	args := []string{
		"--ptrstruct1.s1ptrstruct3",
	}

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig(),
		Run: func() error {
			return nil
		},
	}

	s := NewStaert(rootCmd)

	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)

	err := s.parseConfigAllSources(rootCmd)
	require.NoError(t, err)

	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
			S1PtrStruct3: &Struct3{
				S3Float64: 11.11,
			},
		},
		DurationField: parse.Duration(time.Second),
	}

	assert.Exactly(t, expected, rootCmd.Config)
}

func Test_parseConfigAllSources_flaegSourceFieldUnderPtrUnderPtrArgs(t *testing.T) {
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	args := []string{
		"--ptrstruct1.s1ptrstruct3.s3float64=55.55",
	}

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig(),
		Run: func() error {
			return nil
		},
	}

	s := NewStaert(rootCmd)

	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)

	err := s.parseConfigAllSources(rootCmd)
	require.NoError(t, err)

	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
			S1PtrStruct3: &Struct3{
				S3Float64: 55.55,
			},
		},
		DurationField: parse.Duration(time.Second),
	}

	assert.Exactly(t, expected, rootCmd.Config)
}

func Test_parseConfigAllSources_mergeFlaegWithoutArgsAndEmptyToml(t *testing.T) {
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	var args []string

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig(),
		Run: func() error {
			return nil
		},
	}

	s := NewStaert(rootCmd)

	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)

	toml := NewTomlSource("nothing", []string{"./fixtures/", "/any/other/path"})
	s.AddSource(toml)

	err := s.parseConfigAllSources(rootCmd)
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

func Test_parseConfigAllSources_mergeFlaegFieldUnderPointerUnderPointerAndEmptyToml(t *testing.T) {
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	args := []string{
		"--ptrstruct1.s1ptrstruct3.s3float64=55.55",
	}

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig(),
		Run: func() error {
			return nil
		},
	}

	s := NewStaert(rootCmd)

	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)

	toml := NewTomlSource("nothing", []string{"./fixtures/", "/any/other/path"})
	s.AddSource(toml)

	err := s.parseConfigAllSources(rootCmd)
	require.NoError(t, err)

	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
			S1PtrStruct3: &Struct3{
				S3Float64: 55.55,
			},
		},
		DurationField: parse.Duration(time.Second),
	}

	assert.Exactly(t, expected, rootCmd.Config)
}

func Test_parseConfigAllSources_mergeOverrideFlaegArgsAndTomlField(t *testing.T) {
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	args := []string{
		"--ptrstruct1.s1int=55",
		"--durationfield=55s",
		"--ptrstruct2.s2string=S2StringFlaeg",
	}

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig(),
		Run: func() error {
			return nil
		},
	}

	s := NewStaert(rootCmd)

	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)

	toml := NewTomlSource("trivial", []string{"./fixtures/", "/any/other/path"})
	s.AddSource(toml)

	err := s.parseConfigAllSources(rootCmd)
	require.NoError(t, err)

	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    28,
			S1String: "S1StringDefaultPointersConfig",
			S1Bool:   true,
		},
		DurationField: parse.Duration(28 * time.Second),
		PtrStruct2: &Struct2{
			S2Int64:  22,
			S2String: "S2StringFlaeg",
			S2Bool:   false,
		},
	}

	assert.Exactly(t, expected, rootCmd.Config)
}

func Test_parseConfigAllSources_mergeEmptyTomlAndFlaegWithoutArgs(t *testing.T) {
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	var args []string

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig(),
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

	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	assert.Exactly(t, expected, rootCmd.Config)
}

func Test_parseConfigAllSources_mergeTomlFieldUnderPointerUnderPointerAndFlaegWithoutArgs(t *testing.T) {
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	var args []string

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig(),
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

func Test_parseConfigAllSources_mergeTomlTrivialAndFlaegOverwriteField(t *testing.T) {
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	args := []string{"--ptrstruct1.s1int=55"}

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig(),
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

func Test_parseConfigAllSources_mergeTomlPointerUnderPointerAndFlaegArgs(t *testing.T) {
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	args := []string{
		"--ptrstruct1.s1int=55",
		"--durationfield=55s",
		"--ptrstruct2.s2string=S2StringFlaeg",
	}

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig(),
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

func TestRun_withoutLoadConfig(t *testing.T) {
	// use buffer instead of stdout
	var b bytes.Buffer

	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	args := []string{
		"--ptrstruct2",
	}

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig(),
		Run: func() error {
			fmt.Fprintf(&b, "Run with config :\n%+v\n", config)
			return nil
		},
	}

	s := NewStaert(rootCmd)

	toml := NewTomlSource("trivial", []string{"./fixtures/", "/any/other/path"})
	s.AddSource(toml)

	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)

	// s.LoadConfig() IS MISSING
	err := s.Run()
	require.NoError(t, err)

	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	assert.Exactly(t, expected, rootCmd.Config)
	assert.Contains(t, b.String(), "Run with config")
}

func TestRun_flaegFieldUnderPtrUnderPtr1(t *testing.T) {
	// use buffer instead of stdout
	var b bytes.Buffer

	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	args := []string{
		"--ptrstruct1.s1ptrstruct3.s3float64=55.55",
	}

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig(),
		Run: func() error {
			fmt.Fprintf(&b, "Run with config :\n%+v\n", config)
			check := &StructPtr{
				PtrStruct1: &Struct1{
					S1Int:    1,
					S1String: "S1StringInitConfig",
					S1PtrStruct3: &Struct3{
						S3Float64: 55.55,
					},
				},
				DurationField: parse.Duration(time.Second),
			}

			if !reflect.DeepEqual(config, check) {
				return fmt.Errorf("expected\t: %+v\ngot\t\t\t: %+v", check, config)
			}
			return nil
		},
	}

	s := NewStaert(rootCmd)
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)

	_, err := s.LoadConfig()
	require.NoError(t, err)

	err = s.Run()
	require.NoError(t, err)

	assert.Contains(t, b.String(), "Run with config")
}

func TestRun_flaegFieldUnderPtrUnderPtr2(t *testing.T) {
	// use buffer instead of stdout
	var b bytes.Buffer

	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	// init version config
	versionConfig := &VersionConfig{"0.1"}

	args := []string{
		"--ptrstruct1.s1ptrstruct3.s3float64=55.55",
	}

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig(),
		Run: func() error {
			fmt.Fprintf(&b, "Run with config :\n%+v\n", config)
			check := &StructPtr{
				PtrStruct1: &Struct1{
					S1Int:    1,
					S1String: "S1StringInitConfig",
					S1PtrStruct3: &Struct3{
						S3Float64: 55.55,
					},
				},
				DurationField: parse.Duration(time.Second),
			}

			if !reflect.DeepEqual(config, check) {
				return fmt.Errorf("expected\t: %+v\ngot\t\t\t: %+v", check, config)
			}
			return nil
		},
	}

	// version command
	versionCmd := &flaeg.Command{
		Name:                  "version",
		Description:           `Print version`,
		Config:                versionConfig,
		DefaultPointersConfig: versionConfig,
		// Test in run
		Run: func() error {
			fmt.Fprintf(&b, "Version %s \n", versionConfig.Version)
			// check
			if versionConfig.Version != "0.1" {
				return fmt.Errorf("expected 0.1 got %s", versionConfig.Version)
			}

			return nil
		},
	}

	s := NewStaert(rootCmd)

	fs := flaeg.New(rootCmd, args)
	fs.AddCommand(versionCmd)
	s.AddSource(fs)

	// check in command run func
	_, err := s.LoadConfig()
	require.NoError(t, err)

	err = s.Run()
	require.NoError(t, err)

	assert.Contains(t, b.String(), "Run with config")
}

func TestRun_flaegVersion2CallVersion(t *testing.T) {
	// use buffer instead of stdout
	var b bytes.Buffer

	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	// init version config
	versionConfig := &VersionConfig{"0.1"}

	args := []string{
		// "--toto",  // it now has effet
		"version", // call Command
		"-v2.2beta",
	}

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig(),
		Run: func() error {
			fmt.Fprintf(&b, "Run with config :\n%+v\n", config)
			check := &StructPtr{
				PtrStruct1: &Struct1{
					S1Int:    1,
					S1String: "S1StringInitConfig",
					S1PtrStruct3: &Struct3{
						S3Float64: 55.55,
					},
				},
				DurationField: parse.Duration(time.Second),
			}

			if !reflect.DeepEqual(config, check) {
				return fmt.Errorf("expected\t: %+v\ngot\t\t\t: %+v", check, config)
			}
			return nil
		},
	}

	// version command
	versionCmd := &flaeg.Command{
		Name:                  "version",
		Description:           `Print version`,
		Config:                versionConfig,
		DefaultPointersConfig: versionConfig,
		// test in run
		Run: func() error {
			fmt.Fprintf(&b, "Version %s \n", versionConfig.Version)
			// expected
			if versionConfig.Version != "2.2beta" {
				return fmt.Errorf("expected 2.2beta got %s", versionConfig.Version)
			}
			return nil

		},
	}

	s := NewStaert(rootCmd)

	fs := flaeg.New(rootCmd, args)
	fs.AddCommand(versionCmd)
	s.AddSource(fs)

	// expected in command run func
	_, err := s.LoadConfig()
	require.NoError(t, err)

	err = s.Run()
	require.NoError(t, err)

	assert.Contains(t, b.String(), "Version 2.2beta")
}

func TestRun_mergeFlaegToml2CallRootCmd(t *testing.T) {
	// use buffer instead of stdout
	var b bytes.Buffer

	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	// init version config
	versionConfig := &VersionConfig{"0.1"}

	args := []string{
		"--ptrstruct1.s1int=55",
		"--durationfield=55s",
		"--ptrstruct2.s2string=S2StringFlaeg",
	}

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig(),
		Run: func() error {
			fmt.Fprintf(&b, "Run with config :\n%+v\n", config)
			// expected
			expected := &StructPtr{
				PtrStruct1: &Struct1{
					S1Int:    28,
					S1String: "S1StringDefaultPointersConfig",
					S1Bool:   true,
				},
				DurationField: parse.Duration(28 * time.Second),
				PtrStruct2: &Struct2{
					S2Int64:  22,
					S2String: "S2StringFlaeg",
					S2Bool:   false,
				},
			}

			if !reflect.DeepEqual(config, expected) {
				return fmt.Errorf("expected\t: %+v\ngot\t\t\t: %+v", expected, config)
			}
			return nil
		},
	}

	// version command
	versionCmd := &flaeg.Command{
		Name:        "version",
		Description: `Print version`,

		Config:                versionConfig,
		DefaultPointersConfig: versionConfig,
		// Test in run
		Run: func() error {
			fmt.Fprintf(&b, "Version %s \n", versionConfig.Version)
			// expected
			if versionConfig.Version != "0.1" {
				return fmt.Errorf("expected 0.1 got %s", versionConfig.Version)
			}
			return nil

		},
	}

	s := NewStaert(rootCmd)

	fs := flaeg.New(rootCmd, args)
	fs.AddCommand(versionCmd)
	s.AddSource(fs)

	toml := NewTomlSource("trivial", []string{"./fixtures/", "/any/other/path"})
	s.AddSource(toml)

	// check in command run func
	_, err := s.LoadConfig()
	require.NoError(t, err)

	err = s.Run()
	require.NoError(t, err)

	assert.Contains(t, b.String(), "Run with config :")
}

func TestRun_flaegTomlSubCommandParseAllSources(t *testing.T) {
	// use buffer instead of stdout
	var b bytes.Buffer

	args := []string{
		"subcmd",
		"--Vstring=toto",
	}

	config := &struct {
		Vstring string `description:"string field"`
		Vint    int    `description:"int field"`
	}{
		Vstring: "tata",
		Vint:    -15,
	}

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: config,
		Run: func() error {
			fmt.Fprintln(&b, "rootCmd")
			fmt.Fprintf(&b, "run with config : %+v\n", config)
			return nil
		},
	}

	subCmd := &flaeg.Command{
		Name:                  "subcmd",
		Description:           "description subcmd",
		Config:                config,
		DefaultPointersConfig: config,
		Run: func() error {
			fmt.Fprintln(&b, "subcmd")
			fmt.Fprintf(&b, "run with config : %+v\n", config)
			return nil
		},
		Metadata: map[string]string{
			"parseAllSources": "true",
		},
	}

	s := NewStaert(rootCmd)

	toml := NewTomlSource("subcmd", []string{"./fixtures/", "/any/other/path"})
	s.AddSource(toml)

	fs := flaeg.New(rootCmd, args)
	fs.AddCommand(subCmd)
	s.AddSource(fs)

	_, err := s.LoadConfig()
	require.NoError(t, err)

	err = s.Run()
	require.NoError(t, err)

	assert.Contains(t, b.String(), "subcmd")
	assert.Contains(t, b.String(), "Vstring:toto")
	assert.Contains(t, b.String(), "Vint:777")
}

func TestLoadConfig_tomlMissingCustomParser(t *testing.T) {
	config := &StructPtrCustom{}

	defaultPointersConfig := &StructPtrCustom{&StructCustomParser{SliceStr{"str1", "str2"}}}

	command := &flaeg.Command{
		Name:                  "MissingCustomParser",
		Description:           "This is an example of description",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig,
		Run: func() error {
			// expected
			check := &StructPtrCustom{&StructCustomParser{SliceStr{"str1", "str2"}}}
			if !reflect.DeepEqual(config, check) {
				return fmt.Errorf("Expected %+v\ngot %+v", check.PtrCustom, config.PtrCustom)
			}
			return nil
		},
	}

	s := NewStaert(command)
	toml := NewTomlSource("missingCustomParser", []string{"fixtures"})
	s.AddSource(toml)

	_, err := s.LoadConfig()
	require.NoError(t, err)

	err = s.Run()
	require.NoError(t, err)

	expected := &StructPtrCustom{&StructCustomParser{SliceStr{"str1", "str2"}}}
	assert.Exactly(t, expected, config)
}

func TestLoadConfig_flaegTomlSubCommandParseAllSourcesShouldError(t *testing.T) {
	// use buffer instead of stdout
	var b bytes.Buffer

	args := []string{
		"subcmd",
		"--Vstring=toto",
	}

	config := &struct {
		Vstring string `description:"string field"`
		Vint    int    `description:"int field"`
	}{
		Vstring: "tata",
		Vint:    -15,
	}

	config2 := &struct {
		Vstring int `description:"int field"` // TO check error
		Vint    int `description:"int field"`
	}{
		Vstring: -1,
		Vint:    -15,
	}

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: config,
		Run: func() error {
			fmt.Fprintln(&b, "rootCmd")
			fmt.Fprintf(&b, "run with config : %+v\n", config)
			return nil
		},
	}

	subCmd := &flaeg.Command{
		Name:                  "subcmd",
		Description:           "description subcmd",
		Config:                config2,
		DefaultPointersConfig: config2,
		Run: func() error {
			fmt.Fprintln(&b, "subcmd")
			fmt.Fprintf(&b, "run with config : %+v\n", config)
			return nil
		},
		Metadata: map[string]string{
			"parseAllSources": "true",
		},
	}

	s := NewStaert(rootCmd)
	toml := NewTomlSource("subcmd", []string{"./fixtures/", "/any/other/path"})
	s.AddSource(toml)
	fs := flaeg.New(rootCmd, args)
	fs.AddCommand(subCmd)
	s.AddSource(fs)

	_, err := s.LoadConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Config type doesn't match with root command config type.")
}
