package api

import (
	"github.com/ck33122/news-aggregator/app"
)

const (
	pageSize = 10
)

func Setup(server *app.Server, actions *app.Actions) {
	server.Get("/channels", func(ctx app.RequestContext) error {
		channels, err := actions.GetAllChannels()
		if err != nil {
			return ctx.WrapActionsError(err)
		}
		return ctx.AnswerJson(mapToChannelsM(channels))
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
		return ctx.AnswerJson(mapToChannelM(channel))
	})
	server.Get("/channels/<id>/posts", func(ctx app.RequestContext) error {
		id, idErr := ctx.UuidParam("id")
		if idErr != nil {
			return idErr
		}
		page, pageErr := ctx.IntQueryParam("p")
		if pageErr != nil || page < 0 {
			page = 0
		}
		posts, getErr := actions.GetPostsByChannelOrdered(id, page, pageSize)
		if getErr != nil {
			return ctx.WrapActionsError(getErr)
		}
		return ctx.AnswerJson(mapToPostsM(posts))
	})
	server.Get("/posts", func(ctx app.RequestContext) error {
		page, pageErr := ctx.IntQueryParam("p")
		if pageErr != nil || page < 0 {
			page = 0
		}
		posts, err := actions.GetPostsOrdered(page, pageSize)
		if err != nil {
			return ctx.WrapActionsError(err)
		}
		return ctx.AnswerJson(mapToPostsM(posts))
	})
	server.Get("/posts/<id>", func(ctx app.RequestContext) error {
		id, idErr := ctx.UuidParam("id")
		if idErr != nil {
			return idErr
		}
		post, getErr := actions.GetPostById(id)
		if getErr != nil {
			return ctx.WrapActionsError(getErr)
		}
		return ctx.AnswerJson(mapToPostM(post))
	})
}
