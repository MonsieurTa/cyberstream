package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

type translator struct {
	key   []byte
	block cipher.Block
	gcm   cipher.AEAD
}

func NewTranslator(aesKey string) (Translator, error) {
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
	return &translator{key, block, gcm}, nil
}

func (t *translator) Encrypt(toEncrypt []byte) (string, error) {
	nonce := make([]byte, t.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	return hex.EncodeToString(t.gcm.Seal(nonce, nonce, toEncrypt, nil)), nil
}

func (t *translator) Decrypt(toDecrypt []byte) (string, error) {
	nonceSize := t.gcm.NonceSize()
	nonce, encrypted := toDecrypt[:nonceSize], toDecrypt[nonceSize:]
	plaintext, err := t.gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(plaintext), nil
}

func (t *translator) EncryptBatch(toEncrypt [][]byte) ([]string, error) {
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
