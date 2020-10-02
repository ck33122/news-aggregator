package migrations

import (
	"log"

	"github.com/go-pg/migrations/v8"
)

func sequentialSqlStatements(stmts ...string) func(migrations.DB) error {
	return func(db migrations.DB) error {
		for _, stmt := range stmts {
			log.Printf("running sql statement: %s", stmt)
			_, err := db.Exec(stmt)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
