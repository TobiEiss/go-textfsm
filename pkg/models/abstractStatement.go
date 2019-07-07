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
	// StateHeader is the "Header" of a new "state"
	StateHeader
)

// AbstractStatement is the raw parsed statement
type AbstractStatement struct {
	Type         StatementType
	VariableName string
	Regex        string
	Actions      []Action
	Record       bool
	Continue     bool
	Clear        bool
	Clearall     bool
	Comment      string
	Filldown     bool
	List         bool
	Required     bool
	StateName    string
	StateCall    string
}

// Action is one regex in a command
type Action struct {
	Value string
	Regex string
}

// Value creates a Val from AbstractStatement
func (statement *AbstractStatement) Value() Val {
	return Val{
		Variable: (*statement).VariableName,
		Regex:    (*statement).Regex,
		Filldown: (*statement).Filldown,
		List:     (*statement).List,
		Required: (*statement).Required,
	}
}

// Command creates a Cmd from AbstractStatement
func (statement *AbstractStatement) Command() Cmd {
	return Cmd{
		Actions:   statement.Actions,
		Record:    statement.Record,
		Continue:  statement.Continue,
		Clear:     statement.Clear,
		Clearall:  statement.Clearall,
		StateCall: statement.StateCall,
	}
}
