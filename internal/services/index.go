package services

import (
	"context"

	"github.com/quietguido/mapnu/internal/repo"
	eventModel "github.com/quietguido/mapnu/internal/repo/event/model"
	"github.com/quietguido/mapnu/internal/services/event"
	"go.uber.org/zap"
)

type EventService interface {
	Create(ctx context.Context, createEvent eventModel.CreateEvent) error
	GetEventById(ctx context.Context, eventId string) (*eventModel.Event, error)
	GetMapForQuadrant(ctx context.Context, mapQuery eventModel.GetMapQueryParams) ([]eventModel.Event, error)
}

type Service struct {
	Event EventService
}

func InitServices(lg *zap.Logger, repos *repo.Repositories) *Service {
	return &Service{
		Event: event.InitService(lg, repos.Event),
	}
}
