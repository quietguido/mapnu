package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/quietguido/mapnu/internal/repo/booking"
	bookingModel "github.com/quietguido/mapnu/internal/repo/booking/model"
	"github.com/quietguido/mapnu/internal/repo/event"
	eventModel "github.com/quietguido/mapnu/internal/repo/event/model"
	"github.com/quietguido/mapnu/internal/repo/user"
	userModel "github.com/quietguido/mapnu/internal/repo/user/model"

	"go.uber.org/zap"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, createEvent eventModel.CreateEvent) (int, error)
	GetEventById(ctx context.Context, eventId int) (*eventModel.Event, error)
	GetMapForQuadrant(ctx context.Context, mapQuery eventModel.GetMapQueryParams) ([]eventModel.Event, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, newUser userModel.CreateUser) error
	GetUserById(ctx context.Context, userId string) (*userModel.User, error)
}

type BookingReposity interface {
	CreateBooking(ctx context.Context, createBooking bookingModel.CreateBooking) (int, error)
	GetBookingById(ctx context.Context, bookingId int) (*bookingModel.Booking, error)
	GetBookingsForUser(ctx context.Context, userID uuid.UUID) ([]bookingModel.Booking, error)
	ChangeBookingStatus(ctx context.Context, bookingId int, status string) error
}

type Repositories struct {
	Event   EventRepository
	User    UserRepository
	Booking BookingReposity
}

func InitRepositories(lg *zap.Logger, db *sqlx.DB) *Repositories {
	return &Repositories{
		Event:   event.NewRepository(lg, db),
		User:    user.NewRepository(lg, db),
		Booking: booking.NewRepository(lg, db),
	}
}
