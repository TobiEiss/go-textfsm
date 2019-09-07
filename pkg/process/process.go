package process

import (
	"fmt"
	"reflect"
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
		lastAddedRow:      make([]interface{}, len(ast.Vals)),
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
	// first active State is always Start,
	machineState := process.machine["Start"]
	stateName := "Start"
	activeState := machineState

	tmp := make([]interface{}, len(process.ast.Vals))
	// iterate lines
	for {
		// get next line
		line, ok := <-in
		if !ok {
			break
		}

	Start:
		// iterate commands of a current activeState commands
		for _, processCommand := range activeState.StateDescription.ProcessCommands {
			// check one command matches to line
			re := regexp.MustCompile(processCommand.MatchingLine)

			// check if line is relevant
			if re.MatchString(line) {

				processLine(line, re, process, stateName, machineState, processCommand)

				// change activeState and back iterating
				if processCommand.Command.StateCall != "" {
					stateName = processCommand.Command.StateCall
					activeState = process.getStateMachine(stateName, tmp)
					goto Start
				}

				if !processCommand.Command.Continue {
					break
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

// get the state-machine, if not exist create one and return it.
func (process *process) getStateMachine(stateName string, tmpRow []interface{}) *statemachine.MachineState {

Start:
	for name, machine := range process.machine {
		if name == stateName {
			return machine
		}
	}
	process.findStateAndAddToMachine(stateName, tmpRow)
	goto Start

}

// process a single line and add records
func processLine(line string, re *regexp.Regexp, process *process, stateName string,
	machineState *statemachine.MachineState, processCommand statemachine.ProcessCommand) {

	submatch := re.FindStringSubmatch(line)
	names := re.SubexpNames()

	// len of submatch and names should be same
	if len(submatch) == len(names) {
		// transform result to map
		result := map[string]interface{}{}
		for index, name := range names {
			if existing, ok := result[name]; ok {
				// is result[name] already a slice?
				if _, ok := result[name].([]string); ok {
					result[name] = append(result[name].([]string), submatch[index])
				} else {
					result[name] = []string{fmt.Sprintf("%v", existing), submatch[index]}
				}
			} else {
				result[name] = submatch[index]
			}
		}

		// add all founded fields to record
		for index, val := range process.ast.Vals {
			if field, ok := result[val.Variable]; ok {
				if val.List && reflect.TypeOf(field).Kind() != reflect.Slice {
					field = []string{fmt.Sprintf("%v", field)}
				}
				machineState.SetRowField(index, field)
			}
		}
	}

	if processCommand.Command.Clearall {
		for index := range process.ast.Vals {
			machineState.SetRowField(index, "")
		}
	}

	if processCommand.Command.Clear {
		for index, val := range process.ast.Vals {
			if !val.Filldown {
				machineState.SetRowField(index, "")
			}
		}
	}

	// if processCommand has Record, add tempRecord to Record
	if processCommand.Command.Record {
		requiredFieldIsEmpty := false

		// iterate all vals
		for index, val := range process.ast.Vals {
			// removeFlag if required-Field is nil
			if machineState.TmpRowField(index) == nil && val.Required {
				requiredFieldIsEmpty = true
				continue
			}
			// add an empty string if tmpRow-Item is nil and val is FILLDOWN
			if machineState.TmpRowField(index) == nil && val.Filldown {
				machineState.SetRowField(index, process.lastAddedRow[index])
			} else if machineState.TmpRowField(index) == nil && !val.Filldown {
				machineState.SetRowField(index, "")
			}
		}

		if !requiredFieldIsEmpty {
			tmp := make([]interface{}, len(machineState.TmpRow()))
			copy(tmp, machineState.TmpRow())
			process.lastAddedRow = tmp
			process.record <- tmp

			machineState.ClearRowField()
		}
	}

}
