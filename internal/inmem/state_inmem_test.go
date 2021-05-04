package inmem

import "testing"

func TestNewStateInMem(t *testing.T) {
	v := NewStateInMem()
	if v == nil {
		t.Errorf("NewStateInMem returned %v\n", v)
	}
}

func TestStateInMem(t *testing.T) {
	tests := [5]string{}
	for i := 0; i < 5; i++ {
		state, err := GenerateState()
		if err != nil {
			t.Errorf("state generation failed: got %s\n", state)
		}
		tests[i] = state
	}

	stateInMem := NewStateInMem()
	for _, v := range tests {
		stateInMem.Save(v)
	}

	for _, v := range tests {
		err := stateInMem.Exist(v)
		if err != nil {
			t.Errorf("could not retrieve %s in store\n", v)
		}
	}

	for _, v := range tests {
		stateInMem.Delete(v)
		err := stateInMem.Exist(v)
		if err == nil {
			t.Errorf("did not delete entry in store: %s\n", v)
		}
	}
}
