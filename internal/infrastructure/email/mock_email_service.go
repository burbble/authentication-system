package email

import (
    "log"
    "authentication-system/internal/application/service"
)

type MockEmailService struct{}

func NewMockEmailService() service.EmailService {
    return &MockEmailService{}
}

func (s *MockEmailService) SendVerificationEmail(to string, code string) error {
    log.Printf("Mock: Sending verification email to %s with code: %s", to, code)
    return nil
}
