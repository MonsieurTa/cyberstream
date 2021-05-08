package inmem

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
)

type StateInMem struct {
	states map[string]bool
}

const (
	DEFAULT_STATE_SIZE = 64
)

func NewStateInMem() *StateInMem {
	return &StateInMem{
		states: make(map[string]bool),
	}
}

func (s *StateInMem) Save(state string) {
	s.states[state] = true
}

func (s *StateInMem) Delete(state string) {
	delete(s.states, state)
}

func (s *StateInMem) Exist(state string) error {
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
