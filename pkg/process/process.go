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
	ast    models.AST
	states []state
}

type commandPart struct {
	regex     string
	valueName string
}

type processCommand struct {
	MatchingLine string
	Command      models.Cmd
}

type state struct {
	processCommands []processCommand
	state           models.State
}

// NewProcess create a new implementation of Process
func NewProcess(ast models.AST) (Process, error) {
	process := &process{states: []state{}, ast: ast}

	// calculate matchingLine
	for _, st := range ast.States {
		currentState := state{state: st}
		for _, command := range st.Commands {
			matchingLine, err := ast.CreateMatchingLine(command)
			if err != nil {
				return nil, err
			}
			currentState.processCommands = append(currentState.processCommands,
				processCommand{MatchingLine: matchingLine, Command: command})
		}
		process.states = append(process.states, currentState)
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

	// machineState is one state in the statemachine
	type machineState struct {
		tmpRecord map[string]interface{}
		state     state
	}
	// statemachine is <nameOfState><machineState>. All machineStates which are in the statemachine are active
	statemachine := map[string]machineState{}

	// add "start"-states to statemachine
	findStateAndAddToMachine := func(stateName string, tmpRecord map[string]interface{}) {
		for _, state := range process.states {
			if state.state.Name == stateName {
				statemachine[state.state.Name] = machineState{tmpRecord: tmpRecord, state: state}
				break
			}
		}
	}
	findStateAndAddToMachine("Start", map[string]interface{}{})

	// keep last added record in mind
	lastAddedRecord := map[string]interface{}{}

	// iterate lines
	for {
		// get next line
		line, ok := <-in
		if !ok {
			break
		}

		// iterate all activestates
		for stateName, activeState := range statemachine {
			// iterate commands of a active state
			for _, processCommand := range activeState.state.processCommands {
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
									if activeState.tmpRecord[names[i]] != nil && activeState.tmpRecord[names[i]] != "" {
										activeState.tmpRecord[names[i]] = append(activeState.tmpRecord[names[i]].([]string), submatch[i])
									} else {
										activeState.tmpRecord[names[i]] = []string{submatch[i]}
									}
								} else {
									activeState.tmpRecord[names[i]] = submatch[i]
								}
							}

						}
					}

					// if processCommand has Record, add tempRecord to Record
					if processCommand.Command.Record {
						// iterate all keys of record and add from tmpRecord
						for colHeader := range record {
							if val, ok := activeState.tmpRecord[colHeader]; ok {
								record[colHeader].Entries = append(record[colHeader].Entries, val)
							} else {
								record[colHeader].Entries = append(record[colHeader].Entries, "")
							}
							if val := process.ast.GetValForValName(colHeader); val != nil && !val.Filldown {
								// clear tempRecord if not filldown-field
								activeState.tmpRecord[colHeader] = ""
								// this is lastaAddedRecord for next state (respect Filldown)
								lastAddedRecord = activeState.tmpRecord
							}
						}
					}

					// if state ends
					if processCommand.Command.StateCall == "Start" {
						delete(statemachine, stateName)
					}

					// if state calls a new state
					if processCommand.Command.StateCall != "" {
						findStateAndAddToMachine(processCommand.Command.StateCall, lastAddedRecord)
					}
				}
			}
		}

	}
	checkRequired(record, process)
	return record
}

// check if all required Fields have Values,
// if not remove the the record from records.
func checkRequired(record map[string]*Column, proc process) {
	var removeIndices []int
	for name, column := range record {
		val := proc.ast.GetValForValName(name)
		if val.Required {
			for idx, value := range column.Entries {
				if value == "" {
					if removeIndices == nil {
						removeIndices = []int{idx}
					} else {
						removeIndices = append(removeIndices, idx)
					}
				}
			}
		}
	}
	for _, idx := range removeIndices {
		for _, column := range record {
			column.Entries = append(column.Entries[:idx], column.Entries[idx+1:]...)

		}
	}

}
