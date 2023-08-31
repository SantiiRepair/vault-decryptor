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

var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "Sub-module that rescue seed phrase from file (.ldb, .log) vault.",
	Run: func(cmd *cobra.Command, args []string) {
		var pbkdf2 []byte
		var vault []Vault
		var payload Payload
		var passwords []string
		var plaintext [][]byte
		var output_csv string

		red := color.New(color.FgRed)
		green := color.New(color.FgGreen)

		this, err := os.Getwd()
		if err != nil {
			red.Printf("[ERROR]: %s", err)
			os.Exit(1)
		}

		ext := cmd.Flag("ext").Value.String()
		key := cmd.Flag("key").Value.String()
		password := cmd.Flag("password").Value.String()
		path := cmd.Flag("path").Value.String()
		output := cmd.Flag("output").Value.String()
		recursive := cmd.Flag("recursive").Value.String()

		if ext == "" {
			red.Println("[ERROR]: Missing argument '--ext' in list.")
			os.Exit(1)
		}
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

			values := misc.ExtractVaultFromFile(string(content))

			json.Unmarshal(values, &payload)

			ivByte, _ := base64.StdEncoding.DecodeString(payload.Iv)
			saltByte, _ := base64.StdEncoding.DecodeString(payload.Salt)
			dataByte, _ := base64.StdEncoding.DecodeString(payload.Data)

			if key != "" {
				pbkdf2 = []byte(key)
			}
			if password != "" {
				pbkdf2 = misc.KeyFromPassword([]byte(password), saltByte)
			}
			text, err := decryptor.WithKey(pbkdf2, dataByte, ivByte)
			if err != nil {
				red.Println("[ERROR]: Incorrect Password. Maybe you'd forget '--key' or '--password' argument.")
				os.Exit(1)
			}

			plaintext = append(plaintext, text)
			passwords = append(passwords, password)
		}

		if recursive == "yes" {
			files, err := misc.PathInfo(path, fmt.Sprintf(".%s", ext))
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

				values := misc.ExtractVaultFromFile(string(content))

				json.Unmarshal(values, &payload)

				ivByte, _ := base64.StdEncoding.DecodeString(payload.Iv)
				saltByte, _ := base64.StdEncoding.DecodeString(payload.Salt)
				dataByte, _ := base64.StdEncoding.DecodeString(payload.Data)

				if key != "" {
					kss, err := os.ReadFile(key)
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
						pbkdf2 = misc.KeyFromPassword([]byte(ks), saltByte)
						text, err := decryptor.WithKey(pbkdf2, dataByte, ivByte)
						if err == nil {
							break
						}

						plaintext = append(plaintext, text)
						bs64key := base64.StdEncoding.EncodeToString([]byte(ks))
						passwords = append(passwords, bs64key)
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
						pbkdf2 = misc.KeyFromPassword([]byte(pswd), saltByte)
						text, err := decryptor.WithKey(pbkdf2, dataByte, ivByte)
						plaintext = append(plaintext, text)
						passwords = append(passwords, pswd)
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
			crecord := []string{"Password", "Address", "Mnemonic", "PrivateKey", "HDPath"}
			wterr := writer.Write(crecord)
			if wterr != nil {
				red.Printf("[ERROR]: %s", wterr)
				os.Exit(1)
			}
		}

		for i, each := range plaintext {
			json.Unmarshal(each, &vault)
			asoc, err := misc.FromMnemonic(string(vault[0].Data.Mnemonic))
			if err != nil {
				red.Printf("[ERROR]: %s", err)
				os.Exit(1)
			}

			record := []string{passwords[i], asoc[0], string(vault[0].Data.Mnemonic), asoc[1], vault[0].Data.HDPath}
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
	rootCmd.AddCommand(fileCmd)
	fileCmd.Flags().StringP("ext", "e", "", "File extension (ldb, log) of vault")
	fileCmd.Flags().StringP("key", "k", "", "PBKDF2 derived key if you have any")
	fileCmd.Flags().StringP("output", "o", "", "Path to where you wanna that be saved CSV file")
	fileCmd.Flags().StringP("path", "p", "", "Path to log or vault, folder or file")
	fileCmd.Flags().StringP("password", "w", "", "Password of your Metamask wallet")
	fileCmd.Flags().StringP("recursive", "r", "", "Iterate over all files in the specified path")
	fileCmd.PersistentFlags().String("file", "", "Usage: vault-decryptor file [-r] [-w] [-p] [-o] [-e]")
}
