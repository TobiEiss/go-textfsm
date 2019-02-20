package reader_test

import (
	"errors"
	"log"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/TobiEiss/go-textfsm/pkg/ast"
	"github.com/TobiEiss/go-textfsm/pkg/reader"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func TestReader(t *testing.T) {
	filepath := basepath + "/../../testfiles/01.txt"
	ch := make(chan string)
	go reader.ReadLineByLine(filepath, ch)

	for {
		line, ok := <-ch
		if !ok {
			break
		}
		log.Println(line)
	}
}

func TestReaderWithCreatingAST(t *testing.T) {
	filepath := basepath + "/../../testfiles/02.txt"
	ch := make(chan string)
	go reader.ReadLineByLine(filepath, ch)

	// create AST
	ast, err := ast.CreateAST(ch)
	if err != nil {
		t.Error(err)
	}

	// check
	if len(ast.Vals) != 5 {
		t.Error(errors.New("ast.vals has not length of 5 like expected"))
	}
	if len(ast.Commands[0].Actions) != 12 {
		t.Error(errors.New("actions has not length of 7 like expected"))
	}
}
