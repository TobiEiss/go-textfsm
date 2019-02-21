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
	tests := []struct {
		TemplateFilePath string
		SourceFilePath   string
		CorrectRecord    map[string]interface{}
	}{
		{
			TemplateFilePath: "/../../testfiles/01.txt",
			SourceFilePath:   "/../../testfiles/src01.txt",
			CorrectRecord: map[string]interface{}{
				"Year":     "2009",
				"Time":     "18:42:41",
				"Timezone": "PST",
				"Month":    "Feb",
				"MonthDay": "8",
			},
		},
		{
			TemplateFilePath: "/../../testfiles/02.txt",
			SourceFilePath:   "/../../testfiles/src02.txt",
			CorrectRecord: map[string]interface{}{
				"ResetReason":    "Reload",
				"Version":        "12.2(31)SGA1",
				"Uptime":         "11 weeks, 4 days, 20 hours, 26 minutes",
				"ConfigRegister": "0x2102",
			},
		},
	}

	// iterate all test.cases
	for index, test := range tests {
		// read template
		filepath := basepath + test.TemplateFilePath
		tmplCh := make(chan string)
		go reader.ReadLineByLine(filepath, tmplCh)

		// read file
		filepathSrc := basepath + test.SourceFilePath
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
		for k, v := range test.CorrectRecord {
			if len(record["Record"][k]) < 1 {
				t.Errorf("%d failed: Values for '%s' are missing", index, k)
			}
			if record["Record"][k][0] != v {
				t.Errorf("%d failed: Field '%s' Value '%s' is not equal expected '%s'", index, k, v, record["Record"][k][0])
			}
		}
	}

}
