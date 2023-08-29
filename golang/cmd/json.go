/*
Copyright © 2023 Santiago Ramirez

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "Run vault-decryptor in json mode.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Vault Decryptor is a CLI library for Go that help to rescue wallet.
This is a tool to retrieve to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("json called")
	},
}

func init() {
	rootCmd.AddCommand(jsonCmd)
	jsonCmd.PersistentFlags().String("foo", "", "A help for foo")
}
