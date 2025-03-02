package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateUser struct {
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
}

// User represents the user table in the database.
type User struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Username string    `json:"username" db:"username"`
	Email    string    `json:"email" db:"email"`
	// Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// -- Create users table
// CREATE TABLE IF NOT EXISTS users (
//     id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
//     username VARCHAR(255) UNIQUE NOT NULL,
//     email VARCHAR(255) UNIQUE NOT NULL,
//     created_at TIMESTAMP
//     WITH
//         TIME ZONE DEFAULT CURRENT_TIMESTAMP
// );
