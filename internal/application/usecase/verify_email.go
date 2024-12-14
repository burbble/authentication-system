package usecase

import (
	"context"
	"fmt"
	"time"
	"authentication-system/internal/application/dto"
	"authentication-system/internal/domain/repository"
)

type VerifyEmailUseCase struct {
	userRepo         repository.UserRepository
	verificationRepo repository.VerificationRepository
}

func NewVerifyEmailUseCase(
	userRepo repository.UserRepository,
	verificationRepo repository.VerificationRepository,
) *VerifyEmailUseCase {
	return &VerifyEmailUseCase{
		userRepo:         userRepo,
		verificationRepo: verificationRepo,
	}
}

func (uc *VerifyEmailUseCase) Execute(ctx context.Context, req dto.VerifyEmailRequest) error {
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if user.EmailVerified {
		return fmt.Errorf("email already verified")
	}

	verificationCode, err := uc.verificationRepo.GetByCode(ctx, req.Code)
	if err != nil {
		return fmt.Errorf("invalid verification code")
	}

	if verificationCode.UserID != user.ID {
		return fmt.Errorf("invalid verification code")
	}

	if verificationCode.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("verification code expired")
	}

	if verificationCode.UsedAt != nil {
		return fmt.Errorf("verification code already used")
	}

	if err := uc.userRepo.UpdateEmailVerification(ctx, user.ID, true); err != nil {
		return fmt.Errorf("failed to update email verification status: %w", err)
	}

	if err := uc.verificationRepo.MarkAsUsed(ctx, verificationCode.ID); err != nil {
		return fmt.Errorf("failed to mark verification code as used: %w", err)
	}

	return nil
}
