package lexer

import (
	"github.com/TobiEiss/go-textfsm/pkg/models"
)

func (parser *Parser) parseVal() (*models.AbstractStatement, error) {
	statement := &models.AbstractStatement{}
	// Next have to be Variable Option or a variable name
	identToken, variable := parser.scanIgnoreWhitespace()
	// If its FILLDOWN -> next ident
	if identToken == FILLDOWN {
		statement.Filldown = true
		identToken, variable = parser.scanIgnoreWhitespace()
	}
	if identToken == LIST {
		statement.List = true
		identToken, variable = parser.scanIgnoreWhitespace()
	}
	if identToken == REQUIRED {
		statement.Required = true
		identToken, variable = parser.scanIgnoreWhitespace()
	}
	if identToken == EOF {
		return nil, &Error{ErrorType: MISSINGARGUMENT, CurrentLine: parser.currentline, ErrorToken: variable}
	}
	if identToken != IDENT {
		return nil, &Error{ErrorType: ILLEGALTOKEN, CurrentLine: parser.currentline, ErrorToken: variable}
	}

	// Next have to be BRACKETLEFT
	if bracketleftToken, val := parser.scanIgnoreWhitespace(); bracketleftToken != BRACKETLEFT {
		return nil, &Error{ErrorType: MISSINGARGUMENT, CurrentLine: parser.currentline, ErrorToken: val}
	}

	// Now the regex
	var bracketCounter = 0
	regex := ""
	for {

		token, val := parser.scan()
		if token == EOF || token == ILLEGAL {
			return nil, &Error{ErrorType: ILLEGALTOKEN, CurrentLine: parser.currentline, ErrorToken: val}
		} else if token == BRACKETLEFT {
			bracketCounter++
		} else if token == BRACKETRIGHT {
			if bracketCounter == 0 {
				break
			}
			bracketCounter--
		}
		regex += val
	}

	statement.Type = models.Value
	statement.VariableName = variable
	statement.Regex = regex

	return statement, nil
}
