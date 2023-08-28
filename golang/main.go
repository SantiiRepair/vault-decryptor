package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
	"github.com/spf13/cobra"
	"github.com/fatih/color"
)

type Payload struct {
	Data string `json:"data"`
	Salt string `json:"salt"`
	Iv   string `json:"iv"`
}

type Vault struct {
	Type string `json:"type"`
	Data struct {
		Mnemonic         []byte `json:"mnemonic"`
		NumberOfAccounts int    `json:"numberOfAccounts"`
		HDPath           string `json:"hdPath"`
	} `json:"data"`
}

func main() {
	var payload Payload
	var wallet = []byte(``)
	var files []string
	json.Unmarshal(wallet, &payload)

	mode := flag.CommandLine.String("mode", "", "Run tool as, log or vault mode")
	password := flag.CommandLine.String("password", "", "Password of asoc metamask")
	path := flag.CommandLine.String("path", "", "Path to log or vault, folder or file")
	flag.Parse()

	iv, _ := base64.StdEncoding.DecodeString(payload.Iv)
	salt, _ := base64.StdEncoding.DecodeString(payload.Salt)
	data, _ := base64.StdEncoding.DecodeString(payload.Data)

	if *mode == "vault" {
		glob, err := pathInfo(*path, []string{".log", ".json"})
		if err != nil {
			fmt.Println(err)
		}
		files = append(files, glob...)
	}

	key := pbkdf2.Key([]byte(*password), salt, 10000, 32, sha256.New)
	block, _ := aes.NewCipher(key)

	// Initialization Vector size is fixed to 16 bytes.
	gcm, _ := cipher.NewGCMWithNonceSize(block, len(iv))
	plaintext, err := gcm.Open(nil, iv, data, nil)

	if err != nil {
		panic(err)
	}

	var vault []Vault
	json.Unmarshal(plaintext, &vault)
	fmt.Println(vault)
	fmt.Println(string(vault[0].Data.Mnemonic))
}
