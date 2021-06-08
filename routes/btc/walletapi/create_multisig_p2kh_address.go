package walletapi

import (
	"btcwalletapi/cryto/multisig"
	"btcwalletapi/http/request"
	"btcwalletapi/http/response"
	"encoding/json"
	"log"
	"net/http"
)

// CreateMultiSigP2SHAddress handle n-out-of-m Multisignature (multi-sig) Pay-To-Script-Hash (P2SH) bitcoin address
func (api *BTCWalletAPI) CreateMultiSigP2SHAddress(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	var reqBody request.MultiSig
	json.NewDecoder(req.Body).Decode(&reqBody)

	// create address
	address, _, err := multisig.GenerateAddress(reqBody.M, reqBody.N, reqBody.PublicKeys)
	if err != nil {
		log.Println(err)
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(response.GetResponse(response.ErrInvalidInput))
		return
	}

	json.NewEncoder(res).Encode(response.Address{
		Address: address,
	})
}
