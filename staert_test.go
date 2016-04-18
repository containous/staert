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
	S1Int    int    `description:"Struct 1 Int"`
	S1String string `description:"Struct 1 String"`
	S1Bool   bool   `description:"Struct 1 Bool"`
}

//Struct2 : trivial Struct
type Struct2 struct {
	S2Int64  int64  `description:"Struct 2 Int64"`
	S2String string `description:"Struct 2 String"`
	S2Bool   bool   `description:"Struct 2 Bool"`
}

func TestFleagSource(t *testing.T) {
	//Init
	defaultStructPtr := StructPtr{
		PtrStruct1: &Struct1{
			1000,
			"S1StringDefault",
			true,
		},
		PtrStruct2: &Struct2{
			2000,
			"S2StringDefault",
			false,
		},
		DurationField: time.Second * 5,
	}
	config := StructPtr{
		PtrStruct1: &Struct1{
			3,
			"S1StringNonDefault",
			false,
		},
		PtrStruct2: nil,
	}
	args := []string{
		// "-h",
		"--durationfield=50s",
		"--ptrstruct2.s2int64=1111",
		// "--ptrstruct2",

	}

	//Test
	s := NewStaert(&config, &defaultStructPtr)
	fs := NewFlaegSource(args, nil)
	s.Add(fs)
	result, err := s.GetConfig()
	if err != nil {
		t.Errorf("Error %s", err.Error())
	}

	//Check
	check := StructPtr{
		PtrStruct1: &Struct1{
			3,
			"S1StringNonDefault",
			false,
		},
		PtrStruct2: &Struct2{
			1111,
			"S2StringDefault",
			false,
		},
		DurationField: time.Second * 50,
	}

	if resultStructPtr, ok := result.(*StructPtr); ok {
		if !reflect.DeepEqual(*resultStructPtr, check) {
			t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check, *resultStructPtr)
		}
	}
}

func TestTomlSource(t *testing.T) {
	//Init
	defaultStructPtr := StructPtr{
		PtrStruct1: &Struct1{
			1000,
			"S1StringDefault",
			true,
		},
		PtrStruct2: &Struct2{
			2000,
			"S2StringDefault",
			false,
		},
		DurationField: time.Second * 5,
	}
	config := StructPtr{
		PtrStruct1: &Struct1{
			3,
			"S1StringNonDefault",
			false,
		},
		PtrStruct2: nil,
	}

	//Test
	s := NewStaert(&config, &defaultStructPtr)
	ts := NewTomlSource("test", []string{".", "$HOME"})
	s.Add(ts)
	result, err := s.GetConfig()
	if err != nil {
		t.Errorf("Error %s", err.Error())
	}
	// fmt.Printf("%+v\n", result)

	//Check
	check := StructPtr{
		PtrStruct1: &Struct1{
			1111,
			"S1StringNonDefault",
			false,
		},
		PtrStruct2: &Struct2{
			2222,
			"S2StringToml",
			true,
		},
		DurationField: time.Second * 9,
	}

	if resultStructPtr, ok := result.(*StructPtr); ok {
		if !reflect.DeepEqual(*resultStructPtr, check) {
			t.Fatalf("\nexpected\t: %+v\ngot\t\t\t: %+v\n", check.PtrStruct1, resultStructPtr.PtrStruct1)
		}
	}
}
