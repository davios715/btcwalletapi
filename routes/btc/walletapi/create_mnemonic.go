package walletapi

import (
	"btcwalletapi/cryto/mnemonic"
	"btcwalletapi/http/response"
	"encoding/json"
	"log"
	"net/http"
)

// CreateMnemonic handle random mnemonic words request, following BIP39 standard
func (api *BTCWalletAPI) CreateMnemonic(res http.ResponseWriter, _ *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	// create mnemonic
	var mnmnic, err = mnemonic.GetMnemonic()
	if err != nil {
		log.Println(err)
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(response.GetResponse(response.ErrInternal))
		return
	}
	json.NewEncoder(res).Encode(response.Mnemonic{
		Mnemonic: mnmnic,
	})
}

