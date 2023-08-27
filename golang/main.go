package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
)

var wallet = []byte(``)

func main() {
	var payload Payload
	json.Unmarshal(wallet, &payload)

	iv, _ := base64.StdEncoding.DecodeString(payload.Iv)
	salt, _ := base64.StdEncoding.DecodeString(payload.Salt)
	data, _ := base64.StdEncoding.DecodeString(payload.Data)

	password := "Elmetamask1"

	key := pbkdf2.Key([]byte(password), salt, 10000, 32, sha256.New)
	block, _ := aes.NewCipher(key)

	// Initialization Vector size is fixed to 16 bytes.
	gcm, _ := cipher.NewGCMWithNonceSize(block, len(iv))
	plaintext, err := gcm.Open(nil, iv, data, nil)

	if err != nil {
		panic(err)
	}

	var vault []Vault
	json.Unmarshal(plaintext, &vault)
	fmt.Println(string(vault[0].Data.Mnemonic))
}
