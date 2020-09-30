package api

import (
	"github.com/ck33122/news-aggregator/app"
)

func Setup(server *app.Server, actions *app.Actions) {
	server.Get("/channels", func(ctx app.RequestContext) error {
		channels, err := actions.GetAllChannels()
		if err != nil {
			return ctx.WrapActionsError(err)
		}
		// TODO map channels to model
		return ctx.AnswerJson(channels)
	})
	server.Get("/channels/<id>", func(ctx app.RequestContext) error {
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
	})
	server.Get("/posts", func(ctx app.RequestContext) error {
		channels, err := actions.GetAllPosts()
		if err != nil {
			return ctx.WrapActionsError(err)
		}
		// TODO map channels to model
		return ctx.AnswerJson(channels)
	})
}
