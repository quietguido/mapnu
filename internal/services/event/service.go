package event

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/quietguido/mapnu/internal/repo"
	"go.uber.org/zap"

	eventModel "github.com/quietguido/mapnu/internal/repo/event/model"
)

type service struct {
	lg   *zap.Logger
	repo repo.EventRepository
}

func InitService(lg *zap.Logger, repo repo.EventRepository) *service {
	return &service{
		lg:   lg,
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, createEvent eventModel.CreateEvent) error {
	return s.repo.CreateEvent(ctx, createEvent)
}

func (s *service) GetEventById(ctx context.Context, eventId string) (*eventModel.Event, error) {
	_, err := uuid.Parse(eventId)
	if err != nil {
		return nil, errors.Wrap(err, "not proper uuid")
	}

	return s.repo.GetEventById(ctx, eventId)
}

func (s *service) GetMapForQuadrant(ctx context.Context, mapQuery eventModel.GetMapQueryParams) ([]eventModel.Event, error) {
	return s.repo.GetMapForQuadrant(ctx, mapQuery)
}
