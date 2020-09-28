package migrations

import (
	"github.com/go-pg/migrations/v8"
)

func init() {
	up := sequentialSqlStatements(
		`create table channels(
			id                uuid primary key,
			title             text not null,
			image             varchar(2048),
			description       text not null
		)`,
		`create table posts(
			id                uuid primary key,
			channel_id        uuid not null,
			publication_date  timestamp with time zone not null,
			title             text not null,
			image             varchar(2048),
			link              varchar(2048),
			description       text not null,
			constraint fk_channel
				foreign key(channel_id)
				references channels(id)
				on delete cascade
		)`,
		`create table rss_post_ids(
			guid              text primary key,
			post_id           uuid not null,
			constraint fk_post
				foreign key(post_id)
				references posts(id)
				on delete cascade
		)`,
	)
	down := sequentialSqlStatements(
		`drop table rss_post_ids`,
		`drop table posts`,
		`drop table channels`,
	)
	migrations.MustRegisterTx(up, down)
}
