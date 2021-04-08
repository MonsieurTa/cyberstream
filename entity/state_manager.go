package entity

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

const (
	DEFAULT_STATE_SIZE = 64
)

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

func GenerateState() (string, error) {
	data := make([]byte, DEFAULT_STATE_SIZE)
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
