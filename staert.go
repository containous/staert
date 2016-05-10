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

// Command structure contains program/command information (command name and description)
// config must be a pointer on the configuration struct to parse (it contains default values of field)
// defaultPointersConfig contains default pointers values: those values are set on pointers fields if their flags are called
// It must be the same type(struct) as config
// Run is the func which launch the program using initialized configuration structure
type Command struct {
	Name                  string
	Description           string
	Run                   func(InitalizedConfig interface{}) error
	sources               []Source
	defaultPointersConfig interface{}
	config                interface{}
}

// Staert contains the struct to configure, thee default values inside structs and the sources
type Staert struct {
	commands []*Command
}

// Source interface must be satisfy to Add any kink of Source to Staert as like as TomlFile or Flaeg
type Source interface {
	Parse(sourceConfig interface{}, defaultPointersConfig interface{}) (interface{}, error)
}

// NewStaert creats and return a pointer on Staert. Need defaultConfig and defaultPointersConfig given by references
func NewStaert(rootCommand *Command) *Staert {
	s := Staert{
		commands: []*Command{rootCommand},
	}
	s.defaultPointersConfig = defaultPointersConfig
	s.config = defaultConfig
	return &s
}

// Add new Source to Staert, give it by reference
func (s *Staert) Add(src Source) {
	s.sources = append(s.sources, src)

}

// GetConfig run sources Parse func in the raw
// it retrurns a reference on the parsed config
func (s *Staert) GetConfig() (interface{}, error) {

	for _, src := range s.sources {
		var err error
		s.config, err = src.Parse(s.config, s.defaultPointersConfig)
		if err != nil {
			return nil, err
		}
	}
	return s.config, nil
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
