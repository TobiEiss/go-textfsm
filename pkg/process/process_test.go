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
		CorrectRecord    map[string]process.Column
	}{
		{
			TemplateFilePath: "/../../testfiles/01.txt",
			SourceFilePath:   "/../../testfiles/src01.txt",
			CorrectRecord: map[string]process.Column{
				"Year":     process.Column{Entries: []interface{}{"2009"}},
				"Time":     process.Column{Entries: []interface{}{"18:42:41"}},
				"Timezone": process.Column{Entries: []interface{}{"PST"}},
				"Month":    process.Column{Entries: []interface{}{"Feb"}},
				"MonthDay": process.Column{Entries: []interface{}{"8"}},
			},
		},
		{
			TemplateFilePath: "/../../testfiles/02.txt",
			SourceFilePath:   "/../../testfiles/src02.txt",
			CorrectRecord: map[string]process.Column{
				"ResetReason":    process.Column{Entries: []interface{}{"Reload"}},
				"Version":        process.Column{Entries: []interface{}{"12.2(31)SGA1"}},
				"Uptime":         process.Column{Entries: []interface{}{"11 weeks, 4 days, 20 hours, 26 minutes"}},
				"ConfigRegister": process.Column{Entries: []interface{}{"0x2102"}},
			},
		},
		{
			TemplateFilePath: "/../../testfiles/03.txt",
			SourceFilePath:   "/../../testfiles/src03.txt",
			CorrectRecord: map[string]process.Column{
				"Slot":        process.Column{Entries: []interface{}{"0", "1", "2", "3"}},
				"State":       process.Column{Entries: []interface{}{"Online", "Online", "Online", "Online"}},
				"Temperature": process.Column{Entries: []interface{}{"24", "25", "24", "23"}},
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
			// check if entries are available
			if len(record["Record"][k].Entries) < 1 {
				t.Errorf("%d failed: Values for '%s' are missing", index, k)
			}

			// check if erntries are correct
			for entryIndex, entrie := range v.Entries {
				if record["Record"][k].Entries[entryIndex] != entrie {
					t.Errorf("%d failed: Field '%s' Value '%s' is not equal expected '%+v'",
						index, k, record["Record"][k].Entries[entryIndex], entrie)
				}
			}
		}
	}

}
