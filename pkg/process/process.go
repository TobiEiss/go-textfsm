package process

import (
	"regexp"

	"github.com/TobiEiss/go-textfsm/pkg/models"
	"github.com/TobiEiss/go-textfsm/pkg/process/statemachine"
)

// Process describes how to implement process an AST
type Process interface {
	Do(in chan string)
}

// machine is <nameOfState><machineState>. All machineStates which are in the statemachine are active
type process struct {
	ast               models.AST
	stateDescriptions []statemachine.StateDescription
	machine           map[string]*statemachine.MachineState
	lastAddedRow      []interface{}
	record            chan<- []interface{}
}

type commandPart struct {
	regex     string
	valueName string
}

// NewProcess create a new implementation of Process
func NewProcess(ast models.AST, out chan<- []interface{}) (Process, error) {
	process := &process{
		stateDescriptions: []statemachine.StateDescription{},
		ast:               ast,
		machine:           map[string]*statemachine.MachineState{},
		lastAddedRow:      []interface{}{},
		record:            out,
	}

	// calculate matchingLine
	for _, st := range ast.States {
		currentState := statemachine.StateDescription{OriginState: st}
		for _, command := range st.Commands {
			matchingLine, err := ast.CreateMatchingLine(command)
			if err != nil {
				return nil, err
			}
			currentState.ProcessCommands = append(currentState.ProcessCommands,
				statemachine.ProcessCommand{MatchingLine: matchingLine, Command: command})
		}
		process.stateDescriptions = append(process.stateDescriptions, currentState)
	}

	// add "start"-states to statemachine
	process.findStateAndAddToMachine("Start", make([]interface{}, len(process.ast.Vals)))

	return process, nil
}

// Do process an ast. Get inputfile as channel line by line
func (process *process) Do(in chan string) {
	// iterate lines
	for {
		// get next line
		line, ok := <-in
		if !ok {
			break
		}
		// iterate all activestates
		for stateName, activeState := range process.machine {
			// iterate commands of a active state
			for _, processCommand := range activeState.StateDescription.ProcessCommands {
				// check one command matches to line
				re := regexp.MustCompile(processCommand.MatchingLine)

				// check if line is relevant
				if re.MatchString(line) {

					submatch := re.FindStringSubmatch(line)
					names := re.SubexpNames()

					// len of submatch and names should be same
					if len(submatch) == len(names) {
						// transform result to map
						result := map[string]interface{}{}
						for index, name := range names {
							result[name] = submatch[index]
						}

						// add all founded fields to record
						for index, val := range process.ast.Vals {
							if field, ok := result[val.Variable]; ok {
								if val.List {
									activeState.AppendToRowField(index, field)
								} else {
									activeState.SetRowField(index, field)
								}
							}
						}
					}

					// if processCommand has Record, add tempRecord to Record
					if processCommand.Command.Record {
						requiredFieldIsEmpty := false

						// iterate all vals
						for index, val := range process.ast.Vals {
							// removeFlag if required-Field is nil
							if activeState.TmpRowField(index) == nil && val.Required {
								requiredFieldIsEmpty = true
							}
							// add an empty string if tmpRow-Item is nil anv val is not FILLDOWN
							if activeState.TmpRowField(index) == nil && !val.Filldown {
								activeState.SetRowField(index, nil)
							}
						}

						if !requiredFieldIsEmpty {
							tmp := make([]interface{}, len(activeState.TmpRow()))
							copy(tmp, activeState.TmpRow())
							process.record <- tmp
							process.lastAddedRow = tmp
						}
					}

					// if state ends
					if processCommand.Command.StateCall == "Start" {
						delete(process.machine, stateName)
					}

					// if state calls a new state
					if processCommand.Command.StateCall != "" {
						process.findStateAndAddToMachine(processCommand.Command.StateCall, process.lastAddedRow)
					}
				}
			}
		}
	}
	close(process.record)
}

func (process *process) findStateAndAddToMachine(stateName string, tmpRow []interface{}) {
	for _, stateDescription := range process.stateDescriptions {
		if stateDescription.OriginState.Name == stateName {
			process.machine[stateDescription.OriginState.Name] = statemachine.NewMachineState(tmpRow, stateDescription)
			break
		}
	}
}
