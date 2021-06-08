package mnemonic

import (
	"github.com/tyler-smith/go-bip39"
)

const entropySize = 256

// GetMnemonic give a random mnemonic words following BIP39 standard
func GetMnemonic() (string, error)  {
	entropy, err := bip39.NewEntropy(entropySize)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}
