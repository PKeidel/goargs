package main

import (
	"fmt"
	"os"

	"github.com/PKeidel/goargs"
)

func main() {
	var args struct {
		Host  string `goargs:"long:host,short:h"`
		User  string
		Port  int `goargs:"short:p"`
		Debug bool
	}

	goargs.Bind(&args)

	goargs.Parse()

	for _, arg := range goargs.DefaultArgs {
		fmt.Printf("args: %s => %s (%v)\n", arg.LongName, arg.Source, arg.Val())
	}

	fmt.Printf("args: %+v\n", args)

	for _, v := range goargs.DefaultArgs {
		fmt.Printf("- %+v\n", v)
	}

	fmt.Println("To rerun, run:", os.Args[0], goargs.String())
}
