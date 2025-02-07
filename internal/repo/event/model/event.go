package model

import (
	"time"

	"github.com/google/uuid"
	// For handling PostGIS geometry types
)

// User represents the user table in the database.
type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Event struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	CreatedBy    uuid.UUID `json:"created_by" db:"created_by"`
	Location_lat float64   `json:"location_lat" db:"location_lat"` // For PostGIS geometry data
	Location_lon float64   `json:"location_lon" db:"location_lon"` // For PostGIS geometry data
	Time         float64   `json:"time" db:"time"`                 // For PostGIS geometry data
	Organizer    string    `json:"organizer" db:"organizer"`
	Upvote       int       `json:"upvote" db:"upvote"`
	Downvote     int       `json:"downvote" db:"downvote"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type CreateEvent struct {
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	CreatedBy    uuid.UUID `json:"created_by" db:"created_by"`
	Location_lat float64   `json:"location_lat" db:"location_lat"` // For PostGIS geometry data
	Location_lon float64   `json:"location_lon" db:"location_lon"` // For PostGIS geometry data
	Time         time.Time `json:"time" db:"time"`                 // For PostGIS geometry data
	Organizer    string    `json:"organizer" db:"organizer"`
	Upvote       int       `json:"upvote" db:"upvote"`
	Downvote     int       `json:"downvote" db:"downvote"`
	// CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
