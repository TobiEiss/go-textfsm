package ast_test

import (
	"testing"

	"github.com/TobiEiss/go-textfsm/pkg/ast"
)

func TestAST(t *testing.T) {
	var tests = []struct {
		Lines            []string
		ExpectedVals     int
		ExpectedCommands int
	}{
		{
			Lines: []string{
				`Value Year (\d+)`,
				`Value Time (..:..:..)`,
				`Start`,
				`^${Time}.* ${Timezone} \w+ ${Month} ${MonthDay} ${Year} -> Record`},
			ExpectedVals:     2,
			ExpectedCommands: 1,
		},
	}

	// iterate all tests
	for index, test := range tests {
		// create chan of lines
		lines := make(chan string)
		go func() {
			for _, line := range test.Lines {
				lines <- line
			}
			close(lines)
		}()

		// create an AST
		ast, err := ast.CreateAST(lines)
		if err != nil {
			t.Error(err)
		}

		// check result
		if len(ast.Vals) != test.ExpectedVals {
			t.Errorf("%d failed: len of ast.vals '%d' is not equal expected len '%d'",
				index, len(ast.Vals), test.ExpectedVals)
		}

		// check commands
		if len(ast.Commands) != test.ExpectedCommands {
			t.Errorf("%d failed: len of ast.commands '%d' is not equal expected len '%d'",
				index, len(ast.Commands), test.ExpectedCommands)
		}
	}
}

func TestParseCommands(t *testing.T) {
	var tests = []struct {
		Command         string
		ExpectedActions []string
	}{
		{
			Command:         `^${Time}.* ${Timezone} \w+ ${Month} ${MonthDay} ${Year} -> Record`,
			ExpectedActions: []string{`^${Time}`, `.*`, `${Timezone}`, `\w+`, `${Month}`, `${MonthDay}`, `${Year}`},
		},
	}

	// iterate all tests
	for index, test := range tests {
		// create chan of lines
		lines := make(chan string)
		go func() {
			lines <- test.Command
			close(lines)
		}()

		// create an AST
		ast, err := ast.CreateAST(lines)
		if err != nil {
			t.Error(err)
		}

		// check
		if len(ast.Commands[0].Actions) != len(test.ExpectedActions) {
			t.Errorf("%d failed: len of actions '%d' is not equal expected len '%d'",
				index, len(ast.Commands[0].Actions), len(test.ExpectedActions))
		}
	}
}
