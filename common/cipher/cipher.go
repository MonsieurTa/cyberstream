package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

type cryptograph struct {
	key   []byte
	block cipher.Block
	gcm   cipher.AEAD
}

func NewCryptograph(aesKey string) (Cryptograph, error) {
	key, err := hex.DecodeString(aesKey)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &cryptograph{key, block, gcm}, nil
}

func (t *cryptograph) Encrypt(toEncrypt []byte) (string, error) {
	nonce := make([]byte, t.gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	return hex.EncodeToString(t.gcm.Seal(nonce, nonce, toEncrypt, nil)), nil
}

func (t *cryptograph) Decrypt(toDecrypt []byte) (string, error) {
	nonceSize := t.gcm.NonceSize()

	nonce, encrypted := toDecrypt[:nonceSize], toDecrypt[nonceSize:]
	plaintext, err := t.gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func (t *cryptograph) EncryptBatch(toEncrypt [][]byte) ([]string, error) {
	rv := make([]string, len(toEncrypt))
	for i, v := range toEncrypt {
		nonce := make([]byte, t.gcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return []string{}, err
		}
		rv[i] = hex.EncodeToString(t.gcm.Seal(nonce, nonce, v, nil))
	}
	return rv, nil
}
