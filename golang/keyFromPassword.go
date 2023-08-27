package main

import (
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/crypto/pbkdf2"
)

func keyFromPassword(password string, salt string) []byte {
	saltByte, _ := base64.StdEncoding.DecodeString(salt)
	key := pbkdf2.Key([]byte(password), saltByte, 10000, 32, sha256.New)
	return key
}
