package main

import (
	"crytowallet/app"
)

func main() {
	var a = app.NewApp()

	//Initial app
	a.Init()

	//Run app
	a.Run()
}
