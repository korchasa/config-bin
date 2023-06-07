package aes

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestEncryptor(t *testing.T) {
    encryptor := NewAESEncryptor()

    unencryptedData := "test data"
    password := "password"

    encryptedData, err := encryptor.Encrypt(unencryptedData, password)
    assert.NoError(t, err, "Encryption should not return an error")
    assert.NotEqual(t, unencryptedData, encryptedData, "Encrypted data should not be equal to unencrypted data")

    decryptedData, err := encryptor.Decrypt(encryptedData, password)
    assert.NoError(t, err, "Decryption should not return an error")
    assert.Equal(t, unencryptedData, decryptedData, "Decrypted data should be equal to original unencrypted data")
}
