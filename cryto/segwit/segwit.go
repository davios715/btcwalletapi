package segwit

import (
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/tyler-smith/go-bip32"
	"math"
	"math/big"
	"strings"
)

type Purpose = uint32

const (
	// PurposeBIP44 44' BIP44
	PurposeBIP44 Purpose = 0x8000002C
	// PurposeBIP49 49' BIP49
	PurposeBIP49 Purpose = 0x80000031
	// PurposeBIP84 84' BIP84
	PurposeBIP84 Purpose = 0x80000054
)

var supportedPurpose = []Purpose{
	PurposeBIP44,
	PurposeBIP49,
	PurposeBIP84,
}

type CoinType = uint32

const (
	CoinTypeBTC CoinType = 0x80000000
)

var supportedCoinTyped = []CoinType{
	CoinTypeBTC,
}

const (
	// Apostrophe 0'
	Apostrophe uint32 = 0x80000000
)

var (
	ErrEmptyPath           = errors.New("empty derivation path")
	ErrInvalidPathPrefix   = errors.New("use 'm/' prefix for absolute paths")
	ErrInvalidPath         = errors.New("invalid derivation path")
	ErrInvalidComponent    = errors.New("invalid component")
	ErrComponentOutOfRange = fmt.Errorf("component out of allowed range [0, %d]", math.MaxUint32)
	ErrUnsupportedCoinType = errors.New("unsupported coinType")
	ErrUnsupportedPurpose  = errors.New("unsupported purpose")
)

type Key struct {
	path     string
	bip32Key *bip32.Key
}

func (k *Key) encode(compress bool) (wif, address, segwitBech32, segwitNested string, err error) {
	prvKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), k.bip32Key.Key)
	return generateFromBytes(prvKey, compress)
}

type KeyManager struct {
	seed []byte
	keys map[string]*bip32.Key
}

func newKeyManager(seed []byte) (*KeyManager, error) {
	km := &KeyManager{
		seed: seed,
		keys: make(map[string]*bip32.Key, 0),
	}
	return km, nil
}

func (km *KeyManager) getSeed() []byte {
	return km.seed
}

func (km *KeyManager) getKey(path string) (*bip32.Key, bool) {
	key, ok := km.keys[path]
	return key, ok
}

func (km *KeyManager) setKey(path string, key *bip32.Key) {
	km.keys[path] = key
}

func (km *KeyManager) getMasterKey() (*bip32.Key, error) {
	path := "m"

	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	key, err := bip32.NewMasterKey(km.getSeed())
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return key, nil
}

func (km *KeyManager) getPurposeKey(purpose uint32) (*bip32.Key, error) {
	path := fmt.Sprintf(`m/%d'`, purpose-Apostrophe)

	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	parent, err := km.getMasterKey()
	if err != nil {
		return nil, err
	}

	key, err = parent.NewChildKey(purpose)
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return key, nil
}

func (km *KeyManager) getCoinTypeKey(purpose, coinType uint32) (*bip32.Key, error) {
	path := fmt.Sprintf(`m/%d'/%d'`, purpose-Apostrophe, coinType-Apostrophe)

	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	parent, err := km.getPurposeKey(purpose)
	if err != nil {
		return nil, err
	}

	key, err = parent.NewChildKey(coinType)
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return key, nil
}

func (km *KeyManager) getAccountKey(purpose, coinType, account uint32) (*bip32.Key, error) {
	path := fmt.Sprintf(`m/%d'/%d'/%d'`, purpose-Apostrophe, coinType-Apostrophe, account-Apostrophe)

	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	parent, err := km.getCoinTypeKey(purpose, coinType)
	if err != nil {
		return nil, err
	}

	key, err = parent.NewChildKey(account)
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return key, nil
}

func (km *KeyManager) getChangeKey(purpose, coinType, account, change uint32) (*bip32.Key, error) {
	path := fmt.Sprintf(`m/%d'/%d'/%d'/%d`, purpose-Apostrophe, coinType-Apostrophe, account-Apostrophe, change)

	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	parent, err := km.getAccountKey(purpose, coinType, account)
	if err != nil {
		return nil, err
	}

	key, err = parent.NewChildKey(change)
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return key, nil
}

func (km *KeyManager) GetKey(purpose, coinType, account, change, index uint32) (*Key, error) {
	path := fmt.Sprintf(`m/%d'/%d'/%d'/%d/%d`, purpose-Apostrophe, coinType-Apostrophe, account-Apostrophe, change, index)

	key, ok := km.getKey(path)
	if ok {
		return &Key{path: path, bip32Key: key}, nil
	}

	parent, err := km.getChangeKey(purpose, coinType, account, change)
	if err != nil {
		return nil, err
	}

	key, err = parent.NewChildKey(index)
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return &Key{path: path, bip32Key: key}, nil
}

func generateFromBytes(prvKey *btcec.PrivateKey, compress bool) (wif, address, segwitBech32, segwitNested string, err error) {
	// generate the wif(wallet import format) string
	btcwif, err := btcutil.NewWIF(prvKey, &chaincfg.MainNetParams, compress)
	if err != nil {
		return "", "", "", "", err
	}
	wif = btcwif.String()

	// generate a normal p2pkh address
	serializedPubKey := btcwif.SerializePubKey()
	addressPubKey, err := btcutil.NewAddressPubKey(serializedPubKey, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", "", "", err
	}
	address = addressPubKey.EncodeAddress()

	// generate a normal p2wkh address from the pubkey hash
	witnessProg := btcutil.Hash160(serializedPubKey)
	addressWitnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", "", "", err
	}
	segwitBech32 = addressWitnessPubKeyHash.EncodeAddress()

	// generate an address which is
	// backwards compatible to Bitcoin nodes running 0.6.0 onwards, but
	// allows us to take advantage of segwit's scripting improvments,
	// and malleability fixes.
	serializedScript, err := txscript.PayToAddrScript(addressWitnessPubKeyHash)
	if err != nil {
		return "", "", "", "", err
	}
	addressScriptHash, err := btcutil.NewAddressScriptHash(serializedScript, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", "", "", err
	}
	segwitNested = addressScriptHash.EncodeAddress()

	return wif, address, segwitBech32, segwitNested, nil
}

// ParseDerivationPath parse a string absolute path to a component slice
func ParseDerivationPath(path string) ([]uint32, error) {
	var result []uint32

	// Handle absolute or relative paths
	components := strings.Split(path, "/")
	switch {
	case len(components) == 0:
		return nil, ErrEmptyPath

	case strings.TrimSpace(components[0]) != "m":
		return nil, ErrInvalidPathPrefix

	default:
		components = components[1:]
	}
	// All remaining components are relative, append one by one
	if len(components) != 5 {
		return nil, ErrInvalidPath
	}
	for _, component := range components {
		// Ignore any user added whitespace
		component = strings.TrimSpace(component)
		var value uint32

		// Handle hardened paths
		if strings.HasSuffix(component, "'") {
			value = 0x80000000
			component = strings.TrimSpace(strings.TrimSuffix(component, "'"))
		}
		// Handle the non hardened component
		bigval, ok := new(big.Int).SetString(component, 0)
		if !ok {
			return nil, ErrInvalidComponent
		}
		max := math.MaxUint32 - value
		if bigval.Sign() < 0 || bigval.Cmp(big.NewInt(int64(max))) > 0 {
			if value == 0 {
				return nil, ErrComponentOutOfRange
			}
			return nil, ErrComponentOutOfRange
		}
		value += uint32(bigval.Uint64())

		// Append and repeat
		result = append(result, value)
	}

	if !contains(supportedCoinTyped, result[1]) {
		return nil, ErrUnsupportedCoinType
	}

	if !contains(supportedPurpose, result[0]) {
		return nil, ErrUnsupportedPurpose
	}

	return result, nil
}

func contains(s []uint32, e uint32) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// GetAddress generate a Hierarchical Deterministic Segregated Witness bitcoin address
func GetAddress(seed []byte, purpose, coinType, account, change, index uint32) (string, error) {
	var err error

	km, err := newKeyManager(seed)
	if err != nil {
		return "", err
	}

	key, err := km.GetKey(purpose, coinType, account, change, index)
	if err != nil {
		return "", err
	}

	_, address, segwitBech32, segwitNested, err := key.encode(false)
	if err != nil {
		return "", err
	}

	var result string
	switch purpose {
	case PurposeBIP44:
		result = address
	case PurposeBIP49:
		result = segwitNested
	case PurposeBIP84:
		result = segwitBech32
	}

	return result, nil
}
