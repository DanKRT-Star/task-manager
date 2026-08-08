package integration

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/DanKRT-Star/task-manager/internal/apperror"
	"github.com/DanKRT-Star/task-manager/internal/handler"
	"github.com/DanKRT-Star/task-manager/internal/repository"
	v1 "github.com/DanKRT-Star/task-manager/internal/route/v1"
	"github.com/DanKRT-Star/task-manager/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_RateLimit_AuthLogin(t *testing.T) {
	cleanTables()

	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	authService := service.NewAuthService(userRepo)
	taskService := service.NewTaskService(taskRepo)
	authHandler := handler.NewAuthHandler(authService)
	taskHandler := handler.NewTaskHandler(taskService)

	rateLimitedApp := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			if appErr, ok := err.(*apperror.AppError); ok {
				return c.Status(appErr.Code).JSON(fiber.Map{"error": appErr.Message})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		},
	})
	// Bật rate limit thật cho app riêng này (true), không đụng tới app toàn cục
	v1.SetupRoutes(rateLimitedApp, authHandler, taskHandler, true)

	loginBody := map[string]string{
		"email":    "ratelimit@example.com",
		"password": "wrongpassword",
	}
	jsonBody, _ := json.Marshal(loginBody)

	var lastStatus int
	for i := 0; i < 12; i++ {
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := rateLimitedApp.Test(req)
		require.NoError(t, err)
		lastStatus = resp.StatusCode
	}

	assert.Equal(t, 429, lastStatus)
}