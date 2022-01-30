package users

import (
	"context"
	"time"

	"github.com/daniilty/sharenote-friends/internal/kafka"
	"github.com/daniilty/sharenote-friends/internal/mongo"
	events "github.com/daniilty/sharenote-kafka-events"
	"go.uber.org/zap"
)

type EventsHandler interface {
	Listen(ctx context.Context)
}

type EventsHandlerImpl struct {
	logger        *zap.SugaredLogger
	timeout       time.Duration
	db            mongo.DB
	kafkaConsumer kafka.Consumer
}

func NewEventsHandler(logger *zap.SugaredLogger, timeout time.Duration, db mongo.DB, consumer kafka.Consumer) EventsHandler {
	return &EventsHandlerImpl{
		logger:        logger,
		timeout:       timeout,
		db:            db,
		kafkaConsumer: consumer,
	}
}

func (e *EventsHandlerImpl) Listen(ctx context.Context) {
	e.logger.Info("Listening for user events.")
	for {
		select {
		case <-ctx.Done():
			e.logger.Infow("Stopping users event handler.")

			return
		default:
			e.handleMessage(ctx)
		}
	}
}

func (e *EventsHandlerImpl) handleMessage(ctx context.Context) {
	event := &events.Event{}

	commit, err := e.kafkaConsumer.UnmarshalMessage(ctx, event)
	if err != nil {
		e.logger.Errorw("Unmarshal kafka message.", "err", err)

		return
	}

	switch event.Type {
	case events.EventTypeUserDelete:
		userDeleteEvent, err := eventDataToUserDeleteEventData(event.Data)
		if err != nil {
			e.logger.Errorw("Convert event data to delete user event.", "err", err)

			err = commit(ctx)
			if err != nil {
				e.logger.Errorw("Commit kafka message.", "err", err)
			}

			return
		}

		e.logger.Infow("Remove user event.", "id", userDeleteEvent.ID)

		err = e.db.RemoveUser(ctx, userDeleteEvent.ID)
		if err != nil {
			e.logger.Errorw("Remove user.", "err", err)

			return
		}
	}

	err = commit(ctx)
	if err != nil {
		e.logger.Errorw("Commit kafka message.", "err", err)
	}
}
