package migrations

import "github.com/go-pg/migrations/v8"

func init() {
	up := sequentialSqlStatements(
		`create index idx_posts_channel_publications on posts (channel_id, publication_date desc)`,
		`create index idx_posts_publication_date on posts (publication_date desc)`,
	)
	down := sequentialSqlStatements(
		`drop index idx_posts_channel_publications`,
		`drop index idx_posts_publication_date`,
	)
	migrations.MustRegisterTx(up, down)
}
