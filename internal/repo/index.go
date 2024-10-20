package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/quietguido/mapnu/internal/repo/event"
	"github.com/quietguido/mapnu/internal/repo/event/model"
	"go.uber.org/zap"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, createEvent model.CreateEvent) error
	GetEventById(ctx context.Context, eventId string) (*model.Event, error)
}

type Repositories struct {
	Event EventRepository
}

func InitRepositories(lg *zap.Logger, db *sqlx.DB) *Repositories {
	return &Repositories{
		Event: event.NewRepository(lg, db),
	}
}
