//go:build debug
// +build debug

package goargs

import "fmt"

func log(msg string) {
	fmt.Println(msg)
}

func logf(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
}
