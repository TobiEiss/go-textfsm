package models

import (
	"errors"
	"fmt"
)

// AST is the abstract command tree
type AST struct {
	Vals     []Val
	Commands []Cmd
}

// Val represent a varible like "Value Year (\d+)"
type Val struct {
	Variable string
	Regex    string
}

// Cmd is one statement after keyword "Start"
type Cmd struct {
	Actions []Action
	Record  string
}

// GetValForValName searches a val for a valName
func (ast AST) GetValForValName(valName string) *Val {
	for _, val := range ast.Vals {
		if val.Variable == valName {
			return &val
		}
	}
	return nil
}

// CreateMatchingLine creates a regex that have to match to a line
func (ast AST) CreateMatchingLine(cmd Cmd) (matchingLine string, err error) {
	// iterate all actions
	for _, action := range cmd.Actions {
		if action.Value != "" {
			if val := ast.GetValForValName(action.Value); val != nil {
				matchingLine += fmt.Sprintf(`(?P<%s>%s)`, val.Variable, val.Regex)
				continue
			}
			return matchingLine, errors.New("Can't find val for ValName" + action.Value)
		}
		matchingLine += action.Regex
	}
	matchingLine += "$"
	return
}
