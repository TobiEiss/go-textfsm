package lexer

import "github.com/TobiEiss/go-textfsm/pkg/models"

func (parser *Parser) parseVal() (*models.AbstractStatement, error) {
	statement := &models.AbstractStatement{}
	// Next have to be a variable name
	identToken, variable := parser.scanIgnoreWhitespace()
	if identToken == EOF {
		return nil, parser.createError(MISSINGARGUMENT)
	}
	if identToken != IDENT {
		return nil, parser.createError(ILLEGALTOKEN)
	}

	// Next have to be BRACKETLEFT
	if bracketleftToken, _ := parser.scanIgnoreWhitespace(); bracketleftToken != BRACKETLEFT {
		return nil, parser.createError(MISSINGARGUMENT)
	}

	// Now the regex
	regex := ""
	for {
		token, val := parser.scanIgnoreWhitespace()
		if token == EOF || token == ILLEGAL {
			return nil, parser.createError(ILLEGALTOKEN)
		} else if token == BRACKETRIGHT {
			break
		}
		regex += val
	}

	statement.Type = models.Value
	statement.VariableName = variable
	statement.Regex = regex

	return statement, nil
}
