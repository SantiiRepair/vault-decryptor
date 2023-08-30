/*
Copyright © 2023 Santiago Ramirez
*/
package cmd

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/SantiiRepair/vault-decryptor/decryptor"
	"github.com/SantiiRepair/vault-decryptor/misc"
	color "github.com/fatih/color"
	"github.com/spf13/cobra"
)

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "Sub-module that rescue seed phrase from json vault.",
	Run: func(cmd *cobra.Command, args []string) {
		var key []byte
		var vault []Vault
		var payload Payload
		var plaintext []byte

		red := color.New(color.FgRed).PrintFunc()
		green := color.New(color.FgGreen).PrintFunc()

		this, err := os.Getwd()
		if err != nil {
			red("[ERROR]: ", err)
			os.Exit(1)
		}

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
				red("[ERROR]: ", err)
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

		if recursive == "yes" {
			glob, err := misc.PathInfo(path, "json")
			if err != nil {
				red("[ERROR]: ", err)
				os.Exit(1)
			}

			for _, file := range glob {
				content, err := os.ReadFile(file)
				if err != nil {
					red("[ERROR]: ", err)
					os.Exit(1)
				}

				json.Unmarshal(content, &payload)

				ivByte, _ := base64.StdEncoding.DecodeString(payload.Iv)
				saltByte, _ := base64.StdEncoding.DecodeString(payload.Salt)
				dataByte, _ := base64.StdEncoding.DecodeString(payload.Data)

				if k != "" {
					psw_file, err := os.ReadFile(k)
					if err != nil {
						red("[ERROR]: ", err)
						os.Exit(1)
					}
					for _, keyByte := range psw_file {
						key = misc.KeyFromPassword([]byte{keyByte}, saltByte)
						plaintext, err = decryptor.WithKey(key, dataByte, ivByte)
						if err != nil {
							continue
						}
					}
				}
				if password != "" {
					psw_file, err := os.ReadFile(password)
					if err != nil {
						red("[ERROR]: ", err)
						os.Exit(1)
					}
					for _, passByte := range psw_file {
						key = misc.KeyFromPassword([]byte{passByte}, saltByte)
						plaintext, err = decryptor.WithKey(key, dataByte, ivByte)
						if err != nil {
							continue
						}
					}
				}
			}
		}

		csv_path := fmt.Sprintf("%s/csv/metamask.csv", this)
		mkerr := os.Mkdir(fmt.Sprintf("%s/csv", this), 0755)
		if !os.IsExist(mkerr) {
			red("[ERROR]: ", mkerr)
			os.Exit(1)
		}

		csv_file, err := os.OpenFile(csv_path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			red("[ERROR]: ", err)
			os.Exit(1)
		}

		defer csv_file.Close()
		fileInfo, err := os.Stat(csv_path)
		if err != nil {
			red("[ERROR]: ", err)
			os.Exit(1)
		}

		json.Unmarshal(plaintext, &vault)
		record := []string{string(vault[0].Data.Mnemonic), vault[0].Data.HDPath}
		writer := csv.NewWriter(csv_file)
		if fileInfo.Size() == 0 {
			crecord := []string{"Mnemonic", "HDPath"}
			wterr := writer.Write(crecord)
			if wterr != nil {
				red("[ERROR]: ", wterr)
				os.Exit(1)
			}
		}

		wterr := writer.Write(record)
		if wterr != nil {
			red("[ERROR]: ", wterr)
			os.Exit(1)
		}

		writer.Flush()
		green("[INFO]: Successfuly saved CSV with new values!")
	},
}

func init() {
	rootCmd.AddCommand(jsonCmd)
	jsonCmd.Flags().StringP("key", "k", "", "PBKDF2 derived key if you have any")
	jsonCmd.Flags().StringP("path", "p", "", "Path to log or vault, folder or file")
	jsonCmd.Flags().StringP("password", "w", "", "Password of your Metamask wallet")
	jsonCmd.Flags().StringP("recursive", "r", "", "Iterate over all files in the specified path")
	jsonCmd.PersistentFlags().String("json", "", "Usage: vault-decryptor json [--r] [--pass] [--path]")
}
