package statemachine

import (
	"sync"

	"github.com/TobiEiss/go-textfsm/pkg/models"
)

// Machine is the machine
type Machine struct {
}

// MachineState is a state for a machine
type MachineState struct {
	tmpRow           []interface{}
	StateDescription StateDescription
	mux              *sync.Mutex
}

// NewMachineState creates a new state
func NewMachineState(tmpRow []interface{}, stateDescription StateDescription) *MachineState {
	return &MachineState{tmpRow: tmpRow, StateDescription: stateDescription, mux: &sync.Mutex{}}
}

// TmpRow returns tmpRow
func (machineState *MachineState) TmpRow() []interface{} {
	machineState.mux.Lock()
	defer machineState.mux.Unlock()
	return machineState.tmpRow
}

// TmpRowField returns an item of tmpRow
func (machineState *MachineState) TmpRowField(index int) interface{} {
	machineState.mux.Lock()
	defer machineState.mux.Unlock()
	return machineState.tmpRow[index]
}

// SetRowField set a specific field to row
func (machineState *MachineState) SetRowField(index int, field interface{}) {
	machineState.mux.Lock()
	machineState.tmpRow[index] = field
	machineState.mux.Unlock()
}

// AppendToRowField appends to specific field
func (machineState *MachineState) AppendToRowField(index int, field interface{}) {
	machineState.mux.Lock()
	if machineState.tmpRow[index] != nil && machineState.tmpRow[index] != "" {
		machineState.tmpRow[index] = append(machineState.tmpRow[index].([]interface{}), field)
	} else {
		machineState.tmpRow[index] = []interface{}{field}
	}
	machineState.mux.Unlock()
}

// StateDescription holds the origin state
type StateDescription struct {
	ProcessCommands []ProcessCommand
	OriginState     models.State
}

// ProcessCommand is command for the process
type ProcessCommand struct {
	MatchingLine string
	Command      models.Cmd
}
