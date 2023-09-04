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

			if key != "" && !strings.Contains(key, ".txt") {
				pbkdf2 = []byte(key)
				text, err := decryptor.WithKey(pbkdf2, dataByte, ivByte)
				if err != nil {
					red.Println("[ERROR]: Incorrect Password. Maybe you forgot '--key' or '--password' argument.")
					os.Exit(1)
				}

				plaintext = append(plaintext, text)
				bs64key := base64.StdEncoding.EncodeToString([]byte(key))
				passwords = append(passwords, bs64key)
			}
			if password != "" && !strings.Contains(password, ".txt") {
				pbkdf2 = misc.KeyFromPassword([]byte(password), saltByte)
				text, err := decryptor.WithKey(pbkdf2, dataByte, ivByte)
				if err != nil {
					red.Println("[ERROR]: Incorrect Password. Maybe you forgot '--key' or '--password' argument.")
					os.Exit(1)
				}

				plaintext = append(plaintext, text)
				passwords = append(passwords, password)
			}
			if key != "" && strings.Contains(key, ".txt") {
				fkey, err := os.ReadFile(key)
				if err != nil {
					red.Printf("[ERROR]: %s", err)
					os.Exit(1)
				}

				lines := strings.Split(string(fkey), "\n")

				for _, ks := range lines {
					pbkdf2 = []byte(ks)
					text, err := decryptor.WithKey(pbkdf2, dataByte, ivByte)
					if text != nil && err == nil {
						bs64key := base64.StdEncoding.EncodeToString(pbkdf2)
						plaintext = append(plaintext, text)
						passwords = append(passwords, bs64key)
						break
					}
				}
			} else if password != "" && strings.Contains(password, ".txt") {
				fkey, err := os.ReadFile(password)
				if err != nil {
					red.Printf("[ERROR]: %s", err)
					os.Exit(1)
				}

				lines := strings.Split(string(fkey), "\n")

				for _, pswd := range lines {
					pbkdf2 = misc.KeyFromPassword([]byte(pswd), saltByte)
					text, err := decryptor.WithKey(pbkdf2, dataByte, ivByte)
					if text != nil && err == nil {
						plaintext = append(plaintext, text)
						passwords = append(passwords, pswd)
						break
					}
				}
			}
			if len(plaintext) == 0 {
				red.Printf("[ERROR]: No vault %s file could be decrypted.", ext)
				os.Exit(1)
			}
		}

		if recursive == "yes" {
			files, err := misc.PathInfo(path, fmt.Sprintf(".%s", ext))
			if err != nil {
				red.Printf("[ERROR]: %s", err)
				os.Exit(1)
			}

			if len(files) <= 1 {
				red.Println("[ERROR]: Found 1 file, more than 1 is expected.")
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
					if err != nil {
						red.Printf("[ERROR]: %s", err)
						os.Exit(1)
					}

					lines := strings.Split(string(kss), "\n")

					if len(lines) <= 1 {
						red.Printf("[ERROR]: Found %d files, then more than 1 key is expected.", len(files))
						os.Exit(1)
					}

					for _, ks := range lines {
						pbkdf2 = []byte(ks)
						text, err := decryptor.WithKey(pbkdf2, dataByte, ivByte)
						if text != nil && err == nil {
							bs64key := base64.StdEncoding.EncodeToString(pbkdf2)
							plaintext = append(plaintext, text)
							passwords = append(passwords, bs64key)
							break
						}
					}
				}

				if password != "" {
					pswds, err := os.ReadFile(password)
					if err != nil {
						red.Printf("[ERROR]: %s", err)
						os.Exit(1)
					}

					lines := strings.Split(string(pswds), "\n")

					if len(lines) <= 1 {
						red.Printf("[ERROR]: Found %d files, then more than 1 password is expected.", len(files))
						os.Exit(1)
					}

					for _, pswd := range lines {
						pbkdf2 = misc.KeyFromPassword([]byte(pswd), saltByte)
						text, err := decryptor.WithKey(pbkdf2, dataByte, ivByte)
						if text != nil && err == nil {
							plaintext = append(plaintext, text)
							passwords = append(passwords, pswd)
							break
						}
					}
				}

				if len(plaintext) == 0 {
					fmt.Println(plaintext)
					red.Printf("[ERROR]: No vault %s file could be decrypted.", ext)
					os.Exit(1)
				}
			}
		}

		if output == "." {
			output_csv = fmt.Sprintf("%s/output.csv", this)
		} else if output != "." {
			output_csv = fmt.Sprintf("%s/output.csv", output)
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
			crecord := []string{"Password", "Address", "Mnemonic", "PrivateKey"}
			wterr := writer.Write(crecord)
			if wterr != nil {
				red.Printf("[ERROR]: %s", wterr)
				os.Exit(1)
			}
		}

		for i, each := range plaintext {
			json.Unmarshal(each, &vault)
			asoc, err := misc.FromMnemonic(string(vault[0].Data.Mnemonic), passwords[i])
			if err != nil {
				red.Printf("[ERROR]: %s", err)
				os.Exit(1)
			}

			record := []string{passwords[i], asoc[0], string(vault[0].Data.Mnemonic), asoc[1]}
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
