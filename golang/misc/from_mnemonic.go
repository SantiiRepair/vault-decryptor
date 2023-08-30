package misc

import (
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	color "github.com/fatih/color"
	bip39 "github.com/tyler-smith/go-bip39"
)

func FromMnemonic(mnemonic string) ([]string, error) {
	red := color.New(color.FgRed)
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	privateKey, err := masterKey.ECPrivKey()
	privateKeyECDSA := privateKey.ToECDSA()
	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		red.Println("[ERROR]: Could not convert the public key type.")
		os.Exit(1)
	}

	address, privateKey := crypto.PubkeyToAddress(publicKeyECDSA), hexutil.Encode(crypto.FromECDSA(privateKeyECDSA))[2:]
	return []string{address, privateKey}, nil
}
