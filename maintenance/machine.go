// Package maintenance models the lifecycle of a node maintenance attempt.
package maintenance

import "fmt"

// State is a phase of node maintenance.
type State string

const (
	Ready       State = "READY"
	Draining    State = "DRAINING"
	Upgrading   State = "UPGRADING"
	Verifying   State = "VERIFYING"
	RetryWait   State = "RETRY_WAIT"
	Quarantined State = "QUARANTINED"
)

// allowedTransitions is the complete set of legal state changes.
var allowedTransitions = map[State]map[State]struct{}{
	Ready:       {Draining: {}},
	Draining:    {Upgrading: {}, RetryWait: {}, Quarantined: {}},
	Upgrading:   {Verifying: {}, RetryWait: {}, Quarantined: {}},
	Verifying:   {Ready: {}, RetryWait: {}, Quarantined: {}},
	RetryWait:   {Draining: {}, Quarantined: {}},
	Quarantined: {Ready: {}},
}

// InvalidTransitionError identifies a rejected state change.
type InvalidTransitionError struct {
	From State
	To   State
}

func (e *InvalidTransitionError) Error() string {
	return fmt.Sprintf("invalid maintenance transition: %s -> %s", e.From, e.To)
}

// Machine holds the current maintenance state. A new machine starts ready.
type Machine struct {
	state State
}

func NewMachine() *Machine {
	return &Machine{state: Ready}
}

func (m *Machine) State() State {
	return m.state
}

// Transition changes state only when the requested edge is legal.
func (m *Machine) Transition(to State) error {
	if _, ok := allowedTransitions[m.state][to]; !ok {
		return &InvalidTransitionError{From: m.state, To: to}
	}
	m.state = to
	return nil
}
