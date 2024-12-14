package service

type EmailService interface {
	SendVerificationEmail(to string, code string) error
}
