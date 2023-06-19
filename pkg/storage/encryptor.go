package storage

type Encryptor interface {
	Encrypt(UnencryptedData string, password string) (EncryptedData string, err error)
	Decrypt(EncryptedData string, password string) (UnencryptedData string, err error)
}
