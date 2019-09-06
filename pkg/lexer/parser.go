package lexer

import (
	"fmt"
	"strings"

	"github.com/TobiEiss/go-textfsm/pkg/models"
)

// Parser represents a parser.
type Parser struct {
	currentline string // only for error-logging
	scanner     *Scanner
	buf         struct {
		token      Token  // last read token
		literal    string // last read literal
		buffersize int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(val string) *Parser {
	reader := strings.NewReader(val)
	return &Parser{scanner: NewScanner(reader), currentline: val}
}

// ParseStatement parses a statement.
func (parser *Parser) ParseStatement() (*models.AbstractStatement, error) {
	// Find first Token
	token, val := parser.scanIgnoreWhitespace()

	// check if this a keyword
	if isTokenAKeyWord(token) {
		// Kind of keyword
		switch token {
		case VALUE:
			return parser.parseVal()
		case START:
			return &models.AbstractStatement{Type: models.StateHeader, StateName: "Start"}, nil
		}
		return nil, &Error{ErrorType: ILLEGALTOKEN, CurrentLine: parser.currentline, ErrorToken: val}
	}

	// if token is not a keyword it have to be a command
	if token == CIRCUMFLEX {
		return parser.parseCmd()
	}

	// check if this is a comment
	if token == HASH {
		return parser.parseCmmt()
	}

	// should be an state-header
	if token == IDENT {
		return parser.parseStateHeader(val)
	}

	// if this is a nil-line -> continue
	if token == EOF {
		return nil, nil
	}

	return nil, fmt.Errorf("Can't parse line: %s", parser.currentline)
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (parser *Parser) scanIgnoreWhitespace() (token Token, literal string) {
	token, literal = parser.scan()
	for token == WHITESPACE {
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

// unscan pushes the previously read token back onto the buffer.
func (parser *Parser) unscan() {
	parser.buf.buffersize = 1
}

func isTokenAKeyWord(token Token) bool {
	for _, value := range KeyWordMap {
		if value == token {
			return true
		}
	}
	return false
}
