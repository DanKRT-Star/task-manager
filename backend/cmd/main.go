package main

import (
	"log"
	"os"
	"time"

	apperror "github.com/DanKRT-Star/task-manager/internal/apperror"
	"github.com/DanKRT-Star/task-manager/internal/config"
	"github.com/DanKRT-Star/task-manager/internal/handler"
	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/DanKRT-Star/task-manager/internal/repository"
	v1 "github.com/DanKRT-Star/task-manager/internal/route/v1"
	"github.com/DanKRT-Star/task-manager/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/joho/godotenv"

	_ "github.com/DanKRT-Star/task-manager/docs"
	"github.com/gofiber/contrib/v3/swaggo"
)

// @title           Task Manager API
// @version         1.0
// @description     REST API for a task management application with JWT authentication.
// @host            localhost:3000
// @BasePath        /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Type "Bearer" followed by a space and the JWT token.

func main() {
	// Load biến môi trường từ file .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Kết nối database
	config.ConnectDatabase()

	// Tự động tạo bảng theo model
	config.DB.AutoMigrate(&model.User{})
	config.DB.AutoMigrate(&model.Project{})
	config.DB.AutoMigrate(&model.Task{})
	config.DB.AutoMigrate(&model.ProjectMember{})

	// Wiring: repository -> service -> handler
	userRepo := repository.NewUserRepository(config.DB)
	taskRepo := repository.NewTaskRepository(config.DB)
	projectRepo := repository.NewProjectRepository(config.DB)
	memberRepo := repository.NewProjectMemberRepository(config.DB)

	authService := service.NewAuthService(userRepo)
	taskService := service.NewTaskService(taskRepo, memberRepo)
	projectService := service.NewProjectService(projectRepo, memberRepo, userRepo)

	authHandler := handler.NewAuthHandler(authService)
	taskHandler := handler.NewTaskHandler(taskService)
	projectHandler := handler.NewProjectHandler(projectService)

	// Khởi tạo Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			if appErr, ok := err.(*apperror.AppError); ok {
				return c.Status(appErr.Code).JSON(fiber.Map{"error": appErr.Message})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		},
	})

	// Cấu hình CORS để cho phép truy cập từ frontend
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:3001", "https://task-manager-phi-one-73.vercel.app"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	// Giới hạn chung cho toàn bộ API — chống spam/DDoS cơ bản
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	}))

	// Route test kiểm tra server sống
	app.Get("/health", func(c fiber.Ctx) error {
		sqlDB, err := config.DB.DB()
		if err != nil || sqlDB.Ping() != nil {
			return c.Status(503).JSON(fiber.Map{
				"status":   "error",
				"database": "unreachable",
			})
		}

		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// Route Swagger
	app.Get("/swagger/*", swaggo.New(swaggo.Config{}))

	// Đăng ký toàn bộ route API v1
	v1.SetupRoutes(app, authHandler, taskHandler, projectHandler, true)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}