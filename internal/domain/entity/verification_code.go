package entity

import (
	"time"
	"github.com/google/uuid"
)

type VerificationCode struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Code      string
	Type      string
	ExpiresAt time.Time
	CreatedAt time.Time
	UsedAt    *time.Time
}
