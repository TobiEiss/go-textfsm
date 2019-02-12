package ast_test

import (
	"testing"

	"github.com/TobiEiss/go-textfsm/pkg/ast"
	"github.com/TobiEiss/go-textfsm/pkg/models"
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
		ExpectedActions []models.Action
	}{
		{
			Command: `^${Time}.* ${Timezone} \w+ ${Month} ${MonthDay} ${Year} -> Record`,
			ExpectedActions: []models.Action{
				models.Action{Value: "Time"},
				models.Action{Regex: ".*"},
				models.Action{Value: "Timezone"},
				models.Action{Regex: `\w+`},
				models.Action{Value: "Month"},
				models.Action{Value: "MonthDay"},
				models.Action{Value: "Year"},
			},
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

		// check len of actions
		if len(ast.Commands[0].Actions) != len(test.ExpectedActions) {
			t.Errorf("%d failed: len of actions '%d' is not equal expected len '%d'",
				index, len(ast.Commands[0].Actions), len(test.ExpectedActions))
		}

		// check actions
		for i := 0; i < len(test.ExpectedActions); i++ {
			if test.ExpectedActions[i] != ast.Commands[0].Actions[i] {
				t.Errorf("%d failed: action '%s' is not equal to expected Action '%s'",
					index, test.ExpectedActions[i], ast.Commands[0].Actions[i])
			}
		}
	}
}
