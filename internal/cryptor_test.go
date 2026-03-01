package internal

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestCryptor(t *testing.T) {
	t.Run("successfuly encrypts and decrypts", func(t *testing.T) {
		hexKey := "6bbfec8f5e6a9da42da189e099a90c70"
		cryptor := NewCryptor(hexKey)

		encryptionResult, err := cryptor.Encrypt("dupamaryna")

		assert.NilError(t, err)
		assert.Check(t, encryptionResult.Ciphertext != "dupamaryna")
		assert.Check(t, encryptionResult.Nonce != "")

		plaintext, err := cryptor.Decrypt(encryptionResult)
		assert.NilError(t, err)
		assert.Check(t, plaintext == "dupamaryna")
	})
}
