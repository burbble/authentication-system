package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"authentication-system/internal/domain/entity"
)

type VerificationRepository struct {
	db *Postgres
}

func NewVerificationRepository(db *Postgres) *VerificationRepository {
	return &VerificationRepository{
		db: db,
	}
}

func (r *VerificationRepository) Create(ctx context.Context, code *entity.VerificationCode) error {
	query := `
		INSERT INTO verification_codes (id, user_id, code, type, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	
	_, err := r.db.GetDB().Exec(ctx, query,
		code.ID,
		code.UserID,
		code.Code,
		code.Type,
		code.ExpiresAt,
		code.CreatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create verification code: %w", err)
	}
	
	return nil
}

func (r *VerificationRepository) GetByCode(ctx context.Context, code string) (*entity.VerificationCode, error) {
	query := `
		SELECT id, user_id, code, type, expires_at, created_at, used_at
		FROM verification_codes
		WHERE code = $1 AND used_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`
	
	vc := &entity.VerificationCode{}
	err := r.db.GetDB().QueryRow(ctx, query, code).Scan(
		&vc.ID,
		&vc.UserID,
		&vc.Code,
		&vc.Type,
		&vc.ExpiresAt,
		&vc.CreatedAt,
		&vc.UsedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("verification code not found")
		}
		return nil, fmt.Errorf("failed to get verification code: %w", err)
	}
	
	return vc, nil
}

func (r *VerificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.VerificationCode, error) {
	query := `
		SELECT id, user_id, code, type, expires_at, created_at, used_at
		FROM verification_codes
		WHERE user_id = $1 AND used_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`
	
	vc := &entity.VerificationCode{}
	err := r.db.GetDB().QueryRow(ctx, query, userID).Scan(
		&vc.ID,
		&vc.UserID,
		&vc.Code,
		&vc.Type,
		&vc.ExpiresAt,
		&vc.CreatedAt,
		&vc.UsedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("verification code not found")
		}
		return nil, fmt.Errorf("failed to get verification code: %w", err)
	}
	
	return vc, nil
}

func (r *VerificationRepository) MarkAsUsed(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE verification_codes
		SET used_at = $1
		WHERE id = $2 AND used_at IS NULL
	`
	
	result, err := r.db.GetDB().Exec(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to mark verification code as used: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("verification code not found or already used")
	}
	
	return nil
}
