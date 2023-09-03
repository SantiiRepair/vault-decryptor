package misc

import (
	color "github.com/fatih/color"
	"github.com/miguelmota/go-ethereum-hdwallet"
	"os"
)

func FromMnemonic(mnemonic string, hdpath string) ([]string, error) {
	red := color.New(color.FgRed)
	wallet, err := hdwallet.(mnemonic)
	if err != nil {
		red.Sprintf("[ERROR]: %s", err)
		os.Exit(1)
	}

	path := hdwallet.MustParseDerivationPath(hdpath)
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
