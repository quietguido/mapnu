package booking

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/quietguido/mapnu/mainservice/internal/repo"
	"go.uber.org/zap"

	bookingModel "github.com/quietguido/mapnu/mainservice/internal/repo/booking/model"
)

const (
	ConfimedBookingStatus = "confirmed"
	PendingBookingStatus  = "pending"
	RejectedBookingStatus = "rejected"
)

type service struct {
	lg          *zap.Logger
	bookingRepo repo.BookingReposity
	eventRepo   repo.EventRepository
}

func InitService(
	lg *zap.Logger,
	bookingRepo repo.BookingReposity,
	eventRepo repo.EventRepository,
) *service {
	return &service{
		lg:          lg,
		bookingRepo: bookingRepo,
		eventRepo:   eventRepo,
	}
}

func (s *service) Create(ctx context.Context, createBooking bookingModel.CreateBooking) (int, error) {
	return s.bookingRepo.CreateBooking(ctx, createBooking)
}

func (s *service) GetBookingById(ctx context.Context, bookingId int) (*bookingModel.Booking, error) {
	return s.bookingRepo.GetBookingById(ctx, bookingId)
}

func (s *service) GetBookingsForUser(ctx context.Context, userId uuid.UUID) ([]bookingModel.Booking, error) {
	return s.bookingRepo.GetBookingsForUser(ctx, userId)
}

func (s *service) ChangeBookingStatus(ctx context.Context, changeBookingStatus bookingModel.ChangeBookingStatus) error {
	if !checkBookingStatus(changeBookingStatus.BookingStatus) {
		return errors.New("Incorrect booking status")
	}

	booking, err := s.bookingRepo.GetBookingById(ctx, changeBookingStatus.BookingID)
	if err != nil {
		return err
	}

	event, err := s.eventRepo.GetEventById(ctx, int(booking.EventID))
	if err != nil {
		return err
	}

	if event.CreatedBy == nil || *event.CreatedBy != changeBookingStatus.EventHolderUserId {
		return errors.New("Booking does not belong to user")
	}

	err = s.bookingRepo.ChangeBookingStatus(ctx, changeBookingStatus.BookingID, changeBookingStatus.BookingStatus)
	return err
}

func (s *service) GetBookingApplicationsForOrganizer(ctx context.Context, userId uuid.UUID) ([]bookingModel.Booking, error) {
	//get events for userid
	//get bookings for some eventIDS

	return nil, nil
}

func checkBookingStatus(bookingStatus string) bool {
	switch bookingStatus {
	case ConfimedBookingStatus, PendingBookingStatus, RejectedBookingStatus:
		return true
	default:
		return false
	}
}
