/*
Copyright © 2023 Santiago Ramirez
*/
package main

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/SantiiRepair/vault-decryptor/decryptor"
	"github.com/SantiiRepair/vault-decryptor/misc"
	color "github.com/fatih/color"
	"github.com/spf13/cobra"
)

func main() {
	var r string
	var mode string
	var pass string
	var path string
	var plaintext []byte
	red := color.New(color.FgRed).PrintFunc()
	green := color.New(color.FgRed).PrintFunc()

	var rootCmd = &cobra.Command{
		Use:   "vault-decryptor",
		Short: "A fast, local Metamask Vault Decryptor in the command line.",
		Long:  "Vault Decryptor is a cli tool that allows you to decrypt vault data of Metamask Extension, this work by entering vault data path and password of the wallet extension, then if the data entered in the arguments are correct it creates a csv file with the seed phrases of the wallet.",
		Run: func(cmd *cobra.Command, args []string) {
			var vault []Vault
			var payload Payload

			if r == "" {
				red("Missing argument '--r' in list.")
				os.Exit(1)
			}
			if mode == "" {
				red("Missing argument '--mode' in list.")
				os.Exit(1)
			}
			if pass == "" && path == "" {
				red("Missing argument '--pass' in list.")
				os.Exit(1)
			}
			if path == "" {
				red("Missing argument '--path' in list.")
				os.Exit(1)
			}

			if r == "no" {
				if mode == "log" {
				}
				if mode == "json" {
					content, err := os.ReadFile(path)
					if err != nil {
						red(err)
						os.Exit(1)
					}

					json.Unmarshal(content, &payload)

					ivByte, _ := base64.StdEncoding.DecodeString(payload.Iv)
					saltByte, _ := base64.StdEncoding.DecodeString(payload.Salt)
					dataByte, _ := base64.StdEncoding.DecodeString(payload.Data)

					key := misc.KeyFromPassword([]byte(pass), saltByte)
					plaintext = decryptor.WithKey(key, dataByte, ivByte)
				}
			}

			if mode == "log" {
				glob, err := misc.PathInfo(path, strings.ToLower(filepath.Ext(path)))
				if err != nil {
					red(err)
				}

				for _, file := range glob {
					content, err := os.ReadFile(file)
					if err != nil {
						red(err)
						os.Exit(1)
					}

					json.Unmarshal(content, &payload)

					ivByte, _ := base64.StdEncoding.DecodeString(payload.Iv)
					saltByte, _ := base64.StdEncoding.DecodeString(payload.Salt)
					dataByte, _ := base64.StdEncoding.DecodeString(payload.Data)

					key := misc.KeyFromPassword([]byte(pass), saltByte)
					plaintext = decryptor.WithKey(key, dataByte, ivByte)
				}
			}

			if mode == "json" {
				glob, err := misc.PathInfo(path, strings.ToLower(filepath.Ext(path)))
				if err != nil {
					red(err)
				}

				for _, file := range glob {
					content, err := os.ReadFile(file)
					if err != nil {
						red(err)
						os.Exit(1)
					}

					json.Unmarshal(content, &payload)

					ivByte, _ := base64.StdEncoding.DecodeString(payload.Iv)
					saltByte, _ := base64.StdEncoding.DecodeString(payload.Salt)
					dataByte, _ := base64.StdEncoding.DecodeString(payload.Data)

					key := misc.KeyFromPassword([]byte(pass), saltByte)
					plaintext = decryptor.WithKey(key, dataByte, ivByte)
				}
			}

			json.Unmarshal(plaintext, &vault)
			green(vault)
			// fmt.Println(string(vault[0].Data.Mnemonic))
		},
	}

	rootCmd.PersistentFlags().String("vault-decryptor", "", "Usage: vault-decryptor [--mode] [--r] [--pass] [--path]")

}

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
