package misc

import (
	"os"

	color "github.com/fatih/color"
	"github.com/miguelmota/go-ethereum-hdwallet"
)

func FromMnemonic(mnemonic string, password string) ([]string, error) {
	red := color.New(color.FgRed)
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		red.Sprintf("[ERROR]: %s", err)
		os.Exit(1)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		red.Sprintf("[ERROR]: %s", err)
		os.Exit(1)
	}

	a := account.Address.Hex()
	p, err := wallet.PrivateKeyHex(account)
	if err != nil {
		red.Sprintf("[ERROR]: %s", err)
		os.Exit(1)
	}

	return []string{a, p}, nil
}
