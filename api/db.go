package api

import (
	"github.com/ck33122/news-aggregator/app"
	"github.com/go-pg/pg/v10"
)

var (
	getAllChannelsStatement *pg.Stmt
	getChannelByIdStatement *pg.Stmt
)

func prepareStatements() error {
	db := app.GetDB()
	var err error

	getAllChannelsStatement, err = db.Prepare(`select * from channels`)
	if err != nil {
		return err
	}

	getChannelByIdStatement, err = db.Prepare(`select * from channels where id = $1`)
	if err != nil {
		return err
	}

	return err
}
