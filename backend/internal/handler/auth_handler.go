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
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "Registration payload"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
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
// @Description  Returns a short-lived access token and a long-lived refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login payload"
// @Success      200 {object} dto.LoginResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := validator.Validate.Struct(req); err != nil {
		return apperror.BadRequest(validator.FormatValidationError(err))
	}

	accessToken, refreshToken, err := h.AuthService.Login(req.Email, req.Password)
	if err != nil {
		return apperror.Unauthorized(err.Error())
	}

	return c.JSON(dto.LoginResponse{
		Message:      "login successful",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Exchange a valid refresh token for a new access + refresh token pair (rotation)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RefreshTokenRequest true "Refresh token payload"
// @Success      200 {object} dto.RefreshTokenResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.Bind().Body(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := validator.Validate.Struct(req); err != nil {
		return apperror.BadRequest(validator.FormatValidationError(err))
	}

	accessToken, refreshToken, err := h.AuthService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		return apperror.Unauthorized(err.Error())
	}

	return c.JSON(dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Logout godoc
// @Summary      Log out
// @Description  Revokes the given refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LogoutRequest true "Refresh token to revoke"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c fiber.Ctx) error {
	var req dto.LogoutRequest
	if err := c.Bind().Body(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := validator.Validate.Struct(req); err != nil {
		return apperror.BadRequest(validator.FormatValidationError(err))
	}

	if err := h.AuthService.Logout(req.RefreshToken); err != nil {
		return apperror.Internal("failed to logout")
	}

	return c.JSON(fiber.Map{"message": "logged out successfully"})
}

// GetMe godoc
// @Summary      Get current user
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} model.User
// @Failure      401 {object} map[string]string
// @Router       /auth/me [get]
func (h *AuthHandler) GetMe(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	user, err := h.AuthService.GetMe(userID)
	if err != nil {
		return apperror.Internal("failed to fetch current user")
	}

	return c.JSON(user)
}