package app

import (
	"crytowallet/config"
	btcwalletapi "crytowallet/routes/btc/walletapi"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// App Web app struct
type App struct {
	// Router
	router *mux.Router
	// Configuration
	config config.Config
}

func (a *App) GetRouter() *mux.Router{
	return a.router
}

func (a *App) Run(){
	fmt.Println("Starting the application...")

	fmt.Printf("Listening at :%s...\n", a.config.Application.HttP.Port)
	var err = http.ListenAndServe(
		fmt.Sprintf(":%s", a.config.Application.HttP.Port),
		handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}))(a.GetRouter()),
	)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func (a *App) Init(){
	a.register()
}

func (a *App) register(){
	// Register wallet api
	var api = btcwalletapi.BTCWalletAPI{}
	api.Register(a)
}

func NewApp() App {
	r := mux.NewRouter()

	var conf, err = config.GetConfig()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	return App{
		router: r,
		config: conf,
	}
}