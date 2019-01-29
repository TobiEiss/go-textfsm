package lexer

// ErrorType represent the type of the error
type ErrorType string

const (
	ILLEGALTOKEN    ErrorType = "illegal token"
	MISSINGARGUMENT ErrorType = "missing argument"
)

// Error is a error-type only for
type Error struct {
	error
	ErrorType ErrorType
}

func (e *Error) Error() string {
	return string(e.ErrorType)
}
