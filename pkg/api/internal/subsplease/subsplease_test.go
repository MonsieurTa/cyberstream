package subsplease

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubsPlease(t *testing.T) {
	repo := NewSubsPlease()
	assert.NotNil(t, repo)

	_, err := repo.Latest()
	assert.Nil(t, err)
}
