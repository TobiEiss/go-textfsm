# go-textfsm 
[![Build Status](https://travis-ci.org/TobiEiss/go-textfsm.svg?branch=master)](https://travis-ci.org/TobiEiss/go-textfsm)
[![GolangCI](https://golangci.com/badges/github.com/TobiEiss/go-textfsm.svg)](https://golangci.com)

This library is an golang implementation of [TextFSM](https://github.com/google/textfsm).  
If you miss something (there are definitely something) create an issue!

# getting started

Given is following source-file `source.txt`
```
18:42:41.321 PST Sun Feb 8 2009
```

Goal is to extract several values with following template-file `template.txt`
```
Value Year (\d+)
Value MonthDay (\d+)
Value Month (\w+)
Value Timezone (\S+)
Value Time (..:..:..)

Start
  ^${Time}.* ${Timezone} \w+ ${Month} ${MonthDay} ${Year} -> Record
```

You can use the following code to
1. read source-file `source.txt`
2. read template-file `template.txt`
3. process this files
4. print out result

Folder-structure for the example is like this:
```
.
├── source.txt
├── template.txt
├── main.go
```

Here the `main.go`:
```go
package main

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/TobiEiss/go-textfsm/pkg/ast"
	"github.com/TobiEiss/go-textfsm/pkg/process"
	"github.com/TobiEiss/go-textfsm/pkg/reader"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func main() {
	// read template
	filepath := basepath + "/template.txt"
	tmplCh := make(chan string)
	go reader.ReadLineByLine(filepath, tmplCh)

	// read file
	filepathSrc := basepath + "/source.txt"
	srcCh := make(chan string)
	go reader.ReadLineByLine(filepathSrc, srcCh)

	// create AST
	ast, err := ast.CreateAST(tmplCh)
	if err != nil {
		// handle error
	}

	// process ast
	record := make(chan []interface{})
	process, err := process.NewProcess(ast, record)
	if err != nil {
		// handle error
	}
	go process.Do(srcCh)

	// print record
	for {
		// get next row
		row, ok := <-record
		if !ok {
			break
		}

		log.Println(row)
	}
}
```

Find more examples how to build template-files here: [TextFSM-Wiki](https://github.com/google/textfsm/wiki/TextFSM)
