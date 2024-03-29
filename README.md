# goargs

Use command-line arguments as easy as:
```go
var args struct {
	Host string `goargs:"short:h"`
	User string
}
goargs.Bind(&args)
goargs.Parse()

fmt.Printf("Host: %s\n", args.Host)
```

<details>
  <summary>See full file</summary>

```go
package main

import (
	"fmt"
	"os"

	"github.com/PKeidel/goargs"
)

var args struct {
	Host string `goargs:"short:h"`
	User string
}

func main() {
	goargs.Bind(&args)
	goargs.Parse()

	fmt.Printf("Host: %s\n", args.Host)
	fmt.Printf("User: %s\n", args.User)

	fmt.Println("To rerun, run:", os.Args[0], goargs.String())
}
```
  
</details>

<br>

Now you can run your program for example with: `program --user PKeidel -h host01`

`strcase.ToKebab` is applied to all "long" names. So a value like `UserName` will become `--user-name` as an arg.

## Tags
There are a few tags to set/overrite the long form, short form and the required flag per field. `long` defaults to the field name in lowercase.
```go
var args struct {
	Host string `goargs:"long:host,short:h,required"`
	User string
}
```

## String()
`goargs.String()` returns a command-line variant of the current values. This can be usefull if the values are combined from different inputs but you want to rerun your program with the exact same arguments.

```go
fmt.Println("To rerun, run:", os.Args[0], goargs.String())
// To rerun, run: program --host host01
```

## Sources
`goargs` will search for the values in your `os.Args`, `os.Environ` and in the future in config files (json, yaml).

The order in which values overrite eachother can be configured with `goargs.DefaultSourceOrder` and defaults to: `SourceFile, SourceEnv, SourceArgs` with args overriding all other values.

## Without structs
You can also register single values. For example:
```go
var host string

func main() {
	goargs.WithString("host", &host)
	goargs.Parse()

	fmt.Printf("Host: %s\n", host)
}
```

To set for example the `Required` flag, you need to use the reference returned from the `With*` functions:
```go
argHost := goargs.WithStringS("host", "h", &host)
argHost.Required = true
goargs.Parse()
```

## Roadmap
```
[x] Support strings
[x] Support ints
[x] Support flags (bool)
[x] Support require flag
[x] Support os.Args parsing
[x] Support env parsing
[ ] Support multi flag to allow arrays like `--host host01 host02`
[ ] Support config file parsing
[ ] Detect undefined args
[ ] Detect wrong types
[x] Support externally changed values
[ ] Support values with spaces, quotes, ...
[ ] Print useage
```