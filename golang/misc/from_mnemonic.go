package misc

import (
	color "github.com/fatih/color"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"github.com/ethereum/go-ethereum"
	"os"
)

type PrivateKey struct {
	Mnemonic string `json:"mnemonic"`
	Key      string `json:"key"`
}

func FromMnemonic(mnemonic string, password string) ([]string, error) {
	red := color.New(color.FgRed)
	seed := bip39.NewSeed(mnemonic, password)
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		red.Sprintf("[ERROR]: %s", err)
		os.Exit(1)
	}

	privateKey, err := masterKey.NewChildKey(0)
	if err != nil {
		red.Sprintf("[ERROR]: %s", err)
		os.Exit(1)
	}

	derivedKey, err := privateKey.NewChildKey(0)
	if err != nil {
		red.Sprintf("[ERROR]: %s", err)
		os.Exit(1)
	}

	a := account.Address.Hex()
	p := derivedKey.String()

	return []string{a, p}, nil
}
