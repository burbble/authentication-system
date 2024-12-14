package handler

import (
	"net/http"
	"authentication-system/internal/interfaces/http/middleware"
	"authentication-system/internal/domain/repository"
)

type Handler struct {
	userHandler *UserHandler
	userRepo    repository.UserRepository
}

func NewHandler(userHandler *UserHandler, userRepo repository.UserRepository) *Handler {
	return &Handler{
		userHandler: userHandler,
		userRepo:    userRepo,
	}
}

func (h *Handler) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	publicMux := http.NewServeMux()
	publicMux.HandleFunc("/api/v1/register", h.userHandler.Register)
	publicMux.HandleFunc("/api/v1/verify-email", h.userHandler.VerifyEmail)

	mux.Handle("/api/v1/register", publicMux)
	mux.Handle("/api/v1/verify-email", publicMux)

	handler := middleware.Logger(mux)
	handler = middleware.Recover(handler)
	handler = middleware.CORS(handler)

	return handler
}
