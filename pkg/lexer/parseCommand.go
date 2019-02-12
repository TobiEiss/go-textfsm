package lexer

import "github.com/TobiEiss/go-textfsm/pkg/models"

func (parser *Parser) parseCmd() (*models.AbstractStatement, error) {
	statement := &models.AbstractStatement{Type: models.Command, Actions: []models.Action{}}

	// function to parse a value-call
	parseValue := func(parser *Parser) (valname string, err error) {
		// need a "{""
		if token, _ := parser.scan(); token != CURLYBRACKETLEFT {
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
			case WHITESPACE:
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
		token, val := parser.scanIgnoreWhitespace()

		switch token {
		case DOLAR:
			valueName, err := parseValue(parser)
			if err != nil {
				return statement, err
			}
			statement.Actions = append(statement.Actions, models.Action{Value: valueName})
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

			// now follows the record-name
			recordname := ""
			for {
				token, val := parser.scan()
				switch token {
				case WHITESPACE:
					statement.Record = recordname
					return statement, nil
				case EOF:
					statement.Record = recordname
					return statement, nil
				default:
					recordname += val
				}
			}
		default:
			statement.Actions = append(statement.Actions, models.Action{Regex: parseRegex(val, parser)})
		}
	}
}
