package migrations

import (
	"github.com/ck33122/news-aggregator/app"
	"github.com/go-pg/migrations/v8"
)

func sequentialSqlStatements(stmts ...string) func(migrations.DB) error {
	return func(db migrations.DB) error {
		log := app.GetLog().Sugar()
		for _, stmt := range stmts {
			log.Debugf("running sql statement: %s", stmt)
			_, err := db.Exec(stmt)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
