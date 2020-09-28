package api

import (
	"github.com/ck33122/news-aggregator/app"
)

func Setup(server *app.Server) {
	app.Must(prepareStatements())
	server.
		GetHandler("/channels", getChannels).
		GetHandler("/channels/<id>", getChannel)
}
