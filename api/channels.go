package api

import (
	"github.com/ck33122/news-aggregator/app"
)

var (
	emptyChannels = make([]app.Channel, 0)
)

func getChannels(ctx app.RequestContext) error {
	var channels []app.Channel
	result, err := getAllChannelsStatement.Query(&channels)
	if err != nil {
		return ctx.WrapDBError("get channels", err)
	}
	if result.RowsReturned() == 0 {
		channels = emptyChannels
	}
	return ctx.AnswerJson(channels)
}

func getChannel(ctx app.RequestContext) error {
	id, err := ctx.UuidParam("id")
	if err != nil {
		return err
	}
	var channel app.Channel
	_, err = getChannelByIdStatement.QueryOne(&channel, id)
	if err != nil {
		return ctx.WrapDBError("get channel by id", err)
	}
	return ctx.AnswerJson(&channel)
}
