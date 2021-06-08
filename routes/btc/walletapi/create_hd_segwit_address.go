package walletapi

import (
	"btcwalletapi/cryto/segwit"
	"btcwalletapi/http/request"
	"btcwalletapi/http/response"
	"encoding/json"
	"log"
	"net/http"
)

// CreateHDSegWitAddress handle Hierarchical Deterministic (HD) Segregated Witness (SegWit) bitcoin address request
func (api *BTCWalletAPI) CreateHDSegWitAddress(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var reqBody request.HDSegWit
	json.NewDecoder(req.Body).Decode(&reqBody)

	// decode path
	var derivationPath, err = segwit.ParseDerivationPath(reqBody.Path)
	if err != nil {
		log.Println(err)
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(response.GetResponse(response.ErrInvalidPath))
		return
	}

	// create address
	address, err := segwit.GetAddress(
		reqBody.Seed,
		derivationPath[0], derivationPath[1], derivationPath[2], derivationPath[3], derivationPath[4],
	)
	if err != nil {
		log.Println(err)
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(response.GetResponse(response.ErrInternal))
		return
	}

	json.NewEncoder(res).Encode(response.Address{
		Address: address,
	})
}
