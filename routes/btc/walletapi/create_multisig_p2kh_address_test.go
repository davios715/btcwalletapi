package walletapi

import (
	"bytes"
	"crytowallet/http/response"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestRoute_CreateMultiSigP2SHAddress_ReturnNormal(t *testing.T) {
	var api = BTCWalletAPI{}
	params := struct {
		M          int      `json:"m"`
		N          int      `json:"n"`
		PublicKeys []string `json:"public_keys"`
	}{
		M: 2,
		N: 3,
		PublicKeys: []string{"04a882d414e478039cd5b52a92ffb13dd5e6bd4515497439dffd691a0f12af9575fa349b5694ed3155b136f09e63975a1700c9f4d4df849323dac06cf3bd6458cd","046ce31db9bdd543e72fe3039a1f1c047dab87037c36a669ff90e28da1848f640de68c2fe913d363a51154a0c62d7adea1b822d05035077418267b1a1379790187","0411ffd36c70776538d079fbae117dc38effafb33304af83ce4894589747aee1ef992f63280567f52f5ba870678b4ab4ff6c8ea600bd217870a8b4f1f09f3a8e83"},
	}
	paramsByte, _ := json.Marshal(params)
	var r = httptest.NewRequest("POST", "/", bytes.NewBuffer(paramsByte))
	var w = httptest.NewRecorder()

	api.CreateMultiSigP2SHAddress(w, r)

	var res response.Address
	var err = json.NewDecoder(w.Body).Decode(&res)

	assert.NoError(t, err, "Expected no error: valid response struct")

	var expectedAddress = "347N1Thc213QqfYCz3PZkjoJpNv5b14kBd"

	assert.Equal(t, expectedAddress, res.Address, "Incorrect address")
}

func TestRoute_CreateMultiSigP2SHAddress_ReturnInvalidInputError(t *testing.T) {
	var api = BTCWalletAPI{}
	params := struct {
		M          int      `json:"m"`
		N          int      `json:"n"`
		PublicKeys []string `json:"public_keys"`
	}{
		M: 4,
		N: 3,
		PublicKeys: []string{"04a882d414e478039cd5b52a92ffb13dd5e6bd4515497439dffd691a0f12af9575fa349b5694ed3155b136f09e63975a1700c9f4d4df849323dac06cf3bd6458cd","046ce31db9bdd543e72fe3039a1f1c047dab87037c36a669ff90e28da1848f640de68c2fe913d363a51154a0c62d7adea1b822d05035077418267b1a1379790187","0411ffd36c70776538d079fbae117dc38effafb33304af83ce4894589747aee1ef992f63280567f52f5ba870678b4ab4ff6c8ea600bd217870a8b4f1f09f3a8e83"},
	}
	paramsByte, _ := json.Marshal(params)
	var r = httptest.NewRequest("POST", "/", bytes.NewBuffer(paramsByte))
	var w = httptest.NewRecorder()

	api.CreateMultiSigP2SHAddress(w, r)

	var res response.ErrorResponse
	var err = json.NewDecoder(w.Body).Decode(&res)

	assert.NoError(t, err, "Expected no error: valid response struct")

	var expectedCode = "INVALID_INPUT"

	assert.Equal(t, expectedCode, res.Code, "Expected error error: invalid input")
}
