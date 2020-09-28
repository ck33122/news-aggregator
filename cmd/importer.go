package main

import (
	"github.com/ck33122/news-aggregator/app"
)

func main() {
	app.Init("news-aggregator-importer")
	defer app.Destroy()

	app.GetLog().Debug("running daemon")
	app.GetLog().Info("WOW!")
}
