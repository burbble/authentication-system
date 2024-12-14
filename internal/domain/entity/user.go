package entity

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID
	Username      string
	Name          string
	Email         string
	PasswordHash  string
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
