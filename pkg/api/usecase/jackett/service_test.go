package jackett

import (
	"testing"

	"github.com/MonsieurTa/hypertube/common/infrastructure/repository"
	"github.com/stretchr/testify/assert"
)

func TestConfiguredIndexers(t *testing.T) {
	repo := repository.NewJackett()

	s := NewService(repo)
	_, err := s.ConfiguredIndexers()
	assert.Nil(t, err)
}
