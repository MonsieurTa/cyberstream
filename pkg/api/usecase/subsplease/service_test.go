package subsplease

import (
	"os"
	"testing"

	repo "github.com/MonsieurTa/hypertube/pkg/api/internal/subsplease"
	"github.com/stretchr/testify/assert"
)

func TestSubsPlease(t *testing.T) {
	os.Setenv("AES_KEY", "4737EF03C50EEC6C0651757F44B38DF1A66178C30DE9578A3F290C87EBBE8CE4")

	repo := repo.NewSubsPlease()
	assert.NotNil(t, repo)

	rv, err := repo.Latest()
	assert.Nil(t, err)

	expectedSize := len(rv)

	service := NewService(repo)
	movies, err := service.Latest()
	assert.Nil(t, err)
	assert.Equal(t, expectedSize, len(movies))
}
