package models

import (
	"errors"
	"fmt"
)

// AST is the abstract command tree
type AST struct {
	Vals   []Val
	States []State
}

// Val represent a varible like "Value Year (\d+)"
type Val struct {
	Variable string
	Regex    string
	Filldown bool
	List     bool
	Required bool
}

// State represent a state like "Start"
type State struct {
	Name     string
	Commands []Cmd
}

// Cmd is one statement after keyword "Start"
type Cmd struct {
	Actions   []Action
	Vals      []*Val
	Record    bool
	StateCall string
}

// GetValForValName and index searches a val for a valName
func (ast AST) GetValForValName(valName string) (*Val, int) {
	for index, val := range ast.Vals {
		if val.Variable == valName {
			return &val, index
		}
	}
	return nil, -1
}

// CreateMatchingLine creates a regex that have to match to a line
func (ast AST) CreateMatchingLine(cmd Cmd) (matchingLine string, err error) {
	cmd.Vals = []*Val{}
	// iterate all actions
	for _, action := range cmd.Actions {
		if action.Value != "" {
			if val, _ := ast.GetValForValName(action.Value); val != nil {
				matchingLine += fmt.Sprintf(`(?P<%s>%s)`, val.Variable, val.Regex)
				cmd.Vals = append(cmd.Vals, val)
				continue
			}
			return matchingLine, errors.New("Can't find val for ValName" + action.Value)
		}
		matchingLine += action.Regex
	}
	return
}
