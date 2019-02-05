package ast_test

import (
	"testing"

	"github.com/TobiEiss/go-textfsm/pkg/ast"
)

func TestAST(t *testing.T) {
	var tests = []struct {
		Lines []string
	}{
		{Lines: []string{
			`Value Year (\d+)`,
			`Value Time (..:..:..)`},
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
		if len(ast.Vals) != len(test.Lines) {
			t.Errorf("%d failed: len of ast.vals '%d' is not equal expected len '%d'",
				index, len(ast.Vals), len(test.Lines))
		}
	}
}
