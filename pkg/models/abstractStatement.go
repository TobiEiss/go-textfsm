package models

// StatementType describes what statement type the current is
type StatementType int

const (
	// Value describes the type of statement as variable
	Value StatementType = iota
)

// AbstractStatement is the raw parsed statement
type AbstractStatement struct {
	Type         StatementType
	VariableName string
	Regex        string
}

// Value creates a Value from AbstractStatemente
func (statement *AbstractStatement) Value() Val {
	return Val{
		Variable: (*statement).VariableName,
		Regex:    (*statement).Regex,
	}
}
