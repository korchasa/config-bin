package aes

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "errors"
    "io"
)

type AESEncryptor struct{}

func NewAESEncryptor() *AESEncryptor {
    return &AESEncryptor{}
}

func (e *AESEncryptor) Encrypt(data string, password string) (string, error) {
    hash := createPasswordHash(password)

    block, err := aes.NewCipher(hash)
    if err != nil {
        return "", err
    }

    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, aesGCM.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    ciphertext := aesGCM.Seal(nonce, nonce, []byte(data), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (e *AESEncryptor) Decrypt(data string, password string) (string, error) {
    hash := createPasswordHash(password)

    block, err := aes.NewCipher(hash)
    if err != nil {
        return "", err
    }

    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    ciphertext, err := base64.StdEncoding.DecodeString(data)
    if err != nil {
        return "", err
    }

    nonceSize := aesGCM.NonceSize()
    if len(ciphertext) < nonceSize {
        return "", errors.New("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}

func createPasswordHash(password string) []byte {
    // Create a new hash.
    h := sha256.New()
    // Write password to hash.
    h.Write([]byte(password))
    // Use the sum as the key.
    key := h.Sum(nil)
    // Return the key as a string.
    return key
}
