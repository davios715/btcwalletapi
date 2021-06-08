package walletapi

import (
	"btcwalletapi/http/response"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/tyler-smith/go-bip39"
	"net/http/httptest"
	"testing"
)

func TestRoute_CreateMnemonic_ReturnNormal(t *testing.T) {
	var api = BTCWalletAPI{}

	var r = httptest.NewRequest("GET", "/", nil)
	var w = httptest.NewRecorder()

	api.CreateMnemonic(w, r)

	var res response.Mnemonic
	var err = json.NewDecoder(w.Body).Decode(&res)

	assert.NoError(t, err, "Expected no error: valid response struct")

	_, err = bip39.EntropyFromMnemonic(res.Mnemonic)

	assert.NoError(t, err, "Expected no error: valid mnemonic")
}