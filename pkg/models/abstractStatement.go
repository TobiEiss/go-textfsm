package models

// StatementType describes what statement type the current is
type StatementType int

const (
	// Value describes the type of statement as variable
	Value StatementType = iota + 1
	// Start is the start statement
	Start
	// Command represent a command
	Command
	// Comment is a statement to ignore
	Comment
)

// AbstractStatement is the raw parsed statement
type AbstractStatement struct {
	Type         StatementType
	VariableName string
	Regex        string
	Actions      []Action
	Record       string
	Comment      string
	Filldown     bool
	List         bool
}

// Action is one regex in a command
type Action struct {
	Value string
	Regex string
}

// Value creates a Val from AbstractStatemente
func (statement *AbstractStatement) Value() Val {
	return Val{
		Variable: (*statement).VariableName,
		Regex:    (*statement).Regex,
		Filldown: (*statement).Filldown,
		List:     (*statement).List,
	}
}

// Command creates a Cmd from AbstractStatemente
func (statement *AbstractStatement) Command() Cmd {
	return Cmd{
		Actions: statement.Actions,
		Record:  statement.Record,
	}
}
