package cipher

type Cryptograph interface {
	Encrypt(data []byte) (string, error)
	EncryptBatch(data [][]byte) ([]string, error)
	Decrypt(data []byte) (string, error)
}
