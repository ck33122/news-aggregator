package main

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/go-co-op/gocron"
	"github.com/mmcdole/gofeed"

	"github.com/ck33122/news-aggregator/app"
	"github.com/ck33122/news-aggregator/app/actions"
)

func main() {
	app.Init("news-aggregator-importer")
	defer app.Destroy()
	actions.Init()

	scheduler := gocron.NewScheduler(time.UTC)
	_, err := scheduler.Every(10).Minutes().StartImmediately().Do(updateRssTask)
	app.Must(err)

	scheduler.StartBlocking()
}

func updateRssTask() {
	log := app.GetLog()
	log.Info("running job")

	for _, im := range app.GetConfig().Import {
		log.Info("importing rss", zap.String("address", im.Address), zap.String("id", im.Id.String()))
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		feedParser := gofeed.NewParser()
		feed, err := feedParser.ParseURLWithContext(im.Address, ctx)
		if err != nil {
			log.Error("error parsing rss", zap.Error(err))
			continue
		}

		channel := newChannelFromFeed(&im, feed)
		if synced := syncChannelInfo(channel); !synced {
			continue
		}

		for _, item := range feed.Items {
			syncItem(channel, item)
		}
	}
}

func syncChannelInfo(channelFromFeed app.Channel) bool {
	log := app.GetLog()
	db := app.GetDB()

	tx, err := db.Begin()
	if err != nil {
		log.Error("can't begin transaction", zap.Error(err))
		return false
	}
	defer tx.Rollback()

	dbChannel, actionErr := actions.GetChannelByIdTx(tx, channelFromFeed.Id)
	if actionErr != nil {
		log.Debug("some error getting channel, trying to add", zap.Error(actionErr))
		actionErr = actions.AddChannelTx(tx, channelFromFeed)
		if actionErr != nil {
			log.Error("error adding channel", zap.Error(actionErr))
			return false
		}
		log.Debug("add channel ok")
	} else if dbChannel.Title != channelFromFeed.Title ||
		dbChannel.Description != channelFromFeed.Description ||
		dbChannel.Image != channelFromFeed.Image {
		log.Debug("updating channel info")
		if actionErr = actions.UpdateChannelTx(tx, channelFromFeed); actionErr != nil {
			log.Error("error updating channel", zap.Error(actionErr))
			return false
		}
		log.Debug("update ok")
	} else {
		log.Debug("channel info is up to date")
	}

	if err := tx.Commit(); err != nil {
		log.Error("commit error", zap.Error(err))
	}

	return true
}

func syncItem(channel app.Channel, item *gofeed.Item) {
	log := app.GetLog()
	db := app.GetDB()

	log.Debug("syncing item", zap.String("title", item.Title), zap.String("channelId", channel.Id.String()))

	tx, err := db.Begin()
	if err != nil {
		log.Error("can't begin transaction", zap.Error(err))
		return
	}
	defer tx.Rollback()

	// var image string
	// if item.Image != nil {
	// 	image = item.Image.URL
	// }

	// log.Debug("item",
	// 	zap.String("guid", item.GUID),
	// 	zap.String("published", item.Published),
	// 	zap.String("title", item.Title),
	// 	zap.String("image", image),
	// 	zap.String("link", item.Link),
	// 	zap.String("description", item.Description),
	// )

	if err := tx.Commit(); err != nil {
		log.Error("commit error", zap.Error(err))
	}
}

func newChannelFromFeed(channelConfig *app.ImportChannelConfig, feed *gofeed.Feed) app.Channel {
	var image string
	if feed.Image != nil {
		image = feed.Image.URL
	}
	return app.Channel{
		Id:          channelConfig.Id,
		Title:       feed.Title,
		Image:       image,
		Description: feed.Description,
	}
}
