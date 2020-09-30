package actions

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/ck33122/news-aggregator/app"

	"github.com/go-pg/pg/v10"
)

type ActionError struct {
	message  string
	notFound bool
}

func (err *ActionError) Error() string {
	return err.message
}

func (err *ActionError) IsNotFound() bool {
	return err.notFound
}

func wrapDbError(action string, err error) *ActionError {
	if err.Error() == pg.ErrNoRows.Error() {
		return &ActionError{
			message:  fmt.Sprintf("%s: not found", action),
			notFound: true,
		}
	}
	message := fmt.Sprintf("%s: unknown error", action)
	app.GetLog().Error(message, zap.Error(err), zap.String("action", action))
	return &ActionError{
		message:  message,
		notFound: false,
	}
}
