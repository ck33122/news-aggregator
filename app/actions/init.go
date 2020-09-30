package actions

import (
	"github.com/ck33122/news-aggregator/app"
	"github.com/go-pg/pg/v10"
)

var (
	getAllChannelsStatement *pg.Stmt
	getChannelByIdStatement *pg.Stmt
)

func Init() {
	getAllChannelsStatement = mustPrepareDbStatement(`select * from channels`)
	getChannelByIdStatement = mustPrepareDbStatement(`select * from channels where id = $1`)
}

func mustPrepareDbStatement(stmt string) *pg.Stmt {
	res, err := app.GetDB().Prepare(stmt)
	app.Must(err)
	return res
}
