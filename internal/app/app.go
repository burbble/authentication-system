package app

import (
	"authentication-system/internal/application/usecase"
	"authentication-system/internal/config"
	"authentication-system/internal/infrastructure/email"
	"authentication-system/internal/infrastructure/persistence/postgres"
	"authentication-system/internal/interfaces/http/handler"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"authentication-system/internal/application/service"
)

type App struct {
	cfg         *config.Config
	httpServer  *http.Server
	connections []io.Closer
}

func NewApp(cfg *config.Config) (*App, error) {
	app := &App{
		cfg:         cfg,
		connections: make([]io.Closer, 0),
	}

	db, err := initDatabase(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to init database: %w", err)
	}
	app.connections = append(app.connections, db)

	userRepo := postgres.NewUserRepository(db.(*postgres.Postgres))
	verificationRepo := postgres.NewVerificationRepository(db.(*postgres.Postgres))

	var emailService service.EmailService
	if cfg.UseMockEmail {
		emailService = email.NewMockEmailService()
	} else {
		emailService = email.NewSMTPEmailService(email.SMTPConfig{
			Host:     cfg.SMTP.Host,
			Port:     cfg.SMTP.Port,
			Username: cfg.SMTP.Username,
			Password: cfg.SMTP.Password,
			From:     cfg.SMTP.From,
		})
	}

	registerUseCase := usecase.NewRegisterUserUseCase(userRepo, verificationRepo, emailService)
	verifyEmailUseCase := usecase.NewVerifyEmailUseCase(userRepo, verificationRepo)

	userHandler := handler.NewUserHandler(registerUseCase, verifyEmailUseCase)
	handlers := handler.NewHandler(userHandler, userRepo)

	app.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:      handlers.RegisterRoutes(),
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	return app, nil
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		log.Printf("HTTP server starting on %s:%s", a.cfg.HTTP.Host, a.cfg.HTTP.Port)
		if err := a.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("Failed to serve HTTP: %v", err)
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case v := <-quit:
		log.Printf("Received signal: %v", v)
	case <-ctx.Done():
		log.Printf("Context cancelled: %v", ctx.Err())
	}

	return a.Shutdown()
}

func (a *App) Shutdown() error {
	if err := a.httpServer.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	for _, conn := range a.connections {
		if err := conn.Close(); err != nil {
			log.Printf("Failed to close connection: %v", err)
		}
	}

	return nil
}

func initDatabase(cfg config.DatabaseConfig) (io.Closer, error) {
	return postgres.NewPostgres(
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)
}
