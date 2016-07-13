# Stært
[![Travis branch](https://img.shields.io/travis/containous/staert/master.svg)](https://travis-ci.org/containous/staert)
[![Coverage Status](https://coveralls.io/repos/github/containous/staert/badge.svg?branch=master)](https://coveralls.io/github/containous/staert?branch=master)
[![license](https://img.shields.io/github/license/containous/staert.svg)](https://github.com/containous/staert/blob/master/LICENSE.md)

Stært is a Go library for loading and merging a program configuration structure from many sources.

## Overview
Stært was born in order to merge two sources of Configuration ([Flæg](https://github.com/containous/flaeg), [Toml](http://github.com/BurntSushi/toml)).
Now it also supports [Key-Value Store](#kvstore).
We developed [Flæg](https://github.com/containous/flaeg) and Stært in order to simplify configuration maintenance on [Træfik](https://github.com/containous/traefik).

## Features
 - Load your Configuration structure from many sources
 - Keep your Configuration structure values unchanged if no overwriting (support defaults values)
 - Three native sources :
	- Command line arguments using [flæg](https://github.com/containous/flaeg) package
	- TOML config file using [toml](http://github.com/BurntSushi/toml) package
	- [Key-Value Store](#kvstore) using [libkv](https://github.com/docker/libkv) and [mapstructure](https://github.com/mitchellh/mapstructure) packages
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
package example

import (
	"fmt"
	"github.com/containous/flaeg"
	"github.com/containous/staert"
	"os"
)

//Configuration is a struct which contains all differents type to field
type Configuration struct {
	IntField     int                      `description:"An integer field"`
	StringField  string                   `description:"A string field"`
	PointerField *PointerSubConfiguration `description:"A pointer field"`
}

//PointerSubConfiguration is a SubStructure Configuration
type PointerSubConfiguration struct {
	BoolField  bool    `description:"A boolean field"`
	FloatField float64 `description:"A float field"`
}
```

Let's initialize it: 
```go
 func main() {
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
    //Create Command
    command:=&flaeg.Command{
        Name:"example",
        Description:"This is an example of description",
        Config:config,
        DefaultPointersConfig:defaultPointersConfig,
        Run: func() error {
            fmt.Printf("Run example with the config :\n%+v\n",config)
 			fmt.Printf("PointerField contains:%+v\n", config.PointerField)
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
     f:=flaeg.New(command, os.Args[1:])
```
### Add sources
Add TOML and flæg sources
```go
    s.AddSource(toml)
    s.AddSource(f)
``` 
NB : You can change order, so that, flaeg configuration will overwrite toml one 
### Load your configuration
Just call LoadConfig function :
```go
	loadedConfig, err := s.LoadConfig();
    if err != nil {
		//OOPS
	}
	//DO WHAT YOU WANT WITH loadedConfig 
	//OR CALL RUN FUNC
``` 

### You can call Run
Run function will call the func `run()` from the command :
```go
    if err := s.Run(); err != nil {
		//OOPS
	}
 }
``` 
 NB : If you didn't call `LoadConfig()` before, your func `run()` will use your original configuration
### Let's run example

TOML file `./toml/example.toml` :
```toml
IntField= 2
[PointerField]
``` 
We can run the example program using folowing CLI arguments :
```
$ ./example --stringfield=owerwrittenFromFlag --pointerfield.floatfield=55.55
Run example with the config :
&{IntField:2 StringField:owerwrittenFromFlag PointerField:0xc82000ec80}
PointerField contains:&{BoolField:true FloatField:55.55}

```

## Full example : 
[Tagoæl](https://github.com/debovema/tagoael) is a trivial example which shows how Stært can be use.
This funny golang progam takes its configuration from both TOML and Flaeg sources to display messages.
```shell
$ ./tagoael -h
tagoæl is an enhanced Hello World program to display messages with
an advanced configuration mechanism provided by flæg & stært.

flæg:   https://github.com/containous/flaeg
stært:  https://github.com/containous/staert
tagoæl: https://github.com/debovema/tagoael


Usage: tagoael [--flag=flag_argument] [-f[flag_argument]] ...     set flag_argument to flag(s)
   or: tagoael [--flag[=true|false| ]] [-f[true|false| ]] ...     set true/false to boolean flag(s)

Flags:
        -c, --commandlineoverridesconfigfile               Whether configuration from command line overrides configuration from configuration file or not. (default "true")
        --configfile                                       Configuration file to use (TOML). (default "tagoael")
        -i, --displayindex                                 Whether to display index of each message (default "false")
        -m, --messagetodisplay                             Message to display (default "HELLO WOLRD")
        -n, --numbertodisplay                              Number of messages to display (default "1000")
        -h, --help                                         Print Help (this message) and exit
```
Thank you [@debovema](https://github.com/debovema) for this work :)

## KvStore
As with Flæg and Toml sources, the configuration structure can be loaded from a Key-Value Store.
The package [libkv](https://github.com/docker/libkv) provides connection to many KV Store like `Consul`, `Etcd` or `Zookeeper`.

The whole configuration structure is stored, using architecture like this pattern :
 - Key : `<prefix1>/<prefix2>/.../<fieldNameLevel1>/<fieldNameLevel2>/.../<fieldName>`
 - Value : `<value>`
 
It handles :
 - All [mapstructure](https://github.com/mitchellh/mapstructure) features(`bool`, `int`, ... , Squashed Embedded Sub `struct`, Pointer).
 - Maps with pattern : `.../<MapFieldName>/<mapKey>` -> `<mapValue>` (Struct as key not supported)
 - Slices (and Arrays) with pattern : `.../<SliceFieldName>/<SliceIndex>` -> `<value>`

Note : Hopefully, we provide the function `StoreConfig` to store your configuration structure ;)

### KvSource
KvSource implements Source: 

```go
type KvSource struct {
	store.Store
	Prefix string // like this "prefix" (witout the /)
}
```

### Initialize 
It can be initialized like this :
```go
	kv, err := staert.NewKvSource(backend store.Backend, addrs []string, options *store.Config, prefix string)
```

### LoadConfig
You can directly load data from the KV Store into the config structure (given by reference)
```go
	config := &ConfigStruct{} // Here your configuration structure by reference
	err := kv.Parse(config)
	//DO WHAT YOU WANT WITH config
```

### Add to Stært sources
Or you can add this source to Stært, as with other sources
```go
	s.AddSource(kv)
```

### StoreConfig
You can also store your whole configuration structure into the KV Store :
```go
	// We assume that config is Initialazed
	err := s.StoreConfig(config)
```


## Contributing
1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D
