package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	fmt.Println("hello world")
	Stub(1)
}

type Struct interface {
	struct{}
}

func Stub[T any](s T) T {
	switch s.(type) {
	case int:
		spew.Dump("int")
	}
	return s
}