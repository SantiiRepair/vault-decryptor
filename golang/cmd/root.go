/*
Copyright © 2023 Santiago Ramirez
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vault-decryptor",
	Short: "A fast, local Metamask Vault Decryptor in the command line.",
	Long:  "Vault Decryptor is a cli tool that allows you to decrypt vault data of Metamask Extension, this work by entering vault data path and password of the wallet extension, then if the data entered in the arguments are correct it creates a csv file with the seed phrases of the wallet."
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vault-decryptor.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

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
