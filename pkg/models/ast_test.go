package models_test

import (
	"testing"

	"github.com/TobiEiss/go-textfsm/pkg/models"
)

func TestASTMatchingLine(t *testing.T) {
	ast := models.AST{
		Vals: []models.Val{
			{
				Regex:    `.*`,
				Variable: "MyVal1",
			},
		},
		States: []models.State{
			{
				Commands: []models.Cmd{
					{
						Actions: []models.Action{
							{Regex: "./"},
							{Value: "MyVal1"},
						},
					},
				},
			},
		},
	}

	line, err := ast.CreateMatchingLine(ast.States[0].Commands[0])
	if err != nil {
		t.Error(err)
	}

	if line != "./(?P<MyVal1>.*)$" {
		t.Errorf("matchingline is not expected. It is: %s", line)
	}
}
