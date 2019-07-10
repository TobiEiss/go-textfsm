package main

import (
	"fmt"
	"os"

	"github.com/TobiEiss/go-textfsm/pkg/ast"
	"github.com/TobiEiss/go-textfsm/pkg/process"
	"github.com/TobiEiss/go-textfsm/pkg/reader"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		warning(err.Error())
	}

	args := os.Args[1:]

	// check args
	if len(args) < 2 {
		warning("Please specify template and input file like 'go-textfsm <templatefile.tmpl> <inputfile>'")
	}

	// extract filepath
	templateFilePath := dir + "/" + args[0]
	if _, err := os.Stat(templateFilePath); os.IsNotExist(err) {
		warning("Template doesn't exist: " + err.Error())
	}

	// extract sourceFile
	sourceFilePath := dir + "/" + args[1]
	if _, err := os.Stat(sourceFilePath); os.IsNotExist(err) {
		warning("Template doesn't exist: " + err.Error())
	}

	// read template
	tmplCh := make(chan string)
	go reader.ReadLineByLine(templateFilePath, tmplCh)

	srcCh := make(chan string)
	go reader.ReadLineByLine(sourceFilePath, srcCh)

	// create AST
	ast, err := ast.CreateAST(tmplCh)
	if err != nil {
		warning(err.Error())
	}

	// process ast
	record := make(chan []interface{})
	process, err := process.NewProcess(ast, record)
	if err != nil {
		warning(err.Error())
	}
	go process.Do(srcCh)

	// print to console
	for {
		// get next row
		row, ok := <-record
		if !ok {
			break
		}

		fmt.Printf("%+q\n", row)
	}
}

func warning(message string) {
	fmt.Printf("⚠️  %s\n", message)
	os.Exit(1)
}
