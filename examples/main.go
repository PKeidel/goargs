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
		Port int `goargs:"short:p"`
	}

	goargs.Bind(&args)

	goargs.EnvPrefix = ""
    goargs.Parse()

	for k, v := range goargs.ArgsByLongname {
		fmt.Printf("args: %s => %s\n", k, v.Source)
	}

	fmt.Printf("args: %+v\n", args)

    fmt.Println("To rerun, run:", os.Args[0], goargs.String())
}
