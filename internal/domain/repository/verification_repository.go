package repository

import (
	"context"
	"authentication-system/internal/domain/entity"
	"github.com/google/uuid"
)

type VerificationRepository interface {
	Create(ctx context.Context, code *entity.VerificationCode) error
	GetByCode(ctx context.Context, code string) (*entity.VerificationCode, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.VerificationCode, error)
	MarkAsUsed(ctx context.Context, id uuid.UUID) error
}
