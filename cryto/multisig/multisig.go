package multisig

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

// OP_1 through OP_16
const (
	OP_1 = 81 + iota
	OP_2 //82
	OP_3 //83
	OP_4 //..
	OP_5
	OP_6
	OP_7
	OP_8
	OP_9
	OP_10
	OP_11
	OP_12
	OP_13
	OP_14 //..
	OP_15 //95
	OP_16 //96
)

// OP codes other than OP_1 through OP_16, used in P2SH Multisig transanctions.
const (
	OP_CHECKMULTISIG = 174
)

var (
	ErrOffendPubKey = errors.New("offending publicKey")
	ErrNRange       = errors.New("N must be between 1 and 7 (inclusive) for valid, standard P2SH multisig transaction as per Bitcoin protocol")
	ErrMRange       = errors.New("M must be between 1 and N (inclusive)")
	ErrNumOfPubKeys = func(n int, m int, numOfPubKeys int) error {
		return fmt.Errorf("need exactly %d public keys to create P2SH address for %d-of-%d multisig transaction. Only %d keys provided", n, m, n, numOfPubKeys)
	}
	ErrEmptyBytes    = errors.New("empty bytes")
	ErrEmptyPubKey   = errors.New("public key cannot be empty")
	ErrInvalidPubKey = errors.New("public key invalid")
)

func GenerateAddress(flagM int, flagN int, publicKeyStrings []string) (string, string, error) {
	var err error

	publicKeys := make([][]byte, len(publicKeyStrings))
	for i, publicKeyString := range publicKeyStrings {
		publicKeyString = strings.TrimSpace(publicKeyString)
		publicKeys[i], err = hex.DecodeString(publicKeyString)
		if err != nil {
			return "", "", ErrOffendPubKey
		}
	}
	// create redeemScript from public keys
	redeemScript, err := newMOfNRedeemScript(flagM, flagN, publicKeys)
	if err != nil {
		return "", "", err
	}
	redeemScriptHash, err := hash160(redeemScript)
	if err != nil {
		return "", "", err
	}
	// get P2SH address by base58 encoding with P2SH prefix 0x05
	P2SHAddress := base58.CheckEncode(redeemScriptHash, 5)
	// get redeemScript in Hex
	redeemScriptHex := hex.EncodeToString(redeemScript)

	return P2SHAddress, redeemScriptHex, err
}

func newMOfNRedeemScript(m int, n int, publicKeys [][]byte) ([]byte, error) {
	// validate inputs
	if n < 1 || n > 7 {
		return nil, ErrNRange
	}
	if m < 1 || m > n {
		return nil, ErrMRange
	}
	if len(publicKeys) != n {
		return nil, ErrNumOfPubKeys(n, m, len(publicKeys))
	}

	// get OP Code for m and n.
	// 81 is OP_1, 82 is OP_2 etc.
	// 80 is not a valid OP_Code, so we floor at 81
	mOPCode := OP_1 + (m - 1)
	nOPCode := OP_1 + (n - 1)

	// multisig redeemScript format:
	// <OP_m> <A pubkey> <B pubkey> <C pubkey>... <OP_n> OP_CHECKMULTISIG
	var redeemScript bytes.Buffer
	redeemScript.WriteByte(byte(mOPCode))
	for _, publicKey := range publicKeys {
		err := isPublicKeyValid(publicKey)
		if err != nil {
			return nil, err
		}
		redeemScript.WriteByte(byte(len(publicKey)))
		redeemScript.Write(publicKey)
	}
	redeemScript.WriteByte(byte(nOPCode))
	redeemScript.WriteByte(byte(OP_CHECKMULTISIG))
	return redeemScript.Bytes(), nil
}

// hash160 performs the same operations as OP_HASH160 in Bitcoin Script
func hash160(data []byte) ([]byte, error) {
	// validate inputs
	if data == nil {
		return nil, ErrEmptyBytes
	}

	var shaHash = sha256.New()
	shaHash.Write(data)
	var hash = shaHash.Sum(nil)
	var ripemd160Hash = ripemd160.New()
	ripemd160Hash.Write(hash)

	return ripemd160Hash.Sum(nil), nil
}

// isPublicKeyValid validate publicKey
func isPublicKeyValid(publicKey []byte) error {
	switch {
	case publicKey == nil:
		return ErrEmptyPubKey
	case len(publicKey) != 65:
		return ErrInvalidPubKey
	case publicKey[0] != byte(4):
		return ErrInvalidPubKey
	default:
		return nil
	}
}
