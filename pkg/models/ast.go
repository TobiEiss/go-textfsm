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

// GetValForValName searches a val for a valName
func (ast AST) GetValForValName(valName string) *Val {
	for _, val := range ast.Vals {
		if val.Variable == valName {
			return &val
		}
	}
	return nil
}
