package main

import (
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	spew.Dump(strings.SplitN("xx,xx", ",", 3))
}
