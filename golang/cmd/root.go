package cmd

import (
	"os"

	"github.com/spf13/cobra"
)
var rootCmd = &cobra.Command{
	Use:   "vault-decryptor",
	Short: "A fast, local Metamask Vault Decryptor in the command line.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetUsageTemplate("Usage: vault-decryptor [mode] [-r] [OPTIONS]")
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
