package staert

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/containous/flaeg"
	"github.com/containous/flaeg/parse"
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

func TestFlaegSourceNoArgs(t *testing.T) {
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

	var args []string
	s := NewStaert(rootCmd)
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)

	if err := s.parseConfigAllSources(rootCmd); err != nil {
		t.Errorf("Error %v", err)
	}

	// expected
	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	if !reflect.DeepEqual(rootCmd.Config, expected) {
		t.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected, rootCmd.Config)
	}
}

func TestFlaegSourcePtrUnderPtrArgs(t *testing.T) {
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
		"--ptrstruct1.s1ptrstruct3",
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
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)

	if err := s.parseConfigAllSources(rootCmd); err != nil {
		t.Errorf("Error %v", err)
	}

	// expected
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

	if !reflect.DeepEqual(rootCmd.Config, expected) {
		result, ok := rootCmd.Config.(*StructPtr)
		if ok {
			fmt.Printf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected.PtrStruct1, result.PtrStruct1)
		}

		t.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected, rootCmd.Config)
	}
}

func TestFlaegSourceFieldUnderPtrUnderPtrArgs(t *testing.T) {
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
		"--ptrstruct1.s1ptrstruct3.s3float64=55.55",
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
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)

	if err := s.parseConfigAllSources(rootCmd); err != nil {
		t.Errorf("Error %v", err)
	}

	// expected
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

	if !reflect.DeepEqual(rootCmd.Config, expected) {
		t.Errorf("expected\t: %+v\ngot\t\t\t: %+v\n", expected, rootCmd.Config)
	}
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
	toml := NewTomlSource("nothing", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)

	if err := s.parseConfigAllSources(rootCmd); err != nil {
		t.Errorf("Error %v", err)
	}

	// expected
	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	if !reflect.DeepEqual(rootCmd.Config, expected) {
		t.Errorf("expected\t: %+v\ngot\t\t\t: %+v\n", expected, rootCmd.Config)
	}
}

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
	toml := NewTomlSource("trivial", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)

	if err := s.parseConfigAllSources(rootCmd); err != nil {
		t.Errorf("Error %v", err)
	}

	// expected
	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    28,
			S1String: "S1StringDefaultPointersConfig",
			S1Bool:   true,
		},
		DurationField: parse.Duration(28 * time.Second),
	}

	if !reflect.DeepEqual(rootCmd.Config, expected) {
		t.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected, rootCmd.Config)
		t.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected.PtrStruct1, config.PtrStruct1)
	}
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
	toml := NewTomlSource("pointer", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)

	if err := s.parseConfigAllSources(rootCmd); err != nil {
		t.Errorf("Error %v", err)
	}

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

	if !reflect.DeepEqual(rootCmd.Config, expected) {
		t.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected, rootCmd.Config)
	}
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
	toml := NewTomlSource("fieldUnderPointer", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)

	if err := s.parseConfigAllSources(rootCmd); err != nil {
		t.Errorf("Error %v", err)
	}

	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    11,
			S1String: "S1StringDefaultPointersConfig",
			S1Bool:   true,
		},
		DurationField: parse.Duration(42 * time.Second),
	}

	if !reflect.DeepEqual(rootCmd.Config, expected) {
		t.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected, config)
		t.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected.PtrStruct1, config.PtrStruct1)
	}
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
	toml := NewTomlSource("pointerUnderPointer", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)

	if err := s.parseConfigAllSources(rootCmd); err != nil {
		t.Errorf("Error %v", err)
	}

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

	if !reflect.DeepEqual(rootCmd.Config, expected) {
		t.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected, rootCmd.Config)
		t.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected.PtrStruct1, config.PtrStruct1)
		t.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected.PtrStruct1.S1PtrStruct3, config.PtrStruct1.S1PtrStruct3)
	}
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
	toml := NewTomlSource("fieldUnderPtrUnderPtr", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)

	if err := s.parseConfigAllSources(rootCmd); err != nil {
		t.Errorf("Error %v", err)
	}

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

	if !reflect.DeepEqual(rootCmd.Config, expected) {
		t.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected, rootCmd.Config)
	}
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
	toml := NewTomlSource("nothing", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)

	if err := s.parseConfigAllSources(rootCmd); err != nil {
		t.Errorf("Error %v", err)
	}

	// expected
	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	if !reflect.DeepEqual(rootCmd.Config, expected) {
		t.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected, rootCmd.Config)
	}
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
	toml := NewTomlSource("fieldUnderPtrUnderPtr", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)

	if err := s.parseConfigAllSources(rootCmd); err != nil {
		t.Errorf("Error %v", err)
	}

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

	if !reflect.DeepEqual(rootCmd.Config, expected) {
		t.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected, rootCmd.Config)
	}
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
	toml := NewTomlSource("trivial", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)

	if err := s.parseConfigAllSources(rootCmd); err != nil {
		t.Errorf("Error %v", err)
	}

	// expected
	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    55,
			S1String: "S1StringDefaultPointersConfig",
			S1Bool:   true,
		},
		DurationField: parse.Duration(28 * time.Second),
	}

	if !reflect.DeepEqual(rootCmd.Config, expected) {
		t.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected, rootCmd.Config)
		t.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected.PtrStruct1, config.PtrStruct1)

	}
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
	toml := NewTomlSource("pointerUnderPointer", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)

	if err := s.parseConfigAllSources(rootCmd); err != nil {
		t.Errorf("Error %v", err)
	}

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

	if !reflect.DeepEqual(rootCmd.Config, expected) {
		t.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected, rootCmd.Config)
	}
}

func TestMergeFlaegNoArgsTomlNothing(t *testing.T) {
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
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)
	toml := NewTomlSource("nothing", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)

	if err := s.parseConfigAllSources(rootCmd); err != nil {
		t.Errorf("Error %v", err)
	}

	// expected
	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	if !reflect.DeepEqual(rootCmd.Config, expected) {
		t.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected, rootCmd.Config)
	}
}

func TestMergeFlaegFieldUnderPointerUnderPointerTomlNothing(t *testing.T) {
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
		"--ptrstruct1.s1ptrstruct3.s3float64=55.55",
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
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)
	toml := NewTomlSource("nothing", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)

	if err := s.parseConfigAllSources(rootCmd); err != nil {
		t.Errorf("Error %v", err)
	}

	// expected
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

	if !reflect.DeepEqual(rootCmd.Config, expected) {
		t.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected, rootCmd.Config)
	}
}

func TestMergeFlaegManyArgsTomlOverwriteField(t *testing.T) {
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
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)
	toml := NewTomlSource("trivial", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)

	if err := s.parseConfigAllSources(rootCmd); err != nil {
		t.Errorf("Error %v", err)
	}

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

	if !reflect.DeepEqual(rootCmd.Config, expected) {
		t.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected, rootCmd.Config)
	}
}

func TestRunFlaegFieldUnderPtrUnderPtr1Command(t *testing.T) {
	// use buffer instead of stdout
	var b bytes.Buffer

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
		"--ptrstruct1.s1ptrstruct3.s3float64=55.55",
	}

	// Test
	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig,
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
				return fmt.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, config)
			}
			return nil
		},
	}

	s := NewStaert(rootCmd)
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)
	_, err := s.LoadConfig()
	if err != nil {
		t.Errorf("Error %v", err)
	}

	if err := s.Run(); err != nil {
		t.Errorf("Error %v", err)
	}

	// expected buffer
	expectedOutput := `Run with config`
	if !strings.Contains(b.String(), expectedOutput) {
		t.Errorf("Error output doesn't contain %s,\ngot: %s", expectedOutput, &b)
	}
}

// Version Config
type VersionConfig struct {
	Version string `short:"v" description:"Version"`
}

func TestRunFlaegFieldUnderPtrUnderPtr2Command(t *testing.T) {
	// use buffer instead of stdout
	var b bytes.Buffer

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

	// init version config
	versionConfig := &VersionConfig{"0.1"}

	args := []string{
		"--ptrstruct1.s1ptrstruct3.s3float64=55.55",
	}

	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig,
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
				return fmt.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, config)
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

	// test
	s := NewStaert(rootCmd)
	fs := flaeg.New(rootCmd, args)
	fs.AddCommand(versionCmd)
	s.AddSource(fs)

	// check in command run func
	_, err := s.LoadConfig()
	if err != nil {
		t.Errorf("Error %v", err)
	}

	if err := s.Run(); err != nil {
		t.Errorf("Error %v", err)
	}

	expectedOutput := `Run with config`
	if !strings.Contains(b.String(), expectedOutput) {
		t.Errorf("Error output doesn't contain %s,\ngot: %s", expectedOutput, &b)
	}
}

func TestRunFlaegVersion2CommandCallVersion(t *testing.T) {
	// use buffer instead of stdout
	var b bytes.Buffer

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
		DefaultPointersConfig: defaultPointersConfig,
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
				return fmt.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, config)
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

	// Test
	s := NewStaert(rootCmd)
	fs := flaeg.New(rootCmd, args)
	fs.AddCommand(versionCmd)
	s.AddSource(fs)

	// expected in command run func
	_, err := s.LoadConfig()
	if err != nil {
		t.Errorf("Error %v", err)
	}

	if err := s.Run(); err != nil {
		t.Errorf("Error %v", err)
	}

	expectedOutput := `Version 2.2beta`
	if !strings.Contains(b.String(), expectedOutput) {
		t.Errorf("Error output doesn't contain %s,\ngot: %s", expectedOutput, &b)
	}
}

func TestRunMergeFlaegToml2CommmandCallRootCmd(t *testing.T) {
	// use buffer instead of stdout
	var b bytes.Buffer

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

	// init version config
	versionConfig := &VersionConfig{"0.1"}

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
			fmt.Fprintf(&b, "Run with config :\n%+v\n", config)
			// expected
			check := &StructPtr{
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

			if !reflect.DeepEqual(config, check) {
				return fmt.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, config)
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
	toml := NewTomlSource("trivial", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)

	// check in command run func
	_, err := s.LoadConfig()
	if err != nil {
		t.Errorf("Error %v", err)
	}

	if err := s.Run(); err != nil {
		t.Errorf("Error %v", err)
	}

	expectedOutput := `Run with config :`
	if !strings.Contains(b.String(), expectedOutput) {
		t.Errorf("Error output doesn't contain %s,\ngot: %s", expectedOutput, &b)
	}
}

func TestTomlSourceErrorFileNotFound(t *testing.T) {
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

	expectedCmd := *rootCmd
	s := NewStaert(rootCmd)
	toml := NewTomlSource("nothing", []string{"../path", "/any/other/path"})
	s.AddSource(toml)

	// expected
	if err := s.parseConfigAllSources(rootCmd); err != nil {
		t.Errorf("No Error expected\nGot Error : %s", err)
	}

	if !reflect.DeepEqual(expectedCmd.Config, rootCmd.Config) {
		t.Errorf("Expected %+v \nGot %+v", expectedCmd.Config, rootCmd.Config)
	}

	if !reflect.DeepEqual(expectedCmd.DefaultPointersConfig, rootCmd.DefaultPointersConfig) {
		t.Errorf("Expected %+v \nGot %+v", expectedCmd.DefaultPointersConfig, rootCmd.DefaultPointersConfig)
	}
}

func TestPreprocessDir(t *testing.T) {
	thisPath, err := filepath.Abs(".")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	checkMap := map[string]string{
		".":                   thisPath,
		"dir1/dir2":           thisPath + "/dir1/dir2",
		"/etc/test":           "/etc/test",
		"/etc/dir1/file1.ext": "/etc/dir1/file1.ext",
	}

	for in, check := range checkMap {
		out, err := preprocessDir(in)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		//always check against the absolute path
		checkAbs, _ := filepath.Abs(check)

		if checkAbs != out {
			t.Errorf("input %s\nexpected %s\n got %s", in, checkAbs, out)
		}
	}
}

func TestPreprocessDirEnvVariablesExpansions(t *testing.T) {
	expected, _ := filepath.Abs("/some/path/my/path")
	in := "$TEST_ENV_VARIABLE/my/path"
	os.Setenv("TEST_ENV_VARIABLE", "/some/path")

	out, err := preprocessDir(in)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if out != expected {
		t.Errorf("input %s\nexpected %s\n got %s", in, expected, out)
	}
}

func TestFindFile(t *testing.T) {
	result := findFile("nothing", []string{"", "$HOME/test", "toml"})

	// expected
	thisPath, err := filepath.Abs(".")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := filepath.Join(thisPath, "toml", "nothing.toml")
	if result != expected {
		t.Errorf("Expected %s\ngot %s", expected, result)
	}
}

type SliceStr []string

type StructPtrCustom struct {
	PtrCustom *StructCustomParser `description:"Ptr on StructCustomParser"`
}

type StructCustomParser struct {
	CustomField SliceStr `description:"CustomField which requires custom parser"`
}

func TestTomlMissingCustomParser(t *testing.T) {
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
	toml := NewTomlSource("missingCustomParser", []string{"toml"})
	s.AddSource(toml)

	_, err := s.LoadConfig()
	if err != nil {
		t.Errorf("Error %v", err)
	}

	if err := s.Run(); err != nil {
		t.Errorf("Error :%s", err)
	}

	// expected
	expected := &StructPtrCustom{&StructCustomParser{SliceStr{"str1", "str2"}}}
	if !reflect.DeepEqual(config, expected) {
		t.Errorf("Expected %+v\ngot %+v", expected.PtrCustom, config.PtrCustom)
	}
}

func TestFindFileSliceFileAndDirLastIf(t *testing.T) {
	thisPath, _ := filepath.Abs(".")

	expected := filepath.Join(thisPath, "/toml/trivial.toml")
	result := findFile("trivial", []string{"./toml/", "/any/other/path"})

	if result != expected {
		t.Errorf("Expected %s\nGot %s", expected, result)
	}
}

func TestFindFileSliceFileAndDirFirstIf(t *testing.T) {
	inFilename := ""
	inDirNfile := []string{"$PWD/toml/nothing.toml"}

	thisPath, _ := filepath.Abs(".")
	expected := filepath.Join(thisPath, "/toml/nothing.toml")
	result := findFile(inFilename, inDirNfile)

	if result != expected {
		t.Errorf("Expected %s\nGot %s", expected, result)
	}
}

func TestRunWithoutLoadConfig(t *testing.T) {
	// use buffer instead of stdout
	var b bytes.Buffer

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
		"--ptrstruct2",
	}

	// Test
	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig,
		Run: func() error {
			fmt.Fprintf(&b, "Run with config :\n%+v\n", config)
			return nil
		},
	}

	s := NewStaert(rootCmd)
	toml := NewTomlSource("trivial", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)
	// s.LoadConfig() IS MISSING
	s.Run()

	expected := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: parse.Duration(time.Second),
	}

	if !reflect.DeepEqual(rootCmd.Config, expected) {
		t.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", expected, rootCmd.Config)
	}

	expectedOutput := `Run with config`
	if !strings.Contains(b.String(), expectedOutput) {
		t.Errorf("Error output doesn't contain %s,\ngot: %s", expectedOutput, &b)
	}
}

func TestFlaegTomlSubCommandParseAllSources(t *testing.T) {
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
	toml := NewTomlSource("subcmd", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)
	fs := flaeg.New(rootCmd, args)
	fs.AddCommand(subCmd)
	s.AddSource(fs)

	_, err := s.LoadConfig()
	if err != nil {
		t.Errorf("Error %v", err)
	}

	if err = s.Run(); err != nil {
		t.Errorf("Error %v", err)
	}

	// Test
	if !strings.Contains(b.String(), "subcmd") ||
		!strings.Contains(b.String(), "Vstring:toto") ||
		!strings.Contains(b.String(), "Vint:777") {
		t.Errorf("expected: subcmd, Vstring = toto, Vint = 777\n got %s", b.String())
	}
}

func TestFlaegTomlSubCommandParseAllSourcesShouldError(t *testing.T) {
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
	toml := NewTomlSource("subcmd", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)
	fs := flaeg.New(rootCmd, args)
	fs.AddCommand(subCmd)
	s.AddSource(fs)

	_, err := s.LoadConfig()

	errExp := "Config type doesn't match with root command config type."
	if err == nil || !strings.Contains(err.Error(), errExp) {
		t.Errorf("Experted error %s\n got : %s", errExp, err)
	}
}
