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

var jsonCmd = &cobra.Command{
	Use:   "vault-decryptor",
	Short: "A fast, local Metamask Vault Decryptor in the command line.",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
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
		if password == "" || k == "" {
			red("[ERROR]: Missing argument '--key' or '--password' in list.")
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

			key := misc.KeyFromPassword([]byte(password), saltByte)
			plaintext = decryptor.WithKey(key, dataByte, ivByte)
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

			key := misc.KeyFromPassword([]byte(password), saltByte)
			plaintext = decryptor.WithKey(key, dataByte, ivByte)

		}

		json.Unmarshal(plaintext, &vault)
		green(vault)
		// fmt.Println(string(vault[0].Data.Mnemonic))
	},
}

func init() {
	rootCmd.AddCommand(jsonCmd)
	jsonCmd.Flags().StringP("key", "k", "", "PBKDF2 derived key if you have any")
	jsonCmd.Flags().StringP("path", "path", "", "Path to log or vault, folder or file")
	jsonCmd.Flags().StringP("password", "pass", "", "Password of your Metamask wallet")
	jsonCmd.Flags().StringP("recursive", "r", "", "Iterate over all files in the specified path")
	jsonCmd.PersistentFlags().String("vault-decryptor", "", "Usage: vault-decryptor json [--r] [--pass] [--path]")
}
