package actions

import (
	"github.com/ck33122/news-aggregator/app"
	"github.com/go-pg/pg/v10"
	uuid "github.com/satori/go.uuid"
)

var (
	emptyChannels = make([]app.Channel, 0)
)

func GetAllChannels() ([]app.Channel, *ActionError) {
	var channels []app.Channel
	result, err := getAllChannelsStatement.Query(&channels)
	if err != nil {
		return nil, wrapDbError("get channels", err)
	}
	if result.RowsReturned() == 0 {
		channels = emptyChannels
	}
	return channels, nil
}

func GetChannelById(id uuid.UUID) (*app.Channel, *ActionError) {
	var channel app.Channel
	if _, err := getChannelByIdStatement.QueryOne(&channel, id); err != nil {
		return nil, wrapDbError("get channel by id", err)
	}
	app.GetLog().Debug("get channel ok!")
	return &channel, nil
}

func GetChannelByIdTx(tx *pg.Tx, id uuid.UUID) (*app.Channel, *ActionError) {
	channel := app.Channel{Id: id}
	if err := tx.Model(&channel).WherePK().Select(); err != nil {
		return nil, wrapDbError("get channel by id", err)
	}
	app.GetLog().Debug("get channel ok!")
	return &channel, nil
}

func AddChannelTx(tx *pg.Tx, channel app.Channel) *ActionError {
	if _, err := tx.Model(&channel).Insert(); err != nil {
		return wrapDbError("insert channel", err)
	}
	return nil
}

func UpdateChannelTx(tx *pg.Tx, channel app.Channel) *ActionError {
	_, err := tx.Model(&channel).WherePK().Update()
	if err != nil {
		return wrapDbError("update channel", err)
	}
	return nil
}
