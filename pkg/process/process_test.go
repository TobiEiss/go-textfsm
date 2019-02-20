package process_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/TobiEiss/go-textfsm/pkg/ast"
	"github.com/TobiEiss/go-textfsm/pkg/process"
	"github.com/TobiEiss/go-textfsm/pkg/reader"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func TestProcessAST(t *testing.T) {
	// read template
	filepath := basepath + "/../../testfiles/02.txt"
	tmplCh := make(chan string)
	go reader.ReadLineByLine(filepath, tmplCh)

	// read file
	filepathSrc := basepath + "/../../testfiles/src01.txt"
	srcCh := make(chan string)
	go reader.ReadLineByLine(filepathSrc, srcCh)

	// create AST
	ast, err := ast.CreateAST(tmplCh)
	if err != nil {
		t.Error(err)
	}

	// process ast
	process, err := process.NewProcess(ast)
	if err != nil {
		t.Error(err)
	}

	if !process.Do(srcCh) {
		t.Error("can't find matching line")
	}
}
