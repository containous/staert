package staert

import (
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

//Struct1 : trivial Struct
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
	S2Float64 float64 `description:"Struct 3 float64"`
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
				S2Float64: 11.11,
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

func TestFleagSourcePtrArgs(t *testing.T) {
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
				S2Float64: 11.11,
			},
		},
		PtrStruct2: &Struct2{
			S2Int64:  22,
			S2String: "S2StringDefaultPointersConfig",
			S2Bool:   false,
		},
	}
	args := []string{
		"--ptrstruct1",
		"--ptrstruct2",
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
			S1Int:        11,
			S1String:     "S1StringDefaultPointersConfig",
			S1Bool:       true, //Wait for Flaeg ISSUE6
			S1PtrStruct3: nil,
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
		t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, resultStructPtr)
	}
}
