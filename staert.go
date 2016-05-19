package staert

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/containous/flaeg"
	"os"
	"path/filepath"
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
func (s *Staert) getConfig(cmd *flaeg.Command) error {
	for _, src := range s.sources {
		var err error
		_, err = src.Parse(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

// Run calls the Run func of the command with the parsed config
func (s *Staert) Run() error {
	cmd := s.command
	for _, src := range s.sources {
		//Type assertion
		f, ok := src.(*flaeg.Flaeg)
		if ok {
			if fCmd, err := f.GetCommand(); err != nil {
				return err
			} else if cmd != fCmd {
				if err := f.Run(); err != nil {
					return err
				}
				return nil
			}
		}
	}
	if err := s.getConfig(cmd); err != nil {
		return err
	}
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
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
	return &TomlSource{filename, directories, ""}
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

func (ts *TomlSource) findFile() error {
	for _, d := range ts.directories {
		if d != "" {
			dir, err := preprocessDir(d)
			if err != nil {
				return err
			}
			fullpath := dir + "/" + ts.filename + ".toml"
			// fmt.Printf("Lookup fullpath %s\n", fullpath)
			// Test if the file exits
			if _, err := os.Stat(fullpath); err == nil {
				//Turn fullpath in absolute representation of path

				// fmt.Printf("File in fullpath %s exists\n", fullpath)
				ts.fullpath = fullpath
				return nil
			}
		}
	}
	return fmt.Errorf("No file %s.toml found in directories %+v", ts.filename, ts.directories)
}

// Parse calls Flaeg Load Function
func (ts *TomlSource) Parse(cmd *flaeg.Command) (*flaeg.Command, error) {
	if err := ts.findFile(); err != nil {
		return nil, err
	}
	metadata, err := toml.DecodeFile(ts.fullpath, cmd.Config)
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
	err = flaeg.Load(cmd.Config, cmd.DefaultPointersConfig, flaegArgs)
	//if err!= missing parser err
	if err != nil && !strings.Contains(err.Error(), "No parser for type") {
		return nil, err
	}
	return cmd, nil
}
