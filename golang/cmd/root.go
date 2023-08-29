/*
Copyright © 2023 Santiago Ramirez
*/
package cmd

import (
	"os"

	color "github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vault-decryptor",
	Short: "A fast, local Metamask Vault Decryptor in the command line.",
	Run: func(cmd *cobra.Command, args []string) {
		magenta := color.New(color.FgMagenta).PrintFunc()
		magenta(ascii)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetUsageTemplate("Usage: vault-decryptor [mode] [--r] [--pass] [--path]")
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
