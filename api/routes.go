package api

import (
	"github.com/ck33122/news-aggregator/app"
	"github.com/ck33122/news-aggregator/app/actions"
)

func Setup(server *app.Server) {
	server.GetHandler("/channels", getChannels)
	server.GetHandler("/channels/<id>", getChannel)
}

func getChannels(ctx app.RequestContext) error {
	channels, err := actions.GetAllChannels()
	if err != nil {
		return ctx.WrapActionsError(err)
	}
	// TODO map channels to model
	return ctx.AnswerJson(channels)
}

func getChannel(ctx app.RequestContext) error {
	id, idErr := ctx.UuidParam("id")
	if idErr != nil {
		return idErr
	}
	channel, getErr := actions.GetChannelById(id)
	if getErr != nil {
		return ctx.WrapActionsError(getErr)
	}
	// TODO map channel to model
	return ctx.AnswerJson(channel)
}
