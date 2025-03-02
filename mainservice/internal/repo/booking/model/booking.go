package model

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	BookingID     int64     `json:"booking_id" db:"booking_id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	EventID       int64     `json:"event_id" db:"event_id"`
	BookingStatus string    `json:"booking_status" db:"booking_status"`
	Visibility    string    `json:"visibility" db:"visibility"`
	BookedAt      time.Time `json:"booked_at" db:"booked_at"`
}

type CreateBooking struct {
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	EventID    int64     `json:"event_id" db:"event_id"`
	Visibility string    `json:"visibility" db:"visibility"` // "public", "private"
}

type ChangeBookingStatus struct {
	EventHolderUserId uuid.UUID `json:"event_holder_user_id" db:"event_holder_user_id"`
	BookingID         int       `json:"booking_id" db:"booking_id"`
	BookingStatus     string    `json:"booking_status" db:"booking_status"`
}
