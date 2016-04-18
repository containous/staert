package staert

import (
	"github.com/cocap10/flaeg"
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

// New creats and return a pointer on Staert. Need defaultConfig and defaultPointersConfig given by references
func (s *Staert) New(defaultConfig interface{}, defaultPointersConfig interface{}) *Staert {
	s.DefaultPointersConfig = defaultPointersConfig
	s.Config = defaultConfig
	return s
}

// Add new Source to Staert, give it by reference
func (s *Staert) Add(src Source) error {
	s.Sources = append(s.Sources, src)
	return nil

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
	customParsers map[reflect.Type]flaeg.Parser
	args          []string
}

// AddParsers adds custom parsers, used in Flaeg
func (fs *FlaegSource) AddParsers(customParsers map[reflect.Type]flaeg.Parser) {
	fs.customParsers = customParsers
}

// AddArgs adds custom parsers, used in Flaeg
func (fs *FlaegSource) AddArgs(args []string) {
	fs.args = args
}

// Parse calls Flaeg Load Function
func (fs *FlaegSource) Parse(sourceConfig interface{}, defaultPointersConfig interface{}) (interface{}, error) {
	if err := flaeg.LoadWithParsers(sourceConfig, defaultPointersConfig, fs.args, fs.customParsers); err != nil {
		return nil, err
	}
	return sourceConfig, nil
}
