package walletapi

import "github.com/gorilla/mux"

// BTCWalletAPI struct to build the DI
type BTCWalletAPI struct {
	app app
}

type app interface {
	GetRouter() *mux.Router
}

// Register register routes in an app and reserve for DI
func (api *BTCWalletAPI) Register(a app) {
	// @Title Wallet API
	// @Version 1.0
	// @Description A BTC Wallet API
	// @BasePath /api/v1/btc/wallet/
	apiV1 := a.GetRouter().PathPrefix("/api/v1/btc/wallet").Subrouter()

	// CreateMnemonic
	// @Summary		Generate a mnemonic words
	// @Description Generate a random mnemonic words following BIP39 standard
	// @Produce		json
	// @Success		200 (object) http.response.Mnemonic
	// @Failure		500 (object) http.response.ErrorResponse
	apiV1.HandleFunc("/mnemonic", api.CreateMnemonic).Methods("GET")

	// CreateHDSegWitAddress
	// @Summary		Create a mnemonic words
	// @Description Generate a Hierarchical Deterministic (HD) Segregated Witness (SegWit)
	//				bitcoin address from a given seed and path
	// @Accept		json http.request.HDSegWit
	// @Produce		json
	// @Success		200 (object) http.response.Address
	// @Failure		409 (object) http.response.ErrorResponse
	// @Failure		500 (object) http.response.ErrorResponse
	apiV1.HandleFunc("/hd/segwit", api.CreateHDSegWitAddress).Methods("POST")

	// CreateMultiSigP2SHAddress
	// @Summary		Create a mnemonic words
	// @Description Generate an n-out-of-m Multisignature (multi-sig) Pay-To-Script-Hash (P2SH)
	//				bitcoin address, where n, m and public keys can be specified
	// @Accept		json http.request.MultiSig
	// @Produce		json
	// @Success		200 (object) http.response.Address
	// @Failure		409 (object) http.response.ErrorResponse
	// @Failure		500 (object) http.response.ErrorResponse
	apiV1.HandleFunc("/multisig", api.CreateMultiSigP2SHAddress).Methods("POST")
}