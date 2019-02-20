package process

import (
	"errors"
	"log"
	"regexp"

	"github.com/TobiEiss/go-textfsm/pkg/models"
)

// Process describes how to implement process an AST
type Process interface {
	Do(chan string) bool
}

type process struct {
	ast      models.AST
	commands []processCommand
}

type processCommand struct {
	Command     models.Cmd
	LineCommand string
}

// NewProcess create a new implementation of Process
func NewProcess(ast models.AST) (Process, error) {
	process := &process{commands: []processCommand{}, ast: ast}

	// calculate lineCommand
	for _, command := range ast.Commands {
		lineCommand := ""
		for _, action := range command.Actions {
			if action.Value != "" {
				if val := ast.GetValForValName(action.Value); val != nil {
					lineCommand += val.Regex
					continue
				}
				return nil, errors.New("Can't find val for ValName")
			}
			lineCommand += action.Regex + action.Value
		}
		process.commands = append(process.commands, processCommand{Command: command, LineCommand: lineCommand})
	}

	return process, nil
}

// Do process an ast. Get inputfile as channel line by line
func (process process) Do(in chan string) bool {
	for {
		// get next line
		line, ok := <-in
		if !ok {
			break
		}

		// check if one command matches to line
		for _, processCommand := range process.commands {
			re := regexp.MustCompile(processCommand.LineCommand)
			log.Println(processCommand.Command.Actions)
			return re.MatchString(line)
		}
	}
	return false
}
