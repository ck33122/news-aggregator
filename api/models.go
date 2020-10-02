package api

import (
	"time"

	"github.com/ck33122/news-aggregator/app"
)

type PostM struct {
	Id          string `json:"id"`
	ChannelId   string `json:"channel_id"`
	Published   string `json:"published"`
	Title       string `json:"title"`
	Image       string `json:"image,omitempty"`
	Link        string `json:"link"`
	Description string `json:"description"`
}

type ChannelM struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Image       string `json:"image,omitempty"`
	Description string `json:"description"`
}

func mapToPostM(post *app.Post) PostM {
	return PostM{
		Id:          post.Id.String(),
		ChannelId:   post.ChannelId.String(),
		Published:   post.PublicationDate.Format(time.RFC3339),
		Title:       post.Title,
		Image:       post.Image,
		Link:        post.Link,
		Description: post.Description,
	}
}

func mapToPostsM(posts []app.Post) []PostM {
	res := make([]PostM, len(posts))
	for idx, post := range posts {
		res[idx] = mapToPostM(&post)
	}
	return res
}

func mapToChannelM(channel *app.Channel) ChannelM {
	return ChannelM{
		Id:          channel.Id.String(),
		Title:       channel.Title,
		Image:       channel.Image,
		Description: channel.Description,
	}
}

func mapToChannelsM(channels []app.Channel) []ChannelM {
	res := make([]ChannelM, len(channels))
	for idx, channel := range channels {
		res[idx] = mapToChannelM(&channel)
	}
	return res
}
