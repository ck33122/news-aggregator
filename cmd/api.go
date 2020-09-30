package main

import (
	"github.com/ck33122/news-aggregator/api"
	"github.com/ck33122/news-aggregator/app"
	"github.com/ck33122/news-aggregator/app/actions"
)

func main() {
	app.Init("news-aggregator-api")
	defer app.Destroy()
	actions.Init()
	server := app.NewServer()
	api.Setup(server)
	server.Run(app.GetConfig().Api.Listen)
}
