package models

// Val represent a varible like "Value Year (\d+)"
type Val struct {
	Variable string
	Regex    string
}

// Statement is one statement after keyword "Start"
type Statement struct {
	Action string
}

// AST is the abstract statement tree
type AST struct {
	Vals       []Val
	Statements []Statement
}
