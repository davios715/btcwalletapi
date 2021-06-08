package main

import (
	"btcwalletapi/app"
)

func main() {
	var a = app.NewApp()

	//Initial app
	a.Init()

	//Run app
	a.Run()
}
