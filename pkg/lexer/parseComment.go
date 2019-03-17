package lexer

import "github.com/TobiEiss/go-textfsm/pkg/models"

// I've no idea why parse comments..
func (parser *Parser) parseCmmt() (*models.AbstractStatement, error) {
	statement := &models.AbstractStatement{Type: models.Comment, Comment: ""}

	// iterate
	for {
		token, val := parser.scan()
		if token == EOF {
			return statement, nil
		}
		statement.Comment += val
	}
}
