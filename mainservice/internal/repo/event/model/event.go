package model

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	EventID      int64      `json:"event_id" db:"event_id"`               // BIGSERIAL Primary Key
	Name         string     `json:"name" db:"name"`                       // VARCHAR(255) NOT NULL
	Description  string     `json:"description" db:"description"`         // TEXT
	CreatedBy    *uuid.UUID `json:"created_by,omitempty" db:"created_by"` // UUID (Nullable, FK to users)
	Location_lat float64    `json:"location_lat" db:"location_lat"`       // For PostGIS geometry data
	Location_lon float64    `json:"location_lon" db:"location_lon"`       // For PostGIS geometry data
	StartDate    time.Time  `json:"start_date" db:"start_date"`           // TIMESTAMP WITH TIME ZONE NOT NULL
	Organizer    string     `json:"organizer" db:"organizer"`             // VARCHAR(255) NOT NULL
	Upvote       int        `json:"upvote" db:"upvote"`                   // INTEGER DEFAULT 0
	Downvote     int        `json:"downvote" db:"downvote"`               // INTEGER DEFAULT 0
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`           // TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
}

// ✅ CreateEvent struct (for inserting new events)
type CreateEvent struct {
	Name         string     `json:"name" db:"name"`                       // VARCHAR(255) NOT NULL
	Description  string     `json:"description" db:"description"`         // TEXT
	CreatedBy    *uuid.UUID `json:"created_by,omitempty" db:"created_by"` // UUID (Nullable, FK to users)
	Location_lat float64    `json:"location_lat" db:"location_lat"`       // For PostGIS geometry data
	Location_lon float64    `json:"location_lon" db:"location_lon"`       // For PostGIS geometry data
	StartDate    time.Time  `json:"start_date" db:"start_date"`           // TIMESTAMP WITH TIME ZONE NOT NULL
	Organizer    string     `json:"organizer" db:"organizer"`             // VARCHAR(255) NOT NULL
}
