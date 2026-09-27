package maintenance

import (
	"errors"
	"testing"
)

func TestNewMachine(t *testing.T) {
	if got := NewMachine().State(); got != Ready {
		t.Fatalf("initial state = %s, want %s", got, Ready)
	}
}

func TestLegalTransitions(t *testing.T) {
	tests := []struct {
		from State
		to   State
	}{
		{Ready, Draining},
		{Draining, Upgrading},
		{Draining, RetryWait},
		{Draining, Quarantined},
		{Upgrading, Verifying},
		{Upgrading, RetryWait},
		{Upgrading, Quarantined},
		{Verifying, Ready},
		{Verifying, RetryWait},
		{Verifying, Quarantined},
		{RetryWait, Draining},
		{RetryWait, Quarantined},
		{Quarantined, Ready},
	}

	for _, tt := range tests {
		t.Run(string(tt.from)+"_to_"+string(tt.to), func(t *testing.T) {
			m := &Machine{state: tt.from}
			if err := m.Transition(tt.to); err != nil {
				t.Fatalf("Transition(%s) returned error: %v", tt.to, err)
			}
			if got := m.State(); got != tt.to {
				t.Errorf("state = %s, want %s", got, tt.to)
			}
		})
	}
}

func TestIllegalTransitions(t *testing.T) {
	tests := []struct {
		from State
		to   State
	}{
		{Ready, Ready},
		{Ready, Upgrading},
		{Draining, Ready},
		{Upgrading, Draining},
		{Verifying, Upgrading},
		{RetryWait, Ready},
		{Quarantined, Draining},
		{Ready, State("UNKNOWN")},
		{State("UNKNOWN"), Ready},
	}

	for _, tt := range tests {
		t.Run(string(tt.from)+"_to_"+string(tt.to), func(t *testing.T) {
			m := &Machine{state: tt.from}
			err := m.Transition(tt.to)
			var invalid *InvalidTransitionError
			if !errors.As(err, &invalid) {
				t.Fatalf("error = %v, want *InvalidTransitionError", err)
			}
			if invalid.From != tt.from || invalid.To != tt.to {
				t.Errorf("error = %+v, want %s -> %s", invalid, tt.from, tt.to)
			}
			if got := m.State(); got != tt.from {
				t.Errorf("state = %s after rejected transition, want %s", got, tt.from)
			}
		})
	}
}
