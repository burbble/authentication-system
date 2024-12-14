package errors

type AppError struct {
    Message string `json:"message"`
}

func (e *AppError) Error() string {
    return e.Message
}

func NewEmailAlreadyRegisteredError() error {
    return &AppError{Message: "email already registered"}
}

func NewUsernameAlreadyTakenError() error {
    return &AppError{Message: "username already taken"}
}
