package models

// AST is the abstract command tree
type AST struct {
	Vals     []Val
	Commands []Cmd
}

// Val represent a varible like "Value Year (\d+)"
type Val struct {
	Variable string
	Regex    string
}

// Cmd is one statement after keyword "Start"
type Cmd struct {
	Actions []Action
	Record  string
}
