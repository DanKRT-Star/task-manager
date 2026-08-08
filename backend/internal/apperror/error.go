package apperror

import "github.com/gofiber/fiber/v3"

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func BadRequest(message string) *AppError {
	return New(fiber.StatusBadRequest, message)
}

func Unauthorized(message string) *AppError {
	return New(fiber.StatusUnauthorized, message)
}

func NotFound(message string) *AppError {
	return New(fiber.StatusNotFound, message)
}

func Internal(message string) *AppError {
	return New(fiber.StatusInternalServerError, message)
}