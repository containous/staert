package staert

import (
	"bytes"
	"fmt"
	"github.com/containous/flaeg"
	"reflect"
	"strings"
	"testing"
	"time"
)

//StructPtr : Struct with pointers
type StructPtr struct {
	PtrStruct1    *Struct1      `description:"Enable Struct1"`
	PtrStruct2    *Struct2      `description:"Enable Struct1"`
	DurationField time.Duration `description:"Duration Field"`
}

//Struct1 : Struct with pointer
type Struct1 struct {
	S1Int        int      `description:"Struct 1 Int"`
	S1String     string   `description:"Struct 1 String"`
	S1Bool       bool     `description:"Struct 1 Bool"`
	S1PtrStruct3 *Struct3 `description:"Enable Struct3"`
}

//Struct2 : trivial Struct
type Struct2 struct {
	S2Int64  int64  `description:"Struct 2 Int64"`
	S2String string `description:"Struct 2 String"`
	S2Bool   bool   `description:"Struct 2 Bool"`
}

//Struct3 : trivial Struct
type Struct3 struct {
	S3Float64 float64 `description:"Struct 3 float64"`
}

func TestFleagSourceNoArgs(t *testing.T) {
	//use buffer instead of stdout
	var b bytes.Buffer
	//Init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
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
	args := []string{}

	//Test
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
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)
	if err := s.getConfig(rootCmd); err != nil {
		t.Errorf("Error %s", err.Error())
	}

	//Check
	check := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
	}

	if !reflect.DeepEqual(rootCmd.Config, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, rootCmd.Config)
	}
}

func TestFleagSourcePtrUnderPtrArgs(t *testing.T) {
	//use buffer instead of stdout
	var b bytes.Buffer
	//Init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
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

	//Test
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
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)
	if err := s.getConfig(rootCmd); err != nil {
		t.Errorf("Error %s", err.Error())
	}
	//Check
	check := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
			S1PtrStruct3: &Struct3{
				S3Float64: 11.11,
			},
		},
		DurationField: time.Second,
	}

	if !reflect.DeepEqual(rootCmd.Config, check) {
		result, ok := rootCmd.Config.(*StructPtr)
		if ok {
			fmt.Printf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check.PtrStruct1, result.PtrStruct1)
		}

		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, rootCmd.Config)
	}
}

func TestFleagSourceFieldUnderPtrUnderPtrArgs(t *testing.T) {
	//use buffer instead of stdout
	var b bytes.Buffer
	//Init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
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

	//Test
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
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)
	if err := s.getConfig(rootCmd); err != nil {
		t.Errorf("Error %s", err.Error())
	}
	//Check
	check := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
			S1PtrStruct3: &Struct3{
				S3Float64: 55.55,
			},
		},
		DurationField: time.Second,
	}

	if !reflect.DeepEqual(rootCmd.Config, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, rootCmd.Config)
	}
}

func TestTomlSourceNothing(t *testing.T) {
	//use buffer instead of stdout
	var b bytes.Buffer
	//Init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
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

	//Test
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
	toml := NewTomlSource("nothing", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)

	if err := s.getConfig(rootCmd); err != nil {
		t.Errorf("Error %s", err.Error())
	}
	//Check
	check := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
	}

	if !reflect.DeepEqual(rootCmd.Config, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, rootCmd.Config)
	}
}

func TestTomlSourceTrivial(t *testing.T) {
	//use buffer instead of stdout
	var b bytes.Buffer
	//Init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
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

	//Test
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

	if err := s.getConfig(rootCmd); err != nil {
		t.Errorf("Error %s", err.Error())
	}
	//Check
	check := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    28,
			S1String: "S1StringInitConfig",
		},
		DurationField: 28 * time.Nanosecond,
	}

	if !reflect.DeepEqual(rootCmd.Config, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, rootCmd.Config)
	}
}

func TestTomlSourcePointer(t *testing.T) {
	//use buffer instead of stdout
	var b bytes.Buffer
	//Init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
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

	//Test
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
	toml := NewTomlSource("pointer", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)

	if err := s.getConfig(rootCmd); err != nil {
		t.Errorf("Error %s", err.Error())
	}
	//Check
	check := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		PtrStruct2: &Struct2{
			S2Int64:  22,
			S2String: "S2StringDefaultPointersConfig",
			S2Bool:   false,
		},
		DurationField: time.Second,
	}

	if !reflect.DeepEqual(rootCmd.Config, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, rootCmd.Config)
	}
}

func TestTomlSourcePointerUnderPointer(t *testing.T) {
	//use buffer instead of stdout
	var b bytes.Buffer
	//Init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
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

	//Test
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
	toml := NewTomlSource("pointerUnderPointer", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)

	if err := s.getConfig(rootCmd); err != nil {
		t.Errorf("Error %s", err.Error())
	}
	//Check
	check := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
			S1PtrStruct3: &Struct3{
				S3Float64: 11.11,
			},
		},
		DurationField: time.Second,
	}

	if !reflect.DeepEqual(rootCmd.Config, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, rootCmd.Config)
	}
}

func TestTomlSourceFieldUnderPointerUnderPointer(t *testing.T) {
	//use buffer instead of stdout
	var b bytes.Buffer
	//Init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
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

	//Test
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
	toml := NewTomlSource("fieldUnderPtrUnderPtr", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)

	if err := s.getConfig(rootCmd); err != nil {
		t.Errorf("Error %s", err.Error())
	}
	//Check
	check := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
			S1PtrStruct3: &Struct3{
				S3Float64: 28.28,
			},
		},
		DurationField: time.Second,
	}

	if !reflect.DeepEqual(rootCmd.Config, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, rootCmd.Config)
	}
}

func TestMergeTomlNothingFlaegNoArgs(t *testing.T) {
	//use buffer instead of stdout
	var b bytes.Buffer
	//Init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
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

	args := []string{}

	//Test
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
	toml := NewTomlSource("nothing", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)
	if err := s.getConfig(rootCmd); err != nil {
		t.Errorf("Error %s", err.Error())
	}

	//Check
	check := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
	}

	if !reflect.DeepEqual(rootCmd.Config, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, rootCmd.Config)
	}

}

func TestMergeTomlFieldUnderPointerUnderPointerFlaegNoArgs(t *testing.T) {
	//use buffer instead of stdout
	var b bytes.Buffer
	//Init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
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

	args := []string{}

	//Test
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
	toml := NewTomlSource("fieldUnderPtrUnderPtr", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)
	if err := s.getConfig(rootCmd); err != nil {
		t.Errorf("Error %s", err.Error())
	}

	//Check
	check := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
			S1PtrStruct3: &Struct3{
				S3Float64: 28.28,
			},
		},
		DurationField: time.Second,
	}

	if !reflect.DeepEqual(rootCmd.Config, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, rootCmd.Config)
	}

}

func TestMergeTomlTrivialFlaegOverwriteField(t *testing.T) {
	//use buffer instead of stdout
	var b bytes.Buffer
	//Init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
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

	//Test
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
	if err := s.getConfig(rootCmd); err != nil {
		t.Errorf("Error %s", err.Error())
	}
	//Check
	check := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    55,
			S1String: "S1StringInitConfig",
		},
		DurationField: 28 * time.Nanosecond,
	}

	if !reflect.DeepEqual(rootCmd.Config, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, rootCmd.Config)
	}

}

func TestMergeTomlPointerUnderPointerFlaegManyArgs(t *testing.T) {
	//use buffer instead of stdout
	var b bytes.Buffer
	//Init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
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

	//Test
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
	toml := NewTomlSource("pointerUnderPointer", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)
	if err := s.getConfig(rootCmd); err != nil {
		t.Errorf("Error %s", err.Error())
	}

	//Check
	check := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    55,
			S1String: "S1StringInitConfig",
			S1PtrStruct3: &Struct3{
				S3Float64: 11.11,
			},
		},
		PtrStruct2: &Struct2{
			S2Int64:  22,
			S2String: "S2StringFlaeg",
			S2Bool:   false,
		},
		DurationField: time.Second * 55,
	}
	if !reflect.DeepEqual(rootCmd.Config, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, rootCmd.Config)
	}

}

func TestMergeFlaegNoArgsTomlNothing(t *testing.T) {
	//use buffer instead of stdout
	var b bytes.Buffer
	//Init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
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

	args := []string{}

	//Test
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
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)
	toml := NewTomlSource("nothing", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)
	if err := s.getConfig(rootCmd); err != nil {
		t.Errorf("Error %s", err.Error())
	}

	//Check
	check := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
	}

	if !reflect.DeepEqual(rootCmd.Config, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, rootCmd.Config)
	}

}

func TestMergeFlaegFieldUnderPointerUnderPointerTomlNothing(t *testing.T) {
	//use buffer instead of stdout
	var b bytes.Buffer
	//Init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
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
	//Test
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
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)
	toml := NewTomlSource("nothing", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)
	if err := s.getConfig(rootCmd); err != nil {
		t.Errorf("Error %s", err.Error())
	}

	//Check
	check := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
			S1PtrStruct3: &Struct3{
				S3Float64: 55.55,
			},
		},
		DurationField: time.Second,
	}

	if !reflect.DeepEqual(rootCmd.Config, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, rootCmd.Config)
	}

}

func TestMergeFlaegManyArgsTomlOverwriteField(t *testing.T) {
	//use buffer instead of stdout
	var b bytes.Buffer
	//Init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
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
	//Test
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
	fs := flaeg.New(rootCmd, args)
	s.AddSource(fs)
	toml := NewTomlSource("trivial", []string{"./toml/", "/any/other/path"})
	s.AddSource(toml)
	if err := s.getConfig(rootCmd); err != nil {
		t.Errorf("Error %s", err.Error())
	}

	//Check
	check := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    28,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Nanosecond * 28,
		PtrStruct2: &Struct2{
			S2Int64:  22,
			S2String: "S2StringFlaeg",
			S2Bool:   false,
		},
	}

	if !reflect.DeepEqual(rootCmd.Config, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, rootCmd.Config)
	}

}

func TestRunFleagFieldUnderPtrUnderPtr1Command(t *testing.T) {
	//use buffer instead of stdout
	var b bytes.Buffer
	//Init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
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

	//Test
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
				DurationField: time.Second,
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
	if err := s.Run(); err != nil {
		t.Fatalf("Error %s", err.Error())
	}
	//check buffer
	checkB := `Run with config`
	if !strings.Contains(b.String(), checkB) {
		t.Fatalf("Error output doesn't contain %s,\ngot: %s", checkB, &b)
	}
}

//Version Config
type VersionConfig struct {
	Version string `short:"v" description:"Version"`
}

func TestRunFleagFieldUnderPtrUnderPtr2Command(t *testing.T) {
	//use buffer instead of stdout
	var b bytes.Buffer
	//Init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
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
	//init version config
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
				DurationField: time.Second,
			}

			if !reflect.DeepEqual(config, check) {
				return fmt.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, config)
			}
			return nil
		},
	}
	//vesion command
	versionCmd := &flaeg.Command{
		Name:        "version",
		Description: `Print version`,

		Config:                versionConfig,
		DefaultPointersConfig: versionConfig,
		//test in run
		Run: func() error {
			fmt.Fprintf(&b, "Version %s \n", versionConfig.Version)
			//CHECK
			if versionConfig.Version != "0.1" {
				return fmt.Errorf("expected 0.1 got %s", versionConfig.Version)
			}
			return nil

		},
	}
	//TEST
	s := NewStaert(rootCmd)
	fs := flaeg.New(rootCmd, args)
	fs.AddCommand(versionCmd)
	s.AddSource(fs)
	//check in command run func
	if err := s.Run(); err != nil {
		t.Fatalf("Error %s", err.Error())
	}
	//check buffer
	checkB := `Run with config`
	if !strings.Contains(b.String(), checkB) {
		t.Fatalf("Error output doesn't contain %s,\ngot: %s", checkB, &b)
	}
}

func TestRunFleagVersion2CommandCallVersion(t *testing.T) {
	//use buffer instead of stdout
	var b bytes.Buffer
	//Init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
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
	//init version config
	versionConfig := &VersionConfig{"0.1"}

	args := []string{
		"--toto",  //no effect
		"version", //call Command
		// "-v2.2beta",
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
				DurationField: time.Second,
			}

			if !reflect.DeepEqual(config, check) {
				return fmt.Errorf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, config)
			}
			return nil
		},
	}
	//vesion command
	versionCmd := &flaeg.Command{
		Name:        "version",
		Description: `Print version`,

		Config:                versionConfig,
		DefaultPointersConfig: versionConfig,
		//test in run
		Run: func() error {
			fmt.Fprintf(&b, "Version %s \n", versionConfig.Version)
			//CHECK
			if versionConfig.Version != "0.1" {
				return fmt.Errorf("expected 0.1 got %s", versionConfig.Version)
			}
			return nil

		},
	}
	//TEST
	s := NewStaert(rootCmd)
	fs := flaeg.New(rootCmd, args)
	fs.AddCommand(versionCmd)
	s.AddSource(fs)
	//check in command run func
	if err := s.Run(); err != nil {
		t.Fatalf("Error %s", err.Error())
	}
	//check buffer
	checkB := `Version 0.1`
	if !strings.Contains(b.String(), checkB) {
		t.Fatalf("Error output doesn't contain %s,\ngot: %s", checkB, &b)
	}
}

func TestRunMergeFlaegToml2CommmandCallRootCmd(t *testing.T) {
	//use buffer instead of stdout
	var b bytes.Buffer
	//Init
	config := &StructPtr{
		PtrStruct1: &Struct1{
			S1Int:    1,
			S1String: "S1StringInitConfig",
		},
		DurationField: time.Second,
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
	//init version config
	versionConfig := &VersionConfig{"0.1"}

	args := []string{
		"--ptrstruct1.s1int=55",
		"--durationfield=55s",
		"--ptrstruct2.s2string=S2StringFlaeg",
	}
	//Test
	rootCmd := &flaeg.Command{
		Name:                  "test",
		Description:           "description test",
		Config:                config,
		DefaultPointersConfig: defaultPointersConfig,
		Run: func() error {
			fmt.Fprintf(&b, "Run with config :\n%+v\n", config)
			//Check
			check := &StructPtr{
				PtrStruct1: &Struct1{
					S1Int:    28,
					S1String: "S1StringInitConfig",
				},
				DurationField: time.Nanosecond * 28,
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
	//vesion command
	versionCmd := &flaeg.Command{
		Name:        "version",
		Description: `Print version`,

		Config:                versionConfig,
		DefaultPointersConfig: versionConfig,
		//test in run
		Run: func() error {
			fmt.Fprintf(&b, "Version %s \n", versionConfig.Version)
			//CHECK
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
	//check in command run func
	if err := s.Run(); err != nil {
		t.Fatalf("Error %s", err.Error())
	}
	//check buffer
	checkB := `Run with config :`
	if !strings.Contains(b.String(), checkB) {
		t.Fatalf("Error output doesn't contain %s,\ngot: %s", checkB, &b)
	}

}
