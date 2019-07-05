package lexer

import (
	"github.com/TobiEiss/go-textfsm/pkg/models"
)

func (parser *Parser) parseCmd() (*models.AbstractStatement, error) {
	statement := &models.AbstractStatement{Type: models.Command, Actions: []models.Action{}}

	// function to parse a value-call
	parseValue := func(parser *Parser) (valname string, err error) {
		// need a "{""
		if token, _ := parser.scan(); token != CURLYBRACKETLEFT && token != DOLAR {
			return "", parser.createError(ILLEGALTOKEN)
		}

		for {
			token, val := parser.scan()
			switch token {
			case EOF:
				err = parser.createError(ILLEGALTOKEN)
				return
			case CURLYBRACKETRIGHT:
				return
			default:
				valname += val
			}
		}
	}

	// function to parse a "normal regex"
	parseRegex := func(startIdent string, parser *Parser) (regex string) {
		regex += startIdent
		for {
			token, val := parser.scan()
			switch token {
			case DOLAR:
				parser.unscan()
				return
			case WHITESPACE:
				parser.unscan()
				return
			case EOF:
				return
			default:
				regex += val
			}
		}
	}

	// iterate
	for {
		token, val := parser.scan()
		switch token {
		case DOLAR:
			// check if EOL (i.e. $$) is the case.
			if token, _ := parser.scan(); token == DOLAR {
				statement.Actions = append(statement.Actions, models.Action{Regex: parseRegex(val, parser)})
			} else {
				parser.unscan()
				valueName, err := parseValue(parser)
				if err != nil {
					return statement, err
				}
				statement.Actions = append(statement.Actions, models.Action{Value: valueName})
			}
		case EOF:
			return statement, nil
		case MINUS:
			// check if after "-" comes a ">"
			if token, val := parser.scan(); token != BIGGER {
				statement.Actions = append(statement.Actions, models.Action{Regex: parseRegex("-"+val, parser)})
			}

			// check if after "->" comes whitespace
			if token, val := parser.scan(); token != WHITESPACE {
				statement.Actions = append(statement.Actions, models.Action{Regex: parseRegex("->"+val, parser)})
			}

			// delete last WHITESPACE if there is one
			if statement.Actions[len(statement.Actions)-1].Regex == " " {
				statement.Actions = statement.Actions[:len(statement.Actions)-1]
			}

			for {
				token, val := parser.scan()
				switch token {
				case IDENT:
					statement.StateCall = val
				case RECORD:
					statement.Record = true
				case CONTINUE:
					statement.Continue = true
				case EOF:
					return statement, nil
				}
			}
		case WHITESPACE:
			statement.Actions = append(statement.Actions, models.Action{Regex: " "})
		default:
			statement.Actions = append(statement.Actions, models.Action{Regex: parseRegex(val, parser)})
		}
	}
}
