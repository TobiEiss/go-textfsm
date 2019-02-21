package ast

import (
	"github.com/TobiEiss/go-textfsm/pkg/lexer"
	"github.com/TobiEiss/go-textfsm/pkg/models"
)

// CreateAST get a chan with lines to parse
func CreateAST(lines chan string) (models.AST, error) {
	// create ast
	ast := models.AST{
		Commands: []models.Cmd{},
		Vals:     []models.Val{},
	}

	for {
		// get next line
		line, ok := <-lines
		if !ok {
			break
		}

		// check what is line
		parser := lexer.NewParser(line)
		as, err := parser.ParseStatement()
		if err != nil {
			return ast, err
		}

		// "as" can be nil, if there is a ""-line (empty)
		if as != nil {
			switch as.Type {
			case models.Value:
				ast.Vals = append(ast.Vals, as.Value())
			case models.Command:
				ast.Commands = append(ast.Commands, as.Command())
			}
		}
	}

	return ast, nil
}
