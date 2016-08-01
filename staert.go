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

// Source interface must be satisfy to Add any kink of Source to Staert as like as TomlFile or Flaeg
type Source interface {
	Parse(cmd *flaeg.Command) (*flaeg.Command, error)
}

// Staert contains the struct to configure, thee default values inside structs and the sources
type Staert struct {
	command *flaeg.Command
	sources []Source
}

// NewStaert creats and return a pointer on Staert. Need defaultConfig and defaultPointersConfig given by references
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

// getConfig for a flaeg.Command run sources Parse func in the raw
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
		//Type assertion
		f, ok := src.(*flaeg.Flaeg)
		if ok {
			if fCmd, err := f.GetCommand(); err != nil {
				return nil, err
			} else if s.command != fCmd {
				//IF fleag sub-command
				if fCmd.Metadata["parseAllSources"] == "true" {
					//IF parseAllSources
					fCmdConfigType := reflect.TypeOf(fCmd.Config)
					sCmdConfigType := reflect.TypeOf(s.command.Config)
					if fCmdConfigType != sCmdConfigType {
						return nil, fmt.Errorf("Command %s : Config type doesn't match with root command config type. Expected %s got %s", fCmd.Name, sCmdConfigType.Name(), fCmdConfigType.Name())
					}
					s.command = fCmd
				} else {
					// ELSE (not parseAllSources)
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

//TomlSource impement Source
type TomlSource struct {
	filename     string
	dirNfullpath []string
	fullpath     string
}

// NewTomlSource creats and return a pointer on TomlSource.
// Parameter filename is the file name (without extension type, ".toml" will be added)
// dirNfullpath may contain directories or fullpath to the file.
func NewTomlSource(filename string, dirNfullpath []string) *TomlSource {
	return &TomlSource{filename, dirNfullpath, ""}
}

// ConfigFileUsed return config file used
func (ts *TomlSource) ConfigFileUsed() string {
	return ts.fullpath
}

func preprocessDir(dirIn string) (string, error) {
	dirOut := dirIn
	if strings.HasPrefix(dirIn, "$") {
		end := strings.Index(dirIn, string(os.PathSeparator))
		if end == -1 {
			end = len(dirIn)
		}
		dirOut = os.Getenv(dirIn[1:end]) + dirIn[end:]
	}
	dirOut, err := filepath.Abs(dirOut)
	return dirOut, err
}

func findFile(filename string, dirNfile []string) string {
	for _, df := range dirNfile {
		if df != "" {
			fullpath, _ := preprocessDir(df)
			if fileinfo, err := os.Stat(fullpath); err == nil && !fileinfo.IsDir() {
				return fullpath
			}
			fullpath = fullpath + "/" + filename + ".toml"
			if fileinfo, err := os.Stat(fullpath); err == nil && !fileinfo.IsDir() {
				return fullpath
			}
		}
	}
	return ""
}

// Parse calls toml.DecodeFile() func
func (ts *TomlSource) Parse(cmd *flaeg.Command) (*flaeg.Command, error) {
	ts.fullpath = findFile(ts.filename, ts.dirNfullpath)
	if len(ts.fullpath) < 2 {
		return cmd, nil
	}
	metadata, err := toml.DecodeFile(ts.fullpath, cmd.Config)
	if err != nil {
		return nil, err
	}
	flags, err := flaeg.GetFlags(cmd.Config)
	if err != nil {
		return nil, err
	}
	flaegArgs, err := generateArgs(metadata, flags)
	if err != nil {
		return nil, err
	}

	// fmt.Println(flaegArgs)
	err = flaeg.Load(cmd.Config, cmd.DefaultPointersConfig, flaegArgs)
	//if err!= missing parser err
	if err != nil && err != flaeg.ErrParserNotFound {
		return nil, err
	}
	return cmd, nil
}

func generateArgs(metadata toml.MetaData, flags []string) ([]string, error) {
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
				match := false
				for _, flag := range flags {
					if flag == strings.ToLower(key.String()) {
						match = true
						break
					}
				}
				if match {
					flaegArgs = append(flaegArgs, "--"+strings.ToLower(key.String()))
				}
			}
		}
	}
	return flaegArgs, nil
}
