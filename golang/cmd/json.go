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
	"strings"

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
		var plaintext [][]byte
		var output_csv string

		red := color.New(color.FgRed)
		green := color.New(color.FgGreen)

		this, err := os.Getwd()
		if err != nil {
			red.Printf("[ERROR]: %s", err)
			os.Exit(1)
		}

		k := cmd.Flag("key").Value.String()
		password := cmd.Flag("password").Value.String()
		path := cmd.Flag("path").Value.String()
		output := cmd.Flag("output").Value.String()
		recursive := cmd.Flag("recursive").Value.String()

		if recursive == "" {
			red.Println("[ERROR]: Missing argument '--recursive' in list.")
			os.Exit(1)
		}
		if path == "" {
			red.Println("[ERROR]: Missing argument '--path' in list.")
			os.Exit(1)
		}
		if output == "" {
			red.Println("[ERROR]: Missing argument '--output' in list.")
			os.Exit(1)
		}

		if recursive == "no" {
			content, err := os.ReadFile(path)
			if err != nil {
				red.Printf("[ERROR]: %s", err)
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
			text, err := decryptor.WithKey(key, dataByte, ivByte)
			plaintext = append(plaintext, text)
			if err != nil {
				red.Println("[ERROR]: Incorrect Password. Maybe you'd forget '--key' or '--password' argument.")
				os.Exit(1)
			}
		}

		if recursive == "yes" {
			files, err := misc.PathInfo(path, ".json")
			if err != nil {
				red.Printf("[ERROR]: %s", err)
				os.Exit(1)
			}

			if len(files) <= 1 {
				red.Println("[ERROR]: Found 1 file, expected more than 1.")
				os.Exit(1)
			}

			for i := 0; i < len(files); i++ {
				content, err := os.ReadFile(files[i])
				if err != nil {
					red.Printf("[ERROR]: %s", err)
					os.Exit(1)
				}

				json.Unmarshal(content, &payload)

				ivByte, _ := base64.StdEncoding.DecodeString(payload.Iv)
				saltByte, _ := base64.StdEncoding.DecodeString(payload.Salt)
				dataByte, _ := base64.StdEncoding.DecodeString(payload.Data)

				if k != "" {
					kss, err := os.ReadFile(k)
					lines := strings.Split(string(kss), "\n")
					if err != nil {
						red.Printf("[ERROR]: %s", err)
						os.Exit(1)
					}

					if len(lines) <= 1 {
						red.Printf("[ERROR]: Found %d files, then more than 1 key is expected.", len(files))
						os.Exit(1)
					}

					for _, ks := range lines {
						key = misc.KeyFromPassword([]byte(ks), saltByte)
						text, err := decryptor.WithKey(key, dataByte, ivByte)
						plaintext = append(plaintext, text)
						if err == nil {
							break
						}
					}
				}
				if password != "" {
					pswds, err := os.ReadFile(password)
					lines := strings.Split(string(pswds), "\n")
					if err != nil {
						red.Printf("[ERROR]: %s", err)
						os.Exit(1)
					}

					if len(lines) <= 1 {
						red.Printf("[ERROR]: Found %d files, then more than 1 password is expected.", len(files))
						os.Exit(1)
					}

					for _, pswd := range lines {
						key = misc.KeyFromPassword([]byte(pswd), saltByte)
						text, err := decryptor.WithKey(key, dataByte, ivByte)
						plaintext = append(plaintext, text)
						if err == nil {
							break
						}
					}
				}

				if len(plaintext) == 0 {
					red.Println("[ERROR]: No vault json could be decrypted.")
					os.Exit(1)
				}
			}
		}

		if strings.Contains(output, "/") {
			output_csv = fmt.Sprintf("%s/output.csv", output)
		} else if !strings.Contains(output, "/") {
			output_csv = fmt.Sprintf("%s/output.csv", this)
		}
		mkerr := os.Mkdir(output, 0755)
		if !os.IsExist(mkerr) {
			red.Printf("[ERROR]: %s", mkerr)
			os.Exit(1)
		}

		csv_file, err := os.OpenFile(output_csv, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			red.Printf("[ERROR]: %s", err)
			os.Exit(1)
		}

		defer csv_file.Close()
		fileInfo, err := os.Stat(output_csv)
		if err != nil {
			red.Printf("[ERROR]: %s", err)
			os.Exit(1)
		}

		writer := csv.NewWriter(csv_file)
		if fileInfo.Size() == 0 {
			crecord := []string{"Mnemonic", "HDPath"}
			wterr := writer.Write(crecord)
			if wterr != nil {
				red.Printf("[ERROR]: %s", wterr)
				os.Exit(1)
			}
		}

		for _, each := range plaintext {
			json.Unmarshal(each, &vault)
			record := []string{string(vault[0].Data.Mnemonic), vault[0].Data.HDPath}
			wterr := writer.Write(record)
			if wterr != nil {
				red.Printf("[ERROR]: %s", wterr)
				os.Exit(1)
			}

			writer.Flush()
		}

		green.Println("[INFO]: Successfuly saved CSV with new values!")
	},
}

func init() {
	rootCmd.AddCommand(jsonCmd)
	jsonCmd.Flags().StringP("key", "k", "", "PBKDF2 derived key if you have any")
	jsonCmd.Flags().StringP("output", "o", "", "Path to where you wanna that be saved CSV file")
	jsonCmd.Flags().StringP("path", "p", "", "Path to log or vault, folder or file")
	jsonCmd.Flags().StringP("password", "w", "", "Password of your Metamask wallet")
	jsonCmd.Flags().StringP("recursive", "r", "", "Iterate over all files in the specified path")
	jsonCmd.PersistentFlags().String("json", "", "Usage: vault-decryptor json [--r] [--pass] [--path]")
}
