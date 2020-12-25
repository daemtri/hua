package huarpc

import (
	"embed"
	"go/parser"
	"go/token"

	"github.com/davecgh/go-spew/spew"
)

func Extend(protocol embed.FS) {
	tokenFs := token.NewFileSet()
	calc, err := protocol.ReadFile("calcservice.api.go")
	if err != nil {
		panic(err)
	}
	f, err := parser.ParseFile(tokenFs, "calcservice.api.go", calc, parser.ParseComments|parser.AllErrors)
	if err != nil {
		panic(err)
	}
	spew.Dump(f.Comments)
}
