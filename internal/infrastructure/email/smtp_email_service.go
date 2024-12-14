package email

import (
	"fmt"
	"net/smtp"
	"authentication-system/internal/application/service"
)

type SMTPConfig struct {
    Host     string
    Port     string
    Username string
    Password string
    From     string
}

type SMTPEmailService struct {
    config SMTPConfig
}

func NewSMTPEmailService(config SMTPConfig) service.EmailService {
    return &SMTPEmailService{
        config: config,
    }
}

func (s *SMTPEmailService) SendVerificationEmail(to string, code string) error {
    auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
    
    subject := "Email Verification"
    body := fmt.Sprintf(`
        <html>
            <body>
                <h2>Email Verification</h2>
                <p>Your verification code is: <strong>%s</strong></p>
                <p>This code will expire in 24 hours.</p>
            </body>
        </html>
    `, code)

    msg := []byte(fmt.Sprintf("To: %s\r\n"+
        "From: %s\r\n"+
        "MIME-Version: 1.0\r\n"+
        "Content-Type: text/html; charset=UTF-8\r\n"+
        "Subject: %s\r\n\r\n"+
        "%s", to, s.config.From, subject, body))

    addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
    if err := smtp.SendMail(addr, auth, s.config.From, []string{to}, msg); err != nil {
        return fmt.Errorf("failed to send email: %w", err)
    }

    return nil
}
