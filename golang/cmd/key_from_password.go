package cmd

import (
	"crypto/sha256"
	"golang.org/x/crypto/pbkdf2"
)

func keyFromPassword(passwordByte []byte, saltByte []byte) []byte {
	key := pbkdf2.Key(passwordByte, saltByte, 10000, 32, sha256.New)
	return key
}
