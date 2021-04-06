package fortytwo

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
)

type StateManager struct {
	states map[string]bool
}

func NewStateManager() *StateManager {
	return &StateManager{
		states: make(map[string]bool),
	}
}

func (s *StateManager) SaveStateInMemory(state string) {
	s.states[state] = true
}

func (s *StateManager) DeleteStateInMemory(state string) {
	delete(s.states, state)
}

func (s *StateManager) ValidateState(state string) error {
	if _, ok := s.states[state]; !ok {
		return errors.New("state not recognized")
	}
	return nil
}

func generateState(n int) (string, error) {
	data := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}
	h := sha256.New()
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
