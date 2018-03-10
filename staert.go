package staert

import (
	"fmt"
	"reflect"

	"github.com/containous/flaeg"
)

// Source interface must be satisfy to Add any kink of Source to Staert as like as TomlFile or Flaeg
type Source interface {
	Parse(cmd *flaeg.Command) (*flaeg.Command, error)
}

// Staert contains the struct to configure, thee default values inside structs and the sources
type Staert struct {
	command *flaeg.Command
	sources []Source
}

// NewStaert creates and return a pointer on Staert. Need defaultConfig and defaultPointersConfig given by references
func NewStaert(rootCommand *flaeg.Command) *Staert {
	s := Staert{
		command: rootCommand,
	}
	return &s
}

// AddSource adds new Source to Staert, give it by reference
func (s *Staert) AddSource(src Source) {
	s.sources = append(s.sources, src)
}

// parseConfigAllSources getConfig for a flaeg.Command run sources Parse func in the raw
func (s *Staert) parseConfigAllSources(cmd *flaeg.Command) error {
	for _, src := range s.sources {
		var err error
		_, err = src.Parse(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadConfig check which command is called and parses config
// It returns the the parsed config or an error if it fails
func (s *Staert) LoadConfig() (interface{}, error) {
	for _, src := range s.sources {
		// Type assertion
		f, ok := src.(*flaeg.Flaeg)
		if ok {
			fCmd, err := f.GetCommand()
			if err != nil {
				return nil, err
			}

			if s.command != fCmd {
				// if fleag sub-command
				if fCmd.Metadata["parseAllSources"] == "true" {
					// if parseAllSources
					fCmdConfigType := reflect.TypeOf(fCmd.Config)
					sCmdConfigType := reflect.TypeOf(s.command.Config)
					if fCmdConfigType != sCmdConfigType {
						return nil, fmt.Errorf("command %s : Config type doesn't match with root command config type. Expected %s got %s",
							fCmd.Name, sCmdConfigType.Name(), fCmdConfigType.Name())
					}
					s.command = fCmd
				} else {
					// (not parseAllSources)
					s.command, err = f.Parse(fCmd)
					return s.command.Config, err
				}
			}
		}
	}
	err := s.parseConfigAllSources(s.command)
	return s.command.Config, err
}

// Run calls the Run func of the command
// Warning, Run doesn't parse the config
func (s *Staert) Run() error {
	return s.command.Run()
}
