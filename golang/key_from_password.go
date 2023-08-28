package main

import (
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/crypto/pbkdf2"
)

func keyFromPassword(password string, saltByte []byte) []byte {
	key := pbkdf2.Key([]byte(password), saltByte, 10000, 32, sha256.New)
	return key
}
