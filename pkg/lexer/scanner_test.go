package lexer_test

import (
	"strings"
	"testing"

	"github.com/TobiEiss/go-textfsm/pkg/lexer"
)

func TestScanner(t *testing.T) {
	// Testdata
	var tests = []struct {
		sequence string
		token    lexer.Token
		expected string
	}{
		// Special tokens (EOF, ILLEGAL, WHITESPACE)
		{sequence: ``, token: lexer.EOF},
		{`#`, lexer.HASH, `#`},
		{` `, lexer.WHITESPACE, " "},
		{"\t", lexer.WHITESPACE, "\t"},
		{"\n", lexer.WHITESPACE, "\n"},

		// misc chars
		{"(", lexer.BRACKETLEFT, "("},
		{")", lexer.BRACKETRIGHT, ")"},

		// keywords
		{"Value", lexer.VALUE, "Value"},

		// Identifiers
		{`foo`, lexer.IDENT, `foo`},
		{`abc_def-123`, lexer.IDENT, `abc_def`},
		{`23`, lexer.IDENT, `23`},

		// linguistics
		{`Ä`, lexer.IDENT, `Ä`},
		{`ö`, lexer.IDENT, `ö`},
	}

	// run test
	for i, test := range tests {
		scanner := lexer.NewScanner(strings.NewReader(test.sequence))
		token, literal := scanner.Scan()
		if test.token != token {
			t.Errorf("%d. %q token mismatch: exp=%q got=%q <%q>", i, test.sequence, test.token, token, literal)
		} else if test.expected != literal {
			t.Errorf("%d. %q literal mismatch: exp=%q got=%q", i, test.sequence, test.expected, literal)
		}
	}
}
