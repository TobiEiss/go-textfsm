package lexer_test

import (
	"testing"

	"github.com/TobiEiss/go-textfsm/pkg/lexer"
	"github.com/TobiEiss/go-textfsm/pkg/models"
)

func TestParserValue(t *testing.T) {
	var tests = []struct {
		ValStr      string
		ErrorType   lexer.ErrorType
		ExpectedVal models.Val
	}{
		// legal cases
		{ValStr: `Value Year (\d+)`, ExpectedVal: models.Val{Variable: "Year", Regex: `\d+`}},
		{ValStr: `Value MonthDay (\d+)`, ExpectedVal: models.Val{Variable: "MonthDay", Regex: `\d+`}},
		{ValStr: `Value Month (\w+)`, ExpectedVal: models.Val{Variable: "Month", Regex: `\w+`}},
		{ValStr: `Value Timezone (\S+)`, ExpectedVal: models.Val{Variable: "Timezone", Regex: `\S+`}},
		{ValStr: `Value Time (..:..:..)`, ExpectedVal: models.Val{Variable: "Time", Regex: `..:..:..`}},

		// illegal cases
		{ValStr: "Valuee ", ErrorType: lexer.ILLEGALTOKEN},
		{ValStr: "Value", ErrorType: lexer.MISSINGARGUMENT},
		{ValStr: "abc", ErrorType: lexer.ILLEGALTOKEN},
		{ValStr: "Value Value", ErrorType: lexer.ILLEGALTOKEN},
		{ValStr: `Value Year`, ErrorType: lexer.MISSINGARGUMENT},
		{ValStr: `Value (`, ErrorType: lexer.ILLEGALTOKEN},
	}

	// iterate all tests
	for index, test := range tests {
		stmt, err := lexer.NewParser(test.ValStr).ParseStatement()

		// expected Error-Case
		if lexerError, ok := err.(*lexer.Error); ok {
			if lexerError.ErrorType != test.ErrorType {
				t.Errorf("%d failed: Not expected error: %s", index, err)
			}
			// error was expected
		}

		// no-Error-Case - check expected val
		if err == nil {
			val := stmt.Value()
			if val.Variable != test.ExpectedVal.Variable {
				t.Errorf("%d failed: Variable '%s' is not equal expected Variable '%s'",
					index, val.Variable, test.ExpectedVal.Variable)
			}

			if stmt.Regex != test.ExpectedVal.Regex {
				t.Errorf("%d failed: Regex '%s' is not equal expected Regex '%s'",
					index, stmt.Regex, test.ExpectedVal.Regex)
			}
		}
	}
}
