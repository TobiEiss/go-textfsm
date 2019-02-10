package models

// Val represent a varible like "Value Year (\d+)"
type Val struct {
	Variable string
	Regex    string
}

// Command is one statement after keyword "Start"
type Command struct {
	Action string
}

// AST is the abstract command tree
type AST struct {
	Vals    []Val
	Command []Command
}
