package ast

import (
	"github.com/TobiEiss/go-textfsm/pkg/lexer"
	"github.com/TobiEiss/go-textfsm/pkg/models"
)

// CreateAST get a chan with lines to parse
func CreateAST(lines chan string) (models.AST, error) {
	// create ast
	ast := models.AST{
		States: []models.State{},
		Vals:   []models.Val{},
	}

	// this is the "current recording function"
	currentState := models.State{Commands: []models.Cmd{}}

	// iterate all lines
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
				currentState.Commands = append(currentState.Commands, as.Command())
			case models.StateHeader:
				if len(currentState.Commands) > 0 {
					ast.States = append(ast.States, currentState)
				}
				currentState = models.State{Commands: []models.Cmd{}, Name: as.StateName}
			}
		}
	}

	// finally add currentState to states
	ast.States = append(ast.States, currentState)

	// if there is only the "Start"-state, last command have be with a record
	if len(ast.States) == 1 {
		ast.States[0].Commands[len(currentState.Commands)-1].Record = true
	}

	return ast, nil
}
