package staert

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/containous/flaeg"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// Staert contains the struct to configure, thee default values inside structs and the sources
type Staert struct {
	DefaultPointersConfig interface{}
	Config                interface{}
	Sources               []Source
}

// Source interface must be satisfy to Add any kink of Source to Staert as like as TomlFile or Flaeg
type Source interface {
	Parse(sourceConfig interface{}, defaultPointersConfig interface{}) (interface{}, error)
}

// NewStaert creats and return a pointer on Staert. Need defaultConfig and defaultPointersConfig given by references
func NewStaert(defaultConfig interface{}, defaultPointersConfig interface{}) *Staert {
	s := Staert{}
	s.DefaultPointersConfig = defaultPointersConfig
	s.Config = defaultConfig
	return &s
}

// Add new Source to Staert, give it by reference
func (s *Staert) Add(src Source) {
	s.Sources = append(s.Sources, src)

}

// GetConfig run sources Parse func in the raw
// it retrurns a reference on the parsed config
func (s *Staert) GetConfig() (interface{}, error) {

	for _, src := range s.Sources {
		var err error
		s.Config, err = src.Parse(s.Config, s.DefaultPointersConfig)
		if err != nil {
			return nil, err
		}
	}
	return s.Config, nil
}

//FlaegSource impement Source
type FlaegSource struct {
	args          []string
	customParsers map[reflect.Type]flaeg.Parser
}

// NewFlaegSource creats and return a pointer on FlaegSource. Need args in a slice of string. customParsers should be nil if none
func NewFlaegSource(args []string, customParsers map[reflect.Type]flaeg.Parser) *FlaegSource {
	return &FlaegSource{args, customParsers}
}

// Parse calls Flaeg Load Function
func (fs *FlaegSource) Parse(sourceConfig interface{}, defaultPointersConfig interface{}) (interface{}, error) {
	if err := flaeg.LoadWithParsers(sourceConfig, defaultPointersConfig, fs.args, fs.customParsers); err != nil {
		return nil, err
	}
	return sourceConfig, nil
}

//TomlSource impement Source
type TomlSource struct {
	filename    string
	directories []string
	fullpath    string
}

// NewTomlSource creats and return a pointer on TomlSource. Parameter filename is the name of the file without neither fullpath not extension type and directories is a slice of paths
// (staert look for the toml file from the first directory to last one)
func NewTomlSource(filename string, directories []string) *TomlSource {
	//TODO trim path /, VAR ENV
	return &TomlSource{filename, directories, ""}
}

func (ts *TomlSource) findFile() error {
	for _, dir := range ts.directories {
		fullpath := string(dir[:len(dir)-1]) + strings.Trim(dir[len(dir)-1:], "/") + "/" + ts.filename + ".toml"
		// fmt.Printf("Lookup fullpath %s\n", fullpath)
		// Test if the file exits
		if _, err := os.Stat(fullpath); err == nil {
			//Turn fullpath in absolute representation of path
			fullpath, err = filepath.Abs(fullpath)
			if err != nil {
				return err
			}
			// fmt.Printf("File in fullpath %s exists\n", fullpath)
			ts.fullpath = fullpath
			return nil
		}
	}
	return fmt.Errorf("No file %s.toml found in directories %+v", ts.filename, ts.directories)
}

// Parse calls Flaeg Load Function
func (ts *TomlSource) Parse(sourceConfig interface{}, defaultPointersConfig interface{}) (interface{}, error) {
	if err := ts.findFile(); err != nil {
		return nil, err
	}
	metadata, err := toml.DecodeFile(ts.fullpath, sourceConfig)
	if err != nil {
		return nil, err
	}
	flaegArgs := []string{}
	keys := metadata.Keys()
	for i, key := range keys {
		// fmt.Println(key)
		if metadata.Type(key.String()) == "Hash" {
			//Ptr case
			// fmt.Printf("%s is a ptr\n", key)
			hasUnderField := false
			for j := i; j < len(keys); j++ {
				// fmt.Printf("%s =? %s\n", keys[j].String(), "."+key.String())
				if strings.Contains(keys[j].String(), key.String()+".") {
					hasUnderField = true
					break
				}
			}
			if !hasUnderField {
				flaegArgs = append(flaegArgs, "--"+strings.ToLower(key.String()))
			}
		}
	}
	// fmt.Println(flaegArgs)
	f := NewFlaegSource(flaegArgs, nil) //FIX ME : add custom parsers here
	sourceConfig, err = f.Parse(sourceConfig, defaultPointersConfig)
	if err != nil {
		return nil, err
	}
	return sourceConfig, nil
}
