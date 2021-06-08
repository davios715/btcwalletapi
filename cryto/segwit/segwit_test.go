package segwit

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"strconv"
	"testing"
)

func TestGetAddress(t *testing.T) {
	var seed = []byte{244, 184, 4, 62, 59, 59, 77, 11, 158, 60, 124, 218, 129, 214, 134, 140, 51, 26, 174, 204, 128, 85, 93, 199, 178, 208, 237, 206, 107, 115, 234, 80, 169, 29, 103, 88, 111, 116, 97, 205, 70, 202, 204, 238, 110, 36, 10, 89, 138, 154, 170, 48, 99, 205, 217, 190, 198, 90, 61, 36, 211, 170, 85, 27}
	var address, err = GetAddress(seed, 0x80000000+84, 0x80000000, 0x80000000, 0, 0)

	var expected = "bc1qjslpmfrrmhu4hwmsd7an9n3hctys6452xln2ek"
	assert.NoError(t, err, "Expected no error: valid input")
	assert.Equal(t, expected, address, "Incorrect address")
}

func TestParseDerivationPath(t *testing.T) {
	var path = ""
	var _, err = ParseDerivationPath(path)

	assert.Error(t, err, "Expected error: empty path")

	path = "m/84'/0'/0'/0"
	_, err = ParseDerivationPath(path)

	assert.Error(t, err, "Expected error: short path")

	path = "m/84'/a'/0'/0/0"
	_, err = ParseDerivationPath(path)

	assert.Error(t, err, "Expected error: invalid character")

	path = fmt.Sprintf("m/84'/a'/%s'/0/0", strconv.FormatUint(math.MaxUint64, 10))
	_, err = ParseDerivationPath(path)

	assert.Error(t, err, "Expected error: component out of allowed range")

	path = "m/84'/-1'/0'/0/0"
	_, err = ParseDerivationPath(path)

	assert.Error(t, err, "Expected error: negative component")

	path = "m/50'/0'/0'/0/0"
	_, err = ParseDerivationPath(path)

	assert.Error(t, err, "Expected error: unsupported purpose")

	path = "m/50'/2'/0'/0/0"
	_, err = ParseDerivationPath(path)

	assert.Error(t, err, "Expected error: unsupported coin type")

	path = "m/84'/0'/0'/0/0"
	d, err := ParseDerivationPath(path)
	var expected = []uint32{0x80000054, 0x80000000, 0x80000000, 0x0, 0x0}

	assert.NoError(t, err, "Expected no error: valid path")
	assert.Equal(t, expected, d, "Incorrect components")
}
