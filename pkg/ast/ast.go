package ast

import (
	"github.com/TobiEiss/go-textfsm/pkg/lexer"
	"github.com/TobiEiss/go-textfsm/pkg/models"
)

// CreateAST get a chan with lines to parse
func CreateAST(lines chan string) (models.AST, error) {
	ast := models.AST{
		Statements: []models.Statement{},
		Vals:       []models.Val{},
	}

	// collect all vals
	for {
		// get next line
		line, ok := <-lines
		if !ok {
			break
		}

		// try to get val
		val, err := lexer.NewParser(line).ParseVal()
		if err != nil {
			return ast, err
		}
		ast.Vals = append(ast.Vals, *val)
	}

	return ast, nil
}
