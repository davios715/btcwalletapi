package mnemonic

import (
	"github.com/stretchr/testify/assert"
	"github.com/tyler-smith/go-bip39"
	"testing"
)

func TestGetMnemonic_ReturnNormal(t *testing.T) {
	var result, err = GetMnemonic()

	assert.Empty(t, err)

	_, err = bip39.EntropyFromMnemonic(result)

	assert.NoError(t, err, "Expected no error: valid mnemonic")

	var isEntrophSizeValid = (entropySize%32) == 0 && entropySize >= 128 && entropySize <= 256

	assert.True(t, isEntrophSizeValid, "Invalid entropySize")
}
