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
	projectRepo := repository.NewProjectRepository(db)
	memberRepo := repository.NewProjectMemberRepository(db)
	activityRepo := repository.NewActivityLogRepository(db)
	epicRepo := repository.NewEpicRepository(db)
	milestoneRepo := repository.NewMilestoneRepository(db)
	sprintRepo := repository.NewSprintRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	labelRepo := repository.NewLabelRepository(db)

	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	authService := service.NewAuthService(userRepo, refreshTokenRepo)
	taskService := service.NewTaskService(taskRepo, memberRepo, activityRepo)
	projectService := service.NewProjectService(projectRepo, memberRepo, userRepo)
	epicService := service.NewEpicService(epicRepo, memberRepo)
	milestoneService := service.NewMilestoneService(milestoneRepo, memberRepo)
	sprintService := service.NewSprintService(sprintRepo, memberRepo)
	commentService := service.NewCommentService(commentRepo, taskRepo, memberRepo)
	labelService := service.NewLabelService(labelRepo, taskRepo, memberRepo)
	activityLogService := service.NewActivityLogService(activityRepo, taskRepo, memberRepo)

	authHandler := handler.NewAuthHandler(authService)
	taskHandler := handler.NewTaskHandler(taskService)
	projectHandler := handler.NewProjectHandler(projectService)
	epicHandler := handler.NewEpicHandler(epicService)
	milestoneHandler := handler.NewMilestoneHandler(milestoneService)
	sprintHandler := handler.NewSprintHandler(sprintService)
	commentHandler := handler.NewCommentHandler(commentService)
	labelHandler := handler.NewLabelHandler(labelService)
	activityLogHandler := handler.NewActivityLogHandler(activityLogService)

	rateLimitedApp := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			if appErr, ok := err.(*apperror.AppError); ok {
				return c.Status(appErr.Code).JSON(fiber.Map{"error": appErr.Message})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		},
	})
	v1.SetupRoutes(rateLimitedApp, authHandler, taskHandler, projectHandler, epicHandler, milestoneHandler, sprintHandler, commentHandler, labelHandler, activityLogHandler, true)

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
