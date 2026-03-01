package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

type Cryptor struct {
	key []byte
}

type EncryptionResult struct {
	Ciphertext string
	Nonce      string
}

func NewCryptor(hexKey string) Cryptor {
	key, _ := hex.DecodeString(hexKey)
	return Cryptor{key: key}
}

func (c Cryptor) Encrypt(plaintext string) (EncryptionResult, error) {
	bytetext := []byte(plaintext)

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return EncryptionResult{}, fmt.Errorf("Cipher creation failed: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return EncryptionResult{}, fmt.Errorf("GCM creation failed: %w", err)
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return EncryptionResult{}, fmt.Errorf("Nonce creation failed: %w", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, bytetext, nil)
	return EncryptionResult{hex.EncodeToString(ciphertext), hex.EncodeToString(nonce)}, nil
}

func (c Cryptor) Decrypt(ec EncryptionResult) (string, error) {
	ciphertext, _ := hex.DecodeString(ec.Ciphertext)
	nonce, _ := hex.DecodeString(ec.Nonce)

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", fmt.Errorf("Cipher creation failed: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("GCM creation failed: %w", err)
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("GCM decryption failed: %w", err)
	}

	return string(plaintext), nil
}
