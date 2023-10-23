# goargs

```go
package main

import (
	"fmt"
	"os"

	"github.com/PKeidel/goargs"
)

func main() {
	var args struct {
		Host string `goargs:"long:host,short:h"`
		User string
	}

	goargs.Bind(&args)

	goargs.EnvPrefix = ""
	goargs.ParseWith(
		[]string{"test", "--host", "localhost"},
		os.Environ(),
	)

	fmt.Printf("Host: %s\n", args.Host)
	fmt.Printf("User: %s\n", args.User)

	fmt.Println("To rerun, run:", os.Args[0], goargs.String())
}
```
