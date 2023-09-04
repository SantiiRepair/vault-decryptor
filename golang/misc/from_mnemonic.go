package misc

import (
	"crypto/ecdsa"
	"os"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	color "github.com/fatih/color"
	bip39 "github.com/tyler-smith/go-bip39"
)

func FromMnemonic(mnemonic string, password string) ([]string, error) {
	red := color.New(color.FgRed)
	seed := bip39.NewSeed(mnemonic, password)
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		red.Sprintf("[ERROR]: %s", err)
		os.Exit(1)
	}

	privateKey, err := masterKey.ECPrivKey()
	privateKeyECDSA := privateKey.ToECDSA()
	publicKey := privateKeyECDSA.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		red.Println("[ERROR]: Could not convert the public key type.")
		os.Exit(1)
	}

	a, p := crypto.PubkeyToAddress(*publicKeyECDSA).String(), hexutil.Encode(crypto.FromECDSA(privateKeyECDSA))[2:]

	return []string{a, p}, nil
}
