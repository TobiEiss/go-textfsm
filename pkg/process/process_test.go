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
		ExpectedHeader   []string
		ExpectedRows     [][]interface{}
	}{
		{ // index 0
			TemplateFilePath: "/../../testfiles/01.txt",
			SourceFilePath:   "/../../testfiles/src01.txt",
			ExpectedHeader:   []string{"Year", "Time", "Timezone", "Month", "MonthDay"},
			ExpectedRows: [][]interface{}{
				{"2009", "18:42:41", "PST", "Feb", "8"},
			},
		},
		{ // index 1
			TemplateFilePath: "/../../testfiles/02.txt",
			SourceFilePath:   "/../../testfiles/src02.txt",
			ExpectedHeader:   []string{"ResetReason", "Version", "Uptime", "ConfigRegister"},
			ExpectedRows: [][]interface{}{
				{"Reload", "12.2(31)SGA1", "11 weeks, 4 days, 20 hours, 26 minutes", "0x2102"},
			},
		},
		{ // index 2
			TemplateFilePath: "/../../testfiles/03.txt",
			SourceFilePath:   "/../../testfiles/src03.txt",
			ExpectedHeader:   []string{"Slot", "State", "Temperature"},
			ExpectedRows: [][]interface{}{
				{"0", "Online", "24"},
				{"1", "Online", "25"},
				{"2", "Online", "24"},
				{"3", "Online", "23"},
			},
		},
		{ // index 3
			TemplateFilePath: "/../../testfiles/04.txt",
			SourceFilePath:   "/../../testfiles/src04.txt",
			ExpectedHeader:   []string{"Slot", "State", "Temperature"},
			ExpectedRows: [][]interface{}{
				{"0", "Online", "24"},
				{"1", "Online", "25"},
				{"2", "Online", "24"},
				{"3", "Online", "23"},
				{"4", "Empty", ""},
				{"5", "Empty", ""},
				{"6", "Empty", ""},
				{"7", "Empty", ""},
			},
		},
		{ // index 4
			TemplateFilePath: "/../../testfiles/05.txt",
			SourceFilePath:   "/../../testfiles/src05.txt",
			ExpectedHeader:   []string{"Chassis", "Slot", "State", "Temperature", "DRAM", "Buffer"},
			ExpectedRows: [][]interface{}{
				{"lcc0-re0", "0", "Online", "24", "512", "52"},
				{"", "1", "Online", "23", "256", "53"},
				{"", "2", "Online", "23", "256", "49"},
				{"", "3", "Online", "21", "256", "49"},
				{"", "4", "Empty", "", "", ""},
				{"", "5", "Empty", "", "", ""},
				{"", "6", "Empty", "", "", ""},
				{"", "7", "Empty", "", "", ""},
				{"lcc1-re1", "0", "Online", "20", "256", "50"},
				{"", "1", "Online", "20", "256", "49"},
				{"", "2", "Online", "21", "256", "49"},
				{"", "3", "Online", "20", "256", "49"},
				{"", "4", "Online", "18", "256", "49"},
				{"", "5", "Empty", "", "", ""},
				{"", "6", "Empty", "", "", ""},
				{"", "7", "Empty", "", "", ""},
			},
		},
		{ // index 5
			TemplateFilePath: "/../../testfiles/06.txt",
			SourceFilePath:   "/../../testfiles/src06.txt",
			ExpectedHeader:   []string{"Chassis", "Slot", "State", "Temperature", "DRAM", "Buffer"},
			ExpectedRows: [][]interface{}{
				{"lcc0-re0", "0", "Online", "24", "512", "52"},
				{"lcc0-re0", "1", "Online", "23", "256", "53"},
				{"lcc0-re0", "2", "Online", "23", "256", "49"},
				{"lcc0-re0", "3", "Online", "21", "256", "49"},
				{"lcc0-re0", "4", "Empty", "", "", ""},
				{"lcc0-re0", "5", "Empty", "", "", ""},
				{"lcc0-re0", "6", "Empty", "", "", ""},
				{"lcc0-re0", "7", "Empty", "", "", ""},
				{"lcc1-re1", "0", "Online", "20", "256", "50"},
				{"lcc1-re1", "1", "Online", "20", "256", "49"},
				{"lcc1-re1", "2", "Online", "21", "256", "49"},
				{"lcc1-re1", "3", "Online", "20", "256", "49"},
				{"lcc1-re1", "4", "Online", "18", "256", "49"},
				{"lcc1-re1", "5", "Empty", "", "", ""},
				{"lcc1-re1", "6", "Empty", "", "", ""},
				{"lcc1-re1", "7", "Empty", "", "", ""},
			},
		},
		{ // index 6
			TemplateFilePath: "/../../testfiles/07.txt",
			SourceFilePath:   "/../../testfiles/src07.txt",
			ExpectedHeader:   []string{"Name", "Status", "Protocol"},
			ExpectedRows: [][]interface{}{
				{"Gi0/1", "up", []string{"tcp", "udp", "arp"}},
				{"Gi0/2", "down", []string{"https", "udp", "bgp"}},
				{"Gi0/3", "down", []string{"tcp", "udp", "ospf"}},
				{"Gi0/4", "up", []string{"ip", "http", "rip"}},
			},
		},
		{ // index 7
			TemplateFilePath: "/../../testfiles/08.txt",
			SourceFilePath:   "/../../testfiles/src08.txt",
			ExpectedHeader:   []string{"Protocol", "Type", "Prefix", "Gateway", "Distance", "Metric", "LastChange"},
			ExpectedRows: [][]interface{}{
				{"B", "EX", "0.0.0.0/0", "192.0.2.73", "20", "100", "4w0d"},
				{"B", "IN", "192.0.2.76/30", "203.0.113.183", "200", "100", "4w2d"},
				{"B", "IN", "192.0.2.204/30", "203.0.113.183", "200", "100", "4w2d"},
				{"B", "IN", "192.0.2.80/30", "203.0.113.183", "200", "100", "4w2d"},
				{"B", "IN", "192.0.2.208/30", "203.0.113.183", "200", "100", "4w2d"},
			},
		},
		{ // index 8
			TemplateFilePath: "/../../testfiles/09.txt",
			SourceFilePath:   "/../../testfiles/src09.txt",
			ExpectedHeader:   []string{"Protocol", "Type", "Prefix", "Gateway", "Distance", "Metric", "LastChange"},
			ExpectedRows: [][]interface{}{
				{"B", "EX", "0.0.0.0/0", []string{"192.0.2.73"}, "20", "100", "4w0d"},
				{"B", "IN", "192.0.2.76/30", []string{"203.0.113.183"}, "200", "100", "4w2d"},
				{"B", "IN", "192.0.2.204/30", []string{"203.0.113.183"}, "200", "100", "4w2d"},
				{"B", "IN", "192.0.2.80/30", []string{"203.0.113.183"}, "200", "100", "4w2d"},
				{"B", "IN", "192.0.2.208/30", []string{"203.0.113.183"}, "200", "100", "4w2d"},
			},
		},
		{ // index 9
			TemplateFilePath: "/../../testfiles/10a.txt",
			SourceFilePath:   "/../../testfiles/src10.txt",
			ExpectedHeader:   []string{"Iface", "Name", "Status", "Error"},
			ExpectedRows: [][]interface{}{
				{"Gi0/1", "wan1", "up", "output Queue errors"},
				{"Gi0/2", "wan2", "down", "input Queue errors"},
				{"Gi0/3", "inside", "down", "input Queue errors"},
				{"Gi0/6", "dmz", "down", "input Queue errors"},
			},
		},
		{ // index 10
			TemplateFilePath: "/../../testfiles/10b.txt",
			SourceFilePath:   "/../../testfiles/src10.txt",
			ExpectedHeader:   []string{"Iface", "Name", "Status", "Error"},
			ExpectedRows: [][]interface{}{
				{"Gi0/1", "wan1", "up", "output Queue errors"},
				{"Gi0/2", "wan2", "down", "input Queue errors"},
				{"Gi0/3", "inside", "down", "input Queue errors"},
				{"Gi0/4", "", "up", "output Queue errors"},
				{"Gi0/6", "dmz", "down", "input Queue errors"},
			},
		},
		{ // index 11
			TemplateFilePath: "/../../testfiles/11.txt",
			SourceFilePath:   "/../../testfiles/src11.txt",
			ExpectedHeader:   []string{"Port", "Name", "Status", "Vlan", "Duplex", "Speed", "Type"},
			ExpectedRows: [][]interface{}{
				{"Gi1/0/2", "AccessPoint", "connected", "8", "a-full", "a-1000", "10/100/1000BaseTX"},
				{"Gi1/0/3", "John's Office", "notconnect", "1", "auto", "auto", "10/100/1000BaseTX"},
				{"Gi1/0/4", "SingleName", "connected", "1", "a-full", "a-100", "10/100/1000BaseTX"},
			},
		},
		{ // index 12
			TemplateFilePath: "/../../testfiles/12.txt",
			SourceFilePath:   "/../../testfiles/src12.txt",
			ExpectedHeader:   []string{"Port", "Name", "Status", "Vlan", "Duplex", "Speed", "Type"},
			ExpectedRows: [][]interface{}{
				{"Gi1/0/1", "Cpu1", "notconnect", "1", "auto", "auto", "10/100/1000BaseTX"},
				{"Gi1/0/2", "AccessPoint", "connected", "8", "a-full", "a-1000", "10/100/1000BaseTX"},
				{"Gi1/0/3", "John's", "notconnect", "1", "auto", "auto", "10/100/1000BaseTX"},
				{"Gi1/0/4", "SingleName", "connected", "1", "a-full", "a-100", "10/100/1000BaseTX"},
			},
		},
		{ // index 13
			TemplateFilePath: "/../../testfiles/13a.txt",
			SourceFilePath:   "/../../testfiles/src13.txt",
			ExpectedHeader:   []string{"Ifname", "Name", "Status", "Index"},
			ExpectedRows: [][]interface{}{
				{"Gi0/1", "", "up", "101"},
				{"Gi0/2", "", "down", "898"},
				{"Gi0/3", "", "down", "666"},
				{"Gi0/4", "", "up", "999"},
				{"Gi0/6", "", "down", "100"},
			},
		},
		{ // index 14
			TemplateFilePath: "/../../testfiles/13b.txt",
			SourceFilePath:   "/../../testfiles/src13.txt",
			ExpectedHeader:   []string{"Ifname", "Name", "Status", "Index"},
			ExpectedRows: [][]interface{}{
				{"", "", "up", "101"},
				{"", "", "down", "898"},
				{"", "", "down", "666"},
				{"", "", "up", "999"},
				{"", "", "down", "100"},
			},
		},
		{ // index 15
			TemplateFilePath: "/../../testfiles/13c.txt",
			SourceFilePath:   "/../../testfiles/src13.txt",
			ExpectedHeader:   []string{"Ifname", "Name", "Status", "Index"},
			ExpectedRows: [][]interface{}{
				{"", "", "up", "101"},
				{"", "", "down", "898"},
				{"", "", "down", "666"},
				{"", "", "up", "999"},
				{"", "", "down", "100"},
			},
		},
		{ // index 16
			TemplateFilePath: "/../../testfiles/14.txt",
			SourceFilePath:   "/../../testfiles/src14.txt",
			ExpectedHeader:   []string{"Interface", "Description", "UnnumInterface", "Destination"},
			ExpectedRows: [][]interface{}{
				{"5455", "RTRA->RTRB->LAB", "Loopback65", "8.8.8.8"},
			},
		},
		{ // index 17
			TemplateFilePath: "/../../testfiles/15.txt",
			SourceFilePath:   "/../../testfiles/src15.txt",
			ExpectedHeader:   []string{"FirstValue", "SecondValue"},
			ExpectedRows: [][]interface{}{
				{"100", "2"},
			},
		},
		{ // index 17
			TemplateFilePath: "/../../testfiles/16.txt",
			SourceFilePath:   "/../../testfiles/src16.txt",
			ExpectedHeader:   []string{"FirstValue", "SecondValue"},
			ExpectedRows: [][]interface{}{
				{"100", "2"},
			},
		},
	}

	// iterate all test.cases
	for testIndex, test := range tests {
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
		record := make(chan []interface{})
		process, err := process.NewProcess(ast, record)
		if err != nil {
			t.Error(err)
		}
		go process.Do(srcCh)

		// check
		indexRow := 0
		for {
			// get next row
			row, ok := <-record
			if !ok {
				break
			}
			for colIndex, colHeader := range test.ExpectedHeader {
				_, astIndex := ast.GetValForValName(colHeader)

				if reflect.TypeOf(test.ExpectedRows[indexRow][colIndex]).Kind() == reflect.Slice {
					if !isEqual(row[astIndex].([]string), test.ExpectedRows[indexRow][colIndex].([]string)) {
						t.Errorf("%d failed: Field '%s' in row %d with value '%s' is not equal expected '%s'",
							testIndex, colHeader, indexRow, row[astIndex], test.ExpectedRows[indexRow][colIndex])

					}
				} else {
					if row[astIndex] != test.ExpectedRows[indexRow][colIndex] {
						t.Errorf("%d failed: Field '%s' in row %d with value '%s' is not equal expected '%s'",
							testIndex, colHeader, indexRow, row[astIndex], test.ExpectedRows[indexRow][colIndex])
					}
				}
			}
			indexRow++
		}

		if indexRow != len(test.ExpectedRows) {
			t.Errorf("%d failed: %d expected rows - but here are %d rows", testIndex, len(test.ExpectedRows), indexRow)
		}
	}

}

func TestInline(t *testing.T) {
	tests := []struct {
		template       string
		source         string
		ExpectedHeader []string
		ExpectedRows   [][]interface{}
	}{
		{
			template: `Value Year (\d+)
			Value MonthDay (\d+)
			Value Month (\w+)
			Value Timezone (\S+)
			Value Time (..:..:..)
			
			Start
			  ^${Time}.* ${Timezone} \w+ ${Month} ${MonthDay} ${Year} -> Record`,
			source:         `18:42:41.321 PST Sun Feb 8 2009`,
			ExpectedHeader: []string{"Year", "Time", "Timezone", "Month", "MonthDay"},
			ExpectedRows: [][]interface{}{
				{"2009", "18:42:41", "PST", "Feb", "8"},
			},
		},
	}

	for testIndex, test := range tests {
		tmplCh := make(chan string)
		go reader.ReadLineByLineFileAsString(test.template, tmplCh)

		srcCh := make(chan string)
		go reader.ReadLineByLineFileAsString(test.template, srcCh)

		// create AST
		ast, err := ast.CreateAST(tmplCh)
		if err != nil {
			t.Error(err)
		}

		// process ast
		record := make(chan []interface{})
		process, err := process.NewProcess(ast, record)
		if err != nil {
			t.Error(err)
		}
		go process.Do(srcCh)

		// check
		indexRow := 0
		for {
			// get next row
			row, ok := <-record
			if !ok {
				break
			}
			for colIndex, colHeader := range test.ExpectedHeader {
				_, astIndex := ast.GetValForValName(colHeader)

				if reflect.TypeOf(test.ExpectedRows[indexRow][colIndex]).Kind() == reflect.Slice {
					if !isEqual(row[astIndex].([]string), test.ExpectedRows[indexRow][colIndex].([]string)) {
						t.Errorf("%d failed: Field '%s' in row %d with value '%s' is not equal expected '%s'",
							testIndex, colHeader, indexRow, row[astIndex], test.ExpectedRows[indexRow][colIndex])

					}
				} else {
					if row[astIndex] != test.ExpectedRows[indexRow][colIndex] {
						t.Errorf("%d failed: Field '%s' in row %d with value '%s' is not equal expected '%s'",
							testIndex, colHeader, indexRow, row[astIndex], test.ExpectedRows[indexRow][colIndex])
					}
				}
			}
			indexRow++
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
