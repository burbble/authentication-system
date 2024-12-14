package handler

import (
	"encoding/json"
	"net/http"
	"github.com/go-playground/validator/v10"
	"authentication-system/internal/application/dto"
	"authentication-system/internal/application/usecase"
)

type UserHandler struct {
	registerUseCase    *usecase.RegisterUserUseCase
	verifyEmailUseCase *usecase.VerifyEmailUseCase
	validator         *validator.Validate
}

func NewUserHandler(
	registerUseCase *usecase.RegisterUserUseCase,
	verifyEmailUseCase *usecase.VerifyEmailUseCase,
) *UserHandler {
	return &UserHandler{
		registerUseCase:    registerUseCase,
		verifyEmailUseCase: verifyEmailUseCase,
		validator:         validator.New(),
	}
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dto.ErrorResponse{
		Message: message,
	})
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.registerUseCase.Execute(r.Context(), req)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.verifyEmailUseCase.Execute(r.Context(), req); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
