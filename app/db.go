package app

import (
	"errors"

	"github.com/go-pg/pg/v10"
	"go.uber.org/zap"
)

func NewDB(log *zap.Logger, cfg *Config) (*pg.DB, error) {
	log.Info(
		"connecting to database",
		zap.String("database", cfg.Database.Database),
		zap.String("address", cfg.Database.Address),
	)
	db := pg.Connect(&pg.Options{
		ApplicationName: cfg.AppName,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		Database:        cfg.Database.Database,
		Addr:            cfg.Database.Address,
	})

	// check connection
	var n int
	if _, err := db.QueryOne(pg.Scan(&n), "select 1"); err != nil {
		return nil, err
	}
	if n != 1 {
		return nil, errors.New("error query select 1, result was not 1")
	}

	return db, nil
}
