# Stært
[![Travis branch](https://img.shields.io/travis/containous/staert/master.svg)](https://travis-ci.org/containous/staert)
[![Coverage Status](https://coveralls.io/repos/github/containous/staert/badge.svg?branch=master)](https://coveralls.io/github/containous/staert?branch=master)
[![license](https://img.shields.io/github/license/containous/staert.svg)](https://github.com/containous/staert/blob/master/LICENSE.md)

Stært is a Go library for loading and merging a program configuration structure from many sources.

## Overview
Stært born in order to merge two sources of Configuration ([Flæg](https://github.com/containous/flaeg), [Toml](http://github.com/BurntSushi/toml))

## Features
 - Load your Configuration structure from many sources
 - Keep your Configuration structure values unchanged if no overwriting (support defaults values)
 - Two native sources :
	- [Flæg](https://github.com/containous/flaeg)
	- [Toml](http://github.com/BurntSushi/toml)
 - An Interface to add your own sources
 - Handle pointers field :
	- You can give a structure of default values for pointers
    - Same comportment as [Flæg](https://github.com/containous/flaeg)
 - Stært is Command oriented
    - It use `flaeg.Command`
    - Same comportment as [Flæg](https://github.com/containous/flaeg) commands

## Getting Started
### The configuration
It works on your own Configuration structure, like this one :
```go
import (
	"github.com/containous/flaeg"
    "github.com/containous/staert"
)

//Configuration is a struct which contains all differents type to field
type Configuration struct {
	IntField     int                      `description:"An integer field"`
	StringField  string                   `description:"A string field"`
	PointerField *PointerSubConfiguration `description:"A pointer field"`
}

//PointerSubConfiguration is a SubStructure Configuration
type PointerSubConfiguration struct {
	BoolField  `description:"A boolean field"`
	FloatField `description:"A float field"`
}
```

Let's initialize it: 
```go
//Init with default value
	config := &Configuration{
		IntField:    1,
		StringField: "init",
		PointerField: &PointerSubConfiguration{
			FloatField: 1.1,
		},
	}
	//Set default pointers value
	defaultPointersConfig := &Configuration{
		PointerField: &PointerSubConfiguration{
			BoolField:  true,
			FloatField: 99.99,
		},
	}
```

### The command
Stært uses `flaeg.Command` Structure, like this :
```go
//Creat Command
    command:=&flaeg.Command{
        Name:"example",
        Description:"This is an example of description",
        Config:config,
        DefaultPointersConfig:defaultPointersConfig,
        Run: func() error {
            fmt.Printf("Run example with the config :\n%+v\n",config)
            return nil
        }
    }
```

### Use stært with sources
Init Stært
```go
    s:=staert.NewStaert(command)
```
Init TOML source
```go
     toml:=staert.NewTomlSource("example", []string{"./toml/", "/any/other/path"})
```
Init Flæg source
```go
     f:=flaeg.New(command, os.Args[1])
```
### Add sources
Add TOML and flæg sources
```go
    s.AddSource(toml)
    s.AddSource(f)
``` 

### Now run
JUST RUN
```go
    if err := s.Run(); err != nil {
		//OUPS
	}
``` 

## Contributing
1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D
