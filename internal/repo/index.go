package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/quietguido/mapnu/internal/repo/event"
	eventModel "github.com/quietguido/mapnu/internal/repo/event/model"
	"go.uber.org/zap"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, createEvent eventModel.CreateEvent) error
	GetEventById(ctx context.Context, eventId string) (*eventModel.Event, error)
	GetMapForQuadrant(ctx context.Context, mapQuery eventModel.GetMapQueryParams) ([]eventModel.Event, error)
}

type Repositories struct {
	Event EventRepository
}

func InitRepositories(lg *zap.Logger, db *sqlx.DB) *Repositories {
	return &Repositories{
		Event: event.NewRepository(lg, db),
	}
}
