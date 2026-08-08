package handler

import (
	"github.com/DanKRT-Star/task-manager/internal/apperror"
	"github.com/DanKRT-Star/task-manager/internal/dto"
	"github.com/DanKRT-Star/task-manager/internal/service"
	"github.com/DanKRT-Star/task-manager/internal/validator"
	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}

	if err := validator.Validate.Struct(req); err != nil {
		return apperror.BadRequest(validator.FormatValidationError(err))
	}

	user, err := h.AuthService.Register(req.UserName, req.Email, req.Password)
	if err != nil {
		return apperror.BadRequest(err.Error())
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "registered successfully",
		"user":    user,
	})
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}

	if err := validator.Validate.Struct(req); err != nil {
		return apperror.BadRequest(validator.FormatValidationError(err))
	}

	token, err := h.AuthService.Login(req.Email, req.Password)
	if err != nil {
		return apperror.Unauthorized(err.Error())
	}

	return c.JSON(dto.LoginResponse{
		Message: "login successful",
		Token:   token,
	})
}

func (h *AuthHandler) GetMe(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok || userID == 0 {
		return apperror.Unauthorized("unauthorized")
	}

	user, err := h.AuthService.GetMe(userID)
	if err != nil {
		return apperror.NotFound("user not found")
	}

	return c.JSON(user)
}
