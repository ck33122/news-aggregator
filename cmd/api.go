package main

import (
	"github.com/ck33122/news-aggregator/api"
	"github.com/ck33122/news-aggregator/app"
)

func main() {
	app.Init("news-aggregator-api")
	defer app.Destroy()
	server := app.NewServer()
	api.Setup(server)
	server.Run(app.GetConfig().Api.Listen)
}
