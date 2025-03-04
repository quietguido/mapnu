package services

import (
	"context"
	"github.com/quietguido/mapnu/mainservice/internal/services/oauth"

	"github.com/google/uuid"
	"github.com/quietguido/mapnu/mainservice/internal/repo"
	bookingModel "github.com/quietguido/mapnu/mainservice/internal/repo/booking/model"
	eventModel "github.com/quietguido/mapnu/mainservice/internal/repo/event/model"
	userModel "github.com/quietguido/mapnu/mainservice/internal/repo/user/model"
	"github.com/quietguido/mapnu/mainservice/internal/services/booking"
	"github.com/quietguido/mapnu/mainservice/internal/services/event"
	"github.com/quietguido/mapnu/mainservice/internal/services/user"
	"go.uber.org/zap"
)

type EventService interface {
	Create(ctx context.Context, createEvent eventModel.CreateEvent) (int, error)
	GetEventById(ctx context.Context, eventId int) (*eventModel.Event, error)
	GetMapForQuadrant(ctx context.Context, mapQuery eventModel.GetMapQueryParams) ([]eventModel.Event, error)
}

type UserService interface {
	CreateUser(ctx context.Context, newUser userModel.CreateUser) error
	GetUserById(ctx context.Context, userId string) (*userModel.User, error)
}

type BookingService interface {
	Create(ctx context.Context, createBooking bookingModel.CreateBooking) (int, error)
	GetBookingById(ctx context.Context, bookingId int) (*bookingModel.Booking, error)
	GetBookingsForUser(ctx context.Context, userID uuid.UUID) ([]bookingModel.Booking, error)
	ChangeBookingStatus(ctx context.Context, changeBookingStatus bookingModel.ChangeBookingStatus) error
	GetBookingApplicationsForOrganizer(ctx context.Context, userId uuid.UUID) ([]bookingModel.Booking, error)
}

type OAuthService interface {
	VerifyIDToken(ctx context.Context, idToken string) (*oauth.Claims, error)
	GenerateJWT(email, givenName, familyName string) (string, error)
	VerifyJWT(ctx context.Context, tokenString string) (*oauth.Claims, error)
}

type Service struct {
	Event   EventService
	User    UserService
	Booking BookingService
	OAuth   OAuthService
}

func InitServices(lg *zap.Logger, repos *repo.Repositories) *Service {
	return &Service{
		Event: event.InitService(lg, repos.Event),
		User:  user.InitService(lg, repos.User),
		Booking: booking.InitService(
			lg,
			repos.Booking,
			repos.Event,
		),
		OAuth: oauth.NewOAuthService(lg),
	}
}
