package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
	"time"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("usage: generate outputfile <package> <arg definitions>")
		fmt.Println("       generate outputfile args '{s:user:u}'")
		fmt.Printf("\nargs: %#v\n", os.Args)
		os.Exit(1)
	}

	fmt.Println("hi :)")

	f, err := os.Create(os.Args[1] + ".go")
	die(err)
	defer f.Close()

	goargsTemplate.Execute(f, struct {
		Timestamp string
		Package       string
	}{
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Package:   os.Args[2],
	})
}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var goargsTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
//
// This file was generated by robots at {{ .Timestamp }}

package {{ .Package }}

import (
	a "github.com/PKeidel/goargs"
)

var Args struct {
	Host string
}

func init() {
	a.WithArgV(a.TypeString, "host", "h", true, &Args.Host)
}
`))