package process

import (
	"regexp"

	"github.com/TobiEiss/go-textfsm/pkg/models"
)

// Process describes how to implement process an AST
type Process interface {
	Do(chan string) map[string]*Column
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
func (process process) Do(in chan string) map[string]*Column {
	// destination-record
	record := map[string]*Column{}

	// create all records and all columns
	for _, colHeader := range process.ast.Vals {
		record[colHeader.Variable] = &Column{}
	}

	// temp-record
	tmpRecord := map[string]interface{}{}

	for {
		// get next line
		line, ok := <-in
		if !ok {
			break
		}

		for _, processCommand := range process.commands {
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
						if val := process.ast.GetValForValName(names[i]); val != nil {
							if val.List {
								if tmpRecord[names[i]] != nil && tmpRecord[names[i]] != "" {
									tmpRecord[names[i]] = append(tmpRecord[names[i]].([]string), submatch[i])
								} else {
									tmpRecord[names[i]] = []string{submatch[i]}
								}
							} else {
								tmpRecord[names[i]] = submatch[i]
							}
						}

					}
				}

				// if processCommand has Record, add tempRecord to Record
				if processCommand.Command.Record == "Record" {
					// iterate all keys of record and add from tmpRecord
					for colHeader := range record {
						if val, ok := tmpRecord[colHeader]; ok {
							record[colHeader].Entries = append(record[colHeader].Entries, val)
						} else {
							record[colHeader].Entries = append(record[colHeader].Entries, "")
						}
						if val := process.ast.GetValForValName(colHeader); val != nil && !val.Filldown {
							// clear tempRecord if not filldown-field
							tmpRecord[colHeader] = ""
						}
					}
				}
			}
		}

	}
	return record
}
