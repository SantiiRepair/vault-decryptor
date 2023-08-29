/*
Copyright © 2023 Santiago Ramirez
*/
package cmd

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

var logCmd = &cobra.Command{
	Use:   "json",
	Short: "Sub-module that rescue seed phrase from json vault.",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		var key []byte
		var vault []Vault
		var payload Payload
		var plaintext []byte

		red := color.New(color.FgRed).PrintFunc()
		green := color.New(color.FgGreen).PrintFunc()

		k := cmd.Flag("key").Value.String()
		password := cmd.Flag("password").Value.String()
		path := cmd.Flag("path").Value.String()
		recursive := cmd.Flag("recursive").Value.String()

		if recursive == "" {
			red("[ERROR]: Missing argument '--recursive' in list.")
			os.Exit(1)
		}
		if path == "" {
			red("[ERROR]: Missing argument '--path' in list.")
			os.Exit(1)
		}

		if recursive == "no" {
			content, err := os.ReadFile(path)
			if err != nil {
				red(err)
				os.Exit(1)
			}

			json.Unmarshal(content, &payload)

			ivByte, _ := base64.StdEncoding.DecodeString(payload.Iv)
			saltByte, _ := base64.StdEncoding.DecodeString(payload.Salt)
			dataByte, _ := base64.StdEncoding.DecodeString(payload.Data)

			if k != "" {
				key = []byte(k)
			}
			if password != "" {
				key = misc.KeyFromPassword([]byte(password), saltByte)
			}
			plaintext, err = decryptor.WithKey(key, dataByte, ivByte)
			if err != nil {
				red("[ERROR]: Incorrect Password. Maybe you'd forget '--key' or '--password' argument.")
				os.Exit(1)
			}
		}

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
			if k != "" {
				key = []byte(k)
			}
			if password != "" {
				key = misc.KeyFromPassword([]byte(password), saltByte)
			}
			plaintext, err = decryptor.WithKey(key, dataByte, ivByte)
			if err != nil {
				red("[ERROR]: Incorrect Password. Maybe you'd forget '--key' or '--password' argument.")
				os.Exit(1)
			}

		}

		json.Unmarshal(plaintext, &vault)
		output, err := json.Marshal(vault)
		if err != nil {
			red(err)
		}

		green(string(output))
		// fmt.Println(string(vault[0].Data.Mnemonic))
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
	logCmd.Flags().StringP("key", "k", "", "PBKDF2 derived key if you have any")
	logCmd.Flags().StringP("path", "p", "", "Path to log or vault, folder or file")
	logCmd.Flags().StringP("password", "w", "", "Password of your Metamask wallet")
	logCmd.Flags().StringP("recursive", "r", "", "Iterate over all files in the specified path")
	logCmd.PersistentFlags().String("vault-decryptor", "", "Usage: vault-decryptor log [--r] [--pass] [--path]")
}

