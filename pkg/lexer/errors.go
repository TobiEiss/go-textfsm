package lexer

import "fmt"

// ErrorType represent the type of the error
type ErrorType string

const (
	ILLEGALTOKEN    ErrorType = "illegal token"
	MISSINGARGUMENT ErrorType = "missing argument"

	colorBlue   = "\033[1;34m%s\033[0m"
	colorAqua   = "\033[1;36m%s\033[0m"
	colorYellow = "\033[1;33m%s\033[0m"
	colorRed    = "\033[1;31m%s\033[0m"
)

// Error is a error-type only for
type Error struct {
	error
	ErrorType   ErrorType
	CurrentLine string
	ErrorToken  string
}

func (e *Error) Error() string {
	err := red(string(e.ErrorType))
	if e.CurrentLine != "" {
		err = fmt.Sprintf("%s in line: '%s'", err, yellow(e.CurrentLine))
	}
	if e.ErrorToken != "" {
		err = fmt.Sprintf("%s. Problem is token: '%s'", err, yellow(e.ErrorToken))
	}
	return err
}

func red(msg string) string {
	return fmt.Sprintf(colorRed, msg)
}

func blue(msg string) string {
	return fmt.Sprintf(colorBlue, msg)
}

func aqua(msg string) string {
	return fmt.Sprintf(colorAqua, msg)
}

func yellow(msg string) string {
	return fmt.Sprintf(colorYellow, msg)
}
