package main

import "github.com/Torbatti/neshanak"

func main() {
	var app *neshanak.Neshanak

	app = neshanak.New()

	app.Start()
}
