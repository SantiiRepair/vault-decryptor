package misc

import (
	color "github.com/fatih/color"
	"github.com/miguelmota/go-ethereum-hdwallet"
	"os"
)

func FromMnemonic(mnemonic string, hdpath string) ([]string, error) {
	red := color.New(color.FgRed)
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		red.Sprintf("[ERROR]: ", err)
	}

	path := hdwallet.MustParseDerivationPath(hdpath)
	account, err := wallet.Derive(path, false)
	if err != nil {
		red.Sprintf("[ERROR]: ", err)
	}

	address, ethPrivateKey := account.Address.Hex()
	return []string{address, ethPrivateKey}, nil
}
