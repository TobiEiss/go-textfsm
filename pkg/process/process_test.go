package process_test

import (
	"path/filepath"
	"reflect"
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
		{
			TemplateFilePath: "/../../testfiles/04.txt",
			SourceFilePath:   "/../../testfiles/src04.txt",
			CorrectRecord: map[string]process.Column{
				"Slot":        process.Column{Entries: []interface{}{"0", "1", "2", "3", "4", "5", "6", "7"}},
				"State":       process.Column{Entries: []interface{}{"Online", "Online", "Online", "Online", "Empty", "Empty", "Empty", "Empty"}},
				"Temperature": process.Column{Entries: []interface{}{"24", "25", "24", "23", "", "", "", ""}},
			},
		},
		{
			TemplateFilePath: "/../../testfiles/05.txt",
			SourceFilePath:   "/../../testfiles/src05.txt",
			CorrectRecord: map[string]process.Column{
				"Chassis":     process.Column{Entries: []interface{}{"lcc0-re0", "", "", "", "", "", "", "", "lcc1-re1", "", "", "", "", "", "", ""}},
				"Slot":        process.Column{Entries: []interface{}{"0", "1", "2", "3", "4", "5", "6", "7", "0", "1", "2", "3", "4", "5", "6", "7"}},
				"State":       process.Column{Entries: []interface{}{"Online", "Online", "Online", "Online", "Empty", "Empty", "Empty", "Empty", "Online", "Online", "Online", "Online", "Online", "Empty", "Empty"}},
				"Temperature": process.Column{Entries: []interface{}{"24", "23", "23", "21", "", "", "", "", "20", "20", "21", "20", "18", "", "", ""}},
			},
		},
		{
			TemplateFilePath: "/../../testfiles/06.txt",
			SourceFilePath:   "/../../testfiles/src06.txt",
			CorrectRecord: map[string]process.Column{
				"Chassis":     process.Column{Entries: []interface{}{"lcc0-re0", "lcc0-re0", "lcc0-re0", "lcc0-re0", "lcc0-re0", "lcc0-re0", "lcc0-re0", "lcc0-re0", "lcc1-re1", "lcc1-re1", "lcc1-re1", "lcc1-re1", "lcc1-re1", "lcc1-re1", "lcc1-re1", "lcc1-re1"}},
				"Slot":        process.Column{Entries: []interface{}{"0", "1", "2", "3", "4", "5", "6", "7", "0", "1", "2", "3", "4", "5", "6", "7"}},
				"State":       process.Column{Entries: []interface{}{"Online", "Online", "Online", "Online", "Empty", "Empty", "Empty", "Empty", "Online", "Online", "Online", "Online", "Online", "Empty", "Empty"}},
				"Temperature": process.Column{Entries: []interface{}{"24", "23", "23", "21", "", "", "", "", "20", "20", "21", "20", "18", "", "", ""}},
			},
		},
		{
			TemplateFilePath: "/../../testfiles/07.txt",
			SourceFilePath:   "/../../testfiles/src07.txt",
			CorrectRecord: map[string]process.Column{
				"Name":     process.Column{Entries: []interface{}{"Gi0/1", "Gi0/2", "Gi0/3", "Gi0/4"}},
				"Status":   process.Column{Entries: []interface{}{"up", "down", "down", "up"}},
				"Protocol": process.Column{Entries: []interface{}{[]string{"tcp", "udp", "arp"}, []string{"https", "udp", "bgp"}, []string{"tcp", "udp", "ospf"}, []string{"ip", "http", "rip"}}},
			},
		},
		{
			TemplateFilePath: "/../../testfiles/08.txt",
			SourceFilePath:   "/../../testfiles/src08.txt",
			CorrectRecord: map[string]process.Column{
				"Protocol":   process.Column{Entries: []interface{}{"B", "B", "B", "B", "B"}},
				"Type":       process.Column{Entries: []interface{}{"EX", "IN", "IN", "IN", "IN"}},
				"Prefix":     process.Column{Entries: []interface{}{"0.0.0.0/0", "192.0.2.76/30", "192.0.2.204/30", "192.0.2.80/30", "192.0.2.208/30"}},
				"Gateway":    process.Column{Entries: []interface{}{"192.0.2.73", "203.0.113.183", "203.0.113.183", "203.0.113.183", "203.0.113.183"}},
				"Distance":   process.Column{Entries: []interface{}{"20", "200", "200", "200", "200"}},
				"Metric":     process.Column{Entries: []interface{}{"100", "100", "100", "100", "100"}},
				"LastChange": process.Column{Entries: []interface{}{"4w0d", "4w2d", "4w2d", "4w2d", "4w2d"}},
			},
		},
		{
			TemplateFilePath: "/../../testfiles/09.txt",
			SourceFilePath:   "/../../testfiles/src09.txt",
			CorrectRecord: map[string]process.Column{
				"Protocol":   process.Column{Entries: []interface{}{"B", "B", "B", "B", "B"}},
				"Type":       process.Column{Entries: []interface{}{"EX", "IN", "IN", "IN", "IN"}},
				"Prefix":     process.Column{Entries: []interface{}{"0.0.0.0/0", "192.0.2.76/30", "192.0.2.204/30", "192.0.2.80/30", "192.0.2.208/30"}},
				"Gateway":    process.Column{Entries: []interface{}{[]string{"192.0.2.73"}, []string{"203.0.113.183"}, []string{"203.0.113.183"}, []string{"203.0.113.183"}, []string{"203.0.113.183"}}},
				"Distance":   process.Column{Entries: []interface{}{"20", "200", "200", "200", "200"}},
				"Metric":     process.Column{Entries: []interface{}{"100", "100", "100", "100", "100"}},
				"LastChange": process.Column{Entries: []interface{}{"4w0d", "4w2d", "4w2d", "4w2d", "4w2d"}},
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
			if len(record[k].Entries) < len(v.Entries) {
				t.Errorf("%d failed: len of values (%d) for '%s' are less than expected (%d)",
					index, len(record[k].Entries), k, len(v.Entries))
				break
			}

			// check if erntries are correct
			for entryIndex, entrie := range v.Entries {
				if reflect.TypeOf(entrie).Kind() == reflect.Slice {
					if !isEqual(record[k].Entries[entryIndex].([]string), entrie.([]string)) {
						t.Errorf("%d failed: Field '%s' Value '%s' is not equal expected '%+v'",
							index, k, record[k].Entries[entryIndex], entrie)

					}
				} else {
					if record[k].Entries[entryIndex] != entrie {
						t.Errorf("%d failed: Field '%s' Value '%s' is not equal expected '%+v'",
							index, k, record[k].Entries[entryIndex], entrie)
					}
				}
			}
		}
	}

}

func isEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
