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
	filepath := basepath + "/../../testfiles/01.txt"
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

	record := process.Do(srcCh)

	// check
	correctRecord := map[string]interface{}{
		"Year":     "2009",
		"Time":     "18:42:41",
		"Timezone": "PST",
		"Month":    "Feb",
		"MonthDay": "8",
	}

	for k, v := range correctRecord {
		if record["Record"][k] != v {
			t.Errorf("'%s' is not expected '%s' - instead it is '%s'", k, v, record[k])
		}
	}
}
