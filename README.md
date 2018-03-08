# Stært

[![Travis branch](https://img.shields.io/travis/containous/staert/master.svg)](https://travis-ci.org/containous/staert)
[![Coverage Status](https://coveralls.io/repos/github/containous/staert/badge.svg?branch=master)](https://coveralls.io/github/containous/staert?branch=master)
[![license](https://img.shields.io/github/license/containous/staert.svg)](https://github.com/containous/staert/blob/master/LICENSE.md)

Stært is a Go library for loading and merging a program configuration structure from many sources.

## Overview

Stært was born in order to merge two sources of configuration ([Flæg](https://github.com/containous/flaeg), [TOML](http://github.com/BurntSushi/toml)).
Now it also supports [Key-Value Store](#kvstore).

We developed [Flæg](https://github.com/containous/flaeg) and Stært in order to simplify configuration maintenance on [Træfik](https://github.com/containous/traefik).

## Features

- Load your configuration structure from many sources
- Keep your configuration structure values unchanged if no overwriting (support defaults values)
- Three native sources :
	- Command line arguments using [Flæg](https://github.com/containous/flaeg) package
	- TOML config file using [TOML](http://github.com/BurntSushi/toml) package
	- [Key-Value Store](#kvstore) using [libkv](https://github.com/docker/libkv) and [mapstructure](https://github.com/mitchellh/mapstructure) packages
- An interface to add your own sources
- Handle pointers field :
	- You can give a structure of default values for pointers
    - Same comportment as [Flæg](https://github.com/containous/flaeg)
- Stært is command oriented
    - It use `flaeg.Command`
    - Same comportment as [Flæg](https://github.com/containous/flaeg) commands
	- Stært supports only one command (the root-command)
	- Flæg allows you to use many commands
	- Only Flæg will be used if a sub-command is called. (because the configuration type could be different from one command to another)
	- You can add meta-data `"parseAllSources" -> "true"` to a sub-command if you want to parse all sources (it requires the same configuration type on the sub-command and the root-command)  

## Getting Started

### The configuration

It works on your own configuration structure, like this one :

```go
package example

import (
	"fmt"
	"github.com/containous/flaeg"
	"github.com/containous/staert"
	"os"
)

// Configuration is a struct which contains all different type to field
type Configuration struct {
	IntField     int                      `description:"An integer field"`
	StringField  string                   `description:"A string field"`
	PointerField *PointerSubConfiguration `description:"A pointer field"`
}

// PointerSubConfiguration is a SubStructure Configuration
type PointerSubConfiguration struct {
	BoolField  bool    `description:"A boolean field"`
	FloatField float64 `description:"A float field"`
}
```

Let's initialize it:

```go
func main() {
	// Init with default value
	config := &Configuration{
		IntField:    1,
		StringField: "init",
		PointerField: &PointerSubConfiguration{
			FloatField: 1.1,
		},
	}
	// Set default pointers value
	defaultPointersConfig := &Configuration{
		PointerField: &PointerSubConfiguration{
			BoolField:  true,
			FloatField: 99.99,
		},
	}
	//...
}
```

### The command

Stært uses `flaeg.Command` structure, like this:

```go
// Create command
command := &flaeg.Command{
	Name:"example",
	Description:"This is an example of description",
	Config:config,
	DefaultPointersConfig:defaultPointersConfig,
	Run: func() error {
		fmt.Printf("Run example with the config :\n%+v\n", config)
		fmt.Printf("PointerField contains:%+v\n", config.PointerField)
		return nil
	}
}
```

### Use stært with sources

Initialize Stært:

```go
s := staert.NewStaert(command)
```

Initialize TOML source:

```go
toml := staert.NewTomlSource("example", []string{"./toml/", "/any/other/path"})
```

Initialize Flæg source:

```go
f := flaeg.New(command, os.Args[1:])
```

### Add sources

Add TOML and Flæg sources:

```go
s.AddSource(toml)
s.AddSource(f)
```

**NB:** You can change order, so that, Flæg configuration will overwrite TOML one.

### Load your configuration

Just call `LoadConfig` function:

```go
loadedConfig, err := s.LoadConfig();
if err != nil {
	// oops
}
// do what you want with `loadedConfig`
// or call run function
```

### You can call Run

Run function will call `run()` from the command:

```go
if err := s.Run(); err != nil {
	//OOPS
}
```

**NB:** If you didn't call `LoadConfig()` before, your function `run()` will use your original configuration.

### Let's run example

TOML file `./toml/example.toml`:

```toml
IntField = 2
[PointerField]
```

We can run the example program using following CLI arguments:

```
$ ./example --stringfield=owerwrittenFromFlag --pointerfield.floatfield=55.55
Run example with the config :
&{IntField:2 StringField:owerwrittenFromFlag PointerField:0xc82000ec80}
PointerField contains:&{BoolField:true FloatField:55.55}
```

### Full example

[Tagoæl](https://github.com/debovema/tagoael) is a trivial example which shows how Stært can be use.
This funny GoLang program takes its configuration from both TOML and Flæg sources to display messages.

```
$ ./tagoael -h
tagoæl is an enhanced Hello World program to display messages with
an advanced configuration mechanism provided by Flæg & Stært.

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

As with Flæg and TOML sources, the configuration structure can be loaded from a Key-Value Store.
The package [libkv](https://github.com/docker/libkv) provides connection to many KV Store like `Consul`, `Etcd` or `Zookeeper`.

The whole configuration structure is stored, using architecture like this pattern:

- Key: `<prefix1>/<prefix2>/.../<fieldNameLevel1>/<fieldNameLevel2>/.../<fieldName>`
- Value: `<value>`

It handles:

- All [mapstructure](https://github.com/mitchellh/mapstructure) features(`bool`, `int`, ... , Squashed Embedded Sub `struct`, Pointer).
- Maps with pattern : `.../<MapFieldName>/<mapKey>` -> `<mapValue>` (Struct as key not supported)
- Slices (and Arrays) with pattern : `.../<SliceFieldName>/<SliceIndex>` -> `<value>`

**Note:** Hopefully, we provide the function `StoreConfig` to store your configuration structure ;)

### KvSource

KvSource implements Source:

```go
type KvSource struct {
	store.Store
	Prefix string // like this "prefix" (without the /)
}
```

### Initialize

It can be initialized like this:

```go
kv, err := staert.NewKvSource(backend store.Backend, addrs []string, options *store.Config, prefix string)
```

### LoadConfig

You can directly load data from the KV Store into the config structure (given by reference)

```go
config := &ConfigStruct{} // Here your configuration structure by reference
err := kv.Parse(config)
// do what you want with `config`
```

### Add to Stært sources

You can add this source to Stært, as with other sources:

```go
s.AddSource(kv)
```

### StoreConfig

You can also store your whole configuration structure into the KV Store:

```go
// We assume that `config` is initialized
err := kv.StoreConfig(config)
```


## Environment variables

You can extract configuration values from your process environment:

```go
env := staert.NewEnvSource(prefix, separator, parsers)
s.AddSource(env)
```

With this configuration, Stært will fetch from your environment all values according to your configuration structure.
Field names are split by words according to camelCase, we rely on [github.com/fatih/camelcase](https://github.com/fatih/camelcase) to manage this.

An environment variable name follows the pattern bellow:

```html
<PREFIX><SEP><MY><SEP><FIELD><SEP><NAME>
```

For instance if the `prefix` is `MyApp` and the `separator` is the `_` rune we'll have the following mapping:

```go
type Config struct {
    MyStringField string // => MYAPP_MY_STRING_FIELD
    MyIntField    int    // => MYAPP_MY_INT_FIELD
}
```

Type conversion is attempted using the configured `parsers` (see Flæg documentation for custom parsers).
If no parser is found for the given type, or if the parsing fails, an error is raised.

### Initialization

`EnvSource` can be initialized like this:

```go
envsource := staert.NewEnvSource(prefix, separator, parsers)
```

It takes three arguments:

- A prefix used in order to format environment variables names to fetch. If left blank, no prefix will be applied to environment variables
- A separator string. If left blank it will default to the `_` string
- A map of `reflect.Type` to `parse.Parser` to configure special parsing. (see [Flæg custom parser documentation](https://github.com/containous/flaeg#custom-parsers))

### Referenced values

You can use pointers to values in your configuration structs.
These fields will be mapped the same way as values.

Again with the same configuration:

```go
type AppConfig struct {
    Groot *int32 // => MYAPP_GROOT
}
```

### Nested structures

Nested structures are supported both by reference and value.

Field names are then prefixed with the field name referencing the nested structure:

```go
type PtrNestedConfig struct {
    AnArgument string // => MYAPP_FOO_AN_ARGUMENT
}

type ValueNestedConfig struct {
    AnotherArgument string // => MYAPP_BAR_ANOTHER_ARGUMENT
}

type AppConfig struct {
    Foo   *PtrNestedConfig
    Bar   ValueNestedConfig
}
```

### Embedded structures

Embedded structures are supported, and environment variable name generation for a field will have the same behavior than a normal struct field.

For instance if we keep our previous example configuration, we'll obtain the following mapping:

```go
type CommonConfig struct {
    CommonString string // => MYAPP_COMMON_STRING
}

type AppConfig struct {
    CommonConfig
}
```

### Array/Slices

Array elements can be configured by providing an index value between the array field name and the value:

```go
type AppConfig struct {
     Slice []string
}
```

```ini
MYAPP_SLICE_0=foo
MYAPP_SLICE_1=bar
```

The above configuration will give you a slice populated with `[foo, bar]`

Arrays/Slices of pointers is handled the same way.

#### Note

- The index part of the environment variable is only used for **ordering** the elements, **not** for the actual array index
- The size of the slice is increase dynamically based on the number of elements found

### Array/Slices of struct

For slices of struct we can map each field of the struct under an index :

```go
type SliceStruct struct {
    Name string
    Age int
}
type AppConfig struct {
     Slice []SliceStruct
}
```

```ini
MYAPP_SLICE_0_NAME=Bart
MYAPP_SLICE_0_AGE=14
MYAPP_SLICE_1_NAME=Lisa
MYAPP_SLICE_1_AGE=17
```

This will populate the slice with 2 elements : `Bart;14` and `Lisa;17`

### Maps

Maps are handled in a similar way as arrays and slices.

The map key separates the variable name from the value:

```go
type AppConfig struct {
     MyMap map[string]string
}
```

```ini
MYAPP_MY_MAP_key1=foo
MYAPP_MY_MAP_key2=bar
```

The above configuration will give you a map populated with `[[key1, foo], [key2, bar]]`

If the key type or value type is not `string`, a conversion is attempted using the configured parsers.
If no parser is found for the given type, or if the parsing fails, an error is raised.

Maps of structs are handled the same way as slices of structs.

## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D
