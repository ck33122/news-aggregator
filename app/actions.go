package app

import (
	"fmt"

	"github.com/go-pg/pg/v10"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type Actions struct {
	db             *pg.DB
	log            *zap.Logger
	getAllChannels *pg.Stmt
	getAllPosts    *pg.Stmt
	getChannelById *pg.Stmt
}

func NewActions(log *zap.Logger, db *pg.DB) (*Actions, error) {
	actions := Actions{db: db, log: log}
	var err error
	actions.getAllChannels, err = db.Prepare(`select * from channels`)
	if err != nil {
		return nil, err
	}
	actions.getAllPosts, err = db.Prepare(`select * from posts`)
	if err != nil {
		return nil, err
	}
	actions.getChannelById, err = db.Prepare(`select * from channels where id = $1`)
	if err != nil {
		return nil, err
	}
	return &actions, nil
}

func (a *Actions) GetAllChannels() ([]Channel, *ActionError) {
	var channels []Channel
	result, err := a.getAllChannels.Query(&channels)
	if err != nil {
		return nil, a.wrapDbError("get all channels", err)
	}
	if result.RowsReturned() == 0 {
		channels = []Channel{}
	}
	return channels, nil
}

func (a *Actions) GetAllPosts() ([]Post, *ActionError) {
	var posts []Post
	result, err := a.getAllPosts.Query(&posts)
	if err != nil {
		return nil, a.wrapDbError("get all posts", err)
	}
	if result.RowsReturned() == 0 {
		posts = []Post{}
	}
	return posts, nil
}

func (a *Actions) GetChannelById(id uuid.UUID) (*Channel, *ActionError) {
	var channel Channel
	if _, err := a.getChannelById.QueryOne(&channel, id); err != nil {
		return nil, a.wrapDbError("get channel by id", err)
	}
	return &channel, nil
}

func (a *Actions) wrapDbError(action string, err error) *ActionError {
	if err.Error() == pg.ErrNoRows.Error() {
		return &ActionError{
			message:  fmt.Sprintf("%s: not found", action),
			notFound: true,
		}
	}
	message := fmt.Sprintf("%s: unknown error", action)
	a.log.Error(message, zap.Error(err), zap.String("action", action))
	return &ActionError{
		message:  message,
		notFound: false,
	}
}
