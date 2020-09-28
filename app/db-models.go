package app

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// RssPostId is model of mapping guid from RSS item to PostId
type RssPostId struct {
	Guid   string    `sql:",pk"`
	PostId uuid.UUID `sql:",type:uuid"`
	Post   *Post     `pg:"rel:has-one"`
}

// Channel is model representing RSS channel.
type Channel struct {
	Id          uuid.UUID `sql:",pk,type:uuid"`
	Title       string    `pg:",notnull"`
	Image       string
	Description string
}

// Post is model representing RSS item.
type Post struct {
	Id              uuid.UUID `sql:",pk,type:uuid"`
	ChannelId       uuid.UUID
	Channel         *Channel  `pg:"rel:has-one"`
	PublicationDate time.Time `pg:",notnull"` // +order key compose index
	Title           string    `pg:",notnull"` // +tsvector index
	Image           string
	Link            string
	Description     string `pg:",notnull"`
}
