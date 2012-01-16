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
	"encoding/json"
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
var outputJSON = flag.Bool("j", false, "output results as json")

type Parser struct {
	filename string
	first    bool
	summary  *Summary
}

type Offender struct {
	Filename string
	Function string
	Count    int
}

type Summary struct {
	Statement []*Offender
	Param     []*Offender
	Result    []*Offender

	// redundant, but using these for easy json output
	NumAboveStatementThreshold int
	NumAboveParamThreshold     int
	NumAboveResultThreshold    int
}

func (s *Summary) addStatement(filename, function string, count int) {
	s.Statement = append(s.Statement, &Offender{filename, function, count})
	s.NumAboveStatementThreshold++
}

func (s *Summary) addParam(filename, function string, count int) {
	s.Param = append(s.Param, &Offender{filename, function, count})
	s.NumAboveParamThreshold++
}

func (s *Summary) addResult(filename, function string, count int) {
	s.Result = append(s.Result, &Offender{filename, function, count})
	s.NumAboveResultThreshold++
}

func NewParser(filename string, summary *Summary) *Parser {
	return &Parser{filename, true, summary}
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
	if *outputJSON {
		return
	}
	if p.first {
		fmt.Printf("\n%s\n", p.filename)
		p.first = false
	}
}

func (p *Parser) checkFuncLength(x *ast.FuncDecl) {
	numStatements := statementCount(x)
	if numStatements <= *statementThreshold {
		return
	}

	p.summary.addStatement(p.filename, x.Name.String(), numStatements)

	if *outputJSON == false {
		p.outputFilename()
		fmt.Printf("function %s too long: %d\n", x.Name, numStatements)
	}
}

func (p *Parser) checkParamCount(x *ast.FuncDecl) {
	numFields := x.Type.Params.NumFields()
	if numFields <= *paramThreshold {
		return
	}

	p.summary.addParam(p.filename, x.Name.String(), numFields)
	if *outputJSON == false {
		p.outputFilename()
		fmt.Printf("function %s has too many params: %d\n", x.Name, numFields)
	}
}

func (p *Parser) checkResultCount(x *ast.FuncDecl) {
	numResults := x.Type.Results.NumFields()
	if numResults <= *resultThreshold {
		return
	}

	p.summary.addResult(p.filename, x.Name.String(), numResults)
	if *outputJSON == false {
		p.outputFilename()
		fmt.Printf("function %s has too many results: %d\n", x.Name, numResults)
	}
}

func (p *Parser) examineFunc(x *ast.FuncDecl) {
	p.checkFuncLength(x)
	p.checkParamCount(x)
	p.checkResultCount(x)
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

func parseFile(filename string, summary *Summary) {
	parser := NewParser(filename, summary)
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

	summary := new(Summary)

	for _, v := range args {
		parseFile(v, summary)
	}

	if *outputJSON {
		/*
		   buf := new(bytes.Buffer)
		   encoder := json.NewEncoder(buf)
		   err := encoder.Encode(summary)
		   if err != nil {
		           fmt.Println("json encode error:", err)
		   }
		   fmt.Println(string(buf.Bytes()))
		*/
		data, err := json.MarshalIndent(summary, "", "\t")
		if err != nil {
			fmt.Println("json encode error:", err)
		}
		fmt.Println(string(data))

	} else {
		fmt.Println()
		fmt.Println("Number of functions above statement threshold:", summary.NumAboveStatementThreshold)
		fmt.Println("Number of functions above param threshold:", summary.NumAboveParamThreshold)
		fmt.Println("Number of functions above result threshold:", summary.NumAboveResultThreshold)
	}
}
