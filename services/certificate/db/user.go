package db

import (
	"time"
)

// User represents the database schema of users.
type User struct {
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// UserDatabase is the interface that wraps all database operations related to
// users.
type UserDatabase interface {
	AddUser(user *User) error
	GetUser(userUUID string) (*User, error)
	DeleteUser(userUUID string) error
}
