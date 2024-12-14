package usecase

import (
	"context"
	"fmt"
	"time"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"authentication-system/internal/application/dto"
	"authentication-system/internal/application/service"
	"authentication-system/internal/domain/entity"
	"authentication-system/internal/domain/repository"
)

type RegisterUserUseCase struct {
	userRepo         repository.UserRepository
	verificationRepo repository.VerificationRepository
	emailService     service.EmailService
}

func NewRegisterUserUseCase(
	userRepo repository.UserRepository,
	verificationRepo repository.VerificationRepository,
	emailService service.EmailService,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo:         userRepo,
		verificationRepo: verificationRepo,
		emailService:     emailService,
	}
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, req dto.RegisterUserRequest) (*dto.RegisterUserResponse, error) {
	if _, err := uc.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return nil, fmt.Errorf("email already registered")
	}

	if _, err := uc.userRepo.GetByUsername(ctx, req.Username); err == nil {
		return nil, fmt.Errorf("username already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &entity.User{
		ID:            uuid.New(),
		Username:      req.Username,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		EmailVerified: false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	verificationCode := &entity.VerificationCode{
		ID:        uuid.New(),
		UserID:    user.ID,
		Code:      generateVerificationCode(),
		Type:      "email_verification",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := uc.verificationRepo.Create(ctx, verificationCode); err != nil {
		return nil, fmt.Errorf("failed to create verification code: %w", err)
	}

	if err := uc.emailService.SendVerificationEmail(user.Email, verificationCode.Code); err != nil {
		return nil, fmt.Errorf("failed to send verification email: %w", err)
	}

	return &dto.RegisterUserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
	}, nil
}

func generateVerificationCode() string {
	return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
}
