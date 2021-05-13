package cipher

type Translator interface {
	Encrypt(data []byte) (string, error)
	EncryptBatch(data [][]byte) ([]string, error)
	Decrypt(data []byte) (string, error)
}
