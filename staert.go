package staert

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/cocap10/flaeg"
	"os"
	"reflect"
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

// GetConfig new Source to Staert
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
	name        string
	directories []string
	file        string
}

// NewTomlSource creats and return a pointer on TomlSource. Need file name (without path and extension type) and directories paths+
func NewTomlSource(name string, directories []string) *TomlSource {
	//TODO trim path /, VAR ENV
	return &TomlSource{name, directories, ""}
}

func (ts *TomlSource) findFile() error {
	for _, dir := range ts.directories {
		file := dir + "/" + ts.name + ".toml"
		if _, err := os.Stat(file); err == nil {
			ts.file = file
			return nil
		}
	}
	return fmt.Errorf("No file %s.toml found in directories %+v", ts.name, ts.directories)
}

// Parse calls Flaeg Load Function
func (ts *TomlSource) Parse(sourceConfig interface{}, defaultPointersConfig interface{}) (interface{}, error) {
	if err := ts.findFile(); err != nil {
		return nil, err
	}
	if _, err := toml.DecodeFile(ts.file, sourceConfig); err != nil {
		return nil, err
	}
	return sourceConfig, nil
}
