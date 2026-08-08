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

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account with a unique email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "Registration payload"
// @Success      201 {object} map[string]interface{} "message and created user"
// @Failure      400 {object} map[string]string "validation error or email already taken"
// @Router       /auth/register [post]
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

// Login godoc
// @Summary      Log in
// @Description  Authenticate with email and password, returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login payload"
// @Success      200 {object} dto.LoginResponse
// @Failure      400 {object} map[string]string "validation error"
// @Failure      401 {object} map[string]string "invalid credentials"
// @Failure      429 {object} map[string]string "too many requests"
// @Router       /auth/login [post]
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

// GetMe godoc
// @Summary      Get current user info
// @Description  Retrieve information about the currently authenticated user
// @Tags         auth
// @Produce      json
// @Success      200 {object} model.User
// @Failure      401 {object} map[string]string "unauthorized"
// @Failure      404 {object} map[string]string "user not found"
// @Router       /auth/me [get]
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
