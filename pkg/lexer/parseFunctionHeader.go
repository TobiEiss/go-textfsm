package lexer

import (
	"fmt"

	"github.com/TobiEiss/go-textfsm/pkg/models"
)

// I've no idea why parse comments..
func (parser *Parser) parseStateHeader(stateHeader string) (*models.AbstractStatement, error) {
	statement := &models.AbstractStatement{Type: models.StateHeader, StateName: stateHeader}

	// iterate
	token, val := parser.scan()
	if token != EOF {
		return statement, fmt.Errorf("thought '%s' is a function header - I've no idea what '%s' is", stateHeader, val)
	}
	return statement, nil
}
