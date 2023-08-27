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

func key_from_password(password: str) {
	return pbkdf2.Key([]byte(password), salt, 10000, 32, sha256.New)
}
