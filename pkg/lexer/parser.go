package lexer

import (
	"strings"

	"github.com/TobiEiss/go-textfsm/pkg/models"
)

// Parser represents a parser.
type Parser struct {
	scanner *Scanner
	buf     struct {
		token      Token  // last read token
		literal    string // last read literal
		buffersize int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(val string) *Parser {
	reader := strings.NewReader(val)
	return &Parser{scanner: NewScanner(reader)}
}

// ParseStatement parses a statement.
func (parser *Parser) ParseStatement() (*models.AbstractStatement, error) {
	statement := &models.AbstractStatement{}
	// Find first Token
	keywordToken, _ := parser.scanIgnoreWhitespace()
	if !isTokenAKeyWord(keywordToken) {
		return nil, parser.createError(ILLEGALTOKEN)
	}

	// Kind of keyword
	switch keywordToken {
	case VALUE:
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
	}

	return statement, nil
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (parser *Parser) scanIgnoreWhitespace() (token Token, literal string) {
	token, literal = parser.scan()
	if token == WHITESPACE {
		token, literal = parser.scan()
	}
	return
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (parser *Parser) scan() (token Token, literal string) {
	// If we have a token on the buffer, then return it.
	if parser.buf.buffersize != 0 {
		parser.buf.buffersize = 0
		return parser.buf.token, parser.buf.literal
	}

	// Otherwise read the next token from the scanner.
	token, literal = parser.scanner.Scan()

	// Save it to the buffer in case we unscan later.
	parser.buf.token, parser.buf.literal = token, literal
	return
}

func isTokenAKeyWord(token Token) bool {
	for _, value := range KeyWordMap {
		if value == token {
			return true
		}
	}
	return false
}

func (parser *Parser) createError(errorType ErrorType) *Error {
	return &Error{ErrorType: errorType}
}
