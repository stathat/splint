// Copyright 2011 Numrotron Inc.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.
//
// Developed at www.stathat.com by Patrick Crosby
// Contact us on twitter with any questions:  twitter.com/stat_hat

// splint is a little Go application to analyze Go source files.  It finds any functions that are
// too long or have too many parameters or results.
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

var statementThreshold = flag.Int("s", 30, "function statement count threshold")
var paramThreshold = flag.Int("p", 5, "parameter list length threshold")
var resultThreshold = flag.Int("r", 5, "result list length threshold")

type Parser struct {
	filename string
	first    bool
}

func NewParser(filename string) *Parser {
	return &Parser{filename, true}
}

func statementCount(n ast.Node) int {
	total := 0
	counter := func(node ast.Node) bool {
		switch node.(type) {
		case ast.Stmt:
			total += 1
		}
		return true
	}
	ast.Inspect(n, counter)
	return total
}

func (p *Parser) outputFilename() {
	if p.first {
		fmt.Printf("\n%s\n", p.filename)
		p.first = false
	}
}

func (p *Parser) checkFuncLength(x *ast.FuncDecl) {
	numStatements := statementCount(x)
	if numStatements > *statementThreshold {
		p.outputFilename()
		fmt.Printf("function %s too long: %d\n", x.Name, numStatements)
	}
}

func (p *Parser) checkParamCount(x *ast.FuncDecl) {
	numFields := x.Type.Params.NumFields()
	if numFields > *paramThreshold {
		p.outputFilename()
		fmt.Printf("function %s has too many params: %d\n", x.Name, numFields)
	}
}

func (p *Parser) checkResultCount(x *ast.FuncDecl) {
	numResults := x.Type.Results.NumFields()
	if numResults > *resultThreshold {
		p.outputFilename()
		fmt.Printf("function %s has too many params: %d\n", x.Name, numResults)
	}
}

func (p *Parser) examineFunc(x *ast.FuncDecl) {
	p.checkFuncLength(x)
	p.checkParamCount(x)
}

func (p *Parser) examineDecls(tree *ast.File) {
	for _, v := range tree.Decls {
		switch x := v.(type) {
		case *ast.FuncDecl:
			p.examineFunc(x)
		}
	}
}

func (p *Parser) Parse() {
	fileset := token.NewFileSet()
	tree, err := parser.ParseFile(fileset, p.filename, nil, 0)
	if err != nil {
		fmt.Printf("error parsing %s: %s\n", p.filename, err)
		return
	}

	p.examineDecls(tree)
}

func parseFile(filename string) {
	parser := NewParser(filename)
	parser.Parse()
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: splint [options] <go file>...")
		flag.PrintDefaults()
		os.Exit(1)
	}

	for _, v := range args {
		parseFile(v)
	}
}
