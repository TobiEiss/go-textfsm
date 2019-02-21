package process

import (
	"regexp"

	"github.com/TobiEiss/go-textfsm/pkg/models"
)

// Process describes how to implement process an AST
type Process interface {
	Do(chan string) map[string]map[string][]interface{}
}

type process struct {
	ast      models.AST
	commands []processCommand
}

type commandPart struct {
	regex     string
	valueName string
}

type processCommand struct {
	MatchingLine string
	Command      models.Cmd
}

// NewProcess create a new implementation of Process
func NewProcess(ast models.AST) (Process, error) {
	process := &process{commands: []processCommand{}, ast: ast}

	// calculate matchingLine
	for _, command := range ast.Commands {
		matchingLine, err := ast.CreateMatchingLine(command)
		if err != nil {
			return nil, err
		}
		process.commands = append(process.commands,
			processCommand{MatchingLine: matchingLine, Command: command})
	}

	return process, nil
}

// Do process an ast. Get inputfile as channel line by line
func (process process) Do(in chan string) map[string]map[string][]interface{} {
	// destination-record
	record := map[string]map[string][]interface{}{}

	// temp-record
	tmpRecord := map[string][]interface{}{}

	commandIndex := 0
	for {
		// get next line
		line, ok := <-in
		if !ok {
			break
		}

		processCommand := process.commands[commandIndex%len(process.commands)]

		// check one command matches to line
		re := regexp.MustCompile(processCommand.MatchingLine)

		// check if line is relevant
		if re.MatchString(line) {
			submatch := re.FindStringSubmatch(line)
			names := re.SubexpNames()

			// len of submatch and names should be same
			if len(submatch) == len(names) {

				// add all founded fields to record
				for i := 1; i < len(names); i++ {
					tmpRecord[names[i]] = append(tmpRecord[names[i]], submatch[i])
				}
			}

			// if processCommand has Record, add tempRecord to Record
			if processCommand.Command.Record != "" {
				if record[processCommand.Command.Record] == nil {
					record[processCommand.Command.Record] = map[string][]interface{}{}
				}
				// iterate all keys of tmpRecord and append to destination-record
				for k, v := range tmpRecord {
					record[processCommand.Command.Record][k] =
						append(record[processCommand.Command.Record][k], v...)
				}
				// clear tempRecord
				tmpRecord = map[string][]interface{}{}
			}

			// increase the commandIndex to get next round the next command
			commandIndex++
		}

	}
	return record
}
