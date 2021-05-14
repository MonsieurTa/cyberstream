package subsplease

import (
	"encoding/hex"
	"os"
	"strings"
	"testing"

	"github.com/MonsieurTa/hypertube/common/cipher"
	repo "github.com/MonsieurTa/hypertube/pkg/api/internal/subsplease"
	"github.com/stretchr/testify/assert"
)

func TestSubsPlease(t *testing.T) {
	os.Setenv("AES_KEY", "4737ef03c50eec6c0651757f44b38df1a66178c30de9578a3f290c87ebbe8ce4")

	repo := repo.NewSubsPlease()
	assert.NotNil(t, repo)

	rv, err := repo.Latest()
	assert.Nil(t, err)

	expectedSize := len(rv)

	service := NewService(repo)
	movies, err := service.Latest()
	assert.Nil(t, err)
	assert.Equal(t, expectedSize, len(movies))

	c, err := cipher.NewCryptograph(os.Getenv("AES_KEY"))
	assert.Nil(t, err)
	for _, m := range movies {
		encrypted, err := hex.DecodeString(m.Magnet)
		assert.Nil(t, err)

		decrypted, err := c.Decrypt(encrypted)
		assert.Nil(t, err)
		if !strings.HasPrefix(decrypted, "magnet:?") {
			t.Errorf("expected magnet:? got %s\n", decrypted[0:4])
		}
	}
}
