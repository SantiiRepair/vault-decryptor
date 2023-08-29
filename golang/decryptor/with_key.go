package decryptor

import (
	"crypto/aes"
	"crypto/cipher"
)

func WithKey(key []byte, data []byte, iv []byte) ([]byte, error) {
	block, _ := aes.NewCipher(key)

	// Initialization Vector size is fixed to 16 bytes.
	gcm, _ := cipher.NewGCMWithNonceSize(block, len(iv))
	plaintext, err := gcm.Open(nil, iv, data, nil)

	return plaintext, err
}
