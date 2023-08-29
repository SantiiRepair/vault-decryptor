/*
Copyright © 2023 Santiago Ramirez
*/
package cmd

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"github.com/SantiiRepair/vault-decryptor/decryptor"
	"github.com/SantiiRepair/vault-decryptor/misc"
	color "github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "Sub-module that rescue seed phrase from json vault.",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		var key []byte
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
				red("[ERROR]: %s", err)
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
				red("[ERROR]: %s", err)
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
		csv_file, err := os.Open("../csv/account.csv")
		if err == os.ErrNotExist {
			os.Create("../csv/account.csv")
		}
		if err != os.ErrNotExist {
			red("[ERROR]: %s", err)
		}
		writer := csv.NewWriter(csv_file)
		writer.Write([]string(plaintext[:]))
		green("[INFO]: Successfuly saved CSV with new values!")
	},
}

func init() {
	rootCmd.AddCommand(jsonCmd)
	jsonCmd.Flags().StringP("key", "k", "", "PBKDF2 derived key if you have any")
	jsonCmd.Flags().StringP("path", "p", "", "Path to log or vault, folder or file")
	jsonCmd.Flags().StringP("password", "w", "", "Password of your Metamask wallet")
	jsonCmd.Flags().StringP("recursive", "r", "", "Iterate over all files in the specified path")
	jsonCmd.PersistentFlags().String("vault-decryptor", "", "Usage: vault-decryptor json [--r] [--pass] [--path]")
}
