package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/quietguido/mapnu/mainservice/internal/repo"
	"go.uber.org/zap"

	userModel "github.com/quietguido/mapnu/mainservice/internal/repo/user/model"
)

type service struct {
	lg   *zap.Logger
	repo repo.UserRepository
}

func InitService(lg *zap.Logger, repo repo.UserRepository) *service {
	return &service{
		lg:   lg,
		repo: repo,
	}
}

func (s *service) CreateUser(ctx context.Context, createUser userModel.CreateUser) error {
	return s.repo.CreateUser(ctx, createUser)
}

func (s *service) GetUserById(ctx context.Context, userId string) (*userModel.User, error) {
	_, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.Wrap(err, "not proper uuid")
	}

	return s.repo.GetUserById(ctx, userId)
}
