package decryptor

import (
	"fmt"
	"crypto/aes"
	"crypto/cipher"
)

func WithKey(key []byte,data []byte, iv []byte) []byte {
	block, _ := aes.NewCipher(key)

	// Initialization Vector size is fixed to 16 bytes.
	gcm, _ := cipher.NewGCMWithNonceSize(block, len(iv))
	plaintext, err := gcm.Open(nil, iv, data, nil)
	if err != nil {
		fmt.Println(err)
	}

	return plaintext
}