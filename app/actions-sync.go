package app

import "go.uber.org/zap"

func (a *Actions) SyncChannel(channel Channel) *ActionError {
	const actionName = "sync channel"

	tx, err := a.db.Begin()
	if err != nil {
		return a.wrapDbError(actionName, err)
	}
	defer tx.Rollback()

	dbChannel := Channel{Id: channel.Id}
	if err = tx.Model(&dbChannel).WherePK().Select(); err != nil {
		a.log.Debug("some error getting channel, trying to add", zap.Error(err))
		if _, err := tx.Model(&channel).Insert(); err != nil {
			return a.wrapDbError(actionName, err)
		}
		a.log.Debug("add channel ok")
	} else if dbChannel.Title != channel.Title ||
		dbChannel.Description != channel.Description ||
		dbChannel.Image != channel.Image {
		a.log.Debug("updating channel info")
		if _, err = tx.Model(&channel).WherePK().Update(); err != nil {
			return a.wrapDbError(actionName, err)
		}
		a.log.Debug("update ok")
	} else {
		a.log.Debug("channel info is up to date")
	}

	if err := tx.Commit(); err != nil {
		return a.wrapDbError(actionName, err)
	}
	return nil
}

func (a *Actions) SyncPost(channel *Channel, post *Post, guid string) *ActionError {
	const actionName = "sync post"

	tx, err := a.db.Begin()
	if err != nil {
		return a.wrapDbError(actionName, err)
	}
	defer tx.Rollback()

	hasRssPostId := true
	rssPostId := RssPostId{Guid: guid, ChannelId: channel.Id, PostId: post.Id}
	if err = tx.Model(&rssPostId).Where("guid = ?guid and channel_id = ?channel_id").Select(); err != nil {
		a.log.Debug("select rrs post id failed, insertion will be called later", zap.Error(err))
		hasRssPostId = false
	} else {
		post.Id = rssPostId.PostId
	}

	dbPost := Post{Id: post.Id}
	if err = tx.Model(&dbPost).WherePK().Select(); err != nil {
		a.log.Debug("some error getting post, trying to add", zap.Error(err))
		if _, err := tx.Model(post).Insert(); err != nil {
			return a.wrapDbError(actionName, err)
		}
		a.log.Debug("add post ok")
	} else if dbPost.Title != post.Title ||
		dbPost.Description != post.Description ||
		dbPost.Link != post.Link ||
		dbPost.Image != post.Image ||
		!dbPost.PublicationDate.Equal(post.PublicationDate) {
		a.log.Debug("updating post")
		if _, err = tx.Model(post).WherePK().Update(); err != nil {
			return a.wrapDbError(actionName, err)
		}
		a.log.Debug("update ok")
	} else {
		a.log.Debug("post is up to date")
	}

	if !hasRssPostId {
		if _, err = tx.Model(&rssPostId).Insert(); err != nil {
			a.log.Debug("was no rrs post id, insertion failed", zap.Error(err), zap.String("guid", guid))
			return a.wrapDbError(actionName, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return a.wrapDbError(actionName, err)
	}
	return nil
}
