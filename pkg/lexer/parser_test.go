package lexer_test

import (
	"testing"

	"github.com/TobiEiss/go-textfsm/pkg/lexer"
)

func TestParserValue(t *testing.T) {
	var tests = []struct {
		StatementStr      string
		ErrorType         lexer.ErrorType
		ExpectedStatement lexer.Statement
	}{
		// legal cases
		{StatementStr: `Value Year (\d+)`, ExpectedStatement: lexer.Statement{Keyword: lexer.VALUE, Variable: "Year", Regex: `\d+`}},
		{StatementStr: `Value MonthDay (\d+)`, ExpectedStatement: lexer.Statement{Keyword: lexer.VALUE, Variable: "MonthDay", Regex: `\d+`}},
		{StatementStr: `Value Month (\w+)`, ExpectedStatement: lexer.Statement{Keyword: lexer.VALUE, Variable: "Month", Regex: `\w+`}},
		{StatementStr: `Value Timezone (\S+)`, ExpectedStatement: lexer.Statement{Keyword: lexer.VALUE, Variable: "Timezone", Regex: `\S+`}},
		{StatementStr: `Value Time (..:..:..)`, ExpectedStatement: lexer.Statement{Keyword: lexer.VALUE, Variable: "Time", Regex: `..:..:..`}},

		// illegal cases
		{StatementStr: "Valuee ", ErrorType: lexer.ILLEGALTOKEN},
		{StatementStr: "Value", ErrorType: lexer.MISSINGARGUMENT},
		{StatementStr: "abc", ErrorType: lexer.ILLEGALTOKEN},
		{StatementStr: "Value Value", ErrorType: lexer.ILLEGALTOKEN},
		{StatementStr: `Value Year`, ErrorType: lexer.MISSINGARGUMENT},
		{StatementStr: `Value (`, ErrorType: lexer.ILLEGALTOKEN},
	}

	// iterate all tests
	for index, test := range tests {
		stmt, err := lexer.NewParser(test.StatementStr).ParseStatement()

		// expected Error-Case
		if lexerError, ok := err.(*lexer.Error); ok {
			if lexerError.ErrorType != test.ErrorType {
				t.Errorf("%d failed: Not expected error: %s", index, err)
			}
			// error was expected
		}

		// no-Error-Case - check expected statement
		if err == nil {
			if stmt.Variable != test.ExpectedStatement.Variable {
				t.Errorf("%d failed: Variable '%s' is not equal expected Variable '%s'",
					index, stmt.Variable, test.ExpectedStatement.Variable)
			}

			if stmt.Keyword != test.ExpectedStatement.Keyword {
				t.Errorf("%d failed: Keyword '%d' is not equal expected Keyword '%d'",
					index, stmt.Keyword, test.ExpectedStatement.Keyword)
			}

			if stmt.Regex != test.ExpectedStatement.Regex {
				t.Errorf("%d failed: Regex '%s' is not equal expected Regex '%s'",
					index, stmt.Regex, test.ExpectedStatement.Regex)
			}
		}
	}
}
