package repository

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubsPlease(t *testing.T) {
	repo := NewSubsPlease()
	assert.NotNil(t, repo)

	rv, err := repo.Latests()
	assert.Nil(t, err)
	fmt.Printf("%v\n", rv)
}
