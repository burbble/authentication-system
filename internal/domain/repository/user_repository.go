package repository

import (
	"context"
	"authentication-system/internal/domain/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	UpdateEmailVerification(ctx context.Context, id uuid.UUID, verified bool) error
}
