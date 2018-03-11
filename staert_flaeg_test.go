package staert

import (
	"testing"
	"time"

	"github.com/containous/flaeg"
	"github.com/containous/flaeg/parse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
