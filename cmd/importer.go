package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/mmcdole/gofeed"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	"github.com/ck33122/news-aggregator/app"
)

var (
	log            *zap.Logger
	actions        *app.Actions
	config         *app.Config
	rssTimeFormats = []string{time.RFC1123Z, time.RFC822, time.RFC1123}
)

func main() {
	flags := app.NewFlags()
	if flags.ShowHelp {
		flag.Usage()
		return
	}
	var err error
	config, err = app.NewConfig(flags.ConfigPath, "news-aggregator-importer")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating config object: %v", err)
		os.Exit(1)
	}
	log, err = app.NewLog(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating log object: %v", err)
		os.Exit(1)
	}
	db, err := app.NewDB(log, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating database object: %v", err)
		os.Exit(1)
	}
	actions, err = app.NewActions(log, db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating actions object: %v", err)
		os.Exit(1)
	}

	scheduler := gocron.NewScheduler(time.UTC)
	_, err = scheduler.Every(10).Minutes().StartImmediately().Do(updateRssTask)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating job object: %v", err)
		os.Exit(1)
	}

	scheduler.StartBlocking()
}

func updateRssTask() {
	log.Info("running update rss job")

	for _, im := range config.Import {
		log.Info(
			"importing rss",
			zap.String("address", im.Address),
			zap.String("id", im.Id.String()),
		)
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		feedParser := gofeed.NewParser()
		feed, err := feedParser.ParseURLWithContext(im.Address, ctx)
		if err != nil {
			log.Error("error parsing rss", zap.Error(err))
			continue
		}

		channel := newChannelFromFeed(&im, feed)
		if err := actions.SyncChannel(channel); err != nil {
			log.Error("error sync channel", zap.Error(err))
			continue
		}

		for _, item := range feed.Items {
			post, err := newPostFromItem(&channel, item)
			if err != nil {
				log.Error("error create post from rss item", zap.Error(err))
				continue
			}
			guid := item.GUID
			if len(guid) == 0 {
				guid = item.Title
			}
			if len(guid) == 0 {
				log.Error("error create post from rss item: there is not guid or title tags")
				continue
			}
			if err := actions.SyncPost(&channel, post, guid); err != nil {
				log.Error("error sync post", zap.Error(err))
				continue
			}
		}
	}

	log.Info("update rss job done")
}

func newPostFromItem(channel *app.Channel, item *gofeed.Item) (*app.Post, error) {
	var image string
	if item.Image != nil {
		image = item.Image.URL
	}
	if len(image) == 0 {
		for _, enclosure := range item.Enclosures {
			if len(enclosure.URL) != 0 && strings.HasPrefix(enclosure.Type, "image/") {
				image = enclosure.URL
				break
			}
		}
	}

	published, parsed := parseRssTime(item.Published)
	if !parsed {
		return nil, fmt.Errorf("can't parse date %v", item.Published)
	}

	return &app.Post{
		Id:              uuid.NewV4(),
		ChannelId:       channel.Id,
		PublicationDate: *published,
		Title:           strings.TrimSpace(item.Title),
		Image:           image,
		Link:            item.Link,
		Description:     item.Description,
	}, nil
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

func parseRssTime(timeStr string) (*time.Time, bool) {
	for _, fmt := range rssTimeFormats {
		if result, err := time.Parse(fmt, timeStr); err == nil {
			return &result, true
		}
	}
	return nil, false
}
