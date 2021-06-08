package walletapi

import (
	"bytes"
	"crytowallet/http/response"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestRoute_CreateHDSegWitAddress_ReturnNormal(t *testing.T) {
	var api = BTCWalletAPI{}

	params := struct {
		Seed []byte `json:"seed"`
		Path string `json:"path"`
	}{
		Seed: []byte{244, 184, 4, 62, 59, 59, 77, 11, 158, 60, 124, 218, 129, 214, 134, 140, 51, 26, 174, 204, 128, 85, 93, 199, 178, 208, 237, 206, 107, 115, 234, 80, 169, 29, 103, 88, 111, 116, 97, 205, 70, 202, 204, 238, 110, 36, 10, 89, 138, 154, 170, 48, 99, 205, 217, 190, 198, 90, 61, 36, 211, 170, 85, 27},
		Path: "m/84'/0'/0'/0/0",
	}
	paramsByte, _ := json.Marshal(params)
	var r = httptest.NewRequest("POST", "/", bytes.NewBuffer(paramsByte))
	var w = httptest.NewRecorder()

	api.CreateHDSegWitAddress(w, r)

	var res response.Address
	var err = json.NewDecoder(w.Body).Decode(&res)

	assert.NoError(t, err, "Expected no error: valid response struct")

	var expectedAddress = "bc1qjslpmfrrmhu4hwmsd7an9n3hctys6452xln2ek"

	assert.Equal(t, expectedAddress, res.Address, "Incorrect address")
}

func TestRoute_CreateHDSegWitAddress_ReturnInvalidPathError(t *testing.T) {
	var api = BTCWalletAPI{}

	params := struct {
		Seed []byte `json:"seed"`
		Path string `json:"path"`
	}{
		Seed: []byte{244, 184, 4, 62, 59, 59, 77, 11, 158, 60, 124, 218, 129, 214, 134, 140, 51, 26, 174, 204, 128, 85, 93, 199, 178, 208, 237, 206, 107, 115, 234, 80, 169, 29, 103, 88, 111, 116, 97, 205, 70, 202, 204, 238, 110, 36, 10, 89, 138, 154, 170, 48, 99, 205, 217, 190, 198, 90, 61, 36, 211, 170, 85, 27},
		Path: "m/84'/a'/0'/0/0",
	}
	paramsByte, _ := json.Marshal(params)
	var r = httptest.NewRequest("POST", "/", bytes.NewBuffer(paramsByte))
	var w = httptest.NewRecorder()

	api.CreateHDSegWitAddress(w, r)

	var res response.ErrorResponse
	var err = json.NewDecoder(w.Body).Decode(&res)

	assert.NoError(t, err, "Expected no error: valid response struct")

	var expectedCode = "INVALID_PATH"

	assert.Equal(t, expectedCode, res.Code, "Expected error error: invalid path")
}
