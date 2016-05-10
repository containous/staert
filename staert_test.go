package staert

import (
	"fmt"
	"reflect"
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
	s := NewStaert(config, defaultPointersConfig)
	fs := NewFlaegSource(args, nil)
	s.Add(fs)
	result, err := s.GetConfig()
	if err != nil {
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

	//Type assertions
	resultStructPtr, ok := result.(*StructPtr)
	if !ok {
		t.Fatalf("Cannot convert the config into Configuration")
	}
	if !reflect.DeepEqual(resultStructPtr, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, resultStructPtr)
	}

}

func TestFleagSourcePtrUnderPtrArgs(t *testing.T) {
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
	s := NewStaert(config, defaultPointersConfig)
	fs := NewFlaegSource(args, nil)
	s.Add(fs)
	result, err := s.GetConfig()
	if err != nil {
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

	//Type assertions
	resultStructPtr, ok := result.(*StructPtr)
	if !ok {
		t.Fatalf("Cannot convert the result into StructPtr")
	}

	if !reflect.DeepEqual(resultStructPtr, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, resultStructPtr)
	}
}

func TestFleagSourceFieldUnderPtrUnderPtrArgs(t *testing.T) {
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
	s := NewStaert(config, defaultPointersConfig)
	fs := NewFlaegSource(args, nil)
	s.Add(fs)
	result, err := s.GetConfig()
	if err != nil {
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

	//Type assertions
	resultStructPtr, ok := result.(*StructPtr)
	if !ok {
		t.Fatalf("Cannot convert the result into StructPtr")
	}

	if !reflect.DeepEqual(resultStructPtr, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, resultStructPtr)
	}
}

func TestTomlSourceNothing(t *testing.T) {
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
	s := NewStaert(config, defaultPointersConfig)
	toml := NewTomlSource("nothing", []string{"./toml/", "/home/martin/go/src/github.com/containous/staert/toml"})
	s.Add(toml)
	result, err := s.GetConfig()
	if err != nil {
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

	//Type assertions
	resultStructPtr, ok := result.(*StructPtr)
	if !ok {
		t.Fatalf("Cannot convert the result into StructPtr")
	}

	if !reflect.DeepEqual(resultStructPtr, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, resultStructPtr)
	}
}

func TestTomlSourceTrivial(t *testing.T) {
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
	s := NewStaert(config, defaultPointersConfig)
	toml := NewTomlSource("trivial", []string{"./toml/", "/home/martin/go/src/github.com/containous/staert/toml"})
	s.Add(toml)
	result, err := s.GetConfig()
	if err != nil {
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

	//Type assertions
	resultStructPtr, ok := result.(*StructPtr)
	if !ok {
		t.Fatalf("Cannot convert the result into StructPtr")
	}

	if !reflect.DeepEqual(resultStructPtr, check) {
		fmt.Printf("expected\t: %+v\ngot\t\t\t: %+v\n", check.PtrStruct1, resultStructPtr.PtrStruct1)
		fmt.Printf("expected\t: %+v\ngot\t\t\t: %+v\n", check.PtrStruct2, resultStructPtr.PtrStruct2)
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, resultStructPtr)
	}
}

func TestTomlSourcePointer(t *testing.T) {
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
	s := NewStaert(config, defaultPointersConfig)
	toml := NewTomlSource("pointer", []string{"./toml/", "/home/martin/go/src/github.com/containous/staert/toml"})
	s.Add(toml)
	result, err := s.GetConfig()
	if err != nil {
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

	//Type assertions
	resultStructPtr, ok := result.(*StructPtr)
	if !ok {
		t.Fatalf("Cannot convert the result into StructPtr")
	}

	if !reflect.DeepEqual(resultStructPtr, check) {
		fmt.Printf("expected\t: %+v\ngot\t\t\t: %+v\n", check.PtrStruct1, resultStructPtr.PtrStruct1)
		fmt.Printf("expected\t: %+v\ngot\t\t\t: %+v\n", check.PtrStruct2, resultStructPtr.PtrStruct2)
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, resultStructPtr)
	}
}

func TestTomlSourcePointerUnderPointer(t *testing.T) {
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
	s := NewStaert(config, defaultPointersConfig)
	toml := NewTomlSource("pointerUnderPointer", []string{"./toml/", "/home/martin/go/src/github.com/containous/staert/toml"})
	s.Add(toml)
	result, err := s.GetConfig()
	if err != nil {
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

	//Type assertions
	resultStructPtr, ok := result.(*StructPtr)
	if !ok {
		t.Fatalf("Cannot convert the result into StructPtr")
	}

	if !reflect.DeepEqual(resultStructPtr, check) {
		fmt.Printf("expected\t: %+v\ngot\t\t\t: %+v\n", check.PtrStruct1, resultStructPtr.PtrStruct1)
		fmt.Printf("expected\t: %+v\ngot\t\t\t: %+v\n", check.PtrStruct2, resultStructPtr.PtrStruct2)
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, resultStructPtr)
	}
}

func TestTomlSourceFieldUnderPointerUnderPointer(t *testing.T) {
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
	s := NewStaert(config, defaultPointersConfig)
	toml := NewTomlSource("fieldUnderPtrUnderPtr", []string{"/home/martin/go/src/github.com/containous/staert/toml", "./toml/"})
	s.Add(toml)
	result, err := s.GetConfig()
	if err != nil {
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

	//Type assertions
	resultStructPtr, ok := result.(*StructPtr)
	if !ok {
		t.Fatalf("Cannot convert the result into StructPtr")
	}

	if !reflect.DeepEqual(resultStructPtr, check) {
		fmt.Printf("expected\t: %+v\ngot\t\t\t: %+v\n", check.PtrStruct1, resultStructPtr.PtrStruct1)
		fmt.Printf("expected\t: %+v\ngot\t\t\t: %+v\n", check.PtrStruct2, resultStructPtr.PtrStruct2)
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, resultStructPtr)
	}
}

func TestMergeTomlNothingFlaegNoArgs(t *testing.T) {
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
	s := NewStaert(config, defaultPointersConfig)
	toml := NewTomlSource("nothing", []string{"./toml/", "/home/martin/go/src/github.com/containous/staert/toml"})
	s.Add(toml)
	fs := NewFlaegSource(args, nil)
	s.Add(fs)
	result, err := s.GetConfig()
	if err != nil {
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

	//Type assertions
	resultStructPtr, ok := result.(*StructPtr)
	if !ok {
		t.Fatalf("Cannot convert the config into Configuration")
	}
	if !reflect.DeepEqual(resultStructPtr, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, resultStructPtr)
	}

}

func TestMergeTomlFieldUnderPointerUnderPointerFlaegNoArgs(t *testing.T) {
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
	s := NewStaert(config, defaultPointersConfig)
	toml := NewTomlSource("fieldUnderPtrUnderPtr", []string{"./toml/", "/home/martin/go/src/github.com/containous/staert/toml"})
	s.Add(toml)
	fs := NewFlaegSource(args, nil)
	s.Add(fs)
	result, err := s.GetConfig()
	if err != nil {
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

	//Type assertions
	resultStructPtr, ok := result.(*StructPtr)
	if !ok {
		t.Fatalf("Cannot convert the config into Configuration")
	}
	if !reflect.DeepEqual(resultStructPtr, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, resultStructPtr)
	}

}

func TestMergeTomlTrivialFlaegOverwriteField(t *testing.T) {
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
	s := NewStaert(config, defaultPointersConfig)
	toml := NewTomlSource("trivial", []string{"./toml/", "/home/martin/go/src/github.com/containous/staert/toml"})
	s.Add(toml)
	fs := NewFlaegSource(args, nil)
	s.Add(fs)
	result, err := s.GetConfig()
	if err != nil {
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

	//Type assertions
	resultStructPtr, ok := result.(*StructPtr)
	if !ok {
		t.Fatalf("Cannot convert the config into Configuration")
	}
	if !reflect.DeepEqual(resultStructPtr, check) {
		fmt.Printf("expected\t: %+v\ngot\t\t\t: %+v\n", check.PtrStruct1, resultStructPtr.PtrStruct1)
		fmt.Printf("expected\t: %+v\ngot\t\t\t: %+v\n", check.PtrStruct2, resultStructPtr.PtrStruct2)
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, resultStructPtr)
	}

}

func TestMergeTomlPointerUnderPointerFlaegManyArgs(t *testing.T) {
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
	s := NewStaert(config, defaultPointersConfig)
	toml := NewTomlSource("pointerUnderPointer", []string{"./toml/", "/home/martin/go/src/github.com/containous/staert/toml"})
	s.Add(toml)
	fs := NewFlaegSource(args, nil)
	s.Add(fs)
	result, err := s.GetConfig()
	if err != nil {
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
	//Type assertions
	resultStructPtr, ok := result.(*StructPtr)
	if !ok {
		t.Fatalf("Cannot convert the config into Configuration")
	}
	if !reflect.DeepEqual(resultStructPtr, check) {
		fmt.Printf("expected\t: %+v\ngot\t\t\t: %+v\n", check.PtrStruct1, resultStructPtr.PtrStruct1)
		fmt.Printf("expected\t: %+v\ngot\t\t\t: %+v\n", check.PtrStruct2, resultStructPtr.PtrStruct2)
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, resultStructPtr)
	}

}

func TestMergeFlaegNoArgsTomlNothing(t *testing.T) {
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
	s := NewStaert(config, defaultPointersConfig)
	fs := NewFlaegSource(args, nil)
	s.Add(fs)
	toml := NewTomlSource("nothing", []string{"./toml/", "/home/martin/go/src/github.com/containous/staert/toml"})
	s.Add(toml)
	result, err := s.GetConfig()
	if err != nil {
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

	//Type assertions
	resultStructPtr, ok := result.(*StructPtr)
	if !ok {
		t.Fatalf("Cannot convert the config into Configuration")
	}
	if !reflect.DeepEqual(resultStructPtr, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, resultStructPtr)
	}

}

func TestMergeFlaegFieldUnderPointerUnderPointerTomlNothing(t *testing.T) {
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
	s := NewStaert(config, defaultPointersConfig)
	fs := NewFlaegSource(args, nil)
	s.Add(fs)
	toml := NewTomlSource("nothing", []string{"./toml/", "/home/martin/go/src/github.com/containous/staert/toml"})
	s.Add(toml)
	result, err := s.GetConfig()
	if err != nil {
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

	//Type assertions
	resultStructPtr, ok := result.(*StructPtr)
	if !ok {
		t.Fatalf("Cannot convert the config into Configuration")
	}
	if !reflect.DeepEqual(resultStructPtr, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, resultStructPtr)
	}

}

func TestMergeFlaegManyArgsTomlOverwriteField(t *testing.T) {
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
	s := NewStaert(config, defaultPointersConfig)
	fs := NewFlaegSource(args, nil)
	s.Add(fs)
	toml := NewTomlSource("trivial", []string{"./toml/", "/home/martin/go/src/github.com/containous/staert/toml"})
	s.Add(toml)
	result, err := s.GetConfig()
	if err != nil {
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

	//Type assertions
	resultStructPtr, ok := result.(*StructPtr)
	if !ok {
		t.Fatalf("Cannot convert the config into Configuration")
	}
	if !reflect.DeepEqual(resultStructPtr, check) {
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, resultStructPtr)
	}

}
