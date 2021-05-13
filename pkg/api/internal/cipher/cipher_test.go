package cipher

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslator(t *testing.T) {
	aesKey := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, aesKey)
	assert.Nil(t, err)

	translator, err := NewTranslator(hex.EncodeToString(aesKey))
	assert.Nil(t, err)

	sizes := []int{128, 256, 512, 1024}
	randomDatas := make([][]byte, len(sizes))
	for i, size := range sizes {
		randomDatas[i] = make([]byte, size)
		_, err := io.ReadFull(rand.Reader, randomDatas[i])
		assert.Nil(t, err)
	}

	for _, bytes := range randomDatas {
		data, err := translator.Encrypt(bytes)
		assert.Nil(t, err)

		encrypted, err := hex.DecodeString(data)
		assert.Nil(t, err)

		decrypted, err := translator.Decrypt(encrypted)
		assert.Nil(t, err)

		expected := hex.EncodeToString(bytes)
		assert.Equal(t, expected, decrypted)
	}

	encryptedDatas, err := translator.EncryptBatch(randomDatas)
	assert.Nil(t, err)
	for i, encrypted := range encryptedDatas {
		data, err := hex.DecodeString(encrypted)
		assert.Nil(t, err)

		decrypted, err := translator.Decrypt(data)
		assert.Nil(t, err)

		expected := hex.EncodeToString(randomDatas[i])
		assert.Equal(t, expected, decrypted)
	}
}
