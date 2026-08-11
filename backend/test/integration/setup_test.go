package integration

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DanKRT-Star/task-manager/internal/apperror"
	"github.com/DanKRT-Star/task-manager/internal/config"
	"github.com/DanKRT-Star/task-manager/internal/handler"
	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/DanKRT-Star/task-manager/internal/repository"
	v1 "github.com/DanKRT-Star/task-manager/internal/route/v1"
	"github.com/DanKRT-Star/task-manager/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var app *fiber.App
var db *gorm.DB

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../.env.test")
	os.Setenv("JWT_SECRET", "test-secret-key")

	config.ConnectDatabase()
	db = config.DB

	db.Exec("DROP TABLE IF EXISTS project_members, tasks, projects, users CASCADE")
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Project{})
	db.AutoMigrate(&model.Task{})
	db.AutoMigrate(&model.ProjectMember{})

	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	memberRepo := repository.NewProjectMemberRepository(db)

	authService := service.NewAuthService(userRepo)
	taskService := service.NewTaskService(taskRepo, memberRepo)
	projectService := service.NewProjectService(projectRepo, memberRepo, userRepo)

	authHandler := handler.NewAuthHandler(authService)
	taskHandler := handler.NewTaskHandler(taskService)
	projectHandler := handler.NewProjectHandler(projectService)

	app = fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			if appErr, ok := err.(*apperror.AppError); ok {
				return c.Status(appErr.Code).JSON(fiber.Map{"error": appErr.Message})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		},
	})

	v1.SetupRoutes(app, authHandler, taskHandler, projectHandler, false)

	code := m.Run()
	os.Exit(code)
}

// cleanTables xóa sạch data giữa các test case để đảm bảo tính độc lập
func cleanTables() {
	db.Exec("TRUNCATE TABLE project_members, tasks, projects, users RESTART IDENTITY CASCADE")
}

// registerAndGetToken helper dùng chung cho các file test khác trong package
func registerAndGetToken(t *testing.T, email, password string) string {
	t.Helper()

	registerBody := map[string]string{
		"userName": "Test User",
		"email":    email,
		"password": password,
	}
	jsonBody, _ := json.Marshal(registerBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)

	loginBody := map[string]string{
		"email":    email,
		"password": password,
	}
	jsonLogin, _ := json.Marshal(loginBody)

	reqLogin := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(jsonLogin))
	reqLogin.Header.Set("Content-Type", "application/json")
	respLogin, err := app.Test(reqLogin)
	require.NoError(t, err)
	require.Equal(t, 200, respLogin.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(respLogin.Body).Decode(&result)

	return result["token"].(string)
}